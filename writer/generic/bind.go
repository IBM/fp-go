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

package generic

import (
	"github.com/IBM/fp-go/internal/apply"
	C "github.com/IBM/fp-go/internal/chain"
	F "github.com/IBM/fp-go/internal/functor"
	M "github.com/IBM/fp-go/monoid"
	SG "github.com/IBM/fp-go/semigroup"
	T "github.com/IBM/fp-go/tuple"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[GS ~func() T.Tuple3[S, W, SG.Semigroup[W]], W, S any](m M.Monoid[W]) func(S) GS {
	return Of[GS, W, S](m)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[GS1 ~func() T.Tuple3[S1, W, SG.Semigroup[W]], GS2 ~func() T.Tuple3[S2, W, SG.Semigroup[W]], GT ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, S1, S2, A any](
	setter func(A) func(S1) S2,
	f func(S1) GT,
) func(GS1) GS2 {
	return C.Bind(
		Chain[GS2, GS1, func(S1) GS2, W, S1, S2],
		Map[GS2, GT, func(A) S2, W, A, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[GS1 ~func() T.Tuple3[S1, W, SG.Semigroup[W]], GS2 ~func() T.Tuple3[S2, W, SG.Semigroup[W]], W, S1, S2, A any](
	key func(A) func(S1) S2,
	f func(S1) A,
) func(GS1) GS2 {
	return F.Let(
		Map[GS2, GS1, func(S1) S2, W, S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[GS1 ~func() T.Tuple3[S1, W, SG.Semigroup[W]], GS2 ~func() T.Tuple3[S2, W, SG.Semigroup[W]], W, S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(GS1) GS2 {
	return F.LetTo(
		Map[GS2, GS1, func(S1) S2, W, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[GS1 ~func() T.Tuple3[S1, W, SG.Semigroup[W]], GT ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, S1, A any](
	setter func(A) S1,
) func(GT) GS1 {
	return C.BindTo(
		Map[GS1, GT, func(A) S1, W, A, S1],
		setter,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[GS1 ~func() T.Tuple3[S1, W, SG.Semigroup[W]], GS2 ~func() T.Tuple3[S2, W, SG.Semigroup[W]], GT ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, S1, S2, A any](
	setter func(A) func(S1) S2,
	fa GT,
) func(GS1) GS2 {
	return apply.ApS(
		Ap[GS2, func() T.Tuple3[func(A) S2, W, SG.Semigroup[W]], GT, W, A, S2],
		Map[func() T.Tuple3[func(A) S2, W, SG.Semigroup[W]], GS1, func(S1) func(A) S2],
		setter,
		fa,
	)
}
