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
	"testing"

	S "github.com/IBM/fp-go/v2/semigroup"
	"github.com/stretchr/testify/assert"
)

// TestFromSemigroup tests the FromSemigroup function with various semigroups
func TestFromSemigroup(t *testing.T) {
	t.Run("integer addition semigroup", func(t *testing.T) {
		// Create a semigroup for integer addition
		addSemigroup := S.MakeSemigroup(func(a, b int) int {
			return a + b
		})

		// Convert to Kleisli arrow
		addKleisli := FromSemigroup(addSemigroup)

		// Create an endomorphism that adds 5
		addFive := addKleisli(5)

		// Test the endomorphism
		assert.Equal(t, 15, addFive(10), "addFive(10) should equal 15")
		assert.Equal(t, 5, addFive(0), "addFive(0) should equal 5")
		assert.Equal(t, -5, addFive(-10), "addFive(-10) should equal -5")
	})

	t.Run("integer multiplication semigroup", func(t *testing.T) {
		// Create a semigroup for integer multiplication
		mulSemigroup := S.MakeSemigroup(func(a, b int) int {
			return a * b
		})

		// Convert to Kleisli arrow
		mulKleisli := FromSemigroup(mulSemigroup)

		// Create an endomorphism that multiplies by 3
		multiplyByThree := mulKleisli(3)

		// Test the endomorphism
		assert.Equal(t, 15, multiplyByThree(5), "multiplyByThree(5) should equal 15")
		assert.Equal(t, 0, multiplyByThree(0), "multiplyByThree(0) should equal 0")
		assert.Equal(t, -9, multiplyByThree(-3), "multiplyByThree(-3) should equal -9")
	})

	t.Run("string concatenation semigroup", func(t *testing.T) {
		// Create a semigroup for string concatenation
		concatSemigroup := S.MakeSemigroup(func(a, b string) string {
			return a + b
		})

		// Convert to Kleisli arrow
		concatKleisli := FromSemigroup(concatSemigroup)

		// Create an endomorphism that appends "Hello, " (input is on the left)
		appendHello := concatKleisli("Hello, ")

		// Test the endomorphism - input is concatenated on the left, "Hello, " on the right
		assert.Equal(t, "WorldHello, ", appendHello("World"), "appendHello('World') should equal 'WorldHello, '")
		assert.Equal(t, "Hello, ", appendHello(""), "appendHello('') should equal 'Hello, '")
		assert.Equal(t, "GoHello, ", appendHello("Go"), "appendHello('Go') should equal 'GoHello, '")
	})

	t.Run("slice concatenation semigroup", func(t *testing.T) {
		// Create a semigroup for slice concatenation
		sliceSemigroup := S.MakeSemigroup(func(a, b []int) []int {
			result := make([]int, len(a)+len(b))
			copy(result, a)
			copy(result[len(a):], b)
			return result
		})

		// Convert to Kleisli arrow
		sliceKleisli := FromSemigroup(sliceSemigroup)

		// Create an endomorphism that appends [1, 2] (input is on the left)
		appendOneTwo := sliceKleisli([]int{1, 2})

		// Test the endomorphism - input is concatenated on the left, [1,2] on the right
		result1 := appendOneTwo([]int{3, 4, 5})
		assert.Equal(t, []int{3, 4, 5, 1, 2}, result1, "appendOneTwo([3,4,5]) should equal [3,4,5,1,2]")

		result2 := appendOneTwo([]int{})
		assert.Equal(t, []int{1, 2}, result2, "appendOneTwo([]) should equal [1,2]")

		result3 := appendOneTwo([]int{10})
		assert.Equal(t, []int{10, 1, 2}, result3, "appendOneTwo([10]) should equal [10,1,2]")
	})

	t.Run("max semigroup", func(t *testing.T) {
		// Create a semigroup for max operation
		maxSemigroup := S.MakeSemigroup(func(a, b int) int {
			if a > b {
				return a
			}
			return b
		})

		// Convert to Kleisli arrow
		maxKleisli := FromSemigroup(maxSemigroup)

		// Create an endomorphism that takes max with 10
		maxWithTen := maxKleisli(10)

		// Test the endomorphism
		assert.Equal(t, 15, maxWithTen(15), "maxWithTen(15) should equal 15")
		assert.Equal(t, 10, maxWithTen(5), "maxWithTen(5) should equal 10")
		assert.Equal(t, 10, maxWithTen(10), "maxWithTen(10) should equal 10")
		assert.Equal(t, 10, maxWithTen(-5), "maxWithTen(-5) should equal 10")
	})

	t.Run("min semigroup", func(t *testing.T) {
		// Create a semigroup for min operation
		minSemigroup := S.MakeSemigroup(func(a, b int) int {
			if a < b {
				return a
			}
			return b
		})

		// Convert to Kleisli arrow
		minKleisli := FromSemigroup(minSemigroup)

		// Create an endomorphism that takes min with 10
		minWithTen := minKleisli(10)

		// Test the endomorphism
		assert.Equal(t, 5, minWithTen(5), "minWithTen(5) should equal 5")
		assert.Equal(t, 10, minWithTen(15), "minWithTen(15) should equal 10")
		assert.Equal(t, 10, minWithTen(10), "minWithTen(10) should equal 10")
		assert.Equal(t, -5, minWithTen(-5), "minWithTen(-5) should equal -5")
	})
}

// TestFromSemigroupComposition tests that endomorphisms created from semigroups can be composed
func TestFromSemigroupComposition(t *testing.T) {
	t.Run("compose addition endomorphisms", func(t *testing.T) {
		// Create a semigroup for integer addition
		addSemigroup := S.MakeSemigroup(func(a, b int) int {
			return a + b
		})
		addKleisli := FromSemigroup(addSemigroup)

		// Create two endomorphisms
		addFive := addKleisli(5)
		addTen := addKleisli(10)

		// Compose them (RIGHT-TO-LEFT execution)
		composed := MonadCompose(addFive, addTen)

		// Test composition: addTen first, then addFive
		result := composed(3) // 3 + 10 = 13, then 13 + 5 = 18
		assert.Equal(t, 18, result, "composed addition should work correctly")
	})

	t.Run("compose string endomorphisms", func(t *testing.T) {
		// Create a semigroup for string concatenation
		concatSemigroup := S.MakeSemigroup(func(a, b string) string {
			return a + b
		})
		concatKleisli := FromSemigroup(concatSemigroup)

		// Create two endomorphisms
		appendHello := concatKleisli("Hello, ")
		appendExclamation := concatKleisli("!")

		// Compose them (RIGHT-TO-LEFT execution)
		composed := MonadCompose(appendHello, appendExclamation)

		// Test composition: appendExclamation first, then appendHello
		// "World" + "!" = "World!", then "World!" + "Hello, " = "World!Hello, "
		result := composed("World")
		assert.Equal(t, "World!Hello, ", result, "composed string operations should work correctly")
	})
}

// TestFromSemigroupWithMonoid tests using FromSemigroup-created endomorphisms with monoid operations
func TestFromSemigroupWithMonoid(t *testing.T) {
	t.Run("monoid concat with addition endomorphisms", func(t *testing.T) {
		// Create a semigroup for integer addition
		addSemigroup := S.MakeSemigroup(func(a, b int) int {
			return a + b
		})
		addKleisli := FromSemigroup(addSemigroup)

		// Create multiple endomorphisms
		addOne := addKleisli(1)
		addTwo := addKleisli(2)
		addThree := addKleisli(3)

		// Use monoid to combine them
		monoid := Monoid[int]()
		combined := monoid.Concat(monoid.Concat(addOne, addTwo), addThree)

		// Test: RIGHT-TO-LEFT execution: addThree, then addTwo, then addOne
		result := combined(10) // 10 + 3 = 13, 13 + 2 = 15, 15 + 1 = 16
		assert.Equal(t, 16, result, "monoid combination should work correctly")
	})
}

// TestFromSemigroupAssociativity tests that the semigroup associativity is preserved
func TestFromSemigroupAssociativity(t *testing.T) {
	t.Run("addition associativity", func(t *testing.T) {
		// Create a semigroup for integer addition
		addSemigroup := S.MakeSemigroup(func(a, b int) int {
			return a + b
		})
		addKleisli := FromSemigroup(addSemigroup)

		// Create three endomorphisms
		addTwo := addKleisli(2)
		addThree := addKleisli(3)
		addFive := addKleisli(5)

		// Test associativity: (a . b) . c = a . (b . c)
		left := MonadCompose(MonadCompose(addTwo, addThree), addFive)
		right := MonadCompose(addTwo, MonadCompose(addThree, addFive))

		testValue := 10
		assert.Equal(t, left(testValue), right(testValue), "composition should be associative")

		// Both should equal: 10 + 5 + 3 + 2 = 20
		assert.Equal(t, 20, left(testValue), "left composition should equal 20")
		assert.Equal(t, 20, right(testValue), "right composition should equal 20")
	})

	t.Run("string concatenation associativity", func(t *testing.T) {
		// Create a semigroup for string concatenation
		concatSemigroup := S.MakeSemigroup(func(a, b string) string {
			return a + b
		})
		concatKleisli := FromSemigroup(concatSemigroup)

		// Create three endomorphisms
		appendA := concatKleisli("A")
		appendB := concatKleisli("B")
		appendC := concatKleisli("C")

		// Test associativity: (a . b) . c = a . (b . c)
		left := MonadCompose(MonadCompose(appendA, appendB), appendC)
		right := MonadCompose(appendA, MonadCompose(appendB, appendC))

		testValue := "X"
		assert.Equal(t, left(testValue), right(testValue), "string composition should be associative")

		// Both should equal: "X" + "C" + "B" + "A" = "XCBA" (RIGHT-TO-LEFT composition)
		assert.Equal(t, "XCBA", left(testValue), "left composition should equal 'XCBA'")
		assert.Equal(t, "XCBA", right(testValue), "right composition should equal 'XCBA'")
	})
}

// TestFromSemigroupEdgeCases tests edge cases and boundary conditions
func TestFromSemigroupEdgeCases(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		// Test with addition and zero
		addSemigroup := S.MakeSemigroup(func(a, b int) int {
			return a + b
		})
		addKleisli := FromSemigroup(addSemigroup)

		addZero := addKleisli(0)
		assert.Equal(t, 5, addZero(5), "adding zero should not change the value")
		assert.Equal(t, 0, addZero(0), "adding zero to zero should be zero")
		assert.Equal(t, -3, addZero(-3), "adding zero to negative should not change")
	})

	t.Run("empty string", func(t *testing.T) {
		// Test with string concatenation and empty string
		concatSemigroup := S.MakeSemigroup(func(a, b string) string {
			return a + b
		})
		concatKleisli := FromSemigroup(concatSemigroup)

		prependEmpty := concatKleisli("")
		assert.Equal(t, "hello", prependEmpty("hello"), "prepending empty string should not change")
		assert.Equal(t, "", prependEmpty(""), "prepending empty to empty should be empty")
	})

	t.Run("empty slice", func(t *testing.T) {
		// Test with slice concatenation and empty slice
		sliceSemigroup := S.MakeSemigroup(func(a, b []int) []int {
			result := make([]int, len(a)+len(b))
			copy(result, a)
			copy(result[len(a):], b)
			return result
		})
		sliceKleisli := FromSemigroup(sliceSemigroup)

		prependEmpty := sliceKleisli([]int{})
		result := prependEmpty([]int{1, 2, 3})
		assert.Equal(t, []int{1, 2, 3}, result, "prepending empty slice should not change")

		emptyResult := prependEmpty([]int{})
		assert.Equal(t, []int{}, emptyResult, "prepending empty to empty should be empty")
	})
}

// TestFromSemigroupDataLastPrinciple explicitly tests that FromSemigroup follows the "data last" principle
func TestFromSemigroupDataLastPrinciple(t *testing.T) {
	t.Run("data last with string concatenation", func(t *testing.T) {
		// Create a semigroup for string concatenation
		// Concat(a, b) = a + b
		concatSemigroup := S.MakeSemigroup(func(a, b string) string {
			return a + b
		})

		// FromSemigroup uses Bind2of2, which binds the second parameter
		// So FromSemigroup(s)(x) creates: func(input) = Concat(input, x)
		// This is "data last" - the input data comes first, bound value comes last
		kleisli := FromSemigroup(concatSemigroup)

		// Bind "World" as the second parameter
		appendWorld := kleisli("World")

		// When we call appendWorld("Hello"), it computes Concat("Hello", "World")
		// The input "Hello" is the first parameter (data), "World" is the second (bound value)
		result := appendWorld("Hello")
		assert.Equal(t, "HelloWorld", result, "Data last: Concat(input='Hello', bound='World') = 'HelloWorld'")

		// Verify with different input
		result2 := appendWorld("Goodbye")
		assert.Equal(t, "GoodbyeWorld", result2, "Data last: Concat(input='Goodbye', bound='World') = 'GoodbyeWorld'")
	})

	t.Run("data last with integer addition", func(t *testing.T) {
		// Create a semigroup for integer addition
		// Concat(a, b) = a + b
		addSemigroup := S.MakeSemigroup(func(a, b int) int {
			return a + b
		})

		// FromSemigroup binds the second parameter
		// So FromSemigroup(s)(5) creates: func(input) = Concat(input, 5) = input + 5
		kleisli := FromSemigroup(addSemigroup)

		// Bind 5 as the second parameter
		addFive := kleisli(5)

		// When we call addFive(10), it computes Concat(10, 5) = 10 + 5 = 15
		// The input 10 is the first parameter (data), 5 is the second (bound value)
		result := addFive(10)
		assert.Equal(t, 15, result, "Data last: Concat(input=10, bound=5) = 15")
	})

	t.Run("data last with non-commutative operation", func(t *testing.T) {
		// Create a semigroup for a non-commutative operation to clearly show order
		// Concat(a, b) = a - b (subtraction is not commutative)
		subSemigroup := S.MakeSemigroup(func(a, b int) int {
			return a - b
		})

		// FromSemigroup binds the second parameter
		// So FromSemigroup(s)(5) creates: func(input) = Concat(input, 5) = input - 5
		kleisli := FromSemigroup(subSemigroup)

		// Bind 5 as the second parameter
		subtractFive := kleisli(5)

		// When we call subtractFive(10), it computes Concat(10, 5) = 10 - 5 = 5
		// The input 10 is the first parameter (data), 5 is the second (bound value)
		result := subtractFive(10)
		assert.Equal(t, 5, result, "Data last: Concat(input=10, bound=5) = 10 - 5 = 5")

		// If it were "data first" (binding first parameter), we would get:
		// Concat(5, 10) = 5 - 10 = -5, which is NOT what we get
		assert.NotEqual(t, -5, result, "Not data first: result is NOT Concat(bound=5, input=10) = 5 - 10 = -5")
	})

	t.Run("data last with list concatenation", func(t *testing.T) {
		// Create a semigroup for list concatenation
		// Concat(a, b) = a ++ b
		listSemigroup := S.MakeSemigroup(func(a, b []int) []int {
			result := make([]int, len(a)+len(b))
			copy(result, a)
			copy(result[len(a):], b)
			return result
		})

		// FromSemigroup binds the second parameter
		// So FromSemigroup(s)([3,4]) creates: func(input) = Concat(input, [3,4])
		kleisli := FromSemigroup(listSemigroup)

		// Bind [3, 4] as the second parameter
		appendThreeFour := kleisli([]int{3, 4})

		// When we call appendThreeFour([1,2]), it computes Concat([1,2], [3,4]) = [1,2,3,4]
		// The input [1,2] is the first parameter (data), [3,4] is the second (bound value)
		result := appendThreeFour([]int{1, 2})
		assert.Equal(t, []int{1, 2, 3, 4}, result, "Data last: Concat(input=[1,2], bound=[3,4]) = [1,2,3,4]")
	})
}

// BenchmarkFromSemigroup benchmarks the FromSemigroup function
func BenchmarkFromSemigroup(b *testing.B) {
	addSemigroup := S.MakeSemigroup(func(a, b int) int {
		return a + b
	})
	addKleisli := FromSemigroup(addSemigroup)
	addFive := addKleisli(5)

	b.ResetTimer()
	for b.Loop() {
		_ = addFive(10)
	}
}

// BenchmarkFromSemigroupComposition benchmarks composed endomorphisms from semigroups
func BenchmarkFromSemigroupComposition(b *testing.B) {
	addSemigroup := S.MakeSemigroup(func(a, b int) int {
		return a + b
	})
	addKleisli := FromSemigroup(addSemigroup)

	addFive := addKleisli(5)
	addTen := addKleisli(10)
	composed := MonadCompose(addFive, addTen)

	b.ResetTimer()
	for b.Loop() {
		_ = composed(3)
	}
}
