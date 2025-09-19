package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

func usage(o io.Writer) {
	fmt.Fprintln(o)
	fmt.Fprintln(o, "jsonnet-trace {<option>} { <filename> }")
	fmt.Fprintln(o, "  Build a jsonnet file and collect trace information for debugging.")
	fmt.Fprintln(o, "  The built outcome and trace information will be served on localhost:8080.")
	fmt.Fprintln(o, "  Only a single root file is supported (but it can import other files).")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Available options:")
	fmt.Fprintln(o, "  -h / --help                This message")
	fmt.Fprintln(o)
}

func main() {
	showHelp := flag.Bool("help", false, "Show usage info")
	flag.Parse()

	if showHelp != nil && *showHelp {
		usage(os.Stderr)
		return
	}

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "jsonnet-trace can only take a single file")
		usage(os.Stderr)
	}

	filename := flag.Args()[0]

	fmt.Fprintf(os.Stdout, "Building jsonnet file %q...\n", filename)

	result, trace, err := buildWithTrace(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate trace for file %s: %s", filename, err.Error())
	}

	fmt.Print(result)
	fmt.Printf("%v\n", trace)
}

func buildWithTrace(filename string) (string, map[int]*ast.LocationRange, error) {
	vm := jsonnet.MakeTracingVM()
	result, trace, err := vm.EvaluateFileWithTrace(filename)
	if err != nil {
		return "", nil, fmt.Errorf("error generating trace: %w", err)
	}
	return result, trace, nil
}
