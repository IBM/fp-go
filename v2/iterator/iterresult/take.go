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

package iterresult

import (
	"github.com/IBM/fp-go/v2/iterator/itereither"
)

// TakeUntilError takes elements from a SeqResult until the first error is encountered, including that error.
//
// This function creates a transformation that yields all success values from the source sequence
// until an error is encountered. When an error is found, it is included in the output and
// the sequence terminates immediately. This is useful for processing sequences that should
// stop at the first error while still capturing that error.
//
// The operation is lazy and only consumes elements from the source sequence as needed.
// Once an error is encountered, iteration stops immediately without consuming the remaining
// elements from the source.
//
// Marble Diagram:
//
//	Input:  --Ok(1)--Ok(2)--Ok(3)--Err(e)--Ok(4)--Ok(5)-->
//	TakeUntilError
//	Output: --Ok(1)--Ok(2)--Ok(3)--Err(e)|
//	                                (includes error, then stops)
//
// Where Ok(x) represents a success Result and Err(e) represents an error Result.
//
// Type Parameters:
//   - T: The success type
//
// Parameters:
//   - s: The input SeqResult to process
//
// Returns:
//   - SeqResult[T]: A sequence containing all success values up to and including the first error
//
// Example - Stop at first error:
//
//	seq := iter.From(
//	    result.Of(1),
//	    result.Of(2),
//	    result.Left(errors.New("error")),
//	    result.Of(3),
//	)
//	result := TakeUntilError(seq)
//	// yields: Ok(1), Ok(2), Err("error")
//	// Note: Ok(3) is not processed
//
// Example - All success values:
//
//	seq := iter.From(
//	    result.Of(1),
//	    result.Of(2),
//	    result.Of(3),
//	)
//	result := TakeUntilError(seq)
//	// yields: Ok(1), Ok(2), Ok(3)
//	// All elements pass through since there's no error
//
// Example - First element is error:
//
//	seq := iter.From(
//	    result.Left(errors.New("immediate error")),
//	    result.Of(1),
//	    result.Of(2),
//	)
//	result := TakeUntilError(seq)
//	// yields: Err("immediate error")
//	// Stops immediately after the first error
//
// Example - Processing with error handling:
//
//	parseNumbers := func(inputs []string) SeqResult[int] {
//	    seq := iter.MonadMap(
//	        iter.From(inputs...),
//	        result.Eitherize1(strconv.Atoi),
//	    )
//	    return TakeUntilError(seq)
//	}
//	// Processes strings until first parse error, including the error
//
// See Also:
//   - itereither.TakeUntilLeft: The underlying function used
//   - StopOnError: Alias for TakeUntilError
func TakeUntilError[T any](s SeqResult[T]) SeqResult[T] {
	return itereither.TakeUntilLeft(s)
}

// StopOnError is an alias for TakeUntilError.
// It takes elements from a SeqResult until the first error is encountered, including that error.
//
// This function provides a more descriptive name for the same operation as TakeUntilError,
// emphasizing that the sequence stops when an error occurs.
//
// Type Parameters:
//   - T: The success type
//
// Parameters:
//   - s: The input SeqResult to process
//
// Returns:
//   - SeqResult[T]: A sequence containing all success values up to and including the first error
//
// Example:
//
//	seq := iter.From(
//	    result.Of(1),
//	    result.Of(2),
//	    result.Left(errors.New("error")),
//	    result.Of(3),
//	)
//	result := StopOnError(seq)
//	// yields: Ok(1), Ok(2), Err("error")
//
// See Also:
//   - TakeUntilError: The primary function this aliases
func StopOnError[T any](s SeqResult[T]) SeqResult[T] {
	return TakeUntilError(s)
}
