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

package io

import (
	F "github.com/IBM/fp-go/v2/function"
	INTA "github.com/IBM/fp-go/v2/internal/array"
	INTI "github.com/IBM/fp-go/v2/internal/iter"
	INTR "github.com/IBM/fp-go/v2/internal/record"
)

// MonadTraverseArray applies an IO-returning function to each element of an array
// and collects the results into an IO of an array. Executes in parallel by default.
//
// Example:
//
//	fetchUsers := func(id int) io.IO[User] { return fetchUser(id) }
//	users := io.MonadTraverseArray([]int{1, 2, 3}, fetchUsers)
//	result := users() // []User with all fetched users
func MonadTraverseArray[A, B any](tas []A, f Kleisli[A, B]) IO[[]B] {
	return INTA.MonadTraverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		tas,
		f,
	)
}

// TraverseArray returns a function that applies an IO-returning function to each element
// of an array and collects the results. This is the curried version of MonadTraverseArray.
// Executes in parallel by default.
//
// Example:
//
//	fetchUsers := io.TraverseArray(func(id int) io.IO[User] {
//	    return fetchUser(id)
//	})
//	users := fetchUsers([]int{1, 2, 3})
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return INTA.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// TraverseIter applies an IO-returning function to each element of an iterator sequence
// and collects the results into an IO of an iterator sequence. Executes in parallel by default.
//
// This function is useful for processing lazy sequences where each element requires an IO operation.
// The resulting iterator is also lazy and will only execute IO operations when iterated.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that takes an element of type A and returns an IO computation producing B
//
// Returns:
//   - A function that takes an iterator sequence of A and returns an IO of an iterator sequence of B
//
// Example:
//
//	// Fetch user data for each ID in a sequence
//	fetchUser := func(id int) io.IO[User] {
//	    return func() User {
//	        // Simulate fetching user from database
//	        return User{ID: id, Name: fmt.Sprintf("User%d", id)}
//	    }
//	}
//
//	// Create an iterator of user IDs
//	userIDs := func(yield func(int) bool) {
//	    for _, id := range []int{1, 2, 3, 4, 5} {
//	        if !yield(id) { return }
//	    }
//	}
//
//	// Traverse the iterator, fetching each user
//	fetchUsers := io.TraverseIter(fetchUser)
//	usersIO := fetchUsers(userIDs)
//
//	// Execute the IO to get the iterator of users
//	users := usersIO()
//	for user := range users {
//	    fmt.Printf("User: %v\n", user)
//	}
func TraverseIter[A, B any](f Kleisli[A, B]) Kleisli[Seq[A], Seq[B]] {
	return INTI.Traverse[Seq[A]](
		Map[B],

		Of[Seq[B]],
		Map[Seq[B]],
		Ap[Seq[B]],

		f,
	)
}

// SequenceIter converts an iterator sequence of IO computations into an IO of an iterator sequence of results.
// All computations are executed in parallel by default when the resulting IO is invoked.
//
// This is a special case of TraverseIter where the transformation function is the identity.
// It "flips" the nesting of the iterator and IO types, executing all IO operations and collecting
// their results into a lazy iterator.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - as: An iterator sequence where each element is an IO computation
//
// Returns:
//   - An IO computation that, when executed, produces an iterator sequence of results
//
// Example:
//
//	// Create an iterator of IO operations
//	operations := func(yield func(io.IO[int]) bool) {
//	    yield(func() int { return 1 })
//	    yield(func() int { return 2 })
//	    yield(func() int { return 3 })
//	}
//
//	// Sequence the operations
//	resultsIO := io.SequenceIter(operations)
//
//	// Execute all IO operations and get the iterator of results
//	results := resultsIO()
//	for result := range results {
//	    fmt.Printf("Result: %d\n", result)
//	}
//
// Note: The IO operations are executed when resultsIO() is called, not when iterating
// over the results. The resulting iterator is lazy but the computations have already
// been performed.
func SequenceIter[A any](as Seq[IO[A]]) IO[Seq[A]] {
	return INTI.MonadSequence(
		Map(INTI.Of[Seq[A]]),
		ApplicativeMonoid(INTI.Monoid[Seq[A]]()),
		as,
	)
}

// TraverseArrayWithIndex is like TraverseArray but the function also receives the index.
// Executes in parallel by default.
//
// Example:
//
//	numbered := io.TraverseArrayWithIndex(func(i int, s string) io.IO[string] {
//	    return io.Of(fmt.Sprintf("%d: %s", i, s))
//	})
func TraverseArrayWithIndex[A, B any](f func(int, A) IO[B]) Kleisli[[]A, []B] {
	return INTA.TraverseWithIndex[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// SequenceArray converts an array of IO computations into an IO of an array of results.
// All computations are executed in parallel by default.
//
// Example:
//
//	operations := []io.IO[int]{fetchA(), fetchB(), fetchC()}
//	results := io.SequenceArray(operations)
//	values := results() // []int with all results
func SequenceArray[A any](tas []IO[A]) IO[[]A] {
	return MonadTraverseArray(tas, F.Identity[IO[A]])
}

// MonadTraverseRecord applies an IO-returning function to each value in a map
// and collects the results into an IO of a map. Executes in parallel by default.
//
// Example:
//
//	fetchData := func(url string) io.IO[Data] { return fetch(url) }
//	urls := map[string]string{"a": "http://a.com", "b": "http://b.com"}
//	data := io.MonadTraverseRecord(urls, fetchData)
func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f Kleisli[A, B]) IO[map[K]B] {
	return INTR.MonadTraverse(
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		tas,
		f,
	)
}

// TraverseRecord returns a function that applies an IO-returning function to each value
// in a map and collects the results. This is the curried version of MonadTraverseRecord.
// Executes in parallel by default.
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return INTR.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// TraverseRecordWithIndex is like TraverseRecord but the function also receives the key.
// Executes in parallel by default.
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) IO[B]) Kleisli[map[K]A, map[K]B] {
	return INTR.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// SequenceRecord converts a map of IO computations into an IO of a map of results.
// All computations are executed in parallel by default.
//
// Example:
//
//	operations := map[string]io.IO[int]{"a": fetchA(), "b": fetchB()}
//	results := io.SequenceRecord(operations)
//	values := results() // map[string]int with all results
func SequenceRecord[K comparable, A any](tas map[K]IO[A]) IO[map[K]A] {
	return MonadTraverseRecord(tas, F.Identity[IO[A]])
}

// MonadTraverseArraySeq applies an IO-returning function to each element of an array
// and collects the results into an IO of an array. Executes sequentially (one after another).
//
// Example:
//
//	fetchUsers := func(id int) io.IO[User] { return fetchUser(id) }
//	users := io.MonadTraverseArraySeq([]int{1, 2, 3}, fetchUsers)
func MonadTraverseArraySeq[A, B any](tas []A, f Kleisli[A, B]) IO[[]B] {
	return INTA.MonadTraverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		tas,
		f,
	)
}

// TraverseArraySeq returns a function that applies an IO-returning function to each element
// of an array and collects the results. Executes sequentially (one after another).
// Use this when operations must be performed in order or when parallel execution is not desired.
func TraverseArraySeq[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return INTA.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		f,
	)
}

// TraverseArrayWithIndexSeq is like TraverseArraySeq but the function also receives the index.
// Executes sequentially (one after another).
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) IO[B]) Kleisli[[]A, []B] {
	return INTA.TraverseWithIndex[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		f,
	)
}

// SequenceArraySeq converts an array of IO computations into an IO of an array of results.
// All computations are executed sequentially (one after another).
func SequenceArraySeq[A any](tas []IO[A]) IO[[]A] {
	return MonadTraverseArraySeq(tas, F.Identity[IO[A]])
}

// MonadTraverseRecordSeq applies an IO-returning function to each value in a map
// and collects the results into an IO of a map. Executes sequentially.
func MonadTraverseRecordSeq[K comparable, A, B any](tas map[K]A, f Kleisli[A, B]) IO[map[K]B] {
	return INTR.MonadTraverse(
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		tas,
		f,
	)
}

// TraverseRecordSeq returns a function that applies an IO-returning function to each value
// in a map and collects the results. Executes sequentially (one after another).
func TraverseRecordSeq[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return INTR.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		f,
	)
}

// TraverseRecordWithIndeSeq is like TraverseRecordSeq but the function also receives the key.
// Executes sequentially (one after another).
// Note: There's a typo in the function name (Inde instead of Index) for backward compatibility.
func TraverseRecordWithIndeSeq[K comparable, A, B any](f func(K, A) IO[B]) Kleisli[map[K]A, map[K]B] {
	return INTR.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		f,
	)
}

// SequenceRecordSeq converts a map of IO computations into an IO of a map of results.
// All computations are executed sequentially (one after another).
func SequenceRecordSeq[K comparable, A any](tas map[K]IO[A]) IO[map[K]A] {
	return MonadTraverseRecordSeq(tas, F.Identity[IO[A]])
}
