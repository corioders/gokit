package errors

import (
	"fmt"
	"strings"
	"testing"
)

const message = "error from testing!"

func TestNew(t *testing.T) {
	err := New(message)
	errStr := strings.Split(err.Error(), "\n")
	if errStr[0] != message {
		t.Fatal("Expected first line of error to equal message")
	}
}

func TestWrap(t *testing.T) {
	checkWrappedErr := func(t *testing.T, err error, wrappedErr error, wrappedType string) {
		unwrap := func(t *testing.T, err error) error {
			u, ok := err.(interface {
				Unwrap() error
			})
			if !ok {
				t.Fatal("Expected err to implement Unwrap interface")
			}
			return u.Unwrap()
		}

		e := err.(*internalError)

		if unwrap(t, err) != wrappedErr || e.wrappedErr != wrappedErr {
			t.Fatal("Expected return of unwrap to be " + wrappedType)
		}
	}

	t.Run("nil", func(t *testing.T) {
		err := New(message)
		checkWrappedErr(t, err, nil, "nil")
	})

	t.Run("error", func(t *testing.T) {
		t.Run("withStack", func(t *testing.T) {
			wrappedErr := fmt.Errorf(message)
			err := WithStack(wrappedErr)
			checkWrappedErr(t, err, wrappedErr, "wrappedErr")
		})

		t.Run("withMessage", func(t *testing.T) {
			wrappedErr := fmt.Errorf(message)
			err := WithMessage(wrappedErr, message)
			checkWrappedErr(t, err, wrappedErr, "wrappedErr")
		})
	})
}

func TestIgnoreStack(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		err := New(message)
		stackIgnored := IgnoreStack(err)

		internalErr := stackIgnored.(*internalError)
		if internalErr.stack != nil {
			t.Fatal("Expected error's stack to be nil after IgnoreStack")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		err := fmt.Errorf(message)
		stackIgnored := IgnoreStack(err)

		if stackIgnored != err {
			t.Fatal("Expected return error to be the same as passed error if not of type *internalError")
		}
	})
}

func TestGetErrorMessage(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		err := New(message)

		if GetErrorNoStack(err) != message {
			t.Fatal("Expected return value of GetErrorNoStack to be equal to message")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		err := fmt.Errorf(message)

		if GetErrorNoStack(err) != "GetErrorNoStack in gokit/errors, err not of type *internalError" {
			t.Fatal("Expected return value of GetErrorNoStack to be an `wrong type` error")
		}
	})
}

func TestGetErrorStack(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		err := New(message)
		internalErr := err.(*internalError)

		if GetErrorStack(err) != internalErr.stack.String() {
			t.Fatal("Expected return value of GetErrorStack to be equal to err.stack.String()")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		err := fmt.Errorf(message)

		if GetErrorStack(err) != "GetErrorStack in gokit/errors, err not of type *internalError" {
			t.Fatal("Expected return value of GetErrorStack to be an `wrong type` error")
		}
	})
}
