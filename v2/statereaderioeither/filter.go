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

import "github.com/IBM/fp-go/v2/either"

// FilterOrElse filters a StateReaderIOEither value based on a predicate.
// If the predicate returns true for the Right value, it passes through unchanged.
// If the predicate returns false, it transforms the Right value into a Left using onFalse.
// Left values are passed through unchanged.
//
// This is useful for adding validation or constraints to successful stateful computations
// that depend on a context, converting values that don't meet certain criteria into errors.
//
// Parameters:
//   - pred: A predicate function that tests the Right value
//   - onFalse: A function that converts the failing value into an error of type E
//
// Returns:
//   - An Operator that filters StateReaderIOEither values based on the predicate
//
// Example:
//
//	type AppState struct {
//	    Counter int
//	}
//
//	type Config struct {
//	    MaxValue int
//	}
//
//	// Validate that a number doesn't exceed the configured maximum
//	validateMax := func(n int) statereaderioeither.StateReaderIOEither[AppState, Config, string, int] {
//	    isValid := func(val int) bool { return val <= 100 }
//	    onInvalid := func(val int) string {
//	        return fmt.Sprintf("%d exceeds maximum", val)
//	    }
//
//	    filter := statereaderioeither.FilterOrElse[AppState, Config](isValid, onInvalid)
//	    return filter(statereaderioeither.Right[AppState, Config, string](n))
//	}
//
//	state := AppState{Counter: 0}
//	cfg := Config{MaxValue: 100}
//	result := validateMax(42)(state)(cfg)() // Right(Pair(state, 42))
//	result2 := validateMax(150)(state)(cfg)() // Left("150 exceeds maximum")
//
//go:inline
func FilterOrElse[S, R, E, A any](pred Predicate[A], onFalse func(A) E) Operator[S, R, E, A, A] {
	return ChainEitherK[S, R](either.FromPredicate(pred, onFalse))
}
