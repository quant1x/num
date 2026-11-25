// 命令行工具：按 schema 读取 CSV 再按 D1 规则写出，用于跨语言互读互写验证。
package main

import (
	"fmt"
	"os"
	"strings"

	numio "github.com/quant1x/num/io"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: go_roundtrip <schema> <in.csv> <out.csv>")
		os.Exit(2)
	}
	schema, err := parseSchema(os.Args[1])
	if err != nil {
		fail(err)
	}
	table, err := numio.ReadFile(os.Args[2], schema)
	if err != nil {
		fail(err)
	}
	if err := numio.WriteFile(os.Args[3], table); err != nil {
		fail(err)
	}
}

func parseSchema(spec string) ([]numio.Column, error) {
	var schema []numio.Column
	for _, item := range strings.Split(spec, ",") {
		name, code, ok := strings.Cut(item, ":")
		if !ok {
			return nil, fmt.Errorf("schema 项非法: %s", item)
		}
		var kind numio.Kind
		switch code {
		case "f":
			kind = numio.Float64
		case "i":
			kind = numio.Int64
		case "b":
			kind = numio.Bool
		case "s":
			kind = numio.String
		case "t":
			kind = numio.Time
		default:
			return nil, fmt.Errorf("未知类型码: %s", code)
		}
		schema = append(schema, numio.Column{Name: name, Kind: kind})
	}
	if len(schema) == 0 {
		return nil, fmt.Errorf("schema 为空")
	}
	return schema, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
