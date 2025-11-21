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
	G "github.com/IBM/fp-go/v2/readereither/generic"
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
//	// result(cfg) returns result.Of(Tuple1{42})
func SequenceT1[L, A any](a ReaderResult[L, A]) ReaderResult[L, T.Tuple1[A]] {
	return G.SequenceT1[
		ReaderResult[L, A],
		ReaderResult[L, T.Tuple1[A]],
	](a)
}

// SequenceT2 combines two independent ReaderResult computations into a tuple.
// Both computations share the same environment.
//
// Example:
//
//	getPort := readerresult.Asks(func(cfg Config) int { return cfg.Port })
//	getHost := readerresult.Asks(func(cfg Config) string { return cfg.Host })
//	result := readerresult.SequenceT2(getPort, getHost)
//	// result(cfg) returns result.Of(Tuple2{cfg.Port, cfg.Host})
func SequenceT2[L, A, B any](
	a ReaderResult[L, A],
	b ReaderResult[L, B],
) ReaderResult[L, T.Tuple2[A, B]] {
	return G.SequenceT2[
		ReaderResult[L, A],
		ReaderResult[L, B],
		ReaderResult[L, T.Tuple2[A, B]],
	](a, b)
}

// SequenceT3 combines three independent ReaderResult computations into a tuple.
//
// Example:
//
//	getUser := getUserRR(42)
//	getConfig := getConfigRR()
//	getStats := getStatsRR()
//	result := readerresult.SequenceT3(getUser, getConfig, getStats)
//	// result(env) returns result.Of(Tuple3{user, config, stats})
func SequenceT3[L, A, B, C any](
	a ReaderResult[L, A],
	b ReaderResult[L, B],
	c ReaderResult[L, C],
) ReaderResult[L, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		ReaderResult[L, A],
		ReaderResult[L, B],
		ReaderResult[L, C],
		ReaderResult[L, T.Tuple3[A, B, C]],
	](a, b, c)
}

// SequenceT4 combines four independent ReaderResult computations into a tuple.
//
// Example:
//
//	result := readerresult.SequenceT4(getUserRR, getConfigRR, getStatsRR, getMetadataRR)
//	// result(env) returns result.Of(Tuple4{user, config, stats, metadata})
func SequenceT4[L, A, B, C, D any](
	a ReaderResult[L, A],
	b ReaderResult[L, B],
	c ReaderResult[L, C],
	d ReaderResult[L, D],
) ReaderResult[L, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		ReaderResult[L, A],
		ReaderResult[L, B],
		ReaderResult[L, C],
		ReaderResult[L, D],
		ReaderResult[L, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
