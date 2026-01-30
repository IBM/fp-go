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

package readeriooption

import (
	"github.com/IBM/fp-go/v2/internal/apply"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT functions convert multiple ReaderIOOption values into a single ReaderIOOption containing a tuple.
// If any input is None, the entire result is None.
// Otherwise, returns Some containing a tuple of all the unwrapped values.
//
// These functions are useful for combining multiple independent ReaderIOOption computations
// where you need to preserve the individual types of each result.

// SequenceT1 converts a single ReaderIOOption into a ReaderIOOption of a 1-tuple.
// This is mainly useful for consistency with the other SequenceT functions.
//
// Example:
//
//	type Config struct { ... }
//
//	user := readeroption.Of[Config](User{Name: "Alice"})
//	result := readeroption.SequenceT1(user)
//	// result(config) returns option.Some(tuple.MakeTuple1(User{Name: "Alice"}))
func SequenceT1[R, A any](a ReaderIOOption[R, A]) ReaderIOOption[R, T.Tuple1[A]] {
	return apply.SequenceT1(
		Map,
		a,
	)
}

// SequenceT2 combines two ReaderIOOption values into a ReaderIOOption of a 2-tuple.
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
func SequenceT2[R, A, B any](
	a ReaderIOOption[R, A],
	b ReaderIOOption[R, B],
) ReaderIOOption[R, T.Tuple2[A, B]] {
	return apply.SequenceT2(
		Map,
		Ap,
		a,
		b,
	)
}

// SequenceT3 combines three ReaderIOOption values into a ReaderIOOption of a 3-tuple.
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
func SequenceT3[R, A, B, C any](
	a ReaderIOOption[R, A],
	b ReaderIOOption[R, B],
	c ReaderIOOption[R, C],
) ReaderIOOption[R, T.Tuple3[A, B, C]] {
	return apply.SequenceT3(
		Map,
		Ap,
		Ap,
		a,
		b,
		c,
	)
}

// SequenceT4 combines four ReaderIOOption values into a ReaderIOOption of a 4-tuple.
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
func SequenceT4[R, A, B, C, D any](
	a ReaderIOOption[R, A],
	b ReaderIOOption[R, B],
	c ReaderIOOption[R, C],
	d ReaderIOOption[R, D],
) ReaderIOOption[R, T.Tuple4[A, B, C, D]] {
	return apply.SequenceT4(
		Map,
		Ap,
		Ap,
		Ap,
		a,
		b,
		c,
		d,
	)
}
