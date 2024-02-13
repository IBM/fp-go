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

package generic

import (
	"github.com/IBM/fp-go/internal/applicative"
	"github.com/IBM/fp-go/internal/functor"
	"github.com/IBM/fp-go/internal/monad"
	"github.com/IBM/fp-go/internal/pointed"
	M "github.com/IBM/fp-go/monoid"
	P "github.com/IBM/fp-go/pair"
	SG "github.com/IBM/fp-go/semigroup"
)

type writerPointed[GA ~func() P.Pair[A, W], W, A any] struct {
	m M.Monoid[W]
}

type writerFunctor[GB ~func() P.Pair[B, W], GA ~func() P.Pair[A, W], W, A, B any] struct{}

type writerApplicative[GB ~func() P.Pair[B, W], GAB ~func() P.Pair[func(A) B, W], GA ~func() P.Pair[A, W], W, A, B any] struct {
	s SG.Semigroup[W]
	m M.Monoid[W]
}

type writerMonad[GB ~func() P.Pair[B, W], GAB ~func() P.Pair[func(A) B, W], GA ~func() P.Pair[A, W], W, A, B any] struct {
	s SG.Semigroup[W]
	m M.Monoid[W]
}

func (o *writerPointed[GA, W, A]) Of(a A) GA {
	return Of[GA](o.m, a)
}

func (o *writerApplicative[GB, GAB, GA, W, A, B]) Of(a A) GA {
	return Of[GA](o.m, a)
}

func (o *writerMonad[GB, GAB, GA, W, A, B]) Of(a A) GA {
	return Of[GA](o.m, a)
}

func (o *writerFunctor[GB, GA, W, A, B]) Map(f func(A) B) func(GA) GB {
	return Map[GB, GA](f)
}

func (o *writerApplicative[GB, GAB, GA, W, A, B]) Map(f func(A) B) func(GA) GB {
	return Map[GB, GA](f)
}

func (o *writerMonad[GB, GAB, GA, W, A, B]) Map(f func(A) B) func(GA) GB {
	return Map[GB, GA](f)
}

func (o *writerMonad[GB, GAB, GA, W, A, B]) Chain(f func(A) GB) func(GA) GB {
	return Chain[GB, GA](o.s, f)
}

func (o *writerApplicative[GB, GAB, GA, W, A, B]) Ap(fa GA) func(GAB) GB {
	return Ap[GB, GAB, GA](o.s, fa)
}

func (o *writerMonad[GB, GAB, GA, W, A, B]) Ap(fa GA) func(GAB) GB {
	return Ap[GB, GAB, GA](o.s, fa)
}

// Pointed implements the pointed operations for [Writer]
func Pointed[GA ~func() P.Pair[A, W], W, A any](m M.Monoid[W]) pointed.Pointed[A, GA] {
	return &writerPointed[GA, W, A]{
		m: m,
	}
}

// Functor implements the functor operations for [Writer]
func Functor[GB ~func() P.Pair[B, W], GA ~func() P.Pair[A, W], W, A, B any]() functor.Functor[A, B, GA, GB] {
	return &writerFunctor[GB, GA, W, A, B]{}
}

// Applicative implements the applicative operations for [Writer]
func Applicative[GB ~func() P.Pair[B, W], GAB ~func() P.Pair[func(A) B, W], GA ~func() P.Pair[A, W], W, A, B any](m M.Monoid[W]) applicative.Applicative[A, B, GA, GB, GAB] {
	return &writerApplicative[GB, GAB, GA, W, A, B]{
		s: M.ToSemigroup(m),
		m: m,
	}
}

// Monad implements the monadic operations for [Writer]
func Monad[GB ~func() P.Pair[B, W], GAB ~func() P.Pair[func(A) B, W], GA ~func() P.Pair[A, W], W, A, B any](m M.Monoid[W]) monad.Monad[A, B, GA, GB, GAB] {
	return &writerMonad[GB, GAB, GA, W, A, B]{
		s: M.ToSemigroup(m),
		m: m,
	}
}
