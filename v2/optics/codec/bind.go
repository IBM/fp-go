// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codec

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/codec/validate"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/semigroup"
)

// ApSL creates an applicative sequencing operator for codecs using a lens.
//
// This function implements the "ApS" (Applicative Sequencing) pattern for codecs,
// allowing you to build up complex codecs by combining a base codec with a field
// accessed through a lens. It's particularly useful for building struct codecs
// field-by-field in a composable way.
//
// The function combines:
//   - Encoding: Extracts the field value using the lens, encodes it with fa, and
//     combines it with the base encoding using the monoid
//   - Validation: Validates the field using the lens and combines the validation
//     with the base validation
//
// # Type Parameters
//
//   - S: The source struct type (what we're building a codec for)
//   - T: The field type accessed by the lens
//   - O: The output type for encoding (must have a monoid)
//   - I: The input type for decoding
//
// # Parameters
//
//   - m: A Monoid[O] for combining encoded outputs
//   - l: A Lens[S, T] that focuses on a specific field in S
//   - fa: A Type[T, O, I] codec for the field type T
//
// # Returns
//
// An Operator[S, S, O, I] that transforms a base codec by adding the field
// specified by the lens.
//
// # How It Works
//
// 1. **Encoding**: When encoding a value of type S:
//   - Extract the field T using l.Get
//   - Encode T to O using fa.Encode
//   - Combine with the base encoding using the monoid
//
// 2. **Validation**: When validating input I:
//   - Validate the field using fa.Validate through the lens
//   - Combine with the base validation
//
// 3. **Type Checking**: Preserves the base type checker
//
// # Example
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec"
//	    "github.com/IBM/fp-go/v2/optics/lens"
//	    S "github.com/IBM/fp-go/v2/string"
//	)
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	// Lenses for Person fields
//	nameLens := lens.MakeLens(
//	    func(p *Person) string { return p.Name },
//	    func(p *Person, name string) *Person { p.Name = name; return p },
//	)
//
//	// Build a Person codec field by field
//	personCodec := F.Pipe1(
//	    codec.Struct[Person]("Person"),
//	    codec.ApSL(S.Monoid, nameLens, codec.String),
//	    // ... add more fields
//	)
//
// # Use Cases
//
//   - Building struct codecs incrementally
//   - Composing codecs for nested structures
//   - Creating type-safe serialization/deserialization
//   - Implementing Do-notation style codec construction
//
// # Notes
//
//   - The monoid determines how encoded outputs are combined
//   - The lens must be total (handle all cases safely)
//   - This is typically used with other ApS functions to build complete codecs
//   - The name is automatically generated for debugging purposes
//
// See also:
//   - validate.ApSL: The underlying validation combinator
//   - reader.ApplicativeMonoid: The monoid-based applicative instance
//   - Lens: The optic for accessing struct fields
func ApSL[S, T, O, I any](
	m Monoid[O],
	l Lens[S, T],
	fa Type[T, O, I],
) Operator[S, S, O, I] {
	name := fmt.Sprintf("ApS[%s x %s]", l, fa)
	rm := reader.ApplicativeMonoid[S](m)

	encConcat := F.Pipe1(
		F.Flow2(
			l.Get,
			fa.Encode,
		),
		semigroup.AppendTo(rm),
	)

	valConcat := validate.ApSL(l, fa.Validate)

	return func(t Type[S, O, I]) Type[S, O, I] {

		return MakeType(
			name,
			t.Is,
			F.Pipe1(
				t.Validate,
				valConcat,
			),
			encConcat(t.Encode),
		)
	}
}

// ApSO creates an applicative sequencing operator for codecs using an optional.
//
// This function implements the "ApS" (Applicative Sequencing) pattern for codecs
// with optional fields, allowing you to build up complex codecs by combining a base
// codec with a field that may or may not be present. It's particularly useful for
// building struct codecs with optional fields in a composable way.
//
// The function combines:
//   - Encoding: Attempts to extract the optional field value, encodes it if present,
//     and combines it with the base encoding using the monoid. If the field is absent,
//     only the base encoding is used.
//   - Validation: Validates the optional field and combines the validation with the
//     base validation using applicative semantics (error accumulation).
//
// # Type Parameters
//
//   - S: The source struct type (what we're building a codec for)
//   - T: The optional field type accessed by the optional
//   - O: The output type for encoding (must have a monoid)
//   - I: The input type for decoding
//
// # Parameters
//
//   - m: A Monoid[O] for combining encoded outputs
//   - o: An Optional[S, T] that focuses on a field in S that may not exist
//   - fa: A Type[T, O, I] codec for the optional field type T
//
// # Returns
//
// An Operator[S, S, O, I] that transforms a base codec by adding the optional field
// specified by the optional.
//
// # How It Works
//
// 1. **Encoding**: When encoding a value of type S:
//   - Try to extract the optional field T using o.GetOption
//   - If present (Some(T)): Encode T to O using fa.Encode and combine with base using monoid
//   - If absent (None): Return only the base encoding unchanged
//
// 2. **Validation**: When validating input I:
//   - Validate the optional field using fa.Validate through o.Set
//   - Combine with the base validation using applicative semantics
//   - Accumulates all validation errors from both base and field
//
// 3. **Type Checking**: Preserves the base type checker
//
// # Difference from ApSL
//
// Unlike ApSL which works with required fields via Lens, ApSO handles optional fields:
//   - ApSL: Field always exists, always encoded
//   - ApSO: Field may not exist, only encoded when present
//   - ApSO uses Optional.GetOption which returns Option[T]
//   - ApSO gracefully handles missing fields without errors
//
// # Example
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec"
//	    "github.com/IBM/fp-go/v2/optics/optional"
//	    S "github.com/IBM/fp-go/v2/string"
//	)
//
//	type Person struct {
//	    Name      string
//	    Nickname  *string  // Optional field
//	}
//
//	// Optional for Person.Nickname
//	nicknameOpt := optional.MakeOptional(
//	    func(p Person) option.Option[string] {
//	        if p.Nickname != nil {
//	            return option.Some(*p.Nickname)
//	        }
//	        return option.None[string]()
//	    },
//	    func(p Person, nick string) Person {
//	        p.Nickname = &nick
//	        return p
//	    },
//	)
//
//	// Build a Person codec with optional nickname
//	personCodec := F.Pipe1(
//	    codec.Struct[Person]("Person"),
//	    codec.ApSO(S.Monoid, nicknameOpt, codec.String),
//	)
//
//	// Encoding with nickname present
//	p1 := Person{Name: "Alice", Nickname: ptr("Ali")}
//	encoded1 := personCodec.Encode(p1)  // Includes nickname
//
//	// Encoding with nickname absent
//	p2 := Person{Name: "Bob", Nickname: nil}
//	encoded2 := personCodec.Encode(p2)  // No nickname in output
//
// # Use Cases
//
//   - Building struct codecs with optional/nullable fields
//   - Handling pointer fields that may be nil
//   - Composing codecs for structures with optional nested data
//   - Creating flexible serialization that omits absent fields
//
// # Notes
//
//   - The monoid determines how encoded outputs are combined when field is present
//   - When the optional field is absent, encoding returns base encoding unchanged
//   - Validation still accumulates errors even for optional fields
//   - The name is automatically generated for debugging purposes
//
// # See Also
//
//   - ApSL: For required fields using Lens
//   - validate.ApS: The underlying validation combinator
//   - Optional: The optic for accessing optional fields
func ApSO[S, T, O, I any](
	m Monoid[O],
	o Optional[S, T],
	fa Type[T, O, I],
) Operator[S, S, O, I] {
	name := fmt.Sprintf("ApS[%s x %s]", o, fa)

	encConcat := F.Flow2(
		o.GetOption,
		option.Map(F.Flow2(
			fa.Encode,
			semigroup.AppendTo(m),
		)),
	)

	valConcat := validate.ApS(o.Set, fa.Validate)

	return func(t Type[S, O, I]) Type[S, O, I] {

		return MakeType(
			name,
			t.Is,
			F.Pipe1(
				t.Validate,
				valConcat,
			),
			func(s S) O {
				to := t.Encode(s)
				return F.Pipe2(
					encConcat(s),
					option.Flap[O](to),
					option.GetOrElse(lazy.Of(to)),
				)
			},
		)
	}
}
