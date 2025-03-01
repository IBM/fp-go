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

package ioeither

import (
	G "github.com/IBM/fp-go/v2/ioeither/generic"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[E, S any](
	empty S,
) IOEither[E, S] {
	return G.Do[IOEither[E, S], E, S](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) IOEither[E, T],
) func(IOEither[E, S1]) IOEither[E, S2] {
	return G.Bind[IOEither[E, S1], IOEither[E, S2], IOEither[E, T], E, S1, S2, T](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(IOEither[E, S1]) IOEither[E, S2] {
	return G.Let[IOEither[E, S1], IOEither[E, S2], E, S1, S2, T](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(IOEither[E, S1]) IOEither[E, S2] {
	return G.LetTo[IOEither[E, S1], IOEither[E, S2], E, S1, S2, T](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[E, S1, T any](
	setter func(T) S1,
) func(IOEither[E, T]) IOEither[E, S1] {
	return G.BindTo[IOEither[E, S1], IOEither[E, T], E, S1, T](setter)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOEither[E, T],
) func(IOEither[E, S1]) IOEither[E, S2] {
	return G.ApS[IOEither[E, S1], IOEither[E, S2], IOEither[E, T], E, S1, S2, T](setter, fa)
}
