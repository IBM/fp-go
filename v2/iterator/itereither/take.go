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

package itereither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/iterator/iter"
)

// TakeUntilLeft takes elements from a SeqEither until the first Left (error) is encountered, including that Left.
//
// This function creates a transformation that yields all Right values from the source sequence
// until a Left value is encountered. When a Left is found, it is included in the output and
// the sequence terminates immediately. This is useful for processing sequences that should
// stop at the first error while still capturing that error.
//
// The operation is lazy and only consumes elements from the source sequence as needed.
// Once a Left is encountered, iteration stops immediately without consuming the remaining
// elements from the source.
//
// Marble Diagram:
//
//	Input:  --R(1)--R(2)--R(3)--L(e)--R(4)--R(5)-->
//	TakeUntilLeft
//	Output: --R(1)--R(2)--R(3)--L(e)|
//	                              (includes Left, then stops)
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
//
// Type Parameters:
//   - E: The error type (Left)
//   - T: The success type (Right)
//
// Parameters:
//   - s: The input SeqEither to process
//
// Returns:
//   - SeqEither[E, T]: A sequence containing all Right values up to and including the first Left
//
// Example - Stop at first error:
//
//	seq := iter.From(
//	    either.Right[string](1),
//	    either.Right[string](2),
//	    either.Left[int]("error"),
//	    either.Right[string](3),
//	)
//	result := TakeUntilLeft(seq)
//	// yields: Right(1), Right(2), Left("error")
//	// Note: Right(3) is not processed
//
// Example - All Right values:
//
//	seq := iter.From(
//	    either.Right[string](1),
//	    either.Right[string](2),
//	    either.Right[string](3),
//	)
//	result := TakeUntilLeft(seq)
//	// yields: Right(1), Right(2), Right(3)
//	// All elements pass through since there's no Left
//
// Example - First element is Left:
//
//	seq := iter.From(
//	    either.Left[int]("immediate error"),
//	    either.Right[string](1),
//	    either.Right[string](2),
//	)
//	result := TakeUntilLeft(seq)
//	// yields: Left("immediate error")
//	// Stops immediately after the first Left
//
// Example - Processing with error handling:
//
//	parseNumbers := func(inputs []string) SeqEither[error, int] {
//	    seq := iter.MonadMap(
//	        iter.FromSlice(inputs),
//	        result.Eitherize1(strconv.Atoi),
//	    )
//	    return TakeUntilLeft(seq)
//	}
//	// Processes strings until first parse error, including the error
//
// See Also:
//   - iter.TakeWhileInclusive: The underlying iterator function used
//   - either.IsRight: The predicate used to identify Right values
func TakeUntilLeft[E, T any](s SeqEither[E, T]) SeqEither[E, T] {
	return iter.TakeWhileInclusive(either.IsRight[E, T])(s)
}

func StopOnLeft[E, T any](s SeqEither[E, T]) SeqEither[E, T] {
	return TakeUntilLeft(s)
}
