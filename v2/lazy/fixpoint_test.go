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

package lazy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ExampleFixpoint demonstrates computing a factorial using fixpoint.
func ExampleFixpoint() {
	// Define factorial using fixpoint
	factorial := Fixpoint(func(self Lazy[func(int) int]) func(int) int {
		return func(x int) int {
			if x <= 1 {
				return 1
			}
			return x * self()(x-1)
		}
	})

	result := factorial(5)
	fmt.Println(result)

	// Output:
	// 120
}

// ExampleFixpoint_fibonacci demonstrates computing Fibonacci numbers using fixpoint.
func ExampleFixpoint_fibonacci() {
	// Define fibonacci using fixpoint
	fib := Fixpoint(func(self Lazy[func(int) int]) func(int) int {
		return func(n int) int {
			if n <= 1 {
				return n
			}
			return self()(n-1) + self()(n-2)
		}
	})

	result := fib(10)
	fmt.Println(result)

	// Output:
	// 55
}

// ExampleFixpoint_infiniteList demonstrates creating a list using fixpoint.
func ExampleFixpoint_infiniteList() {
	// Define a function that generates a list from n down to 1
	countdown := Fixpoint(func(self Lazy[func(int) []int]) func(int) []int {
		return func(n int) []int {
			if n <= 0 {
				return []int{}
			}
			return append([]int{n}, self()(n-1)...)
		}
	})

	// Generate countdown from 5
	result := countdown(5)
	fmt.Println(result)

	// Output:
	// [5 4 3 2 1]
}

// ExampleFixpoint_constantValue demonstrates computing a constant value using fixpoint.
func ExampleFixpoint_constantValue() {
	// A simple fixpoint that returns a constant
	result := Fixpoint(func(self Lazy[int]) int {
		return 42
	})

	fmt.Println(result)

	// Output:
	// 42
}

func TestFixpoint_Factorial(t *testing.T) {
	t.Run("computes factorial correctly", func(t *testing.T) {
		factorial := Fixpoint(func(self Lazy[func(int) int]) func(int) int {
			return func(x int) int {
				if x <= 1 {
					return 1
				}
				return x * self()(x-1)
			}
		})

		tests := []struct {
			input    int
			expected int
		}{
			{0, 1},
			{1, 1},
			{2, 2},
			{3, 6},
			{4, 24},
			{5, 120},
			{6, 720},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("factorial(%d)", tt.input), func(t *testing.T) {
				result := factorial(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestFixpoint_Fibonacci(t *testing.T) {
	t.Run("computes fibonacci numbers correctly", func(t *testing.T) {
		fib := Fixpoint(func(self Lazy[func(int) int]) func(int) int {
			return func(n int) int {
				if n <= 1 {
					return n
				}
				return self()(n-1) + self()(n-2)
			}
		})

		tests := []struct {
			input    int
			expected int
		}{
			{0, 0},
			{1, 1},
			{2, 1},
			{3, 2},
			{4, 3},
			{5, 5},
			{6, 8},
			{7, 13},
			{8, 21},
			{9, 34},
			{10, 55},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("fib(%d)", tt.input), func(t *testing.T) {
				result := fib(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestFixpoint_ConstantValue(t *testing.T) {
	t.Run("returns constant value", func(t *testing.T) {
		result := Fixpoint(func(self Lazy[int]) int {
			return 42
		})

		assert.Equal(t, 42, result)
	})

	t.Run("returns constant string", func(t *testing.T) {
		result := Fixpoint(func(self Lazy[string]) string {
			return "hello"
		})

		assert.Equal(t, "hello", result)
	})
}

func TestFixpoint_SelfReferential(t *testing.T) {
	t.Run("computes value using self reference", func(t *testing.T) {
		// A function that uses the lazy self reference conditionally
		result := Fixpoint(func(self Lazy[int]) int {
			// This demonstrates that self() would give us the result
			// but we don't call it in this case
			return 100
		})

		assert.Equal(t, 100, result)
	})

	t.Run("computes sum using self reference", func(t *testing.T) {
		// Sum from n down to 0
		sumFrom := Fixpoint(func(self Lazy[func(int) int]) func(int) int {
			return func(x int) int {
				if x <= 0 {
					return 0
				}
				return x + self()(x-1)
			}
		})

		assert.Equal(t, 0, sumFrom(0))
		assert.Equal(t, 1, sumFrom(1))
		assert.Equal(t, 3, sumFrom(2))
		assert.Equal(t, 6, sumFrom(3))
		assert.Equal(t, 10, sumFrom(4))
		assert.Equal(t, 15, sumFrom(5))
	})
}

func TestFixpoint_ComplexTypes(t *testing.T) {
	t.Run("works with slice types", func(t *testing.T) {
		// Generate a list of n elements
		generateList := Fixpoint(func(self Lazy[func(int) []int]) func(int) []int {
			return func(count int) []int {
				if count <= 0 {
					return []int{}
				}
				return append([]int{count}, self()(count-1)...)
			}
		})

		result := generateList(5)
		assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
	})

	t.Run("works with struct types", func(t *testing.T) {
		type Node struct {
			Value int
			Count int
		}

		// Create a node with accumulated count
		createNode := func(value int) Node {
			return Fixpoint(func(self Lazy[func(int, int) Node]) func(int, int) Node {
				return func(v, depth int) Node {
					if depth <= 0 {
						return Node{Value: v, Count: 1}
					}
					prev := self()(v, depth-1)
					return Node{Value: v, Count: prev.Count + 1}
				}
			})(value, 3)
		}

		result := createNode(42)
		assert.Equal(t, Node{Value: 42, Count: 4}, result)
	})
}

func TestFixpoint_EdgeCases(t *testing.T) {
	t.Run("handles zero value types", func(t *testing.T) {
		result := Fixpoint(func(self Lazy[int]) int {
			return 0
		})

		assert.Equal(t, 0, result)
	})

	t.Run("handles empty string", func(t *testing.T) {
		result := Fixpoint(func(self Lazy[string]) string {
			return ""
		})

		assert.Equal(t, "", result)
	})

	t.Run("handles nil slice", func(t *testing.T) {
		result := Fixpoint(func(self Lazy[[]int]) []int {
			return nil
		})

		assert.Nil(t, result)
	})

	t.Run("handles boolean values", func(t *testing.T) {
		resultTrue := Fixpoint(func(self Lazy[bool]) bool {
			return true
		})
		resultFalse := Fixpoint(func(self Lazy[bool]) bool {
			return false
		})

		assert.True(t, resultTrue)
		assert.False(t, resultFalse)
	})
}

func TestFixpoint_RecursiveDepth(t *testing.T) {
	t.Run("handles deep recursion", func(t *testing.T) {
		// Count down from n to 0
		countdown := Fixpoint(func(self Lazy[func(int) int]) func(int) int {
			return func(x int) int {
				if x <= 0 {
					return 0
				}
				return self()(x - 1)
			}
		})

		// Test with a reasonably deep recursion
		result := countdown(100)
		assert.Equal(t, 0, result)
	})

	t.Run("accumulates values through recursion", func(t *testing.T) {
		// Product from 1 to n
		product := Fixpoint(func(self Lazy[func(int) int]) func(int) int {
			return func(x int) int {
				if x <= 1 {
					return 1
				}
				return x * self()(x-1)
			}
		})

		assert.Equal(t, 1, product(1))
		assert.Equal(t, 2, product(2))
		assert.Equal(t, 6, product(3))
		assert.Equal(t, 24, product(4))
		assert.Equal(t, 120, product(5))
	})
}

func TestFixpoint_MultipleParameters(t *testing.T) {
	t.Run("works with functions taking multiple parameters", func(t *testing.T) {
		// Ackermann function using fixpoint
		ackermann := Fixpoint(func(self Lazy[func(int, int) int]) func(int, int) int {
			return func(m, n int) int {
				if m == 0 {
					return n + 1
				}
				if n == 0 {
					return self()(m-1, 1)
				}
				return self()(m-1, self()(m, n-1))
			}
		})

		// Test small values (Ackermann grows very quickly)
		assert.Equal(t, 1, ackermann(0, 0))
		assert.Equal(t, 2, ackermann(0, 1))
		assert.Equal(t, 3, ackermann(1, 1))
		assert.Equal(t, 5, ackermann(2, 1))
		assert.Equal(t, 7, ackermann(2, 2))
	})
}

// Made with Bob
