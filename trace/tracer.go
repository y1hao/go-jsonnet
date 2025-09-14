package trace

import (
	"fmt"

	"github.com/google/go-jsonnet"
)

type tracer struct {
}

func NewTracer() *tracer {
	return &tracer{}
}

func (t *tracer) GenerateTrace(filename string) (string, map[int][]*jsonnet.TraceItem, error) {
	vm := jsonnet.MakeTracingVM()
	trace := map[int][]*jsonnet.TraceItem{}
	result, err := vm.EvaluateFileWithTrace(filename, trace)
	if err != nil {
		return "", nil, fmt.Errorf("error generating trace: %w", err)
	}
	return result, trace, nil
}
