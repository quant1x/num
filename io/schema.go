package io

import (
	"fmt"
	"strings"
)

// ParseSchema 解析命令行的 schema 简写：`name:f|i|b|s|t`，逗号分隔。
//
//	f = Float64，i = Int64，b = Bool，s = String，t = Time。
func ParseSchema(spec string) ([]Column, error) {
	if spec == "" {
		return nil, fmt.Errorf("schema 为空")
	}
	seen := make(map[string]bool)
	var schema []Column
	for _, item := range strings.Split(spec, ",") {
		name, code, ok := strings.Cut(item, ":")
		if !ok {
			return nil, fmt.Errorf("schema 项非法: %s", item)
		}
		if err := validateName(name); err != nil {
			return nil, err
		}
		if seen[name] {
			return nil, fmt.Errorf("schema 列名重复: %s", name)
		}
		seen[name] = true
		var kind Kind
		switch code {
		case "f":
			kind = Float64
		case "i":
			kind = Int64
		case "b":
			kind = Bool
		case "s":
			kind = String
		case "t":
			kind = Time
		default:
			return nil, fmt.Errorf("未知类型码: %s", code)
		}
		schema = append(schema, Column{Name: name, Kind: kind})
	}
	return schema, nil
}
