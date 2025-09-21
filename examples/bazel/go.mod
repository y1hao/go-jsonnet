module example/go-jsonnet-using-bazel

go 1.24.0

toolchain go1.24.2

require github.com/y1hao/go-jsonnet v0.21.0

require (
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace github.com/y1hao/go-jsonnet => ../../
