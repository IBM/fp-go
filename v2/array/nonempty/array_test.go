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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
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

// Made with Bob
