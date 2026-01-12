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

	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/stateless"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	// Option represents an optional value, either Some(value) or None.
	// It is used to handle computations that may or may not return a value,
	// providing a type-safe alternative to nil pointers or sentinel values.
	//
	// Type Parameters:
	//   - A: The type of the value that may be present
	Option[A any] = option.Option[A]

	// Seq is a single-value iterator sequence from Go 1.23+.
	// It represents a lazy sequence of values that can be iterated using range.
	// Operations on Seq are lazy and only execute when the sequence is consumed.
	//
	// Type Parameters:
	//   - T: The type of elements in the sequence
	//
	// Example:
	//
	//	seq := From(1, 2, 3)
	//	for v := range seq {
	//	    fmt.Println(v)
	//	}
	Seq[T any] = I.Seq[T]

	// Seq2 is a key-value iterator sequence from Go 1.23+.
	// It represents a lazy sequence of key-value pairs that can be iterated using range.
	// This is useful for working with map-like data structures in a functional way.
	//
	// Type Parameters:
	//   - K: The type of keys in the sequence
	//   - V: The type of values in the sequence
	//
	// Example:
	//
	//	seq := MonadZip(From(1, 2, 3), From("a", "b", "c"))
	//	for k, v := range seq {
	//	    fmt.Printf("%d: %s\n", k, v)
	//	}
	Seq2[K, V any] = I.Seq2[K, V]

	// Iterator is a stateless iterator type.
	// It provides a functional interface for iterating over collections
	// without maintaining internal state.
	//
	// Type Parameters:
	//   - T: The type of elements produced by the iterator
	Iterator[T any] = stateless.Iterator[T]

	// Predicate is a function that tests a value and returns a boolean.
	// Predicates are commonly used for filtering operations.
	//
	// Type Parameters:
	//   - T: The type of value being tested
	//
	// Example:
	//
	//	isEven := func(x int) bool { return x%2 == 0 }
	//	filtered := Filter(isEven)(From(1, 2, 3, 4))
	Predicate[T any] = predicate.Predicate[T]

	// Kleisli represents a function that takes a value and returns a sequence.
	// This is the monadic bind operation for sequences, also known as flatMap.
	// It's used to chain operations that produce sequences.
	//
	// Type Parameters:
	//   - A: The input type
	//   - B: The element type of the output sequence
	//
	// Example:
	//
	//	duplicate := func(x int) Seq[int] { return From(x, x) }
	//	result := Chain(duplicate)(From(1, 2, 3))
	//	// yields: 1, 1, 2, 2, 3, 3
	Kleisli[A, B any] = func(A) Seq[B]

	// Kleisli2 represents a function that takes a value and returns a key-value sequence.
	// This is the monadic bind operation for key-value sequences.
	//
	// Type Parameters:
	//   - K: The key type in the output sequence
	//   - A: The input type
	//   - B: The value type in the output sequence
	Kleisli2[K, A, B any] = func(A) Seq2[K, B]

	// Operator represents a transformation from one sequence to another.
	// It's a function that takes a Seq[A] and returns a Seq[B].
	// Operators are the building blocks for composing sequence transformations.
	//
	// Type Parameters:
	//   - A: The element type of the input sequence
	//   - B: The element type of the output sequence
	//
	// Example:
	//
	//	double := Map(func(x int) int { return x * 2 })
	//	result := double(From(1, 2, 3))
	//	// yields: 2, 4, 6
	Operator[A, B any] = Kleisli[Seq[A], B]

	// Operator2 represents a transformation from one key-value sequence to another.
	// It's a function that takes a Seq2[K, A] and returns a Seq2[K, B].
	//
	// Type Parameters:
	//   - K: The key type (preserved in the transformation)
	//   - A: The value type of the input sequence
	//   - B: The value type of the output sequence
	Operator2[K, A, B any] = Kleisli2[K, Seq2[K, A], B]

	// Lens is an optic that focuses on a field within a structure.
	// It provides a functional way to get and set values in immutable data structures.
	//
	// Type Parameters:
	//   - S: The structure type
	//   - A: The field type being focused on
	Lens[S, A any] = lens.Lens[S, A]

	// Prism is an optic that focuses on a case of a sum type.
	// It provides a functional way to work with variant types (like Result or Option).
	//
	// Type Parameters:
	//   - S: The sum type
	//   - A: The case type being focused on
	Prism[S, A any] = prism.Prism[S, A]

	// Endomorphism is a function from a type to itself.
	// It represents transformations that preserve the type.
	//
	// Type Parameters:
	//   - A: The type being transformed
	//
	// Example:
	//
	//	increment := func(x int) int { return x + 1 }
	//	result := increment(5) // returns 6
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Pair represents a tuple of two values.
	// It's used to group two related values together.
	//
	// Type Parameters:
	//   - A: The type of the first element
	//   - B: The type of the second element
	//
	// Example:
	//
	//	p := pair.MakePair(1, "hello")
	//	first := pair.Head(p)  // returns 1
	//	second := pair.Tail(p) // returns "hello"
	Pair[A, B any] = pair.Pair[A, B]

	// Void represents the absence of a value, similar to void in other languages.
	// It's used in functions that perform side effects but don't return meaningful values.
	Void = function.Void
)
