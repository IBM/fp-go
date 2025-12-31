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

package result

// FilterOrElse filters a Result value based on a predicate in an idiomatic style.
// If the Result is Ok and the predicate returns true, returns the original Ok.
// If the Result is Ok and the predicate returns false, returns Error with the error from onFalse.
// If the Result is Error, returns the original Error without applying the predicate.
//
// This is the idiomatic version that returns an Operator for use in method chaining.
// It's useful for adding validation to successful results, converting them to errors if they don't meet certain criteria.
//
// Example:
//
//	import (
//		R "github.com/IBM/fp-go/v2/idiomatic/result"
//		N "github.com/IBM/fp-go/v2/number"
//	)
//
//	isPositive := N.MoreThan(0)
//	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }
//
//	result := R.Of(5).
//		Pipe(R.FilterOrElse(isPositive, onNegative)) // Ok(5)
//
//	result2 := R.Of(-3).
//		Pipe(R.FilterOrElse(isPositive, onNegative)) // Error(error: "-3 is not positive")
//
//go:inline
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A] {
	return Chain(FromPredicate(pred, onFalse))
}
