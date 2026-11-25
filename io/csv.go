// Package io 按 docs/DATA-FORMAT.md（D1）读写 CSV：纯列名表头、空字段表缺失、ISO 8601 时间、幂等写入。
//
// 包名与标准库 io 相同，包内不导入标准库 io；使用方可按需起别名。
package io

import (
	"encoding/csv"
	"errors"
	"fmt"
	stdio "io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Kind 是列的数据类型。
type Kind uint8

const (
	Float64 Kind = iota // 浮点，十进制或科学计数文本
	Int64               // 十进制整数
	Bool                // true / false
	String              // UTF-8 文本
	Time                // ISO 8601 日期或带时区时间戳
)

// String 返回列类型的名字。
func (k Kind) String() string {
	switch k {
	case Float64:
		return "float64"
	case Int64:
		return "int64"
	case Bool:
		return "bool"
	case String:
		return "string"
	case Time:
		return "time"
	default:
		return "unknown"
	}
}

// Column 是一列：Name / Kind 是 schema，其余切片按 Kind 承载数据，Missing 标记缺失。
type Column struct {
	Name    string
	Kind    Kind
	Floats  []float64
	Ints    []int64
	Bools   []bool
	Strings []string
	Times   []time.Time
	Missing []bool
}

// Table 是一张表，列顺序即写出顺序。
type Table struct {
	Columns []Column
}

// NewTable 用 schema 建一张空表。
func NewTable(schema ...Column) *Table {
	columns := make([]Column, len(schema))
	for i, col := range schema {
		columns[i] = Column{Name: col.Name, Kind: col.Kind}
	}
	return &Table{Columns: columns}
}

// Rows 返回行数，并校验各列数据长度一致。
func (t *Table) Rows() (int, error) {
	if t == nil || len(t.Columns) == 0 {
		return 0, errors.New("空表")
	}
	rows := -1
	for _, col := range t.Columns {
		size := len(col.Missing)
		if rows < 0 {
			rows = size
		}
		if size != rows {
			return 0, fmt.Errorf("列 %s 长度 %d，与首列 %d 不一致", col.Name, size, rows)
		}
	}
	return rows, nil
}

// ReadFile 按 schema 读取 CSV；header 顺序可以与 schema 不同，按列名对齐。
func ReadFile(path string, schema []Column) (*Table, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(raw) >= 3 && raw[0] == 0xEF && raw[1] == 0xBB && raw[2] == 0xBF {
		return nil, fmt.Errorf("%s: 文件含 UTF-8 BOM", path)
	}

	reader := csv.NewReader(strings.NewReader(string(raw)))
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("%s: 读表头: %w", path, err)
	}
	positions, err := headerPositions(header)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	table := NewTable(schema...)
	for _, col := range table.Columns {
		if _, ok := positions[col.Name]; !ok {
			return nil, fmt.Errorf("%s: 表头缺少列 %s", path, col.Name)
		}
	}
	if len(positions) != len(table.Columns) {
		return nil, fmt.Errorf("%s: 表头有 %d 列，schema 有 %d 列", path, len(positions), len(table.Columns))
	}

	line := 1
	for {
		record, err := reader.Read()
		if errors.Is(err, stdio.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%s: 第 %d 行: %w", path, line+1, err)
		}
		line++
		if len(record) != len(header) {
			return nil, fmt.Errorf("%s: 第 %d 行有 %d 个字段，应为 %d", path, line, len(record), len(header))
		}
		for i := range table.Columns {
			col := &table.Columns[i]
			text := record[positions[col.Name]]
			if err := col.appendText(text); err != nil {
				return nil, fmt.Errorf("%s: 第 %d 行 %s: %w", path, line, col.Name, err)
			}
		}
	}
	return table, nil
}

// WriteFile 按 D1 写 CSV：header 只写列名、缺失写空字段、LF 结尾。
func WriteFile(path string, table *Table) error {
	rows, err := table.Rows()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	header := make([]string, len(table.Columns))
	for i, col := range table.Columns {
		if err := validateName(col.Name); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		header[i] = col.Name
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("%s: 写表头: %w", path, err)
	}
	record := make([]string, len(table.Columns))
	for i := 0; i < rows; i++ {
		for j := range table.Columns {
			record[j] = table.Columns[j].format(i)
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("%s: 第 %d 行: %w", path, i+2, err)
		}
	}
	writer.Flush()
	return writer.Error()
}

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func validateName(name string) error {
	if len(name) > 64 || !namePattern.MatchString(name) {
		return fmt.Errorf("列名非法: %q（要求小写 snake_case，长度 <= 64）", name)
	}
	return nil
}

func headerPositions(header []string) (map[string]int, error) {
	positions := make(map[string]int, len(header))
	for i, name := range header {
		if err := validateName(name); err != nil {
			return nil, err
		}
		if _, ok := positions[name]; ok {
			return nil, fmt.Errorf("表头列名重复: %s", name)
		}
		positions[name] = i
	}
	return positions, nil
}

func (c *Column) appendText(text string) error {
	if text == "" {
		c.Missing = append(c.Missing, true)
		c.appendZero()
		return nil
	}
	c.Missing = append(c.Missing, false)
	switch c.Kind {
	case Float64:
		value, err := parseFloat(text)
		if err != nil {
			return err
		}
		c.Floats = append(c.Floats, value)
	case Int64:
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return fmt.Errorf("非法整数 %q", text)
		}
		c.Ints = append(c.Ints, value)
	case Bool:
		switch text {
		case "true":
			c.Bools = append(c.Bools, true)
		case "false":
			c.Bools = append(c.Bools, false)
		default:
			return fmt.Errorf("非法布尔值 %q", text)
		}
	case String:
		c.Strings = append(c.Strings, text)
	case Time:
		value, err := parseTime(text)
		if err != nil {
			return err
		}
		c.Times = append(c.Times, value)
	default:
		return fmt.Errorf("未知列类型 %d", c.Kind)
	}
	return nil
}

func (c *Column) appendZero() {
	switch c.Kind {
	case Float64:
		c.Floats = append(c.Floats, 0)
	case Int64:
		c.Ints = append(c.Ints, 0)
	case Bool:
		c.Bools = append(c.Bools, false)
	case String:
		c.Strings = append(c.Strings, "")
	case Time:
		c.Times = append(c.Times, time.Time{})
	}
}

func (c *Column) format(i int) string {
	if c.Missing[i] {
		return ""
	}
	switch c.Kind {
	case Float64:
		value := c.Floats[i]
		if math.IsNaN(value) {
			return ""
		}
		if math.IsInf(value, 1) {
			return "inf"
		}
		if math.IsInf(value, -1) {
			return "-inf"
		}
		return strconv.FormatFloat(value, 'g', -1, 64)
	case Int64:
		return strconv.FormatInt(c.Ints[i], 10)
	case Bool:
		if c.Bools[i] {
			return "true"
		}
		return "false"
	case String:
		return c.Strings[i]
	case Time:
		return c.Times[i].UTC().Format(time.RFC3339)
	default:
		return ""
	}
}

func parseFloat(text string) (float64, error) {
	switch text {
	case "inf":
		return math.Inf(1), nil
	case "-inf":
		return math.Inf(-1), nil
	}
	lower := strings.ToLower(text)
	if strings.Contains(lower, "nan") || strings.Contains(lower, "inf") {
		return 0, fmt.Errorf("非法浮点数 %q（缺失写空字段，无穷写 inf / -inf）", text)
	}
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("非法浮点数 %q", text)
	}
	return value, nil
}

func parseTime(text string) (time.Time, error) {
	if value, err := time.Parse("2006-01-02", text); err == nil {
		return value.UTC(), nil
	}
	if value, err := time.Parse(time.RFC3339, text); err == nil {
		return value.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("非法时间 %q（要求 ISO 8601 日期或带时区时间戳）", text)
}
