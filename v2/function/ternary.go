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

package function

// Ternary creates a conditional function that applies different transformations based on a predicate.
//
// This function implements a ternary operator (condition ? trueCase : falseCase) in a functional style.
// It takes a predicate and two transformation functions, returning a new function that applies
// the appropriate transformation based on whether the predicate is satisfied.
//
// Type Parameters:
//   - A: The input type
//   - B: The output type
//
// Parameters:
//   - pred: A predicate function that determines which branch to take
//   - onTrue: The transformation to apply when the predicate returns true
//   - onFalse: The transformation to apply when the predicate returns false
//
// Returns:
//   - A function that conditionally applies onTrue or onFalse based on pred
//
// Example:
//
//	isPositive := N.MoreThan(0)
//	double := N.Mul(2)
//	negate := func(n int) int { return -n }
//
//	transform := Ternary(isPositive, double, negate)
//	result1 := transform(5)   // 10 (positive, so doubled)
//	result2 := transform(-3)  // 3 (negative, so negated)
//
//	// Classify numbers
//	classify := Ternary(
//	    N.MoreThan(0),
//	    Constant1[int, string]("positive"),
//	    Constant1[int, string]("non-positive"),
//	)
//	result := classify(5)   // "positive"
//	result2 := classify(-3) // "non-positive"
func Ternary[A, B any](pred func(A) bool, onTrue, onFalse func(A) B) func(A) B {
	return func(a A) B {
		if pred(a) {
			return onTrue(a)
		}
		return onFalse(a)
	}
}
