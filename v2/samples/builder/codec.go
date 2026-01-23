// Package builder demonstrates codec-based validation and encoding/decoding
// for Person objects using fp-go's optics and validation framework.
//
// This file extends the builder pattern with codec functionality, enabling:
//   - Bidirectional transformation between PartialPerson and Person
//   - Validation with detailed error reporting
//   - Type-safe encoding and decoding operations
package builder

import (
	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/codec"
	"github.com/IBM/fp-go/v2/optics/codec/decode"
	"github.com/IBM/fp-go/v2/optics/codec/validate"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
)

type (
	// PersonCodec is a codec type that handles bidirectional transformation
	// between Person and PartialPerson using endomorphisms.
	//
	// Type parameters:
	//   - A: *Person - The validated target type
	//   - O: Endomorphism[*PartialPerson] - The output encoding type (builder)
	//   - I: Endomorphism[*PartialPerson] - The input decoding type (builder)
	//
	// This codec enables:
	//   - Validation: Converting a PartialPerson builder to a validated Person
	//   - Encoding: Converting a Person back to a PartialPerson builder
	PersonCodec = Type[*Person, Endomorphism[*PartialPerson], Endomorphism[*PartialPerson]]
)

var (
	// nameCodec is a codec for validating and transforming name fields.
	// It uses namePrism to ensure names are non-empty strings.
	//
	// Validation: string -> Result[NonEmptyString]
	// Encoding: NonEmptyString -> string
	nameCodec = codec.FromRefinement(namePrism)

	// ageCodec is a codec for validating and transforming age fields.
	// It uses agePrism to ensure ages meet adult criteria (>= 18).
	//
	// Validation: int -> Result[AdultAge]
	// Encoding: AdultAge -> int
	ageCodec = codec.FromRefinement(agePrism)
)

// makePersonValidate creates a validation function that transforms a PartialPerson
// builder (endomorphism) into a validated Person.
//
// The validation process:
//  1. Applies the builder endomorphism to an empty PartialPerson
//  2. Extracts and validates the Name field using nameCodec
//  3. Extracts and validates the Age field using ageCodec
//  4. Combines all validations using applicative composition
//  5. Returns either a validated Person or a collection of validation errors
//
// This function uses the Reader monad to thread validation context through
// the computation, and ReaderEither to accumulate validation errors.
//
// Returns:
//
//	A Validate function that takes a PartialPerson builder and returns
//	a Reader that produces a Validation result (either errors or a Person)
func makePersonValidate() Validate[Endomorphism[*PartialPerson], *Person] {

	// Create a monoid for combining validation operations
	// This allows multiple field validations to be composed together
	rdrMonoid := validate.ApplicativeMonoid[*PartialPerson](endomorphism.Monoid[*Person]())

	// allOfRdr combines an array of validation readers into a single reader
	allOfRdr := monoid.ConcatAll(rdrMonoid)

	// valName validates the Name field:
	// 1. Extract name from PartialPerson
	// 2. Validate using nameCodec (ensures non-empty)
	// 3. Map to a Person name setter if valid
	valName := F.Flow3(
		partialPersonLenses.name.Get,
		nameCodec.Validate,
		decode.Map[validation.Context](personLenses.Name.Set),
	)

	// valAge validates the Age field:
	// 1. Extract age from PartialPerson
	// 2. Validate using ageCodec (ensures >= 18)
	// 3. Map to a Person age setter if valid
	valAge := F.Flow3(
		partialPersonLenses.age.Get,
		ageCodec.Validate,
		decode.Map[validation.Context](personLenses.Age.Set),
	)

	// Collect all field validators
	vals := A.From(valName, valAge)

	// Combine all validations and apply to an empty Person
	return F.Flow3(
		identity.Flap[*PartialPerson](emptyPartialPerson),
		allOfRdr(vals),
		decode.Map[validation.Context](identity.Flap[*Person](emptyPerson)),
	)
}

// makePersonCodec creates a complete codec for Person objects.
//
// The codec provides:
//   - Type checking: Verifies if a value is a *Person
//   - Validation: Converts PartialPerson builders to validated Person instances
//   - Encoding: Converts Person instances back to PartialPerson builders
//
// This enables bidirectional transformation with validation:
//   - Decode: Endomorphism[*PartialPerson] -> Validation[*Person]
//   - Encode: *Person -> Endomorphism[*PartialPerson]
//
// Returns:
//
//	A PersonCodec that can validate, encode, and decode Person objects
func makePersonCodec() PersonCodec {
	return codec.MakeType(
		"Person",
		codec.Is[*Person](),
		makePersonValidate(),
		buildEndomorphism(),
	)
}
