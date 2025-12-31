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

import "github.com/IBM/fp-go/v2/either"

// FilterOrElse filters a Result value based on a predicate.
// If the Result is Ok (Right) and the predicate returns true, returns the original Ok.
// If the Result is Ok (Right) and the predicate returns false, returns Error (Left) with the error from onFalse.
// If the Result is Error (Left), returns the original Error without applying the predicate.
//
// This is useful for adding validation to successful results, converting them to errors if they don't meet certain criteria.
// Result[T] is an alias for Either[error, T], so this function delegates to either.FilterOrElse.
//
// Example:
//
//	isPositive := N.MoreThan(0)
//	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }
//	filter := result.FilterOrElse(isPositive, onNegative)
//
//	result1 := filter(result.Of(5))  // Ok(5)
//	result2 := filter(result.Of(-3)) // Error(error: "-3 is not positive")
//	result3 := filter(result.Error[int](someError)) // Error(someError)
//
//go:inline
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A] {
	return either.FilterOrElse(pred, onFalse)
}
