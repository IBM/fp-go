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

package array

import (
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

// TestSliceBasicCases tests normal slicing operations
func TestSliceBasicCases(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5}

	t.Run("normal slice from middle", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3}, Slice[int](1, 4)(data))
	})

	t.Run("slice from start", func(t *testing.T) {
		assert.Equal(t, []int{0, 1, 2}, Slice[int](0, 3)(data))
	})

	t.Run("slice to end", func(t *testing.T) {
		assert.Equal(t, []int{3, 4, 5}, Slice[int](3, 6)(data))
	})

	t.Run("slice single element", func(t *testing.T) {
		assert.Equal(t, []int{2}, Slice[int](2, 3)(data))
	})

	t.Run("slice entire array", func(t *testing.T) {
		assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, Slice[int](0, 6)(data))
	})
}

// TestSliceNegativeIndices tests negative index handling (counting from end)
func TestSliceNegativeIndices(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5}

	t.Run("negative start index", func(t *testing.T) {
		// -2 means length + (-2) = 6 - 2 = 4
		assert.Equal(t, []int{4, 5}, Slice[int](-2, 6)(data))
	})

	t.Run("negative end index", func(t *testing.T) {
		// -2 means length + (-2) = 6 - 2 = 4
		assert.Equal(t, []int{0, 1, 2, 3}, Slice[int](0, -2)(data))
	})

	t.Run("both negative indices", func(t *testing.T) {
		// -4 = 2, -2 = 4
		assert.Equal(t, []int{2, 3}, Slice[int](-4, -2)(data))
	})

	t.Run("negative index beyond array start", func(t *testing.T) {
		// -10 would be -4, clamped to 0
		assert.Equal(t, []int{0, 1, 2}, Slice[int](-10, 3)(data))
	})

	t.Run("negative end index beyond array start", func(t *testing.T) {
		// -10 would be -4, clamped to 0
		assert.Equal(t, []int{}, Slice[int](0, -10)(data))
	})
}

// TestSliceEmptyArray tests slicing on empty arrays (totality proof)
func TestSliceEmptyArray(t *testing.T) {
	empty := []int{}

	t.Run("slice empty array with zero indices", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](0, 0)(empty))
	})

	t.Run("slice empty array with positive indices", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](0, 5)(empty))
	})

	t.Run("slice empty array with negative indices", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](-1, -1)(empty))
	})

	t.Run("slice empty array with mixed indices", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](-5, 5)(empty))
	})
}

// TestSliceOutOfBounds tests out-of-bounds scenarios (totality proof)
func TestSliceOutOfBounds(t *testing.T) {
	data := []int{0, 1, 2, 3, 4}

	t.Run("start index beyond array length", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](10, 15)(data))
	})

	t.Run("end index beyond array length", func(t *testing.T) {
		assert.Equal(t, []int{2, 3, 4}, Slice[int](2, 100)(data))
	})

	t.Run("both indices beyond array length", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](10, 20)(data))
	})

	t.Run("start equals array length", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](5, 10)(data))
	})

	t.Run("end equals array length", func(t *testing.T) {
		assert.Equal(t, []int{3, 4}, Slice[int](3, 5)(data))
	})
}

// TestSliceInvalidRanges tests invalid range scenarios (totality proof)
func TestSliceInvalidRanges(t *testing.T) {
	data := []int{0, 1, 2, 3, 4}

	t.Run("start equals end", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](2, 2)(data))
	})

	t.Run("start greater than end", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](4, 2)(data))
	})

	t.Run("start greater than end with negative indices", func(t *testing.T) {
		// -1 = 4, -3 = 2
		assert.Equal(t, []int{}, Slice[int](-1, -3)(data))
	})

	t.Run("zero range at start", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](0, 0)(data))
	})

	t.Run("zero range at end", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](5, 5)(data))
	})
}

// TestSliceEdgeCases tests additional edge cases (totality proof)
func TestSliceEdgeCases(t *testing.T) {
	t.Run("single element array - slice all", func(t *testing.T) {
		data := []int{42}
		assert.Equal(t, []int{42}, Slice[int](0, 1)(data))
	})

	t.Run("single element array - slice none", func(t *testing.T) {
		data := []int{42}
		assert.Equal(t, []int{}, Slice[int](1, 1)(data))
	})

	t.Run("single element array - negative indices", func(t *testing.T) {
		data := []int{42}
		assert.Equal(t, []int{42}, Slice[int](-1, 1)(data))
	})

	t.Run("large array slice", func(t *testing.T) {
		data := MakeBy(1000, func(i int) int { return i })
		result := Slice[int](100, 200)(data)
		assert.Equal(t, 100, len(result))
		assert.Equal(t, 100, result[0])
		assert.Equal(t, 199, result[99])
	})
}

// TestSliceWithDifferentTypes tests that Slice works with different types (totality proof)
func TestSliceWithDifferentTypes(t *testing.T) {
	t.Run("string slice", func(t *testing.T) {
		data := []string{"a", "b", "c", "d", "e"}
		assert.Equal(t, []string{"b", "c", "d"}, Slice[string](1, 4)(data))
	})

	t.Run("float slice", func(t *testing.T) {
		data := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
		assert.Equal(t, []float64{2.2, 3.3}, Slice[float64](1, 3)(data))
	})

	t.Run("bool slice", func(t *testing.T) {
		data := []bool{true, false, true, false}
		assert.Equal(t, []bool{false, true}, Slice[bool](1, 3)(data))
	})

	t.Run("struct slice", func(t *testing.T) {
		type Point struct{ X, Y int }
		data := []Point{{1, 2}, {3, 4}, {5, 6}}
		assert.Equal(t, []Point{{3, 4}}, Slice[Point](1, 2)(data))
	})

	t.Run("pointer slice", func(t *testing.T) {
		a, b, c := 1, 2, 3
		data := []*int{&a, &b, &c}
		result := Slice[*int](1, 3)(data)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, 2, *result[0])
		assert.Equal(t, 3, *result[1])
	})
}

// TestSliceNilArray tests behavior with nil arrays (totality proof)
func TestSliceNilArray(t *testing.T) {
	var nilArray []int

	t.Run("slice nil array with zero indices", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](0, 0)(nilArray))
	})

	t.Run("slice nil array with positive indices", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](0, 5)(nilArray))
	})

	t.Run("slice nil array with negative indices", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](-1, 1)(nilArray))
	})

	t.Run("slice nil array with out of bounds indices", func(t *testing.T) {
		assert.Equal(t, []int{}, Slice[int](10, 20)(nilArray))
	})
}

// TestSliceComposition tests that Slice can be composed with other functions
func TestSliceComposition(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	t.Run("compose multiple slices", func(t *testing.T) {
		// First slice [2:8], then slice [1:4] of result
		slice1 := Slice[int](2, 8)
		slice2 := Slice[int](1, 4)
		result := slice2(slice1(data))
		// [2,3,4,5,6,7] -> [3,4,5]
		assert.Equal(t, []int{3, 4, 5}, result)
	})

	t.Run("slice then map", func(t *testing.T) {
		sliced := Slice[int](2, 5)(data)
		mapped := Map(N.Mul(2))(sliced)
		assert.Equal(t, []int{4, 6, 8}, mapped)
	})

	t.Run("slice then filter", func(t *testing.T) {
		sliced := Slice[int](0, 6)(data)
		filtered := Filter(func(x int) bool { return x%2 == 0 })(sliced)
		assert.Equal(t, []int{0, 2, 4}, filtered)
	})
}

// TestSliceImmutability tests that Slice doesn't modify the original array
func TestSliceImmutability(t *testing.T) {
	original := []int{0, 1, 2, 3, 4}
	originalCopy := []int{0, 1, 2, 3, 4}

	t.Run("slicing doesn't modify original", func(t *testing.T) {
		result := Slice[int](1, 4)(original)
		assert.Equal(t, []int{1, 2, 3}, result)
		assert.Equal(t, originalCopy, original)
	})

	t.Run("slice shares underlying array with original", func(t *testing.T) {
		// Note: Go's slice operation creates a view of the underlying array,
		// not a deep copy. This is expected behavior and matches Go's built-in slice semantics.
		result := Slice[int](1, 4)(original)
		result[0] = 999
		// The original array is affected because slices share the underlying array
		assert.Equal(t, 999, original[1], "Slices share underlying array (expected Go behavior)")
	})
}

// TestSliceTotality is a comprehensive test proving Slice is a total function
// A total function is defined for all possible inputs and never panics
func TestSliceTotality(t *testing.T) {
	testCases := []struct {
		name  string
		data  []int
		low   int
		high  int
		panic bool // Should always be false for a total function
	}{
		// Normal cases
		{"normal range", []int{1, 2, 3, 4, 5}, 1, 3, false},
		{"full range", []int{1, 2, 3}, 0, 3, false},
		{"empty result", []int{1, 2, 3}, 1, 1, false},

		// Edge cases with empty/nil arrays
		{"empty array", []int{}, 0, 0, false},
		{"empty array with indices", []int{}, 1, 5, false},
		{"nil array", nil, 0, 5, false},

		// Negative indices
		{"negative low", []int{1, 2, 3, 4, 5}, -2, 5, false},
		{"negative high", []int{1, 2, 3, 4, 5}, 0, -1, false},
		{"both negative", []int{1, 2, 3, 4, 5}, -3, -1, false},
		{"negative beyond bounds", []int{1, 2, 3}, -100, -50, false},

		// Out of bounds
		{"low beyond length", []int{1, 2, 3}, 10, 20, false},
		{"high beyond length", []int{1, 2, 3}, 1, 100, false},
		{"both beyond length", []int{1, 2, 3}, 10, 20, false},

		// Invalid ranges
		{"low equals high", []int{1, 2, 3}, 2, 2, false},
		{"low greater than high", []int{1, 2, 3}, 3, 1, false},
		{"negative invalid range", []int{1, 2, 3, 4, 5}, -1, -3, false},

		// Extreme values
		{"very large indices", []int{1, 2, 3}, 1000000, 2000000, false},
		{"very negative indices", []int{1, 2, 3}, -1000000, -500000, false},
		{"mixed extreme", []int{1, 2, 3}, -1000000, 1000000, false},

		// Zero values
		{"zero indices", []int{1, 2, 3}, 0, 0, false},
		{"zero low", []int{1, 2, 3}, 0, 3, false},
		{"zero high", []int{1, 2, 3}, 0, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This test proves totality by ensuring no panic occurs
			defer func() {
				if r := recover(); r != nil {
					if !tc.panic {
						t.Errorf("Slice panicked unexpectedly: %v", r)
					}
				} else {
					if tc.panic {
						t.Errorf("Slice should have panicked but didn't")
					}
				}
			}()

			// Execute the function - if it's total, it will never panic
			result := Slice[int](tc.low, tc.high)(tc.data)

			// Additional verification: result should always be a valid slice
			assert.NotNil(t, result, "Result should never be nil")
			assert.True(t, len(result) >= 0, "Result length should be non-negative")
		})
	}
}

// TestSlicePropertyBased tests mathematical properties of Slice
func TestSlicePropertyBased(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	t.Run("identity: Slice(0, len) returns copy of array", func(t *testing.T) {
		result := Slice[int](0, len(data))(data)
		assert.Equal(t, data, result)
	})

	t.Run("empty: Slice(i, i) always returns empty", func(t *testing.T) {
		for i := 0; i <= len(data); i++ {
			result := Slice[int](i, i)(data)
			assert.Equal(t, []int{}, result)
		}
	})

	t.Run("length property: len(Slice(i, j)) = max(0, min(j, len) - max(i, 0))", func(t *testing.T) {
		testCases := []struct{ low, high, expected int }{
			{0, 5, 5},
			{2, 7, 5},
			{5, 5, 0},
			{3, 2, 0},   // invalid range
			{-2, 10, 2}, // -2 becomes 8, so slice [8:10] has length 2
			{0, 100, 10},
		}

		for _, tc := range testCases {
			result := Slice[int](tc.low, tc.high)(data)
			assert.Equal(t, tc.expected, len(result),
				"Slice(%d, %d) should have length %d", tc.low, tc.high, tc.expected)
		}
	})

	t.Run("concatenation: Slice(0,i) + Slice(i,len) = original", func(t *testing.T) {
		for i := 0; i <= len(data); i++ {
			left := Slice[int](0, i)(data)
			right := Slice[int](i, len(data))(data)
			concatenated := append(left, right...)
			assert.Equal(t, data, concatenated)
		}
	})

	t.Run("subset property: all elements in slice are in original", func(t *testing.T) {
		result := Slice[int](2, 7)(data)
		for _, elem := range result {
			found := false
			for _, orig := range data {
				if elem == orig {
					found = true
					break
				}
			}
			assert.True(t, found, "Element %d should be in original array", elem)
		}
	})
}
