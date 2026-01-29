package io

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/stretchr/testify/assert"
)

func toUpper(s string) IO[string] {
	return Of(strings.ToUpper(s))
}

func TestTraverseArray(t *testing.T) {

	src := []string{"a", "b"}

	trv := TraverseArray(toUpper)

	res := trv(src)

	assert.Equal(t, res(), []string{"A", "B"})
}

type (
	customSlice []string
)

func TestTraverseCustomSlice(t *testing.T) {

	src := customSlice{"a", "b"}

	trv := TraverseArray(toUpper)

	res := trv(src)

	assert.Equal(t, res(), []string{"A", "B"})
}

func TestTraverseIter(t *testing.T) {
	t.Run("transforms all elements successfully", func(t *testing.T) {
		// Create an iterator of strings
		input := slices.Values(A.From("hello", "world", "test"))

		// Transform each string to uppercase
		transform := func(s string) IO[string] {
			return Of(strings.ToUpper(s))
		}

		// Traverse the iterator
		traverseFn := TraverseIter(transform)
		resultIO := traverseFn(input)

		// Execute the IO and collect results
		result := resultIO()
		var collected []string
		for s := range result {
			collected = append(collected, s)
		}

		assert.Equal(t, []string{"HELLO", "WORLD", "TEST"}, collected)
	})

	t.Run("works with empty iterator", func(t *testing.T) {
		// Create an empty iterator
		input := func(yield func(string) bool) {}

		transform := func(s string) IO[string] {
			return Of(strings.ToUpper(s))
		}

		traverseFn := TraverseIter(transform)
		resultIO := traverseFn(input)

		result := resultIO()
		var collected []string
		for s := range result {
			collected = append(collected, s)
		}

		assert.Empty(t, collected)
	})

	t.Run("works with single element", func(t *testing.T) {
		input := func(yield func(int) bool) {
			yield(42)
		}

		transform := func(n int) IO[int] {
			return Of(n * 2)
		}

		traverseFn := TraverseIter(transform)
		resultIO := traverseFn(input)

		result := resultIO()
		var collected []int
		for n := range result {
			collected = append(collected, n)
		}

		assert.Equal(t, []int{84}, collected)
	})

	t.Run("preserves order of elements", func(t *testing.T) {
		input := func(yield func(int) bool) {
			for i := 1; i <= 5; i++ {
				if !yield(i) {
					return
				}
			}
		}

		transform := func(n int) IO[string] {
			return Of(fmt.Sprintf("item-%d", n))
		}

		traverseFn := TraverseIter(transform)
		resultIO := traverseFn(input)

		result := resultIO()
		var collected []string
		for s := range result {
			collected = append(collected, s)
		}

		expected := []string{"item-1", "item-2", "item-3", "item-4", "item-5"}
		assert.Equal(t, expected, collected)
	})

	t.Run("handles complex transformations", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}

		input := func(yield func(int) bool) {
			for _, id := range []int{1, 2, 3} {
				if !yield(id) {
					return
				}
			}
		}

		transform := func(id int) IO[User] {
			return Of(User{ID: id, Name: fmt.Sprintf("User%d", id)})
		}

		traverseFn := TraverseIter(transform)
		resultIO := traverseFn(input)

		result := resultIO()
		var collected []User
		for user := range result {
			collected = append(collected, user)
		}

		expected := []User{
			{ID: 1, Name: "User1"},
			{ID: 2, Name: "User2"},
			{ID: 3, Name: "User3"},
		}
		assert.Equal(t, expected, collected)
	})
}

func TestSequenceIter(t *testing.T) {
	t.Run("sequences multiple IO operations", func(t *testing.T) {
		// Create an iterator of IO operations
		input := slices.Values(A.From(Of(1), Of(2), Of(3)))

		// Sequence the operations
		resultIO := SequenceIter(input)

		// Execute and collect results
		result := resultIO()
		var collected []int
		for n := range result {
			collected = append(collected, n)
		}

		assert.Equal(t, []int{1, 2, 3}, collected)
	})

	t.Run("works with empty iterator", func(t *testing.T) {
		input := slices.Values(A.Empty[IO[string]]())

		resultIO := SequenceIter(input)

		result := resultIO()
		var collected []string
		for s := range result {
			collected = append(collected, s)
		}

		assert.Empty(t, collected)
	})

	// TODO!!
	// t.Run("executes all IO operations", func(t *testing.T) {
	// 	// Track execution order
	// 	var executed []int

	// 	input := func(yield func(IO[int]) bool) {
	// 		yield(func() int {
	// 			executed = append(executed, 1)
	// 			return 10
	// 		})
	// 		yield(func() int {
	// 			executed = append(executed, 2)
	// 			return 20
	// 		})
	// 		yield(func() int {
	// 			executed = append(executed, 3)
	// 			return 30
	// 		})
	// 	}

	// 	resultIO := SequenceIter(input)

	// 	// Before execution, nothing should be executed
	// 	assert.Empty(t, executed)

	// 	// Execute the IO
	// 	result := resultIO()

	// 	// Collect results
	// 	var collected []int
	// 	for n := range result {
	// 		collected = append(collected, n)
	// 	}

	// 	// All operations should have been executed
	// 	assert.Equal(t, []int{1, 2, 3}, executed)
	// 	assert.Equal(t, []int{10, 20, 30}, collected)
	// })

	t.Run("works with single IO operation", func(t *testing.T) {
		input := func(yield func(IO[string]) bool) {
			yield(Of("hello"))
		}

		resultIO := SequenceIter(input)

		result := resultIO()
		var collected []string
		for s := range result {
			collected = append(collected, s)
		}

		assert.Equal(t, []string{"hello"}, collected)
	})

	t.Run("preserves order of results", func(t *testing.T) {
		input := func(yield func(IO[int]) bool) {
			for i := 5; i >= 1; i-- {
				n := i // capture loop variable
				yield(func() int { return n * 10 })
			}
		}

		resultIO := SequenceIter(input)

		result := resultIO()
		var collected []int
		for n := range result {
			collected = append(collected, n)
		}

		assert.Equal(t, []int{50, 40, 30, 20, 10}, collected)
	})

	t.Run("works with complex types", func(t *testing.T) {
		type Result struct {
			Value int
			Label string
		}

		input := func(yield func(IO[Result]) bool) {
			yield(Of(Result{Value: 1, Label: "first"}))
			yield(Of(Result{Value: 2, Label: "second"}))
			yield(Of(Result{Value: 3, Label: "third"}))
		}

		resultIO := SequenceIter(input)

		result := resultIO()
		var collected []Result
		for r := range result {
			collected = append(collected, r)
		}

		expected := []Result{
			{Value: 1, Label: "first"},
			{Value: 2, Label: "second"},
			{Value: 3, Label: "third"},
		}
		assert.Equal(t, expected, collected)
	})
}
