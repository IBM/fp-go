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
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
	FC "github.com/IBM/fp-go/internal/functor"
	IO "github.com/IBM/fp-go/io/generic"
	M "github.com/IBM/fp-go/monoid"
	SG "github.com/IBM/fp-go/semigroup"
	T "github.com/IBM/fp-go/tuple"
)

func Tell[GA ~func() T.Tuple3[any, W, SG.Semigroup[W]], W any](s SG.Semigroup[W]) func(W) GA {
	return F.Flow2(
		F.Bind13of3(T.MakeTuple3[any, W, SG.Semigroup[W]])(nil, s),
		IO.Of[GA],
	)
}

func Of[GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, A any](m M.Monoid[W]) func(A) GA {
	return F.Flow2(
		F.Bind23of3(T.MakeTuple3[A, W, SG.Semigroup[W]])(m.Empty(), M.ToSemigroup(m)),
		IO.Of[GA],
	)
}

// Listen modifies the result to include the changes to the accumulator
func Listen[GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], GTA ~func() T.Tuple3[T.Tuple2[A, W], W, SG.Semigroup[W]], W, A any](fa GA) GTA {
	return func() T.Tuple3[T.Tuple2[A, W], W, SG.Semigroup[W]] {
		t := fa()
		return T.MakeTuple3(T.MakeTuple2(t.F1, t.F2), t.F2, t.F3)
	}
}

// Pass applies the returned function to the accumulator
func Pass[GFA ~func() T.Tuple3[T.Tuple2[A, FCT], W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(W) W, W, A any](fa GFA) GA {
	return func() T.Tuple3[A, W, SG.Semigroup[W]] {
		t := fa()
		return T.MakeTuple3(t.F1.F1, t.F1.F2(t.F2), t.F3)
	}
}

func MonadMap[GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(A) B, W, A, B any](fa GA, f FCT) GB {
	return IO.MonadMap[GA, GB](fa, T.Map3(f, F.Identity[W], F.Identity[SG.Semigroup[W]]))
}

func Map[GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(A) B, W, A, B any](f FCT) func(GA) GB {
	return IO.Map[GA, GB](T.Map3(f, F.Identity[W], F.Identity[SG.Semigroup[W]]))
}

func MonadChain[GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(A) GB, W, A, B any](fa GA, f FCT) GB {
	return func() T.Tuple3[B, W, SG.Semigroup[W]] {
		a := fa()
		b := f(a.F1)()

		return T.MakeTuple3(b.F1, b.F3.Concat(a.F2, b.F2), b.F3)
	}
}

func Chain[GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(A) GB, W, A, B any](f FCT) func(GA) GB {
	return F.Bind2nd(MonadChain[GB, GA, FCT, W, A, B], f)
}

func MonadAp[GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], GAB ~func() T.Tuple3[func(A) B, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, A, B any](fab GAB, fa GA) GB {
	return func() T.Tuple3[B, W, SG.Semigroup[W]] {
		f := fab()
		a := fa()

		return T.MakeTuple3(f.F1(a.F1), f.F3.Concat(f.F2, a.F2), f.F3)
	}
}

func Ap[GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], GAB ~func() T.Tuple3[func(A) B, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, A, B any](ga GA) func(GAB) GB {
	return F.Bind2nd(MonadAp[GB, GAB, GA], ga)
}

func MonadChainFirst[GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(A) GB, W, A, B any](ma GA, f FCT) GA {
	return C.MonadChainFirst(
		MonadChain[GA, GA, func(A) GA],
		MonadMap[GA, GB, func(B) A],
		ma,
		f,
	)
}

func ChainFirst[GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(A) GB, W, A, B any](f FCT) func(GA) GA {
	return C.ChainFirst(
		Chain[GA, GA, func(A) GA],
		Map[GA, GB, func(B) A],
		f,
	)
}

func Flatten[GAA ~func() T.Tuple3[GA, W, SG.Semigroup[W]], GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, A any](mma GAA) GA {
	return MonadChain[GA, GAA, func(GA) GA](mma, F.Identity[GA])
}

func Execute[GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, A any](fa GA) W {
	return fa().F2
}

func Evaluate[GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], W, A any](fa GA) A {
	return fa().F1
}

// MonadCensor modifies the final accumulator value by applying a function
func MonadCensor[GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(W) W, W, A any](fa GA, f FCT) GA {
	return IO.MonadMap[GA, GA](fa, T.Map3(F.Identity[A], f, F.Identity[SG.Semigroup[W]]))
}

// Censor modifies the final accumulator value by applying a function
func Censor[GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], FCT ~func(W) W, W, A any](f FCT) func(GA) GA {
	return IO.Map[GA, GA](T.Map3(F.Identity[A], f, F.Identity[SG.Semigroup[W]]))
}

// MonadListens projects a value from modifications made to the accumulator during an action
func MonadListens[GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], GAB ~func() T.Tuple3[T.Tuple2[A, B], W, SG.Semigroup[W]], FCT ~func(W) B, W, A, B any](fa GA, f FCT) GAB {
	return func() T.Tuple3[T.Tuple2[A, B], W, SG.Semigroup[W]] {
		a := fa()
		return T.MakeTuple3(T.MakeTuple2(a.F1, f(a.F2)), a.F2, a.F3)
	}
}

// Listens projects a value from modifications made to the accumulator during an action
func Listens[GA ~func() T.Tuple3[A, W, SG.Semigroup[W]], GAB ~func() T.Tuple3[T.Tuple2[A, B], W, SG.Semigroup[W]], FCT ~func(W) B, W, A, B any](f FCT) func(GA) GAB {
	return F.Bind2nd(MonadListens[GA, GAB, FCT], f)
}

func MonadFlap[FAB ~func(A) B, GFAB ~func() T.Tuple3[FAB, W, SG.Semigroup[W]], GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], W, A, B any](fab GFAB, a A) GB {
	return FC.MonadFlap(
		MonadMap[GB, GFAB, func(FAB) B],
		fab,
		a)
}

func Flap[FAB ~func(A) B, GFAB ~func() T.Tuple3[FAB, W, SG.Semigroup[W]], GB ~func() T.Tuple3[B, W, SG.Semigroup[W]], W, A, B any](a A) func(GFAB) GB {
	return FC.Flap(Map[GB, GFAB, func(FAB) B], a)
}
