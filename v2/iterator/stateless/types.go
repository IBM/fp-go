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

package stateless

import (
	"iter"

	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	// Option represents an optional value that may or may not be present.
	// It's an alias for option.Option[A] and is used to handle nullable values safely.
	Option[A any] = option.Option[A]

	// Lazy represents a lazily evaluated computation that produces a value of type A.
	// It's an alias for lazy.Lazy[A] and defers computation until the value is needed.
	Lazy[A any] = lazy.Lazy[A]

	// Pair represents a tuple of two values of types L and R.
	// It's an alias for pair.Pair[L, R] where L is the head (left) and R is the tail (right).
	Pair[L, R any] = pair.Pair[L, R]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's an alias for predicate.Predicate[A] and is used for filtering and testing operations.
	Predicate[A any] = predicate.Predicate[A]

	// IO represents a lazy computation that performs side effects and produces a value of type A.
	// It's an alias for io.IO[A] and encapsulates effectful operations.
	IO[A any] = io.IO[A]

	// Iterator represents a stateless, pure, functional iterator over a sequence of values.
	// It's defined as a lazy computation that returns an optional pair of (next iterator, current value).
	// The stateless nature means each iteration step produces a new iterator, making it immutable
	// and safe for concurrent use. When the sequence is exhausted, it returns None.
	// The value is placed in the tail position of the pair because that is what the pair monad
	// operates on, allowing monadic operations to transform values while preserving the iterator state.
	Iterator[U any] Lazy[Option[Pair[Iterator[U], U]]]

	// Kleisli represents a Kleisli arrow for the Iterator monad.
	// It's a function from A to Iterator[B], which allows composition of
	// monadic functions that produce iterators. This is the fundamental building
	// block for chaining iterator operations.
	Kleisli[A, B any] = reader.Reader[A, Iterator[B]]

	// Operator is a specialized Kleisli arrow that operates on Iterator values.
	// It transforms an Iterator[A] into an Iterator[B], making it useful for
	// building pipelines of iterator transformations such as map, filter, and flatMap.
	Operator[A, B any] = Kleisli[Iterator[A], B]

	// Seq represents Go's standard library iterator type for single values.
	// It's an alias for iter.Seq[T] and provides interoperability with Go 1.23+ range-over-func.
	Seq[T any] = iter.Seq[T]

	// Seq2 represents Go's standard library iterator type for key-value pairs.
	// It's an alias for iter.Seq2[K, V] and provides interoperability with Go 1.23+ range-over-func
	// for iterating over maps and other key-value structures.
	Seq2[K, V any] = iter.Seq2[K, V]
)
