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

package reader

import (
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// ApplySemigroup lifts a Semigroup[A] into a Semigroup[Reader[R, A]].
// This allows you to combine two Readers that produce semigroup values by combining
// their results using the semigroup's concat operation.
//
// The _map and _ap parameters are the Map and Ap operations for the Reader type,
// typically obtained from the reader package.
//
// Example:
//
//	type Config struct { Multiplier int }
//	// Using the additive semigroup for integers
//	intSemigroup := semigroup.MakeSemigroup(func(a, b int) int { return a + b })
//	readerSemigroup := reader.ApplySemigroup(
//	    reader.MonadMap[Config, int, func(int) int],
//	    reader.MonadAp[int, Config, int],
//	    intSemigroup,
//	)
//
//	r1 := reader.Of[Config](5)
//	r2 := reader.Of[Config](3)
//	combined := readerSemigroup.Concat(r1, r2)
//	result := combined(Config{Multiplier: 1}) // 8
func ApplySemigroup[R, A any](
	_map func(func(R) A, func(A) func(A) A) func(R, func(A) A),
	_ap func(func(R, func(A) A), func(R) A) func(R) A,

	s S.Semigroup[A],
) S.Semigroup[func(R) A] {
	return S.ApplySemigroup(_map, _ap, s)
}

// ApplicativeMonoid lifts a Monoid[A] into a Monoid[Reader[R, A]].
// This allows you to combine Readers that produce monoid values, with an empty/identity Reader.
//
// The _of parameter is the Of operation (pure/return) for the Reader type.
// The _map and _ap parameters are the Map and Ap operations for the Reader type.
//
// Example:
//
//	type Config struct { Prefix string }
//	// Using the string concatenation monoid
//	stringMonoid := monoid.MakeMonoid("", func(a, b string) string { return a + b })
//	readerMonoid := reader.ApplicativeMonoid(
//	    reader.Of[Config, string],
//	    reader.MonadMap[Config, string, func(string) string],
//	    reader.MonadAp[string, Config, string],
//	    stringMonoid,
//	)
//
//	r1 := reader.Asks(func(c Config) string { return c.Prefix })
//	r2 := reader.Of[Config]("hello")
//	combined := readerMonoid.Concat(r1, r2)
//	result := combined(Config{Prefix: ">> "}) // ">> hello"
//	empty := readerMonoid.Empty()(Config{Prefix: "any"}) // ""
func ApplicativeMonoid[R, A any](
	_of func(A) func(R) A,
	_map func(func(R) A, func(A) func(A) A) func(R, func(A) A),
	_ap func(func(R, func(A) A), func(R) A) func(R) A,

	m M.Monoid[A],
) M.Monoid[func(R) A] {
	return M.ApplicativeMonoid(_of, _map, _ap, m)
}
