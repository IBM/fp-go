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

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFlip tests the Flip function with various scenarios
func TestFlip(t *testing.T) {
	t.Run("flips string concatenation", func(t *testing.T) {
		// Create a curried function that formats strings
		format := Curry2(func(a, b string) string {
			return fmt.Sprintf("%s:%s", a, b)
		})

		// Original order: a then b
		assert.Equal(t, "a:b", format("a")("b"))
		assert.Equal(t, "hello:world", format("hello")("world"))

		// Flipped order: b then a
		flipped := Flip(format)
		assert.Equal(t, "b:a", flipped("a")("b"))
		assert.Equal(t, "world:hello", flipped("hello")("world"))
	})

	t.Run("flips numeric operations", func(t *testing.T) {
		// Curried subtraction: subtract(a)(b) = a - b
		subtract := Curry2(func(a, b int) int {
			return a - b
		})

		// Original: 10 - 3 = 7
		assert.Equal(t, 7, subtract(10)(3))

		// Flipped: 3 - 10 = -7
		flipped := Flip(subtract)
		assert.Equal(t, -7, flipped(10)(3))
	})

	t.Run("flips division", func(t *testing.T) {
		// Curried division: divide(a)(b) = a / b
		divide := Curry2(func(a, b float64) float64 {
			return a / b
		})

		// Original: 10 / 2 = 5.0
		assert.Equal(t, 5.0, divide(10)(2))

		// Flipped: 2 / 10 = 0.2
		flipped := Flip(divide)
		assert.Equal(t, 0.2, flipped(10)(2))
	})

	t.Run("flips with partial application", func(t *testing.T) {
		// Curried append-like operation
		prepend := Curry2(func(prefix, text string) string {
			return prefix + text
		})

		// Create specialized functions with original order
		addHello := prepend("Hello, ")
		assert.Equal(t, "Hello, World", addHello("World"))
		assert.Equal(t, "Hello, Go", addHello("Go"))

		// Flip and create specialized functions with reversed order
		flipped := Flip(prepend)
		addToWorld := flipped("World")
		assert.Equal(t, "Hello, World", addToWorld("Hello, "))
		assert.Equal(t, "Goodbye, World", addToWorld("Goodbye, "))
	})

	t.Run("flips with different types", func(t *testing.T) {
		// Curried function with different input types
		repeat := Curry2(func(s string, n int) string {
			result := ""
			for i := 0; i < n; i++ {
				result += s
			}
			return result
		})

		// Original: repeat("x")(3) = "xxx"
		assert.Equal(t, "xxx", repeat("x")(3))
		assert.Equal(t, "abab", repeat("ab")(2))

		// Flipped: repeat(3)("x") = "xxx"
		flipped := Flip(repeat)
		assert.Equal(t, "xxx", flipped(3)("x"))
		assert.Equal(t, "abab", flipped(2)("ab"))
	})

	t.Run("double flip returns to original", func(t *testing.T) {
		// Flipping twice should return to original behavior
		original := Curry2(func(a, b string) string {
			return a + "-" + b
		})

		flipped := Flip(original)
		doubleFlipped := Flip(flipped)

		// Original and double-flipped should behave the same
		assert.Equal(t, original("a")("b"), doubleFlipped("a")("b"))
		assert.Equal(t, "a-b", doubleFlipped("a")("b"))
	})

	t.Run("flips with complex types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		// Curried function creating a person
		makePerson := Curry2(func(name string, age int) Person {
			return Person{Name: name, Age: age}
		})

		// Original order: name then age
		alice := makePerson("Alice")(30)
		assert.Equal(t, "Alice", alice.Name)
		assert.Equal(t, 30, alice.Age)

		// Flipped order: age then name
		flipped := Flip(makePerson)
		bob := flipped(25)("Bob")
		assert.Equal(t, "Bob", bob.Name)
		assert.Equal(t, 25, bob.Age)
	})

	t.Run("flips map operations", func(t *testing.T) {
		// Curried map getter: get(map)(key)
		get := Curry2(func(m map[string]int, key string) int {
			return m[key]
		})

		data := map[string]int{"a": 1, "b": 2, "c": 3}

		// Original: provide map first, then key
		getValue := get(data)
		assert.Equal(t, 1, getValue("a"))
		assert.Equal(t, 2, getValue("b"))

		// Flipped: provide key first, then map
		flipped := Flip(get)
		getA := flipped("a")
		assert.Equal(t, 1, getA(data))

		data2 := map[string]int{"a": 10, "b": 20}
		assert.Equal(t, 10, getA(data2))
	})

	t.Run("flips boolean operations", func(t *testing.T) {
		// Curried logical operation
		implies := Curry2(func(a, b bool) bool {
			return !a || b
		})

		// Test truth table for implication
		assert.True(t, implies(true)(true))   // T → T = T
		assert.False(t, implies(true)(false)) // T → F = F
		assert.True(t, implies(false)(true))  // F → T = T
		assert.True(t, implies(false)(false)) // F → F = T

		// Flipped version (reverse implication)
		flipped := Flip(implies)
		assert.True(t, flipped(true)(true))   // T ← T = T
		assert.True(t, flipped(true)(false))  // T ← F = T
		assert.False(t, flipped(false)(true)) // F ← T = F
		assert.True(t, flipped(false)(false)) // F ← F = T
	})

	t.Run("flips with slice operations", func(t *testing.T) {
		// Curried slice append
		appendTo := Curry2(func(slice []int, elem int) []int {
			return append(slice, elem)
		})

		nums := []int{1, 2, 3}

		// Original: provide slice first, then element
		add4 := appendTo(nums)
		result1 := add4(4)
		assert.Equal(t, []int{1, 2, 3, 4}, result1)

		// Flipped: provide element first, then slice
		flipped := Flip(appendTo)
		appendFive := flipped(5)
		result2 := appendFive(nums)
		assert.Equal(t, []int{1, 2, 3, 5}, result2)
	})
}

// TestFlipProperties tests mathematical properties of Flip
func TestFlipProperties(t *testing.T) {
	t.Run("flip is involutive (flip . flip = id)", func(t *testing.T) {
		// Flipping twice should give back the original function behavior
		original := Curry2(func(a, b int) int {
			return a*10 + b
		})

		flipped := Flip(original)
		doubleFlipped := Flip(flipped)

		// Test with multiple inputs
		testCases := []struct{ a, b int }{
			{1, 2},
			{5, 7},
			{0, 0},
			{-1, 3},
		}

		for _, tc := range testCases {
			assert.Equal(t,
				original(tc.a)(tc.b),
				doubleFlipped(tc.a)(tc.b),
				"flip(flip(f)) should equal f for inputs (%d, %d)", tc.a, tc.b)
		}
	})

	t.Run("flip preserves function composition", func(t *testing.T) {
		// If we have f: A → B → C and g: C → D
		// then g ∘ f(a)(b) = g(f(a)(b))
		// and g ∘ flip(f)(b)(a) = g(flip(f)(b)(a))

		f := Curry2(func(a, b int) int {
			return a + b
		})

		g := func(n int) int {
			return n * 2
		}

		flippedF := Flip(f)

		// Compose g with f
		composed1 := func(a, b int) int {
			return g(f(a)(b))
		}

		// Compose g with flipped f
		composed2 := func(a, b int) int {
			return g(flippedF(b)(a))
		}

		// Both should give the same result
		assert.Equal(t, composed1(3, 5), composed2(3, 5))
		assert.Equal(t, 16, composed1(3, 5)) // (3 + 5) * 2 = 16
	})
}

// BenchmarkFlip benchmarks the Flip function
func BenchmarkFlip(b *testing.B) {
	add := Curry2(func(a, b int) int {
		return a + b
	})

	flipped := Flip(add)

	b.Run("original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = add(i)(i + 1)
		}
	})

	b.Run("flipped", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = flipped(i)(i + 1)
		}
	})
}
