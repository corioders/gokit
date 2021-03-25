package log

import (
	"fmt"
	"io"
	"sync"

	"github.com/corioders/gokit/constant"

	"github.com/logrusorgru/aurora"
)

type Logger interface {
	Info(a ...interface{})
	Error(a ...interface{})
	Child(prefix string) Logger
}

type logger struct {
	output   io.Writer
	outputMu sync.Mutex

	prefix string
}

var (
	statusError = aurora.Red("ERROR" + constant.Delimer).String()
	statusInfo  = aurora.Blue("INFO" + constant.Delimer).String()
)

func New(output io.Writer, prefix string) Logger {
	return &logger{
		outputMu: sync.Mutex{},
		output:   output,

		prefix: aurora.Yellow(prefix + constant.Delimer).String(),
	}
}

func (l *logger) Info(a ...interface{}) {
	message := statusInfo + l.prefix + fmt.Sprintln(a...)
	l.log(message)
}

func (l *logger) Error(a ...interface{}) {
	message := statusError + l.prefix + fmt.Sprintln(a...)
	l.log(message)
}

func (l *logger) Child(prefix string) Logger {
	return New(l.output, l.prefix+prefix)
}

func (l *logger) log(s string) {
	l.outputMu.Lock()
	defer l.outputMu.Unlock()
	l.output.Write([]byte(s))
}
