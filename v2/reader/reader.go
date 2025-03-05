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

package reader

import (
	"github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/functor"
	T "github.com/IBM/fp-go/v2/tuple"
)

// Ask reads the current context
func Ask[R any]() Reader[R, R] {
	return function.Identity[R]
}

// Asks projects a value from the global context in a Reader
func Asks[R, A any](f Reader[R, A]) Reader[R, A] {
	return f
}

func AsksReader[R, A any](f func(R) Reader[R, A]) Reader[R, A] {
	return func(r R) A {
		return f(r)(r)
	}
}

func MonadMap[E, A, B any](fa Reader[E, A], f func(A) B) Reader[E, B] {
	return function.Flow2(fa, f)
}

// Map can be used to turn functions `func(A)B` into functions `(fa F[A])F[B]` whose argument and return types
// use the type constructor `F` to represent some computational context.
func Map[E, A, B any](f func(A) B) Operator[E, A, B] {
	return function.Bind2nd(MonadMap[E, A, B], f)
}

func MonadAp[B, R, A any](fab Reader[R, func(A) B], fa Reader[R, A]) Reader[R, B] {
	return func(r R) B {
		return fab(r)(fa(r))
	}
}

// Ap applies a function to an argument under a type constructor.
func Ap[B, R, A any](fa Reader[R, A]) Operator[R, func(A) B, B] {
	return function.Bind2nd(MonadAp[B, R, A], fa)
}

func Of[R, A any](a A) Reader[R, A] {
	return function.Constant1[R](a)
}

func MonadChain[R, A, B any](ma Reader[R, A], f func(A) Reader[R, B]) Reader[R, B] {
	return func(r R) B {
		return f(ma(r))(r)
	}
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[R, A, B any](f func(A) Reader[R, B]) Operator[R, A, B] {
	return function.Bind2nd(MonadChain[R, A, B], f)
}

func Flatten[R, A any](mma func(R) Reader[R, A]) Reader[R, A] {
	return MonadChain(mma, function.Identity[Reader[R, A]])
}

func Compose[R, B, C any](ab Reader[R, B]) func(Reader[B, C]) Reader[R, C] {
	return func(bc Reader[B, C]) Reader[R, C] {
		return function.Flow2(ab, bc)
	}
}

func First[A, B, C any](pab Reader[A, B]) Reader[T.Tuple2[A, C], T.Tuple2[B, C]] {
	return func(tac T.Tuple2[A, C]) T.Tuple2[B, C] {
		return T.MakeTuple2(pab(tac.F1), tac.F2)
	}
}

func Second[A, B, C any](pbc Reader[B, C]) Reader[T.Tuple2[A, B], T.Tuple2[A, C]] {
	return func(tab T.Tuple2[A, B]) T.Tuple2[A, C] {
		return T.MakeTuple2(tab.F1, pbc(tab.F2))
	}
}

func Promap[E, A, D, B any](f func(D) E, g func(A) B) func(Reader[E, A]) Reader[D, B] {
	return func(fea Reader[E, A]) Reader[D, B] {
		return function.Flow3(f, fea, g)
	}
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[R2, R1, A any](f func(R2) R1) func(Reader[R1, A]) Reader[R2, A] {
	return Compose[R2, R1, A](f)
}

// Read applies a context to a reader to obtain its value
func Read[E, A any](e E) func(Reader[E, A]) A {
	return I.Ap[A](e)
}

func MonadFlap[R, A, B any](fab Reader[R, func(A) B], a A) Reader[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

func Flap[R, A, B any](a A) func(Reader[R, func(A) B]) Reader[R, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}
