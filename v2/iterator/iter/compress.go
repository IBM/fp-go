// Copyright (c) 2025 IBM Corp.
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

package iter

import (
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
)

// Compress filters elements from a sequence based on a corresponding sequence of boolean selectors.
//
// This function takes a sequence of boolean values and returns an operator that filters
// elements from the input sequence. An element is included in the output if and only if
// the corresponding boolean selector is true. The filtering stops when either sequence
// is exhausted.
//
// The implementation works by:
//  1. Zipping the input sequence with the selector sequence
//  2. Converting the Seq2 to a sequence of Pairs
//  3. Filtering to keep only pairs where the boolean (tail) is true
//  4. Extracting the original values (head) from the filtered pairs
//
// RxJS Equivalent: Similar to combining [zip] with [filter] - https://rxjs.dev/api/operators/zip
//
// Type Parameters:
//   - U: The type of elements in the sequence to be filtered
//
// Parameters:
//   - sel: A sequence of boolean values used as selectors
//
// Returns:
//   - An Operator that filters elements based on the selector sequence
//
// Example - Basic filtering:
//
//	data := iter.From(1, 2, 3, 4, 5)
//	selectors := iter.From(true, false, true, false, true)
//	filtered := iter.Compress(selectors)(data)
//	// yields: 1, 3, 5
//
// Example - Shorter selector sequence:
//
//	data := iter.From("a", "b", "c", "d", "e")
//	selectors := iter.From(true, true, false)
//	filtered := iter.Compress(selectors)(data)
//	// yields: "a", "b" (stops when selectors are exhausted)
//
// Example - All false selectors:
//
//	data := iter.From(1, 2, 3)
//	selectors := iter.From(false, false, false)
//	filtered := iter.Compress(selectors)(data)
//	// yields: nothing (empty sequence)
//
// Example - All true selectors:
//
//	data := iter.From(10, 20, 30)
//	selectors := iter.From(true, true, true)
//	filtered := iter.Compress(selectors)(data)
//	// yields: 10, 20, 30 (all elements pass through)
func Compress[U any](sel Seq[bool]) Operator[U, U] {
	return F.Flow3(
		Zip[U](sel),
		ToSeqPair[U, bool],
		FilterMap(F.Flow2(
			O.FromPredicate(P.Tail[U, bool]),
			O.Map(P.Head[U, bool]),
		)),
	)
}
