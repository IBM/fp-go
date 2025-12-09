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
	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
)

// TraverseArray applies a ReaderResult-returning function to each element of an array,
// collecting the results. If any element fails, the entire operation fails with the first error.
//
// Example:
//
//	parseUser := func(id int) readerresult.ReaderResult[DB, User] { ... }
//	ids := []int{1, 2, 3}
//	result := readerresult.TraverseArray[DB](parseUser)(ids)
//	// result(db) returns ([]User, nil) with all users or (nil, error) on first error
//
//go:inline
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return RR.TraverseArray(f)
}

//go:inline
func MonadTraverseArray[A, B any](as []A, f Kleisli[A, B]) ReaderResult[[]B] {
	return RR.MonadTraverseArray(as, f)
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
func TraverseArrayWithIndex[A, B any](f func(int, A) ReaderResult[B]) Kleisli[[]A, []B] {
	return RR.TraverseArrayWithIndex(f)
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
//	// result(cfg) returns ([]int{1, 2, 3}, nil)
//
//go:inline
func SequenceArray[A any](ma []ReaderResult[A]) ReaderResult[[]A] {
	return RR.SequenceArray(ma)
}
