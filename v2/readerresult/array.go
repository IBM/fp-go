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

package readerresult

import (
	G "github.com/IBM/fp-go/v2/readereither/generic"
)

// TraverseArray applies a ReaderResult-returning function to each element of an array,
// collecting the results. If any element fails, the entire operation fails with the first error.
//
// Example:
//
//	parseUser := func(id int) readerresult.ReaderResult[DB, User] { ... }
//	ids := []int{1, 2, 3}
//	result := readerresult.TraverseArray[DB](parseUser)(ids)
//	// result(db) returns result.Result[[]User] with all users or first error
//
//go:inline
func TraverseArray[L, A, B any](f Kleisli[L, A, B]) Kleisli[L, []A, []B] {
	return G.TraverseArray[ReaderResult[L, B], ReaderResult[L, []B], []A](f)
}

// TraverseArrayWithIndex is like TraverseArray but the function also receives the element's index.
// This is useful when the transformation depends on the position in the array.
//
// Example:
//
//	processItem := func(idx int, item string) readerresult.ReaderResult[Config, int] {
//	    return readerresult.Of[Config](idx + len(item))
//	}
//	items := []string{"a", "bb", "ccc"}
//	result := readerresult.TraverseArrayWithIndex[Config](processItem)(items)
//
//go:inline
func TraverseArrayWithIndex[L, A, B any](f func(int, A) ReaderResult[L, B]) Kleisli[L, []A, []B] {
	return G.TraverseArrayWithIndex[ReaderResult[L, B], ReaderResult[L, []B], []A](f)
}

// SequenceArray converts an array of ReaderResult values into a single ReaderResult of an array.
// If any element fails, the entire operation fails with the first error encountered.
// All computations share the same environment.
//
// Example:
//
//	readers := []readerresult.ReaderResult[Config, int]{
//	    readerresult.Of[Config](1),
//	    readerresult.Of[Config](2),
//	    readerresult.Of[Config](3),
//	}
//	result := readerresult.SequenceArray(readers)
//	// result(cfg) returns result.Of([]int{1, 2, 3})
//
//go:inline
func SequenceArray[L, A any](ma []ReaderResult[L, A]) ReaderResult[L, []A] {
	return G.SequenceArray[ReaderResult[L, A], ReaderResult[L, []A]](ma)
}
