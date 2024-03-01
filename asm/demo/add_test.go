package demo

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	a := 1
	b := Add(&a)
	fmt.Println(b)
}
