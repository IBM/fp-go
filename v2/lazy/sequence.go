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

package lazy

import (
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

// SequenceT1 combines a single lazy computation into a lazy tuple.
// This is mainly useful for consistency with the other SequenceT functions.
//
// Example:
//
//	lazy1 := lazy.Of(42)
//	result := lazy.SequenceT1(lazy1)()
//	// result is tuple.Tuple1[int]{F1: 42}
func SequenceT1[A any](a Lazy[A]) Lazy[tuple.Tuple1[A]] {
	return io.SequenceT1(a)
}

// SequenceT2 combines two lazy computations into a lazy tuple of two elements.
// Both computations are evaluated when the result is evaluated.
//
// Example:
//
//	lazy1 := lazy.Of(42)
//	lazy2 := lazy.Of("hello")
//	result := lazy.SequenceT2(lazy1, lazy2)()
//	// result is tuple.Tuple2[int, string]{F1: 42, F2: "hello"}
func SequenceT2[A, B any](a Lazy[A], b Lazy[B]) Lazy[tuple.Tuple2[A, B]] {
	return io.SequenceT2(a, b)
}

// SequenceT3 combines three lazy computations into a lazy tuple of three elements.
// All computations are evaluated when the result is evaluated.
//
// Example:
//
//	lazy1 := lazy.Of(42)
//	lazy2 := lazy.Of("hello")
//	lazy3 := lazy.Of(true)
//	result := lazy.SequenceT3(lazy1, lazy2, lazy3)()
//	// result is tuple.Tuple3[int, string, bool]{F1: 42, F2: "hello", F3: true}
func SequenceT3[A, B, C any](a Lazy[A], b Lazy[B], c Lazy[C]) Lazy[tuple.Tuple3[A, B, C]] {
	return io.SequenceT3(a, b, c)
}

// SequenceT4 combines four lazy computations into a lazy tuple of four elements.
// All computations are evaluated when the result is evaluated.
//
// Example:
//
//	lazy1 := lazy.Of(42)
//	lazy2 := lazy.Of("hello")
//	lazy3 := lazy.Of(true)
//	lazy4 := lazy.Of(3.14)
//	result := lazy.SequenceT4(lazy1, lazy2, lazy3, lazy4)()
//	// result is tuple.Tuple4[int, string, bool, float64]{F1: 42, F2: "hello", F3: true, F4: 3.14}
func SequenceT4[A, B, C, D any](a Lazy[A], b Lazy[B], c Lazy[C], d Lazy[D]) Lazy[tuple.Tuple4[A, B, C, D]] {
	return io.SequenceT4(a, b, c, d)
}
