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

package either

// FilterOrElse filters an Either value based on a predicate.
// If the Either is Right and the predicate returns true, returns the original Right.
// If the Either is Right and the predicate returns false, returns Left with the error from onFalse.
// If the Either is Left, returns the original Left without applying the predicate.
//
// This is useful for adding validation to Right values, converting them to Left if they don't meet certain criteria.
//
// Example:
//
//	isPositive := N.MoreThan(0)
//	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }
//	filter := either.FilterOrElse(isPositive, onNegative)
//
//	result1 := filter(either.Right[error](5))  // Right(5)
//	result2 := filter(either.Right[error](-3)) // Left(error: "-3 is not positive")
//	result3 := filter(either.Left[int](someError)) // Left(someError)
//
//go:inline
func FilterOrElse[E, A any](pred Predicate[A], onFalse func(A) E) Operator[E, A, A] {
	return Chain(FromPredicate(pred, onFalse))
}
