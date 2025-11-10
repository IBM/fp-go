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

package readerioresult

import (
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT1 converts a single ReaderIOResult into a ReaderIOResult of a 1-tuple.
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
func SequenceT1[R, A any](a ReaderIOResult[R, A]) ReaderIOResult[R, T.Tuple1[A]] {
	return RIOE.SequenceT1(a)
}

// SequenceT2 combines two ReaderIOResult computations into a single ReaderIOResult of a 2-tuple.
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
func SequenceT2[R, A, B any](a ReaderIOResult[R, A], b ReaderIOResult[R, B]) ReaderIOResult[R, T.Tuple2[A, B]] {
	return RIOE.SequenceT2(a, b)
}

// SequenceT3 combines three ReaderIOResult computations into a single ReaderIOResult of a 3-tuple.
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
func SequenceT3[R, A, B, C any](a ReaderIOResult[R, A], b ReaderIOResult[R, B], c ReaderIOResult[R, C]) ReaderIOResult[R, T.Tuple3[A, B, C]] {
	return RIOE.SequenceT3(a, b, c)
}

// SequenceT4 combines four ReaderIOResult computations into a single ReaderIOResult of a 4-tuple.
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
func SequenceT4[R, A, B, C, D any](a ReaderIOResult[R, A], b ReaderIOResult[R, B], c ReaderIOResult[R, C], d ReaderIOResult[R, D]) ReaderIOResult[R, T.Tuple4[A, B, C, D]] {
	return RIOE.SequenceT4(a, b, c, d)
}
