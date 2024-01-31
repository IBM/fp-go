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
func Do[GS ~func(R) ET.Either[E, S], R, E, S any](
	empty S,
) GS {
	return Of[GS, E, R, S](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[GS1 ~func(R) ET.Either[E, S1], GS2 ~func(R) ET.Either[E, S2], GT ~func(R) ET.Either[E, T], R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) GT,
) func(GS1) GS2 {
	return C.Bind(
		Chain[GS1, GS2, E, R, S1, S2],
		Map[GT, GS2, E, R, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[GS1 ~func(R) ET.Either[E, S1], GS2 ~func(R) ET.Either[E, S2], R, E, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) func(GS1) GS2 {
	return F.Let(
		Map[GS1, GS2, E, R, S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[GS1 ~func(R) ET.Either[E, S1], GS2 ~func(R) ET.Either[E, S2], R, E, S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(GS1) GS2 {
	return F.LetTo(
		Map[GS1, GS2, E, R, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[GS1 ~func(R) ET.Either[E, S1], GT ~func(R) ET.Either[E, T], R, E, S1, T any](
	setter func(T) S1,
) func(GT) GS1 {
	return C.BindTo(
		Map[GT, GS1, E, R, T, S1],
		setter,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering the context and the value concurrently
func ApS[GS1 ~func(R) ET.Either[E, S1], GS2 ~func(R) ET.Either[E, S2], GT ~func(R) ET.Either[E, T], R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa GT,
) func(GS1) GS2 {
	return A.ApS(
		Ap[GT, GS2, func(R) ET.Either[E, func(T) S2], E, R, T, S2],
		Map[GS1, func(R) ET.Either[E, func(T) S2], E, R, S1, func(T) S2],
		setter,
		fa,
	)
}
