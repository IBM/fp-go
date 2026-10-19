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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	"github.com/IBM/fp-go/v2/internal/readert"
)

// MonadChainLeft recovers from the error of a failed computation.
//
// If ma fails, f is called with the error and its computation runs with the same
// environment in place of the failure. A successful computation passes through without
// calling f.
//
// Type Parameters:
//   - R: The environment required by the computations
//   - A: The type of the success value
//
// Parameters:
//   - ma: The computation to recover
//   - f: The recovery computation for the error
//
// Returns:
//   - ReaderIOResult[R, A]: The result of ma, or of the recovery computation
//
// See Also:
//   - ChainLeft: The curried version
//   - MonadMapLeft: Transforms the error instead of recovering from it
//
//go:inline
func MonadChainLeft[R, A any](ma ReaderIOResult[R, A], f Kleisli[R, error, A]) ReaderIOResult[R, A] {
	return readert.MonadChain[ReaderIOResult[R, A], ReaderIOResult[R, A]](ioresult.MonadChainLeft[A], ma, f)
}

// ChainLeft recovers from the error of a failed computation.
//
// If the computation fails, f is called with the error and its computation runs with the
// same environment in place of the failure. A successful computation passes through
// without calling f.
//
// Type Parameters:
//   - R: The environment required by the computations
//   - A: The type of the success value
//
// Parameters:
//   - f: The recovery computation for the error
//
// Returns:
//   - Operator[R, A, A]: An operator that recovers from errors
//
// See Also:
//   - MonadChainLeft: The uncurried version
//   - OrElse: An alias
//   - CatchAs: Recovers only from errors of a given type
//
//go:inline
func ChainLeft[R, A any](f Kleisli[R, error, A]) Operator[R, A, A] {
	return readert.Chain[ReaderIOResult[R, A], ReaderIOResult[R, A]](ioresult.ChainLeft[A], f)
}

// OrElse is an alias for ChainLeft.
//
// Type Parameters:
//   - R: The environment required by the computations
//   - A: The type of the success value
//
// Parameters:
//   - onLeft: The recovery computation for the error
//
// Returns:
//   - Operator[R, A, A]: An operator that recovers from errors
//
// See Also:
//   - ChainLeft: The function this is an alias for
//
//go:inline
func OrElse[R, A any](onLeft Kleisli[R, error, A]) Operator[R, A, A] {
	return ChainLeft(onLeft)
}

// MonadMapLeft transforms the error of a failed computation.
//
// A successful computation passes through without calling f.
//
// Type Parameters:
//   - R: The environment required by the computation
//   - A: The type of the success value
//
// Parameters:
//   - ma: The computation whose error to transform
//   - f: The function that maps the error to a new error
//
// Returns:
//   - ReaderIOResult[R, A]: A computation that fails with the transformed error
//
// See Also:
//   - MapLeft: The curried version
//
//go:inline
func MonadMapLeft[R, A any](ma ReaderIOResult[R, A], f Endomorphism[error]) ReaderIOResult[R, A] {
	return readert.MonadMap[ReaderIOResult[R, A], ReaderIOResult[R, A]](ioresult.MonadMapLeft[A], ma, f)
}

// MapLeft transforms the error of a failed computation.
//
// A successful computation passes through without calling f.
//
// Type Parameters:
//   - R: The environment required by the computation
//   - A: The type of the success value
//
// Parameters:
//   - f: The function that maps the error to a new error
//
// Returns:
//   - Operator[R, A, A]: An operator that transforms errors
//
// See Also:
//   - MonadMapLeft: The uncurried version
//   - MapErrorAs: Transforms only errors of a given type
//
//go:inline
func MapLeft[R, A any](f Endomorphism[error]) Operator[R, A, A] {
	return readert.Map[ReaderIOResult[R, A], ReaderIOResult[R, A], R](ioresult.MapLeft[A], f)
}

// MonadChainFirstLeft runs f on the error of a failed computation for its side effect and
// keeps the original error.
//
// f runs with the same environment. Its result, success or failure, is discarded. A
// successful computation passes through without calling f.
//
// Type Parameters:
//   - R: The environment required by the computations
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - ma: The computation to observe
//   - f: The side effect for the error
//
// Returns:
//   - ReaderIOResult[R, A]: A computation with the result of ma
//
// See Also:
//   - ChainFirstLeft: The curried version
//   - MonadChainLeft: Recovers from the error instead of observing it
//
//go:inline
func MonadChainFirstLeft[R, A, B any](ma ReaderIOResult[R, A], f Kleisli[R, error, B]) ReaderIOResult[R, A] {
	return readert.MonadChainFirst[ReaderIOResult[R, A], ReaderIOResult[R, B]](ioresult.MonadChainFirstLeft[A, B], ma, f)
}

// ChainFirstLeft runs f on the error of a failed computation for its side effect and keeps
// the original error.
//
// f runs with the same environment. Its result, success or failure, is discarded. A
// successful computation passes through without calling f.
//
// Type Parameters:
//   - A: The type of the success value
//   - R: The environment required by the computations
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - f: The side effect for the error
//
// Returns:
//   - Operator[R, A, A]: An operator that preserves the original value or error
//
// See Also:
//   - MonadChainFirstLeft: The uncurried version
//   - TapLeft: An alias
//   - ChainLeft: Recovers from the error instead of observing it
//
//go:inline
func ChainFirstLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A] {
	return readert.ChainFirst[ReaderIOResult[R, A], ReaderIOResult[R, B]](ioresult.ChainFirstLeft[A, B], f)
}

// MonadTapLeft is an alias for MonadChainFirstLeft.
//
// Type Parameters:
//   - R: The environment required by the computations
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - ma: The computation to observe
//   - f: The side effect for the error
//
// Returns:
//   - ReaderIOResult[R, A]: A computation with the result of ma
//
// See Also:
//   - MonadChainFirstLeft: The function this is an alias for
//
//go:inline
func MonadTapLeft[R, A, B any](ma ReaderIOResult[R, A], f Kleisli[R, error, B]) ReaderIOResult[R, A] {
	return MonadChainFirstLeft(ma, f)
}

// TapLeft is an alias for ChainFirstLeft.
//
// Type Parameters:
//   - A: The type of the success value
//   - R: The environment required by the computations
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - f: The side effect for the error
//
// Returns:
//   - Operator[R, A, A]: An operator that preserves the original value or error
//
// See Also:
//   - ChainFirstLeft: The function this is an alias for
//   - TapErrorAs: Observes only errors of a given type
//
//go:inline
func TapLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A] {
	return ChainFirstLeft[A](f)
}
