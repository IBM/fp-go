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

package iooption

import (
	G "github.com/IBM/fp-go/iooption/generic"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[S any](
	empty S,
) IOOption[S] {
	return G.Do[IOOption[S], S](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) IOOption[T],
) func(IOOption[S1]) IOOption[S2] {
	return G.Bind[IOOption[S1], IOOption[S2], IOOption[T], S1, S2, T](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(IOOption[S1]) IOOption[S2] {
	return G.Let[IOOption[S1], IOOption[S2], S1, S2, T](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(IOOption[S1]) IOOption[S2] {
	return G.LetTo[IOOption[S1], IOOption[S2], S1, S2, T](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) func(IOOption[T]) IOOption[S1] {
	return G.BindTo[IOOption[S1], IOOption[T], S1, T](setter)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOOption[T],
) func(IOOption[S1]) IOOption[S2] {
	return G.ApS[IOOption[S1], IOOption[S2], IOOption[T], S1, S2, T](setter, fa)
}
