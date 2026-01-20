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

	"github.com/IBM/fp-go/v2/eq"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

// Eq implements the equals predicate for values contained in the [ReaderIOResult] monad.
// It creates an equality checker that can compare two ReaderIOResult values by executing them
// with a given context and comparing their results using the provided Either equality checker.
//
// Parameters:
//   - eq: Equality checker for Either[A] values
//
// Returns a function that takes a context and returns an equality checker for ReaderIOResult[A].
//
// Example:
//
//	eqInt := eq.FromEquals(func(a, b either.Either[error, int]) bool {
//	    return either.Eq(eq.FromEquals(func(x, y int) bool { return x == y }))(a, b)
//	})
//	eqRIE := Eq(eqInt)
//	ctx := t.Context()
//	equal := eqRIE(ctx).Equals(Right[int](42), Right[int](42)) // true
//
//go:inline
func Eq[A any](eq eq.Eq[Either[A]]) func(context.Context) eq.Eq[ReaderIOResult[A]] {
	return RIOE.Eq[context.Context](eq)
}
