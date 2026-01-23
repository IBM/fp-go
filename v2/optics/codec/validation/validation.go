package validation

import (
	"fmt"
	"log/slog"

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

// LogValue implements the slog.LogValuer interface for ValidationError.
// It provides structured logging representation of the validation error.
// Returns a slog.Value containing the error details as a group with
// message, value, context path, and optional cause.
//
// This method is called automatically when logging a ValidationError with slog.
//
// Example:
//
//	err := &ValidationError{Value: "abc", Messsage: "expected number"}
//	slog.Error("validation failed", "error", err)
//	// Logs: error={message="expected number" value="abc"}
func (v *ValidationError) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("message", v.Messsage),
		slog.Any("value", v.Value),
	}

	// Add context path if available
	if len(v.Context) > 0 {
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
		attrs = append(attrs, slog.String("path", path))
	}

	// Add cause if present
	if v.Cause != nil {
		attrs = append(attrs, slog.Any("cause", v.Cause))
	}

	return slog.GroupValue(attrs...)
}

// Error implements the error interface for ValidationErrors.
// Returns a generic error message indicating validation errors occurred.
func (ve *validationErrors) Error() string {
	if len(ve.errors) == 0 {
		return "ValidationErrors: no errors"
	}
	if len(ve.errors) == 1 {
		return "ValidationErrors: 1 error"
	}
	return fmt.Sprintf("ValidationErrors: %d errors", len(ve.errors))
}

// Unwrap returns the underlying cause error if present.
// This allows ValidationErrors to work with errors.Is and errors.As.
func (ve *validationErrors) Unwrap() error {
	return ve.cause
}

// String returns a simple string representation of all validation errors.
// Each error is listed on a separate line with its index.
func (ve *validationErrors) String() string {
	if len(ve.errors) == 0 {
		return "ValidationErrors: no errors"
	}

	result := fmt.Sprintf("ValidationErrors (%d):\n", len(ve.errors))
	for i, err := range ve.errors {
		result += fmt.Sprintf("  [%d] %s\n", i, err.String())
	}

	if ve.cause != nil {
		result += fmt.Sprintf("  caused by: %v\n", ve.cause)
	}

	return result
}

// Format implements fmt.Formatter for custom formatting of ValidationErrors.
// Supports verbs: %s, %v, %+v (with additional details)
// %s and %v: compact format with error count
// %+v: verbose format with all error details
func (ve *validationErrors) Format(s fmt.State, verb rune) {
	if len(ve.errors) == 0 {
		fmt.Fprint(s, "ValidationErrors: no errors")
		return
	}

	// For simple format, just show the count
	if verb == 's' || (verb == 'v' && !s.Flag('+')) {
		if len(ve.errors) == 1 {
			fmt.Fprint(s, "ValidationErrors: 1 error")
		} else {
			fmt.Fprintf(s, "ValidationErrors: %d errors", len(ve.errors))
		}
		return
	}

	// Verbose format with all details
	if s.Flag('+') && verb == 'v' {
		fmt.Fprintf(s, "ValidationErrors (%d):\n", len(ve.errors))
		for i, err := range ve.errors {
			fmt.Fprintf(s, "  [%d] ", i)
			err.Format(s, verb)
			fmt.Fprint(s, "\n")
		}

		if ve.cause != nil {
			fmt.Fprintf(s, "  root cause: %+v\n", ve.cause)
		}
	}
}

// LogValue implements the slog.LogValuer interface for ValidationErrors.
// It provides structured logging representation of multiple validation errors.
// Returns a slog.Value containing the error count and individual errors as a group.
//
// This method is called automatically when logging ValidationErrors with slog.
//
// Example:
//
//	errors := &ValidationErrors{Errors: []*ValidationError{{Messsage: "error1"}, {Messsage: "error2"}}}
//	slog.Error("validation failed", "errors", errors)
//	// Logs: errors={count=2 errors=[...]}
func (ve *validationErrors) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.Int("count", len(ve.errors)),
	}

	// Add individual errors as a group
	if len(ve.errors) > 0 {
		errorAttrs := make([]slog.Attr, len(ve.errors))
		for i, err := range ve.errors {
			errorAttrs[i] = slog.Any(fmt.Sprintf("error_%d", i), err)
		}
		attrs = append(attrs, slog.Any("errors", slog.GroupValue(errorAttrs...)))
	}

	// Add cause if present
	if ve.cause != nil {
		attrs = append(attrs, slog.Any("cause", ve.cause))
	}

	return slog.GroupValue(attrs...)
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

// MakeValidationErrors converts a collection of validation errors into a single error.
// It wraps the Errors slice in a ValidationErrors struct that implements the error interface.
// This is useful for converting validation failures into standard Go errors.
//
// Parameters:
//   - errors: A slice of ValidationError pointers representing validation failures
//
// Returns:
//   - An error that contains all the validation errors and can be used with standard error handling
//
// Example:
//
//	errors := Errors{
//	    &ValidationError{Value: "abc", Messsage: "expected number"},
//	    &ValidationError{Value: nil, Messsage: "required field"},
//	}
//	err := MakeValidationErrors(errors)
//	fmt.Println(err) // Output: ValidationErrors: 2 errors
func MakeValidationErrors(errors Errors) error {
	return &validationErrors{errors: errors}
}

// ToResult converts a Validation[T] to a Result[T].
// It transforms the Left side (validation errors) into a standard error using MakeValidationErrors,
// while preserving the Right side (successful value) unchanged.
// This is useful for integrating validation results with code that expects Result types.
//
// Type Parameters:
//   - T: The type of the successfully validated value
//
// Parameters:
//   - val: A Validation[T] which is Either[Errors, T]
//
// Returns:
//   - A Result[T] which is Either[error, T], with validation errors converted to a single error
//
// Example:
//
//	validation := Success[int](42)
//	result := ToResult(validation) // Result containing 42
//
//	validation := Failures[int](Errors{&ValidationError{Messsage: "invalid"}})
//	result := ToResult(validation) // Result containing ValidationErrors error
func ToResult[T any](val Validation[T]) Result[T] {
	return either.MonadMapLeft(val, MakeValidationErrors)
}
