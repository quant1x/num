package contract

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// repoRoot 是 tests/contract 相对仓库根的路径。
const repoRoot = "../.."

// nameMapColumns 是 NAME-MAP.tsv 的列数。
const nameMapColumns = 13

// 契约表里 go 列可用的包限定符及其目录。
var contractPackages = map[string]string{
	"num": ".",
	"x32": "x32",
	"x64": "x64",
}

// 同一功能的写入变体与标量变体，按后缀视为已被基础符号覆盖。
var derivedSuffixes = []string{
	"_Inplace", "_Into", "Inplace",
	"Number", "Number_Inplace", "Number_Into",
}

var allowedNan = map[string]bool{
	"—": true, "propagate": true, "comparison_false": true,
	"no_special": true, "nan_to_zero": true, "default": true,
}

// return_shape 只允许下面几种形态：标量 / 元组 / 下标 / 与输入同形 / 一维或二维定形 / 定宽数值序列。
var shapePattern = regexp.MustCompile(`^(—|scalar|tuple|index|first_len|max_len|same_as\([a-z_]+\)|bool\[(same_as\([a-z_]+\)|[a-z_]+)\]|[a-z0-9]*\[[^\[\]\s]+\](\[[^\[\]\s]+\])*)$`)

func TestNameMapFormat(t *testing.T) {
	rows := readTable(t, "NAME-MAP.tsv", nameMapColumns)
	seen := make(map[string]bool, len(rows))
	for i, row := range rows {
		line := i + 2 // 数据行起始于文件第 2 行（第 1 行是列说明注释）
		semantic := row[0]
		if seen[semantic] {
			t.Errorf("第 %d 行语义名重复: %s", line, semantic)
		}
		seen[semantic] = true
		for j, field := range row {
			if field == "" {
				t.Errorf("第 %d 行第 %d 列为空", line, j+1)
			}
		}
		if !allowedNan[row[3]] {
			t.Errorf("第 %d 行 nan_convention 取值非法: %s", line, row[3])
		}
		if !shapePattern.MatchString(row[4]) {
			t.Errorf("第 %d 行 return_shape 取值非法: %s", line, row[4])
		}
		if row[11] != "go" && row[11] != "python" {
			t.Errorf("第 %d 行 authority 只能是 go（临时基准）或 python: %s", line, row[11])
		}
		if row[11] == "python" && row[5] == "—" {
			t.Errorf("第 %d 行 authority 标了 python，但 python 列没登记符号", line)
		}
	}
}

func TestNameMapCoversGoSymbols(t *testing.T) {
	rows := readTable(t, "NAME-MAP.tsv", nameMapColumns)
	listed := make(map[string]map[string]bool, len(contractPackages))
	for qualifier := range contractPackages {
		listed[qualifier] = make(map[string]bool)
	}
	for _, row := range rows {
		cell := row[8]
		if cell == "—" {
			continue
		}
		for _, item := range strings.Split(cell, ", ") {
			qualifier, name, ok := strings.Cut(item, ":")
			if !ok {
				t.Fatalf("go 列无法解析: %q（行 %q）", item, row[0])
			}
			if _, ok := contractPackages[qualifier]; !ok {
				t.Fatalf("go 列包限定符非法: %q（行 %q）", qualifier, row[0])
			}
			listed[qualifier][name] = true
		}
	}

	excluded := readExcluded(t)

	for _, qualifier := range sortedKeys(contractPackages) {
		dir := filepath.Join(repoRoot, contractPackages[qualifier])
		actual := exportedFuncs(t, dir)
		var missing []string
		for _, name := range actual {
			if coveredByBase(name, listed[qualifier]) {
				continue
			}
			if excluded[qualifier][name] || excluded[qualifier]["*"] {
				continue
			}
			missing = append(missing, qualifier+":"+name)
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			t.Errorf("%s 有公开符号未登记契约表: %s", qualifier, strings.Join(missing, ", "))
		}

		var stale []string
		for name := range listed[qualifier] {
			if !contains(actual, name) {
				stale = append(stale, qualifier+":"+name)
			}
		}
		if len(stale) > 0 {
			sort.Strings(stale)
			t.Errorf("契约表登记的符号不存在: %s", strings.Join(stale, ", "))
		}
	}
}

// TestNameMapRustSymbols 校验 rust 列登记的符号确实出现在 src/ 的 pub fn 中。
func TestNameMapRustSymbols(t *testing.T) {
	rows := readTable(t, "NAME-MAP.tsv", nameMapColumns)
	symbols := rustPublicFunctions(t)
	for i, row := range rows {
		cell := row[7]
		if cell == "—" {
			continue
		}
		for _, name := range strings.Split(cell, ", ") {
			if !symbols[strings.TrimSpace(name)] {
				t.Errorf("第 %d 行 rust 符号 %q 在 src/ 中不存在（语义 %s）", i+2, name, row[0])
			}
		}
	}
}

// rustPublicFunctions 枚举 src/*.rs 里的顶层 pub fn 名称。
func rustPublicFunctions(t *testing.T) map[string]bool {
	t.Helper()
	pattern := regexp.MustCompile(`(?m)^\s*pub fn ([A-Za-z0-9_]+)`)
	out := make(map[string]bool)
	dir := repoRoot
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("读取 %s: %v", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".rs") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatalf("读取 %s: %v", entry.Name(), err)
		}
		for _, match := range pattern.FindAllStringSubmatch(string(data), -1) {
			out[match[1]] = true
		}
	}
	if len(out) == 0 {
		t.Fatal("src/ 未枚举到 pub fn")
	}
	return out
}

// TestNameMapPythonSymbols 校验 python 列登记的 `模块.函数` 确实在 num/ 包里。
func TestNameMapPythonSymbols(t *testing.T) {
	rows := readTable(t, "NAME-MAP.tsv", nameMapColumns)
	for i, row := range rows {
		cell := row[5]
		if cell == "—" {
			continue
		}
		for _, item := range strings.Split(cell, ", ") {
			dot := strings.LastIndex(item, ".")
			if dot <= 0 || dot == len(item)-1 {
				t.Errorf("第 %d 行 python 列无法解析: %q（应为 模块.函数，允许多级模块）", i+2, item)
				continue
			}
			module := strings.ReplaceAll(item[:dot], ".", "/")
			function := item[dot+1:]
			path := filepath.Join(repoRoot, "num", filepath.FromSlash(module)+".py")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("第 %d 行 python 符号 %q：读取 %s 失败: %v", i+2, item, path, err)
				continue
			}
			if !strings.Contains(string(data), "def "+function+"(") {
				t.Errorf("第 %d 行 python 符号 %q：%s 里没有 def %s(", i+2, item, path, function)
			}
		}
	}
}

func TestTableEncoding(t *testing.T) {
	for _, name := range []string{"NAME-MAP.tsv", "NAME-MAP-EXCLUDED.tsv"} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("读取 %s: %v", name, err)
		}
		if bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
			t.Errorf("%s 含 UTF-8 BOM", name)
		}
		if bytes.ContainsRune(data, '\r') {
			t.Errorf("%s 含 CR，行尾必须是 LF", name)
		}
		if len(data) == 0 || data[len(data)-1] != '\n' {
			t.Errorf("%s 必须以单个换行结尾", name)
		}
	}
}

// readTable 读取 TSV，跳过空行与 # 注释行，并要求每行恰好 wantColumns 列。
func readTable(t *testing.T, path string, wantColumns int) [][]string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取 %s: %v", path, err)
	}
	var rows [][]string
	for i, line := range strings.Split(string(data), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != wantColumns {
			t.Fatalf("%s 第 %d 行有 %d 列，应为 %d 列", path, i+1, len(fields), wantColumns)
		}
		rows = append(rows, fields)
	}
	if len(rows) == 0 {
		t.Fatalf("%s 没有数据行", path)
	}
	return rows
}

// readExcluded 读取排除清单，返回 包 -> 符号 的集合；符号 `*` 表示整包排除。
func readExcluded(t *testing.T) map[string]map[string]bool {
	t.Helper()
	rows := readTable(t, "NAME-MAP-EXCLUDED.tsv", 3)
	excluded := make(map[string]map[string]bool, len(rows))
	for _, row := range rows {
		if excluded[row[0]] == nil {
			excluded[row[0]] = make(map[string]bool)
		}
		excluded[row[0]][row[1]] = true
	}
	return excluded
}

// exportedFuncs 返回目录下所有导出顶层函数名（不含方法、不含 _test.go）。
func exportedFuncs(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("读取目录 %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			t.Fatalf("解析 %s: %v", name, err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !fn.Name.IsExported() {
				continue
			}
			names = append(names, fn.Name.Name)
		}
	}
	sort.Strings(names)
	return names
}

func coveredByBase(name string, bases map[string]bool) bool {
	if bases[name] {
		return true
	}
	for base := range bases {
		for _, suffix := range derivedSuffixes {
			if name == base+suffix {
				return true
			}
		}
	}
	return false
}

func contains(names []string, target string) bool {
	i := sort.SearchStrings(names, target)
	return i < len(names) && names[i] == target
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
