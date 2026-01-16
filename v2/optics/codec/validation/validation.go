package validation

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/either"
)

// Error implements the error interface for ValidationError.
// Returns a generic error message indicating this is a validation error.
// For detailed error information, use String() or Format() methods.

// Error implements the error interface for ValidationError.
// Returns a generic error message.
func (v *ValidationError) Error() string {
	return "ValidationError"
}

// Unwrap returns the underlying cause error if present.
// This allows ValidationError to work with errors.Is and errors.As.
func (v *ValidationError) Unwrap() error {
	return v.Cause
}

// String returns a simple string representation of the validation error.
// Returns the error message prefixed with "ValidationError: ".
func (v *ValidationError) String() string {
	return fmt.Sprintf("ValidationError: %s", v.Messsage)
}

// Format implements fmt.Formatter for custom formatting of ValidationError.
// It includes the context path, message, and optionally the cause error.
// Supports verbs: %s, %v, %+v (with additional details)
func (v *ValidationError) Format(s fmt.State, verb rune) {
	// Build the context path
	path := ""
	for i, entry := range v.Context {
		if i > 0 {
			path += "."
		}
		if entry.Key != "" {
			path += entry.Key
		} else {
			path += entry.Type
		}
	}

	// Start with the path if available
	result := ""
	if path != "" {
		result = fmt.Sprintf("at %s: ", path)
	}

	// Add the message
	result += v.Messsage

	// Add the cause if present
	if v.Cause != nil {
		if s.Flag('+') && verb == 'v' {
			// Verbose format with detailed cause
			result += fmt.Sprintf("\n  caused by: %+v", v.Cause)
		} else {
			result += fmt.Sprintf(" (caused by: %v)", v.Cause)
		}
	}

	// Add value information for verbose format
	if s.Flag('+') && verb == 'v' {
		result += fmt.Sprintf("\n  value: %#v", v.Value)
	}

	fmt.Fprint(s, result)
}

// Failures creates a validation failure from a collection of errors.
// Returns a Left Either containing the errors.
func Failures[T any](err Errors) Validation[T] {
	return either.Left[T](err)
}

// FailureWithMessage creates a validation failure with a custom message.
// Returns a Reader that takes a Context and produces a Validation[T] failure.
// This is useful for creating context-aware validation errors.
//
// Example:
//
//	fail := FailureWithMessage[int]("abc", "expected integer")
//	result := fail([]ContextEntry{{Key: "age", Type: "int"}})
func FailureWithMessage[T any](value any, message string) Reader[Context, Validation[T]] {
	return func(context Context) Validation[T] {
		return Failures[T](A.Of(&ValidationError{
			Value:    value,
			Context:  context,
			Messsage: message,
		}))
	}
}

// FailureWithError creates a validation failure with a custom message and underlying cause.
// Returns a Reader that takes an error, then a Context, and produces a Validation[T] failure.
// This is useful for wrapping errors from other operations while maintaining validation context.
//
// Example:
//
//	fail := FailureWithError[int]("abc", "parse failed")
//	result := fail(parseErr)([]ContextEntry{{Key: "count", Type: "int"}})
func FailureWithError[T any](value any, message string) Reader[error, Reader[Context, Validation[T]]] {
	return func(err error) Reader[Context, Validation[T]] {
		return func(context Context) Validation[T] {
			return Failures[T](A.Of(&ValidationError{
				Value:    value,
				Context:  context,
				Messsage: message,
				Cause:    err,
			}))
		}
	}
}

// Success creates a successful validation result.
// Returns a Right Either containing the validated value.
func Success[T any](value T) Validation[T] {
	return either.Of[Errors](value)
}
