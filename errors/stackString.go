package errors

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora"
)

// linesAfter is number of source lines after traced line to display.
const linesAfter = 3

// linesBefore is number of source lines before traced line to display.
const linesBefore = 3

func (s stack) String() string {
	expectedRows := (linesBefore + linesAfter + 3) * len(s)
	rows := make([]string, 0, expectedRows)
	for _, pc := range s {
		frame := newFrame(pc)

		// Don't display runtime.
		if !isValidFrame(frame) {
			break
		}

		message := aurora.Sprintf(aurora.Bold("%s:%d"), frame.Path, frame.Line)
		rows = append(rows, message)
		rows = sourceRows(rows, frame)
	}
	return strings.Join(rows, "\n")
}

func sourceRows(rows []string, frame *frame) []string {
	lines := readLines(frame.Path)
	if lines == nil {
		return rows
	}
	if len(lines) < frame.Line {
		return rows
	}

	current := frame.Line - 1
	start := current - linesBefore
	end := current + linesAfter
	for i := start; i <= end; i++ {
		if i < 0 || i >= len(lines) {
			continue
		}
		line := lines[i]
		var message string
		if i == frame.Line-1 {
			message = aurora.Red(fmt.Sprintf("%d\t%s", i+1, line)).String()
		} else {
			message = aurora.Sprintf("%d\t%s", aurora.Blue(i+1), line)
		}
		rows = append(rows, message)
	}
	return append(rows, "")
}

var cache = sync.Map{}

func readLines(path string) []string {
	var lines []string

	l, ok := cache.Load(path)
	if ok {
		lines = l.([]string)
		return lines
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	lines = strings.Split(string(b), "\n")

	cache.Store(path, lines)
	return lines
}
