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

package readeroption

import (
	G "github.com/IBM/fp-go/v2/readeroption/generic"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT functions convert multiple ReaderOption values into a single ReaderOption containing a tuple.
// If any input is None, the entire result is None.
// Otherwise, returns Some containing a tuple of all the unwrapped values.
//
// These functions are useful for combining multiple independent ReaderOption computations
// where you need to preserve the individual types of each result.

// SequenceT1 converts a single ReaderOption into a ReaderOption of a 1-tuple.
// This is mainly useful for consistency with the other SequenceT functions.
//
// Example:
//
//	type Config struct { ... }
//
//	user := readeroption.Of[Config](User{Name: "Alice"})
//	result := readeroption.SequenceT1(user)
//	// result(config) returns option.Some(tuple.MakeTuple1(User{Name: "Alice"}))
func SequenceT1[E, A any](a ReaderOption[E, A]) ReaderOption[E, T.Tuple1[A]] {
	return G.SequenceT1[
		ReaderOption[E, A],
		ReaderOption[E, T.Tuple1[A]],
	](a)
}

// SequenceT2 combines two ReaderOption values into a ReaderOption of a 2-tuple.
// If either input is None, the result is None.
//
// Example:
//
//	type Config struct { ... }
//
//	user := readeroption.Of[Config](User{Name: "Alice"})
//	count := readeroption.Of[Config](42)
//
//	result := readeroption.SequenceT2(user, count)
//	// result(config) returns option.Some(tuple.MakeTuple2(User{Name: "Alice"}, 42))
//
//	noneUser := readeroption.None[Config, User]()
//	result2 := readeroption.SequenceT2(noneUser, count)
//	// result2(config) returns option.None[tuple.Tuple2[User, int]]()
func SequenceT2[E, A, B any](
	a ReaderOption[E, A],
	b ReaderOption[E, B],
) ReaderOption[E, T.Tuple2[A, B]] {
	return G.SequenceT2[
		ReaderOption[E, A],
		ReaderOption[E, B],
		ReaderOption[E, T.Tuple2[A, B]],
	](a, b)
}

// SequenceT3 combines three ReaderOption values into a ReaderOption of a 3-tuple.
// If any input is None, the result is None.
//
// Example:
//
//	type Config struct { ... }
//
//	user := readeroption.Of[Config](User{Name: "Alice"})
//	count := readeroption.Of[Config](42)
//	active := readeroption.Of[Config](true)
//
//	result := readeroption.SequenceT3(user, count, active)
//	// result(config) returns option.Some(tuple.MakeTuple3(User{Name: "Alice"}, 42, true))
func SequenceT3[E, A, B, C any](
	a ReaderOption[E, A],
	b ReaderOption[E, B],
	c ReaderOption[E, C],
) ReaderOption[E, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		ReaderOption[E, A],
		ReaderOption[E, B],
		ReaderOption[E, C],
		ReaderOption[E, T.Tuple3[A, B, C]],
	](a, b, c)
}

// SequenceT4 combines four ReaderOption values into a ReaderOption of a 4-tuple.
// If any input is None, the result is None.
//
// Example:
//
//	type Config struct { ... }
//
//	user := readeroption.Of[Config](User{Name: "Alice"})
//	count := readeroption.Of[Config](42)
//	active := readeroption.Of[Config](true)
//	score := readeroption.Of[Config](95.5)
//
//	result := readeroption.SequenceT4(user, count, active, score)
//	// result(config) returns option.Some(tuple.MakeTuple4(User{Name: "Alice"}, 42, true, 95.5))
func SequenceT4[E, A, B, C, D any](
	a ReaderOption[E, A],
	b ReaderOption[E, B],
	c ReaderOption[E, C],
	d ReaderOption[E, D],
) ReaderOption[E, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		ReaderOption[E, A],
		ReaderOption[E, B],
		ReaderOption[E, C],
		ReaderOption[E, D],
		ReaderOption[E, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
