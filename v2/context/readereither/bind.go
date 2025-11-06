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

package readereither

import (
	"context"

	G "github.com/IBM/fp-go/v2/readereither/generic"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[S any](
	empty S,
) ReaderEither[S] {
	return G.Do[ReaderEither[S], context.Context, error, S](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) ReaderEither[T],
) func(ReaderEither[S1]) ReaderEither[S2] {
	return G.Bind[ReaderEither[S1], ReaderEither[S2], ReaderEither[T], context.Context, error, S1, S2, T](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(ReaderEither[S1]) ReaderEither[S2] {
	return G.Let[ReaderEither[S1], ReaderEither[S2], context.Context, error, S1, S2, T](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(ReaderEither[S1]) ReaderEither[S2] {
	return G.LetTo[ReaderEither[S1], ReaderEither[S2], context.Context, error, S1, S2, T](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) func(ReaderEither[T]) ReaderEither[S1] {
	return G.BindTo[ReaderEither[S1], ReaderEither[T], context.Context, error, S1, T](setter)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderEither[T],
) func(ReaderEither[S1]) ReaderEither[S2] {
	return G.ApS[ReaderEither[S1], ReaderEither[S2], ReaderEither[T], context.Context, error, S1, S2, T](setter, fa)
}
