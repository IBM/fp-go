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
	P "github.com/IBM/fp-go/pair"
	SG "github.com/IBM/fp-go/semigroup"
)

func Tell[GA ~func() P.Pair[any, W], W any](w W) GA {
	return IO.Of[GA](P.MakePair[any](w, w))
}

func Of[GA ~func() P.Pair[A, W], W, A any](m M.Monoid[W], a A) GA {
	return IO.Of[GA](P.MakePair(a, m.Empty()))
}

// Listen modifies the result to include the changes to the accumulator
func Listen[GA ~func() P.Pair[A, W], GTA ~func() P.Pair[P.Pair[A, W], W], W, A any](fa GA) GTA {
	return func() P.Pair[P.Pair[A, W], W] {
		t := fa()
		return P.MakePair(t, P.Tail(t))
	}
}

// Pass applies the returned function to the accumulator
func Pass[GFA ~func() P.Pair[P.Pair[A, FCT], W], GA ~func() P.Pair[A, W], FCT ~func(W) W, W, A any](fa GFA) GA {
	return func() P.Pair[A, W] {
		t := fa()
		a := P.Head(t)
		return P.MakePair(P.Head(a), P.Tail(a)(P.Tail(t)))
	}
}

func MonadMap[GB ~func() P.Pair[B, W], GA ~func() P.Pair[A, W], FCT ~func(A) B, W, A, B any](fa GA, f FCT) GB {
	return IO.MonadMap[GA, GB](fa, P.Map[W](f))
}

func Map[GB ~func() P.Pair[B, W], GA ~func() P.Pair[A, W], FCT ~func(A) B, W, A, B any](f FCT) func(GA) GB {
	return IO.Map[GA, GB](P.Map[W](f))
}

func MonadChain[GB ~func() P.Pair[B, W], GA ~func() P.Pair[A, W], FCT ~func(A) GB, W, A, B any](s SG.Semigroup[W], fa GA, f FCT) GB {
	return func() P.Pair[B, W] {
		a := fa()
		b := f(P.Head(a))()

		return P.MakePair(P.Head(b), s.Concat(P.Tail(a), P.Tail(b)))
	}
}

func Chain[GB ~func() P.Pair[B, W], GA ~func() P.Pair[A, W], FCT ~func(A) GB, W, A, B any](s SG.Semigroup[W], f FCT) func(GA) GB {
	return func(fa GA) GB {
		return MonadChain(s, fa, f)
	}
}

func MonadAp[GB ~func() P.Pair[B, W], GAB ~func() P.Pair[func(A) B, W], GA ~func() P.Pair[A, W], W, A, B any](s SG.Semigroup[W], fab GAB, fa GA) GB {
	return func() P.Pair[B, W] {
		f := fab()
		a := fa()

		return P.MakePair(P.Head(f)(P.Head(a)), s.Concat(P.Tail(f), P.Tail(a)))
	}
}

func Ap[GB ~func() P.Pair[B, W], GAB ~func() P.Pair[func(A) B, W], GA ~func() P.Pair[A, W], W, A, B any](s SG.Semigroup[W], ga GA) func(GAB) GB {
	return func(fab GAB) GB {
		return MonadAp[GB](s, fab, ga)
	}
}

func MonadChainFirst[GB ~func() P.Pair[B, W], GA ~func() P.Pair[A, W], FCT ~func(A) GB, W, A, B any](s SG.Semigroup[W], ma GA, f FCT) GA {
	return C.MonadChainFirst(
		F.Bind1of3(MonadChain[GA, GA, func(A) GA])(s),
		MonadMap[GA, GB, func(B) A],
		ma,
		f,
	)
}

func ChainFirst[GB ~func() P.Pair[B, W], GA ~func() P.Pair[A, W], FCT ~func(A) GB, W, A, B any](s SG.Semigroup[W], f FCT) func(GA) GA {
	return C.ChainFirst(
		F.Bind1st(Chain[GA, GA, func(A) GA], s),
		Map[GA, GB, func(B) A],
		f,
	)
}

func Flatten[GAA ~func() P.Pair[GA, W], GA ~func() P.Pair[A, W], W, A any](s SG.Semigroup[W], mma GAA) GA {
	return MonadChain[GA, GAA, func(GA) GA](s, mma, F.Identity[GA])
}

func Execute[GA ~func() P.Pair[A, W], W, A any](fa GA) W {
	return P.Tail(fa())
}

func Evaluate[GA ~func() P.Pair[A, W], W, A any](fa GA) A {
	return P.Head(fa())
}

// MonadCensor modifies the final accumulator value by applying a function
func MonadCensor[GA ~func() P.Pair[A, W], FCT ~func(W) W, W, A any](fa GA, f FCT) GA {
	return IO.MonadMap[GA, GA](fa, P.MapTail[A](f))
}

// Censor modifies the final accumulator value by applying a function
func Censor[GA ~func() P.Pair[A, W], FCT ~func(W) W, W, A any](f FCT) func(GA) GA {
	return IO.Map[GA, GA](P.MapTail[A](f))
}

// MonadListens projects a value from modifications made to the accumulator during an action
func MonadListens[GA ~func() P.Pair[A, W], GAB ~func() P.Pair[P.Pair[A, B], W], FCT ~func(W) B, W, A, B any](fa GA, f FCT) GAB {
	return func() P.Pair[P.Pair[A, B], W] {
		a := fa()
		t := P.Tail(a)
		return P.MakePair(P.MakePair(P.Head(a), f(t)), t)
	}
}

// Listens projects a value from modifications made to the accumulator during an action
func Listens[GA ~func() P.Pair[A, W], GAB ~func() P.Pair[P.Pair[A, B], W], FCT ~func(W) B, W, A, B any](f FCT) func(GA) GAB {
	return F.Bind2nd(MonadListens[GA, GAB, FCT], f)
}

func MonadFlap[FAB ~func(A) B, GFAB ~func() P.Pair[FAB, W], GB ~func() P.Pair[B, W], W, A, B any](fab GFAB, a A) GB {
	return FC.MonadFlap(
		MonadMap[GB, GFAB, func(FAB) B],
		fab,
		a)
}

func Flap[FAB ~func(A) B, GFAB ~func() P.Pair[FAB, W], GB ~func() P.Pair[B, W], W, A, B any](a A) func(GFAB) GB {
	return FC.Flap(Map[GB, GFAB, func(FAB) B], a)
}
