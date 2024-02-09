// Copyright (c) 2024 IBM Corp.
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

package pair

import (
	"github.com/IBM/fp-go/internal/applicative"
	"github.com/IBM/fp-go/internal/functor"
	"github.com/IBM/fp-go/internal/monad"
	"github.com/IBM/fp-go/internal/pointed"
	M "github.com/IBM/fp-go/monoid"
	Sg "github.com/IBM/fp-go/semigroup"
)

type (
	pairPointedHead[A, B any] struct {
		m M.Monoid[B]
	}

	pairFunctorHead[A, B, A1 any] struct {
	}

	pairApplicativeHead[A, B, A1 any] struct {
		s Sg.Semigroup[B]
		m M.Monoid[B]
	}

	pairMonadHead[A, B, A1 any] struct {
		s Sg.Semigroup[B]
		m M.Monoid[B]
	}

	pairPointedTail[A, B any] struct {
		m M.Monoid[A]
	}

	pairFunctorTail[A, B, B1 any] struct {
	}

	pairApplicativeTail[A, B, B1 any] struct {
		s Sg.Semigroup[A]
		m M.Monoid[A]
	}

	pairMonadTail[A, B, B1 any] struct {
		s Sg.Semigroup[A]
		m M.Monoid[A]
	}
)

func (o *pairMonadHead[A, B, A1]) Of(a A) Pair[A, B] {
	return MakePair(a, o.m.Empty())
}

func (o *pairMonadHead[A, B, A1]) Map(f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return Map[B](f)
}

func (o *pairMonadHead[A, B, A1]) Chain(f func(A) Pair[A1, B]) func(Pair[A, B]) Pair[A1, B] {
	return Chain[B, A, A1](o.s, f)
}

func (o *pairMonadHead[A, B, A1]) Ap(fa Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B] {
	return Ap[B, A, A1](o.s, fa)
}

func (o *pairPointedHead[A, B]) Of(a A) Pair[A, B] {
	return MakePair(a, o.m.Empty())
}

func (o *pairFunctorHead[A, B, A1]) Map(f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return Map[B, A, A1](f)
}

func (o *pairApplicativeHead[A, B, A1]) Map(f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return Map[B, A, A1](f)
}

func (o *pairApplicativeHead[A, B, A1]) Ap(fa Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B] {
	return Ap[B, A, A1](o.s, fa)
}

func (o *pairApplicativeHead[A, B, A1]) Of(a A) Pair[A, B] {
	return MakePair(a, o.m.Empty())
}

// Monad implements the monadic operations for [Pair]
func Monad[A, B, A1 any](m M.Monoid[B]) monad.Monad[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]] {
	return &pairMonadHead[A, B, A1]{s: M.ToSemigroup(m), m: m}
}

// Pointed implements the pointed operations for [Pair]
func Pointed[A, B any](m M.Monoid[B]) pointed.Pointed[A, Pair[A, B]] {
	return &pairPointedHead[A, B]{m: m}
}

// Functor implements the functor operations for [Pair]
func Functor[A, B, A1 any]() functor.Functor[A, A1, Pair[A, B], Pair[A1, B]] {
	return &pairFunctorHead[A, B, A1]{}
}

// Applicative implements the applicative operations for [Pair]
func Applicative[A, B, A1 any](m M.Monoid[B]) applicative.Applicative[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]] {
	return &pairApplicativeHead[A, B, A1]{s: M.ToSemigroup(m), m: m}
}

// MonadHead implements the monadic operations for [Pair]
func MonadHead[A, B, A1 any](m M.Monoid[B]) monad.Monad[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]] {
	return Monad[A, B, A1](m)
}

// PointedHead implements the pointed operations for [Pair]
func PointedHead[A, B any](m M.Monoid[B]) pointed.Pointed[A, Pair[A, B]] {
	return PointedHead[A, B](m)
}

// FunctorHead implements the functor operations for [Pair]
func FunctorHead[A, B, A1 any]() functor.Functor[A, A1, Pair[A, B], Pair[A1, B]] {
	return Functor[A, B, A1]()
}

// ApplicativeHead implements the applicative operations for [Pair]
func ApplicativeHead[A, B, A1 any](m M.Monoid[B]) applicative.Applicative[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]] {
	return Applicative[A, B, A1](m)
}

func (o *pairMonadTail[A, B, B1]) Of(b B) Pair[A, B] {
	return MakePair(o.m.Empty(), b)
}

func (o *pairMonadTail[A, B, B1]) Map(f func(B) B1) func(Pair[A, B]) Pair[A, B1] {
	return MapTail[A, B, B1](f)
}

func (o *pairMonadTail[A, B, B1]) Chain(f func(B) Pair[A, B1]) func(Pair[A, B]) Pair[A, B1] {
	return ChainTail[A, B, B1](o.s, f)
}

func (o *pairMonadTail[A, B, B1]) Ap(fa Pair[A, B]) func(Pair[A, func(B) B1]) Pair[A, B1] {
	return ApTail[A, B, B1](o.s, fa)
}

func (o *pairPointedTail[A, B]) Of(b B) Pair[A, B] {
	return MakePair(o.m.Empty(), b)
}

func (o *pairFunctorTail[A, B, B1]) Map(f func(B) B1) func(Pair[A, B]) Pair[A, B1] {
	return MapTail[A, B, B1](f)
}

func (o *pairApplicativeTail[A, B, B1]) Map(f func(B) B1) func(Pair[A, B]) Pair[A, B1] {
	return MapTail[A, B, B1](f)
}

func (o *pairApplicativeTail[A, B, B1]) Ap(fa Pair[A, B]) func(Pair[A, func(B) B1]) Pair[A, B1] {
	return ApTail[A, B, B1](o.s, fa)
}

func (o *pairApplicativeTail[A, B, B1]) Of(b B) Pair[A, B] {
	return MakePair(o.m.Empty(), b)
}

// MonadTail implements the monadic operations for [Pair]
func MonadTail[B, A, B1 any](m M.Monoid[A]) monad.Monad[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]] {
	return &pairMonadTail[A, B, B1]{s: M.ToSemigroup(m), m: m}
}

// PointedTail implements the pointed operations for [Pair]
func PointedTail[B, A any](m M.Monoid[A]) pointed.Pointed[B, Pair[A, B]] {
	return &pairPointedTail[A, B]{m: m}
}

// FunctorTail implements the functor operations for [Pair]
func FunctorTail[B, A, B1 any]() functor.Functor[B, B1, Pair[A, B], Pair[A, B1]] {
	return &pairFunctorTail[A, B, B1]{}
}

// ApplicativeTail implements the applicative operations for [Pair]
func ApplicativeTail[B, A, B1 any](m M.Monoid[A]) applicative.Applicative[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]] {
	return &pairApplicativeTail[A, B, B1]{s: M.ToSemigroup(m), m: m}
}
