// 命令行工具：按主键把增量合并进已有 CSV；文件不存在时直接写增量。
package main

import (
	"fmt"
	"os"
	"strings"

	numio "github.com/quant1x/num/io"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "usage: go_merge <schema> <path> <delta.csv> <key>")
		os.Exit(2)
	}
	schema, err := parseSchema(os.Args[1])
	if err != nil {
		fail(err)
	}
	incoming, err := numio.ReadFile(os.Args[3], schema)
	if err != nil {
		fail(err)
	}
	if err := numio.WriteByKey(os.Args[2], incoming, os.Args[4]); err != nil {
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
