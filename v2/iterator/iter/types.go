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

package iter

import (
	I "iter"

	"github.com/IBM/fp-go/v2/iterator/stateless"
	"github.com/IBM/fp-go/v2/optics/lens/option"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	// Option represents an optional value, either Some(value) or None.
	Option[A any] = option.Option[A]

	// Seq is a single-value iterator sequence from Go 1.23+.
	// It represents a lazy sequence of values that can be iterated using range.
	Seq[T any] = I.Seq[T]

	// Seq2 is a key-value iterator sequence from Go 1.23+.
	// It represents a lazy sequence of key-value pairs that can be iterated using range.
	Seq2[K, V any] = I.Seq2[K, V]

	// Iterator is a stateless iterator type.
	Iterator[T any] = stateless.Iterator[T]

	// Predicate is a function that tests a value and returns a boolean.
	Predicate[T any] = predicate.Predicate[T]

	// Kleisli represents a function that takes a value and returns a sequence.
	// This is the monadic bind operation for sequences.
	Kleisli[A, B any] = func(A) Seq[B]

	// Kleisli2 represents a function that takes a value and returns a key-value sequence.
	Kleisli2[K, A, B any] = func(A) Seq2[K, B]

	// Operator represents a transformation from one sequence to another.
	// It's a function that takes a Seq[A] and returns a Seq[B].
	Operator[A, B any] = Kleisli[Seq[A], B]

	// Operator2 represents a transformation from one key-value sequence to another.
	Operator2[K, A, B any] = Kleisli2[K, Seq2[K, A], B]
)
