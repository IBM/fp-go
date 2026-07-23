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
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/semigroup"
)

// Optional lifts a codec Type[A, O, I] into a codec Type[Option[A], O, I] by
// gating decoding on a boolean predicate codec.
//
// When decoding input I:
//  1. The predicate pred is validated against the input first.
//  2. If pred decodes to true, the inner codec onSome is invoked and its result
//     is wrapped in option.Some.
//  3. If pred decodes to false (or the input does not match), the result is
//     option.None and the monoid empty value is used for the output.
//
// When encoding Option[A]:
//  1. If the option is Some(a), a is encoded with onSome.Encode and the boolean
//     true is encoded with pred.Encode; the two outputs are combined with the
//     monoid.
//  2. If the option is None, the monoid empty value is used for the inner
//     encoding and false is encoded with pred.Encode; the two outputs are
//     combined with the monoid.
//
// The resulting codec is named "Optional[<pred> x <onSome>]".
//
// Type Parameters:
//   - A: The type decoded by the inner codec onSome
//   - O: The output type produced by both the predicate and the inner codec
//   - I: The input type consumed by both the predicate and the inner codec
//
// Parameters:
//   - m: A Monoid[O] used to combine the encoded predicate output and the
//     encoded value output
//   - pred: A Type[bool, O, I] that decodes the presence flag (bool) and
//     encodes it back to O
//
// Returns:
//   - An Operator[A, Option[A], O, I] that transforms a Type[A, O, I] into a
//     Type[Option[A], O, I]
//
// See Also:
//   - ApSO: Applicative sequencing for optional struct fields via Optional optic
//   - Do: Entry point for do-notation style codec construction
func Optional[A, O, I any](
	m Monoid[O],
	onSome Type[A, O, I]) Operator[bool, Option[A], O, I] {

	merge := semigroup.AppendTo(m)
	orElse := F.Pipe2(
		option.None[A],
		lazy.Map(validate.Of[I, Option[A]]),
		option.GetOrElse,
	)

	return func(pred Type[bool, O, I]) Type[Option[A], O, I] {

		return MakeType(
			fmt.Sprintf("Optional[%s x %s]", pred, onSome),
			Is[Option[A]](),
			F.Pipe1(
				pred.Validate,
				validate.Chain(F.Flow3(
					option.FromPredicate(reader.Ask[bool]()),
					option.MapTo[bool](F.Pipe1(
						onSome.Validate,
						validate.Map[I](option.Some[A]),
					)),
					orElse,
				)),
			),
			F.Pipe1(
				F.Flow2(
					option.Map(onSome.Encode),
					option.GetOrElse(m.Empty),
				),
				reader.ApS(merge, F.Flow2(
					option.IsSome[A],
					pred.Encode,
				)),
			),
		)
	}
}

// EitherOf lifts two codecs — one for L and one for R — into a codec
// Type[Either[L, R], O, I] by dispatching on a boolean predicate codec.
//
// When decoding input I:
//  1. The predicate pred is validated against the input first.
//  2. If pred decodes to true, the right codec onRight is invoked and its
//     result is wrapped in either.Right.
//  3. If pred decodes to false, the left codec onLeft is invoked and its
//     result is wrapped in either.Left.
//
// When encoding Either[L, R]:
//  1. If the value is Right(r), r is encoded with onRight.Encode and true is
//     encoded with pred.Encode; the two outputs are combined with the monoid.
//  2. If the value is Left(l), l is encoded with onLeft.Encode and false is
//     encoded with pred.Encode; the two outputs are combined with the monoid.
//
// The resulting codec is named "EitherOf[<pred> x (<onLeft>|<onRight>)]".
//
// Type Parameters:
//   - L: The type decoded by the left codec onLeft
//   - R: The type decoded by the right codec onRight
//   - O: The output type produced by the predicate and both branch codecs
//   - I: The input type consumed by the predicate and both branch codecs
//
// Parameters:
//   - m: A Monoid[O] used to combine the encoded predicate output and the
//     encoded branch output
//   - onLeft: A Type[L, O, I] that decodes and encodes the Left branch
//   - onRight: A Type[R, O, I] that decodes and encodes the Right branch
//
// Returns:
//   - An Operator[bool, Either[L, R], O, I] that transforms a
//     Type[bool, O, I] predicate codec into a Type[Either[L, R], O, I]
//
// See Also:
//   - Optional: Boolean-gated codec that lifts a value into Option
//   - Either: Untagged either codec that tries both branches by structure
func EitherOf[L, R, O, I any](
	m Monoid[O],
	onLeft Type[L, O, I],
	onRight Type[R, O, I],
) Operator[bool, either.Either[L, R], O, I] {

	merge := semigroup.AppendTo(m)

	return func(pred Type[bool, O, I]) Type[either.Either[L, R], O, I] {

		return MakeType(
			fmt.Sprintf("EitherOf[%s x (%s|%s)]", pred, onLeft, onRight),
			Is[either.Either[L, R]](),
			F.Pipe1(
				pred.Validate,
				validate.Chain(F.Flow3(
					option.FromPredicate(reader.Ask[bool]()),
					option.MapTo[bool](F.Pipe1(
						onRight.Validate,
						validate.Map[I](either.Right[L, R]),
					)),
					option.GetOrElse(lazy.Of(F.Pipe1(
						onLeft.Validate,
						validate.Map[I](either.Left[R, L]),
					))),
				)),
			),
			F.Pipe1(
				either.Fold(onLeft.Encode, onRight.Encode),
				reader.ApS(merge, F.Flow2(
					either.IsRight[L, R],
					pred.Encode,
				)),
			),
		)
	}
}
