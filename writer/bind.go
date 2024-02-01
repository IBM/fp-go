// Copyright (c) 2023 IBM Corp.
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

package writer

import (
	M "github.com/IBM/fp-go/monoid"
	G "github.com/IBM/fp-go/writer/generic"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[S, W any](m M.Monoid[W]) func(S) Writer[W, S] {
	return G.Do[Writer[W, S], W, S](m)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[S1, S2, T, W any](
	setter func(T) func(S1) S2,
	f func(S1) Writer[W, T],
) func(Writer[W, S1]) Writer[W, S2] {
	return G.Bind[Writer[W, S1], Writer[W, S2], Writer[W, T], W, S1, S2, T](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[W, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(Writer[W, S1]) Writer[W, S2] {
	return G.Let[Writer[W, S1], Writer[W, S2], W, S1, S2, T](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[W, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(Writer[W, S1]) Writer[W, S2] {
	return G.LetTo[Writer[W, S1], Writer[W, S2], W, S1, S2, T](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[W, S1, T any](
	setter func(T) S1,
) func(Writer[W, T]) Writer[W, S1] {
	return G.BindTo[Writer[W, S1], Writer[W, T], W, S1, T](setter)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[S1, S2, T, W any](
	setter func(T) func(S1) S2,
	fa Writer[W, T],
) func(Writer[W, S1]) Writer[W, S2] {
	return G.ApS[Writer[W, S1], Writer[W, S2], Writer[W, T], W, S1, S2, T](setter, fa)
}
