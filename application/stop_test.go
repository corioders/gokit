package application

import (
	"fmt"
	"io"
	"testing"

	"github.com/corioders/gokit/log"
)

func TestStop(t *testing.T) {
	logger := log.New(io.Discard, "")
	application := New(logger)

	t.Run("no error", func(t *testing.T) {

		stopFuncCalled := false
		application.RegisterOnStop("stopTest", func() error {
			stopFuncCalled = true
			return nil
		})

		err := application.Stop()
		if err != nil {
			t.Fatal("Expected no error, but got:", err)
			return
		}

		if !stopFuncCalled {
			t.Fatal("Expected the stop func to be called")
		}
	})

	logger = log.New(io.Discard, "")
	application = New(logger)

	t.Run("with error", func(t *testing.T) {
		expectedErr := fmt.Errorf("test error")

		application.RegisterOnStop("stopTest", func() error {
			return expectedErr
		})

		err := application.Stop()
		if err != expectedErr {
			t.Fatal("Expected err returned by Stop to be the same as error returned by StopFunc, but got:", err)
		}
	})
}
