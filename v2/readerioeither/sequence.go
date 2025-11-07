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

package readerioeither

import (
	G "github.com/IBM/fp-go/v2/readerioeither/generic"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT1 converts a single ReaderIOEither into a ReaderIOEither of a 1-tuple.
// This is useful for uniformly handling computations with different arities.
//
// If the input computation fails, the result will be a Left with the error.
// If it succeeds, the result will be a Right with a tuple containing the value.
//
// Example:
//
//	result := SequenceT1(Of[Config, error](42))
//	// result(cfg)() returns Right(Tuple1{42})
//
//go:inline
func SequenceT1[R, E, A any](a ReaderIOEither[R, E, A]) ReaderIOEither[R, E, T.Tuple1[A]] {
	return G.SequenceT1[
		ReaderIOEither[R, E, A],
		ReaderIOEither[R, E, T.Tuple1[A]],
	](a)
}

// SequenceT2 combines two ReaderIOEither computations into a single ReaderIOEither of a 2-tuple.
// Both computations are executed, and if both succeed, their results are combined into a tuple.
// If either fails, the result is a Left with the first error encountered.
//
// This is useful for running multiple independent computations and collecting their results.
//
// Example:
//
//	result := SequenceT2(
//	    fetchUser(123),
//	    fetchProfile(123),
//	)
//	// result(cfg)() returns Right(Tuple2{user, profile}) or Left(error)
//
//go:inline
func SequenceT2[R, E, A, B any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B]) ReaderIOEither[R, E, T.Tuple2[A, B]] {
	return G.SequenceT2[
		ReaderIOEither[R, E, A],
		ReaderIOEither[R, E, B],
		ReaderIOEither[R, E, T.Tuple2[A, B]],
	](a, b)
}

// SequenceT3 combines three ReaderIOEither computations into a single ReaderIOEither of a 3-tuple.
// All three computations are executed, and if all succeed, their results are combined into a tuple.
// If any fails, the result is a Left with the first error encountered.
//
// Example:
//
//	result := SequenceT3(
//	    fetchUser(123),
//	    fetchProfile(123),
//	    fetchSettings(123),
//	)
//	// result(cfg)() returns Right(Tuple3{user, profile, settings}) or Left(error)
//
//go:inline
func SequenceT3[R, E, A, B, C any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B], c ReaderIOEither[R, E, C]) ReaderIOEither[R, E, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		ReaderIOEither[R, E, A],
		ReaderIOEither[R, E, B],
		ReaderIOEither[R, E, C],
		ReaderIOEither[R, E, T.Tuple3[A, B, C]],
	](a, b, c)
}

// SequenceT4 combines four ReaderIOEither computations into a single ReaderIOEither of a 4-tuple.
// All four computations are executed, and if all succeed, their results are combined into a tuple.
// If any fails, the result is a Left with the first error encountered.
//
// Example:
//
//	result := SequenceT4(
//	    fetchUser(123),
//	    fetchProfile(123),
//	    fetchSettings(123),
//	    fetchPreferences(123),
//	)
//	// result(cfg)() returns Right(Tuple4{user, profile, settings, prefs}) or Left(error)
//
//go:inline
func SequenceT4[R, E, A, B, C, D any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B], c ReaderIOEither[R, E, C], d ReaderIOEither[R, E, D]) ReaderIOEither[R, E, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		ReaderIOEither[R, E, A],
		ReaderIOEither[R, E, B],
		ReaderIOEither[R, E, C],
		ReaderIOEither[R, E, D],
		ReaderIOEither[R, E, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
