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

package semigroup

import (
	"testing"

	M "github.com/IBM/fp-go/v2/magma"
	"github.com/stretchr/testify/assert"
)

// Test basic First semigroup
func TestFirst(t *testing.T) {
	first := First[int]()
	assert.Equal(t, 1, first.Concat(1, 2))
	assert.Equal(t, 10, first.Concat(10, 20))
	assert.Equal(t, "a", First[string]().Concat("a", "b"))
}

// Test basic Last semigroup
func TestLast(t *testing.T) {
	last := Last[int]()
	assert.Equal(t, 2, last.Concat(1, 2))
	assert.Equal(t, 20, last.Concat(10, 20))
	assert.Equal(t, "b", Last[string]().Concat("a", "b"))
}

// Test MakeSemigroup
func TestMakeSemigroup(t *testing.T) {
	// Integer addition semigroup
	add := MakeSemigroup(func(a, b int) int { return a + b })
	assert.Equal(t, 5, add.Concat(2, 3))
	assert.Equal(t, 10, add.Concat(4, 6))

	// String concatenation semigroup
	concat := MakeSemigroup(func(a, b string) string { return a + b })
	assert.Equal(t, "hello", concat.Concat("hel", "lo"))
	assert.Equal(t, "foobar", concat.Concat("foo", "bar"))

	// Max semigroup
	max := MakeSemigroup(func(a, b int) int {
		if a > b {
			return a
		}
		return b
	})
	assert.Equal(t, 10, max.Concat(5, 10))
	assert.Equal(t, 20, max.Concat(20, 15))
}

// Test Reverse semigroup
func TestReverse(t *testing.T) {
	// Subtraction is not commutative, so reverse changes the result
	sub := MakeSemigroup(func(a, b int) int { return a - b })
	reversed := Reverse(sub)

	assert.Equal(t, 7, sub.Concat(10, 3))       // 10 - 3 = 7
	assert.Equal(t, -7, reversed.Concat(10, 3)) // 3 - 10 = -7

	// String concatenation
	concat := MakeSemigroup(func(a, b string) string { return a + b })
	reversedConcat := Reverse(concat)

	assert.Equal(t, "ab", concat.Concat("a", "b"))
	assert.Equal(t, "ba", reversedConcat.Concat("a", "b"))
}

// Test FunctionSemigroup
func TestFunctionSemigroup(t *testing.T) {
	// Base semigroup for integers (addition)
	add := MakeSemigroup(func(a, b int) int { return a + b })

	// Lift to functions
	funcSG := FunctionSemigroup[string](add)

	// Create two functions
	f := func(s string) int { return len(s) }
	g := func(s string) int { return len(s) * 2 }

	// Combine functions
	combined := funcSG.Concat(f, g)

	// Test with different strings
	assert.Equal(t, 15, combined("hello")) // 5 + 10 = 15
	assert.Equal(t, 9, combined("abc"))    // 3 + 6 = 9
	assert.Equal(t, 0, combined(""))       // 0 + 0 = 0
}

// Test FunctionSemigroup with different types
func TestFunctionSemigroupMultipleTypes(t *testing.T) {
	// String concatenation semigroup
	concat := MakeSemigroup(func(a, b string) string { return a + b })

	// Lift to functions from int to string
	funcSG := FunctionSemigroup[int](concat)

	f := func(n int) string { return "a" }
	g := func(n int) string { return "b" }

	combined := funcSG.Concat(f, g)
	assert.Equal(t, "ab", combined(42))
}

// Test ToMagma conversion
func TestToMagma(t *testing.T) {
	sg := MakeSemigroup(func(a, b int) int { return a + b })
	magma := ToMagma(sg)

	// Should work as a magma
	assert.Equal(t, 5, magma.Concat(2, 3))
	assert.Equal(t, 10, magma.Concat(4, 6))

	// Verify it's a Magma interface
	var _ M.Magma[int] = magma
}

// Test ConcatAll
func TestConcatAll(t *testing.T) {
	add := MakeSemigroup(func(a, b int) int { return a + b })
	concatAll := ConcatAll(add)

	// Test with various arrays
	assert.Equal(t, 10, concatAll(0)([]int{1, 2, 3, 4}))
	assert.Equal(t, 20, concatAll(10)([]int{1, 2, 3, 4}))
	assert.Equal(t, 5, concatAll(5)([]int{}))
	assert.Equal(t, 15, concatAll(0)([]int{15}))

	// Test with string concatenation
	concat := MakeSemigroup(func(a, b string) string { return a + b })
	concatAllStr := ConcatAll(concat)

	assert.Equal(t, "hello", concatAllStr("")([]string{"h", "e", "l", "l", "o"}))
	assert.Equal(t, "prefix_abc", concatAllStr("prefix_")([]string{"a", "b", "c"}))
}

// Test MonadConcatAll
func TestMonadConcatAll(t *testing.T) {
	add := MakeSemigroup(func(a, b int) int { return a + b })
	monadConcatAll := MonadConcatAll(add)

	// Test with various arrays
	assert.Equal(t, 10, monadConcatAll([]int{1, 2, 3, 4}, 0))
	assert.Equal(t, 20, monadConcatAll([]int{1, 2, 3, 4}, 10))
	assert.Equal(t, 5, monadConcatAll([]int{}, 5))
	assert.Equal(t, 15, monadConcatAll([]int{15}, 0))

	// Test with multiplication
	mul := MakeSemigroup(func(a, b int) int { return a * b })
	monadConcatAllMul := MonadConcatAll(mul)

	assert.Equal(t, 24, monadConcatAllMul([]int{2, 3, 4}, 1))
	assert.Equal(t, 120, monadConcatAllMul([]int{2, 3, 4, 5}, 1))
}

// Test GenericConcatAll with custom slice type
func TestGenericConcatAll(t *testing.T) {
	type MyInts []int

	add := MakeSemigroup(func(a, b int) int { return a + b })
	concatAll := GenericConcatAll[MyInts](add)

	assert.Equal(t, 6, concatAll(0)(MyInts{1, 2, 3}))
	assert.Equal(t, 16, concatAll(10)(MyInts{1, 2, 3}))
	assert.Equal(t, 5, concatAll(5)(MyInts{}))
}

// Test GenericMonadConcatAll with custom slice type
func TestGenericMonadConcatAll(t *testing.T) {
	type MyInts []int

	add := MakeSemigroup(func(a, b int) int { return a + b })
	monadConcatAll := GenericMonadConcatAll[MyInts](add)

	assert.Equal(t, 6, monadConcatAll(MyInts{1, 2, 3}, 0))
	assert.Equal(t, 16, monadConcatAll(MyInts{1, 2, 3}, 10))
	assert.Equal(t, 5, monadConcatAll(MyInts{}, 5))
}

// Test ApplySemigroup
func TestApplySemigroup(t *testing.T) {
	// Base semigroup for integers
	add := MakeSemigroup(func(a, b int) int { return a + b })

	// Simple HKT simulation using slices
	type HKT []int

	fmap := func(hkt HKT, f func(int) func(int) int) []func(int) int {
		result := make([]func(int) int, len(hkt))
		for i, v := range hkt {
			result[i] = f(v)
		}
		return result
	}

	fap := func(fs []func(int) int, hkt HKT) HKT {
		result := make(HKT, 0)
		for _, f := range fs {
			for _, v := range hkt {
				result = append(result, f(v))
			}
		}
		return result
	}

	applySG := ApplySemigroup(fmap, fap, add)

	hkt1 := HKT{1, 2}
	hkt2 := HKT{3, 4}

	result := applySG.Concat(hkt1, hkt2)
	// Should apply the semigroup operation to all combinations
	assert.NotEmpty(t, result)
}

// Test AltSemigroup
func TestAltSemigroup(t *testing.T) {
	// Simple HKT simulation using Option-like type
	type Option[A any] struct {
		value    A
		hasValue bool
	}

	falt := func(first Option[int], second func() Option[int]) Option[int] {
		if first.hasValue {
			return first
		}
		return second()
	}

	altSG := AltSemigroup(falt)

	some := Option[int]{value: 42, hasValue: true}
	none := Option[int]{hasValue: false}
	other := Option[int]{value: 100, hasValue: true}

	// First has value, should return first
	result1 := altSG.Concat(some, none)
	assert.True(t, result1.hasValue)
	assert.Equal(t, 42, result1.value)

	// First is none, should return second
	result2 := altSG.Concat(none, other)
	assert.True(t, result2.hasValue)
	assert.Equal(t, 100, result2.value)

	// Both have values, should return first
	result3 := altSG.Concat(some, other)
	assert.True(t, result3.hasValue)
	assert.Equal(t, 42, result3.value)
}

// Test associativity law for various semigroups
func TestAssociativityLaw(t *testing.T) {
	testCases := []struct {
		name    string
		sg      Semigroup[int]
		a, b, c int
	}{
		{"Addition", MakeSemigroup(func(a, b int) int { return a + b }), 1, 2, 3},
		{"Multiplication", MakeSemigroup(func(a, b int) int { return a * b }), 2, 3, 4},
		{"Max", MakeSemigroup(func(a, b int) int {
			if a > b {
				return a
			}
			return b
		}), 5, 10, 3},
		{"Min", MakeSemigroup(func(a, b int) int {
			if a < b {
				return a
			}
			return b
		}), 5, 10, 3},
		{"First", First[int](), 1, 2, 3},
		{"Last", Last[int](), 1, 2, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// (a • b) • c
			left := tc.sg.Concat(tc.sg.Concat(tc.a, tc.b), tc.c)
			// a • (b • c)
			right := tc.sg.Concat(tc.a, tc.sg.Concat(tc.b, tc.c))

			assert.Equal(t, left, right, "Associativity law failed for %s", tc.name)
		})
	}
}

// Test associativity law for string semigroups
func TestAssociativityLawString(t *testing.T) {
	concat := MakeSemigroup(func(a, b string) string { return a + b })

	a, b, c := "hello", " ", "world"

	left := concat.Concat(concat.Concat(a, b), c)
	right := concat.Concat(a, concat.Concat(b, c))

	assert.Equal(t, left, right)
	assert.Equal(t, "hello world", left)
}

// Test complex types
func TestComplexTypes(t *testing.T) {
	type Config struct {
		Timeout int
		Retries int
	}

	configSG := MakeSemigroup(func(a, b Config) Config {
		maxTimeout := a.Timeout
		if b.Timeout > maxTimeout {
			maxTimeout = b.Timeout
		}
		return Config{
			Timeout: maxTimeout,
			Retries: a.Retries + b.Retries,
		}
	})

	c1 := Config{Timeout: 30, Retries: 3}
	c2 := Config{Timeout: 60, Retries: 5}
	c3 := Config{Timeout: 45, Retries: 2}

	result := configSG.Concat(configSG.Concat(c1, c2), c3)
	assert.Equal(t, 60, result.Timeout)
	assert.Equal(t, 10, result.Retries)

	// Test associativity
	left := configSG.Concat(configSG.Concat(c1, c2), c3)
	right := configSG.Concat(c1, configSG.Concat(c2, c3))
	assert.Equal(t, left, right)
}

// Test with slices
func TestSliceSemigroup(t *testing.T) {
	sliceConcat := MakeSemigroup(func(a, b []int) []int {
		result := make([]int, len(a)+len(b))
		copy(result, a)
		copy(result[len(a):], b)
		return result
	})

	s1 := []int{1, 2, 3}
	s2 := []int{4, 5}
	s3 := []int{6}

	result := sliceConcat.Concat(sliceConcat.Concat(s1, s2), s3)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)

	// Test associativity
	left := sliceConcat.Concat(sliceConcat.Concat(s1, s2), s3)
	right := sliceConcat.Concat(s1, sliceConcat.Concat(s2, s3))
	assert.Equal(t, left, right)
}

// Test with maps
func TestMapSemigroup(t *testing.T) {
	mapMerge := MakeSemigroup(func(a, b map[string]int) map[string]int {
		result := make(map[string]int)
		for k, v := range a {
			result[k] = v
		}
		for k, v := range b {
			result[k] = v // Later values override
		}
		return result
	})

	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	m3 := map[string]int{"c": 5, "d": 6}

	result := mapMerge.Concat(mapMerge.Concat(m1, m2), m3)
	assert.Equal(t, 1, result["a"])
	assert.Equal(t, 3, result["b"])
	assert.Equal(t, 5, result["c"])
	assert.Equal(t, 6, result["d"])
}

// Benchmark tests
func BenchmarkFirst(b *testing.B) {
	first := First[int]()
	for b.Loop() {
		first.Concat(1, 2)
	}
}

func BenchmarkLast(b *testing.B) {
	last := Last[int]()
	for b.Loop() {
		last.Concat(1, 2)
	}
}

func BenchmarkMakeSemigroupAdd(b *testing.B) {
	add := MakeSemigroup(func(a, b int) int { return a + b })
	for b.Loop() {
		add.Concat(1, 2)
	}
}

func BenchmarkReverse(b *testing.B) {
	sub := MakeSemigroup(func(a, b int) int { return a - b })
	reversed := Reverse(sub)
	for b.Loop() {
		reversed.Concat(10, 3)
	}
}

func BenchmarkConcatAll(b *testing.B) {
	add := MakeSemigroup(func(a, b int) int { return a + b })
	concatAll := ConcatAll(add)
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for b.Loop() {
		concatAll(0)(arr)
	}
}

func BenchmarkFunctionSemigroup(b *testing.B) {
	add := MakeSemigroup(func(a, b int) int { return a + b })
	funcSG := FunctionSemigroup[string](add)

	f := func(s string) int { return len(s) }
	g := func(s string) int { return len(s) * 2 }
	combined := funcSG.Concat(f, g)

	b.ResetTimer()
	for b.Loop() {
		combined("hello")
	}
}
