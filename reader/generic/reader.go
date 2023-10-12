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
	I "github.com/IBM/fp-go/identity/generic"
	FC "github.com/IBM/fp-go/internal/functor"
	T "github.com/IBM/fp-go/tuple"
)

// Reader[R, A] = func(R) A

// MakeReader creates a reader, i.e. a method that accepts a context and that returns a value
func MakeReader[GA ~func(R) A, R, A any](r GA) GA {
	return r
}

// Ask reads the current context
func Ask[GR ~func(R) R, R any]() GR {
	return MakeReader(F.Identity[R])
}

// Asks projects a value from the global context in a Reader
func Asks[GA ~func(R) A, R, A any](f GA) GA {
	return MakeReader(f)
}

func AsksReader[GA ~func(R) A, R, A any](f func(R) GA) GA {
	return MakeReader(func(r R) A {
		return f(r)(r)
	})
}

func MonadMap[GA ~func(E) A, GB ~func(E) B, E, A, B any](fa GA, f func(A) B) GB {
	return MakeReader(F.Flow2(fa, f))
}

// Map can be used to turn functions `func(A)B` into functions `(fa F[A])F[B]` whose argument and return types
// use the type constructor `F` to represent some computational context.
func Map[GA ~func(E) A, GB ~func(E) B, E, A, B any](f func(A) B) func(GA) GB {
	return F.Bind2nd(MonadMap[GA, GB, E, A, B], f)
}

func MonadAp[GA ~func(R) A, GB ~func(R) B, GAB ~func(R) func(A) B, R, A, B any](fab GAB, fa GA) GB {
	return MakeReader(func(r R) B {
		return fab(r)(fa(r))
	})
}

// Ap applies a function to an argument under a type constructor.
func Ap[GA ~func(R) A, GB ~func(R) B, GAB ~func(R) func(A) B, R, A, B any](fa GA) func(GAB) GB {
	return F.Bind2nd(MonadAp[GA, GB, GAB, R, A, B], fa)
}

func Of[GA ~func(R) A, R, A any](a A) GA {
	return F.Constant1[R](a)
}

func MonadChain[GA ~func(R) A, GB ~func(R) B, R, A, B any](ma GA, f func(A) GB) GB {
	return MakeReader(func(r R) B {
		return f(ma(r))(r)
	})
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[GA ~func(R) A, GB ~func(R) B, R, A, B any](f func(A) GB) func(GA) GB {
	return F.Bind2nd(MonadChain[GA, GB, R, A, B], f)
}

func Flatten[GA ~func(R) A, GGA ~func(R) GA, R, A any](mma GGA) GA {
	return MonadChain(mma, F.Identity[GA])
}

func Compose[AB ~func(A) B, BC ~func(B) C, AC ~func(A) C, A, B, C any](ab AB) func(BC) AC {
	return func(bc BC) AC {
		return F.Flow2(ab, bc)
	}
}

func First[GAB ~func(A) B, GABC ~func(T.Tuple2[A, C]) T.Tuple2[B, C], A, B, C any](pab GAB) GABC {
	return MakeReader(func(tac T.Tuple2[A, C]) T.Tuple2[B, C] {
		return T.MakeTuple2(pab(tac.F1), tac.F2)
	})
}

func Second[GBC ~func(B) C, GABC ~func(T.Tuple2[A, B]) T.Tuple2[A, C], A, B, C any](pbc GBC) GABC {
	return MakeReader(func(tab T.Tuple2[A, B]) T.Tuple2[A, C] {
		return T.MakeTuple2(tab.F1, pbc(tab.F2))
	})
}

func Promap[GA ~func(E) A, GB ~func(D) B, E, A, D, B any](f func(D) E, g func(A) B) func(GA) GB {
	return func(fea GA) GB {
		return MakeReader(F.Flow3(f, fea, g))
	}
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[GA1 ~func(R1) A, GA2 ~func(R2) A, R2, R1, A any](f func(R2) R1) func(GA1) GA2 {
	return func(r1 GA1) GA2 {
		return F.Flow2(f, r1)
	}
}

// Read applies a context to a reader to obtain its value
func Read[GA ~func(E) A, E, A any](e E) func(GA) A {
	return I.Ap[GA](e)
}

func MonadFlap[GAB ~func(R) func(A) B, GB ~func(R) B, R, A, B any](fab GAB, a A) GB {
	return FC.MonadFlap(MonadMap[GAB, GB], fab, a)
}

func Flap[GAB ~func(R) func(A) B, GB ~func(R) B, R, A, B any](a A) func(GAB) GB {
	return FC.Flap(MonadMap[GAB, GB], a)
}
