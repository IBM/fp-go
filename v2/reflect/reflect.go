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

// Package reflect provides functional programming utilities for working with Go's reflect.Value type.
// It offers higher-order functions like Map, Reduce, and ReduceWithIndex that operate on
// reflective values representing slices or arrays.
//
// These utilities are particularly useful when working with dynamic types or when implementing
// generic algorithms that need to operate on collections discovered at runtime.
package reflect

import (
	R "reflect"

	"github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
)

func MonadReduceWithIndex[A any](val R.Value, f func(int, A, R.Value) A, initial A) A {

	kind := val.Kind()

	// Check if it supports Len() and Index()
	if kind != R.Slice && kind != R.Array && kind != R.String {
		// Not a sequential iterable, return initial
		return initial
	}

	count := val.Len()
	current := initial
	for i := range count {
		current = f(i, current, val.Index(i))
	}
	return current
}

// ReduceWithIndex applies a reducer function to each element of a reflect.Value (representing a slice or array),
// accumulating a result value. The reducer function receives the current index, the accumulated value,
// and the current element as a reflect.Value.
//
// This is a curried function that first takes the reducer function and initial value,
// then returns a function that accepts the reflect.Value to reduce.
//
// Parameters:
//   - f: A reducer function that takes (index int, accumulator A, element reflect.Value) and returns the new accumulator
//   - initial: The initial value for the accumulation
//
// Returns:
//   - A function that takes a reflect.Value and returns the final accumulated value
//
// Example:
//
//	// Sum all integers in a reflected slice with their indices
//	sumWithIndex := ReduceWithIndex(func(i int, acc int, v reflect.Value) int {
//	    return acc + i + int(v.Int())
//	}, 0)
//	result := sumWithIndex(reflect.ValueOf([]int{10, 20, 30}))
//	// result = 0 + (0+10) + (1+20) + (2+30) = 63
func ReduceWithIndex[A any](f func(int, A, R.Value) A, initial A) func(R.Value) A {
	return func(val R.Value) A {
		return MonadReduceWithIndex(val, f, initial)
	}
}

// Reduce applies a reducer function to each element of a reflect.Value (representing a slice or array),
// accumulating a result value. Unlike ReduceWithIndex, the reducer function does not receive the index.
//
// This is a curried function that first takes the reducer function and initial value,
// then returns a function that accepts the reflect.Value to reduce.
//
// Parameters:
//   - f: A reducer function that takes (accumulator A, element reflect.Value) and returns the new accumulator
//   - initial: The initial value for the accumulation
//
// Returns:
//   - A function that takes a reflect.Value and returns the final accumulated value
//
// Example:
//
//	// Sum all integers in a reflected slice
//	sum := Reduce(func(acc int, v reflect.Value) int {
//	    return acc + int(v.Int())
//	}, 0)
//	result := sum(reflect.ValueOf([]int{10, 20, 30}))
//	// result = 60
func Reduce[A any](f func(A, R.Value) A, initial A) func(R.Value) A {
	return ReduceWithIndex(F.Ignore1of3[int](f), initial)
}

// MonadMapWithIndex is the non-curried version of MapWithIndex. It transforms each element of a
// reflect.Value (representing a slice, array, or string) using the provided function that receives
// both the index and the element, returning a new slice containing the transformed values.
//
// Unlike MapWithIndex which is curried, this function takes both the reflect.Value and the
// transformation function as parameters in a single call. This is useful when you need to pass
// the function directly without partial application.
//
// Parameters:
//   - val: The reflect.Value to map over (must be a slice, array, or string)
//   - f: A transformation function that takes (index int, element reflect.Value) and returns a value of type A
//
// Returns:
//   - A slice of transformed values, or an empty slice if val is not iterable
//
// Example:
//
//	// Transform a reflected slice with index awareness
//	input := reflect.ValueOf([]int{10, 20, 30})
//	result := MonadMapWithIndex(input, func(i int, v reflect.Value) string {
//	    return fmt.Sprintf("[%d]=%d", i, int(v.Int()))
//	})
//	// result = []string{"[0]=10", "[1]=20", "[2]=30"}
func MonadMapWithIndex[A any](val R.Value, f func(int, R.Value) A) []A {

	kind := val.Kind()

	// Check if it supports Len() and Index()
	if kind != R.Slice && kind != R.Array && kind != R.String {
		// Not a sequential iterable, return initial
		return array.Empty[A]()
	}

	l := val.Len()
	res := make([]A, l)
	for i := l - 1; i >= 0; i-- {
		res[i] = f(i, val.Index(i))
	}
	return res
}

// MapWithIndex transforms each element of a reflect.Value (representing a slice or array) using the provided
// function that receives both the index and the element, returning a new slice containing the transformed values.
//
// This is a curried function that first takes the transformation function,
// then returns a function that accepts the reflect.Value to map over.
//
// Parameters:
//   - f: A transformation function that takes (index int, element reflect.Value) and returns a value of type A
//
// Returns:
//   - A function that takes a reflect.Value and returns a slice of transformed values
//
// Example:
//
//	// Create indexed labels from a reflected slice
//	indexedLabels := MapWithIndex(func(i int, v reflect.Value) string {
//	    return fmt.Sprintf("[%d]: %d", i, int(v.Int()))
//	})
//	result := indexedLabels(reflect.ValueOf([]int{10, 20, 30}))
//	// result = []string{"[0]: 10", "[1]: 20", "[2]: 30"}
func MapWithIndex[A any](f func(int, R.Value) A) func(R.Value) []A {
	return F.Bind2nd(MonadMapWithIndex, f)
}

// Map transforms each element of a reflect.Value (representing a slice or array) using the provided
// function, returning a new slice containing the transformed values.
//
// This is a curried function that first takes the transformation function,
// then returns a function that accepts the reflect.Value to map over.
//
// Parameters:
//   - f: A transformation function that takes a reflect.Value and returns a value of type A
//
// Returns:
//   - A function that takes a reflect.Value and returns a slice of transformed values
//
// Example:
//
//	// Extract integers from a reflected slice and double them
//	doubleInts := Map(func(v reflect.Value) int {
//	    return int(v.Int()) * 2
//	})
//	result := doubleInts(reflect.ValueOf([]int{1, 2, 3}))
//	// result = []int{2, 4, 6}
func Map[A any](f func(R.Value) A) func(R.Value) []A {
	return MapWithIndex(F.Ignore1of2[int](f))
}
