package liberror

import (
	"fmt"

	"github.com/go-errors/errors"
)

type Failure interface {
	Stack() string
	Message() string
	Actual() error
	IsOurError() bool
	Error() string
	Data() AdditionalData
	Severity() Severity
}

func WrapMessage(err error, msg string) *OurError {
	if err == nil {
		err = fmt.Errorf("%s", msg)
	}
	return NewOurError(err, msg)
}

func WrapStack(err error, msg string) *OurError {
	if err == nil {
		err = fmt.Errorf("%s", msg)
	}
	ourError := NewOurError(err, msg)
	ourError.stack = errors.Wrap(ourError.Error(), 1)
	return ourError
}

func IsSameError(err error, target OurErrorCode) bool {
	typedError, ok := err.(Failure)
	if !ok {
		return false
	}
	for {
		if !typedError.IsOurError() {
			return false
		}
		if typedError.Message() == target.Code() {
			return true
		}
		typedError, ok = typedError.Actual().(Failure)
		if !ok {
			return false
		}
	}
}
