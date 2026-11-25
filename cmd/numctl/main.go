// numctl 是 num 的 Go 单二进制入口：CSV 交换与数值统计，不依赖 Python 运行时。
package main

import (
	"fmt"
	stdio "io"
	"os"

	"github.com/quant1x/num"
	numio "github.com/quant1x/num/io"
)

const usage = `usage:
  numctl csv roundtrip <schema> <in.csv> <out.csv>
  numctl csv merge     <schema> <path> <delta.csv> <key>
  numctl csv stats     <schema> <file.csv> <column>

schema: name:f|i|b|s|t,...（f=f64，i=i64，b=bool，s=string，t=time）`

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, out stdio.Writer) error {
	if len(args) < 2 || args[0] != "csv" {
		return fmt.Errorf("%s", usage)
	}
	switch args[1] {
	case "roundtrip":
		if len(args) != 5 {
			return fmt.Errorf("%s", usage)
		}
		schema, err := numio.ParseSchema(args[2])
		if err != nil {
			return err
		}
		table, err := numio.ReadFile(args[3], schema)
		if err != nil {
			return err
		}
		return numio.WriteFile(args[4], table)
	case "merge":
		if len(args) != 6 {
			return fmt.Errorf("%s", usage)
		}
		schema, err := numio.ParseSchema(args[2])
		if err != nil {
			return err
		}
		incoming, err := numio.ReadFile(args[4], schema)
		if err != nil {
			return err
		}
		return numio.WriteByKey(args[3], incoming, args[5])
	case "stats":
		if len(args) != 5 {
			return fmt.Errorf("%s", usage)
		}
		return stats(args[2], args[3], args[4], out)
	default:
		return fmt.Errorf("%s", usage)
	}
}

func stats(schemaSpec, path, column string, out stdio.Writer) error {
	schema, err := numio.ParseSchema(schemaSpec)
	if err != nil {
		return err
	}
	table, err := numio.ReadFile(path, schema)
	if err != nil {
		return err
	}
	index := -1
	for i, col := range table.Columns {
		if col.Name == column {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("列不存在: %s", column)
	}
	col := &table.Columns[index]
	values := make([]float64, 0, len(col.Missing))
	for row, missing := range col.Missing {
		if missing {
			continue
		}
		switch col.Kind {
		case numio.Float64:
			values = append(values, col.Floats[row])
		case numio.Int64:
			values = append(values, float64(col.Ints[row]))
		default:
			return fmt.Errorf("列 %s 不是数值列（%s）", column, col.Kind)
		}
	}
	if len(values) == 0 {
		return fmt.Errorf("列 %s 没有有效数值", column)
	}
	_, err = fmt.Fprintf(out, "count=%d\nsum=%.6g\nmean=%.6g\nstd=%.6g\nmin=%.6g\nmax=%.6g\n",
		len(values), num.Sum(values), num.Mean(values), num.Std(values), num.Min(values), num.Max(values))
	return err
}
