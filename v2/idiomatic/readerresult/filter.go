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

package readerresult

import "github.com/IBM/fp-go/v2/either"

// FilterOrElse filters a ReaderResult value based on a predicate in an idiomatic style.
// If the ReaderResult computation succeeds and the predicate returns true, returns the original success value.
// If the ReaderResult computation succeeds and the predicate returns false, returns an error with the error from onFalse.
// If the ReaderResult computation fails, returns the original error without applying the predicate.
//
// This is the idiomatic version that returns an Operator for use in method chaining.
// It's useful for adding validation to successful computations, converting them to errors if they don't meet certain criteria.
//
// Example:
//
//	import (
//		RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
//		N "github.com/IBM/fp-go/v2/number"
//	)
//
//	type Config struct {
//		MaxValue int
//	}
//
//	isPositive := N.MoreThan(0)
//	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }
//
//	result := RR.Of[Config](5).
//		Pipe(RR.FilterOrElse(isPositive, onNegative))(Config{MaxValue: 10}) // Ok(5)
//
//	result2 := RR.Of[Config](-3).
//		Pipe(RR.FilterOrElse(isPositive, onNegative))(Config{MaxValue: 10}) // Error(error: "-3 is not positive")
//
//go:inline
func FilterOrElse[R, A any](pred Predicate[A], onFalse func(A) error) Operator[R, A, A] {
	return ChainEitherK[R](either.FromPredicate(pred, onFalse))
}
