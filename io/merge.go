package io

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// MergeByKey 按主键把增量合并进已有表：已有键更新整行、新键按增量顺序追加。
//
// 主键列必须是 String / Int64 / Time；主键缺失、增量内重复键、已有表内重复键都返回错误。
func MergeByKey(existing, incoming *Table, key string) (*Table, error) {
	if existing == nil || incoming == nil {
		return nil, errors.New("合并需要两张表")
	}
	if err := sameSchema(existing, incoming); err != nil {
		return nil, err
	}
	keyIndex, err := columnIndex(existing, key)
	if err != nil {
		return nil, err
	}
	switch existing.Columns[keyIndex].Kind {
	case String, Int64, Time:
	default:
		return nil, fmt.Errorf("主键列 %s 不能是 %s", key, existing.Columns[keyIndex].Kind)
	}

	existingRows, err := existing.Rows()
	if err != nil {
		return nil, err
	}
	incomingRows, err := incoming.Rows()
	if err != nil {
		return nil, err
	}
	positions, err := keyPositions(existing, key, existingRows)
	if err != nil {
		return nil, err
	}
	if _, err := keyPositions(incoming, key, incomingRows); err != nil {
		return nil, err
	}

	merged := cloneTable(existing)
	for row := 0; row < incomingRows; row++ {
		value, err := keyText(&incoming.Columns[keyIndex], row)
		if err != nil {
			return nil, err
		}
		if position, ok := positions[value]; ok {
			replaceRow(merged, position, incoming, row)
			continue
		}
		positions[value] = len(merged.Columns[0].Missing)
		appendRow(merged, incoming, row)
	}
	return merged, nil
}

// WriteByKey 读取已有文件（不存在则直接写），按主键合并后整文件覆盖；重复执行幂等。
func WriteByKey(path string, incoming *Table, key string) error {
	if _, err := columnIndex(incoming, key); err != nil {
		return err
	}
	existing, err := ReadFile(path, schemaOf(incoming))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if _, err := incoming.Rows(); err != nil {
			return err
		}
		return WriteFile(path, incoming)
	}
	merged, err := MergeByKey(existing, incoming, key)
	if err != nil {
		return err
	}
	return WriteFile(path, merged)
}

func schemaOf(table *Table) []Column {
	schema := make([]Column, len(table.Columns))
	for i, col := range table.Columns {
		schema[i] = Column{Name: col.Name, Kind: col.Kind}
	}
	return schema
}

func sameSchema(left, right *Table) error {
	if len(left.Columns) != len(right.Columns) {
		return fmt.Errorf("列数不一致: %d vs %d", len(left.Columns), len(right.Columns))
	}
	for i := range left.Columns {
		if left.Columns[i].Name != right.Columns[i].Name || left.Columns[i].Kind != right.Columns[i].Kind {
			return fmt.Errorf("第 %d 列模式不一致: %s/%s vs %s/%s",
				i+1, left.Columns[i].Name, left.Columns[i].Kind, right.Columns[i].Name, right.Columns[i].Kind)
		}
	}
	return nil
}

func columnIndex(table *Table, name string) (int, error) {
	for i, col := range table.Columns {
		if col.Name == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("主键列不存在: %s", name)
}

func keyPositions(table *Table, key string, rows int) (map[string]int, error) {
	index, err := columnIndex(table, key)
	if err != nil {
		return nil, err
	}
	positions := make(map[string]int, rows)
	for row := 0; row < rows; row++ {
		value, err := keyText(&table.Columns[index], row)
		if err != nil {
			return nil, err
		}
		if _, ok := positions[value]; ok {
			return nil, fmt.Errorf("主键重复: %s=%s", key, value)
		}
		positions[value] = row
	}
	return positions, nil
}

func keyText(column *Column, row int) (string, error) {
	if column.Missing[row] {
		return "", fmt.Errorf("主键列 %s 第 %d 行缺失", column.Name, row+1)
	}
	switch column.Kind {
	case String:
		return column.Strings[row], nil
	case Int64:
		return strconv.FormatInt(column.Ints[row], 10), nil
	case Time:
		return column.Times[row].UTC().Format(time.RFC3339Nano), nil
	default:
		return "", fmt.Errorf("主键列 %s 类型不支持", column.Name)
	}
}

func cloneTable(table *Table) *Table {
	columns := make([]Column, len(table.Columns))
	for i, col := range table.Columns {
		columns[i] = Column{
			Name:    col.Name,
			Kind:    col.Kind,
			Floats:  append([]float64(nil), col.Floats...),
			Ints:    append([]int64(nil), col.Ints...),
			Bools:   append([]bool(nil), col.Bools...),
			Strings: append([]string(nil), col.Strings...),
			Times:   append([]time.Time(nil), col.Times...),
			Missing: append([]bool(nil), col.Missing...),
		}
	}
	return &Table{Columns: columns}
}

func appendRow(dst, src *Table, row int) {
	for i := range dst.Columns {
		appendCell(&dst.Columns[i], &src.Columns[i], row)
	}
}

func replaceRow(dst *Table, dstRow int, src *Table, srcRow int) {
	for i := range dst.Columns {
		replaceCell(&dst.Columns[i], dstRow, &src.Columns[i], srcRow)
	}
}

func appendCell(dst, src *Column, row int) {
	dst.Missing = append(dst.Missing, src.Missing[row])
	switch dst.Kind {
	case Float64:
		dst.Floats = append(dst.Floats, src.Floats[row])
	case Int64:
		dst.Ints = append(dst.Ints, src.Ints[row])
	case Bool:
		dst.Bools = append(dst.Bools, src.Bools[row])
	case String:
		dst.Strings = append(dst.Strings, src.Strings[row])
	case Time:
		dst.Times = append(dst.Times, src.Times[row])
	}
}

func replaceCell(dst *Column, dstRow int, src *Column, srcRow int) {
	dst.Missing[dstRow] = src.Missing[srcRow]
	switch dst.Kind {
	case Float64:
		dst.Floats[dstRow] = src.Floats[srcRow]
	case Int64:
		dst.Ints[dstRow] = src.Ints[srcRow]
	case Bool:
		dst.Bools[dstRow] = src.Bools[srcRow]
	case String:
		dst.Strings[dstRow] = src.Strings[srcRow]
	case Time:
		dst.Times[dstRow] = src.Times[srcRow]
	}
}
