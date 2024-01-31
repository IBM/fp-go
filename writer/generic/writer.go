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
	IO "github.com/IBM/fp-go/io/generic"
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
	T "github.com/IBM/fp-go/tuple"
)

func Of[GA ~func() T.Tuple2[A, W], W, A any](m M.Monoid[W]) func(A) GA {
	return F.Flow2(
		F.Bind2nd(T.MakeTuple2[A, W], m.Empty()),
		IO.Of[GA],
	)
}

// Listen modifies the result to include the changes to the accumulator
func Listen[GA ~func() T.Tuple2[A, W], GTA ~func() T.Tuple2[T.Tuple2[A, W], W], W, A any](fa GA) GTA {
	return func() T.Tuple2[T.Tuple2[A, W], W] {
		t := fa()
		return T.MakeTuple2(T.MakeTuple2(t.F1, t.F2), t.F2)
	}
}

// Pass applies the returned function to the accumulator
func Pass[GFA ~func() T.Tuple2[T.Tuple2[A, FCT], W], GA ~func() T.Tuple2[A, W], FCT ~func(W) W, W, A any](fa GFA) GA {
	return func() T.Tuple2[A, W] {
		t := fa()
		return T.MakeTuple2(t.F1.F1, t.F1.F2(t.F2))
	}
}

func MonadMap[GB ~func() T.Tuple2[B, W], GA ~func() T.Tuple2[A, W], FCT ~func(A) B, W, A, B any](fa GA, f FCT) GB {
	return IO.MonadMap[GA, GB](fa, T.Map2(f, F.Identity[W]))
}

func Map[GB ~func() T.Tuple2[B, W], GA ~func() T.Tuple2[A, W], FCT ~func(A) B, W, A, B any](f FCT) func(GA) GB {
	return IO.Map[GA, GB](T.Map2(f, F.Identity[W]))
}

func MonadChain[GB ~func() T.Tuple2[B, W], GA ~func() T.Tuple2[A, W], FCT ~func(A) GB, W, A, B any](s S.Semigroup[W]) func(GA, FCT) GB {
	return func(fa GA, f FCT) GB {

		return func() T.Tuple2[B, W] {
			a := fa()
			b := f(a.F1)()

			return T.MakeTuple2(b.F1, s.Concat(a.F2, b.F2))
		}
	}
}

func Chain[GB ~func() T.Tuple2[B, W], GA ~func() T.Tuple2[A, W], FCT ~func(A) GB, W, A, B any](s S.Semigroup[W]) func(FCT) func(GA) GB {
	return F.Curry2(F.Swap(MonadChain[GB, GA, FCT](s)))
}

func MonadAp[GB ~func() T.Tuple2[B, W], GAB ~func() T.Tuple2[func(A) B, W], GA ~func() T.Tuple2[A, W], W, A, B any](s S.Semigroup[W]) func(GAB, GA) GB {
	return func(fab GAB, fa GA) GB {
		return func() T.Tuple2[B, W] {
			f := fab()
			a := fa()

			return T.MakeTuple2(f.F1(a.F1), s.Concat(f.F2, a.F2))
		}
	}
}

func Ap[GB ~func() T.Tuple2[B, W], GAB ~func() T.Tuple2[func(A) B, W], GA ~func() T.Tuple2[A, W], W, A, B any](s S.Semigroup[W]) func(GA) func(GAB) GB {
	return F.Curry2(F.Swap(MonadAp[GB, GAB, GA](s)))
}

func MonadChainFirst[GB ~func() T.Tuple2[B, W], GA ~func() T.Tuple2[A, W], FCT ~func(A) GB, W, A, B any](s S.Semigroup[W]) func(GA, FCT) GA {
	chain := MonadChain[GA, GA, func(A) GA](s)
	return func(ma GA, f FCT) GA {
		return chain(ma, func(a A) GA {
			return MonadMap[GA](f(a), F.Constant1[B](a))
		})
	}
}

func ChainFirst[GB ~func() T.Tuple2[B, W], GA ~func() T.Tuple2[A, W], FCT ~func(A) GB, W, A, B any](s S.Semigroup[W]) func(FCT) func(GA) GA {
	return F.Curry2(F.Swap(MonadChainFirst[GB, GA, FCT](s)))
}

func Flatten[GAA ~func() T.Tuple2[GA, W], GA ~func() T.Tuple2[A, W], W, A any](s S.Semigroup[W]) func(GAA) GA {
	chain := MonadChain[GA, GAA, func(GA) GA](s)
	return func(mma GAA) GA {
		return chain(mma, F.Identity[GA])
	}
}

func Execute[GA ~func() T.Tuple2[A, W], W, A any](fa GA) W {
	return T.Second(fa())
}

func Evaluate[GA ~func() T.Tuple2[A, W], W, A any](fa GA) A {
	return T.First(fa())
}

// MonadCensor modifies the final accumulator value by applying a function
func MonadCensor[GA ~func() T.Tuple2[A, W], FCT ~func(W) W, W, A any](fa GA, f FCT) GA {
	return IO.MonadMap[GA, GA](fa, T.Map2(F.Identity[A], f))
}

// Censor modifies the final accumulator value by applying a function
func Censor[GA ~func() T.Tuple2[A, W], FCT ~func(W) W, W, A any](f FCT) func(GA) GA {
	return IO.Map[GA, GA](T.Map2(F.Identity[A], f))
}

// MonadListens projects a value from modifications made to the accumulator during an action
func MonadListens[GA ~func() T.Tuple2[A, W], GAB ~func() T.Tuple2[T.Tuple2[A, B], W], FCT ~func(W) B, W, A, B any](fa GA, f FCT) GAB {
	return func() T.Tuple2[T.Tuple2[A, B], W] {
		a := fa()
		return T.MakeTuple2(T.MakeTuple2(a.F1, f(a.F2)), a.F2)
	}
}

// Listens projects a value from modifications made to the accumulator during an action
func Listens[GA ~func() T.Tuple2[A, W], GAB ~func() T.Tuple2[T.Tuple2[A, B], W], FCT ~func(W) B, W, A, B any](f FCT) func(GA) GAB {
	return F.Bind2nd(MonadListens[GA, GAB, FCT], f)
}
