package codec

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/internal/formatting"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/codec/decode"
	"github.com/IBM/fp-go/v2/optics/codec/validate"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/decoder"
	"github.com/IBM/fp-go/v2/optics/encoder"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerresult"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Formattable represents a type that can be formatted as a string representation.
	// It provides a way to obtain a human-readable description of a type or value.
	Formattable = formatting.Formattable

	// ReaderResult represents a computation that depends on an environment R,
	// produces a value A, and may fail with an error.
	ReaderResult[R, A any] = readerresult.ReaderResult[R, A]

	// Lazy represents a lazily evaluated value.
	Lazy[A any] = lazy.Lazy[A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	Reader[R, A any] = reader.Reader[R, A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Result represents a computation that may fail with an error.
	Result[A any] = result.Result[A]

	// Codec combines a Decoder and an Encoder for bidirectional transformations.
	// It can decode input I to type A and encode type A to output O.
	//
	// This is a simple struct that pairs a decoder with an encoder, providing
	// the basic building blocks for bidirectional data transformation. Unlike
	// the Type interface, Codec is a concrete struct without validation context
	// or type checking capabilities.
	//
	// Type Parameters:
	//   - I: The input type to decode from
	//   - O: The output type to encode to
	//   - A: The intermediate type (decoded to, encoded from)
	//
	// Fields:
	//   - Decode: A decoder that transforms I to A
	//   - Encode: An encoder that transforms A to O
	//
	// Example:
	//   A Codec[string, string, int] can decode strings to integers and
	//   encode integers back to strings.
	//
	// Note: For most use cases, prefer using the Type interface which provides
	// additional validation and type checking capabilities.
	Codec[I, O, A any] struct {
		Decode decoder.Decoder[I, A]
		Encode encoder.Encoder[O, A]
	}

	// Validation represents the result of a validation operation that may contain
	// validation errors or a successfully validated value of type A.
	Validation[A any] = validation.Validation[A]

	// Context provides contextual information for validation operations,
	// such as the current path in a nested structure.
	Context = validation.Context

	// Validate is a function that validates input I to produce type A.
	// It takes an input and returns a Reader that depends on the validation Context.
	//
	// The Validate type is the core validation abstraction, defined as:
	//   Reader[I, Decode[Context, A]]
	//
	// This means:
	//  1. It takes an input of type I
	//  2. Returns a Reader that depends on validation Context
	//  3. That Reader produces a Validation[A] (Either[Errors, A])
	//
	// This layered structure allows validators to:
	//   - Access the input value
	//   - Track validation context (path in nested structures)
	//   - Accumulate multiple validation errors
	//   - Compose with other validators
	//
	// Example:
	//   A Validate[string, int] takes a string and returns a context-aware
	//   function that validates and converts it to an integer.
	Validate[I, A any] = validate.Validate[I, A]

	// Decode is a function that decodes input I to type A with validation.
	// It returns a Validation result directly.
	//
	// The Decode type is defined as:
	//   Reader[I, Validation[A]]
	//
	// This is simpler than Validate as it doesn't require explicit context passing.
	// The context is typically created automatically when the decoder is invoked.
	//
	// Decode is used when:
	//   - You don't need to manually manage validation context
	//   - You want a simpler API for basic validation
	//   - You're working at the top level of validation
	//
	// Example:
	//   A Decode[string, int] takes a string and returns a Validation[int]
	//   which is Either[Errors, int].
	Decode[I, A any] = decode.Decode[I, A]

	// Encode is a function that encodes type A to output O.
	//
	// Encode is simply a Reader[A, O], which is a function from A to O.
	// Encoders are pure functions with no error handling - they assume
	// the input is valid.
	//
	// Encoding is the inverse of decoding:
	//   - Decoding: I -> Validation[A] (may fail)
	//   - Encoding: A -> O (always succeeds)
	//
	// Example:
	//   An Encode[int, string] takes an integer and returns its string
	//   representation.
	Encode[A, O any] = Reader[A, O]

	// Decoder is an interface for types that can decode and validate input.
	//
	// A Decoder transforms input of type I into a validated value of type A,
	// providing detailed error information when validation fails. It supports
	// both context-aware validation (via Validate) and direct decoding (via Decode).
	//
	// Type Parameters:
	//   - I: The input type to decode from
	//   - A: The target type to decode to
	//
	// Methods:
	//   - Name(): Returns a descriptive name for this decoder (used in error messages)
	//   - Validate(I): Returns a context-aware validation function that can track
	//     the path through nested structures
	//   - Decode(I): Directly decodes input to a Validation result with a fresh context
	//
	// The Validate method is more flexible as it returns a Reader that can be called
	// with different contexts, while Decode is a convenience method that creates a
	// new context automatically.
	//
	// Example:
	//   A Decoder[string, int] can decode strings to integers with validation.
	Decoder[I, A any] interface {
		Name() string
		Validate(I) Decode[Context, A]
		Decode(I) Validation[A]
	}

	// Encoder is an interface for types that can encode values.
	//
	// An Encoder transforms values of type A into output format O. This is the
	// inverse operation of decoding, allowing bidirectional transformations.
	//
	// Type Parameters:
	//   - A: The source type to encode from
	//   - O: The output type to encode to
	//
	// Methods:
	//   - Encode(A): Transforms a value of type A into output format O
	//
	// Encoders are pure functions with no validation or error handling - they
	// assume the input is valid. Validation should be performed during decoding.
	//
	// Example:
	//   An Encoder[int, string] can encode integers to their string representation.
	Encoder[A, O any] interface {
		// Encode transforms a value of type A into output format O.
		Encode(A) O
	}

	// Type is a bidirectional codec that combines encoding, decoding, validation,
	// and type checking capabilities. It represents a complete specification of
	// how to work with a particular type.
	//
	// Type is the central abstraction in the codec package, providing:
	//   - Decoding: Transform input I to validated type A
	//   - Encoding: Transform type A to output O
	//   - Validation: Context-aware validation with detailed error reporting
	//   - Type Checking: Runtime type verification via Is()
	//   - Formatting: Human-readable type descriptions via Name()
	//
	// Type Parameters:
	//   - A: The target type (what we decode to and encode from)
	//   - O: The output type (what we encode to)
	//   - I: The input type (what we decode from)
	//
	// Common patterns:
	//   - Type[A, A, A]: Identity codec (no transformation)
	//   - Type[A, string, string]: String-based serialization
	//   - Type[A, any, any]: Generic codec accepting any input/output
	//   - Type[A, JSON, JSON]: JSON codec
	//
	// Methods:
	//   - Name(): Returns the codec's descriptive name
	//   - Validate(I): Returns context-aware validation function
	//   - Decode(I): Decodes input with automatic context creation
	//   - Encode(A): Encodes value to output format
	//   - AsDecoder(): Returns this Type as a Decoder interface
	//   - AsEncoder(): Returns this Type as an Encoder interface
	//   - Is(any): Checks if a value can be converted to type A
	//
	// Example usage:
	//   intCodec := codec.Int()                    // Type[int, int, any]
	//   stringCodec := codec.String()              // Type[string, string, any]
	//   intFromString := codec.IntFromString()     // Type[int, string, string]
	//
	//   // Decode
	//   result := intFromString.Decode("42")       // Validation[int]
	//
	//   // Encode
	//   str := intFromString.Encode(42)            // "42"
	//
	//   // Type check
	//   isInt := intCodec.Is(42)                   // Right(42)
	//   notInt := intCodec.Is("42")                // Left(error)
	//
	// Composition:
	//   Types can be composed using operators like Alt, Map, Chain, and Pipe
	//   to build complex codecs from simpler ones.
	Type[A, O, I any] interface {
		Formattable
		Decoder[I, A]
		Encoder[A, O]
		AsDecoder() Decoder[I, A]
		AsEncoder() Encoder[A, O]
		Is(any) Result[A]
	}

	// Endomorphism represents a function from type A to itself (A -> A).
	// It forms a monoid under function composition.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Pair represents a tuple of two values of types L and R.
	Pair[L, R any] = pair.Pair[L, R]

	// Prism is an optic that focuses on a part of a sum type S that may or may not
	// contain a value of type A. It provides a way to preview and review values.
	Prism[S, A any] = prism.Prism[S, A]

	// Refinement represents the concept that B is a specialized type of A.
	// It's an alias for Prism[A, B], providing a semantic name for type refinement operations.
	//
	// A refinement allows you to:
	//   - Preview: Try to extract a B from an A (may fail if A is not a B)
	//   - Review: Inject a B back into an A
	//
	// This is useful for working with subtypes, validated types, or constrained types.
	//
	// Example:
	//   - Refinement[int, PositiveInt] - refines int to positive integers only
	//   - Refinement[string, NonEmptyString] - refines string to non-empty strings
	//   - Refinement[any, User] - refines any to User type
	Refinement[A, B any] = Prism[A, B]

	// Kleisli represents a Kleisli arrow in the codec context.
	// It's a function that takes a value of type A and returns a codec Type[B, O, I].
	//
	// This is the fundamental building block for codec transformations and compositions.
	// Kleisli arrows allow you to:
	//   - Chain codec operations
	//   - Build dependent codecs (where the next codec depends on the previous result)
	//   - Create codec pipelines
	//
	// Type Parameters:
	//   - A: The input type to the function
	//   - B: The target type that the resulting codec decodes to
	//   - O: The output type that the resulting codec encodes to
	//   - I: The input type that the resulting codec decodes from
	//
	// Example:
	//   A Kleisli[string, int, string, string] takes a string and returns a codec
	//   that can decode strings to ints and encode ints to strings.
	Kleisli[A, B, O, I any] = Reader[A, Type[B, O, I]]

	// Operator is a specialized Kleisli arrow that transforms codecs.
	// It takes a codec Type[A, O, I] and returns a new codec Type[B, O, I].
	//
	// Operators are the primary way to build codec transformation pipelines.
	// They enable functional composition of codec transformations using F.Pipe.
	//
	// Type Parameters:
	//   - A: The source type that the input codec decodes to
	//   - B: The target type that the output codec decodes to
	//   - O: The output type (same for both input and output codecs)
	//   - I: The input type (same for both input and output codecs)
	//
	// Common operators include:
	//   - Map: Transforms the decoded value
	//   - Chain: Sequences dependent codec operations
	//   - Alt: Provides alternative fallback codecs
	//   - Refine: Adds validation constraints
	//
	// Example:
	//   An Operator[int, PositiveInt, int, any] transforms a codec that decodes
	//   to int into a codec that decodes to PositiveInt (with validation).
	//
	// Usage with F.Pipe:
	//   codec := F.Pipe2(
	//       baseCodec,
	//       operator1,  // Operator[A, B, O, I]
	//       operator2,  // Operator[B, C, O, I]
	//   )
	Operator[A, B, O, I any] = Kleisli[Type[A, O, I], B, O, I]

	// Monoid represents an algebraic structure with an associative binary operation
	// and an identity element.
	//
	// A Monoid[A] provides:
	//   - Empty(): Returns the identity element
	//   - Concat(A, A): Combines two values associatively
	//
	// Monoid laws:
	//   1. Left Identity: Concat(Empty(), a) = a
	//   2. Right Identity: Concat(a, Empty()) = a
	//   3. Associativity: Concat(Concat(a, b), c) = Concat(a, Concat(b, c))
	//
	// In the codec context, monoids are used to:
	//   - Combine multiple codecs with specific semantics
	//   - Build codec chains with fallback behavior (AltMonoid)
	//   - Aggregate validation results (ApplicativeMonoid)
	//   - Compose codec transformations
	//
	// Example monoids for codecs:
	//   - AltMonoid: First success wins (alternative semantics)
	//   - ApplicativeMonoid: Combines successful results using inner monoid
	//   - AlternativeMonoid: Combines applicative and alternative behaviors
	Monoid[A any] = monoid.Monoid[A]
)
