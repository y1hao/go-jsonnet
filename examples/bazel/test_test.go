package main

import (
	"fmt"
	"testing"

	gjs "github.com/google/go-jsonnet"
)

func TestThings(t *testing.T) {
	vm := gjs.MakeTracingVM()
	out, err := vm.EvaluateFileWithTrace("testdata/comprehension.jsonnet", map[int][]*gjs.TraceItem{})
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%s", out)
	}
}
