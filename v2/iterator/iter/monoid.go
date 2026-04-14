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
	G "github.com/IBM/fp-go/v2/internal/iter"
	M "github.com/IBM/fp-go/v2/monoid"
)

// Monoid returns a Monoid instance for Seq[T].
// The monoid's concat operation concatenates sequences sequentially, and the empty value is an empty sequence.
//
// Marble Diagram:
//
//	Seq1:   --1--2--3--|
//	Seq2:   --4--5--6--|
//	Concat: --1--2--3--4--5--6--|
//
// Example:
//
//	m := Monoid[int]()
//	seq1 := From(1, 2)
//	seq2 := From(3, 4)
//	result := m.Concat(seq1, seq2)
//	// yields: 1, 2, 3, 4
//
//go:inline
func Monoid[T any]() M.Monoid[Seq[T]] {
	return G.Monoid[Seq[T]]()
}

// ConcatMonoid returns a Monoid instance for Seq[T] that concatenates sequences sequentially.
// This is an alias for Monoid that makes the sequential concatenation behavior explicit.
//
// A Monoid is an algebraic structure with an associative binary operation (concat)
// and an identity element (empty). For sequences, the concat operation appends one
// sequence after another in deterministic order, and the identity is an empty sequence.
//
// This monoid is useful for functional composition patterns where you need to combine
// multiple sequences sequentially using monoid operations like Reduce, FoldMap, or when
// working with monadic operations that require a monoid instance.
//
// Marble Diagram (Sequential Concatenation):
//
//	Seq1:   --1--2--3--|
//	Seq2:   --4--5--6--|
//	Concat: --1--2--3--4--5--6--|
//	        (deterministic order)
//
// Marble Diagram (vs MergeMonoid):
//
//	ConcatMonoid:
//	  Seq1:   --1--2--3--|
//	  Seq2:              --4--5--6--|
//	  Result: --1--2--3--4--5--6--|
//
//	MergeMonoid:
//	  Seq1:   --1--2--3--|
//	  Seq2:   --4--5--6--|
//	  Result: --1-4-2-5-3-6--|
//	          (non-deterministic)
//
// Type Parameters:
//   - T: The type of elements in the sequences
//
// Returns:
//   - Monoid[Seq[T]]: A monoid instance with:
//   - Concat: Appends sequences sequentially (deterministic order)
//   - Empty: Returns an empty sequence
//
// Properties:
//   - Identity: concat(empty, x) = concat(x, empty) = x
//   - Associativity: concat(concat(a, b), c) = concat(a, concat(b, c))
//   - Deterministic: Elements always appear in the order of the input sequences
//
// Comparison with MergeMonoid:
//
// ConcatMonoid and MergeMonoid serve different purposes:
//
//   - ConcatMonoid: Sequential concatenation
//
//   - Order: Deterministic - elements from first sequence, then second, etc.
//
//   - Concurrency: No concurrency - sequences are processed one after another
//
//   - Performance: Lower overhead, no goroutines or channels
//
//   - Use when: Order matters, no I/O operations, or simplicity is preferred
//
//   - MergeMonoid: Concurrent merging
//
//   - Order: Non-deterministic - elements interleaved based on timing
//
//   - Concurrency: Spawns goroutines for each sequence
//
//   - Performance: Better for I/O-bound operations, higher overhead for CPU-bound
//
//   - Use when: Order doesn't matter, parallel I/O, or concurrent processing needed
//
// Example Usage:
//
//	// Create a monoid for concatenating integer sequences
//	monoid := ConcatMonoid[int]()
//
//	// Use with Reduce to concatenate multiple sequences
//	sequences := []Seq[int]{
//	    From(1, 2, 3),
//	    From(4, 5, 6),
//	    From(7, 8, 9),
//	}
//	concatenated := MonadReduce(From(sequences...), monoid.Concat, monoid.Empty)
//	// yields: 1, 2, 3, 4, 5, 6, 7, 8, 9 (deterministic order)
//
// Example with Empty Identity:
//
//	monoid := ConcatMonoid[int]()
//	seq := From(1, 2, 3)
//
//	// Concatenating with empty is identity
//	result1 := monoid.Concat(monoid.Empty, seq)  // same as seq
//	result2 := monoid.Concat(seq, monoid.Empty)  // same as seq
//
// Example with FoldMap:
//
//	// Convert each number to a sequence and concatenate all results
//	monoid := ConcatMonoid[int]()
//	numbers := From(1, 2, 3)
//	result := MonadFoldMap(numbers, func(n int) Seq[int] {
//	    return From(n, n*10, n*100)
//	}, monoid)
//	// yields: 1, 10, 100, 2, 20, 200, 3, 30, 300 (deterministic order)
//
// Example Comparing ConcatMonoid vs MergeMonoid:
//
//	seq1 := From(1, 2, 3)
//	seq2 := From(4, 5, 6)
//
//	// ConcatMonoid: Sequential, deterministic
//	concatMonoid := ConcatMonoid[int]()
//	concat := concatMonoid.Concat(seq1, seq2)
//	// Always yields: 1, 2, 3, 4, 5, 6
//
//	// MergeMonoid: Concurrent, non-deterministic
//	mergeMonoid := MergeMonoid[int](10)
//	merged := mergeMonoid.Concat(seq1, seq2)
//	// May yield: 1, 4, 2, 5, 3, 6 (order varies)
//
// See Also:
//   - Monoid: The base monoid function (alias)
//   - MergeMonoid: Concurrent merging monoid
//   - MonadChain: Sequential flattening of sequences
//   - Empty: Creates an empty sequence
//
//go:inline
func ConcatMonoid[T any]() M.Monoid[Seq[T]] {
	return Monoid[T]()
}
