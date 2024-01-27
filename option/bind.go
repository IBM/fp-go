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

package option

import (
	A "github.com/IBM/fp-go/internal/apply"
	C "github.com/IBM/fp-go/internal/chain"
	F "github.com/IBM/fp-go/internal/functor"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[S any](
	empty S,
) Option[S] {
	return Of(empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[S1, S2, A any](
	setter func(A) func(S1) S2,
	f func(S1) Option[A],
) func(Option[S1]) Option[S2] {
	return C.Bind(
		Chain[S1, S2],
		Map[A, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, B any](
	key func(B) func(S1) S2,
	f func(S1) B,
) func(Option[S1]) Option[S2] {
	return F.Let(
		Map[S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(Option[S1]) Option[S2] {
	return F.LetTo(
		Map[S1, S2],
		key,
		b,
	)
}

// BindTo attaches a value to a context [S1] to produce a context [S2]
func BindTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Option[T],
) func(Option[S1]) Option[S2] {
	return C.BindTo(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		fa,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Option[T],
) func(Option[S1]) Option[S2] {
	return A.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
		setter,
		fa,
	)
}
