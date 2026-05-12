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

package iterresult

import (
	"github.com/IBM/fp-go/v2/iterator/itereither"
)

// FilterOrElse filters a SeqResult value based on a predicate.
// If the predicate returns true for the Right value, it passes through unchanged.
// If the predicate returns false, it transforms the Right value into a Left using onFalse.
// Left values are passed through unchanged.
//
// This is useful for adding validation or constraints to successful computations,
// converting values that don't meet certain criteria into errors.
//
// Marble diagram:
//
//	Input:  ---R(5)---R(-3)---L(e)---R(10)---|
//	pred(x) = x > 0
//	onFalse(x) = "negative: " + x
//	Output: ---R(5)---L("negative: -3")---L(e)---R(10)---|
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
// Values that fail the predicate are converted to Left.
//
// Parameters:
//   - pred: A predicate function that tests the Right value
//   - onFalse: A function that converts the failing value into an error of type E
//
// Returns:
//   - An Operator that filters SeqResult values based on the predicate
//
// Example:
//
//	// Validate that a number is positive
//	isPositive := func(x int) bool { return x > 0 }
//	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }
//
//	validatePositive := FilterOrElse(isPositive, onNegative)
//
//	result1 := validatePositive(Right(42))  // Right(42)
//	result2 := validatePositive(Right(-5))  // Left(error: "-5 is not positive")
//	result3 := validatePositive(Left[int](errors.New("error"))) // Left(error: "error")
//
//go:inline
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A] {
	return itereither.FilterOrElse(pred, onFalse)
}
