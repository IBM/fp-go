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

package endomorphism

import (
	"github.com/IBM/fp-go/v2/function"
)

// Curry2 curries a binary function that returns an endomorphism.
//
// This function takes a binary function f(T0, T1) T1 and converts it into a
// curried form that returns an endomorphism. The result is a function that
// takes the first argument and returns an endomorphism (a function T1 -> T1).
//
// Deprecated: This function is no longer needed. Use function.Curry2 directly
// from the function package instead.
//
// Parameters:
//   - f: A binary function where the return type matches the second parameter type
//
// Returns:
//   - A curried function that takes T0 and returns an Endomorphism[T1]
//
// Example:
//
//	// A binary function that adds two numbers
//	add := func(x, y int) int { return x + y }
//	curriedAdd := endomorphism.Curry2(add)
//	addFive := curriedAdd(5) // Returns an endomorphism that adds 5
//	result := addFive(10)    // Returns: 15
func Curry2[FCT ~func(T0, T1) T1, T0, T1 any](f FCT) func(T0) Endomorphism[T1] {
	return function.Curry2(f)
}

// Curry3 curries a ternary function that returns an endomorphism.
//
// This function takes a ternary function f(T0, T1, T2) T2 and converts it into
// a curried form. The result is a function that takes the first argument, returns
// a function that takes the second argument, and finally returns an endomorphism
// (a function T2 -> T2).
//
// Deprecated: This function is no longer needed. Use function.Curry3 directly
// from the function package instead.
//
// Parameters:
//   - f: A ternary function where the return type matches the third parameter type
//
// Returns:
//   - A curried function that takes T0, then T1, and returns an Endomorphism[T2]
//
// Example:
//
//	// A ternary function
//	combine := func(x, y, z int) int { return x + y + z }
//	curriedCombine := endomorphism.Curry3(combine)
//	addTen := curriedCombine(5)(5) // Returns an endomorphism that adds 10
//	result := addTen(20)           // Returns: 30
func Curry3[FCT ~func(T0, T1, T2) T2, T0, T1, T2 any](f FCT) func(T0) func(T1) Endomorphism[T2] {
	return function.Curry3(f)
}
