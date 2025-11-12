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

package function

// Flip reverses the order of parameters of a curried function.
//
// Given a curried function f that takes T1 then T2 and returns R,
// Flip returns a new curried function that takes T2 then T1 and returns R.
// This is useful when you have a curried function but need to apply its
// arguments in a different order.
//
// Mathematical notation:
//   - Given: f: T1 → T2 → R
//   - Returns: g: T2 → T1 → R where g(t2)(t1) = f(t1)(t2)
//
// Type Parameters:
//   - T1: The type of the first parameter (becomes second after flip)
//   - T2: The type of the second parameter (becomes first after flip)
//   - R: The return type
//
// Parameters:
//   - f: A curried function taking T1 then T2 and returning R
//
// Returns:
//   - A new curried function taking T2 then T1 and returning R
//
// Relationship to Swap:
//
// Flip is the curried version of Swap. While Swap works with binary functions,
// Flip works with curried functions:
//   - Swap: func(T1, T2) R → func(T2, T1) R
//   - Flip: func(T1) func(T2) R → func(T2) func(T1) R
//
// Example - Basic usage:
//
//	// Create a curried division function
//	divide := Curry2(func(a, b float64) float64 { return a / b })
//	// divide(10)(2) = 5.0 (10 / 2)
//
//	// Flip the parameter order
//	divideFlipped := Flip(divide)
//	// divideFlipped(10)(2) = 0.2 (2 / 10)
//
// Example - String formatting:
//
//	// Curried string formatter: format(template)(value)
//	format := Curry2(func(template, value string) string {
//	    return fmt.Sprintf(template, value)
//	})
//
//	// Normal order: template first, then value
//	result1 := format("Hello, %s!")("World")  // "Hello, World!"
//
//	// Flipped order: value first, then template
//	formatFlipped := Flip(format)
//	result2 := formatFlipped("Hello, %s!")("World")  // "Hello, World!"
//
//	// Useful for partial application in different order
//	greetWorld := format("Hello, %s!")
//	greetWorld("Alice")  // "Hello, Alice!"
//
//	formatAlice := formatFlipped("Alice")
//	formatAlice("Hello, %s!")  // "Hello, Alice!"
//
// Example - Practical use case with map operations:
//
//	// Curried map lookup: getFrom(map)(key)
//	getFrom := Curry2(func(m map[string]int, key string) int {
//	    return m[key]
//	})
//
//	data := map[string]int{"a": 1, "b": 2, "c": 3}
//
//	// Create a getter for this specific map
//	getValue := getFrom(data)
//	getValue("a")  // 1
//
//	// Flip to create key-first version: get(key)(map)
//	get := Flip(getFrom)
//	getA := get("a")
//	getA(data)  // 1
//
// Example - Combining with other functional patterns:
//
//	// Curried append: append(slice)(element)
//	appendTo := Curry2(func(slice []int, elem int) []int {
//	    return append(slice, elem)
//	})
//
//	// Flip to get: prepend(element)(slice)
//	prepend := Flip(appendTo)
//
//	nums := []int{1, 2, 3}
//	add4 := appendTo(nums)
//	result1 := add4(4)  // [1, 2, 3, 4]
//
//	prependZero := prepend(0)
//	result2 := prependZero(nums)  // [1, 2, 3, 0]
//
// See also:
//   - Swap: For flipping parameters of non-curried binary functions
//   - Curry2: For converting binary functions to curried form
//   - Uncurry2: For converting curried functions back to binary form
func Flip[T1, T2, R any](f func(T1) func(T2) R) func(T2) func(T1) R {
	return func(t2 T2) func(T1) R {
		return func(t1 T1) R {
			return f(t1)(t2)
		}
	}
}
