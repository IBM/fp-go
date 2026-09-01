// Package builder demonstrates the builder pattern using functional programming concepts
// from fp-go, including validation and transformation of data structures.
package builder

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/internal/common"
)

//go:generate go run ../../main.go lens --dir . --filename gen_lens.go

type (
	// Endomorphism represents a function from type A to type A.
	// It is an alias for endomorphism.Endomorphism[A].
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Prism represents an optic that focuses on a subset of values of type S that can be
	// converted to type A. It provides bidirectional transformation with validation.
	// It is an alias for common.Prism[S, A].
	Prism[S, A any] = common.Prism[S, A]

	// Lens represents an optic that focuses on a field of type A within a structure of type S.
	// It provides getter and setter operations for immutable updates.
	// It is an alias for common.Lens[S, A].
	Lens[S, A any] = common.Lens[S, A]

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
