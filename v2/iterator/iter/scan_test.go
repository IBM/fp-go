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

package iter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestScanBasic tests basic Scan functionality
func TestScanBasic(t *testing.T) {
	t.Run("running sum of integers", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{1, 3, 6, 10, 15}, result)
	})

	t.Run("running product", func(t *testing.T) {
		seq := From(2, 3, 4)
		scanned := Scan(func(acc, x int) int { return acc * x }, 1)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{2, 6, 24}, result)
	})

	t.Run("string concatenation", func(t *testing.T) {
		seq := From("a", "b", "c")
		scanned := Scan(func(acc, x string) string { return acc + x }, "")
		result := toSlice(scanned(seq))
		assert.Equal(t, []string{"a", "ab", "abc"}, result)
	})

	t.Run("string concatenation with separator", func(t *testing.T) {
		seq := From("hello", "world", "test")
		scanned := Scan(func(acc, x string) string {
			if acc == "" {
				return x
			}
			return acc + "-" + x
		}, "")
		result := toSlice(scanned(seq))
		assert.Equal(t, []string{"hello", "hello-world", "hello-world-test"}, result)
	})

	t.Run("single element", func(t *testing.T) {
		seq := From(42)
		scanned := Scan(func(acc, x int) int { return acc + x }, 10)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{52}, result)
	})

	t.Run("two elements", func(t *testing.T) {
		seq := From(5, 10)
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{5, 15}, result)
	})
}

// TestScanEmpty tests Scan with empty sequences
func TestScanEmpty(t *testing.T) {
	t.Run("empty integer sequence", func(t *testing.T) {
		seq := Empty[int]()
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		result := toSlice(scanned(seq))
		assert.Empty(t, result)
	})

	t.Run("empty string sequence", func(t *testing.T) {
		seq := Empty[string]()
		scanned := Scan(func(acc, x string) string { return acc + x }, "start")
		result := toSlice(scanned(seq))
		assert.Empty(t, result)
	})
}

// TestScanWithDifferentTypes tests Scan with different input/output types
func TestScanWithDifferentTypes(t *testing.T) {
	t.Run("int to string accumulation", func(t *testing.T) {
		seq := From(1, 2, 3)
		scanned := Scan(func(acc string, x int) string {
			return fmt.Sprintf("%s%d", acc, x)
		}, "")
		result := toSlice(scanned(seq))
		assert.Equal(t, []string{"1", "12", "123"}, result)
	})

	t.Run("string to int length accumulation", func(t *testing.T) {
		seq := From("a", "bb", "ccc")
		scanned := Scan(func(acc int, x string) int {
			return acc + len(x)
		}, 0)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{1, 3, 6}, result)
	})

	t.Run("accumulate into slice", func(t *testing.T) {
		seq := From(1, 2, 3)
		scanned := Scan(func(acc []int, x int) []int {
			return append(acc, x)
		}, []int{})
		result := toSlice(scanned(seq))
		assert.Equal(t, [][]int{
			{1},
			{1, 2},
			{1, 2, 3},
		}, result)
	})
}

// TestScanStateful tests Scan with stateful operations
func TestScanStateful(t *testing.T) {
	t.Run("tracking maximum", func(t *testing.T) {
		seq := From(3, 1, 4, 1, 5, 9, 2)
		scanned := Scan(func(acc, x int) int {
			if x > acc {
				return x
			}
			return acc
		}, 0)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{3, 3, 4, 4, 5, 9, 9}, result)
	})

	t.Run("tracking minimum", func(t *testing.T) {
		seq := From(5, 3, 8, 1, 4, 2)
		scanned := Scan(func(acc, x int) int {
			if acc == 0 || x < acc {
				return x
			}
			return acc
		}, 0)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{5, 3, 3, 1, 1, 1}, result)
	})

	t.Run("counting occurrences", func(t *testing.T) {
		seq := From(1, 2, 1, 3, 1, 2)
		scanned := Scan(func(acc map[int]int, x int) map[int]int {
			newMap := make(map[int]int)
			for k, v := range acc {
				newMap[k] = v
			}
			newMap[x]++
			return newMap
		}, map[int]int{})
		result := toSlice(scanned(seq))
		assert.Len(t, result, 6)
		assert.Equal(t, 1, result[0][1])
		assert.Equal(t, 1, result[1][2])
		assert.Equal(t, 2, result[2][1])
		assert.Equal(t, 3, result[4][1])
	})
}

// TestScanWithComplexTypes tests Scan with complex data types
func TestScanWithComplexTypes(t *testing.T) {
	type Point struct {
		X, Y int
	}

	t.Run("accumulate points", func(t *testing.T) {
		seq := From(Point{1, 0}, Point{0, 1}, Point{2, 2})
		scanned := Scan(func(acc, p Point) Point {
			return Point{acc.X + p.X, acc.Y + p.Y}
		}, Point{0, 0})
		result := toSlice(scanned(seq))
		assert.Equal(t, []Point{
			{1, 0},
			{1, 1},
			{3, 3},
		}, result)
	})

	t.Run("accumulate struct fields", func(t *testing.T) {
		type Data struct {
			Value int
			Count int
		}
		seq := From(5, 10, 15)
		scanned := Scan(func(acc Data, x int) Data {
			return Data{
				Value: acc.Value + x,
				Count: acc.Count + 1,
			}
		}, Data{0, 0})
		result := toSlice(scanned(seq))
		assert.Equal(t, []Data{
			{5, 1},
			{15, 2},
			{30, 3},
		}, result)
	})
}

// TestScanWithChainedOperations tests Scan combined with other operations
func TestScanWithChainedOperations(t *testing.T) {
	t.Run("scan then map", func(t *testing.T) {
		seq := From(1, 2, 3, 4)
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		mapped := MonadMap(scanned(seq), func(x int) int { return x * 2 })
		result := toSlice(mapped)
		assert.Equal(t, []int{2, 6, 12, 20}, result)
	})

	t.Run("map then scan", func(t *testing.T) {
		seq := From(1, 2, 3, 4)
		mapped := MonadMap(seq, func(x int) int { return x * 2 })
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		result := toSlice(scanned(mapped))
		assert.Equal(t, []int{2, 6, 12, 20}, result)
	})

	t.Run("scan then filter", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		filtered := MonadFilter(scanned(seq), func(x int) bool { return x%2 == 0 })
		result := toSlice(filtered)
		assert.Equal(t, []int{6, 10}, result)
	})

	t.Run("scan then take", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		taken := Take[int](3)(scanned(seq))
		result := toSlice(taken)
		assert.Equal(t, []int{1, 3, 6}, result)
	})
}

// TestScanWithCycle tests Scan with infinite sequences
func TestScanWithCycle(t *testing.T) {
	t.Run("scan cycled sequence with take", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		taken := Take[int](10)(scanned(cycled))
		result := toSlice(taken)
		// 1, 3, 6, 7, 9, 12, 13, 15, 18, 19
		assert.Equal(t, []int{1, 3, 6, 7, 9, 12, 13, 15, 18, 19}, result)
	})
}

// TestScanEarlyTermination tests that Scan respects early termination
func TestScanEarlyTermination(t *testing.T) {
	t.Run("terminates when yield returns false", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)

		count := 0
		for v := range scanned(seq) {
			count++
			if v >= 6 {
				break
			}
		}

		assert.Equal(t, 3, count) // Should stop at 6 (1+2+3)
	})
}

// TestScanWithInitialValue tests different initial values
func TestScanWithInitialValue(t *testing.T) {
	t.Run("non-zero initial value", func(t *testing.T) {
		seq := From(1, 2, 3)
		scanned := Scan(func(acc, x int) int { return acc + x }, 10)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{11, 13, 16}, result)
	})

	t.Run("negative initial value", func(t *testing.T) {
		seq := From(1, 2, 3)
		scanned := Scan(func(acc, x int) int { return acc + x }, -10)
		result := toSlice(scanned(seq))
		assert.Equal(t, []int{-9, -7, -4}, result)
	})

	t.Run("string initial value", func(t *testing.T) {
		seq := From("a", "b", "c")
		scanned := Scan(func(acc, x string) string { return acc + x }, "start:")
		result := toSlice(scanned(seq))
		assert.Equal(t, []string{"start:a", "start:ab", "start:abc"}, result)
	})
}

// TestScanLargeSequence tests Scan with larger sequences
func TestScanLargeSequence(t *testing.T) {
	t.Run("scan large sequence", func(t *testing.T) {
		data := make([]int, 100)
		for i := range data {
			data[i] = i + 1
		}
		seq := From(data...)
		scanned := Scan(func(acc, x int) int { return acc + x }, 0)
		result := toSlice(scanned(seq))

		assert.Len(t, result, 100)
		// Sum of 1 to n is n*(n+1)/2
		assert.Equal(t, 5050, result[99]) // Sum of 1 to 100
		assert.Equal(t, 1, result[0])
		assert.Equal(t, 3, result[1])
		assert.Equal(t, 6, result[2])
	})
}

// Benchmark tests
func BenchmarkScan(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	scanned := Scan(func(acc, x int) int { return acc + x }, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range scanned(seq) {
		}
	}
}

func BenchmarkScanLarge(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i + 1
	}
	seq := From(data...)
	scanned := Scan(func(acc, x int) int { return acc + x }, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range scanned(seq) {
		}
	}
}

// Example tests for documentation
func ExampleScan() {
	seq := From(1, 2, 3, 4, 5)
	runningSum := Scan(func(acc, x int) int { return acc + x }, 0)
	result := runningSum(seq)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 3 6 10 15
}

func ExampleScan_runningProduct() {
	seq := From(2, 3, 4)
	runningProduct := Scan(func(acc, x int) int { return acc * x }, 1)
	result := runningProduct(seq)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 2 6 24
}

func ExampleScan_stringConcatenation() {
	seq := From("a", "b", "c")
	concat := Scan(func(acc, x string) string { return acc + x }, "")
	result := concat(seq)

	for v := range result {
		fmt.Printf("%s ", v)
	}
	// Output: a ab abc
}

func ExampleScan_trackingMaximum() {
	seq := From(3, 1, 4, 1, 5, 9, 2)
	maxSoFar := Scan(func(acc, x int) int {
		if x > acc {
			return x
		}
		return acc
	}, 0)
	result := maxSoFar(seq)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 3 3 4 4 5 9 9
}

func ExampleScan_empty() {
	seq := Empty[int]()
	runningSum := Scan(func(acc, x int) int { return acc + x }, 0)
	result := runningSum(seq)

	count := 0
	for range result {
		count++
	}
	fmt.Printf("Count: %d\n", count)
	// Output: Count: 0
}
