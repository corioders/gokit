package application

import (
	"io"
	"testing"

	"github.com/corioders/gokit/log"
)

func TestGetLogger(t *testing.T) {
	logger := log.New(io.Discard, "")
	application := New(logger)

	if application.GetLogger() != logger {
		t.Fatal("Expected GetLogger() to return the same logger that was passed into New")
	}
}
