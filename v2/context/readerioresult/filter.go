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

package readerioresult

import (
	"context"

	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

// FilterOrElse filters a ReaderIOResult value based on a predicate.
// This is a convenience wrapper around readerioresult.FilterOrElse that fixes
// the context type to context.Context.
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
//   - An Operator that filters ReaderIOResult values based on the predicate
//
// Example:
//
//	// Validate that a number is positive
//	isPositive := N.MoreThan(0)
//	onNegative := func(n int) error { return fmt.Errorf("%d is not positive", n) }
//
//	filter := readerioresult.FilterOrElse(isPositive, onNegative)
//	result := filter(readerioresult.Right(42))(t.Context())()
//
//go:inline
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A] {
	return RIOR.FilterOrElse[context.Context](pred, onFalse)
}
