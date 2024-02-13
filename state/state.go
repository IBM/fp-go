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

package state

import (
	P "github.com/IBM/fp-go/pair"
	R "github.com/IBM/fp-go/reader"
	G "github.com/IBM/fp-go/state/generic"
)

type State[S, A any] R.Reader[S, P.Pair[A, S]]

func Get[S any]() State[S, S] {
	return G.Get[State[S, S]]()
}

func Gets[FCT ~func(S) A, A, S any](f FCT) State[S, A] {
	return G.Gets[State[S, A]](f)
}

func Put[S any]() State[S, any] {
	return G.Put[State[S, any]]()
}

func Modify[FCT ~func(S) S, S any](f FCT) State[S, any] {
	return G.Modify[State[S, any]](f)
}

func Of[S, A any](a A) State[S, A] {
	return G.Of[State[S, A]](a)
}

func MonadMap[S any, FCT ~func(A) B, A, B any](fa State[S, A], f FCT) State[S, B] {
	return G.MonadMap[State[S, B], State[S, A]](fa, f)
}

func Map[S any, FCT ~func(A) B, A, B any](f FCT) func(State[S, A]) State[S, B] {
	return G.Map[State[S, B], State[S, A]](f)
}

func MonadChain[S any, FCT ~func(A) State[S, B], A, B any](fa State[S, A], f FCT) State[S, B] {
	return G.MonadChain[State[S, B], State[S, A]](fa, f)
}

func Chain[S any, FCT ~func(A) State[S, B], A, B any](f FCT) func(State[S, A]) State[S, B] {
	return G.Chain[State[S, B], State[S, A]](f)
}

func MonadAp[S, A, B any](fab State[S, func(A) B], fa State[S, A]) State[S, B] {
	return G.MonadAp[State[S, B], State[S, func(A) B], State[S, A]](fab, fa)
}

func Ap[S, A, B any](ga State[S, A]) func(State[S, func(A) B]) State[S, B] {
	return G.Ap[State[S, B], State[S, func(A) B], State[S, A]](ga)
}

func MonadChainFirst[S any, FCT ~func(A) State[S, B], A, B any](ma State[S, A], f FCT) State[S, A] {
	return G.MonadChainFirst[State[S, B], State[S, A]](ma, f)
}

func ChainFirst[S any, FCT ~func(A) State[S, B], A, B any](f FCT) func(State[S, A]) State[S, A] {
	return G.ChainFirst[State[S, B], State[S, A]](f)
}

func Flatten[S, A any](mma State[S, State[S, A]]) State[S, A] {
	return G.Flatten[State[S, State[S, A]], State[S, A]](mma)
}

func Execute[A, S any](s S) func(State[S, A]) S {
	return G.Execute[State[S, A]](s)
}

func Evaluate[A, S any](s S) func(State[S, A]) A {
	return G.Evaluate[State[S, A]](s)
}

func MonadFlap[FAB ~func(A) B, S, A, B any](fab State[S, FAB], a A) State[S, B] {
	return G.MonadFlap[FAB, State[S, FAB], State[S, B], S, A, B](fab, a)
}

func Flap[S, A, B any](a A) func(State[S, func(A) B]) State[S, B] {
	return G.Flap[func(A) B, State[S, func(A) B], State[S, B]](a)
}
