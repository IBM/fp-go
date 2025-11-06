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

package array

import (
	M "github.com/IBM/fp-go/v2/monoid"
)

// ConcatAll concatenates all elements of an array using the provided Monoid.
// This reduces the array to a single value by repeatedly applying the Monoid's concat operation.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/monoid"
//
//	// Sum all numbers
//	sumAll := array.ConcatAll(monoid.MonoidSum[int]())
//	result := sumAll([]int{1, 2, 3, 4, 5}) // 15
//
//	// Concatenate all strings
//	concatStrings := array.ConcatAll(monoid.MonoidString())
//	result2 := concatStrings([]string{"Hello", " ", "World"}) // "Hello World"
//
//go:inline
func ConcatAll[A any](m M.Monoid[A]) func([]A) A {
	return Reduce(m.Concat, m.Empty())
}
