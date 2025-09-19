package main

import (
	"fmt"
	"testing"

	gjs "github.com/google/go-jsonnet"
)

func TestThings(t *testing.T) {
	vm := gjs.MakeTracingVM()
	out, trace, err := vm.EvaluateFileWithTrace("testdata/multi-conditionals/0.jsonnet")
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%s", out)
	}
	fmt.Printf("%v\n", trace)
}
