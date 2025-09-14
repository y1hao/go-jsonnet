package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-jsonnet/trace"
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
	tracer := trace.NewTracer()

	fmt.Fprintf(os.Stdout, "Building jsonnet file %q...\n", filename)

	result, frames, err := tracer.GenerateTrace(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate trace for file %s: %s", filename, err.Error())
	}

	fmt.Fprintln(os.Stdout, "Build succeeded. Serving trace information on")
	fmt.Fprintln(os.Stdout, "    http://localhost:8080")

	handler := trace.NewServer(filename, result, frames)
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to display trace info: ", err.Error())
	}
}
