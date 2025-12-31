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

package ioeither

import "github.com/IBM/fp-go/v2/either"

// FilterOrElse filters an IOEither value based on a predicate.
// If the predicate returns true for the Right value, it passes through unchanged.
// If the predicate returns false, it transforms the Right value into a Left using onFalse.
// Left values are passed through unchanged.
//
// This is useful for adding validation or constraints to successful computations,
// converting values that don't meet certain criteria into errors.
//
// Parameters:
//   - pred: A predicate function that tests the Right value
//   - onFalse: A function that converts the failing value into an error of type E
//
// Returns:
//   - An Operator that filters IOEither values based on the predicate
//
// Example:
//
//	// Validate that a number is positive
//	isPositive := N.MoreThan(0)
//	onNegative := S.Format[int]("%d is not positive")
//
//	validatePositive := ioeither.FilterOrElse(isPositive, onNegative)
//
//	result1 := validatePositive(ioeither.Right[string](42))()  // Right(42)
//	result2 := validatePositive(ioeither.Right[string](-5))()  // Left("-5 is not positive")
//	result3 := validatePositive(ioeither.Left[int]("error"))() // Left("error")
//
//go:inline
func FilterOrElse[E, A any](pred Predicate[A], onFalse func(A) E) Operator[E, A, A] {
	return ChainEitherK(either.FromPredicate(pred, onFalse))
}
