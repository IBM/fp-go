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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBind1st tests the Bind1st function with various scenarios
func TestBind1st(t *testing.T) {
	t.Run("binds first parameter of multiplication", func(t *testing.T) {
		multiply := func(a, b int) int { return a * b }
		double := Bind1st(multiply, 2)
		triple := Bind1st(multiply, 3)

		assert.Equal(t, 10, double(5))
		assert.Equal(t, 20, double(10))
		assert.Equal(t, 15, triple(5))
		assert.Equal(t, 30, triple(10))
	})

	t.Run("binds first parameter of division", func(t *testing.T) {
		divide := func(a, b float64) float64 { return a / b }
		divideBy10 := Bind1st(divide, 10.0)
		divideBy5 := Bind1st(divide, 5.0)

		assert.Equal(t, 5.0, divideBy10(2.0))
		assert.Equal(t, 2.0, divideBy10(5.0))
		assert.Equal(t, 1.0, divideBy5(5.0))
	})

	t.Run("binds first parameter of subtraction", func(t *testing.T) {
		subtract := func(a, b int) int { return a - b }
		subtract10From := Bind1st(subtract, 10)

		assert.Equal(t, 7, subtract10From(3))   // 10 - 3
		assert.Equal(t, 0, subtract10From(10))  // 10 - 10
		assert.Equal(t, -5, subtract10From(15)) // 10 - 15
	})

	t.Run("binds first parameter of string concatenation", func(t *testing.T) {
		concat := func(a, b string) string { return a + b }
		addHello := Bind1st(concat, "Hello ")
		addPrefix := Bind1st(concat, "Prefix: ")

		assert.Equal(t, "Hello World", addHello("World"))
		assert.Equal(t, "Hello Go", addHello("Go"))
		assert.Equal(t, "Prefix: Test", addPrefix("Test"))
	})

	t.Run("binds first parameter with different types", func(t *testing.T) {
		repeat := func(s string, n int) string {
			return strings.Repeat(s, n)
		}
		repeatX := Bind1st(repeat, "x")
		repeatAB := Bind1st(repeat, "ab")

		assert.Equal(t, "xxx", repeatX(3))
		assert.Equal(t, "xxxxx", repeatX(5))
		assert.Equal(t, "abab", repeatAB(2))
	})

	t.Run("binds first parameter with complex types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		format := func(p Person, suffix string) string {
			return fmt.Sprintf("%s (%d) %s", p.Name, p.Age, suffix)
		}

		alice := Person{Name: "Alice", Age: 30}
		formatAlice := Bind1st(format, alice)

		assert.Equal(t, "Alice (30) is here", formatAlice("is here"))
		assert.Equal(t, "Alice (30) says hello", formatAlice("says hello"))
	})

	t.Run("binds first parameter with slice operations", func(t *testing.T) {
		appendSlice := func(slice []int, elem int) []int {
			return append(slice, elem)
		}

		nums := []int{1, 2, 3}
		appendToNums := Bind1st(appendSlice, nums)

		result1 := appendToNums(4)
		assert.Equal(t, []int{1, 2, 3, 4}, result1)

		result2 := appendToNums(5)
		assert.Equal(t, []int{1, 2, 3, 5}, result2)
	})

	t.Run("binds first parameter with map operations", func(t *testing.T) {
		getFromMap := func(m map[string]int, key string) int {
			return m[key]
		}

		data := map[string]int{"a": 1, "b": 2, "c": 3}
		getFromData := Bind1st(getFromMap, data)

		assert.Equal(t, 1, getFromData("a"))
		assert.Equal(t, 2, getFromData("b"))
		assert.Equal(t, 3, getFromData("c"))
	})

	t.Run("creates specialized comparison functions", func(t *testing.T) {
		greaterThan := func(a, b int) bool { return a > b }
		greaterThan10 := Bind1st(greaterThan, 10)
		greaterThan5 := Bind1st(greaterThan, 5)

		assert.True(t, greaterThan10(3))   // 10 > 3
		assert.False(t, greaterThan10(15)) // 10 > 15
		assert.True(t, greaterThan5(3))    // 5 > 3
		assert.False(t, greaterThan5(10))  // 5 > 10
	})
}

// TestBind2nd tests the Bind2nd function with various scenarios
func TestBind2nd(t *testing.T) {
	t.Run("binds second parameter of multiplication", func(t *testing.T) {
		multiply := func(a, b int) int { return a * b }
		double := Bind2nd(multiply, 2)
		triple := Bind2nd(multiply, 3)

		assert.Equal(t, 10, double(5))
		assert.Equal(t, 20, double(10))
		assert.Equal(t, 15, triple(5))
		assert.Equal(t, 30, triple(10))
	})

	t.Run("binds second parameter of division", func(t *testing.T) {
		divide := func(a, b float64) float64 { return a / b }
		halve := Bind2nd(divide, 2.0)
		third := Bind2nd(divide, 3.0)

		assert.Equal(t, 5.0, halve(10.0))
		assert.Equal(t, 2.5, halve(5.0))
		assert.InDelta(t, 3.333, third(10.0), 0.001)
	})

	t.Run("binds second parameter of subtraction", func(t *testing.T) {
		subtract := func(a, b int) int { return a - b }
		decrementBy5 := Bind2nd(subtract, 5)
		decrementBy10 := Bind2nd(subtract, 10)

		assert.Equal(t, 5, decrementBy5(10))  // 10 - 5
		assert.Equal(t, 0, decrementBy5(5))   // 5 - 5
		assert.Equal(t, 0, decrementBy10(10)) // 10 - 10
		assert.Equal(t, -5, decrementBy10(5)) // 5 - 10
	})

	t.Run("binds second parameter of string concatenation", func(t *testing.T) {
		concat := func(a, b string) string { return a + b }
		addWorld := Bind2nd(concat, " World")
		addSuffix := Bind2nd(concat, "!")

		assert.Equal(t, "Hello World", addWorld("Hello"))
		assert.Equal(t, "Goodbye World", addWorld("Goodbye"))
		assert.Equal(t, "Hello!", addSuffix("Hello"))
	})

	t.Run("binds second parameter with different types", func(t *testing.T) {
		repeat := func(s string, n int) string {
			return strings.Repeat(s, n)
		}
		repeatThrice := Bind2nd(repeat, 3)
		repeatTwice := Bind2nd(repeat, 2)

		assert.Equal(t, "xxx", repeatThrice("x"))
		assert.Equal(t, "ababab", repeatThrice("ab"))
		assert.Equal(t, "aa", repeatTwice("a"))
	})

	t.Run("binds second parameter with complex types", func(t *testing.T) {
		type Config struct {
			Debug bool
			Port  int
		}

		format := func(name string, cfg Config) string {
			return fmt.Sprintf("%s: debug=%v, port=%d", name, cfg.Debug, cfg.Port)
		}

		prodConfig := Config{Debug: false, Port: 8080}
		formatWithProd := Bind2nd(format, prodConfig)

		assert.Equal(t, "API: debug=false, port=8080", formatWithProd("API"))
		assert.Equal(t, "Web: debug=false, port=8080", formatWithProd("Web"))
	})

	t.Run("binds second parameter with slice operations", func(t *testing.T) {
		appendElem := func(slice []int, elem int) []int {
			return append(slice, elem)
		}

		append5 := Bind2nd(appendElem, 5)

		result1 := append5([]int{1, 2, 3})
		assert.Equal(t, []int{1, 2, 3, 5}, result1)

		result2 := append5([]int{10, 20})
		assert.Equal(t, []int{10, 20, 5}, result2)
	})

	t.Run("creates specialized comparison functions", func(t *testing.T) {
		greaterThan := func(a, b int) bool { return a > b }
		greaterThan10 := Bind2nd(greaterThan, 10)
		greaterThan5 := Bind2nd(greaterThan, 5)

		assert.False(t, greaterThan10(3)) // 3 > 10
		assert.True(t, greaterThan10(15)) // 15 > 10
		assert.False(t, greaterThan5(3))  // 3 > 5
		assert.True(t, greaterThan5(10))  // 10 > 5
	})

	t.Run("binds second parameter for power function", func(t *testing.T) {
		power := func(base, exp float64) float64 {
			result := 1.0
			for i := 0; i < int(exp); i++ {
				result *= base
			}
			return result
		}

		square := Bind2nd(power, 2.0)
		cube := Bind2nd(power, 3.0)

		assert.Equal(t, 25.0, square(5.0))
		assert.Equal(t, 100.0, square(10.0))
		assert.Equal(t, 125.0, cube(5.0))
		assert.Equal(t, 8.0, cube(2.0))
	})
}

// TestBind1stVsBind2nd tests the difference between Bind1st and Bind2nd
func TestBind1stVsBind2nd(t *testing.T) {
	t.Run("demonstrates difference with non-commutative operations", func(t *testing.T) {
		subtract := func(a, b int) int { return a - b }

		// Bind1st: fixes first parameter (a)
		subtract10From := Bind1st(subtract, 10) // 10 - b
		assert.Equal(t, 7, subtract10From(3))   // 10 - 3 = 7

		// Bind2nd: fixes second parameter (b)
		decrementBy10 := Bind2nd(subtract, 10) // a - 10
		assert.Equal(t, -7, decrementBy10(3))  // 3 - 10 = -7
	})

	t.Run("demonstrates difference with division", func(t *testing.T) {
		divide := func(a, b float64) float64 { return a / b }

		// Bind1st: fixes numerator
		divide10By := Bind1st(divide, 10.0)   // 10 / b
		assert.Equal(t, 5.0, divide10By(2.0)) // 10 / 2 = 5

		// Bind2nd: fixes denominator
		divideBy10 := Bind2nd(divide, 10.0)   // a / 10
		assert.Equal(t, 0.2, divideBy10(2.0)) // 2 / 10 = 0.2
	})

	t.Run("demonstrates equivalence with commutative operations", func(t *testing.T) {
		add := func(a, b int) int { return a + b }

		// For commutative operations, both should give same result
		add5First := Bind1st(add, 5)  // 5 + b
		add5Second := Bind2nd(add, 5) // a + 5

		assert.Equal(t, 8, add5First(3))
		assert.Equal(t, 8, add5Second(3))
		assert.Equal(t, add5First(10), add5Second(10))
	})
}

// TestSK tests the SK combinator function
func TestSK(t *testing.T) {
	t.Run("returns second argument ignoring first", func(t *testing.T) {
		assert.Equal(t, "hello", SK(42, "hello"))
		assert.Equal(t, 100, SK(true, 100))
		assert.Equal(t, 3.14, SK("test", 3.14))
		assert.Equal(t, false, SK(123, false))
	})

	t.Run("works with nil values", func(t *testing.T) {
		var nilPtr *int
		assert.Nil(t, SK("ignored", nilPtr))
		assert.Equal(t, 42, SK(nilPtr, 42))
	})

	t.Run("works with complex types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		alice := Person{Name: "Alice", Age: 30}
		bob := Person{Name: "Bob", Age: 25}

		result := SK(alice, bob)
		assert.Equal(t, "Bob", result.Name)
		assert.Equal(t, 25, result.Age)
	})

	t.Run("works with slices", func(t *testing.T) {
		slice1 := []int{1, 2, 3}
		slice2 := []string{"a", "b", "c"}

		result := SK(slice1, slice2)
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("works with maps", func(t *testing.T) {
		map1 := map[string]int{"a": 1}
		map2 := map[int]string{1: "one"}

		result := SK(map1, map2)
		assert.Equal(t, map[int]string{1: "one"}, result)
	})

	t.Run("behaves identically to Second", func(t *testing.T) {
		// SK should be identical to Second function
		testCases := []struct {
			first  any
			second any
		}{
			{42, "hello"},
			{true, 100},
			{"test", 3.14},
			{[]int{1, 2}, []string{"a", "b"}},
		}

		for _, tc := range testCases {
			assert.Equal(t,
				Second(tc.first, tc.second),
				SK(tc.first, tc.second),
				"SK should behave like Second")
		}
	})

	t.Run("demonstrates K combinator property", func(t *testing.T) {
		// SK is the K combinator applied to the second argument
		// K x y = x, so SK x y = K y x = y
		// This means SK always returns its second argument

		// Test with various types
		assert.Equal(t, 42, SK("anything", 42))
		assert.Equal(t, "result", SK(999, "result"))
		assert.True(t, SK(false, true))
	})
}

// TestBindComposition tests composition of bind operations
func TestBindComposition(t *testing.T) {
	t.Run("composes multiple Bind1st operations", func(t *testing.T) {
		add := func(a, b int) int { return a + b }
		multiply := func(a, b int) int { return a * b }

		add5 := Bind1st(add, 5)
		double := Bind1st(multiply, 2)

		// Compose: first add 5, then double
		result := double(add5(3)) // (3 + 5) * 2 = 16
		assert.Equal(t, 16, result)
	})

	t.Run("composes Bind1st and Bind2nd", func(t *testing.T) {
		subtract := func(a, b int) int { return a - b }

		subtract10From := Bind1st(subtract, 10) // 10 - b
		decrementBy5 := Bind2nd(subtract, 5)    // a - 5

		// Apply both transformations
		result1 := decrementBy5(subtract10From(3)) // (10 - 3) - 5 = 2
		assert.Equal(t, 2, result1)

		result2 := subtract10From(decrementBy5(8)) // 10 - (8 - 5) = 7
		assert.Equal(t, 7, result2)
	})

	t.Run("creates pipeline with bound functions", func(t *testing.T) {
		multiply := func(a, b int) int { return a * b }
		add := func(a, b int) int { return a + b }

		double := Bind2nd(multiply, 2)
		add10 := Bind2nd(add, 10)

		// Pipeline: input -> double -> add10
		pipeline := func(n int) int {
			return add10(double(n))
		}

		assert.Equal(t, 20, pipeline(5))  // (5 * 2) + 10 = 20
		assert.Equal(t, 30, pipeline(10)) // (10 * 2) + 10 = 30
	})
}

// TestBindWithHigherOrderFunctions tests bind with higher-order functions
func TestBindWithHigherOrderFunctions(t *testing.T) {
	t.Run("binds function parameter", func(t *testing.T) {
		applyTwice := func(f func(int) int, n int) int {
			return f(f(n))
		}

		increment := func(n int) int { return n + 1 }
		applyIncrementTwice := Bind1st(applyTwice, increment)

		assert.Equal(t, 7, applyIncrementTwice(5)) // increment(increment(5)) = 7
	})

	t.Run("binds value for higher-order function", func(t *testing.T) {
		applyFunc := func(f func(int) int, n int) int {
			return f(n)
		}

		applyTo10 := Bind2nd(applyFunc, 10)

		double := func(n int) int { return n * 2 }
		square := func(n int) int { return n * n }

		assert.Equal(t, 20, applyTo10(double))  // double(10) = 20
		assert.Equal(t, 100, applyTo10(square)) // square(10) = 100
	})
}

// BenchmarkBind1st benchmarks the Bind1st function
func BenchmarkBind1st(b *testing.B) {
	multiply := func(a, b int) int { return a * b }
	double := Bind1st(multiply, 2)

	b.Run("direct call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = multiply(2, i)
		}
	})

	b.Run("bound function", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = double(i)
		}
	})
}

// BenchmarkBind2nd benchmarks the Bind2nd function
func BenchmarkBind2nd(b *testing.B) {
	multiply := func(a, b int) int { return a * b }
	double := Bind2nd(multiply, 2)

	b.Run("direct call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = multiply(i, 2)
		}
	})

	b.Run("bound function", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = double(i)
		}
	})
}

// BenchmarkSK benchmarks the SK combinator
func BenchmarkSK(b *testing.B) {
	b.Run("SK with ints", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = SK(i, i+1)
		}
	})

	b.Run("Second with ints", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Second(i, i+1)
		}
	})
}
