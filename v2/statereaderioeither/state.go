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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/statet"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/readerioeither"
)

// Left creates a StateReaderIOEither that represents a failed computation with the given error.
// The error value is immediately available and does not depend on state or context.
//
// Example:
//
//	result := statereaderioeither.Left[AppState, Config, string](errors.New("validation failed"))
//	// Returns a failed computation that ignores state and context
func Left[S, R, A, E any](e E) StateReaderIOEither[S, R, E, A] {
	return function.Constant1[S](readerioeither.Left[R, Pair[S, A]](e))
}

// Right creates a StateReaderIOEither that represents a successful computation with the given value.
// The value is wrapped and the state is passed through unchanged.
//
// Example:
//
//	result := statereaderioeither.Right[AppState, Config, error](42)
//	// Returns a successful computation containing 42
func Right[S, R, E, A any](a A) StateReaderIOEither[S, R, E, A] {
	return statet.Of[StateReaderIOEither[S, R, E, A]](readerioeither.Of[R, E, Pair[S, A]], a)
}

// Of creates a StateReaderIOEither that represents a successful computation with the given value.
// This is the monadic return/pure operation for StateReaderIOEither.
// Equivalent to [Right].
//
// Example:
//
//	result := statereaderioeither.Of[AppState, Config, error](42)
//	// Returns a successful computation containing 42
func Of[S, R, E, A any](a A) StateReaderIOEither[S, R, E, A] {
	return Right[S, R, E](a)
}

// MonadMap transforms the success value of a StateReaderIOEither using the provided function.
// If the computation fails, the error is propagated unchanged.
// The state is threaded through the computation.
// This is the functor map operation.
//
// Example:
//
//	result := statereaderioeither.MonadMap(
//	    statereaderioeither.Of[AppState, Config, error](21),
//	    N.Mul(2),
//	) // Result contains 42
func MonadMap[S, R, E, A, B any](fa StateReaderIOEither[S, R, E, A], f func(A) B) StateReaderIOEither[S, R, E, B] {
	return statet.MonadMap[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]](
		readerioeither.MonadMap[R, E, Pair[S, A], Pair[S, B]],
		fa,
		f,
	)
}

// Map is the curried version of [MonadMap].
// Returns a function that transforms a StateReaderIOEither.
//
// Example:
//
//	double := statereaderioeither.Map[AppState, Config, error](N.Mul(2))
//	result := function.Pipe1(statereaderioeither.Of[AppState, Config, error](21), double)
func Map[S, R, E, A, B any](f func(A) B) Operator[S, R, E, A, B] {
	return statet.Map[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]](
		readerioeither.Map[R, E, Pair[S, A], Pair[S, B]],
		f,
	)
}

// MonadChain sequences two computations, passing the result of the first to a function
// that produces the second computation. This is the monadic bind operation.
// The state is threaded through both computations.
//
// Example:
//
//	result := statereaderioeither.MonadChain(
//	    statereaderioeither.Of[AppState, Config, error](5),
//	    func(x int) statereaderioeither.StateReaderIOEither[AppState, Config, error, string] {
//	        return statereaderioeither.Of[AppState, Config, error](fmt.Sprintf("value: %d", x))
//	    },
//	)
func MonadChain[S, R, E, A, B any](fa StateReaderIOEither[S, R, E, A], f Kleisli[S, R, E, A, B]) StateReaderIOEither[S, R, E, B] {
	return statet.MonadChain(
		readerioeither.MonadChain[R, E, Pair[S, A], Pair[S, B]],
		fa,
		f,
	)
}

// Chain is the curried version of [MonadChain].
// Returns a function that sequences computations.
//
// Example:
//
//	stringify := statereaderioeither.Chain(func(x int) statereaderioeither.StateReaderIOEither[AppState, Config, error, string] {
//	    return statereaderioeither.Of[AppState, Config, error](fmt.Sprintf("%d", x))
//	})
//	result := function.Pipe1(statereaderioeither.Of[AppState, Config, error](42), stringify)
func Chain[S, R, E, A, B any](f Kleisli[S, R, E, A, B]) Operator[S, R, E, A, B] {
	return statet.Chain[StateReaderIOEither[S, R, E, A]](
		readerioeither.Chain[R, E, Pair[S, A], Pair[S, B]],
		f,
	)
}

// MonadAp applies a function wrapped in a StateReaderIOEither to a value wrapped in a StateReaderIOEither.
// If either the function or the value fails, the error is propagated.
// The state is threaded through both computations sequentially.
// This is the applicative apply operation.
//
// Example:
//
//	fab := statereaderioeither.Of[AppState, Config, error](N.Mul(2))
//	fa := statereaderioeither.Of[AppState, Config, error](21)
//	result := statereaderioeither.MonadAp(fab, fa) // Result contains 42
func MonadAp[B, S, R, E, A any](fab StateReaderIOEither[S, R, E, func(A) B], fa StateReaderIOEither[S, R, E, A]) StateReaderIOEither[S, R, E, B] {
	return statet.MonadAp[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]](
		readerioeither.MonadMap[R, E, Pair[S, A], Pair[S, B]],
		readerioeither.MonadChain[R, E, Pair[S, func(A) B], Pair[S, B]],
		fab,
		fa,
	)
}

// Ap is the curried version of [MonadAp].
// Returns a function that applies a wrapped function to the given wrapped value.
func Ap[B, S, R, E, A any](fa StateReaderIOEither[S, R, E, A]) Operator[S, R, E, func(A) B, B] {
	return statet.Ap[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B], StateReaderIOEither[S, R, E, func(A) B]](
		readerioeither.Map[R, E, Pair[S, A], Pair[S, B]],
		readerioeither.Chain[R, E, Pair[S, func(A) B], Pair[S, B]],
		fa,
	)
}

// FromReaderIOEither lifts a ReaderIOEither into a StateReaderIOEither.
// The state is passed through unchanged.
//
// Example:
//
//	rioe := readerioeither.Of[Config, error](42)
//	result := statereaderioeither.FromReaderIOEither[AppState](rioe)
func FromReaderIOEither[S, R, E, A any](fa ReaderIOEither[R, E, A]) StateReaderIOEither[S, R, E, A] {
	return statet.FromF[StateReaderIOEither[S, R, E, A]](
		readerioeither.MonadMap[R, E, A],
		fa,
	)
}

// FromReaderEither lifts a ReaderEither into a StateReaderIOEither.
// The state is passed through unchanged.
func FromReaderEither[S, R, E, A any](fa ReaderEither[R, E, A]) StateReaderIOEither[S, R, E, A] {
	return FromReaderIOEither[S](readerioeither.FromReaderEither(fa))
}

// FromIOEither lifts an IOEither into a StateReaderIOEither.
// The state is passed through unchanged and the context is ignored.
func FromIOEither[S, R, E, A any](fa IOEither[E, A]) StateReaderIOEither[S, R, E, A] {
	return FromReaderIOEither[S](readerioeither.FromIOEither[R](fa))
}

// FromState lifts a State computation into a StateReaderIOEither.
// The computation cannot fail (uses the error type parameter).
func FromState[R, E, S, A any](sa State[S, A]) StateReaderIOEither[S, R, E, A] {
	return statet.FromState[StateReaderIOEither[S, R, E, A]](readerioeither.Of[R, E, Pair[S, A]], sa)
}

// FromIO lifts an IO computation into a StateReaderIOEither.
// The state is passed through unchanged and the context is ignored.
func FromIO[S, R, E, A any](fa IO[A]) StateReaderIOEither[S, R, E, A] {
	return FromReaderIOEither[S](readerioeither.FromIO[R, E](fa))
}

// FromReader lifts a Reader into a StateReaderIOEither.
// The state is passed through unchanged.
func FromReader[S, E, R, A any](fa Reader[R, A]) StateReaderIOEither[S, R, E, A] {
	return FromReaderIOEither[S](readerioeither.FromReader[E](fa))
}

// FromEither lifts an Either into a StateReaderIOEither.
// The state is passed through unchanged and the context is ignored.
//
// Example:
//
//	result := statereaderioeither.FromEither[AppState, Config](either.Right[error](42))
func FromEither[S, R, E, A any](ma Either[E, A]) StateReaderIOEither[S, R, E, A] {
	return either.MonadFold(ma, Left[S, R, A, E], Right[S, R, E, A])
}

// Combinators

// Local runs a computation with a modified context.
// The function f transforms the context before passing it to the computation.
//
// Example:
//
//	// Modify config before running computation
//	withTimeout := statereaderioeither.Local[AppState, error, int](
//	    func(cfg Config) Config { return Config{...cfg, Timeout: 60} }
//	)
//	result := withTimeout(computation)
func Local[S, E, A, B, R1, R2 any](f func(R2) R1) func(StateReaderIOEither[S, R1, E, A]) StateReaderIOEither[S, R2, E, A] {
	return func(ma StateReaderIOEither[S, R1, E, A]) StateReaderIOEither[S, R2, E, A] {
		return function.Flow2(ma, readerioeither.Local[E, Pair[S, A]](f))
	}
}

// Asks creates a computation that derives a value from the context.
// The function receives the context and returns a StateReaderIOEither.
//
// Example:
//
//	getTimeout := statereaderioeither.Asks[AppState, Config, error, int](
//	    func(cfg Config) statereaderioeither.StateReaderIOEither[AppState, Config, error, int] {
//	        return statereaderioeither.Of[AppState, Config, error](cfg.Timeout)
//	    },
//	)
func Asks[
	S, R, E, A any,
](f func(R) StateReaderIOEither[S, R, E, A]) StateReaderIOEither[S, R, E, A] {
	return func(s S) ReaderIOEither[R, E, Pair[S, A]] {
		return func(r R) IOEither[E, Pair[S, A]] {
			return f(r)(s)(r)
		}
	}
}

// FromEitherK lifts an Either-returning function into a Kleisli arrow for StateReaderIOEither.
//
// Example:
//
//	validate := func(x int) either.Either[error, int] {
//	    if x > 0 { return either.Right[error](x) }
//	    return either.Left[int](errors.New("negative"))
//	}
//	kleisli := statereaderioeither.FromEitherK[AppState, Config](validate)
func FromEitherK[S, R, E, A, B any](f either.Kleisli[E, A, B]) Kleisli[S, R, E, A, B] {
	return function.Flow2(
		f,
		FromEither[S, R, E, B],
	)
}

// FromIOK lifts an IO-returning function into a Kleisli arrow for StateReaderIOEither.
func FromIOK[S, R, E, A, B any](f func(A) IO[B]) Kleisli[S, R, E, A, B] {
	return function.Flow2(
		f,
		FromIO[S, R, E, B],
	)
}

// FromIOEitherK lifts an IOEither-returning function into a Kleisli arrow for StateReaderIOEither.
func FromIOEitherK[
	S, R, E, A, B any,
](f ioeither.Kleisli[E, A, B]) Kleisli[S, R, E, A, B] {
	return function.Flow2(
		f,
		FromIOEither[S, R, E, B],
	)
}

// FromReaderIOEitherK lifts a ReaderIOEither-returning function into a Kleisli arrow for StateReaderIOEither.
func FromReaderIOEitherK[S, R, E, A, B any](f readerioeither.Kleisli[R, E, A, B]) Kleisli[S, R, E, A, B] {
	return function.Flow2(
		f,
		FromReaderIOEither[S, R, E, B],
	)
}

// MonadChainReaderIOEitherK chains a StateReaderIOEither with a ReaderIOEither-returning function.
func MonadChainReaderIOEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f readerioeither.Kleisli[R, E, A, B]) StateReaderIOEither[S, R, E, B] {
	return MonadChain(ma, FromReaderIOEitherK[S](f))
}

// ChainReaderIOEitherK is the curried version of [MonadChainReaderIOEitherK].
func ChainReaderIOEitherK[S, R, E, A, B any](f readerioeither.Kleisli[R, E, A, B]) Operator[S, R, E, A, B] {
	return Chain(FromReaderIOEitherK[S](f))
}

// MonadChainIOEitherK chains a StateReaderIOEither with an IOEither-returning function.
func MonadChainIOEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f ioeither.Kleisli[E, A, B]) StateReaderIOEither[S, R, E, B] {
	return MonadChain(ma, FromIOEitherK[S, R](f))
}

// ChainIOEitherK is the curried version of [MonadChainIOEitherK].
func ChainIOEitherK[S, R, E, A, B any](f ioeither.Kleisli[E, A, B]) Operator[S, R, E, A, B] {
	return Chain(FromIOEitherK[S, R](f))
}

// MonadChainEitherK chains a StateReaderIOEither with an Either-returning function.
func MonadChainEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f either.Kleisli[E, A, B]) StateReaderIOEither[S, R, E, B] {
	return MonadChain(ma, FromEitherK[S, R](f))
}

// ChainEitherK is the curried version of [MonadChainEitherK].
func ChainEitherK[S, R, E, A, B any](f either.Kleisli[E, A, B]) Operator[S, R, E, A, B] {
	return Chain(FromEitherK[S, R](f))
}
