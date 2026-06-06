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
	return ConcatMonoid[T](defaultBufferSize)
}

func makeMonoid[T any](concAll Operator[Seq[T], T]) M.Monoid[Seq[T]] {
	return M.MakeMonoid(
		func(l, r Seq[T]) Seq[T] {
			return concAll(From(l, r))
		},
		Empty[T](),
	)
}

// MonoidSeq returns a Monoid[Seq[T]] whose Concat operation uses purely
// sequential nested iteration — no goroutines, no channels. Each left-hand
// sequence is fully drained before the right-hand sequence is started.
//
// Use this when inner sequences are cheap and synchronous and goroutine
// overhead must be avoided.
//
// Properties:
//   - Identity: Concat(Empty(), x) = Concat(x, Empty()) = x
//   - Associativity: output order is always left-to-right
//
// See Also:
//   - Monoid: Concurrent producers, same output order (uses defaultBufferSize)
//   - MonoidPar: Always concurrent, never selects the sequential path
//   - ConcatMonoid: Concurrent monoid with configurable buffer size
func MonoidSeq[T any]() M.Monoid[Seq[T]] {
	return makeMonoid(ConcatAllSeq[T]())
}

// MonoidPar returns a Monoid[Seq[T]] whose Concat operation runs both
// sequences concurrently via ConcatAllPar with defaultBufferSize, draining
// them in left-to-right order. Unlike Monoid, it never selects the
// goroutine-free sequential path.
//
// Use this when you need concurrent producers regardless of buffer size,
// for example to ensure forward progress in I/O-bound pipelines.
//
// Properties:
//   - Identity: Concat(Empty(), x) = Concat(x, Empty()) = x
//   - Associativity: output order is always left-to-right
//
// See Also:
//   - Monoid: Dispatches to sequential when bufSize == 1
//   - MonoidSeq: Always sequential, no goroutines
//   - ConcatMonoid: Configurable buffer size
func MonoidPar[T any]() M.Monoid[Seq[T]] {
	return makeMonoid(ConcatAllPar[T](defaultBufferSize))
}

// ConcatMonoid returns a Monoid[Seq[T]] whose Concat operation runs both
// sequences concurrently (each in its own goroutine) but drains them in order,
// guaranteeing deterministic output. bufSize is forwarded to ConcatBuf and
// controls channel buffer capacity (negative → 0, unbuffered).
//
// This is the order-preserving counterpart of MergeMonoid: both use goroutines,
// but ConcatMonoid's drain is sequential while MergeMonoid's is not.
//
// Properties:
//   - Identity: Concat(Empty(), x) = Concat(x, Empty()) = x
//   - Associativity: output order is always left-to-right
//
// Marble Diagram:
//
//	ConcatMonoid.Concat(seq1, seq2):
//	  seq1 (goroutine): --1--2--3--|
//	  seq2 (goroutine): --4--5--6--|   (concurrent)
//	  output:           --1--2--3--4--5--6--|   (in order)
//
//	MergeMonoid.Concat(seq1, seq2):
//	  output:           --1-4-2-5-3-6--|   (non-deterministic)
//
// See Also:
//   - MergeMonoid: Non-deterministic order, same concurrency model
//   - ConcatBuf: The underlying implementation
func ConcatMonoid[T any](bufSize int) M.Monoid[Seq[T]] {
	return makeMonoid(ConcatAll[T](defaultBufferSize))
}
