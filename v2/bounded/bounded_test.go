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

package bounded

import (
	"testing"

	"github.com/IBM/fp-go/v2/ord"
	"github.com/stretchr/testify/assert"
)

func TestMakeBounded(t *testing.T) {
	t.Run("creates bounded instance with correct top and bottom", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)

		assert.Equal(t, 100, b.Top())
		assert.Equal(t, 0, b.Bottom())
	})

	t.Run("preserves ordering from Ord", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)

		assert.Equal(t, -1, b.Compare(5, 10))
		assert.Equal(t, 0, b.Compare(5, 5))
		assert.Equal(t, 1, b.Compare(10, 5))
	})

	t.Run("preserves equality from Ord", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)

		assert.True(t, b.Equals(5, 5))
		assert.False(t, b.Equals(5, 10))
	})
}

func TestClamp(t *testing.T) {
	t.Run("returns value within bounds unchanged", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)
		clamp := Clamp(b)

		assert.Equal(t, 50, clamp(50))
		assert.Equal(t, 0, clamp(0))
		assert.Equal(t, 100, clamp(100))
	})

	t.Run("clamps value above top to top", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)
		clamp := Clamp(b)

		assert.Equal(t, 100, clamp(150))
		assert.Equal(t, 100, clamp(200))
	})

	t.Run("clamps value below bottom to bottom", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)
		clamp := Clamp(b)

		assert.Equal(t, 0, clamp(-10))
		assert.Equal(t, 0, clamp(-100))
	})

	t.Run("works with float64", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[float64](), 1.0, 0.0)
		clamp := Clamp(b)

		assert.Equal(t, 0.5, clamp(0.5))
		assert.Equal(t, 1.0, clamp(1.5))
		assert.Equal(t, 0.0, clamp(-0.5))
	})

	t.Run("works with strings", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[string](), "z", "a")
		clamp := Clamp(b)

		assert.Equal(t, "m", clamp("m"))
		assert.Equal(t, "z", clamp("zzz"))
		assert.Equal(t, "a", clamp("A"))
	})
}

func TestReverse(t *testing.T) {
	t.Run("reverses the ordering", func(t *testing.T) {
		original := MakeBounded(ord.FromStrictCompare[int](), 100, 0)
		reversed := Reverse(original)

		// In original: 5 < 10, so Compare(5, 10) = -1
		assert.Equal(t, -1, original.Compare(5, 10))

		// In reversed: 5 > 10, so Compare(5, 10) = 1
		assert.Equal(t, 1, reversed.Compare(5, 10))
	})

	t.Run("swaps top and bottom values", func(t *testing.T) {
		original := MakeBounded(ord.FromStrictCompare[int](), 100, 0)
		reversed := Reverse(original)

		// Reverse swaps the bounds
		assert.Equal(t, original.Bottom(), reversed.Top())
		assert.Equal(t, original.Top(), reversed.Bottom())
	})

	t.Run("double reverse returns to original ordering", func(t *testing.T) {
		original := MakeBounded(ord.FromStrictCompare[int](), 100, 0)
		reversed := Reverse(original)
		doubleReversed := Reverse(reversed)

		assert.Equal(t, original.Compare(5, 10), doubleReversed.Compare(5, 10))
		assert.Equal(t, original.Compare(10, 5), doubleReversed.Compare(10, 5))
	})

	t.Run("preserves equality", func(t *testing.T) {
		original := MakeBounded(ord.FromStrictCompare[int](), 100, 0)
		reversed := Reverse(original)

		assert.Equal(t, original.Equals(5, 5), reversed.Equals(5, 5))
		assert.Equal(t, original.Equals(5, 10), reversed.Equals(5, 10))
	})
}

func TestBoundedLaws(t *testing.T) {
	t.Run("bottom is less than or equal to all values", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)

		testValues := []int{0, 25, 50, 75, 100}
		for _, v := range testValues {
			assert.True(t, b.Compare(b.Bottom(), v) <= 0,
				"Bottom (%d) should be <= %d", b.Bottom(), v)
		}
	})

	t.Run("top is greater than or equal to all values", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)

		testValues := []int{0, 25, 50, 75, 100}
		for _, v := range testValues {
			assert.True(t, b.Compare(b.Top(), v) >= 0,
				"Top (%d) should be >= %d", b.Top(), v)
		}
	})

	t.Run("bottom is less than or equal to top", func(t *testing.T) {
		b := MakeBounded(ord.FromStrictCompare[int](), 100, 0)

		assert.True(t, b.Compare(b.Bottom(), b.Top()) <= 0,
			"Bottom should be <= Top")
	})
}

// Example tests
func ExampleMakeBounded() {
	// Create a bounded type for percentages (0-100)
	percentage := MakeBounded(
		ord.FromStrictCompare[int](),
		100, // top
		0,   // bottom
	)

	println(percentage.Top())    // 100
	println(percentage.Bottom()) // 0

	// Output:
}

func ExampleClamp() {
	// Create bounded type for percentages
	percentage := MakeBounded(
		ord.FromStrictCompare[int](),
		100, // top
		0,   // bottom
	)

	clamp := Clamp(percentage)

	println(clamp(50))  // 50 (within bounds)
	println(clamp(150)) // 100 (clamped to top)
	println(clamp(-10)) // 0 (clamped to bottom)

	// Output:
}

func ExampleReverse() {
	original := MakeBounded(
		ord.FromStrictCompare[int](),
		100, // top
		0,   // bottom
	)

	reversed := Reverse(original)

	// Ordering is reversed
	println(original.Compare(5, 10)) // -1 (5 < 10)
	println(reversed.Compare(5, 10)) // 1 (5 > 10 in reversed)

	// Bounds are swapped
	println(reversed.Top())    // 0
	println(reversed.Bottom()) // 100

	// Output:
}
