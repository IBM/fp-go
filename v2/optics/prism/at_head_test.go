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

package prism

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Head — unit tests
// ---------------------------------------------------------------------------

// TestHead_GetOption_Success verifies Some is returned for non-empty slices.
func TestHead_GetOption_Success(t *testing.T) {
	p := Head[int]()

	t.Run("single-element slice", func(t *testing.T) {
		got := p.GetOption([]int{42})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("multi-element slice returns first", func(t *testing.T) {
		got := p.GetOption([]int{10, 20, 30})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 10, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("string type", func(t *testing.T) {
		ps := Head[string]()
		got := ps.GetOption([]string{"hello", "world"})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant(""))(got))
	})
}

// TestHead_GetOption_Failure verifies None is returned for empty and nil slices.
func TestHead_GetOption_Failure(t *testing.T) {
	p := Head[int]()

	t.Run("empty slice", func(t *testing.T) {
		assert.True(t, O.IsNone(p.GetOption([]int{})))
	})

	t.Run("nil slice", func(t *testing.T) {
		assert.True(t, O.IsNone(p.GetOption(nil)))
	})
}

// TestHead_ReverseGet verifies that ReverseGet wraps a value in a singleton slice.
func TestHead_ReverseGet(t *testing.T) {
	t.Run("int value", func(t *testing.T) {
		assert.Equal(t, []int{99}, Head[int]().ReverseGet(99))
	})

	t.Run("string value", func(t *testing.T) {
		assert.Equal(t, []string{"hi"}, Head[string]().ReverseGet("hi"))
	})
}

// TestHead_PrismLaws verifies the two standard prism laws.
func TestHead_PrismLaws(t *testing.T) {
	p := Head[int]()

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		// ReverseGet produces [a], so GetOption must return Some(a).
		for _, a := range []int{0, 1, -5, 100} {
			got := p.GetOption(p.ReverseGet(a))
			assert.True(t, O.IsSome(got), "expected Some for value %d", a)
			assert.Equal(t, a, O.GetOrElse(F.Constant(-9999))(got))
		}
	})

	t.Run("law 2: if GetOption(s)==Some(a) then ReverseGet(a) is valid []T", func(t *testing.T) {
		s := []int{7, 8, 9}
		got := p.GetOption(s)
		if O.IsSome(got) {
			a := O.GetOrElse(F.Constant(0))(got)
			reconstructed := p.ReverseGet(a)
			// ReverseGet of the extracted head is a non-empty slice with that head.
			assert.Equal(t, O.Some(a), p.GetOption(reconstructed))
		}
	})
}

// TestHead_WithSet verifies the Set helper works correctly through Head.
func TestHead_WithSet(t *testing.T) {
	p := Head[int]()

	t.Run("set on non-empty slice replaces head", func(t *testing.T) {
		result := Set[[]int](99)(p)([]int{1, 2, 3})
		// Set extracts head (1), replaces with 99, then calls ReverseGet(99) == [99].
		// The original tail is lost because ReverseGet only produces a singleton.
		assert.Equal(t, []int{99}, result)
	})

	t.Run("set on empty slice returns slice unchanged", func(t *testing.T) {
		result := Set[[]int](99)(p)([]int{})
		assert.Equal(t, []int{}, result)
	})
}

// TestHead_Composition verifies Head can be composed with other prisms.
func TestHead_Composition(t *testing.T) {
	t.Run("compose Head with ParseInt to read first element as number", func(t *testing.T) {
		// Prism[[]string, string] ∘ Prism[string, int] → Prism[[]string, int]
		composed := Compose[[]string](ParseInt())(Head[string]())

		got := composed.GetOption([]string{"42", "not-a-number"})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(got))

		none := composed.GetOption([]string{"abc"})
		assert.True(t, O.IsNone(none))

		empty := composed.GetOption([]string{})
		assert.True(t, O.IsNone(empty))
	})
}

// ---------------------------------------------------------------------------
// At — unit tests
// ---------------------------------------------------------------------------

// TestAt_GetOption_Success verifies Some is returned for valid indices.
func TestAt_GetOption_Success(t *testing.T) {
	at := At[int]()

	t.Run("index 0", func(t *testing.T) {
		got := at(0).GetOption([]int{10, 20, 30})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 10, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("middle index", func(t *testing.T) {
		got := at(2).GetOption([]int{10, 20, 30})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 30, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("last valid index", func(t *testing.T) {
		got := at(4).GetOption([]int{1, 2, 3, 4, 5})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 5, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("string type", func(t *testing.T) {
		got := At[string]()(1).GetOption([]string{"a", "b", "c"})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, "b", O.GetOrElse(F.Constant(""))(got))
	})
}

// TestAt_GetOption_Failure verifies None is returned when index is out of bounds or negative.
func TestAt_GetOption_Failure(t *testing.T) {
	at := At[int]()

	t.Run("index equals length", func(t *testing.T) {
		assert.True(t, O.IsNone(at(3).GetOption([]int{1, 2, 3})))
	})

	t.Run("index beyond length", func(t *testing.T) {
		assert.True(t, O.IsNone(at(100).GetOption([]int{1, 2, 3})))
	})

	t.Run("negative index", func(t *testing.T) {
		assert.True(t, O.IsNone(at(-1).GetOption([]int{1, 2, 3})))
	})

	t.Run("empty slice", func(t *testing.T) {
		assert.True(t, O.IsNone(at(0).GetOption([]int{})))
	})

	t.Run("nil slice", func(t *testing.T) {
		assert.True(t, O.IsNone(at(0).GetOption(nil)))
	})
}

// TestAt_ReverseGet verifies the reconstructed slice has the value at position i.
func TestAt_ReverseGet(t *testing.T) {
	at := At[int]()

	t.Run("index 0 produces singleton", func(t *testing.T) {
		result := at(0).ReverseGet(42)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, 42, result[0])
	})

	t.Run("index 2 produces length-3 slice with value at position 2", func(t *testing.T) {
		result := at(2).ReverseGet(99)
		assert.Equal(t, 3, len(result))
		assert.Equal(t, 99, result[2])
		// positions 0 and 1 are zero-valued
		assert.Equal(t, 0, result[0])
		assert.Equal(t, 0, result[1])
	})

	t.Run("negative index produces nil", func(t *testing.T) {
		result := at(-1).ReverseGet(7)
		assert.Nil(t, result)
	})
}

// TestAt_PrismLaws verifies the two standard prism laws for non-negative indices.
func TestAt_PrismLaws(t *testing.T) {
	at := At[int]()

	for _, i := range []int{0, 1, 3} {
		i := i
		p := at(i)
		t.Run(fmt.Sprintf("law 1 i=%d: GetOption(ReverseGet(a)) == Some(a)", i), func(t *testing.T) {
			a := 42
			got := p.GetOption(p.ReverseGet(a))
			assert.True(t, O.IsSome(got))
			assert.Equal(t, a, O.GetOrElse(F.Constant(-1))(got))
		})

		t.Run(fmt.Sprintf("law 2 i=%d: if GetOption(s)==Some(a) then s[i]==a", i), func(t *testing.T) {
			s := []int{10, 20, 30, 40, 50}
			got := p.GetOption(s)
			if O.IsSome(got) {
				a := O.GetOrElse(F.Constant(-1))(got)
				assert.Equal(t, s[i], a)
			}
		})
	}
}

// TestAt_WithSet verifies Set works correctly through At.
func TestAt_WithSet(t *testing.T) {
	at := At[int]()

	t.Run("set at valid index replaces element in reconstructed slice", func(t *testing.T) {
		// Set extracts element at index 1, replaces it with 99,
		// then calls ReverseGet(99) which produces a length-2 slice.
		result := Set[[]int](99)(at(1))([]int{10, 20, 30})
		assert.Equal(t, 2, len(result))
		assert.Equal(t, 99, result[1])
	})

	t.Run("set on slice too short leaves it unchanged", func(t *testing.T) {
		// Index 5 is out of bounds; GetOption returns None so Set is a no-op.
		original := []int{1, 2, 3}
		result := Set[[]int](99)(at(5))(original)
		assert.Equal(t, original, result)
	})

	t.Run("set at negative index leaves slice unchanged", func(t *testing.T) {
		original := []int{1, 2, 3}
		result := Set[[]int](99)(at(-1))(original)
		assert.Equal(t, original, result)
	})
}

// ---------------------------------------------------------------------------
// Examples
// ---------------------------------------------------------------------------

// ExampleHead_found demonstrates extracting the first element from a non-empty slice.
func ExampleHead_found() {
	p := Head[int]()
	fmt.Println(p.GetOption([]int{10, 20, 30}))
	// Output:
	// Some[int](10)
}

// ExampleHead_empty demonstrates that Head returns None for an empty slice.
func ExampleHead_empty() {
	p := Head[int]()
	fmt.Println(p.GetOption([]int{}))
	// Output:
	// None[int]
}

// ExampleHead_reverseGet demonstrates that ReverseGet wraps a value in a singleton slice.
func ExampleHead_reverseGet() {
	p := Head[string]()
	fmt.Println(p.ReverseGet("hello"))
	// Output:
	// [hello]
}

// ExampleAt_found demonstrates retrieving the element at a given index.
func ExampleAt_found() {
	p := At[string]()(2)
	fmt.Println(p.GetOption([]string{"a", "b", "c", "d"}))
	// Output:
	// Some[string](c)
}

// ExampleAt_outOfBounds demonstrates that At returns None when the index exceeds the slice length.
func ExampleAt_outOfBounds() {
	p := At[int]()(10)
	fmt.Println(p.GetOption([]int{1, 2, 3}))
	// Output:
	// None[int]
}

// ExampleAt_negative demonstrates that a negative index always yields None.
func ExampleAt_negative() {
	p := At[int]()(-1)
	fmt.Println(p.GetOption([]int{1, 2, 3}))
	// Output:
	// None[int]
}

// ExampleAt_reverseGet demonstrates that ReverseGet produces a slice with the value at position i.
func ExampleAt_reverseGet() {
	p := At[int]()(2)
	fmt.Println(p.ReverseGet(99))
	// Output:
	// [0 0 99]
}
