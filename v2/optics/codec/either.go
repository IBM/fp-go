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

// encodeEither creates an encoder for Either[L, R] values by folding over the
// two branch encoders.
//
// Type Parameters:
//   - R: The Right value type
//   - L: The Left value type
//   - O: The common output type
//   - I: The input type (unused in encoding, carried for type consistency)
//
// Parameters:
//   - leftItem: codec whose Encode is applied to Left values
//   - rightItem: codec whose Encode is applied to Right values
//
// Returns:
//   - Encode[Either[L, R], O]: a function that dispatches to the correct encoder
func encodeEither[R, L, O, I any](
	leftItem Type[L, O, I],
	rightItem Type[R, O, I],
) Encode[either.Either[L, R], O] {
	return either.Fold(
		leftItem.Encode,
		rightItem.Encode,
	)
}

// validateEither creates a validator for Either[L, R] values.
//
// The strategy tries the Right branch first; if it succeeds the result is
// wrapped in either.Right.  If the Right branch fails, the Left branch is
// tried; if it succeeds the result is wrapped in either.Left.  If both
// branches fail, errors from both attempts are accumulated.
//
// Type Parameters:
//   - R: The Right value type
//   - L: The Left value type
//   - O: The output type (unused in validation, carried for type consistency)
//   - I: The input type to validate
//
// Parameters:
//   - leftItem: codec used to validate the Left branch
//   - rightItem: codec used to validate the Right branch (tried first)
//
// Returns:
//   - Validate[I, Either[L, R]]: context-aware validator that returns Either[L, R]
//
// See Also:
//   - AltW: the public codec combinator built on top of this function
func validateEither[R, L, O, I any](
	leftItem Type[L, O, I],
	rightItem Type[R, O, I],
) Validate[I, either.Either[L, R]] {

	valRight := F.Pipe1(
		rightItem.Validate,
		validate.Map[I, R](either.Right[L]),
	)

	valLeft := F.Pipe1(
		leftItem.Validate,
		validate.Map[I, L](either.Left[R]),
	)

	return F.Pipe1(
		valRight,
		validate.Alt(lazy.Of(valLeft)),
	)
}

// AltW lifts two codecs of different decoded types into a single codec whose
// decoded type is Either[L, R].  The "W" suffix signals widening: the result
// type Either[L, R] is strictly wider than either branch alone.
//
// AltW is the widening counterpart of Alt.  Where Alt combines two
// Type[A, O, I] codecs (same decoded type, same encoder used), AltW combines
// a Type[R, O, I] and a Type[L, O, I] whose decoded types differ, producing a
// Type[Either[L, R], O, I].
//
// When decoding input I:
//  1. The Right branch (rightItem) is tried first.
//  2. If it succeeds the value is wrapped in either.Right.
//  3. If it fails the Left branch (leftItem) is tried.
//  4. If the Left branch succeeds the value is wrapped in either.Left.
//  5. If both fail, errors from both branches are accumulated.
//
// When encoding Either[L, R]:
//   - Left(l)  is encoded with leftItem.Encode.
//   - Right(r) is encoded with rightItem.Encode.
//
// The resulting codec is named "AltW[<leftItem>, <rightItem>]".
//
// Type Parameters:
//   - R: The Right decoded type (first explicit type argument)
//   - L: The Left decoded type
//   - O: The common output type for both branches
//   - I: The common input type for both branches
//
// Parameters:
//   - leftItem: A Type[L, O, I] codec for the Left branch
//
// Returns:
//   - Operator[R, Either[L, R], O, I]: an operator that accepts the Right
//     branch codec and produces the combined Either codec
//
// See Also:
//   - Alt: non-widening alternative that keeps the same decoded type
//   - EitherOf: boolean-predicate-gated Either codec
func AltW[R, L, O, I any](
	leftItem Type[L, O, I],
) Operator[R, either.Either[L, R], O, I] {
	return func(rightItem Type[R, O, I]) Type[either.Either[L, R], O, I] {
		return MakeType(
			fmt.Sprintf("AltW[%s, %s]", leftItem, rightItem),
			Is[either.Either[L, R]](),
			validateEither[R, L, O, I](leftItem, rightItem),
			encodeEither[R, L, O, I](leftItem, rightItem),
		)
	}
}
