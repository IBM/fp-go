// Copyright (c) 2023 - 2025 IBM Corp.
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
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/reader"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

// withContextK makes the computations returned by f fail with the cause of a cancelled
// context instead of running.
func withContextK[S, A, B any](f Kleisli[S, A, B]) Kleisli[S, A, B] {
	return F.Flow2(f, reader.Map[S](RIORES.WithContext[Pair[S, B]]))
}

// chainLeft is ChainLeft without the context check, for handlers that check the context
// themselves.
func chainLeft[S, A any](f Kleisli[S, error, A]) Operator[S, A, A] {
	return readert.Chain[StateReaderIOResult[S, A], StateReaderIOResult[S, A]](RIOR.ChainLeft[context.Context, Pair[S, A]], f)
}

// chainFirstLeft is ChainFirstLeft without the context check, for handlers that check the
// context themselves.
func chainFirstLeft[A, S, B any](f Kleisli[S, error, B]) Operator[S, A, A] {
	return readert.ChainFirst[StateReaderIOResult[S, A], StateReaderIOResult[S, B]](RIOR.ChainFirstLeft[Pair[S, A], context.Context, Pair[S, B]], f)
}

// MonadChainLeft recovers from the error of a failed computation.
//
// If ma fails, f is called with the error and its computation runs in place of the failure.
// A failure carries no state, so the recovery computation starts from the state that ma
// started from. The handler does not run if the context is already cancelled. A successful
// computation passes through without calling f.
//
// Type Parameters:
//   - S: The type of the state
//   - A: The type of the success value
//
// Parameters:
//   - ma: The computation to recover
//   - f: The recovery computation for the error
//
// Returns:
//   - StateReaderIOResult[S, A]: The result of ma, or of the recovery computation
//
// See Also:
//   - ChainLeft: The curried version
//   - MonadMapLeft: Transforms the error instead of recovering from it
//
//go:inline
func MonadChainLeft[S, A any](ma StateReaderIOResult[S, A], f Kleisli[S, error, A]) StateReaderIOResult[S, A] {
	return readert.MonadChain[StateReaderIOResult[S, A], StateReaderIOResult[S, A]](RIOR.MonadChainLeft[context.Context, Pair[S, A]], ma, withContextK(f))
}

// ChainLeft recovers from the error of a failed computation.
//
// If the computation fails, f is called with the error and its computation runs in place of
// the failure. A failure carries no state, so the recovery computation starts from the
// state that the failed computation started from. The handler does not run if the context
// is already cancelled. A successful computation passes through without calling f.
//
// Type Parameters:
//   - S: The type of the state
//   - A: The type of the success value
//
// Parameters:
//   - f: The recovery computation for the error
//
// Returns:
//   - Operator[S, A, A]: An operator that recovers from errors
//
// See Also:
//   - MonadChainLeft: The uncurried version
//   - OrElse: An alias
//   - CatchAs: Recovers only from errors of a given type
//
//go:inline
func ChainLeft[S, A any](f Kleisli[S, error, A]) Operator[S, A, A] {
	return chainLeft(withContextK(f))
}

// OrElse is an alias for ChainLeft.
//
// Type Parameters:
//   - S: The type of the state
//   - A: The type of the success value
//
// Parameters:
//   - onLeft: The recovery computation for the error
//
// Returns:
//   - Operator[S, A, A]: An operator that recovers from errors
//
// See Also:
//   - ChainLeft: The function this is an alias for
//
//go:inline
func OrElse[S, A any](onLeft Kleisli[S, error, A]) Operator[S, A, A] {
	return ChainLeft(onLeft)
}

// MonadMapLeft transforms the error of a failed computation.
//
// A successful computation passes through without calling f.
//
// Type Parameters:
//   - S: The type of the state
//   - A: The type of the success value
//
// Parameters:
//   - ma: The computation whose error to transform
//   - f: The function that maps the error to a new error
//
// Returns:
//   - StateReaderIOResult[S, A]: A computation that fails with the transformed error
//
// See Also:
//   - MapLeft: The curried version
//
//go:inline
func MonadMapLeft[S, A any](ma StateReaderIOResult[S, A], f Endomorphism[error]) StateReaderIOResult[S, A] {
	return readert.MonadMap[StateReaderIOResult[S, A], StateReaderIOResult[S, A]](RIORES.MonadMapLeft[Pair[S, A]], ma, f)
}

// MapLeft transforms the error of a failed computation.
//
// A successful computation passes through without calling f.
//
// Type Parameters:
//   - S: The type of the state
//   - A: The type of the success value
//
// Parameters:
//   - f: The function that maps the error to a new error
//
// Returns:
//   - Operator[S, A, A]: An operator that transforms errors
//
// See Also:
//   - MonadMapLeft: The uncurried version
//   - MapErrorAs: Transforms only errors of a given type
//
//go:inline
func MapLeft[S, A any](f Endomorphism[error]) Operator[S, A, A] {
	return readert.Map[StateReaderIOResult[S, A], StateReaderIOResult[S, A], S](RIORES.MapLeft[Pair[S, A]], f)
}

// MonadChainFirstLeft runs f on the error of a failed computation for its side effect and
// keeps the original error.
//
// f starts from the state that ma started from, and does not run if the context is already
// cancelled. Its result, success or failure, and its state are discarded. A successful
// computation passes through without calling f.
//
// Type Parameters:
//   - S: The type of the state
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - ma: The computation to observe
//   - f: The side effect for the error
//
// Returns:
//   - StateReaderIOResult[S, A]: A computation with the result of ma
//
// See Also:
//   - ChainFirstLeft: The curried version
//   - MonadChainLeft: Recovers from the error instead of observing it
//
//go:inline
func MonadChainFirstLeft[S, A, B any](ma StateReaderIOResult[S, A], f Kleisli[S, error, B]) StateReaderIOResult[S, A] {
	return readert.MonadChainFirst[StateReaderIOResult[S, A], StateReaderIOResult[S, B]](RIOR.MonadChainFirstLeft[Pair[S, A], context.Context, Pair[S, B]], ma, withContextK(f))
}

// ChainFirstLeft runs f on the error of a failed computation for its side effect and keeps
// the original error.
//
// f starts from the state that the failed computation started from, and does not run if
// the context is already cancelled. Its result, success or failure, and its state are
// discarded. A successful computation passes through without calling f.
//
// Type Parameters:
//   - A: The type of the success value
//   - S: The type of the state
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - f: The side effect for the error
//
// Returns:
//   - Operator[S, A, A]: An operator that preserves the original value or error
//
// See Also:
//   - MonadChainFirstLeft: The uncurried version
//   - TapLeft: An alias
//   - ChainLeft: Recovers from the error instead of observing it
//
//go:inline
func ChainFirstLeft[A, S, B any](f Kleisli[S, error, B]) Operator[S, A, A] {
	return chainFirstLeft[A](withContextK(f))
}

// MonadTapLeft is an alias for MonadChainFirstLeft.
//
// Type Parameters:
//   - S: The type of the state
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - ma: The computation to observe
//   - f: The side effect for the error
//
// Returns:
//   - StateReaderIOResult[S, A]: A computation with the result of ma
//
// See Also:
//   - MonadChainFirstLeft: The function this is an alias for
//
//go:inline
func MonadTapLeft[S, A, B any](ma StateReaderIOResult[S, A], f Kleisli[S, error, B]) StateReaderIOResult[S, A] {
	return MonadChainFirstLeft(ma, f)
}

// TapLeft is an alias for ChainFirstLeft.
//
// Type Parameters:
//   - A: The type of the success value
//   - S: The type of the state
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - f: The side effect for the error
//
// Returns:
//   - Operator[S, A, A]: An operator that preserves the original value or error
//
// See Also:
//   - ChainFirstLeft: The function this is an alias for
//   - TapErrorAs: Observes only errors of a given type
//
//go:inline
func TapLeft[A, S, B any](f Kleisli[S, error, B]) Operator[S, A, A] {
	return ChainFirstLeft[A](f)
}
