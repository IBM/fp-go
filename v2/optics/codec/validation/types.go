package validation

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Result[A any] = result.Result[A]

	// Either represents a value that can be one of two types: Left (error) or Right (success).
	Either[E, A any] = either.Either[E, A]

	// ContextEntry represents a single entry in the validation context path.
	// It tracks the location and type information during nested validation.
	ContextEntry struct {
		Key    string // The key or field name (for objects/maps)
		Type   string // The expected type name
		Actual any    // The actual value being validated
	}

	// Context is a stack of ContextEntry values representing the path through
	// nested structures during validation. Used to provide detailed error messages.
	Context = []ContextEntry

	// ValidationError represents a single validation failure with context.
	ValidationError struct {
		Value    any     // The value that failed validation
		Context  Context // The path to the value in nested structures
		Messsage string  // Human-readable error message
		Cause    error
	}

	// Errors is a collection of validation errors.
	Errors = []*ValidationError

	ValidationErrors struct {
		Errors Errors
		Cause  error
	}

	// Validation represents the result of a validation operation.
	// Left contains validation errors, Right contains the successfully validated value.
	Validation[A any] = Either[Errors, A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	Reader[R, A any] = reader.Reader[R, A]

	Kleisli[A, B any] = Reader[A, Validation[B]]

	Operator[A, B any] = Kleisli[Validation[A], B]

	Monoid[A any] = monoid.Monoid[A]
)
