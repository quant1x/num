//go:build !noasm && !appengine
// +build !noasm,!appengine

package demo

import "unsafe"

// 自动生成add.s
//go:generate clang -S -DENABLE_AVX2 -target x86_64-unknown-none -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -O3 -fno-builtin -ffast-math -mavx cpp/add.c -o cpp/add.s
//go:generate c2goasm -a cpp/add.s add_amd64.s

//go:noescape
func _Add(fl, result unsafe.Pointer)

func Add(fl *int) int {
	var _f4 int
	_Add(unsafe.Pointer(fl), unsafe.Pointer(&_f4))
	return _f4
}
