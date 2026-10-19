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

// Package catch implements the catch-by-type combinators (CatchAs, CatchIs, MapErrorAs and
// TapErrorAs) once for all monads whose error channel is error. Each monad passes in its
// own ChainLeft, TapLeft, MapLeft and Left.
package catch

import (
	"errors"

	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
)

// ErrorAs returns the E found in the chain of err as the Right value, or err itself as the
// Left value.
func ErrorAs[E error](err error) ET.Either[error, E] {
	var e E
	if errors.As(err, &e) {
		return ET.Right[error](e)
	}
	return ET.Left[E](err)
}

// ErrorIs returns a function that returns err as the Right value if errors.Is(err, target)
// holds, and as the Left value otherwise.
func ErrorIs(target error) func(error) ET.Either[error, error] {
	return ET.FromPredicate(F.Bind2nd(errors.Is, target), F.Identity[error])
}

// OnMatch returns an error handler that applies f to the errors selected by match and
// fails with the original error otherwise.
//
//go:inline
func OnMatch[E, HKTB any](
	fromLeft func(error) HKTB,
	match func(error) ET.Either[error, E],
	f func(E) HKTB,
) func(error) HKTB {
	return F.Flow2(match, ET.Fold(fromLeft, f))
}

// OnMatchI is the variant of OnMatch for handlers that return a (value, error) tuple. It
// returns the zero value and the original error for the errors not selected by match.
func OnMatchI[E, B any](
	match func(error) ET.Either[error, E],
	f func(E) (B, error),
) func(error) (B, error) {
	return func(err error) (B, error) {
		e, cause := ET.Unwrap(match(err))
		if cause != nil {
			var b B
			return b, cause
		}
		return f(e)
	}
}

// MapAs returns an error mapping that replaces errors of type E by f applied to them and
// keeps all other errors.
//
//go:inline
func MapAs[E error](f func(E) error) func(error) error {
	return OnMatch(F.Identity[error], ErrorAs[E], f)
}

// CatchAs recovers from errors of type E with f.
//
//go:inline
func CatchAs[E error, HKTB, OP any](
	chainLeft func(func(error) HKTB) OP,
	fromLeft func(error) HKTB,
	f func(E) HKTB,
) OP {
	return chainLeft(OnMatch(fromLeft, ErrorAs[E], f))
}

// CatchIs recovers from errors that match target with f.
//
//go:inline
func CatchIs[HKTB, OP any](
	chainLeft func(func(error) HKTB) OP,
	fromLeft func(error) HKTB,
	target error,
	f func(error) HKTB,
) OP {
	return chainLeft(OnMatch(fromLeft, ErrorIs(target), f))
}

// MapErrorAs replaces errors of type E by f applied to them.
//
//go:inline
func MapErrorAs[E error, OP any](
	mapLeft func(func(error) error) OP,
	f func(E) error,
) OP {
	return mapLeft(MapAs(f))
}

// TapErrorAs runs f for its side effect on errors of type E and keeps the original error.
//
//go:inline
func TapErrorAs[E error, HKTB, OP any](
	tapLeft func(func(error) HKTB) OP,
	fromLeft func(error) HKTB,
	f func(E) HKTB,
) OP {
	return tapLeft(OnMatch(fromLeft, ErrorAs[E], f))
}
