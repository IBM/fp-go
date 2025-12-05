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

package readerio

import (
	G "github.com/IBM/fp-go/v2/readerio/generic"
)

// TraverseArray transforms each element of an array using a function that returns a ReaderIO,
// then collects the results into a single ReaderIO containing an array.
//
// All transformations are executed sequentially.
//
// Type parameters:
//   - R: The context type
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that transforms each element into a ReaderIO
//
// Returns:
//
//	A function that takes an array and returns a ReaderIO of an array
//
// Example:
//
//	fetchUsers := TraverseArray(func(id int) ReaderIO[Config, User] {
//	    return fetchUser(id)
//	})
//	result := fetchUsers([]int{1, 2, 3})
//	// result(cfg)() returns [user1, user2, user3]
//
//go:inline
func TraverseArray[R, A, B any](f func(A) ReaderIO[R, B]) func([]A) ReaderIO[R, []B] {
	return G.TraverseArray[ReaderIO[R, B], ReaderIO[R, []B], IO[B], IO[[]B], []A](f)
}

// TraverseArrayWithIndex is like TraverseArray but the transformation function also receives the index.
//
// This is useful when the transformation depends on the element's position in the array.
//
// Type parameters:
//   - R: The context type
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that transforms each element and its index into a ReaderIO
//
// Returns:
//
//	A function that takes an array and returns a ReaderIO of an array
//
// Example:
//
//	processWithIndex := TraverseArrayWithIndex(func(i int, val string) ReaderIO[Config, string] {
//	    return Of[Config](fmt.Sprintf("%d: %s", i, val))
//	})
//
//go:inline
func TraverseArrayWithIndex[R, A, B any](f func(int, A) ReaderIO[R, B]) func([]A) ReaderIO[R, []B] {
	return G.TraverseArrayWithIndex[ReaderIO[R, B], ReaderIO[R, []B], IO[B], IO[[]B], []A](f)
}

// SequenceArray converts an array of ReaderIO into a ReaderIO of an array.
//
// This is useful when you have multiple independent computations and want to execute them all
// and collect their results.
//
// Type parameters:
//   - R: The context type
//   - A: The element type
//
// Parameters:
//   - ma: An array of ReaderIO computations
//
// Returns:
//
//	A ReaderIO that produces an array of results
//
// Example:
//
//	computations := []ReaderIO[Config, int]{
//	    fetchCount("users"),
//	    fetchCount("posts"),
//	    fetchCount("comments"),
//	}
//	result := SequenceArray(computations)
//	// result(cfg)() returns [userCount, postCount, commentCount]
//
//go:inline
func SequenceArray[R, A any](ma []ReaderIO[R, A]) ReaderIO[R, []A] {
	return G.SequenceArray[ReaderIO[R, A], ReaderIO[R, []A]](ma)
}
