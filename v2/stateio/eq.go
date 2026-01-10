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

// Eq constructs an equality checker for StateIO values.
// It takes an equality checker for IO[Pair[S, A]] and returns a function that,
// given an initial state S, produces an equality checker for StateIO[S, A].
//
// Two StateIO values are considered equal if, when executed with the same initial state,
// they produce equal IO[Pair[S, A]] results.
//
// Example:
//
//	eqIO := io.FromStrictEquals[Pair[AppState, int]]()
//	eqStateIO := Eq[AppState, int](eqIO)
//	initialState := AppState{}
//	areEqual := eqStateIO(initialState).Equals(stateIO1, stateIO2)
func Eq[
	S, A any](eqr eq.Eq[IO[Pair[S, A]]]) func(S) eq.Eq[StateIO[S, A]] {
	return func(s S) eq.Eq[StateIO[S, A]] {
		return eq.FromEquals(func(l, r StateIO[S, A]) bool {
			return eqr.Equals(l(s), r(s))
		})
	}
}

// FromStrictEquals constructs an equality checker for StateIO values where both
// the state S and value A are comparable types.
//
// This is a convenience function that uses Go's built-in equality (==) for comparison.
// It returns a function that, given an initial state, produces an equality checker
// for StateIO[S, A].
//
// Example:
//
//	eqStateIO := FromStrictEquals[AppState, int]()
//	initialState := AppState{}
//	areEqual := eqStateIO(initialState).Equals(stateIO1, stateIO2)
func FromStrictEquals[
	S, A comparable]() func(S) eq.Eq[StateIO[S, A]] {
	return function.Pipe1(
		io.FromStrictEquals[Pair[S, A]](),
		Eq[S, A],
	)
}
