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

package statereaderioresult

import (
	"context"

	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/statet"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/result"
)

// Left creates a StateReaderIOResult that represents a failed computation with the given error.
// The error value is immediately available and does not depend on state or context.
//
// Example:
//
//	result := statereaderioresult.Left[AppState, string](errors.New("validation failed"))
//	// Returns a failed computation that ignores state and context
func Left[S, A any](e error) StateReaderIOResult[S, A] {
	return function.Constant1[S](RIORES.Left[Pair[S, A]](e))
}

// Right creates a StateReaderIOResult that represents a successful computation with the given value.
// The value is wrapped and the state is passed through unchanged.
//
// Example:
//
//	result := statereaderioresult.Right[AppState](42)
//	// Returns a successful computation containing 42
func Right[S, A any](a A) StateReaderIOResult[S, A] {
	return statet.Of[StateReaderIOResult[S, A]](RIORES.Of[Pair[S, A]], a)
}

// Of creates a StateReaderIOResult that represents a successful computation with the given value.
// This is the monadic return/pure operation for StateReaderIOResult.
// Equivalent to [Right].
//
// Example:
//
//	result := statereaderioresult.Of[AppState](42)
//	// Returns a successful computation containing 42
func Of[S, A any](a A) StateReaderIOResult[S, A] {
	return Right[S](a)
}

// MonadMap transforms the success value of a StateReaderIOResult using the provided function.
// If the computation fails, the error is propagated unchanged.
// The state is threaded through the computation.
// This is the functor map operation.
//
// Example:
//
//	result := statereaderioresult.MonadMap(
//	    statereaderioresult.Of[AppState](21),
//	    N.Mul(2),
//	) // Result contains 42
func MonadMap[S, A, B any](fa StateReaderIOResult[S, A], f func(A) B) StateReaderIOResult[S, B] {
	return statet.MonadMap[StateReaderIOResult[S, A], StateReaderIOResult[S, B]](
		RIORES.MonadMap[Pair[S, A], Pair[S, B]],
		fa,
		f,
	)
}

// Map is the curried version of [MonadMap].
// Returns a function that transforms a StateReaderIOResult.
//
// Example:
//
//	double := statereaderioresult.Map[AppState](N.Mul(2))
//	result := function.Pipe1(statereaderioresult.Of[AppState](21), double)
func Map[S, A, B any](f func(A) B) Operator[S, A, B] {
	return statet.Map[StateReaderIOResult[S, A], StateReaderIOResult[S, B]](
		RIORES.Map[Pair[S, A], Pair[S, B]],
		f,
	)
}

// MonadChain sequences two computations, passing the result of the first to a function
// that produces the second computation. This is the monadic bind operation.
// The state is threaded through both computations.
//
// Example:
//
//	result := statereaderioresult.MonadChain(
//	    statereaderioresult.Of[AppState](5),
//	    func(x int) statereaderioresult.StateReaderIOResult[AppState, string] {
//	        return statereaderioresult.Of[AppState](fmt.Sprintf("value: %d", x))
//	    },
//	)
func MonadChain[S, A, B any](fa StateReaderIOResult[S, A], f Kleisli[S, A, B]) StateReaderIOResult[S, B] {
	return statet.MonadChain(
		RIORES.MonadChain[Pair[S, A], Pair[S, B]],
		fa,
		f,
	)
}

// Chain is the curried version of [MonadChain].
// Returns a function that sequences computations.
//
// Example:
//
//	stringify := statereaderioresult.Chain[AppState](func(x int) statereaderioresult.StateReaderIOResult[AppState, string] {
//	    return statereaderioresult.Of[AppState](fmt.Sprintf("%d", x))
//	})
//	result := function.Pipe1(statereaderioresult.Of[AppState](42), stringify)
func Chain[S, A, B any](f Kleisli[S, A, B]) Operator[S, A, B] {
	return statet.Chain[StateReaderIOResult[S, A]](
		RIORES.Chain[Pair[S, A], Pair[S, B]],
		f,
	)
}

// MonadAp applies a function wrapped in a StateReaderIOResult to a value wrapped in a StateReaderIOResult.
// If either the function or the value fails, the error is propagated.
// The state is threaded through both computations sequentially.
// This is the applicative apply operation.
//
// Example:
//
//	fab := statereaderioresult.Of[AppState](N.Mul(2))
//	fa := statereaderioresult.Of[AppState](21)
//	result := statereaderioresult.MonadAp(fab, fa) // Result contains 42
func MonadAp[B, S, A any](fab StateReaderIOResult[S, func(A) B], fa StateReaderIOResult[S, A]) StateReaderIOResult[S, B] {
	return statet.MonadAp[StateReaderIOResult[S, A], StateReaderIOResult[S, B]](
		RIORES.MonadMap[Pair[S, A], Pair[S, B]],
		RIORES.MonadChain[Pair[S, func(A) B], Pair[S, B]],
		fab,
		fa,
	)
}

// Ap is the curried version of [MonadAp].
// Returns a function that applies a wrapped function to the given wrapped value.
func Ap[B, S, A any](fa StateReaderIOResult[S, A]) Operator[S, func(A) B, B] {
	return statet.Ap[StateReaderIOResult[S, A], StateReaderIOResult[S, B], StateReaderIOResult[S, func(A) B]](
		RIORES.Map[Pair[S, A], Pair[S, B]],
		RIORES.Chain[Pair[S, func(A) B], Pair[S, B]],
		fa,
	)
}

// FromReaderIOResult lifts a ReaderIOResult into a StateReaderIOResult.
// The state is passed through unchanged.
//
// Example:
//
//	riores := readerioresult.Of(42)
//	result := statereaderioresult.FromReaderIOResult[AppState](riores)
func FromReaderIOResult[S, A any](fa ReaderIOResult[A]) StateReaderIOResult[S, A] {
	return statet.FromF[StateReaderIOResult[S, A]](
		RIORES.MonadMap[A],
		fa,
	)
}

// FromIOResult lifts an IOResult into a StateReaderIOResult.
// The state is passed through unchanged and the context is ignored.
func FromIOResult[S, A any](fa IOResult[A]) StateReaderIOResult[S, A] {
	return FromReaderIOResult[S](RIORES.FromIOResult(fa))
}

// FromState lifts a State computation into a StateReaderIOResult.
// The computation cannot fail (uses the error type).
func FromState[S, A any](sa State[S, A]) StateReaderIOResult[S, A] {
	return statet.FromState[StateReaderIOResult[S, A]](RIORES.Of[Pair[S, A]], sa)
}

// FromIO lifts an IO computation into a StateReaderIOResult.
// The state is passed through unchanged and the context is ignored.
func FromIO[S, A any](fa IO[A]) StateReaderIOResult[S, A] {
	return FromReaderIOResult[S](RIORES.FromIO(fa))
}

// FromResult lifts a Result into a StateReaderIOResult.
// The state is passed through unchanged and the context is ignored.
//
// Example:
//
//	result := statereaderioresult.FromResult[AppState](result.Of(42))
func FromResult[S, A any](ma Result[A]) StateReaderIOResult[S, A] {
	return result.Fold(Left[S, A], Right[S, A])(ma)
}

// Combinators

// Local runs a computation with a modified context.
// The function f transforms the context before passing it to the computation.
//
// Example:
//
//	// Modify context before running computation
//	withTimeout := statereaderioresult.Local[AppState](
//	    func(ctx context.Context) context.Context {
//	        ctx, _ = context.WithTimeout(ctx, 60*time.Second)
//	        return ctx
//	    }
//	)
//	result := withTimeout(computation)
func Local[S, A any](f func(context.Context) context.Context) func(StateReaderIOResult[S, A]) StateReaderIOResult[S, A] {
	return func(ma StateReaderIOResult[S, A]) StateReaderIOResult[S, A] {
		return function.Flow2(ma, RIOR.Local[Pair[S, A]](f))
	}
}

// Asks creates a computation that derives a value from the context.
// The function receives the context and returns a StateReaderIOResult.
//
// Example:
//
//	getValue := statereaderioresult.Asks[AppState, string](
//	    func(ctx context.Context) statereaderioresult.StateReaderIOResult[AppState, string] {
//	        return statereaderioresult.Of[AppState](ctx.Value("key").(string))
//	    },
//	)
func Asks[S, A any](f func(context.Context) StateReaderIOResult[S, A]) StateReaderIOResult[S, A] {
	return func(s S) ReaderIOResult[Pair[S, A]] {
		return func(ctx context.Context) IOResult[Pair[S, A]] {
			return f(ctx)(s)(ctx)
		}
	}
}

// FromResultK lifts a Result-returning function into a Kleisli arrow for StateReaderIOResult.
//
// Example:
//
//	validate := func(x int) result.Result[int] {
//	    if x > 0 { return result.Of(x) }
//	    return result.Error[int](errors.New("negative"))
//	}
//	kleisli := statereaderioresult.FromResultK[AppState](validate)
func FromResultK[S, A, B any](f func(A) Result[B]) Kleisli[S, A, B] {
	return function.Flow2(
		f,
		FromResult[S, B],
	)
}

// FromIOK lifts an IO-returning function into a Kleisli arrow for StateReaderIOResult.
func FromIOK[S, A, B any](f func(A) IO[B]) Kleisli[S, A, B] {
	return function.Flow2(
		f,
		FromIO[S, B],
	)
}

// FromIOResultK lifts an IOResult-returning function into a Kleisli arrow for StateReaderIOResult.
func FromIOResultK[S, A, B any](f func(A) IOResult[B]) Kleisli[S, A, B] {
	return function.Flow2(
		f,
		FromIOResult[S, B],
	)
}

// FromReaderIOResultK lifts a ReaderIOResult-returning function into a Kleisli arrow for StateReaderIOResult.
func FromReaderIOResultK[S, A, B any](f func(A) ReaderIOResult[B]) Kleisli[S, A, B] {
	return function.Flow2(
		f,
		FromReaderIOResult[S, B],
	)
}

// MonadChainReaderIOResultK chains a StateReaderIOResult with a ReaderIOResult-returning function.
func MonadChainReaderIOResultK[S, A, B any](ma StateReaderIOResult[S, A], f func(A) ReaderIOResult[B]) StateReaderIOResult[S, B] {
	return MonadChain(ma, FromReaderIOResultK[S](f))
}

// ChainReaderIOResultK is the curried version of [MonadChainReaderIOResultK].
func ChainReaderIOResultK[S, A, B any](f func(A) ReaderIOResult[B]) Operator[S, A, B] {
	return Chain(FromReaderIOResultK[S](f))
}

// MonadChainIOResultK chains a StateReaderIOResult with an IOResult-returning function.
func MonadChainIOResultK[S, A, B any](ma StateReaderIOResult[S, A], f func(A) IOResult[B]) StateReaderIOResult[S, B] {
	return MonadChain(ma, FromIOResultK[S](f))
}

// ChainIOResultK is the curried version of [MonadChainIOResultK].
func ChainIOResultK[S, A, B any](f func(A) IOResult[B]) Operator[S, A, B] {
	return Chain(FromIOResultK[S](f))
}

// MonadChainResultK chains a StateReaderIOResult with a Result-returning function.
func MonadChainResultK[S, A, B any](ma StateReaderIOResult[S, A], f func(A) Result[B]) StateReaderIOResult[S, B] {
	return MonadChain(ma, FromResultK[S](f))
}

// ChainResultK is the curried version of [MonadChainResultK].
func ChainResultK[S, A, B any](f func(A) Result[B]) Operator[S, A, B] {
	return Chain(FromResultK[S](f))
}
