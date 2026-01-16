// Copyright (c) 2025 IBM Corp.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIdentity tests the Identity function
func TestIdentity(t *testing.T) {
	t.Run("returns int unchanged", func(t *testing.T) {
		assert.Equal(t, 42, Identity(42))
		assert.Equal(t, 0, Identity(0))
		assert.Equal(t, -10, Identity(-10))
	})

	t.Run("returns string unchanged", func(t *testing.T) {
		assert.Equal(t, "hello", Identity("hello"))
		assert.Equal(t, "", Identity(""))
	})

	t.Run("returns bool unchanged", func(t *testing.T) {
		assert.True(t, Identity(true))
		assert.False(t, Identity(false))
	})

	t.Run("returns struct unchanged", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		p := Person{Name: "Alice", Age: 30}
		assert.Equal(t, p, Identity(p))
	})
}

// TestConstant tests the Constant function
func TestConstant(t *testing.T) {
	t.Run("returns constant int", func(t *testing.T) {
		getFortyTwo := Constant(42)
		assert.Equal(t, 42, getFortyTwo())
		assert.Equal(t, 42, getFortyTwo())
	})

	t.Run("returns constant string", func(t *testing.T) {
		getMessage := Constant("Hello")
		assert.Equal(t, "Hello", getMessage())
	})

	t.Run("returns constant bool", func(t *testing.T) {
		getTrue := Constant(true)
		assert.True(t, getTrue())
	})
}

// TestConstant1 tests the Constant1 function
func TestConstant1(t *testing.T) {
	t.Run("ignores input and returns constant", func(t *testing.T) {
		alwaysZero := Constant1[string](0)
		assert.Equal(t, 0, alwaysZero("anything"))
		assert.Equal(t, 0, alwaysZero("something else"))
		assert.Equal(t, 0, alwaysZero(""))
	})

	t.Run("works with different types", func(t *testing.T) {
		defaultName := Constant1[int]("Unknown")
		assert.Equal(t, "Unknown", defaultName(42))
		assert.Equal(t, "Unknown", defaultName(0))
	})
}

// TestConstant2 tests the Constant2 function
func TestConstant2(t *testing.T) {
	t.Run("ignores both inputs and returns constant", func(t *testing.T) {
		alwaysTrue := Constant2[int, string](true)
		assert.True(t, alwaysTrue(42, "test"))
		assert.True(t, alwaysTrue(0, ""))
	})

	t.Run("works with different types", func(t *testing.T) {
		alwaysPi := Constant2[string, bool](3.14)
		assert.Equal(t, 3.14, alwaysPi("test", true))
	})
}

// TestIsNil tests the IsNil function
func TestIsNil(t *testing.T) {
	t.Run("returns true for nil pointer", func(t *testing.T) {
		var ptr *int
		assert.True(t, IsNil(ptr))

		var strPtr *string
		assert.True(t, IsNil(strPtr))
	})

	t.Run("returns false for non-nil pointer", func(t *testing.T) {
		value := 42
		assert.False(t, IsNil(&value))

		str := "hello"
		assert.False(t, IsNil(&str))
	})
}

// TestIsNonNil tests the IsNonNil function
func TestIsNonNil(t *testing.T) {
	t.Run("returns false for nil pointer", func(t *testing.T) {
		var ptr *int
		assert.False(t, IsNonNil(ptr))
	})

	t.Run("returns true for non-nil pointer", func(t *testing.T) {
		value := 42
		assert.True(t, IsNonNil(&value))

		str := "hello"
		assert.True(t, IsNonNil(&str))
	})
}

// TestSwap tests the Swap function
func TestSwap(t *testing.T) {
	t.Run("swaps parameters of subtraction", func(t *testing.T) {
		subtract := func(a, b int) int { return a - b }
		swapped := Swap(subtract)

		assert.Equal(t, 7, subtract(10, 3)) // 10 - 3
		assert.Equal(t, -7, swapped(10, 3)) // 3 - 10
	})

	t.Run("swaps parameters of division", func(t *testing.T) {
		divide := func(a, b float64) float64 { return a / b }
		swapped := Swap(divide)

		assert.Equal(t, 5.0, divide(10, 2))  // 10 / 2
		assert.Equal(t, 0.2, swapped(10, 2)) // 2 / 10
	})

	t.Run("swaps parameters of string concatenation", func(t *testing.T) {
		concat := func(a, b string) string { return a + b }
		swapped := Swap(concat)

		assert.Equal(t, "HelloWorld", concat("Hello", "World"))
		assert.Equal(t, "WorldHello", swapped("Hello", "World"))
	})
}

// TestFirst tests the First function
func TestFirst(t *testing.T) {
	t.Run("returns first of two ints", func(t *testing.T) {
		assert.Equal(t, 42, First(42, 100))
		assert.Equal(t, 0, First(0, 1))
	})

	t.Run("returns first of two strings", func(t *testing.T) {
		assert.Equal(t, "hello", First("hello", "world"))
	})

	t.Run("returns first of mixed types", func(t *testing.T) {
		assert.Equal(t, 42, First(42, "hello"))
		assert.True(t, First(true, 100))
	})
}

// TestSecond tests the Second function
func TestSecond(t *testing.T) {
	t.Run("returns second of two ints", func(t *testing.T) {
		assert.Equal(t, 100, Second(42, 100))
		assert.Equal(t, 1, Second(0, 1))
	})

	t.Run("returns second of two strings", func(t *testing.T) {
		assert.Equal(t, "world", Second("hello", "world"))
	})

	t.Run("returns second of mixed types", func(t *testing.T) {
		assert.Equal(t, "hello", Second(42, "hello"))
		assert.Equal(t, 100, Second(true, 100))
	})
}

// TestTernary tests the Ternary function
func TestTernary(t *testing.T) {
	t.Run("applies onTrue when predicate is true", func(t *testing.T) {
		isPositive := func(n int) bool { return n > 0 }
		double := func(n int) int { return n * 2 }
		negate := func(n int) int { return -n }

		transform := Ternary(isPositive, double, negate)

		assert.Equal(t, 10, transform(5))
		assert.Equal(t, 20, transform(10))
	})

	t.Run("applies onFalse when predicate is false", func(t *testing.T) {
		isPositive := func(n int) bool { return n > 0 }
		double := func(n int) int { return n * 2 }
		negate := func(n int) int { return -n }

		transform := Ternary(isPositive, double, negate)

		assert.Equal(t, 3, transform(-3))
		assert.Equal(t, 5, transform(-5))
		assert.Equal(t, 0, transform(0))
	})

	t.Run("works with string classification", func(t *testing.T) {
		isPositive := func(n int) bool { return n > 0 }
		classify := Ternary(
			isPositive,
			Constant1[int]("positive"),
			Constant1[int]("non-positive"),
		)

		assert.Equal(t, "positive", classify(5))
		assert.Equal(t, "non-positive", classify(-3))
		assert.Equal(t, "non-positive", classify(0))
	})
}

// TestRef tests the Ref function
func TestRef(t *testing.T) {
	t.Run("creates pointer to int", func(t *testing.T) {
		value := 42
		ptr := Ref(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, 42, *ptr)
	})

	t.Run("creates pointer to string", func(t *testing.T) {
		str := "hello"
		ptr := Ref(str)
		assert.NotNil(t, ptr)
		assert.Equal(t, "hello", *ptr)
	})

	t.Run("creates pointer to struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		p := Person{Name: "Alice", Age: 30}
		ptr := Ref(p)
		assert.NotNil(t, ptr)
		assert.Equal(t, "Alice", ptr.Name)
		assert.Equal(t, 30, ptr.Age)
	})
}

// TestDeref tests the Deref function
func TestDeref(t *testing.T) {
	t.Run("dereferences int pointer", func(t *testing.T) {
		value := 42
		ptr := &value
		assert.Equal(t, 42, Deref(ptr))
	})

	t.Run("dereferences string pointer", func(t *testing.T) {
		str := "hello"
		ptr := &str
		assert.Equal(t, "hello", Deref(ptr))
	})

	t.Run("round trip with Ref", func(t *testing.T) {
		original := "test"
		copy := Deref(Ref(original))
		assert.Equal(t, original, copy)
	})
}

// TestToAny tests the ToAny function
func TestToAny(t *testing.T) {
	t.Run("converts int to any", func(t *testing.T) {
		value := 42
		anyValue := ToAny(value)
		assert.Equal(t, any(42), anyValue)
	})

	t.Run("converts string to any", func(t *testing.T) {
		str := "hello"
		anyStr := ToAny(str)
		assert.Equal(t, any("hello"), anyStr)
	})

	t.Run("converts bool to any", func(t *testing.T) {
		b := true
		anyBool := ToAny(b)
		assert.Equal(t, any(true), anyBool)
	})
}

// TestConstNil tests the ConstNil function
func TestConstNil(t *testing.T) {
	t.Run("returns nil int pointer", func(t *testing.T) {
		nilInt := ConstNil[int]()
		assert.Nil(t, nilInt)
		assert.True(t, IsNil(nilInt))
	})

	t.Run("returns nil string pointer", func(t *testing.T) {
		nilString := ConstNil[string]()
		assert.Nil(t, nilString)
		assert.True(t, IsNil(nilString))
	})

	t.Run("returns nil struct pointer", func(t *testing.T) {
		type Person struct {
			Name string
		}
		nilPerson := ConstNil[Person]()
		assert.Nil(t, nilPerson)
	})
}

// TestConstTrue tests the ConstTrue constant
func TestConstTrue(t *testing.T) {
	t.Run("always returns true", func(t *testing.T) {
		assert.True(t, ConstTrue())
		assert.True(t, ConstTrue())
	})
}

// TestConstFalse tests the ConstFalse constant
func TestConstFalse(t *testing.T) {
	t.Run("always returns false", func(t *testing.T) {
		assert.False(t, ConstFalse())
		assert.False(t, ConstFalse())
	})
}

// TestSwitch tests the Switch function
func TestSwitch(t *testing.T) {
	type Animal struct {
		Type string
		Name string
	}

	getType := func(a Animal) string { return a.Type }

	handlers := map[string]func(Animal) string{
		"dog": func(a Animal) string { return a.Name + " barks" },
		"cat": func(a Animal) string { return a.Name + " meows" },
	}

	defaultHandler := func(a Animal) string {
		return a.Name + " makes a sound"
	}

	makeSound := Switch(getType, handlers, defaultHandler)

	t.Run("applies handler for dog", func(t *testing.T) {
		dog := Animal{Type: "dog", Name: "Rex"}
		assert.Equal(t, "Rex barks", makeSound(dog))
	})

	t.Run("applies handler for cat", func(t *testing.T) {
		cat := Animal{Type: "cat", Name: "Whiskers"}
		assert.Equal(t, "Whiskers meows", makeSound(cat))
	})

	t.Run("applies default handler for unknown type", func(t *testing.T) {
		bird := Animal{Type: "bird", Name: "Tweety"}
		assert.Equal(t, "Tweety makes a sound", makeSound(bird))
	})
}

// TestPipeAndFlow tests basic Pipe and Flow functions
func TestPipeAndFlow(t *testing.T) {
	t.Run("Pipe1 applies function", func(t *testing.T) {
		double := func(n int) int { return n * 2 }
		result := Pipe1(5, double)
		assert.Equal(t, 10, result)
	})

	t.Run("Pipe3 composes functions left-to-right", func(t *testing.T) {
		add1 := func(n int) int { return n + 1 }
		double := func(n int) int { return n * 2 }
		square := func(n int) int { return n * n }

		// (5 + 1) * 2 = 12, then 12 * 12 = 144
		result := Pipe3(5, add1, double, square)
		assert.Equal(t, 144, result)
	})

	t.Run("Flow3 creates composed function", func(t *testing.T) {
		add1 := func(n int) int { return n + 1 }
		double := func(n int) int { return n * 2 }
		square := func(n int) int { return n * n }

		// Flow3 composes left-to-right like Pipe3
		// Flow3(f1, f2, f3)(x) = f3(f2(f1(x)))
		// So Flow3(add1, double, square)(5) = square(double(add1(5)))
		// = square(double(6)) = square(12) = 144
		composed := Flow3(add1, double, square)
		result := composed(5)
		assert.Equal(t, 144, result)
	})
}

// TestCurry tests currying functions
func TestCurry(t *testing.T) {
	t.Run("Curry2 curries binary function", func(t *testing.T) {
		add := func(a, b int) int { return a + b }
		curriedAdd := Curry2(add)

		add5 := curriedAdd(5)
		assert.Equal(t, 8, add5(3))
		assert.Equal(t, 10, add5(5))
	})
}
