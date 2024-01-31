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
	ET "github.com/IBM/fp-go/either"
	A "github.com/IBM/fp-go/internal/apply"
	C "github.com/IBM/fp-go/internal/chain"
	F "github.com/IBM/fp-go/internal/functor"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[GRS ~func(R) GS, GS ~func() ET.Either[E, S], R, E, S any](
	empty S,
) GRS {
	return Of[GRS, GS, R, E, S](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[GRS1 ~func(R) GS1, GRS2 ~func(R) GS2, GRT ~func(R) GT, GS1 ~func() ET.Either[E, S1], GS2 ~func() ET.Either[E, S2], GT ~func() ET.Either[E, T], R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) GRT,
) func(GRS1) GRS2 {
	return C.Bind(
		Chain[GRS1, GRS2, GS1, GS2, R, E, S1, S2],
		Map[GRT, GRS2, GT, GS2, R, E, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[GRS1 ~func(R) GS1, GRS2 ~func(R) GS2, GS1 ~func() ET.Either[E, S1], GS2 ~func() ET.Either[E, S2], R, E, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) func(GRS1) GRS2 {
	return F.Let(
		Map[GRS1, GRS2, GS1, GS2, R, E, S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[GRS1 ~func(R) GS1, GRS2 ~func(R) GS2, GS1 ~func() ET.Either[E, S1], GS2 ~func() ET.Either[E, S2], R, E, S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(GRS1) GRS2 {
	return F.LetTo(
		Map[GRS1, GRS2, GS1, GS2, R, E, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[GRS1 ~func(R) GS1, GRT ~func(R) GT, GS1 ~func() ET.Either[E, S1], GT ~func() ET.Either[E, T], R, E, S1, T any](
	setter func(T) S1,
) func(GRT) GRS1 {
	return C.BindTo(
		Map[GRT, GRS1, GT, GS1, R, E, T, S1],
		setter,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[GRTS1 ~func(R) GTS1, GRS1 ~func(R) GS1, GRS2 ~func(R) GS2, GRT ~func(R) GT, GTS1 ~func() ET.Either[E, func(T) S2], GS1 ~func() ET.Either[E, S1], GS2 ~func() ET.Either[E, S2], GT ~func() ET.Either[E, T], R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa GRT,
) func(GRS1) GRS2 {
	return A.ApS(
		Ap[GRT, GRS2, GRTS1, GT, GS2, GTS1, R, E, T, S2],
		Map[GRS1, GRTS1, GS1, GTS1, R, E, S1, func(T) S2],
		setter,
		fa,
	)
}
