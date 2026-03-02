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

// Do creates the initial empty codec to be used as the starting point for
// do-notation style codec construction.
//
// This is the entry point for building up a struct codec field-by-field using
// the applicative and monadic sequencing operators ApSL, ApSO, and Bind.
// It wraps Empty and lifts a lazily-evaluated default Pair[O, A] into a
// Type[A, O, I] that ignores its input and always succeeds with the default value.
//
// # Type Parameters
//
//   - I: The input type for decoding (what the codec reads from)
//   - A: The target struct type being built up (what the codec decodes to)
//   - O: The output type for encoding (what the codec writes to)
//
// # Parameters
//
//   - e: A Lazy[Pair[O, A]] providing the initial default values:
//   - pair.Head(e()): The default encoded output O (e.g. an empty monoid value)
//   - pair.Tail(e()): The initial zero value of the struct A (e.g. MyStruct{})
//
// # Returns
//
//   - A Type[A, O, I] that always decodes to the default A and encodes to the
//     default O, regardless of input. This is then transformed by chaining
//     ApSL, ApSO, or Bind operators to add fields one by one.
//
// # Example Usage
//
// Building a struct codec using do-notation style:
//
//	import (
//	    "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/lazy"
//	    "github.com/IBM/fp-go/v2/optics/codec"
//	    "github.com/IBM/fp-go/v2/optics/lens"
//	    "github.com/IBM/fp-go/v2/pair"
//	    S "github.com/IBM/fp-go/v2/string"
//	)
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLens(
//	    func(p Person) string { return p.Name },
//	    func(p Person, name string) Person { p.Name = name; return p },
//	)
//	ageLens := lens.MakeLens(
//	    func(p Person) int { return p.Age },
//	    func(p Person, age int) Person { p.Age = age; return p },
//	)
//
//	personCodec := F.Pipe2(
//	    codec.Do[any, Person, string](lazy.Of(pair.MakePair("", Person{}))),
//	    codec.ApSL(S.Monoid, nameLens, codec.String()),
//	    codec.ApSL(S.Monoid, ageLens, codec.Int()),
//	)
//
// # Notes
//
//   - Do is typically the first call in a codec pipeline, followed by ApSL, ApSO, or Bind
//   - The lazy pair should use the monoid's empty value for O and the zero value for A
//   - For convenience, use Struct to create the initial codec for named struct types
//
// # See Also
//
//   - Empty: The underlying codec constructor that Do delegates to
//   - ApSL: Applicative sequencing for required struct fields via Lens
//   - ApSO: Applicative sequencing for optional struct fields via Optional
//   - Bind: Monadic sequencing for context-dependent field codecs
//
//go:inline
func Do[I, A, O any](e Lazy[Pair[O, A]]) Type[A, O, I] {
	return Empty[I](e)
}

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

// Bind creates a monadic sequencing operator for codecs using a lens and a Kleisli arrow.
//
// This function implements the "Bind" (monadic bind / chain) pattern for codecs,
// allowing you to build up complex codecs where the codec for a field depends on
// the current decoded value of the struct. Unlike ApSL which uses a fixed field
// codec, Bind accepts a Kleisli arrow — a function from the current struct value S
// to a Type[T, O, I] — enabling context-sensitive codec construction.
//
// The function combines:
//   - Encoding: Evaluates the Kleisli arrow f on the current struct value s to obtain
//     the field codec, extracts the field T using the lens, encodes it with that codec,
//     and combines it with the base encoding using the monoid.
//   - Validation: Validates the base struct first (monadic sequencing), then uses the
//     Kleisli arrow to obtain the field codec for the decoded struct value, and validates
//     the field through the lens. Errors are propagated but NOT accumulated (fail-fast
//     semantics, unlike ApSL which accumulates errors).
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
//   - f: A Kleisli[S, T, O, I] — a function from S to Type[T, O, I] — that produces
//     the field codec based on the current struct value
//
// # Returns
//
// An Operator[S, S, O, I] that transforms a base codec by adding the field
// specified by the lens, where the field codec is determined by the Kleisli arrow.
//
// # How It Works
//
// 1. **Encoding**: When encoding a value of type S:
//   - Evaluate f(s) to obtain the field codec fa
//   - Extract the field T using l.Get
//   - Encode T to O using fa.Encode
//   - Combine with the base encoding using the monoid
//
// 2. **Validation**: When validating input I:
//   - Run the base validation to obtain a decoded S (fail-fast: stop on base failure)
//   - For the decoded S, evaluate f(s) to obtain the field codec fa
//   - Validate the input I using fa.Validate
//   - Set the validated T into S using l.Set
//
// 3. **Type Checking**: Preserves the base type checker
//
// # Difference from ApSL
//
// Unlike ApSL which uses a fixed field codec:
//   - ApSL: Field codec is fixed at construction time; errors are accumulated
//   - Bind: Field codec depends on the current struct value (Kleisli arrow); validation
//     uses monadic sequencing (fail-fast on base failure)
//   - Bind is more powerful but less parallel than ApSL
//
// # Example
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec"
//	    "github.com/IBM/fp-go/v2/optics/lens"
//	    S "github.com/IBM/fp-go/v2/string"
//	)
//
//	type Config struct {
//	    Mode  string
//	    Value int
//	}
//
//	modeLens := lens.MakeLens(
//	    func(c Config) string { return c.Mode },
//	    func(c Config, mode string) Config { c.Mode = mode; return c },
//	)
//
//	// Build a Config codec where the Value codec depends on the Mode
//	configCodec := F.Pipe1(
//	    codec.Struct[Config]("Config"),
//	    codec.Bind(S.Monoid, modeLens, func(c Config) codec.Type[string, string, any] {
//	        return codec.String()
//	    }),
//	)
//
// # Use Cases
//
//   - Building codecs where a field's codec depends on another field's value
//   - Implementing discriminated unions or tagged variants
//   - Context-sensitive validation (e.g., validate field B differently based on field A)
//   - Dependent type-like patterns in codec construction
//
// # Notes
//
//   - The monoid determines how encoded outputs are combined
//   - The lens must be total (handle all cases safely)
//   - Validation uses monadic (fail-fast) sequencing: if the base codec fails,
//     the Kleisli arrow is never evaluated
//   - The name is automatically generated for debugging purposes
//
// See also:
//   - ApSL: Applicative sequencing with a fixed lens codec (error accumulation)
//   - Kleisli: The function type from S to Type[T, O, I]
//   - decode.Bind: The underlying decode-level bind combinator
func Bind[S, T, O, I any](
	m Monoid[O],
	l Lens[S, T],
	f Kleisli[S, T, O, I],
) Operator[S, S, O, I] {
	name := fmt.Sprintf("Bind[%s]", l)
	val := F.Curry2(Type[T, O, I].Validate)

	return func(t Type[S, O, I]) Type[S, O, I] {

		return MakeType(
			name,
			t.Is,
			F.Pipe1(
				t.Validate,
				validate.Bind(l.Set, F.Flow2(f, val)),
			),
			func(s S) O {
				return m.Concat(t.Encode(s), f(s).Encode(l.Get(s)))
			},
		)
	}
}
