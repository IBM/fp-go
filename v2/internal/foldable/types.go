// Copyright (c) 2024 - 2025 IBM Corp.
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

package foldable

import (
	M "github.com/IBM/fp-go/v2/monoid"
)

// Foldable represents a data structure that can be folded/reduced to a single value.
//
// Foldable provides operations to collapse a structure containing multiple values
// into a single summary value by applying a combining function.
//
// Type Parameters:
//   - A: The type of elements in the structure
//   - B: The type of the accumulated result
//   - HKTA: The higher-kinded type containing A
type Foldable[A, B, HKTA any] interface {
	// Reduce folds the structure from left to right using a binary function and initial value.
	Reduce(func(B, A) B, B) func(HKTA) B

	// ReduceRight folds the structure from right to left using a binary function and initial value.
	ReduceRight(func(B, A) B, B) func(HKTA) B

	// FoldMap maps each element to a monoid and combines them using the monoid's operation.
	FoldMap(m M.Monoid[B]) func(func(A) B) func(HKTA) B
}
