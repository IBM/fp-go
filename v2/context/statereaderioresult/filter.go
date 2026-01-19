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

	"github.com/IBM/fp-go/v2/statereaderioeither"
)

// FilterOrElse filters a StateReaderIOResult value based on a predicate.
// This is a convenience wrapper around statereaderioeither.FilterOrElse that fixes
// the context type to context.Context and the error type to error.
//
// If the predicate returns true for the Right value, it passes through unchanged.
// If the predicate returns false, it transforms the Right value into a Left (error) using onFalse.
// Left values are passed through unchanged.
//
// Parameters:
//   - pred: A predicate function that tests the Right value
//   - onFalse: A function that converts the failing value into an error
//
// Returns:
//   - An Operator that filters StateReaderIOResult values based on the predicate
//
// Example:
//
//	type AppState struct {
//	    Counter int
//	}
//
//	// Validate that a number is positive
//	isPositive := N.MoreThan(0)
//	onNegative := func(n int) error { return fmt.Errorf("%d is not positive", n) }
//
//	filter := statereaderioresult.FilterOrElse[AppState](isPositive, onNegative)
//	result := filter(statereaderioresult.Right[AppState](42))(AppState{})(t.Context())()
//
//go:inline
func FilterOrElse[S, A any](pred Predicate[A], onFalse func(A) error) Operator[S, A, A] {
	return statereaderioeither.FilterOrElse[S, context.Context](pred, onFalse)
}
