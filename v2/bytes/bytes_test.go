// Copyright (c) 2023 IBM Corp.
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
