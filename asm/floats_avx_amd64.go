//go:build !noasm && !appengine
// +build !noasm,!appengine

// AUTO-GENERATED BY GOAT -- DO NOT EDIT

package asm

import "unsafe"

//go:noescape
func ___mm256_mul_const_add_to(a, b, c, n unsafe.Pointer)

//go:noescape
func ___mm256_mul_const_to(a, b, c, n unsafe.Pointer)

//go:noescape
func ___mm256_mul_const(a, b, n unsafe.Pointer)

//go:noescape
func ___mm256_mul_to(a, b, c, n unsafe.Pointer)

//go:noescape
func ___mm256_dot(a, b, n, ret unsafe.Pointer)
