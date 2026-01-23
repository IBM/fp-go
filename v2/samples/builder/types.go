// Package builder demonstrates the builder pattern using functional programming concepts
// from fp-go, including validation and transformation of data structures.
package builder

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/codec"
	"github.com/IBM/fp-go/v2/optics/codec/validate"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/result"
)

//go:generate go run ../../main.go lens --dir . --filename gen_lens.go

type (
	// Endomorphism represents a function from type A to type A.
	// It is an alias for endomorphism.Endomorphism[A].
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Result represents a computation that may succeed with a value of type A or fail with an error.
	// It is an alias for result.Result[A].
	Result[A any] = result.Result[A]

	// Option represents an optional value of type A that may or may not be present.
	// It is an alias for option.Option[A].
	Option[A any] = option.Option[A]

	// ReaderOption represents a computation that depends on an environment R and produces
	// an optional value of type A. It is an alias for readeroption.ReaderOption[R, A].
	ReaderOption[R, A any] = readeroption.ReaderOption[R, A]

	// Reader represents a computation that depends on an environment R and produces
	// a value of type A. It is an alias for reader.Reader[R, A].
	Reader[R, A any] = reader.Reader[R, A]

	// Prism represents an optic that focuses on a subset of values of type S that can be
	// converted to type A. It provides bidirectional transformation with validation.
	// It is an alias for prism.Prism[S, A].
	Prism[S, A any] = prism.Prism[S, A]

	// Lens represents an optic that focuses on a field of type A within a structure of type S.
	// It provides getter and setter operations for immutable updates.
	// It is an alias for lens.Lens[S, A].
	Lens[S, A any] = lens.Lens[S, A]

	// Type represents a codec that handles bidirectional transformation between types.
	// A: The validated target type
	// O: The output encoding type
	// I: The input decoding type
	// It is an alias for codec.Type[A, O, I].
	Type[A, O, I any] = codec.Type[A, O, I]

	// Validate represents a validation function that transforms input I into a validated result A.
	// It returns a Validation that contains either the validated value or validation errors.
	// It is an alias for validate.Validate[I, A].
	Validate[I, A any] = validate.Validate[I, A]

	// Validation represents the result of a validation operation.
	// It contains either a validated value of type A (Right) or validation errors (Left).
	// It is an alias for validation.Validation[A].
	Validation[A any] = validation.Validation[A]

	// Encode represents an encoding function that transforms a value of type A into type O.
	// It is used in codecs for the reverse direction of validation.
	// It is an alias for codec.Encode[A, O].
	Encode[A, O any] = codec.Encode[A, O]

	// NonEmptyString is a string type that represents a validated non-empty string.
	// It is used to ensure that string fields contain meaningful data.
	NonEmptyString string

	// AdultAge is an unsigned integer type that represents a validated age
	// that meets adult criteria (typically >= 18).
	AdultAge uint
)

// PartialPerson represents a person record with unvalidated fields.
// This type is typically used as an intermediate representation before
// validation is applied to create a Person instance.
//
// The fp-go:Lens directive generates lens functions for accessing and
// modifying the fields of this struct in a functional way.
//
// fp-go:Lens
type PartialPerson struct {
	// name is the person's name as a raw string, which may be empty or invalid.
	name string

	// age is the person's age as a raw integer, which may be negative or otherwise invalid.
	age int
}

// Person represents a person record with validated fields.
// All fields in this type have been validated and are guaranteed to meet
// specific business rules (non-empty name, adult age).
//
// The fp-go:Lens directive generates lens functions for accessing and
// modifying the fields of this struct in a functional way.
//
// fp-go:Lens
type Person struct {
	// Name is the person's validated name, guaranteed to be non-empty.
	Name NonEmptyString

	// Age is the person's validated age, guaranteed to meet adult criteria.
	Age AdultAge
}
