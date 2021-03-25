package errors

import (
	"runtime"
	"strings"
)

// stack represents a stack of program counters.
type stack []uintptr

func callers(skip int) stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3+skip, pcs[:])
	var st stack = pcs[0:n]
	return st
}

type frame struct {
	Func string
	Path string
	Line int
}

func newFrame(pc uintptr) *frame {
	fn := runtime.FuncForPC(pc)
	path, line := fn.FileLine(pc - 1)
	return &frame{
		Func: fn.Name(),
		Line: line,
		Path: path,
	}
}

func isValidPC(pc uintptr) bool {
	fn := runtime.FuncForPC(pc)
	path, _ := fn.FileLine(pc - 1)

	return isValidFramePath(path)
}

func isValidFrame(f *frame) bool {
	return isValidFramePath(f.Path)
}

func isValidFramePath(path string) bool {
	
	// we don't want to include runtime in stack traces
	return !strings.Contains(path, "libexec/src/runtime")
}
