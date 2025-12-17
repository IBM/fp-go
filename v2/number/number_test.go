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

package number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test MagmaSub
func TestMagmaSub(t *testing.T) {
	subMagma := MagmaSub[int]()

	tests := []struct {
		name     string
		first    int
		second   int
		expected int
	}{
		{"positive numbers", 10, 3, 7},
		{"negative result", 3, 10, -7},
		{"with zero", 5, 0, 5},
		{"zero minus number", 0, 5, -5},
		{"negative numbers", -5, -3, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := subMagma.Concat(tt.first, tt.second)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test MagmaSub with floats
func TestMagmaSub_Float(t *testing.T) {
	subMagma := MagmaSub[float64]()

	result := subMagma.Concat(10.5, 3.2)
	assert.InDelta(t, 7.3, result, 0.0001)

	result = subMagma.Concat(3.2, 10.5)
	assert.InDelta(t, -7.3, result, 0.0001)
}

// Test MagmaDiv
func TestMagmaDiv(t *testing.T) {
	divMagma := MagmaDiv[int]()

	tests := []struct {
		name     string
		first    int
		second   int
		expected int
	}{
		{"simple division", 10, 2, 5},
		{"division with remainder", 10, 3, 3},
		{"one divided by itself", 5, 5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := divMagma.Concat(tt.first, tt.second)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test MagmaDiv with floats
func TestMagmaDiv_Float(t *testing.T) {
	divMagma := MagmaDiv[float64]()

	result := divMagma.Concat(10.0, 2.0)
	assert.Equal(t, 5.0, result)

	result = divMagma.Concat(10.0, 3.0)
	assert.InDelta(t, 3.333333, result, 0.0001)

	result = divMagma.Concat(1.0, 2.0)
	assert.Equal(t, 0.5, result)
}

// Test SemigroupSum
func TestSemigroupSum(t *testing.T) {
	sumSemigroup := SemigroupSum[int]()

	tests := []struct {
		name     string
		first    int
		second   int
		expected int
	}{
		{"positive numbers", 5, 3, 8},
		{"with zero", 5, 0, 5},
		{"negative numbers", -5, -3, -8},
		{"mixed signs", 10, -3, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sumSemigroup.Concat(tt.first, tt.second)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test associativity
	a, b, c := 1, 2, 3
	assert.Equal(t,
		sumSemigroup.Concat(sumSemigroup.Concat(a, b), c),
		sumSemigroup.Concat(a, sumSemigroup.Concat(b, c)),
	)
}

// Test SemigroupSum with floats
func TestSemigroupSum_Float(t *testing.T) {
	sumSemigroup := SemigroupSum[float64]()

	result := sumSemigroup.Concat(3.14, 2.86)
	assert.InDelta(t, 6.0, result, 0.0001)

	result = sumSemigroup.Concat(-1.5, 2.5)
	assert.Equal(t, 1.0, result)
}

// Test SemigroupSum with complex numbers
func TestSemigroupSum_Complex(t *testing.T) {
	sumSemigroup := SemigroupSum[complex128]()

	c1 := complex(1, 2)
	c2 := complex(3, 4)
	result := sumSemigroup.Concat(c1, c2)
	expected := complex(4, 6)
	assert.Equal(t, expected, result)
}

// Test SemigroupProduct
func TestSemigroupProduct(t *testing.T) {
	prodSemigroup := SemigroupProduct[int]()

	tests := []struct {
		name     string
		first    int
		second   int
		expected int
	}{
		{"positive numbers", 5, 3, 15},
		{"with one", 5, 1, 5},
		{"with zero", 5, 0, 0},
		{"negative numbers", -5, -3, 15},
		{"mixed signs", 5, -3, -15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prodSemigroup.Concat(tt.first, tt.second)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test associativity
	a, b, c := 2, 3, 4
	assert.Equal(t,
		prodSemigroup.Concat(prodSemigroup.Concat(a, b), c),
		prodSemigroup.Concat(a, prodSemigroup.Concat(b, c)),
	)
}

// Test SemigroupProduct with floats
func TestSemigroupProduct_Float(t *testing.T) {
	prodSemigroup := SemigroupProduct[float64]()

	result := prodSemigroup.Concat(2.5, 4.0)
	assert.Equal(t, 10.0, result)

	result = prodSemigroup.Concat(0.5, 10.0)
	assert.Equal(t, 5.0, result)
}

// Test MonoidSum
func TestMonoidSum(t *testing.T) {
	sumMonoid := MonoidSum[int]()

	// Test concat
	assert.Equal(t, 8, sumMonoid.Concat(5, 3))
	assert.Equal(t, 0, sumMonoid.Concat(5, -5))

	// Test empty
	assert.Equal(t, 0, sumMonoid.Empty())

	// Test identity laws
	testValues := []int{0, 1, -1, 5, 10, -10, 100}
	for _, x := range testValues {
		// Left identity: Empty() + x = x
		assert.Equal(t, x, sumMonoid.Concat(sumMonoid.Empty(), x),
			"Left identity failed for %d", x)

		// Right identity: x + Empty() = x
		assert.Equal(t, x, sumMonoid.Concat(x, sumMonoid.Empty()),
			"Right identity failed for %d", x)
	}

	// Test associativity
	assert.Equal(t,
		sumMonoid.Concat(sumMonoid.Concat(1, 2), 3),
		sumMonoid.Concat(1, sumMonoid.Concat(2, 3)),
	)
}

// Test MonoidSum with floats
func TestMonoidSum_Float(t *testing.T) {
	sumMonoid := MonoidSum[float64]()

	assert.InDelta(t, 6.0, sumMonoid.Concat(3.14, 2.86), 0.0001)
	assert.Equal(t, 0.0, sumMonoid.Empty())

	// Test identity
	x := 5.5
	assert.Equal(t, x, sumMonoid.Concat(sumMonoid.Empty(), x))
	assert.Equal(t, x, sumMonoid.Concat(x, sumMonoid.Empty()))
}

// Test MonoidProduct
func TestMonoidProduct(t *testing.T) {
	prodMonoid := MonoidProduct[int]()

	// Test concat
	assert.Equal(t, 15, prodMonoid.Concat(5, 3))
	assert.Equal(t, 0, prodMonoid.Concat(5, 0))

	// Test empty
	assert.Equal(t, 1, prodMonoid.Empty())

	// Test identity laws
	testValues := []int{1, 2, 3, 5, 10}
	for _, x := range testValues {
		// Left identity: Empty() * x = x
		assert.Equal(t, x, prodMonoid.Concat(prodMonoid.Empty(), x),
			"Left identity failed for %d", x)

		// Right identity: x * Empty() = x
		assert.Equal(t, x, prodMonoid.Concat(x, prodMonoid.Empty()),
			"Right identity failed for %d", x)
	}

	// Test associativity
	assert.Equal(t,
		prodMonoid.Concat(prodMonoid.Concat(2, 3), 4),
		prodMonoid.Concat(2, prodMonoid.Concat(3, 4)),
	)
}

// Test MonoidProduct with floats
func TestMonoidProduct_Float(t *testing.T) {
	prodMonoid := MonoidProduct[float64]()

	assert.Equal(t, 10.0, prodMonoid.Concat(2.5, 4.0))
	assert.Equal(t, 1.0, prodMonoid.Empty())

	// Test identity
	x := 5.5
	assert.Equal(t, x, prodMonoid.Concat(prodMonoid.Empty(), x))
	assert.Equal(t, x, prodMonoid.Concat(x, prodMonoid.Empty()))
}

// Test Add curried function
func TestAdd(t *testing.T) {
	add5 := Add(5)

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive", 10, 15},
		{"zero", 0, 5},
		{"negative", -3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := add5(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Add with floats
func TestAdd_Float(t *testing.T) {
	add2_5 := Add(2.5)

	assert.Equal(t, 7.5, add2_5(5.0))
	assert.Equal(t, 2.5, add2_5(0.0))
	assert.InDelta(t, 5.64, add2_5(3.14), 0.0001)
}

// Test Sub curried function
func TestSub(t *testing.T) {
	sub3 := Sub(3)

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive result", 10, 7},
		{"zero result", 3, 0},
		{"negative result", 1, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sub3(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Sub with floats
func TestSub_Float(t *testing.T) {
	sub2_5 := Sub(2.5)

	assert.Equal(t, 2.5, sub2_5(5.0))
	assert.Equal(t, -2.5, sub2_5(0.0))
	assert.InDelta(t, 0.64, sub2_5(3.14), 0.0001)
}

// Test Mul curried function
func TestMul(t *testing.T) {
	double := Mul(2)

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive", 5, 10},
		{"zero", 0, 0},
		{"negative", -3, -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := double(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Mul with floats
func TestMul_Float(t *testing.T) {
	triple := Mul(3.0)

	assert.Equal(t, 15.0, triple(5.0))
	assert.Equal(t, 0.0, triple(0.0))
	assert.InDelta(t, 9.42, triple(3.14), 0.0001)
}

// Test Div curried function
func TestDiv(t *testing.T) {
	divBy2 := Div(2)

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"even number", 10, 5},
		{"odd number", 9, 4},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := divBy2(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Div with floats
func TestDiv_Float(t *testing.T) {
	half := Div(2.0)

	assert.Equal(t, 5.0, half(10.0))
	assert.Equal(t, 2.5, half(5.0))
	assert.InDelta(t, 1.57, half(3.14), 0.0001)
}

// Test Inc function
func TestInc(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive", 5, 6},
		{"zero", 0, 1},
		{"negative", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Inc(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Inc with floats
func TestInc_Float(t *testing.T) {
	assert.Equal(t, 6.5, Inc(5.5))
	assert.Equal(t, 1.0, Inc(0.0))
	assert.InDelta(t, 4.14, Inc(3.14), 0.0001)
}

// Test Min function
func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a smaller", 3, 5, 3},
		{"b smaller", 5, 3, 3},
		{"equal", 5, 5, 5},
		{"negative", -5, -3, -5},
		{"mixed signs", -5, 3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Min with floats
func TestMin_Float(t *testing.T) {
	assert.Equal(t, 2.5, Min(2.5, 7.8))
	assert.Equal(t, 2.5, Min(7.8, 2.5))
	assert.Equal(t, 5.5, Min(5.5, 5.5))
	assert.Equal(t, -3.14, Min(-3.14, 2.71))
}

// Test Max function
func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a larger", 5, 3, 5},
		{"b larger", 3, 5, 5},
		{"equal", 5, 5, 5},
		{"negative", -5, -3, -3},
		{"mixed signs", -5, 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Max with floats
func TestMax_Float(t *testing.T) {
	assert.Equal(t, 7.8, Max(2.5, 7.8))
	assert.Equal(t, 7.8, Max(7.8, 2.5))
	assert.Equal(t, 5.5, Max(5.5, 5.5))
	assert.Equal(t, 2.71, Max(-3.14, 2.71))
}

// Benchmark tests
func BenchmarkMonoidSum(b *testing.B) {
	sumMonoid := MonoidSum[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sumMonoid.Concat(i, i+1)
	}
}

func BenchmarkMonoidProduct(b *testing.B) {
	prodMonoid := MonoidProduct[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prodMonoid.Concat(i+1, i+2)
	}
}

func BenchmarkAdd(b *testing.B) {
	add5 := Add(5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = add5(i)
	}
}

func BenchmarkMin(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Min(i, i+1)
	}
}

func BenchmarkMax(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Max(i, i+1)
	}
}

// Test MoreThan curried function
func TestMoreThan(t *testing.T) {
	moreThan10 := MoreThan(10)

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"greater than threshold", 15, true},
		{"less than threshold", 5, false},
		{"equal to threshold", 10, false},
		{"much greater", 100, true},
		{"negative value", -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := moreThan10(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test MoreThan with floats
func TestMoreThan_Float(t *testing.T) {
	moreThan5_5 := MoreThan(5.5)

	assert.True(t, moreThan5_5(6.0))
	assert.False(t, moreThan5_5(5.0))
	assert.False(t, moreThan5_5(5.5))
	assert.True(t, moreThan5_5(10.5))
	assert.False(t, moreThan5_5(5.4))
}

// Test MoreThan with negative numbers
func TestMoreThan_Negative(t *testing.T) {
	moreThanNeg5 := MoreThan(-5)

	assert.True(t, moreThanNeg5(0))
	assert.True(t, moreThanNeg5(-4))
	assert.False(t, moreThanNeg5(-5))
	assert.False(t, moreThanNeg5(-10))
}

// Test LessThan curried function
func TestLessThan(t *testing.T) {
	lessThan10 := LessThan(10)

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"less than threshold", 5, true},
		{"greater than threshold", 15, false},
		{"equal to threshold", 10, false},
		{"much less", -10, true},
		{"zero", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lessThan10(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test LessThan with floats
func TestLessThan_Float(t *testing.T) {
	lessThan5_5 := LessThan(5.5)

	assert.True(t, lessThan5_5(5.0))
	assert.False(t, lessThan5_5(6.0))
	assert.False(t, lessThan5_5(5.5))
	assert.True(t, lessThan5_5(2.5))
	assert.True(t, lessThan5_5(5.4))
}

// Test LessThan with negative numbers
func TestLessThan_Negative(t *testing.T) {
	lessThanNeg5 := LessThan(-5)

	assert.False(t, lessThanNeg5(0))
	assert.False(t, lessThanNeg5(-4))
	assert.False(t, lessThanNeg5(-5))
	assert.True(t, lessThanNeg5(-10))
}

// Test MoreThan and LessThan together for range checking
func TestMoreThanLessThan_Range(t *testing.T) {
	// Check if value is in range (10, 20) - exclusive
	moreThan10 := MoreThan(10)
	lessThan20 := LessThan(20)

	inRange := func(x int) bool {
		return moreThan10(x) && lessThan20(x)
	}

	assert.True(t, inRange(15))
	assert.False(t, inRange(10))
	assert.False(t, inRange(20))
	assert.False(t, inRange(5))
	assert.False(t, inRange(25))
}

// Benchmark tests for comparison functions
func BenchmarkMoreThan(b *testing.B) {
	moreThan10 := MoreThan(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = moreThan10(i)
	}
}

func BenchmarkLessThan(b *testing.B) {
	lessThan10 := LessThan(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lessThan10(i)
	}
}
