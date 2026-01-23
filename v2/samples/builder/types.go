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

	Prism[S, A any] = prism.Prism[S, A]
	Lens[S, A any]  = lens.Lens[S, A]

	Type[A, O, I any]  = codec.Type[A, O, I]
	Validate[I, A any] = validate.Validate[I, A]
	Validation[A any]  = validation.Validation[A]
	Encode[A, O any]   = codec.Encode[A, O]

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
	// Name is the person's name as a raw string, which may be empty or invalid.
	Name string

	// Age is the person's age as a raw integer, which may be negative or otherwise invalid.
	Age int
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
