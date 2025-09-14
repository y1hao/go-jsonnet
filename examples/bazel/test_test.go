package main

import (
	"fmt"
	"testing"

	gjs "github.com/google/go-jsonnet"
)

func TestThings(t *testing.T) {
	vm := gjs.MakeVM()
	out, err := vm.EvaluateFile("example.jsonnet")
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%s", out)
	}
}
