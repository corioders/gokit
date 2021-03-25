package errors

import (
	"github.com/corioders/gokit/constant"
)

type internalError struct {
	message    string
	wrappedErr error
	stack      stack
}

func (e *internalError) Error() string {
	return e.error(true)
}

func (e *internalError) error(withStack bool) string {
	var str string
	if e.message != "" {
		str += e.message

		if e.wrappedErr != nil {
			str += constant.Delimer
		}
	}

	if e.wrappedErr != nil {
		if wrappedErr, ok := e.wrappedErr.(*internalError); ok {
			str += wrappedErr.error(false)
		} else {
			str += e.wrappedErr.Error()
		}
	}

	if withStack && e.stack != nil {
		str += "\n" + e.stack.String()
	}

	return str
}

func (e *internalError) Unwrap() error {
	return e.wrappedErr
}

func New(message string) error {
	return &internalError{
		message: message,
		stack:   callers(0),
	}
}

// WithMessage returns err with message and internally calls Withstack.
func WithMessage(err error, message string) error {
	e := withStack(err)
	e.message = message
	return e
}

func WithStack(err error) error {
	return withStack(err)
}

func withStack(err error) *internalError {
	const skip = 1

	internalErr, ok := err.(*internalError)
	if !ok || internalErr.stack == nil {
		return &internalError{
			wrappedErr: err,
			stack:      callers(skip),
		}
	}

	currentStack := callers(skip)
	// Insert part of currentStack after fist matched element of e.stack.
	// The part inserted is the part before fist matched element. This is needed because
	// just appending current stack won't preserve order and will be very hard to read.
	matchedCurrentStackIndex := -1
	matchedErrorStackIndex := -1
	for i := 0; i < len(currentStack); i++ {
		for j := 0; j < len(internalErr.stack); j++ {
			if currentStack[i] == internalErr.stack[j] {
				matchedCurrentStackIndex = i
				matchedErrorStackIndex = j
				break
			}
		}
		if matchedCurrentStackIndex != -1 {
			break
		}
	}

	if matchedCurrentStackIndex == -1 || matchedErrorStackIndex == -1 {
		return &internalError{
			wrappedErr: internalErr,
			stack:      callers(skip),
		}
	}

	// this occurs when errors.New is called from init func,
	// then the stacks cannot be matched in other place than runtime
	if !isValidPC(currentStack[matchedCurrentStackIndex]) {
		lastValidStackIndex := -1
		for i := len(internalErr.stack) - 1; i >= 0; i-- {
			if isValidPC(internalErr.stack[i]) {
				lastValidStackIndex = i
				break
			}
		}

		stack := make(stack, 0, len(currentStack)+lastValidStackIndex+1)

		stack = append(stack, internalErr.stack[:lastValidStackIndex+1]...)
		stack = append(stack, currentStack...)
		return &internalError{
			wrappedErr: internalErr,
			stack:      stack,
		}
	}

	stack := make(stack, 0, len(internalErr.stack)+matchedCurrentStackIndex)

	stack = append(stack, internalErr.stack[:matchedErrorStackIndex]...)
	stack = append(stack, currentStack[:matchedCurrentStackIndex]...)
	stack = append(stack, internalErr.stack[matchedErrorStackIndex:]...)
	return &internalError{
		wrappedErr: internalErr,
		stack:      stack,
	}
}

// IgnoreStack ignores stack of err,
// note that err should be possessed from New or WithMessage func.
func IgnoreStack(err error) error {
	e, ok := err.(*internalError)
	if !ok {
		return err
	}

	e.stack = nil
	return e
}

// GetErrorMessage returns error message of wrapped error
// note that err should be possessed from New or WithMessage func.
func GetErrorNoStack(err error) string {
	e, ok := err.(*internalError)
	if !ok {
		return "GetErrorNoStack in gokit/errors, err not of type *internalError"
	}

	return e.error(false)
}

// GetErrorStack returns stack of err
// note that err should be possessed from New or WithMessage func.
func GetErrorStack(err error) string {
	e, ok := err.(*internalError)
	if !ok {
		return "GetErrorStack in gokit/errors, err not of type *internalError"
	}

	if e.stack == nil {
		return ""
	}

	return e.stack.String()
}
