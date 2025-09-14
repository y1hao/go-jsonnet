package trace

type tracer struct {
}

func NewTracer() *tracer {
	return &tracer{}
}

func (t *tracer) GenerateTrace(filename string) (string, map[int][]*TraceFrame, error) {
	return "Test\nline1\nline2\nline3\n", map[int][]*TraceFrame{
		0: {{Filename: "xyz", StartLine: 1, EndLine: 10}, {Filename: "xyz", StartLine: 1, EndLine: 10}},
		1: {{Filename: "123", StartLine: 1, EndLine: 10}},
		2: {{Filename: "456", StartLine: 1, EndLine: 10}},
		3: {{Filename: "789", StartLine: 1, EndLine: 10}},
	}, nil
}
