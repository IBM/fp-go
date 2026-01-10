// Copyright (c) 2024 - 2025 IBM Corp.
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

package stateio

import (
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
)

// Eq implements the equals predicate for values contained in the [StateIO] monad
func Eq[
	S, A any](eqr eq.Eq[IO[Pair[S, A]]]) func(S) eq.Eq[StateIO[S, A]] {
	return func(s S) eq.Eq[StateIO[S, A]] {
		return eq.FromEquals(func(l, r StateIO[S, A]) bool {
			return eqr.Equals(l(s), r(s))
		})
	}
}

// FromStrictEquals constructs an [eq.Eq] from the canonical comparison function
func FromStrictEquals[
	S, A comparable]() func(S) eq.Eq[StateIO[S, A]] {
	return function.Pipe1(
		io.FromStrictEquals[Pair[S, A]](),
		Eq[S, A],
	)
}
