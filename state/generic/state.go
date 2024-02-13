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
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
	FC "github.com/IBM/fp-go/internal/functor"
	P "github.com/IBM/fp-go/pair"
)

var (
	undefined any = struct{}{}
)

func Get[GA ~func(S) P.Pair[S, S], S any]() GA {
	return P.Of[S]
}

func Gets[GA ~func(S) P.Pair[A, S], FCT ~func(S) A, A, S any](f FCT) GA {
	return func(s S) P.Pair[A, S] {
		return P.MakePair(f(s), s)
	}
}

func Put[GA ~func(S) P.Pair[any, S], S any]() GA {
	return F.Bind1st(P.MakePair[any, S], undefined)
}

func Modify[GA ~func(S) P.Pair[any, S], FCT ~func(S) S, S any](f FCT) GA {
	return F.Flow2(
		f,
		F.Bind1st(P.MakePair[any, S], undefined),
	)
}

func Of[GA ~func(S) P.Pair[A, S], S, A any](a A) GA {
	return F.Bind1st(P.MakePair[A, S], a)
}

func MonadMap[GB ~func(S) P.Pair[B, S], GA ~func(S) P.Pair[A, S], FCT ~func(A) B, S, A, B any](fa GA, f FCT) GB {
	return func(s S) P.Pair[B, S] {
		p2 := fa(s)
		return P.MakePair(f(P.Head(p2)), P.Tail(p2))
	}
}

func Map[GB ~func(S) P.Pair[B, S], GA ~func(S) P.Pair[A, S], FCT ~func(A) B, S, A, B any](f FCT) func(GA) GB {
	return F.Bind2nd(MonadMap[GB, GA, FCT, S, A, B], f)
}

func MonadChain[GB ~func(S) P.Pair[B, S], GA ~func(S) P.Pair[A, S], FCT ~func(A) GB, S, A, B any](fa GA, f FCT) GB {
	return func(s S) P.Pair[B, S] {
		a := fa(s)
		return f(P.Head(a))(P.Tail(a))
	}
}

func Chain[GB ~func(S) P.Pair[B, S], GA ~func(S) P.Pair[A, S], FCT ~func(A) GB, S, A, B any](f FCT) func(GA) GB {
	return F.Bind2nd(MonadChain[GB, GA, FCT, S, A, B], f)
}

func MonadAp[GB ~func(S) P.Pair[B, S], GAB ~func(S) P.Pair[func(A) B, S], GA ~func(S) P.Pair[A, S], S, A, B any](fab GAB, fa GA) GB {
	return func(s S) P.Pair[B, S] {
		f := fab(s)
		a := fa(P.Tail(f))

		return P.MakePair(P.Head(f)(P.Head(a)), P.Tail(a))
	}
}

func Ap[GB ~func(S) P.Pair[B, S], GAB ~func(S) P.Pair[func(A) B, S], GA ~func(S) P.Pair[A, S], S, A, B any](ga GA) func(GAB) GB {
	return F.Bind2nd(MonadAp[GB, GAB, GA, S, A, B], ga)
}

func MonadChainFirst[GB ~func(S) P.Pair[B, S], GA ~func(S) P.Pair[A, S], FCT ~func(A) GB, S, A, B any](ma GA, f FCT) GA {
	return C.MonadChainFirst(
		MonadChain[GA, GA, func(A) GA],
		MonadMap[GA, GB, func(B) A],
		ma,
		f,
	)
}

func ChainFirst[GB ~func(S) P.Pair[B, S], GA ~func(S) P.Pair[A, S], FCT ~func(A) GB, S, A, B any](f FCT) func(GA) GA {
	return C.ChainFirst(
		Chain[GA, GA, func(A) GA],
		Map[GA, GB, func(B) A],
		f,
	)
}

func Flatten[GAA ~func(S) P.Pair[GA, S], GA ~func(S) P.Pair[A, S], S, A any](mma GAA) GA {
	return MonadChain[GA, GAA, func(GA) GA](mma, F.Identity[GA])
}

func Execute[GA ~func(S) P.Pair[A, S], S, A any](s S) func(GA) S {
	return func(fa GA) S {
		return P.Tail(fa(s))
	}
}

func Evaluate[GA ~func(S) P.Pair[A, S], S, A any](s S) func(GA) A {
	return func(fa GA) A {
		return P.Head(fa(s))
	}
}

func MonadFlap[FAB ~func(A) B, GFAB ~func(S) P.Pair[FAB, S], GB ~func(S) P.Pair[B, S], S, A, B any](fab GFAB, a A) GB {
	return FC.MonadFlap(
		MonadMap[GB, GFAB, func(FAB) B],
		fab,
		a)
}

func Flap[FAB ~func(A) B, GFAB ~func(S) P.Pair[FAB, S], GB ~func(S) P.Pair[B, S], S, A, B any](a A) func(GFAB) GB {
	return FC.Flap(Map[GB, GFAB, func(FAB) B], a)
}
