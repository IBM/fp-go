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

package readerioeither

import "github.com/IBM/fp-go/v2/either"

// FilterOrElse filters a ReaderIOEither value based on a predicate.
// If the predicate returns true for the Right value, it passes through unchanged.
// If the predicate returns false, it transforms the Right value into a Left using onFalse.
// Left values are passed through unchanged.
//
// This is useful for adding validation or constraints to successful computations that
// depend on a context, converting values that don't meet certain criteria into errors.
//
// Parameters:
//   - pred: A predicate function that tests the Right value
//   - onFalse: A function that converts the failing value into an error of type E
//
// Returns:
//   - An Operator that filters ReaderIOEither values based on the predicate
//
// Example:
//
//	type Config struct {
//	    MaxValue int
//	}
//
//	// Validate that a number doesn't exceed the configured maximum
//	validateMax := func(cfg Config) readerioeither.ReaderIOEither[Config, string, int] {
//	    isValid := func(n int) bool { return n <= cfg.MaxValue }
//	    onInvalid := func(n int) string {
//	        return fmt.Sprintf("%d exceeds max %d", n, cfg.MaxValue)
//	    }
//
//	    filter := readerioeither.FilterOrElse(isValid, onInvalid)
//	    return filter(readerioeither.Right[Config, string](42))
//	}
//
//	cfg := Config{MaxValue: 100}
//	result := validateMax(cfg)(cfg)() // Right(42)
//
//	cfg2 := Config{MaxValue: 10}
//	result2 := validateMax(cfg2)(cfg2)() // Left("42 exceeds max 10")
//
//go:inline
func FilterOrElse[R, E, A any](pred Predicate[A], onFalse func(A) E) Operator[R, E, A, A] {
	return ChainEitherK[R](either.FromPredicate(pred, onFalse))
}
