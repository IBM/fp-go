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

package state

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/pair"
)

var (
	undefined any = struct{}{}
)

//go:inline
func Get[S any]() State[S, S] {
	return pair.Of[S]
}

func Gets[FCT ~func(S) A, A, S any](f FCT) State[S, A] {
	return func(s S) Pair[S, A] {
		return pair.MakePair(s, f(s))
	}
}

//go:inline
func Put[S any]() State[S, any] {
	return function.Bind2nd(pair.MakePair[S, any], undefined)
}

func Modify[FCT ~func(S) S, S any](f FCT) State[S, any] {
	return function.Flow2(
		f,
		function.Bind2nd(pair.MakePair[S, any], undefined),
	)
}

//go:inline
func Of[S, A any](a A) State[S, A] {
	return function.Bind2nd(pair.MakePair[S, A], a)
}

func MonadMap[S any, FCT ~func(A) B, A, B any](fa State[S, A], f FCT) State[S, B] {
	return func(s S) Pair[S, B] {
		p2 := fa(s)
		return pair.MakePair(pair.Head(p2), f(pair.Tail(p2)))
	}
}

//go:inline
func Map[S any, FCT ~func(A) B, A, B any](f FCT) Operator[S, A, B] {
	return function.Bind2nd(MonadMap[S, FCT, A, B], f)
}

func MonadChain[S any, FCT ~func(A) State[S, B], A, B any](fa State[S, A], f FCT) State[S, B] {
	return func(s S) Pair[S, B] {
		a := fa(s)
		return f(pair.Tail(a))(pair.Head(a))
	}
}

//go:inline
func Chain[S any, FCT ~func(A) State[S, B], A, B any](f FCT) Operator[S, A, B] {
	return function.Bind2nd(MonadChain[S, FCT, A, B], f)
}

func MonadAp[B, S, A any](fab State[S, func(A) B], fa State[S, A]) State[S, B] {
	return func(s S) Pair[S, B] {
		f := fab(s)
		a := fa(pair.Head(f))

		return pair.MakePair(pair.Head(a), pair.Tail(f)(pair.Tail(a)))
	}
}

//go:inline
func Ap[B, S, A any](ga State[S, A]) Operator[S, func(A) B, B] {
	return function.Bind2nd(MonadAp[B, S, A], ga)
}

func MonadChainFirst[S any, FCT ~func(A) State[S, B], A, B any](ma State[S, A], f FCT) State[S, A] {
	return chain.MonadChainFirst(
		MonadChain[S, func(A) State[S, A], A, A],
		MonadMap[S, func(B) A],
		ma,
		f,
	)
}

func ChainFirst[S any, FCT ~func(A) State[S, B], A, B any](f FCT) Operator[S, A, A] {
	return chain.ChainFirst(
		Chain[S, func(A) State[S, A], A, A],
		Map[S, func(B) A],
		f,
	)
}

//go:inline
func Flatten[S, A any](mma State[S, State[S, A]]) State[S, A] {
	return MonadChain(mma, function.Identity[State[S, A]])
}

func Execute[A, S any](s S) func(State[S, A]) S {
	return func(fa State[S, A]) S {
		return pair.Head(fa(s))
	}
}

func Evaluate[A, S any](s S) func(State[S, A]) A {
	return func(fa State[S, A]) A {
		return pair.Tail(fa(s))
	}
}

func MonadFlap[FAB ~func(A) B, S, A, B any](fab State[S, FAB], a A) State[S, B] {
	return functor.MonadFlap(
		MonadMap[S, func(FAB) B],
		fab,
		a)
}

func Flap[S, A, B any](a A) Operator[S, func(A) B, B] {
	return functor.Flap(
		Map[S, func(func(A) B) B],
		a)
}
