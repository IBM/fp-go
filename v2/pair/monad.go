// Copyright (c) 2024 - 2025 IBM Corp.
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
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/semigroup"
)

type (
	pairPointedHead[A, B any] struct {
		m monoid.Monoid[B]
	}

	pairFunctorHead[A, B, A1 any] struct {
	}

	pairApplicativeHead[A, B, A1 any] struct {
		s Semigroup[B]
		m monoid.Monoid[B]
	}

	pairMonadHead[A, B, A1 any] struct {
		s Semigroup[B]
		m monoid.Monoid[B]
	}

	pairPointedTail[A, B any] struct {
		m monoid.Monoid[A]
	}

	pairFunctorTail[A, B, B1 any] struct {
	}

	pairApplicativeTail[A, B, B1 any] struct {
		s semigroup.Semigroup[A]
		m monoid.Monoid[A]
	}

	pairMonadTail[A, B, B1 any] struct {
		s semigroup.Semigroup[A]
		m monoid.Monoid[A]
	}
)

func (o *pairMonadHead[A, B, A1]) Of(a A) Pair[A, B] {
	return MakePair(a, o.m.Empty())
}

func (o *pairMonadHead[A, B, A1]) Map(f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return MapHead[B](f)
}

func (o *pairMonadHead[A, B, A1]) Chain(f func(A) Pair[A1, B]) func(Pair[A, B]) Pair[A1, B] {
	return ChainHead(o.s, f)
}

func (o *pairMonadHead[A, B, A1]) Ap(fa Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B] {
	return ApHead[B, A, A1](o.s, fa)
}

func (o *pairPointedHead[A, B]) Of(a A) Pair[A, B] {
	return MakePair(a, o.m.Empty())
}

func (o *pairFunctorHead[A, B, A1]) Map(f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return MapHead[B](f)
}

func (o *pairApplicativeHead[A, B, A1]) Map(f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return MapHead[B](f)
}

func (o *pairApplicativeHead[A, B, A1]) Ap(fa Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B] {
	return ApHead[B, A, A1](o.s, fa)
}

func (o *pairApplicativeHead[A, B, A1]) Of(a A) Pair[A, B] {
	return MakePair(a, o.m.Empty())
}

// MonadHead implements the monadic operations for [Pair] operating on the head value.
// Requires a monoid for the tail type to provide an identity element for the Of operation
// and a semigroup for combining tail values in Chain and Ap operations.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	stringMonoid := M.MakeMonoid(
//	    func(a, b string) string { return a + b },
//	    "",
//	)
//	monad := pair.MonadHead[int, string, int](stringMonoid)
//	p := monad.Of(42)  // Pair[int, string]{42, ""}
func MonadHead[A, B, A1 any](m monoid.Monoid[B]) monad.Monad[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]] {
	return &pairMonadHead[A, B, A1]{s: monoid.ToSemigroup(m), m: m}
}

// PointedHead implements the pointed operations for [Pair] operating on the head value.
// Requires a monoid for the tail type to provide an identity element.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	stringMonoid := M.MakeMonoid(
//	    func(a, b string) string { return a + b },
//	    "",
//	)
//	pointed := pair.PointedHead[int, string](stringMonoid)
//	p := pointed.Of(42)  // Pair[int, string]{42, ""}
func PointedHead[A, B any](m monoid.Monoid[B]) pointed.Pointed[A, Pair[A, B]] {
	return &pairPointedHead[A, B]{m: m}
}

// FunctorHead implements the functor operations for [Pair] operating on the head value.
//
// Example:
//
//	functor := pair.FunctorHead[int, string, string]()
//	mapper := functor.Map(func(n int) string { return fmt.Sprintf("%d", n) })
//	p := pair.MakePair(42, "hello")
//	p2 := mapper(p)  // Pair[string, string]{"42", "hello"}
func FunctorHead[A, B, A1 any]() functor.Functor[A, A1, Pair[A, B], Pair[A1, B]] {
	return &pairFunctorHead[A, B, A1]{}
}

// ApplicativeHead implements the applicative operations for [Pair] operating on the head value.
// Requires a monoid for the tail type to provide an identity element for the Of operation
// and a semigroup for combining tail values in the Ap operation.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	stringMonoid := M.MakeMonoid(
//	    func(a, b string) string { return a + b },
//	    "",
//	)
//	applicative := pair.ApplicativeHead[int, string, string](stringMonoid)
//	pf := applicative.Of(func(n int) string { return fmt.Sprintf("%d", n) })
//	pv := pair.MakePair(42, "!")
//	result := applicative.Ap(pv)(pf)  // Pair[string, string]{"42", "!"}
func ApplicativeHead[A, B, A1 any](m monoid.Monoid[B]) applicative.Applicative[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]] {
	return &pairApplicativeHead[A, B, A1]{s: monoid.ToSemigroup(m), m: m}
}

func (o *pairMonadTail[A, B, B1]) Of(b B) Pair[A, B] {
	return MakePair(o.m.Empty(), b)
}

func (o *pairMonadTail[A, B, B1]) Map(f func(B) B1) Operator[A, B, B1] {
	return MapTail[A](f)
}

func (o *pairMonadTail[A, B, B1]) Chain(f Kleisli[A, B, B1]) Operator[A, B, B1] {
	return ChainTail(o.s, f)
}

func (o *pairMonadTail[A, B, B1]) Ap(fa Pair[A, B]) Operator[A, func(B) B1, B1] {
	return ApTail[A, B, B1](o.s, fa)
}

func (o *pairPointedTail[A, B]) Of(b B) Pair[A, B] {
	return MakePair(o.m.Empty(), b)
}

func (o *pairFunctorTail[A, B, B1]) Map(f func(B) B1) Operator[A, B, B1] {
	return MapTail[A](f)
}

func (o *pairApplicativeTail[A, B, B1]) Map(f func(B) B1) Operator[A, B, B1] {
	return MapTail[A](f)
}

func (o *pairApplicativeTail[A, B, B1]) Ap(fa Pair[A, B]) Operator[A, func(B) B1, B1] {
	return ApTail[A, B, B1](o.s, fa)
}

func (o *pairApplicativeTail[A, B, B1]) Of(b B) Pair[A, B] {
	return MakePair(o.m.Empty(), b)
}

// MonadTail implements the monadic operations for [Pair] operating on the tail value.
// Requires a monoid for the head type to provide an identity element for the Of operation
// and a semigroup for combining head values in Chain and Ap operations.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	intSum := M.MonoidSum[int]()
//	monad := pair.MonadTail[string, int, int](intSum)
//	p := monad.Of("hello")  // Pair[int, string]{0, "hello"}
func MonadTail[B, A, B1 any](m monoid.Monoid[A]) monad.Monad[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]] {
	return &pairMonadTail[A, B, B1]{s: monoid.ToSemigroup(m), m: m}
}

// PointedTail implements the pointed operations for [Pair] operating on the tail value.
// Requires a monoid for the head type to provide an identity element.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	intSum := M.MonoidSum[int]()
//	pointed := pair.PointedTail[string, int](intSum)
//	p := pointed.Of("hello")  // Pair[int, string]{0, "hello"}
func PointedTail[B, A any](m monoid.Monoid[A]) pointed.Pointed[B, Pair[A, B]] {
	return &pairPointedTail[A, B]{m: m}
}

// FunctorTail implements the functor operations for [Pair] operating on the tail value.
//
// Example:
//
//	functor := pair.FunctorTail[string, int, int]()
//	mapper := functor.Map(S.Size)
//	p := pair.MakePair(5, "hello")
//	p2 := mapper(p)  // Pair[int, int]{5, 5}
func FunctorTail[B, A, B1 any]() functor.Functor[B, B1, Pair[A, B], Pair[A, B1]] {
	return &pairFunctorTail[A, B, B1]{}
}

// ApplicativeTail implements the applicative operations for [Pair] operating on the tail value.
// Requires a monoid for the head type to provide an identity element for the Of operation
// and a semigroup for combining head values in the Ap operation.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	intSum := M.MonoidSum[int]()
//	applicative := pair.ApplicativeTail[string, int, int](intSum)
//	pf := applicative.Of(S.Size)
//	pv := pair.MakePair(5, "hello")
//	result := applicative.Ap(pv)(pf)  // Pair[int, int]{5, 5}
func ApplicativeTail[B, A, B1 any](m monoid.Monoid[A]) applicative.Applicative[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]] {
	return &pairApplicativeTail[A, B, B1]{s: monoid.ToSemigroup(m), m: m}
}

// Monad implements the monadic operations for [Pair] operating on the tail value (alias for MonadTail).
// This is the default monad instance for Pair.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	intSum := M.MonoidSum[int]()
//	monad := pair.Monad[string, int, int](intSum)
//	p := monad.Of("hello")  // Pair[int, string]{0, "hello"}
func Monad[B, A, B1 any](m monoid.Monoid[A]) monad.Monad[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]] {
	return MonadTail[B, A, B1](m)
}

// Pointed implements the pointed operations for [Pair] operating on the tail value (alias for PointedTail).
// This is the default pointed instance for Pair.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	intSum := M.MonoidSum[int]()
//	pointed := pair.Pointed[string, int](intSum)
//	p := pointed.Of("hello")  // Pair[int, string]{0, "hello"}
func Pointed[B, A any](m monoid.Monoid[A]) pointed.Pointed[B, Pair[A, B]] {
	return PointedTail[B](m)
}

// Functor implements the functor operations for [Pair] operating on the tail value (alias for FunctorTail).
// This is the default functor instance for Pair.
//
// Example:
//
//	functor := pair.Functor[string, int, int]()
//	mapper := functor.Map(S.Size)
//	p := pair.MakePair(5, "hello")
//	p2 := mapper(p)  // Pair[int, int]{5, 5}
func Functor[B, A, B1 any]() functor.Functor[B, B1, Pair[A, B], Pair[A, B1]] {
	return FunctorTail[B, A, B1]()
}

// Applicative implements the applicative operations for [Pair] operating on the tail value (alias for ApplicativeTail).
// This is the default applicative instance for Pair.
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	intSum := M.MonoidSum[int]()
//	applicative := pair.Applicative[string, int, int](intSum)
//	pf := applicative.Of(S.Size)
//	pv := pair.MakePair(5, "hello")
//	result := applicative.Ap(pv)(pf)  // Pair[int, int]{5, 5}
func Applicative[B, A, B1 any](m monoid.Monoid[A]) applicative.Applicative[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]] {
	return ApplicativeTail[B, A, B1](m)
}
