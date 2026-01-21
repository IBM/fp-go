package codec

import (
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/prism"
)

// TypeToPrism converts a Type codec into a Prism optic.
//
// A Type[A, S, S] represents a bidirectional codec that can decode S to A (with validation)
// and encode A back to S. A Prism[S, A] is an optic that can optionally extract an A from S
// and always construct an S from an A.
//
// This conversion bridges the codec and optics worlds, allowing you to use validation-based
// codecs as prisms for functional optics composition.
//
// Type Parameters:
//   - S: The source/encoded type (both input and output)
//   - A: The decoded/focus type
//
// Parameters:
//   - t: A Type[A, S, S] codec where:
//   - Decode: S → Validation[A] (may fail with validation errors)
//   - Encode: A → S (always succeeds)
//   - Name: Provides a descriptive name for the type
//
// Returns:
//   - A Prism[S, A] where:
//   - GetOption: S → Option[A] (Some if decode succeeds, None if validation fails)
//   - ReverseGet: A → S (uses the codec's Encode function)
//   - Name: Inherited from the Type's name
//
// The conversion works as follows:
//   - GetOption: Decodes the value and converts validation result to Option
//     (Right(a) → Some(a), Left(errors) → None)
//   - ReverseGet: Directly uses the Type's Encode function
//   - Name: Preserves the Type's descriptive name
//
// Example:
//
//	// Create a codec for positive integers
//	positiveInt := codec.MakeType[int, int, int](
//	    "PositiveInt",
//	    func(i any) result.Result[int] { ... },
//	    func(i int) codec.Validate[int] {
//	        if i <= 0 {
//	            return validation.FailureWithMessage(i, "must be positive")
//	        }
//	        return validation.Success(i)
//	    },
//	    func(i int) int { return i },
//	)
//
//	// Convert to prism
//	prism := codec.TypeToPrism(positiveInt)
//
//	// Use as prism
//	value := prism.GetOption(42)   // Some(42) - validation succeeds
//	value = prism.GetOption(-5)    // None - validation fails
//	result := prism.ReverseGet(10) // 10 - encoding always succeeds
//
// Use cases:
//   - Composing codecs with other optics (lenses, prisms, traversals)
//   - Using validation logic in optics pipelines
//   - Building complex data transformations with functional composition
//   - Integrating type-safe parsing with optics-based data access
//
// Note: The prism's GetOption will return None for any validation failure,
// discarding the specific error details. If you need error information,
// use the Type's Decode method directly instead.
func TypeToPrism[S, A any](t Type[A, S, S]) Prism[S, A] {
	return prism.MakePrismWithName(
		F.Flow2(
			t.Decode,
			either.ToOption,
		),
		t.Encode,
		t.Name(),
	)
}
