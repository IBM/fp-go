//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package magma

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeMagma(t *testing.T) {
	t.Run("integer addition", func(t *testing.T) {
		addMagma := MakeMagma(func(a, b int) int {
			return a + b
		})

		assert.Equal(t, 8, addMagma.Concat(5, 3))
		assert.Equal(t, 0, addMagma.Concat(-5, 5))
		assert.Equal(t, 10, addMagma.Concat(7, 3))
	})

	t.Run("integer multiplication", func(t *testing.T) {
		mulMagma := MakeMagma(func(a, b int) int {
			return a * b
		})

		assert.Equal(t, 15, mulMagma.Concat(5, 3))
		assert.Equal(t, 0, mulMagma.Concat(0, 5))
		assert.Equal(t, 21, mulMagma.Concat(7, 3))
	})

	t.Run("string concatenation", func(t *testing.T) {
		stringMagma := MakeMagma(func(a, b string) string {
			return a + b
		})

		assert.Equal(t, "HelloWorld", stringMagma.Concat("Hello", "World"))
		assert.Equal(t, "ab", stringMagma.Concat("a", "b"))
	})

	t.Run("max operation", func(t *testing.T) {
		maxMagma := MakeMagma(func(a, b int) int {
			if a > b {
				return a
			}
			return b
		})

		assert.Equal(t, 5, maxMagma.Concat(5, 3))
		assert.Equal(t, 7, maxMagma.Concat(2, 7))
		assert.Equal(t, 10, maxMagma.Concat(10, 10))
	})
}

func TestFirst(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		m := First[string]()

		assert.Equal(t, "a", m.Concat("a", "b"))
		assert.Equal(t, "first", m.Concat("first", "second"))
		assert.Equal(t, "", m.Concat("", "something"))
	})

	t.Run("int", func(t *testing.T) {
		m := First[int]()

		assert.Equal(t, 1, m.Concat(1, 2))
		assert.Equal(t, 42, m.Concat(42, 100))
		assert.Equal(t, 0, m.Concat(0, 5))
	})
}

func TestSecond(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		m := Second[string]()

		assert.Equal(t, "b", m.Concat("a", "b"))
		assert.Equal(t, "second", m.Concat("first", "second"))
		assert.Equal(t, "something", m.Concat("", "something"))
	})

	t.Run("int", func(t *testing.T) {
		m := Second[int]()

		assert.Equal(t, 2, m.Concat(1, 2))
		assert.Equal(t, 100, m.Concat(42, 100))
		assert.Equal(t, 5, m.Concat(0, 5))
	})
}

func TestReverse(t *testing.T) {
	t.Run("subtraction", func(t *testing.T) {
		subMagma := MakeMagma(func(a, b int) int {
			return a - b
		})

		reversedMagma := Reverse(subMagma)

		assert.Equal(t, 7, subMagma.Concat(10, 3))       // 10 - 3
		assert.Equal(t, -7, reversedMagma.Concat(10, 3)) // 3 - 10
	})

	t.Run("division", func(t *testing.T) {
		divMagma := MakeMagma(func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		})

		reversedMagma := Reverse(divMagma)

		assert.Equal(t, 5, divMagma.Concat(10, 2))      // 10 / 2
		assert.Equal(t, 0, reversedMagma.Concat(10, 2)) // 2 / 10
	})

	t.Run("string concatenation", func(t *testing.T) {
		stringMagma := MakeMagma(func(a, b string) string {
			return a + b
		})

		reversedMagma := Reverse(stringMagma)

		assert.Equal(t, "ab", stringMagma.Concat("a", "b"))
		assert.Equal(t, "ba", reversedMagma.Concat("a", "b"))
	})
}

func TestFilterFirst(t *testing.T) {
	addMagma := MakeMagma(func(a, b int) int {
		return a + b
	})

	t.Run("positive first", func(t *testing.T) {
		filteredMagma := FilterFirst(func(n int) bool {
			return n > 0
		})(addMagma)

		assert.Equal(t, 8, filteredMagma.Concat(5, 3))  // 5 is positive: 5 + 3
		assert.Equal(t, 3, filteredMagma.Concat(-5, 3)) // -5 is negative: return 3
		assert.Equal(t, 10, filteredMagma.Concat(7, 3)) // 7 is positive: 7 + 3
	})

	t.Run("even first", func(t *testing.T) {
		filteredMagma := FilterFirst(func(n int) bool {
			return n%2 == 0
		})(addMagma)

		assert.Equal(t, 7, filteredMagma.Concat(4, 3)) // 4 is even: 4 + 3
		assert.Equal(t, 3, filteredMagma.Concat(5, 3)) // 5 is odd: return 3
		assert.Equal(t, 5, filteredMagma.Concat(2, 3)) // 2 is even: 2 + 3
	})
}

func TestFilterSecond(t *testing.T) {
	addMagma := MakeMagma(func(a, b int) int {
		return a + b
	})

	t.Run("positive second", func(t *testing.T) {
		filteredMagma := FilterSecond(func(n int) bool {
			return n > 0
		})(addMagma)

		assert.Equal(t, 8, filteredMagma.Concat(5, 3))  // 3 is positive: 5 + 3
		assert.Equal(t, 5, filteredMagma.Concat(5, -3)) // -3 is negative: return 5
		assert.Equal(t, 10, filteredMagma.Concat(7, 3)) // 3 is positive: 7 + 3
	})

	t.Run("even second", func(t *testing.T) {
		filteredMagma := FilterSecond(func(n int) bool {
			return n%2 == 0
		})(addMagma)

		assert.Equal(t, 9, filteredMagma.Concat(5, 4)) // 4 is even: 5 + 4
		assert.Equal(t, 5, filteredMagma.Concat(5, 3)) // 3 is odd: return 5
		assert.Equal(t, 7, filteredMagma.Concat(5, 2)) // 2 is even: 5 + 2
	})
}

func TestEndo(t *testing.T) {
	addMagma := MakeMagma(func(a, b int) int {
		return a + b
	})

	t.Run("double before adding", func(t *testing.T) {
		doubledMagma := Endo(func(n int) int {
			return n * 2
		})(addMagma)

		assert.Equal(t, 14, doubledMagma.Concat(3, 4)) // (3*2) + (4*2) = 6 + 8
		assert.Equal(t, 10, doubledMagma.Concat(2, 3)) // (2*2) + (3*2) = 4 + 6
		assert.Equal(t, 0, doubledMagma.Concat(0, 0))  // (0*2) + (0*2) = 0
	})

	t.Run("square before adding", func(t *testing.T) {
		squaredMagma := Endo(func(n int) int {
			return n * n
		})(addMagma)

		assert.Equal(t, 25, squaredMagma.Concat(3, 4)) // (3*3) + (4*4) = 9 + 16
		assert.Equal(t, 13, squaredMagma.Concat(2, 3)) // (2*2) + (3*3) = 4 + 9
		assert.Equal(t, 2, squaredMagma.Concat(1, 1))  // (1*1) + (1*1) = 2
	})

	t.Run("negate before adding", func(t *testing.T) {
		negatedMagma := Endo(func(n int) int {
			return -n
		})(addMagma)

		assert.Equal(t, -7, negatedMagma.Concat(3, 4)) // (-3) + (-4) = -7
		assert.Equal(t, -5, negatedMagma.Concat(2, 3)) // (-2) + (-3) = -5
		assert.Equal(t, 0, negatedMagma.Concat(0, 0))  // 0 + 0 = 0
	})
}

func TestConcatAll(t *testing.T) {
	t.Run("sum integers", func(t *testing.T) {
		addMagma := MakeMagma(func(a, b int) int {
			return a + b
		})

		numbers := []int{1, 2, 3, 4, 5}
		result := ConcatAll(addMagma)(0)(numbers)
		assert.Equal(t, 15, result)
	})

	t.Run("multiply integers", func(t *testing.T) {
		mulMagma := MakeMagma(func(a, b int) int {
			return a * b
		})

		numbers := []int{2, 3, 4}
		result := ConcatAll(mulMagma)(1)(numbers)
		assert.Equal(t, 24, result)
	})

	t.Run("max of integers", func(t *testing.T) {
		maxMagma := MakeMagma(func(a, b int) int {
			if a > b {
				return a
			}
			return b
		})

		numbers := []int{3, 7, 2, 9, 1}
		result := ConcatAll(maxMagma)(0)(numbers)
		assert.Equal(t, 9, result)
	})

	t.Run("concatenate strings", func(t *testing.T) {
		stringMagma := MakeMagma(func(a, b string) string {
			return a + b
		})

		words := []string{"Hello", " ", "World"}
		result := ConcatAll(stringMagma)("")(words)
		assert.Equal(t, "Hello World", result)
	})

	t.Run("empty slice", func(t *testing.T) {
		addMagma := MakeMagma(func(a, b int) int {
			return a + b
		})

		numbers := []int{}
		result := ConcatAll(addMagma)(42)(numbers)
		assert.Equal(t, 42, result)
	})
}

func TestMonadConcatAll(t *testing.T) {
	t.Run("sum integers", func(t *testing.T) {
		addMagma := MakeMagma(func(a, b int) int {
			return a + b
		})

		numbers := []int{1, 2, 3, 4, 5}
		result := MonadConcatAll(addMagma)(numbers, 0)
		assert.Equal(t, 15, result)
	})

	t.Run("multiply integers", func(t *testing.T) {
		mulMagma := MakeMagma(func(a, b int) int {
			return a * b
		})

		numbers := []int{2, 3, 4}
		result := MonadConcatAll(mulMagma)(numbers, 1)
		assert.Equal(t, 24, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		addMagma := MakeMagma(func(a, b int) int {
			return a + b
		})

		numbers := []int{}
		result := MonadConcatAll(addMagma)(numbers, 100)
		assert.Equal(t, 100, result)
	})
}

func TestGenericConcatAll(t *testing.T) {
	type IntSlice []int

	t.Run("custom slice type", func(t *testing.T) {
		addMagma := MakeMagma(func(a, b int) int {
			return a + b
		})

		numbers := IntSlice{1, 2, 3, 4, 5}
		result := GenericConcatAll[IntSlice](addMagma)(0)(numbers)
		assert.Equal(t, 15, result)
	})

	t.Run("regular slice", func(t *testing.T) {
		mulMagma := MakeMagma(func(a, b int) int {
			return a * b
		})

		numbers := []int{2, 3, 4}
		result := GenericConcatAll[[]int](mulMagma)(1)(numbers)
		assert.Equal(t, 24, result)
	})
}

func TestGenericMonadConcatAll(t *testing.T) {
	type IntSlice []int

	t.Run("custom slice type", func(t *testing.T) {
		addMagma := MakeMagma(func(a, b int) int {
			return a + b
		})

		numbers := IntSlice{1, 2, 3, 4, 5}
		result := GenericMonadConcatAll[IntSlice](addMagma)(numbers, 0)
		assert.Equal(t, 15, result)
	})

	t.Run("regular slice", func(t *testing.T) {
		mulMagma := MakeMagma(func(a, b int) int {
			return a * b
		})

		numbers := []int{2, 3, 4}
		result := GenericMonadConcatAll[[]int](mulMagma)(numbers, 1)
		assert.Equal(t, 24, result)
	})
}

// Test practical examples
func TestPracticalExamples(t *testing.T) {
	t.Run("min magma", func(t *testing.T) {
		minMagma := MakeMagma(func(a, b int) int {
			if a < b {
				return a
			}
			return b
		})

		numbers := []int{3, 7, 2, 9, 1}
		minimum := ConcatAll(minMagma)(100)(numbers)
		assert.Equal(t, 1, minimum)
	})

	t.Run("string join with separator", func(t *testing.T) {
		joinMagma := MakeMagma(func(a, b string) string {
			if a == "" {
				return b
			}
			if b == "" {
				return a
			}
			return a + ", " + b
		})

		words := []string{"apple", "banana", "cherry"}
		result := ConcatAll(joinMagma)("")(words)
		assert.Equal(t, "apple, banana, cherry", result)
	})

	t.Run("boolean AND", func(t *testing.T) {
		andMagma := MakeMagma(func(a, b bool) bool {
			return a && b
		})

		values := []bool{true, true, true}
		result := ConcatAll(andMagma)(true)(values)
		assert.True(t, result)

		values2 := []bool{true, false, true}
		result2 := ConcatAll(andMagma)(true)(values2)
		assert.False(t, result2)
	})

	t.Run("boolean OR", func(t *testing.T) {
		orMagma := MakeMagma(func(a, b bool) bool {
			return a || b
		})

		values := []bool{false, false, false}
		result := ConcatAll(orMagma)(false)(values)
		assert.False(t, result)

		values2 := []bool{false, true, false}
		result2 := ConcatAll(orMagma)(false)(values2)
		assert.True(t, result2)
	})
}
