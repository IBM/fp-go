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
	"github.com/IBM/fp-go/v2/reader"
)

// validateAlt creates a validation function that tries the first codec's validation,
// and if it fails, tries the second codec's validation as a fallback.
//
// This is an internal helper function that implements the Alternative pattern for
// codec validation. It combines two codec validators using the validate.Alt operation,
// which provides error recovery and fallback logic.
//
// # Type Parameters
//
//   - A: The target type that both codecs decode to
//   - O: The output type that both codecs encode to
//   - I: The input type that both codecs decode from
//
// # Parameters
//
//   - first: The primary codec whose validation is tried first
//   - second: A lazy codec that serves as the fallback. It's only evaluated if the
//     first validation fails.
//
// # Returns
//
// A Validate[I, A] function that tries the first codec's validation, falling back
// to the second if needed. If both fail, errors from both are aggregated.
//
// # Behavior
//
//   - **First succeeds**: Returns the first result, second is never evaluated
//   - **First fails, second succeeds**: Returns the second result
//   - **Both fail**: Aggregates errors from both validators
//
// # Notes
//
//   - The second codec is lazily evaluated for efficiency
//   - This function is used internally by MonadAlt and Alt
//   - The validation context is threaded through both validators
//   - Errors are accumulated using the validation error monoid
func validateAlt[A, O, I any](
	first Type[A, O, I],
	second Lazy[Type[A, O, I]],
) Validate[I, A] {

	return F.Pipe1(
		first.Validate,
		validate.Alt(F.Pipe1(
			second,
			lazy.Map(F.Flip(reader.Curry(Type[A, O, I].Validate))),
		)),
	)
}

// MonadAlt creates a new codec that tries the first codec, and if it fails during
// validation, tries the second codec as a fallback.
//
// This function implements the Alternative typeclass pattern for codecs, enabling
// "try this codec, or else try that codec" logic. It's particularly useful for:
//   - Handling multiple valid input formats
//   - Providing backward compatibility with legacy formats
//   - Implementing graceful degradation in parsing
//   - Supporting union types or polymorphic data
//
// The resulting codec uses the first codec's encoder and combines both validators
// using the Alternative pattern. If both validations fail, errors from both are
// aggregated for comprehensive error reporting.
//
// # Type Parameters
//
//   - A: The target type that both codecs decode to
//   - O: The output type that both codecs encode to
//   - I: The input type that both codecs decode from
//
// # Parameters
//
//   - first: The primary codec to try first. Its encoder is used for the result.
//   - second: A lazy codec that serves as the fallback. It's only evaluated if the
//     first validation fails.
//
// # Returns
//
// A new Type[A, O, I] that combines both codecs with Alternative semantics.
//
// # Behavior
//
// **Validation**:
//   - **First succeeds**: Returns the first result, second is never evaluated
//   - **First fails, second succeeds**: Returns the second result
//   - **Both fail**: Aggregates errors from both validators
//
// **Encoding**:
//   - Always uses the first codec's encoder
//   - This assumes both codecs encode to the same output format
//
// **Type Checking**:
//   - Uses the generic Is[A]() type checker
//   - Validates that values are of type A
//
// # Example: Multiple Input Formats
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec"
//	)
//
//	// Accept integers as either strings or numbers
//	intFromString := codec.IntFromString()
//	intFromNumber := codec.Int()
//
//	// Try parsing as string first, fall back to number
//	flexibleInt := codec.MonadAlt(
//	    intFromString,
//	    func() codec.Type[int, any, any] { return intFromNumber },
//	)
//
//	// Can now decode both "42" and 42
//	result1 := flexibleInt.Decode("42")   // Success(42)
//	result2 := flexibleInt.Decode(42)     // Success(42)
//
// # Example: Backward Compatibility
//
//	// Support both old and new configuration formats
//	newConfigCodec := codec.Struct(/* new format */)
//	oldConfigCodec := codec.Struct(/* old format */)
//
//	// Try new format first, fall back to old format
//	configCodec := codec.MonadAlt(
//	    newConfigCodec,
//	    func() codec.Type[Config, any, any] { return oldConfigCodec },
//	)
//
//	// Automatically handles both formats
//	config := configCodec.Decode(input)
//
// # Example: Error Aggregation
//
//	// Both validations will fail for invalid input
//	result := flexibleInt.Decode("not a number")
//	// Result contains errors from both string and number parsing attempts
//
// # Notes
//
//   - The second codec is lazily evaluated for efficiency
//   - First success short-circuits evaluation (second not called)
//   - Errors are aggregated when both fail
//   - The resulting codec's name is "Alt[<first codec name>]"
//   - Both codecs must have compatible input and output types
//   - The first codec's encoder is always used
//
// # See Also
//
//   - Alt: The curried, point-free version
//   - validate.MonadAlt: The underlying validation operation
//   - Either: For codecs that decode to Either[L, R] types
func MonadAlt[A, O, I any](first Type[A, O, I], second Lazy[Type[A, O, I]]) Type[A, O, I] {
	return MakeType(
		fmt.Sprintf("Alt[%s]", first.Name()),
		Is[A](),
		validateAlt(first, second),
		first.Encode,
	)
}

// Alt creates an operator that adds alternative fallback logic to a codec.
//
// This is the curried, point-free version of MonadAlt. It returns a function that
// can be applied to codecs to add fallback behavior. This style is particularly
// useful for building codec transformation pipelines using function composition.
//
// Alt implements the Alternative typeclass pattern, enabling "try this codec, or
// else try that codec" logic in a composable way.
//
// # Type Parameters
//
//   - A: The target type that both codecs decode to
//   - O: The output type that both codecs encode to
//   - I: The input type that both codecs decode from
//
// # Parameters
//
//   - second: A lazy codec that serves as the fallback. It's only evaluated if the
//     first codec's validation fails.
//
// # Returns
//
// An Operator[A, A, O, I] that transforms codecs by adding alternative fallback logic.
// This operator can be applied to any Type[A, O, I] to create a new codec with
// fallback behavior.
//
// # Behavior
//
// When the returned operator is applied to a codec:
//   - **First succeeds**: Returns the first result, second is never evaluated
//   - **First fails, second succeeds**: Returns the second result
//   - **Both fail**: Aggregates errors from both validators
//
// # Example: Point-Free Style
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/optics/codec"
//	)
//
//	// Create a reusable fallback operator
//	withNumberFallback := codec.Alt(func() codec.Type[int, any, any] {
//	    return codec.Int()
//	})
//
//	// Apply it to different codecs
//	flexibleInt1 := withNumberFallback(codec.IntFromString())
//	flexibleInt2 := withNumberFallback(codec.IntFromHex())
//
// # Example: Pipeline Composition
//
//	// Build a codec pipeline with multiple fallbacks
//	flexibleCodec := F.Pipe2(
//	    primaryCodec,
//	    codec.Alt(func() codec.Type[T, O, I] { return fallback1 }),
//	    codec.Alt(func() codec.Type[T, O, I] { return fallback2 }),
//	)
//	// Tries primary, then fallback1, then fallback2
//
// # Example: Reusable Transformations
//
//	// Create a transformation that adds JSON fallback
//	withJSONFallback := codec.Alt(func() codec.Type[Config, string, string] {
//	    return codec.JSONCodec[Config]()
//	})
//
//	// Apply to multiple codecs
//	yamlWithFallback := withJSONFallback(yamlCodec)
//	tomlWithFallback := withJSONFallback(tomlCodec)
//
// # Notes
//
//   - This is the point-free style version of MonadAlt
//   - Useful for building transformation pipelines with F.Pipe
//   - The second codec is lazily evaluated for efficiency
//   - First success short-circuits evaluation
//   - Errors are aggregated when both fail
//   - Can be composed with other codec operators
//
// # See Also
//
//   - MonadAlt: The direct application version
//   - validate.Alt: The underlying validation operation
//   - F.Pipe: For composing multiple operators
func Alt[A, O, I any](second Lazy[Type[A, O, I]]) Operator[A, A, O, I] {
	return F.Bind2nd(MonadAlt, second)
}
