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

package nonempty

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	STR "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestToNonEmptyArray tests the ToNonEmptyArray function
func TestToNonEmptyArray(t *testing.T) {
	t.Run("Convert non-empty slice of integers", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From(0)))(result)
		assert.Equal(t, 3, Size(nea))
		assert.Equal(t, 1, Head(nea))
		assert.Equal(t, 3, Last(nea))
	})

	t.Run("Convert empty slice returns None", func(t *testing.T) {
		input := []int{}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsNone(result))
	})

	t.Run("Convert single element slice", func(t *testing.T) {
		input := []string{"hello"}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From("")))(result)
		assert.Equal(t, 1, Size(nea))
		assert.Equal(t, "hello", Head(nea))
	})

	t.Run("Convert non-empty slice of strings", func(t *testing.T) {
		input := []string{"a", "b", "c", "d"}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From("")))(result)
		assert.Equal(t, 4, Size(nea))
		assert.Equal(t, "a", Head(nea))
		assert.Equal(t, "d", Last(nea))
	})

	t.Run("Convert nil slice returns None", func(t *testing.T) {
		var input []int
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsNone(result))
	})

	t.Run("Convert slice with struct elements", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		input := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From(Person{})))(result)
		assert.Equal(t, 2, Size(nea))
		assert.Equal(t, "Alice", Head(nea).Name)
	})

	t.Run("Convert slice with pointer elements", func(t *testing.T) {
		val1, val2 := 10, 20
		input := []*int{&val1, &val2}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From[*int](nil)))(result)
		assert.Equal(t, 2, Size(nea))
		assert.Equal(t, 10, *Head(nea))
	})

	t.Run("Convert large slice", func(t *testing.T) {
		input := make([]int, 1000)
		for i := range input {
			input[i] = i
		}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From(0)))(result)
		assert.Equal(t, 1000, Size(nea))
		assert.Equal(t, 0, Head(nea))
		assert.Equal(t, 999, Last(nea))
	})

	t.Run("Convert slice with float64 elements", func(t *testing.T) {
		input := []float64{1.5, 2.5, 3.5}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From(0.0)))(result)
		assert.Equal(t, 3, Size(nea))
		assert.Equal(t, 1.5, Head(nea))
	})

	t.Run("Convert slice with boolean elements", func(t *testing.T) {
		input := []bool{true, false, true}
		result := ToNonEmptyArray(input)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From(false)))(result)
		assert.Equal(t, 3, Size(nea))
		assert.True(t, Head(nea))
	})
}

// TestToNonEmptyArrayWithOption tests ToNonEmptyArray with Option operations
func TestToNonEmptyArrayWithOption(t *testing.T) {
	t.Run("Chain with Map to process elements", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := F.Pipe2(
			input,
			ToNonEmptyArray[int],
			O.Map(Map(func(x int) int { return x * 2 })),
		)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From(0)))(result)
		assert.Equal(t, 2, Head(nea))
		assert.Equal(t, 6, Last(nea))
	})

	t.Run("Chain with Map to get head", func(t *testing.T) {
		input := []string{"first", "second", "third"}
		result := F.Pipe2(
			input,
			ToNonEmptyArray[string],
			O.Map(Head[string]),
		)

		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, "first", value)
	})

	t.Run("GetOrElse with default value for empty slice", func(t *testing.T) {
		input := []int{}
		defaultValue := From(42)
		result := F.Pipe2(
			input,
			ToNonEmptyArray[int],
			O.GetOrElse(F.Constant(defaultValue)),
		)

		assert.Equal(t, 1, Size(result))
		assert.Equal(t, 42, Head(result))
	})

	t.Run("GetOrElse with default value for non-empty slice", func(t *testing.T) {
		input := []int{1, 2, 3}
		defaultValue := From(42)
		result := F.Pipe2(
			input,
			ToNonEmptyArray[int],
			O.GetOrElse(F.Constant(defaultValue)),
		)

		assert.Equal(t, 3, Size(result))
		assert.Equal(t, 1, Head(result))
	})

	t.Run("Fold with Some case", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := F.Pipe2(
			input,
			ToNonEmptyArray[int],
			O.Fold(
				F.Constant(0),
				func(nea NonEmptyArray[int]) int { return Head(nea) },
			),
		)

		assert.Equal(t, 1, result)
	})

	t.Run("Fold with None case", func(t *testing.T) {
		input := []int{}
		result := F.Pipe2(
			input,
			ToNonEmptyArray[int],
			O.Fold(
				F.Constant(-1),
				func(nea NonEmptyArray[int]) int { return Head(nea) },
			),
		)

		assert.Equal(t, -1, result)
	})
}

// TestToNonEmptyArrayComposition tests composing ToNonEmptyArray with other operations
func TestToNonEmptyArrayComposition(t *testing.T) {
	t.Run("Compose with filter-like operation", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		// Filter even numbers then convert
		filtered := []int{}
		for _, x := range input {
			if x%2 == 0 {
				filtered = append(filtered, x)
			}
		}
		result := ToNonEmptyArray(filtered)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From(0)))(result)
		assert.Equal(t, 2, Size(nea))
		assert.Equal(t, 2, Head(nea))
	})

	t.Run("Compose with map operation before conversion", func(t *testing.T) {
		input := []int{1, 2, 3}
		// Map then convert
		mapped := make([]int, len(input))
		for i, x := range input {
			mapped[i] = x * 10
		}
		result := ToNonEmptyArray(mapped)

		assert.True(t, O.IsSome(result))
		nea := O.GetOrElse(F.Constant(From(0)))(result)
		assert.Equal(t, 10, Head(nea))
		assert.Equal(t, 30, Last(nea))
	})

	t.Run("Chain multiple Option operations", func(t *testing.T) {
		input := []int{5, 10, 15}
		result := F.Pipe3(
			input,
			ToNonEmptyArray[int],
			O.Map(Map(func(x int) int { return x / 5 })),
			O.Map(func(nea NonEmptyArray[int]) int {
				return Head(nea) + Last(nea)
			}),
		)

		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(0))(result)
		assert.Equal(t, 4, value) // 1 + 3
	})
}

// TestToNonEmptyArrayUseCases demonstrates practical use cases
func TestToNonEmptyArrayUseCases(t *testing.T) {
	t.Run("Validate user input has at least one item", func(t *testing.T) {
		// Simulate user input
		userInput := []string{"item1", "item2"}

		result := ToNonEmptyArray(userInput)
		if O.IsSome(result) {
			nea := O.GetOrElse(F.Constant(From("")))(result)
			firstItem := Head(nea)
			assert.Equal(t, "item1", firstItem)
		} else {
			t.Fatal("Expected Some but got None")
		}
	})

	t.Run("Process only non-empty collections", func(t *testing.T) {
		processItems := func(items []int) Option[int] {
			return F.Pipe2(
				items,
				ToNonEmptyArray[int],
				O.Map(func(nea NonEmptyArray[int]) int {
					// Safe to use Head since we know it's non-empty
					return Head(nea) * 2
				}),
			)
		}

		result1 := processItems([]int{5, 10, 15})
		assert.True(t, O.IsSome(result1))
		assert.Equal(t, 10, O.GetOrElse(F.Constant(0))(result1))

		result2 := processItems([]int{})
		assert.True(t, O.IsNone(result2))
	})

	t.Run("Convert API response to NonEmptyArray", func(t *testing.T) {
		// Simulate API response
		type APIResponse struct {
			Items []string
		}

		response := APIResponse{Items: []string{"data1", "data2", "data3"}}

		result := F.Pipe2(
			response.Items,
			ToNonEmptyArray[string],
			O.Map(func(nea NonEmptyArray[string]) string {
				return "First item: " + Head(nea)
			}),
		)

		assert.True(t, O.IsSome(result))
		message := O.GetOrElse(F.Constant("No items"))(result)
		assert.Equal(t, "First item: data1", message)
	})

	t.Run("Ensure collection is non-empty before processing", func(t *testing.T) {
		calculateAverage := func(numbers []float64) Option[float64] {
			return F.Pipe2(
				numbers,
				ToNonEmptyArray[float64],
				O.Map(func(nea NonEmptyArray[float64]) float64 {
					sum := 0.0
					for _, n := range nea {
						sum += n
					}
					return sum / float64(Size(nea))
				}),
			)
		}

		result1 := calculateAverage([]float64{10.0, 20.0, 30.0})
		assert.True(t, O.IsSome(result1))
		assert.Equal(t, 20.0, O.GetOrElse(F.Constant(0.0))(result1))

		result2 := calculateAverage([]float64{})
		assert.True(t, O.IsNone(result2))
	})

	t.Run("Safe head extraction with type guarantee", func(t *testing.T) {
		getFirstOrDefault := func(items []string, defaultValue string) string {
			return F.Pipe2(
				items,
				ToNonEmptyArray[string],
				O.Fold(
					F.Constant(defaultValue),
					Head[string],
				),
			)
		}

		result1 := getFirstOrDefault([]string{"a", "b", "c"}, "default")
		assert.Equal(t, "a", result1)

		result2 := getFirstOrDefault([]string{}, "default")
		assert.Equal(t, "default", result2)
	})
}

// TestOf tests the Of function
func TestOf(t *testing.T) {
	t.Run("Create single element array with int", func(t *testing.T) {
		arr := Of(42)
		assert.Equal(t, 1, Size(arr))
		assert.Equal(t, 42, Head(arr))
	})

	t.Run("Create single element array with string", func(t *testing.T) {
		arr := Of("hello")
		assert.Equal(t, 1, Size(arr))
		assert.Equal(t, "hello", Head(arr))
	})

	t.Run("Create single element array with struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		person := Person{Name: "Alice", Age: 30}
		arr := Of(person)
		assert.Equal(t, 1, Size(arr))
		assert.Equal(t, "Alice", Head(arr).Name)
	})
}

// TestFrom tests the From function
func TestFrom(t *testing.T) {
	t.Run("Create array with single element", func(t *testing.T) {
		arr := From(1)
		assert.Equal(t, 1, Size(arr))
		assert.Equal(t, 1, Head(arr))
	})

	t.Run("Create array with multiple elements", func(t *testing.T) {
		arr := From(1, 2, 3, 4, 5)
		assert.Equal(t, 5, Size(arr))
		assert.Equal(t, 1, Head(arr))
		assert.Equal(t, 5, Last(arr))
	})

	t.Run("Create array with strings", func(t *testing.T) {
		arr := From("a", "b", "c")
		assert.Equal(t, 3, Size(arr))
		assert.Equal(t, "a", Head(arr))
		assert.Equal(t, "c", Last(arr))
	})
}

// TestIsEmpty tests the IsEmpty function
func TestIsEmpty(t *testing.T) {
	t.Run("IsEmpty always returns false", func(t *testing.T) {
		arr := From(1, 2, 3)
		assert.False(t, IsEmpty(arr))
	})

	t.Run("IsEmpty returns false for single element", func(t *testing.T) {
		arr := Of(1)
		assert.False(t, IsEmpty(arr))
	})
}

// TestIsNonEmpty tests the IsNonEmpty function
func TestIsNonEmpty(t *testing.T) {
	t.Run("IsNonEmpty always returns true", func(t *testing.T) {
		arr := From(1, 2, 3)
		assert.True(t, IsNonEmpty(arr))
	})

	t.Run("IsNonEmpty returns true for single element", func(t *testing.T) {
		arr := Of(1)
		assert.True(t, IsNonEmpty(arr))
	})
}

// TestMonadMap tests the MonadMap function
func TestMonadMap(t *testing.T) {
	t.Run("Map integers to doubles", func(t *testing.T) {
		arr := From(1, 2, 3, 4)
		result := MonadMap(arr, func(x int) int { return x * 2 })
		assert.Equal(t, 4, Size(result))
		assert.Equal(t, 2, Head(result))
		assert.Equal(t, 8, Last(result))
	})

	t.Run("Map strings to lengths", func(t *testing.T) {
		arr := From("a", "bb", "ccc")
		result := MonadMap(arr, func(s string) int { return len(s) })
		assert.Equal(t, 3, Size(result))
		assert.Equal(t, 1, Head(result))
		assert.Equal(t, 3, Last(result))
	})

	t.Run("Map single element", func(t *testing.T) {
		arr := Of(5)
		result := MonadMap(arr, func(x int) int { return x * 10 })
		assert.Equal(t, 1, Size(result))
		assert.Equal(t, 50, Head(result))
	})
}

// TestMap tests the Map function
func TestMap(t *testing.T) {
	t.Run("Curried map with integers", func(t *testing.T) {
		double := Map(func(x int) int { return x * 2 })
		arr := From(1, 2, 3)
		result := double(arr)
		assert.Equal(t, 3, Size(result))
		assert.Equal(t, 2, Head(result))
		assert.Equal(t, 6, Last(result))
	})

	t.Run("Curried map with strings", func(t *testing.T) {
		toUpper := Map(func(s string) string { return s + "!" })
		arr := From("hello", "world")
		result := toUpper(arr)
		assert.Equal(t, 2, Size(result))
		assert.Equal(t, "hello!", Head(result))
		assert.Equal(t, "world!", Last(result))
	})
}

// TestReduce tests the Reduce function
func TestReduce(t *testing.T) {
	t.Run("Sum integers", func(t *testing.T) {
		sum := Reduce(func(acc int, x int) int { return acc + x }, 0)
		arr := From(1, 2, 3, 4, 5)
		result := sum(arr)
		assert.Equal(t, 15, result)
	})

	t.Run("Concatenate strings", func(t *testing.T) {
		concat := Reduce(func(acc string, x string) string { return acc + x }, "")
		arr := From("a", "b", "c")
		result := concat(arr)
		assert.Equal(t, "abc", result)
	})

	t.Run("Product of numbers", func(t *testing.T) {
		product := Reduce(func(acc int, x int) int { return acc * x }, 1)
		arr := From(2, 3, 4)
		result := product(arr)
		assert.Equal(t, 24, result)
	})

	t.Run("Reduce single element", func(t *testing.T) {
		sum := Reduce(func(acc int, x int) int { return acc + x }, 10)
		arr := Of(5)
		result := sum(arr)
		assert.Equal(t, 15, result)
	})
}

// TestReduceRight tests the ReduceRight function
func TestReduceRight(t *testing.T) {
	t.Run("Concatenate strings right to left", func(t *testing.T) {
		concat := ReduceRight(func(x string, acc string) string { return acc + x }, "")
		arr := From("a", "b", "c")
		result := concat(arr)
		assert.Equal(t, "cba", result)
	})

	t.Run("Build list right to left", func(t *testing.T) {
		buildList := ReduceRight(func(x int, acc []int) []int { return append(acc, x) }, []int{})
		arr := From(1, 2, 3)
		result := buildList(arr)
		assert.Equal(t, []int{3, 2, 1}, result)
	})
}

// TestTail tests the Tail function
func TestTail(t *testing.T) {
	t.Run("Get tail of multi-element array", func(t *testing.T) {
		arr := From(1, 2, 3, 4)
		tail := Tail(arr)
		assert.Equal(t, 3, len(tail))
		assert.Equal(t, []int{2, 3, 4}, tail)
	})

	t.Run("Get tail of single element array", func(t *testing.T) {
		arr := Of(1)
		tail := Tail(arr)
		assert.Equal(t, 0, len(tail))
		assert.Equal(t, []int{}, tail)
	})

	t.Run("Get tail of two element array", func(t *testing.T) {
		arr := From(1, 2)
		tail := Tail(arr)
		assert.Equal(t, 1, len(tail))
		assert.Equal(t, []int{2}, tail)
	})
}

// TestHead tests the Head function
func TestHead(t *testing.T) {
	t.Run("Get head of multi-element array", func(t *testing.T) {
		arr := From(1, 2, 3)
		head := Head(arr)
		assert.Equal(t, 1, head)
	})

	t.Run("Get head of single element array", func(t *testing.T) {
		arr := Of(42)
		head := Head(arr)
		assert.Equal(t, 42, head)
	})

	t.Run("Get head of string array", func(t *testing.T) {
		arr := From("first", "second", "third")
		head := Head(arr)
		assert.Equal(t, "first", head)
	})
}

// TestFirst tests the First function
func TestFirst(t *testing.T) {
	t.Run("First is alias for Head", func(t *testing.T) {
		arr := From(1, 2, 3)
		assert.Equal(t, Head(arr), First(arr))
	})

	t.Run("Get first element", func(t *testing.T) {
		arr := From("a", "b", "c")
		first := First(arr)
		assert.Equal(t, "a", first)
	})
}

// TestLast tests the Last function
func TestLast(t *testing.T) {
	t.Run("Get last of multi-element array", func(t *testing.T) {
		arr := From(1, 2, 3, 4, 5)
		last := Last(arr)
		assert.Equal(t, 5, last)
	})

	t.Run("Get last of single element array", func(t *testing.T) {
		arr := Of(42)
		last := Last(arr)
		assert.Equal(t, 42, last)
	})

	t.Run("Get last of string array", func(t *testing.T) {
		arr := From("first", "second", "third")
		last := Last(arr)
		assert.Equal(t, "third", last)
	})
}

// TestSize tests the Size function
func TestSize(t *testing.T) {
	t.Run("Size of multi-element array", func(t *testing.T) {
		arr := From(1, 2, 3, 4, 5)
		size := Size(arr)
		assert.Equal(t, 5, size)
	})

	t.Run("Size of single element array", func(t *testing.T) {
		arr := Of(1)
		size := Size(arr)
		assert.Equal(t, 1, size)
	})

	t.Run("Size of large array", func(t *testing.T) {
		elements := make([]int, 1000)
		arr := From(1, elements...)
		size := Size(arr)
		assert.Equal(t, 1001, size)
	})
}

// TestFlatten tests the Flatten function
func TestFlatten(t *testing.T) {
	t.Run("Flatten nested arrays", func(t *testing.T) {
		nested := From(From(1, 2), From(3, 4), From(5))
		flat := Flatten(nested)
		assert.Equal(t, 5, Size(flat))
		assert.Equal(t, 1, Head(flat))
		assert.Equal(t, 5, Last(flat))
	})

	t.Run("Flatten single nested array", func(t *testing.T) {
		nested := Of(From(1, 2, 3))
		flat := Flatten(nested)
		assert.Equal(t, 3, Size(flat))
		assert.Equal(t, []int{1, 2, 3}, []int(flat))
	})

	t.Run("Flatten arrays of different sizes", func(t *testing.T) {
		nested := From(Of(1), From(2, 3, 4), From(5, 6))
		flat := Flatten(nested)
		assert.Equal(t, 6, Size(flat))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, []int(flat))
	})
}

// TestMonadChain tests the MonadChain function
func TestMonadChain(t *testing.T) {
	t.Run("Chain with duplication", func(t *testing.T) {
		arr := From(1, 2, 3)
		result := MonadChain(arr, func(x int) NonEmptyArray[int] {
			return From(x, x*10)
		})
		assert.Equal(t, 6, Size(result))
		assert.Equal(t, []int{1, 10, 2, 20, 3, 30}, []int(result))
	})

	t.Run("Chain with expansion", func(t *testing.T) {
		arr := From(1, 2)
		result := MonadChain(arr, func(x int) NonEmptyArray[int] {
			return From(x, x+1, x+2)
		})
		assert.Equal(t, 6, Size(result))
		assert.Equal(t, []int{1, 2, 3, 2, 3, 4}, []int(result))
	})

	t.Run("Chain single element", func(t *testing.T) {
		arr := Of(5)
		result := MonadChain(arr, func(x int) NonEmptyArray[int] {
			return From(x, x*2)
		})
		assert.Equal(t, 2, Size(result))
		assert.Equal(t, []int{5, 10}, []int(result))
	})
}

// TestChain tests the Chain function
func TestChain(t *testing.T) {
	t.Run("Curried chain with duplication", func(t *testing.T) {
		duplicate := Chain(func(x int) NonEmptyArray[int] {
			return From(x, x)
		})
		arr := From(1, 2, 3)
		result := duplicate(arr)
		assert.Equal(t, 6, Size(result))
		assert.Equal(t, []int{1, 1, 2, 2, 3, 3}, []int(result))
	})

	t.Run("Curried chain with transformation", func(t *testing.T) {
		expand := Chain(func(x int) NonEmptyArray[string] {
			return Of(fmt.Sprintf("%d", x))
		})
		arr := From(1, 2, 3)
		result := expand(arr)
		assert.Equal(t, 3, Size(result))
		assert.Equal(t, "1", Head(result))
	})
}

// TestMonadAp tests the MonadAp function
func TestMonadAp(t *testing.T) {
	t.Run("Apply functions to values", func(t *testing.T) {
		fns := From(
			func(x int) int { return x * 2 },
			func(x int) int { return x + 10 },
		)
		vals := From(1, 2)
		result := MonadAp(fns, vals)
		assert.Equal(t, 4, Size(result))
		assert.Equal(t, []int{2, 4, 11, 12}, []int(result))
	})

	t.Run("Apply single function to multiple values", func(t *testing.T) {
		fns := Of(func(x int) int { return x * 3 })
		vals := From(1, 2, 3)
		result := MonadAp(fns, vals)
		assert.Equal(t, 3, Size(result))
		assert.Equal(t, []int{3, 6, 9}, []int(result))
	})
}

// TestAp tests the Ap function
func TestAp(t *testing.T) {
	t.Run("Curried apply", func(t *testing.T) {
		vals := From(1, 2)
		applyTo := Ap[int](vals)
		fns := From(
			func(x int) int { return x * 2 },
			func(x int) int { return x + 10 },
		)
		result := applyTo(fns)
		assert.Equal(t, 4, Size(result))
		assert.Equal(t, []int{2, 4, 11, 12}, []int(result))
	})
}

// TestFoldMap tests the FoldMap function
func TestFoldMap(t *testing.T) {
	t.Run("FoldMap with sum semigroup", func(t *testing.T) {
		sumSemigroup := N.SemigroupSum[int]()
		arr := From(1, 2, 3, 4)
		result := FoldMap[int](sumSemigroup)(func(x int) int { return x * 2 })(arr)
		assert.Equal(t, 20, result) // (1*2) + (2*2) + (3*2) + (4*2) = 20
	})

	t.Run("FoldMap with string concatenation", func(t *testing.T) {
		concatSemigroup := STR.Semigroup
		arr := From(1, 2, 3)
		result := FoldMap[int](concatSemigroup)(func(x int) string { return fmt.Sprintf("%d", x) })(arr)
		assert.Equal(t, "123", result)
	})
}

// TestFold tests the Fold function
func TestFold(t *testing.T) {
	t.Run("Fold with sum semigroup", func(t *testing.T) {
		sumSemigroup := N.SemigroupSum[int]()
		arr := From(1, 2, 3, 4, 5)
		result := Fold(sumSemigroup)(arr)
		assert.Equal(t, 15, result)
	})

	t.Run("Fold with string concatenation", func(t *testing.T) {
		concatSemigroup := STR.Semigroup
		arr := From("a", "b", "c")
		result := Fold(concatSemigroup)(arr)
		assert.Equal(t, "abc", result)
	})

	t.Run("Fold single element", func(t *testing.T) {
		sumSemigroup := N.SemigroupSum[int]()
		arr := Of(42)
		result := Fold(sumSemigroup)(arr)
		assert.Equal(t, 42, result)
	})
}

// TestPrepend tests the Prepend function
func TestPrepend(t *testing.T) {
	t.Run("Prepend to multi-element array", func(t *testing.T) {
		arr := From(2, 3, 4)
		prepend1 := Prepend(1)
		result := prepend1(arr)
		assert.Equal(t, 4, Size(result))
		assert.Equal(t, 1, Head(result))
		assert.Equal(t, 4, Last(result))
	})

	t.Run("Prepend to single element array", func(t *testing.T) {
		arr := Of(2)
		prepend1 := Prepend(1)
		result := prepend1(arr)
		assert.Equal(t, 2, Size(result))
		assert.Equal(t, []int{1, 2}, []int(result))
	})

	t.Run("Prepend string", func(t *testing.T) {
		arr := From("world")
		prependHello := Prepend("hello")
		result := prependHello(arr)
		assert.Equal(t, 2, Size(result))
		assert.Equal(t, "hello", Head(result))
	})
}

// TestExtract tests the Extract function
func TestExtract(t *testing.T) {
	t.Run("Extract from multi-element array", func(t *testing.T) {
		arr := From(1, 2, 3)
		result := Extract(arr)
		assert.Equal(t, 1, result)
	})

	t.Run("Extract from single element array", func(t *testing.T) {
		arr := Of(42)
		result := Extract(arr)
		assert.Equal(t, 42, result)
	})

	t.Run("Extract is same as Head", func(t *testing.T) {
		arr := From("a", "b", "c")
		assert.Equal(t, Head(arr), Extract(arr))
	})
}

// TestExtend tests the Extend function
func TestExtend(t *testing.T) {
	t.Run("Extend with sum of suffixes", func(t *testing.T) {
		arr := From(1, 2, 3, 4)
		sumSuffix := Extend(func(xs NonEmptyArray[int]) int {
			sum := 0
			for _, x := range xs {
				sum += x
			}
			return sum
		})
		result := sumSuffix(arr)
		assert.Equal(t, 4, Size(result))
		assert.Equal(t, []int{10, 9, 7, 4}, []int(result))
	})

	t.Run("Extend with head of suffixes", func(t *testing.T) {
		arr := From(1, 2, 3)
		getHeads := Extend(Head[int])
		result := getHeads(arr)
		assert.Equal(t, 3, Size(result))
		assert.Equal(t, []int{1, 2, 3}, []int(result))
	})

	t.Run("Extend with size of suffixes", func(t *testing.T) {
		arr := From("a", "b", "c", "d")
		getSizes := Extend(Size[string])
		result := getSizes(arr)
		assert.Equal(t, 4, Size(result))
		assert.Equal(t, []int{4, 3, 2, 1}, []int(result))
	})

	t.Run("Extend single element", func(t *testing.T) {
		arr := Of(5)
		double := Extend(func(xs NonEmptyArray[int]) int {
			return Head(xs) * 2
		})
		result := double(arr)
		assert.Equal(t, 1, Size(result))
		assert.Equal(t, 10, Head(result))
	})
}
