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

package reader

import (
	INTI "github.com/IBM/fp-go/v2/internal/iter"
)

// TraverseIter traverses an iterator sequence, applying a Reader-producing function to each element
// and collecting the results in a Reader that produces an iterator.
//
// This function transforms a sequence of values through a function that produces Readers,
// then "flips" the nesting so that instead of having an iterator of Readers, you get a
// single Reader that produces an iterator of values. All Readers share the same environment R.
//
// This is particularly useful when you have a collection of values that each need to be
// transformed using environment-dependent logic, and you want to defer the environment
// injection until the final execution.
//
// Type Parameters:
//   - R: The shared environment/context type
//   - A: The input element type in the iterator
//   - B: The output element type in the resulting iterator
//
// Parameters:
//   - f: A Kleisli arrow that transforms each element A into a Reader[R, B]
//
// Returns:
//   - A Kleisli arrow that takes an iterator of A and returns a Reader producing an iterator of B
//
// Example:
//
//	type Config struct { Multiplier int }
//
//	// Function that creates a Reader for each number
//	multiplyByConfig := func(x int) reader.Reader[Config, int] {
//	    return func(c Config) int { return x * c.Multiplier }
//	}
//
//	// Create an iterator of numbers
//	numbers := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	    yield(3)
//	}
//
//	// Traverse the iterator
//	traversed := reader.TraverseIter(multiplyByConfig)(numbers)
//
//	// Execute with config
//	result := traversed(Config{Multiplier: 10})
//	// result is an iterator that yields: 10, 20, 30
func TraverseIter[R, A, B any](f Kleisli[R, A, B]) Kleisli[R, Seq[A], Seq[B]] {
	return INTI.Traverse[Seq[A]](
		Map[R, B],

		Of[R, Seq[B]],
		Map[R, Seq[B]],
		Ap[Seq[B]],

		f,
	)
}

// SequenceIter sequences an iterator of Readers into a Reader that produces an iterator.
//
// This function "flips" the nesting of an iterator and Reader types. Given an iterator
// where each element is a Reader[R, A], it produces a single Reader[R, Seq[A]] that,
// when executed with an environment, evaluates all the Readers with that environment
// and collects their results into an iterator.
//
// This is a special case of TraverseIter where the transformation function is the identity.
// All Readers in the input iterator share the same environment R and are evaluated with it.
//
// Type Parameters:
//   - R: The shared environment/context type
//   - A: The result type produced by each Reader
//
// Parameters:
//   - as: An iterator sequence where each element is a Reader[R, A]
//
// Returns:
//   - A Reader that, when executed, produces an iterator of all the Reader results
//
// Example:
//
//	type Config struct { Base int }
//
//	// Create an iterator of Readers
//	readers := func(yield func(reader.Reader[Config, int]) bool) {
//	    yield(func(c Config) int { return c.Base + 1 })
//	    yield(func(c Config) int { return c.Base + 2 })
//	    yield(func(c Config) int { return c.Base + 3 })
//	}
//
//	// Sequence the iterator
//	sequenced := reader.SequenceIter(readers)
//
//	// Execute with config
//	result := sequenced(Config{Base: 10})
//	// result is an iterator that yields: 11, 12, 13
func SequenceIter[R, A any](as Seq[Reader[R, A]]) Reader[R, Seq[A]] {
	return INTI.MonadSequence(
		Map[R](INTI.Of[Seq[A]]),
		ApplicativeMonoid[R](INTI.Monoid[Seq[A]]()),
		as,
	)
}
