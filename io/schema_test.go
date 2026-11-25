package io

import "testing"

func TestParseSchema(t *testing.T) {
	schema, err := ParseSchema("account:s,close:f,volume:i,flag:b,ts:t")
	if err != nil {
		t.Fatalf("解析: %v", err)
	}
	if len(schema) != 5 {
		t.Fatalf("列数 %d", len(schema))
	}
	want := []Kind{String, Float64, Int64, Bool, Time}
	for i, col := range schema {
		if col.Kind != want[i] {
			t.Fatalf("第 %d 列类型 %s，期望 %s", i+1, col.Kind, want[i])
		}
	}
}

func TestParseSchemaRejectsBadInput(t *testing.T) {
	for _, spec := range []string{"", "close", "Close:f", "a:f,a:i", "a:x"} {
		if _, err := ParseSchema(spec); err == nil {
			t.Fatalf("schema %q 应被拒绝", spec)
		}
	}
}
