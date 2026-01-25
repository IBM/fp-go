package iter

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/function"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestLast test getting the last element from a non-empty sequence
func TestLastSimple(t *testing.T) {

	t.Run("returns last element from integer sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		last := Last(seq)
		assert.Equal(t, O.Of(3), last)
	})

	t.Run("returns last element from string sequence", func(t *testing.T) {
		seq := From("a", "b", "c")
		last := Last(seq)
		assert.Equal(t, O.Of("c"), last)
	})

	t.Run("returns last element from single element sequence", func(t *testing.T) {
		seq := From(42)
		last := Last(seq)
		assert.Equal(t, O.Of(42), last)
	})

	t.Run("returns last element from large sequence", func(t *testing.T) {
		seq := From(100, 200, 300, 400, 500)
		last := Last(seq)
		assert.Equal(t, O.Of(500), last)
	})
}

// TestLastEmpty tests getting the last element from an empty sequence
func TestLastEmpty(t *testing.T) {

	t.Run("returns None for empty integer sequence", func(t *testing.T) {
		seq := Empty[int]()
		last := Last(seq)
		assert.Equal(t, O.None[int](), last)
	})

	t.Run("returns None for empty string sequence", func(t *testing.T) {
		seq := Empty[string]()
		last := Last(seq)
		assert.Equal(t, O.None[string](), last)
	})

	t.Run("returns None for empty struct sequence", func(t *testing.T) {
		type TestStruct struct {
			Value int
		}
		seq := Empty[TestStruct]()
		last := Last(seq)
		assert.Equal(t, O.None[TestStruct](), last)
	})

	t.Run("returns None for empty sequence of functions", func(t *testing.T) {
		type TestFunc func(int)
		seq := Empty[TestFunc]()
		last := Last(seq)
		assert.Equal(t, O.None[TestFunc](), last)
	})
}

// TestLastWithComplex tests Last with complex types
func TestLastWithComplex(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("returns last person", func(t *testing.T) {
		seq := From(
			Person{"Alice", 30},
			Person{"Bob", 25},
			Person{"Charlie", 35},
		)
		last := Last(seq)
		expected := O.Of(Person{"Charlie", 35})
		assert.Equal(t, expected, last)
	})

	t.Run("returns last pointer", func(t *testing.T) {
		p1 := &Person{"Alice", 30}
		p2 := &Person{"Bob", 25}
		seq := From(p1, p2)
		last := Last(seq)
		assert.Equal(t, O.Of(p2), last)
	})
}

func TestLastWithFunctions(t *testing.T) {

	t.Run("return function", func(t *testing.T) {

		want := "last"
		f1 := function.Constant("first")
		f2 := function.Constant("last")
		seq := From(f1, f2)

		getLast := function.Flow2(
			Last,
			O.Map(funcReader),
		)
		assert.Equal(t, O.Of(want), getLast(seq))
	})
}

func funcReader(f func() string) string {
	return f()
}

// TestLastWithChan tests Last with channels
func TestLastWithChan(t *testing.T) {
	t.Run("return function", func(t *testing.T) {
		want := 30
		seq := From(intChan(10),
			intChan(20),
			intChan(want))

		getLast := function.Flow2(
			Last,
			O.Map(chanReader[int]),
		)
		assert.Equal(t, O.Of(want), getLast(seq))

	})
}

func chanReader[T any](c <-chan T) T {
	return <-c
}

func intChan(val int) <-chan int {
	ch := make(chan int, 1)
	ch <- val
	close(ch)
	return ch
}

// TestLastWithChainedOperations tests Last with multiple chained operations
func TestLastWithChainedOperations(t *testing.T) {
	t.Run("chains filter, map, and last", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		filtered := MonadFilter(seq, N.MoreThan(5))
		mapped := MonadMap(filtered, N.Mul(10))
		result := Last(mapped)
		assert.Equal(t, O.Of(100), result)
	})

	t.Run("chains map and filter", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		mapped := MonadMap(seq, N.Mul(2))
		filtered := MonadFilter(mapped, N.MoreThan(5))
		result := Last(filtered)
		assert.Equal(t, O.Of(10), result)
	})
}

// TestLastWithReplicate tests Last with replicated values
func TestLastWithReplicate(t *testing.T) {
	t.Run("returns last from replicated sequence", func(t *testing.T) {
		seq := Replicate(5, 42)
		last := Last(seq)
		assert.Equal(t, O.Of(42), last)
	})

	t.Run("returns None from zero replications", func(t *testing.T) {
		seq := Replicate(0, 42)
		last := Last(seq)
		assert.Equal(t, O.None[int](), last)
	})
}

// TestLastWithMakeBy tests Last with MakeBy
func TestLastWithMakeBy(t *testing.T) {
	t.Run("returns last generated element", func(t *testing.T) {
		seq := MakeBy(5, func(i int) int { return i * i })
		last := Last(seq)
		assert.Equal(t, O.Of(16), last)
	})

	t.Run("returns None for zero elements", func(t *testing.T) {
		seq := MakeBy(0, F.Identity[int])
		last := Last(seq)
		assert.Equal(t, O.None[int](), last)
	})
}

// TestLastWithPrepend tests Last with Prepend
func TestLastWithPrepend(t *testing.T) {
	t.Run("returns last element, not prepended", func(t *testing.T) {
		seq := From(2, 3, 4)
		prepended := Prepend(1)(seq)
		last := Last(prepended)
		assert.Equal(t, O.Of(4), last)
	})

	t.Run("returns prepended element from empty sequence", func(t *testing.T) {
		seq := Empty[int]()
		prepended := Prepend(42)(seq)
		last := Last(prepended)
		assert.Equal(t, O.Of(42), last)
	})
}

// TestLastWithAppend tests Last with Append
func TestLastWithAppend(t *testing.T) {
	t.Run("returns appended element", func(t *testing.T) {
		seq := From(1, 2, 3)
		appended := Append(4)(seq)
		last := Last(appended)
		assert.Equal(t, O.Of(4), last)
	})

	t.Run("returns appended element from empty sequence", func(t *testing.T) {
		seq := Empty[int]()
		appended := Append(42)(seq)
		last := Last(appended)
		assert.Equal(t, O.Of(42), last)
	})
}

// TestLastWithChain tests Last with Chain (flatMap)
func TestLastWithChain(t *testing.T) {
	t.Run("returns last from chained sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		chained := MonadChain(seq, func(x int) Seq[int] {
			return From(x, x*10)
		})
		last := Last(chained)
		assert.Equal(t, O.Of(30), last)
	})

	t.Run("returns None when chain produces empty", func(t *testing.T) {
		seq := From(1, 2, 3)
		chained := MonadChain(seq, func(x int) Seq[int] {
			return Empty[int]()
		})
		last := Last(chained)
		assert.Equal(t, O.None[int](), last)
	})
}

// TestLastWithFlatten tests Last with Flatten
func TestLastWithFlatten(t *testing.T) {
	t.Run("returns last from flattened sequence", func(t *testing.T) {
		nested := From(From(1, 2), From(3, 4), From(5))
		flattened := Flatten(nested)
		last := Last(flattened)
		assert.Equal(t, O.Of(5), last)
	})

	t.Run("returns None from empty nested sequence", func(t *testing.T) {
		nested := Empty[Seq[int]]()
		flattened := Flatten(nested)
		last := Last(flattened)
		assert.Equal(t, O.None[int](), last)
	})
}

// Example tests for documentation
func ExampleLast() {
	seq := From(1, 2, 3, 4, 5)
	last := Last(seq)

	if value, ok := O.Unwrap(last); ok {
		fmt.Printf("Last element: %d\n", value)
	}
	// Output: Last element: 5
}

func ExampleLast_empty() {
	seq := Empty[int]()
	last := Last(seq)

	if _, ok := O.Unwrap(last); !ok {
		fmt.Println("Sequence is empty")
	}
	// Output: Sequence is empty
}

func ExampleLast_functions() {
	f1 := function.Constant("first")
	f2 := function.Constant("middle")
	f3 := function.Constant("last")
	seq := From(f1, f2, f3)

	last := Last(seq)

	if fn, ok := O.Unwrap(last); ok {
		result := fn()
		fmt.Printf("Last function result: %s\n", result)
	}
	// Output: Last function result: last
}
