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
	"github.com/IBM/fp-go/v2/internal/apply"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT functions convert multiple ReaderResult values into a single ReaderResult of a tuple.
// These are useful for combining independent computations that share the same environment.
// If any computation fails, the entire sequence fails with the first error.

// SequenceT1 wraps a single ReaderResult value in a Tuple1.
// This is primarily for consistency with the other SequenceT functions.
//
// Example:
//
//	rr := readerresult.Of[Config](42)
//	result := readerresult.SequenceT1(rr)
//	// result(cfg) returns (Tuple1{42}, nil)
//
//go:inline
func SequenceT1[R, A any](a ReaderResult[R, A]) ReaderResult[R, T.Tuple1[A]] {
	return apply.SequenceT1(
		Map[R, A, T.Tuple1[A]],

		a)
}

// SequenceT2 combines two independent ReaderResult computations into a tuple.
// Both computations share the same environment.
//
// Example:
//
//	getPort := readerresult.Asks(func(cfg Config) int { return cfg.Port })
//	getHost := readerresult.Asks(func(cfg Config) string { return cfg.Host })
//	result := readerresult.SequenceT2(getPort, getHost)
//	// result(cfg) returns (Tuple2{cfg.Port, cfg.Host}, nil)
//
//go:inline
func SequenceT2[R, A, B any](
	a ReaderResult[R, A],
	b ReaderResult[R, B],
) ReaderResult[R, T.Tuple2[A, B]] {
	return apply.SequenceT2(
		Map[R, A, func(B) T.Tuple2[A, B]],
		Ap[T.Tuple2[A, B], R, B],
		a, b)
}

// SequenceT3 combines three independent ReaderResult computations into a tuple.
//
// Example:
//
//	getUser := getUserRR(42)
//	getConfig := getConfigRR()
//	getStats := getStatsRR()
//	result := readerresult.SequenceT3(getUser, getConfig, getStats)
//	// result(env) returns (Tuple3{user, config, stats}, nil)
//
//go:inline
func SequenceT3[R, A, B, C any](
	a ReaderResult[R, A],
	b ReaderResult[R, B],
	c ReaderResult[R, C],
) ReaderResult[R, T.Tuple3[A, B, C]] {
	return apply.SequenceT3(
		Map[R, A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[func(C) T.Tuple3[A, B, C], R, B],
		Ap[T.Tuple3[A, B, C], R, C],

		a, b, c,
	)
}

// SequenceT4 combines four independent ReaderResult computations into a tuple.
//
// Example:
//
//	result := readerresult.SequenceT4(getUserRR, getConfigRR, getStatsRR, getMetadataRR)
//	// result(env) returns (Tuple4{user, config, stats, metadata}, nil)
//
//go:inline
func SequenceT4[R, A, B, C, D any](
	a ReaderResult[R, A],
	b ReaderResult[R, B],
	c ReaderResult[R, C],
	d ReaderResult[R, D],
) ReaderResult[R, T.Tuple4[A, B, C, D]] {
	return apply.SequenceT4(
		Map[R, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[func(C) func(D) T.Tuple4[A, B, C, D], R, B],
		Ap[func(D) T.Tuple4[A, B, C, D], R, C],
		Ap[T.Tuple4[A, B, C, D], R, D],

		a, b, c, d,
	)
}
