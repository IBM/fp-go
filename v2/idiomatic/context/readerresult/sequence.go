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

package readerresult

import (
	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT1 wraps a single ReaderResult in a Tuple1.
//
// This is mainly for consistency with the other SequenceT functions.
//
// Type Parameters:
//   - A: The value type
//
// Parameters:
//   - a: A ReaderResult[A]
//
// Returns:
//   - A ReaderResult[Tuple1[A]]
//
// Example:
//
//	rr := readerresult.Right(42)
//	result := readerresult.SequenceT1(rr)
//	tuple, err := result(ctx)  // Returns (Tuple1{42}, nil)
//
//go:inline
func SequenceT1[A any](a ReaderResult[A]) ReaderResult[T.Tuple1[A]] {
	return RR.SequenceT1(a)
}

// SequenceT2 combines two independent ReaderResult computations into a tuple.
//
// Both computations are executed with the same context. If either fails,
// the entire operation fails with the first error encountered.
//
// Type Parameters:
//   - A: The first value type
//   - B: The second value type
//
// Parameters:
//   - a: The first ReaderResult
//   - b: The second ReaderResult
//
// Returns:
//   - A ReaderResult[Tuple2[A, B]] containing both results
//
// Example:
//
//	getUser := readerresult.Right(User{ID: 1})
//	getConfig := readerresult.Right(Config{Port: 8080})
//	result := readerresult.SequenceT2(getUser, getConfig)
//	tuple, err := result(ctx)  // Returns (Tuple2{User, Config}, nil)
//
//go:inline
func SequenceT2[A, B any](
	a ReaderResult[A],
	b ReaderResult[B],
) ReaderResult[T.Tuple2[A, B]] {
	return RR.SequenceT2(a, b)
}

// SequenceT3 combines three independent ReaderResult computations into a tuple.
//
// All computations are executed with the same context. If any fails,
// the entire operation fails with the first error encountered.
//
// Type Parameters:
//   - A: The first value type
//   - B: The second value type
//   - C: The third value type
//
// Parameters:
//   - a: The first ReaderResult
//   - b: The second ReaderResult
//   - c: The third ReaderResult
//
// Returns:
//   - A ReaderResult[Tuple3[A, B, C]] containing all three results
//
//go:inline
func SequenceT3[A, B, C any](
	a ReaderResult[A],
	b ReaderResult[B],
	c ReaderResult[C],
) ReaderResult[T.Tuple3[A, B, C]] {
	return RR.SequenceT3(a, b, c)
}

// SequenceT4 combines four independent ReaderResult computations into a tuple.
//
// All computations are executed with the same context. If any fails,
// the entire operation fails with the first error encountered.
//
// Type Parameters:
//   - A: The first value type
//   - B: The second value type
//   - C: The third value type
//   - D: The fourth value type
//
// Parameters:
//   - a: The first ReaderResult
//   - b: The second ReaderResult
//   - c: The third ReaderResult
//   - d: The fourth ReaderResult
//
// Returns:
//   - A ReaderResult[Tuple4[A, B, C, D]] containing all four results
//
//go:inline
func SequenceT4[A, B, C, D any](
	a ReaderResult[A],
	b ReaderResult[B],
	c ReaderResult[C],
	d ReaderResult[D],
) ReaderResult[T.Tuple4[A, B, C, D]] {
	return RR.SequenceT4(a, b, c, d)
}
