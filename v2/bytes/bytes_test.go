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

package bytes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	t.Run("returns empty byte slice", func(t *testing.T) {
		result := Empty()
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("is identity for Monoid", func(t *testing.T) {
		data := []byte("test")

		// Left identity: empty + data = data
		left := Monoid.Concat(Empty(), data)
		assert.Equal(t, data, left)

		// Right identity: data + empty = data
		right := Monoid.Concat(data, Empty())
		assert.Equal(t, data, right)
	})
}

func TestToString(t *testing.T) {
	t.Run("converts byte slice to string", func(t *testing.T) {
		result := ToString([]byte("hello"))
		assert.Equal(t, "hello", result)
	})

	t.Run("handles empty byte slice", func(t *testing.T) {
		result := ToString([]byte{})
		assert.Equal(t, "", result)
	})

	t.Run("handles binary data", func(t *testing.T) {
		data := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello"
		result := ToString(data)
		assert.Equal(t, "Hello", result)
	})
}

func TestSize(t *testing.T) {
	t.Run("returns size of byte slice", func(t *testing.T) {
		assert.Equal(t, 0, Size([]byte{}))
		assert.Equal(t, 5, Size([]byte("hello")))
		assert.Equal(t, 10, Size([]byte("0123456789")))
	})

	t.Run("handles binary data", func(t *testing.T) {
		data := []byte{0x00, 0x01, 0x02, 0x03}
		assert.Equal(t, 4, Size(data))
	})
}

func TestMonoidConcat(t *testing.T) {
	t.Run("concatenates two byte slices", func(t *testing.T) {
		result := Monoid.Concat([]byte("Hello"), []byte(" World"))
		assert.Equal(t, []byte("Hello World"), result)
	})

	t.Run("handles empty slices", func(t *testing.T) {
		result1 := Monoid.Concat([]byte{}, []byte("test"))
		assert.Equal(t, []byte("test"), result1)

		result2 := Monoid.Concat([]byte("test"), []byte{})
		assert.Equal(t, []byte("test"), result2)

		result3 := Monoid.Concat([]byte{}, []byte{})
		assert.Equal(t, []byte{}, result3)
	})

	t.Run("associativity: (a + b) + c = a + (b + c)", func(t *testing.T) {
		a := []byte("a")
		b := []byte("b")
		c := []byte("c")

		left := Monoid.Concat(Monoid.Concat(a, b), c)
		right := Monoid.Concat(a, Monoid.Concat(b, c))

		assert.Equal(t, left, right)
	})
}

func TestConcatAll(t *testing.T) {
	t.Run("concatenates multiple byte slices", func(t *testing.T) {
		result := ConcatAll(
			[]byte("Hello"),
			[]byte(" "),
			[]byte("World"),
			[]byte("!"),
		)
		assert.Equal(t, []byte("Hello World!"), result)
	})

	t.Run("handles empty input", func(t *testing.T) {
		result := ConcatAll()
		assert.Equal(t, []byte{}, result)
	})

	t.Run("handles single slice", func(t *testing.T) {
		result := ConcatAll([]byte("test"))
		assert.Equal(t, []byte("test"), result)
	})

	t.Run("handles slices with empty elements", func(t *testing.T) {
		result := ConcatAll(
			[]byte("a"),
			[]byte{},
			[]byte("b"),
			[]byte{},
			[]byte("c"),
		)
		assert.Equal(t, []byte("abc"), result)
	})

	t.Run("handles binary data", func(t *testing.T) {
		result := ConcatAll(
			[]byte{0x01, 0x02},
			[]byte{0x03, 0x04},
			[]byte{0x05},
		)
		assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05}, result)
	})
}

func TestOrd(t *testing.T) {
	t.Run("compares byte slices lexicographically", func(t *testing.T) {
		// "abc" < "abd"
		assert.Equal(t, -1, Ord.Compare([]byte("abc"), []byte("abd")))

		// "abd" > "abc"
		assert.Equal(t, 1, Ord.Compare([]byte("abd"), []byte("abc")))

		// "abc" == "abc"
		assert.Equal(t, 0, Ord.Compare([]byte("abc"), []byte("abc")))
	})

	t.Run("handles different lengths", func(t *testing.T) {
		// "ab" < "abc"
		assert.Equal(t, -1, Ord.Compare([]byte("ab"), []byte("abc")))

		// "abc" > "ab"
		assert.Equal(t, 1, Ord.Compare([]byte("abc"), []byte("ab")))
	})

	t.Run("handles empty slices", func(t *testing.T) {
		// "" < "a"
		assert.Equal(t, -1, Ord.Compare([]byte{}, []byte("a")))

		// "a" > ""
		assert.Equal(t, 1, Ord.Compare([]byte("a"), []byte{}))

		// "" == ""
		assert.Equal(t, 0, Ord.Compare([]byte{}, []byte{}))
	})

	t.Run("Equals method works", func(t *testing.T) {
		assert.True(t, Ord.Equals([]byte("test"), []byte("test")))
		assert.False(t, Ord.Equals([]byte("test"), []byte("Test")))
		assert.True(t, Ord.Equals([]byte{}, []byte{}))
	})

	t.Run("handles binary data", func(t *testing.T) {
		assert.Equal(t, -1, Ord.Compare([]byte{0x01}, []byte{0x02}))
		assert.Equal(t, 1, Ord.Compare([]byte{0x02}, []byte{0x01}))
		assert.Equal(t, 0, Ord.Compare([]byte{0x01, 0x02}, []byte{0x01, 0x02}))
	})
}

// TestOrdProperties tests mathematical properties of Ord
func TestOrdProperties(t *testing.T) {
	t.Run("reflexivity: x == x", func(t *testing.T) {
		testCases := [][]byte{
			[]byte{},
			[]byte("a"),
			[]byte("test"),
			[]byte{0x01, 0x02, 0x03},
		}

		for _, tc := range testCases {
			assert.Equal(t, 0, Ord.Compare(tc, tc),
				"Compare(%v, %v) should be 0", tc, tc)
			assert.True(t, Ord.Equals(tc, tc),
				"Equals(%v, %v) should be true", tc, tc)
		}
	})

	t.Run("antisymmetry: if x <= y and y <= x then x == y", func(t *testing.T) {
		testCases := []struct {
			a, b []byte
		}{
			{[]byte("abc"), []byte("abc")},
			{[]byte{}, []byte{}},
			{[]byte{0x01}, []byte{0x01}},
		}

		for _, tc := range testCases {
			cmp1 := Ord.Compare(tc.a, tc.b)
			cmp2 := Ord.Compare(tc.b, tc.a)

			if cmp1 <= 0 && cmp2 <= 0 {
				assert.True(t, Ord.Equals(tc.a, tc.b),
					"If %v <= %v and %v <= %v, they should be equal", tc.a, tc.b, tc.b, tc.a)
			}
		}
	})

	t.Run("transitivity: if x <= y and y <= z then x <= z", func(t *testing.T) {
		x := []byte("a")
		y := []byte("b")
		z := []byte("c")

		cmpXY := Ord.Compare(x, y)
		cmpYZ := Ord.Compare(y, z)
		cmpXZ := Ord.Compare(x, z)

		if cmpXY <= 0 && cmpYZ <= 0 {
			assert.True(t, cmpXZ <= 0,
				"If %v <= %v and %v <= %v, then %v <= %v", x, y, y, z, x, z)
		}
	})

	t.Run("totality: either x <= y or y <= x", func(t *testing.T) {
		testCases := []struct {
			a, b []byte
		}{
			{[]byte("abc"), []byte("abd")},
			{[]byte("xyz"), []byte("abc")},
			{[]byte{}, []byte("a")},
			{[]byte{0x01}, []byte{0x02}},
		}

		for _, tc := range testCases {
			cmp1 := Ord.Compare(tc.a, tc.b)
			cmp2 := Ord.Compare(tc.b, tc.a)

			assert.True(t, cmp1 <= 0 || cmp2 <= 0,
				"Either %v <= %v or %v <= %v must be true", tc.a, tc.b, tc.b, tc.a)
		}
	})
}

// TestEdgeCases tests edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	t.Run("very large byte slices", func(t *testing.T) {
		large := make([]byte, 1000000)
		for i := range large {
			large[i] = byte(i % 256)
		}

		size := Size(large)
		assert.Equal(t, 1000000, size)

		str := ToString(large)
		assert.Equal(t, 1000000, len(str))
	})

	t.Run("concatenating many slices", func(t *testing.T) {
		slices := make([][]byte, 100)
		for i := range slices {
			slices[i] = []byte{byte(i)}
		}

		result := ConcatAll(slices...)
		assert.Equal(t, 100, Size(result))
	})

	t.Run("null bytes in slice", func(t *testing.T) {
		data := []byte{0x00, 0x01, 0x00, 0x02}
		size := Size(data)
		assert.Equal(t, 4, size)

		str := ToString(data)
		assert.Equal(t, 4, len(str))
	})

	t.Run("comparing slices with null bytes", func(t *testing.T) {
		a := []byte{0x00, 0x01}
		b := []byte{0x00, 0x02}
		assert.Equal(t, -1, Ord.Compare(a, b))
	})
}

// TestMonoidConcatPerformance tests concatenation performance characteristics
func TestMonoidConcatPerformance(t *testing.T) {
	t.Run("ConcatAll vs repeated Concat", func(t *testing.T) {
		slices := [][]byte{
			[]byte("a"),
			[]byte("b"),
			[]byte("c"),
			[]byte("d"),
			[]byte("e"),
		}

		// Using ConcatAll
		result1 := ConcatAll(slices...)

		// Using repeated Concat
		result2 := Monoid.Empty()
		for _, s := range slices {
			result2 = Monoid.Concat(result2, s)
		}

		assert.Equal(t, result1, result2)
		assert.Equal(t, []byte("abcde"), result1)
	})
}

// TestRoundTrip tests round-trip conversions
func TestRoundTrip(t *testing.T) {
	t.Run("string to bytes to string", func(t *testing.T) {
		original := "Hello, World! 世界"
		bytes := []byte(original)
		result := ToString(bytes)
		assert.Equal(t, original, result)
	})

	t.Run("bytes to string to bytes", func(t *testing.T) {
		original := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
		str := ToString(original)
		result := []byte(str)
		assert.Equal(t, original, result)
	})
}

// TestConcatAllVariadic tests ConcatAll with various argument counts
func TestConcatAllVariadic(t *testing.T) {
	t.Run("zero arguments", func(t *testing.T) {
		result := ConcatAll()
		assert.Equal(t, []byte{}, result)
	})

	t.Run("one argument", func(t *testing.T) {
		result := ConcatAll([]byte("test"))
		assert.Equal(t, []byte("test"), result)
	})

	t.Run("two arguments", func(t *testing.T) {
		result := ConcatAll([]byte("hello"), []byte("world"))
		assert.Equal(t, []byte("helloworld"), result)
	})

	t.Run("many arguments", func(t *testing.T) {
		result := ConcatAll(
			[]byte("a"),
			[]byte("b"),
			[]byte("c"),
			[]byte("d"),
			[]byte("e"),
			[]byte("f"),
			[]byte("g"),
			[]byte("h"),
			[]byte("i"),
			[]byte("j"),
		)
		assert.Equal(t, []byte("abcdefghij"), result)
	})
}

// Benchmark tests
func BenchmarkToString(b *testing.B) {
	data := []byte("Hello, World!")

	b.Run("small", func(b *testing.B) {
		for b.Loop() {
			_ = ToString(data)
		}
	})

	b.Run("large", func(b *testing.B) {
		large := make([]byte, 10000)
		for i := range large {
			large[i] = byte(i % 256)
		}
		b.ResetTimer()
		for b.Loop() {
			_ = ToString(large)
		}
	})
}

func BenchmarkSize(b *testing.B) {
	data := []byte("Hello, World!")

	for b.Loop() {
		_ = Size(data)
	}
}

func BenchmarkMonoidConcat(b *testing.B) {
	a := []byte("Hello")
	c := []byte(" World")

	b.Run("small slices", func(b *testing.B) {
		for b.Loop() {
			_ = Monoid.Concat(a, c)
		}
	})

	b.Run("large slices", func(b *testing.B) {
		large1 := make([]byte, 10000)
		large2 := make([]byte, 10000)
		b.ResetTimer()
		for b.Loop() {
			_ = Monoid.Concat(large1, large2)
		}
	})
}

func BenchmarkConcatAll(b *testing.B) {
	slices := [][]byte{
		[]byte("Hello"),
		[]byte(" "),
		[]byte("World"),
		[]byte("!"),
	}

	b.Run("few slices", func(b *testing.B) {
		for b.Loop() {
			_ = ConcatAll(slices...)
		}
	})

	b.Run("many slices", func(b *testing.B) {
		many := make([][]byte, 100)
		for i := range many {
			many[i] = []byte{byte(i)}
		}
		b.ResetTimer()
		for b.Loop() {
			_ = ConcatAll(many...)
		}
	})
}

func BenchmarkOrdCompare(b *testing.B) {
	a := []byte("abc")
	c := []byte("abd")

	b.Run("equal", func(b *testing.B) {
		for b.Loop() {
			_ = Ord.Compare(a, a)
		}
	})

	b.Run("different", func(b *testing.B) {
		for b.Loop() {
			_ = Ord.Compare(a, c)
		}
	})

	b.Run("large slices", func(b *testing.B) {
		large1 := make([]byte, 10000)
		large2 := make([]byte, 10000)
		large2[9999] = 1
		b.ResetTimer()
		for b.Loop() {
			_ = Ord.Compare(large1, large2)
		}
	})
}

// Example tests
func ExampleEmpty() {
	empty := Empty()
	println(len(empty)) // 0

	// Output:
}

func ExampleToString() {
	str := ToString([]byte("hello"))
	println(str) // hello

	// Output:
}

func ExampleSize() {
	size := Size([]byte("hello"))
	println(size) // 5

	// Output:
}

func ExampleConcatAll() {
	result := ConcatAll(
		[]byte("Hello"),
		[]byte(" "),
		[]byte("World"),
	)
	println(string(result)) // Hello World

	// Output:
}

func ExampleMonoid_concat() {
	result := Monoid.Concat([]byte("Hello"), []byte(" World"))
	println(string(result)) // Hello World

	// Output:
}

func ExampleOrd_compare() {
	cmp := Ord.Compare([]byte("abc"), []byte("abd"))
	println(cmp) // -1 (abc < abd)

	// Output:
}
