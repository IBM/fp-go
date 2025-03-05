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
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[E, S any](
	empty S,
) IOEither[E, S] {
	return Of[E](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) IOEither[E, T],
) Operator[E, S1, S2] {
	return chain.Bind(
		Chain[E, S1, S2],
		Map[E, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[E, S1, S2] {
	return functor.Let(
		Map[E, S1, S2],
		setter,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[E, S1, S2] {
	return functor.LetTo(
		Map[E, S1, S2],
		setter,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[E, S1, T any](
	setter func(T) S1,
) Operator[E, T, S1] {
	return chain.BindTo(
		Map[E, T, S1],
		setter,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOEither[E, T],
) Operator[E, S1, S2] {
	return apply.ApS(
		Ap[S2, E, T],
		Map[E, S1, func(T) S2],
		setter,
		fa,
	)
}
