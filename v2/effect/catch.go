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

package effect

import (
	"errors"

	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
)

// isError selects the errors that match target according to errors.Is. Any other error
// is returned unchanged as the Left value.
func isError(target error) result.Kleisli[error, error] {
	return result.FromPredicate(F.Bind2nd(errors.Is, target), F.Identity[error])
}

// asError selects the errors for which errors.As finds an E in the chain. Any other
// error is returned unchanged as the Left value.
func asError[E error](err error) Result[E] {
	var e E
	if errors.As(err, &e) {
		return result.Of(e)
	}
	return result.Left[E](err)
}

// onMatch returns an error handler that applies f to the errors selected by match and
// fails with the original error otherwise.
func onMatch[C, E, B any](match result.Kleisli[error, E], f Kleisli[C, E, B]) Kleisli[C, error, B] {
	return F.Flow2(match, result.Fold(Fail[C, B], f))
}

// CatchAs recovers from errors of type E.
//
// If the effect fails with an error for which errors.As finds an E in its chain, f is
// called with that E and its effect replaces the failure. Other errors, and successful
// values, pass through unchanged. This corresponds to Effect.catchTag in effect-ts, with
// the error type in place of the tag.
//
// C and A are usually inferred from f, so a call site reads
// CatchAs[*NotFoundError](handler).
//
// Type Parameters:
//   - E: The error type to catch, a concrete error type or an interface that implements error
//   - C: The dependencies required by the effects
//   - A: The type of the success value
//
// Parameters:
//   - f: The recovery effect for an error of type E
//
// Returns:
//   - Operator[C, A, A]: An operator that recovers from errors of type E
//
// See Also:
//   - CatchIs: Recovers from errors that match a sentinel value
//   - ChainLeft: Recovers from any error
//   - MapErrorAs: Transforms errors of type E
//
//go:inline
func CatchAs[E error, C, A any](f Kleisli[C, E, A]) Operator[C, A, A] {
	return ChainLeft(onMatch(asError[E], f))
}

// CatchIs recovers from errors that match target.
//
// If the effect fails with an error for which errors.Is(err, target) holds, f is called
// with the original error and its effect replaces the failure. Other errors, and
// successful values, pass through unchanged.
//
// Type Parameters:
//   - C: The dependencies required by the effects
//   - A: The type of the success value
//
// Parameters:
//   - target: The sentinel error to match with errors.Is
//   - f: The recovery effect, called with the original error
//
// Returns:
//   - Operator[C, A, A]: An operator that recovers from errors that match target
//
// See Also:
//   - CatchAs: Recovers from errors of a given type
//   - ChainLeft: Recovers from any error
//
//go:inline
func CatchIs[C, A any](target error, f Kleisli[C, error, A]) Operator[C, A, A] {
	return ChainLeft(onMatch(isError(target), f))
}

// MapErrorAs transforms errors of type E.
//
// If the effect fails with an error for which errors.As finds an E in its chain, the
// failure is replaced by f applied to that E. Other errors, and successful values, pass
// through unchanged.
//
// Because C and A only appear in the result type, they must be given explicitly, for
// example MapErrorAs[*NotFoundError, Deps, User](toHTTPError).
//
// Type Parameters:
//   - E: The error type to transform
//   - C: The dependencies required by the effect
//   - A: The type of the success value
//
// Parameters:
//   - f: The function that maps an error of type E to a new error
//
// Returns:
//   - Operator[C, A, A]: An operator that transforms errors of type E
//
// See Also:
//   - CatchAs: Recovers from errors of type E
//
//go:inline
func MapErrorAs[E error, C, A any](f func(E) error) Operator[C, A, A] {
	return readerreaderioresult.MapLeft[C, A](F.Flow2(asError[E], result.Fold(F.Identity[error], f)))
}

// TapErrorAs runs an effect for its side effect when the effect fails with an error of
// type E.
//
// If the effect fails with an error for which errors.As finds an E in its chain, f is
// called with that E and its effect is executed. The original error always propagates:
// f cannot recover from the failure, and a failure of f is discarded. Other errors, and
// successful values, pass through without calling f.
//
// Because A only appears in the result type, it must be given explicitly, for example
// TapErrorAs[*NotFoundError, Deps, User](logNotFound).
//
// Type Parameters:
//   - E: The error type to observe
//   - C: The dependencies required by the effects
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - f: The side effect for an error of type E
//
// Returns:
//   - Operator[C, A, A]: An operator that preserves the original value or error
//
// See Also:
//   - TapLeft: Observes any error
//   - CatchAs: Recovers from errors of type E
//
//go:inline
func TapErrorAs[E error, C, A, B any](f Kleisli[C, E, B]) Operator[C, A, A] {
	return TapLeft[A](onMatch(asError[E], f))
}
