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

package readereither

import (
	G "github.com/IBM/fp-go/v2/readereither/generic"
)

// TraverseArray transforms each element of an array using a function that returns a ReaderEither,
// then collects the results into a single ReaderEither containing an array.
//
// If any transformation fails, the entire operation fails with the first error encountered.
// All transformations are executed sequentially.
//
// Type parameters:
//   - E: The context type
//   - L: The error type
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that transforms each element into a ReaderEither
//
// Returns:
//
//	A function that takes an array and returns a ReaderEither of an array
//
// Example:
//
//	fetchUsers := TraverseArray(func(id int) ReaderEither[Config, error, User] {
//	    return fetchUser(id)
//	})
//	result := fetchUsers([]int{1, 2, 3})
//	// result(cfg) returns Right([user1, user2, user3]) or Left(error)
//
//go:inline
func TraverseArray[E, L, A, B any](f func(A) ReaderEither[E, L, B]) func([]A) ReaderEither[E, L, []B] {
	return G.TraverseArray[ReaderEither[E, L, B], ReaderEither[E, L, []B], []A](f)
}

// TraverseArrayWithIndex is like TraverseArray but the transformation function also receives the index.
//
// This is useful when the transformation depends on the element's position in the array.
//
// Type parameters:
//   - E: The context type
//   - L: The error type
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that transforms each element and its index into a ReaderEither
//
// Returns:
//
//	A function that takes an array and returns a ReaderEither of an array
//
// Example:
//
//	processWithIndex := TraverseArrayWithIndex(func(i int, val string) ReaderEither[Config, error, string] {
//	    return Of[Config, error](fmt.Sprintf("%d: %s", i, val))
//	})
//
//go:inline
func TraverseArrayWithIndex[E, L, A, B any](f func(int, A) ReaderEither[E, L, B]) func([]A) ReaderEither[E, L, []B] {
	return G.TraverseArrayWithIndex[ReaderEither[E, L, B], ReaderEither[E, L, []B], []A](f)
}

// SequenceArray converts an array of ReaderEither into a ReaderEither of an array.
//
// This is useful when you have multiple independent computations and want to execute them all
// and collect their results. If any computation fails, the entire operation fails with the first error.
//
// Type parameters:
//   - E: The context type
//   - L: The error type
//   - A: The element type
//
// Parameters:
//   - ma: An array of ReaderEither computations
//
// Returns:
//
//	A ReaderEither that produces an array of results
//
// Example:
//
//	computations := []ReaderEither[Config, error, int]{
//	    fetchCount("users"),
//	    fetchCount("posts"),
//	    fetchCount("comments"),
//	}
//	result := SequenceArray(computations)
//	// result(cfg) returns Right([userCount, postCount, commentCount]) or Left(error)
//
//go:inline
func SequenceArray[E, L, A any](ma []ReaderEither[E, L, A]) ReaderEither[E, L, []A] {
	return G.SequenceArray[ReaderEither[E, L, A], ReaderEither[E, L, []A]](ma)
}
