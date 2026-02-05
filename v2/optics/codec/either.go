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

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/codec/validate"
)

// encodeEither creates an encoder for Either[A, B] values.
//
// This function produces an encoder that handles both Left and Right cases of an Either value.
// It uses the provided codecs to encode the Left (A) and Right (B) values respectively.
//
// # Type Parameters
//
//   - A: The type of the Left value
//   - B: The type of the Right value
//   - O: The output type after encoding
//   - I: The input type for validation (not used in encoding)
//
// # Parameters
//
//   - leftItem: The codec for encoding Left values of type A
//   - rightItem: The codec for encoding Right values of type B
//
// # Returns
//
// An Encode function that takes an Either[A, B] and returns O by encoding
// either the Left or Right value using the appropriate codec.
//
// # Example
//
//	stringCodec := String()
//	intCodec := Int()
//	encoder := encodeEither(stringCodec, intCodec)
//
//	// Encode a Left value
//	leftResult := encoder(either.Left[int]("error"))
//	// leftResult contains the encoded string "error"
//
//	// Encode a Right value
//	rightResult := encoder(either.Right[string](42))
//	// rightResult contains the encoded int 42
//
// # Notes
//
//   - Uses either.Fold to pattern match on the Either value
//   - Left values are encoded using leftItem.Encode
//   - Right values are encoded using rightItem.Encode
func encodeEither[A, B, O, I any](
	leftItem Type[A, O, I],
	rightItem Type[B, O, I],
) Encode[either.Either[A, B], O] {
	return either.Fold(
		leftItem.Encode,
		rightItem.Encode,
	)
}

// validateEither creates a validator for Either[A, B] values.
//
// This function produces a validator that attempts to validate the input as both
// a Left (A) and Right (B) value. The validation strategy is:
//  1. First, try to validate as a Right value (B)
//  2. If Right validation succeeds, return Either.Right[A](B)
//  3. If Right validation fails, try to validate as a Left value (A)
//  4. If Left validation succeeds, return Either.Left[B](A)
//  5. If both validations fail, concatenate all errors from both attempts
//
// This approach ensures that the validator tries both branches and provides
// comprehensive error information when both fail.
//
// # Type Parameters
//
//   - A: The type of the Left value
//   - B: The type of the Right value
//   - O: The output type after encoding (not used in validation)
//   - I: The input type to validate
//
// # Parameters
//
//   - leftItem: The codec for validating Left values of type A
//   - rightItem: The codec for validating Right values of type B
//
// # Returns
//
// A Validate function that takes an input I and returns a Decode function.
// The Decode function takes a Context and returns a Validation[Either[A, B]].
//
// # Validation Logic
//
// The validator follows this decision tree:
//
//	Input I
//	  |
//	  +--> Validate as Right (B)
//	         |
//	         +-- Success --> Return Either.Right[A](B)
//	         |
//	         +-- Failure --> Validate as Left (A)
//	                           |
//	                           +-- Success --> Return Either.Left[B](A)
//	                           |
//	                           +-- Failure --> Return all errors (Left + Right)
//
// # Example
//
//	stringCodec := String()
//	intCodec := Int()
//	validator := validateEither(stringCodec, intCodec)
//
//	// Validate a string (will succeed as Left)
//	result1 := validator("hello")(validation.Context{})
//	// result1 is Success(Either.Left[int]("hello"))
//
//	// Validate an int (will succeed as Right)
//	result2 := validator(42)(validation.Context{})
//	// result2 is Success(Either.Right[string](42))
//
//	// Validate something that's neither (will fail with both errors)
//	result3 := validator([]int{1, 2, 3})(validation.Context{})
//	// result3 is Failure with errors from both string and int validation
//
// # Notes
//
//   - Prioritizes Right validation over Left validation
//   - Accumulates errors from both branches when both fail
//   - Uses the validation context to provide detailed error messages
//   - The validator is lazy: it only evaluates Left if Right fails
func validateEither[A, B, O, I any](
	leftItem Type[A, O, I],
	rightItem Type[B, O, I],
) Validate[I, either.Either[A, B]] {

	valRight := F.Pipe1(
		rightItem.Validate,
		validate.Map[I, B](either.Right[A]),
	)

	valLeft := F.Pipe1(
		leftItem.Validate,
		validate.Map[I, A](either.Left[B]),
	)

	return F.Pipe1(
		valRight,
		validate.Alt(lazy.Of(valLeft)),
	)
}

// Either creates a codec for Either[A, B] values.
//
// This function constructs a complete codec that can encode, decode, and validate
// Either values. An Either represents a value that can be one of two types: Left (A)
// or Right (B). This is commonly used for error handling, where Left represents an
// error and Right represents a success value.
//
// The codec handles both branches of the Either type using the provided codecs for
// each branch. During validation, it attempts to validate the input as both types
// and succeeds if either validation passes.
//
// # Type Parameters
//
//   - A: The type of the Left value
//   - B: The type of the Right value
//   - O: The output type after encoding
//   - I: The input type for validation
//
// # Parameters
//
//   - leftItem: The codec for handling Left values of type A
//   - rightItem: The codec for handling Right values of type B
//
// # Returns
//
// A Type[either.Either[A, B], O, I] that can encode, decode, and validate Either values.
//
// # Codec Behavior
//
// Encoding:
//   - Left values are encoded using leftItem.Encode
//   - Right values are encoded using rightItem.Encode
//
// Validation:
//   - First attempts to validate as Right (B)
//   - If Right fails, attempts to validate as Left (A)
//   - If both fail, returns all accumulated errors
//   - If either succeeds, returns the corresponding Either value
//
// Type Checking:
//   - Uses Is[either.Either[A, B]]() to verify the value is an Either
//
// Naming:
//   - The codec name is "Either[<leftName>, <rightName>]"
//   - Example: "Either[string, int]"
//
// # Example
//
//	// Create a codec for Either[string, int]
//	stringCodec := String()
//	intCodec := Int()
//	eitherCodec := Either(stringCodec, intCodec)
//
//	// Encode a Left value
//	leftEncoded := eitherCodec.Encode(either.Left[int]("error"))
//	// leftEncoded contains the encoded string
//
//	// Encode a Right value
//	rightEncoded := eitherCodec.Encode(either.Right[string](42))
//	// rightEncoded contains the encoded int
//
//	// Decode/validate an input
//	result := eitherCodec.Decode("hello")
//	// result is Success(Either.Left[int]("hello"))
//
//	result2 := eitherCodec.Decode(42)
//	// result2 is Success(Either.Right[string](42))
//
//	// Get the codec name
//	name := eitherCodec.Name()
//	// name is "Either[string, int]"
//
// # Use Cases
//
//   - Error handling: Either[Error, Value]
//   - Alternative values: Either[DefaultValue, CustomValue]
//   - Union types: Either[TypeA, TypeB]
//   - Validation results: Either[ValidationError, ValidatedValue]
//
// # Notes
//
//   - The codec prioritizes Right validation over Left validation
//   - Both branches must have compatible encoding output types (O)
//   - Both branches must have compatible validation input types (I)
//   - The codec name includes the names of both branch codecs
//   - This is a building block for more complex sum types
func Either[A, B, O, I any](
	leftItem Type[A, O, I],
	rightItem Type[B, O, I],
) Type[either.Either[A, B], O, I] {
	return MakeType(
		fmt.Sprintf("Either[%s, %s]", leftItem.Name(), rightItem.Name()),
		Is[either.Either[A, B]](),
		validateEither(leftItem, rightItem),
		encodeEither(leftItem, rightItem),
	)
}
