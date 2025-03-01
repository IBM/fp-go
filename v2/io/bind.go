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

package io

import (
	G "github.com/IBM/fp-go/v2/io/generic"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[S any](
	empty S,
) IO[S] {
	return G.Do[IO[S], S](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) IO[T],
) func(IO[S1]) IO[S2] {
	return G.Bind[IO[S1], IO[S2], IO[T], S1, S2, T](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(IO[S1]) IO[S2] {
	return G.Let[IO[S1], IO[S2], S1, S2, T](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(IO[S1]) IO[S2] {
	return G.LetTo[IO[S1], IO[S2], S1, S2, T](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) func(IO[T]) IO[S1] {
	return G.BindTo[IO[S1], IO[T], S1, T](setter)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IO[T],
) func(IO[S1]) IO[S2] {
	return G.ApS[IO[S1], IO[S2], IO[T], S1, S2, T](setter, fa)
}
