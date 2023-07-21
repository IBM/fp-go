package errors

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
)

// OnNone generates a nullary function that produces a formatted error
func OnNone(msg string, args ...any) func() error {
	return func() error {
		return fmt.Errorf(msg, args...)
	}
}

// OnError generates a unary function that produces a formatted error. The argument
// to that function is the root cause of the error and the message will be augmented with
// a format string containing %w
func OnError(msg string, args ...any) func(error) error {
	return func(err error) error {
		return fmt.Errorf(msg+", Caused By: %w", A.ArrayConcatAll(args, A.Of[any](err))...)
	}
}

// ToString converts an error to a string
func ToString(err error) string {
	return err.Error()
}
