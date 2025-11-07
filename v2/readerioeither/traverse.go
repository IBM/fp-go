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

package readerioeither

import (
	G "github.com/IBM/fp-go/v2/readerioeither/generic"
)

// TraverseArray transforms each element of an array using a function that returns a ReaderIOEither,
// then collects the results into a single ReaderIOEither containing an array.
//
// If any transformation fails, the entire operation fails with the first error encountered.
// All transformations are executed sequentially.
//
// Type parameters:
//   - R: The context type
//   - E: The error type
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that transforms each element into a ReaderIOEither
//
// Returns:
//
//	A function that takes an array and returns a ReaderIOEither of an array
//
// Example:
//
//	fetchUsers := TraverseArray(func(id int) ReaderIOEither[Config, error, User] {
//	    return fetchUser(id)
//	})
//	result := fetchUsers([]int{1, 2, 3})
//	// result(cfg)() returns Right([user1, user2, user3]) or Left(error)
//
//go:inline
func TraverseArray[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func([]A) ReaderIOEither[R, E, []B] {
	return G.TraverseArray[ReaderIOEither[R, E, B], ReaderIOEither[R, E, []B], IOEither[E, B], IOEither[E, []B], []A](f)
}

// TraverseArrayWithIndex is like TraverseArray but the transformation function also receives the index.
//
// This is useful when the transformation depends on the element's position in the array.
//
// Type parameters:
//   - R: The context type
//   - E: The error type
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that transforms each element and its index into a ReaderIOEither
//
// Returns:
//
//	A function that takes an array and returns a ReaderIOEither of an array
//
// Example:
//
//	processWithIndex := TraverseArrayWithIndex(func(i int, val string) ReaderIOEither[Config, error, string] {
//	    return Of[Config, error](fmt.Sprintf("%d: %s", i, val))
//	})
//
//go:inline
func TraverseArrayWithIndex[R, E, A, B any](f func(int, A) ReaderIOEither[R, E, B]) func([]A) ReaderIOEither[R, E, []B] {
	return G.TraverseArrayWithIndex[ReaderIOEither[R, E, B], ReaderIOEither[R, E, []B], IOEither[E, B], IOEither[E, []B], []A](f)
}

// SequenceArray converts an array of ReaderIOEither into a ReaderIOEither of an array.
//
// This is useful when you have multiple independent computations and want to execute them all
// and collect their results. If any computation fails, the entire operation fails with the first error.
//
// Type parameters:
//   - R: The context type
//   - E: The error type
//   - A: The element type
//
// Parameters:
//   - ma: An array of ReaderIOEither computations
//
// Returns:
//
//	A ReaderIOEither that produces an array of results
//
// Example:
//
//	computations := []ReaderIOEither[Config, error, int]{
//	    fetchCount("users"),
//	    fetchCount("posts"),
//	    fetchCount("comments"),
//	}
//	result := SequenceArray(computations)
//	// result(cfg)() returns Right([userCount, postCount, commentCount]) or Left(error)
//
//go:inline
func SequenceArray[R, E, A any](ma []ReaderIOEither[R, E, A]) ReaderIOEither[R, E, []A] {
	return G.SequenceArray[ReaderIOEither[R, E, A], ReaderIOEither[R, E, []A]](ma)
}

// TraverseRecord transforms each value in a map using a function that returns a ReaderIOEither,
// then collects the results into a single ReaderIOEither containing a map.
//
// If any transformation fails, the entire operation fails with the first error encountered.
// The keys are preserved in the output map.
//
// Type parameters:
//   - R: The context type
//   - K: The key type (must be comparable)
//   - E: The error type
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A function that transforms each value into a ReaderIOEither
//
// Returns:
//
//	A function that takes a map and returns a ReaderIOEither of a map
//
// Example:
//
//	enrichUsers := TraverseRecord(func(user User) ReaderIOEither[Config, error, EnrichedUser] {
//	    return enrichUser(user)
//	})
//	result := enrichUsers(map[string]User{"alice": user1, "bob": user2})
//
//go:inline
func TraverseRecord[R any, K comparable, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func(map[K]A) ReaderIOEither[R, E, map[K]B] {
	return G.TraverseRecord[ReaderIOEither[R, E, B], ReaderIOEither[R, E, map[K]B], IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// TraverseRecordWithIndex is like TraverseRecord but the transformation function also receives the key.
//
// This is useful when the transformation depends on the key associated with each value.
//
// Type parameters:
//   - R: The context type
//   - K: The key type (must be comparable)
//   - E: The error type
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A function that transforms each key-value pair into a ReaderIOEither
//
// Returns:
//
//	A function that takes a map and returns a ReaderIOEither of a map
//
// Example:
//
//	processWithKey := TraverseRecordWithIndex(func(key string, val int) ReaderIOEither[Config, error, string] {
//	    return Of[Config, error](fmt.Sprintf("%s: %d", key, val))
//	})
//
//go:inline
func TraverseRecordWithIndex[R any, K comparable, E, A, B any](f func(K, A) ReaderIOEither[R, E, B]) func(map[K]A) ReaderIOEither[R, E, map[K]B] {
	return G.TraverseRecordWithIndex[ReaderIOEither[R, E, B], ReaderIOEither[R, E, map[K]B], IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// SequenceRecord converts a map of ReaderIOEither into a ReaderIOEither of a map.
//
// This is useful when you have multiple independent computations keyed by some identifier
// and want to execute them all and collect their results. If any computation fails,
// the entire operation fails with the first error.
//
// Type parameters:
//   - R: The context type
//   - K: The key type (must be comparable)
//   - E: The error type
//   - A: The value type
//
// Parameters:
//   - ma: A map of ReaderIOEither computations
//
// Returns:
//
//	A ReaderIOEither that produces a map of results
//
// Example:
//
//	computations := map[string]ReaderIOEither[Config, error, int]{
//	    "users": fetchCount("users"),
//	    "posts": fetchCount("posts"),
//	}
//	result := SequenceRecord(computations)
//	// result(cfg)() returns Right(map[string]int{"users": 100, "posts": 50}) or Left(error)
//
//go:inline
func SequenceRecord[R any, K comparable, E, A any](ma map[K]ReaderIOEither[R, E, A]) ReaderIOEither[R, E, map[K]A] {
	return G.SequenceRecord[ReaderIOEither[R, E, A], ReaderIOEither[R, E, map[K]A]](ma)
}
