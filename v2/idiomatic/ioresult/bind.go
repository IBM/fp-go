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

package ioresult

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	L "github.com/IBM/fp-go/v2/optics/lens"
)

// Do starts a do-notation computation with an initial state.
// This is the entry point for building complex computations using the do-notation style.
//
//go:inline
func Do[S any](
	empty S,
) IOResult[S] {
	return Of(empty)
}

// Bind adds a computation step in do-notation, extending the state with a new field.
// The setter function determines how the new value is added to the state.
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return chain.Bind(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		f,
	)
}

// Let adds a pure transformation step in do-notation.
// Unlike Bind, the function does not return an IOResult, making it suitable for pure computations.
//
//go:inline
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return functor.Let(
		Map[S1, S2],
		setter,
		f,
	)
}

// LetTo adds a constant value to the state in do-notation.
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return functor.LetTo(
		Map[S1, S2],
		setter,
		b,
	)
}

// BindTo wraps a value in an initial state structure.
// This is typically the first operation after creating an IOResult in do-notation.
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return chain.BindTo(
		Map[T, S1],
		setter,
	)
}

// ApS applies an IOResult to extend the state in do-notation.
// This is used to add independent computations that don't depend on previous results.
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOResult[T],
) Operator[S1, S2] {
	return apply.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL applies an IOResult using a lens to update a specific field in the state.
func ApSL[S, T any](
	lens L.Lens[S, T],
	fa IOResult[T],
) Operator[S, S] {
	return ApS(lens.Set, fa)
}

// BindL binds a computation using a lens to focus on a specific field.
func BindL[S, T any](
	lens L.Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// LetL applies a pure transformation using a lens to update a specific field.
func LetL[S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Operator[S, S] {
	return Let(lens.Set, F.Flow2(lens.Get, f))
}

// LetToL sets a field to a constant value using a lens.
func LetToL[S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[S, S] {
	return LetTo(lens.Set, b)
}
