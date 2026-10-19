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

package readerreaderioresult

import (
	"github.com/IBM/fp-go/v2/internal/catch"
)

// CatchAs recovers from errors of type E.
//
// If the computation fails with an error for which errors.As finds an E in its chain, f is
// called with that E and its computation replaces the failure. Other errors, and successful
// values, pass through unchanged. This corresponds to Effect.catchTag in effect-ts, with
// the error type in place of the tag.
//
// R and A are usually inferred from f, so a call site reads
// CatchAs[*NotFoundError](handler).
//
// Type Parameters:
//   - E: The error type to catch, a concrete error type or an interface that implements error
//   - R: The environment required by the computations
//   - A: The type of the success value
//
// Parameters:
//   - f: The recovery computation for an error of type E
//
// Returns:
//   - Operator[R, A, A]: An operator that recovers from errors of type E
//
// See Also:
//   - CatchIs: Recovers from errors that match a sentinel value
//   - ChainLeft: Recovers from any error
//   - MapErrorAs: Transforms errors of type E
//
//go:inline
func CatchAs[E error, R, A any](f Kleisli[R, E, A]) Operator[R, A, A] {
	return catch.CatchAs(ChainLeft[R, A], Left[R, A], f)
}

// CatchIs recovers from errors that match target.
//
// If the computation fails with an error for which errors.Is(err, target) holds, f is called
// with the original error and its computation replaces the failure. Other errors, and
// successful values, pass through unchanged.
//
// Type Parameters:
//   - R: The environment required by the computations
//   - A: The type of the success value
//
// Parameters:
//   - target: The sentinel error to match with errors.Is
//   - f: The recovery computation, called with the original error
//
// Returns:
//   - Operator[R, A, A]: An operator that recovers from errors that match target
//
// See Also:
//   - CatchAs: Recovers from errors of a given type
//   - ChainLeft: Recovers from any error
//
//go:inline
func CatchIs[R, A any](target error, f Kleisli[R, error, A]) Operator[R, A, A] {
	return catch.CatchIs(ChainLeft[R, A], Left[R, A], target, f)
}

// MapErrorAs transforms errors of type E.
//
// If the computation fails with an error for which errors.As finds an E in its chain, the
// failure is replaced by f applied to that E. Other errors, and successful values, pass
// through unchanged.
//
// Because R and A only appear in the result type, they must be given explicitly, for
// example MapErrorAs[*NotFoundError, Config, User](toHTTPError).
//
// Type Parameters:
//   - E: The error type to transform
//   - R: The environment required by the computation
//   - A: The type of the success value
//
// Parameters:
//   - f: The function that maps an error of type E to a new error
//
// Returns:
//   - Operator[R, A, A]: An operator that transforms errors of type E
//
// See Also:
//   - CatchAs: Recovers from errors of type E
//
//go:inline
func MapErrorAs[E error, R, A any](f func(E) error) Operator[R, A, A] {
	return catch.MapErrorAs(MapLeft[R, A], f)
}

// TapErrorAs runs a computation for its side effect when the computation fails with an
// error of type E.
//
// If the computation fails with an error for which errors.As finds an E in its chain, f is
// called with that E and its computation is executed. The original error always propagates:
// f cannot recover from the failure, and a failure of f is discarded. Other errors, and
// successful values, pass through without calling f.
//
// Because A only appears in the result type, it must be given explicitly, for example
// TapErrorAs[*NotFoundError, Config, User](logNotFound).
//
// Type Parameters:
//   - E: The error type to observe
//   - R: The environment required by the computations
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - f: The side effect for an error of type E
//
// Returns:
//   - Operator[R, A, A]: An operator that preserves the original value or error
//
// See Also:
//   - TapLeft: Observes any error
//   - CatchAs: Recovers from errors of type E
//
//go:inline
func TapErrorAs[E error, R, A, B any](f Kleisli[R, E, B]) Operator[R, A, A] {
	return catch.TapErrorAs(TapLeft[A, R, B], Left[R, B], f)
}
