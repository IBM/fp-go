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

package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Custom string type for testing generic constraints
type MyString string

func TestToBytes(t *testing.T) {
	t.Run("regular string", func(t *testing.T) {
		result := ToBytes("hello")
		expected := []byte{'h', 'e', 'l', 'l', 'o'}
		assert.Equal(t, expected, result)
	})

	t.Run("empty string", func(t *testing.T) {
		result := ToBytes("")
		assert.Equal(t, []byte{}, result)
	})

	t.Run("custom string type", func(t *testing.T) {
		result := ToBytes(MyString("test"))
		expected := []byte{'t', 'e', 's', 't'}
		assert.Equal(t, expected, result)
	})

	t.Run("unicode string", func(t *testing.T) {
		result := ToBytes("你好")
		// UTF-8 encoding: 你 = E4 BD A0, 好 = E5 A5 BD
		assert.Equal(t, 6, len(result))
	})
}

func TestToRunes(t *testing.T) {
	t.Run("regular string", func(t *testing.T) {
		result := ToRunes("hello")
		expected := []rune{'h', 'e', 'l', 'l', 'o'}
		assert.Equal(t, expected, result)
	})

	t.Run("empty string", func(t *testing.T) {
		result := ToRunes("")
		assert.Equal(t, []rune{}, result)
	})

	t.Run("custom string type", func(t *testing.T) {
		result := ToRunes(MyString("test"))
		expected := []rune{'t', 'e', 's', 't'}
		assert.Equal(t, expected, result)
	})

	t.Run("unicode string", func(t *testing.T) {
		result := ToRunes("你好")
		assert.Equal(t, 2, len(result))
		assert.Equal(t, '你', result[0])
		assert.Equal(t, '好', result[1])
	})

	t.Run("mixed ascii and unicode", func(t *testing.T) {
		result := ToRunes("Hello世界")
		assert.Equal(t, 7, len(result))
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		assert.True(t, IsEmpty(""))
	})

	t.Run("non-empty string", func(t *testing.T) {
		assert.False(t, IsEmpty("hello"))
	})

	t.Run("whitespace string", func(t *testing.T) {
		assert.False(t, IsEmpty(" "))
		assert.False(t, IsEmpty("\t"))
		assert.False(t, IsEmpty("\n"))
	})

	t.Run("custom string type empty", func(t *testing.T) {
		assert.True(t, IsEmpty(MyString("")))
	})

	t.Run("custom string type non-empty", func(t *testing.T) {
		assert.False(t, IsEmpty(MyString("test")))
	})
}

func TestIsNonEmpty(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		assert.False(t, IsNonEmpty(""))
	})

	t.Run("non-empty string", func(t *testing.T) {
		assert.True(t, IsNonEmpty("hello"))
	})

	t.Run("whitespace string", func(t *testing.T) {
		assert.True(t, IsNonEmpty(" "))
		assert.True(t, IsNonEmpty("\t"))
		assert.True(t, IsNonEmpty("\n"))
	})

	t.Run("custom string type empty", func(t *testing.T) {
		assert.False(t, IsNonEmpty(MyString("")))
	})

	t.Run("custom string type non-empty", func(t *testing.T) {
		assert.True(t, IsNonEmpty(MyString("test")))
	})

	t.Run("single character", func(t *testing.T) {
		assert.True(t, IsNonEmpty("a"))
	})
}

func TestSize(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		assert.Equal(t, 0, Size(""))
	})

	t.Run("ascii string", func(t *testing.T) {
		assert.Equal(t, 5, Size("hello"))
		assert.Equal(t, 11, Size("hello world"))
	})

	t.Run("custom string type", func(t *testing.T) {
		assert.Equal(t, 4, Size(MyString("test")))
	})

	t.Run("unicode string - returns byte count", func(t *testing.T) {
		// Note: Size returns byte length, not rune count
		assert.Equal(t, 6, Size("你好"))   // 2 Chinese characters = 6 bytes in UTF-8
		assert.Equal(t, 5, Size("café")) // 'c', 'a', 'f' = 3 bytes, 'é' = 2 bytes in UTF-8
	})

	t.Run("single character", func(t *testing.T) {
		assert.Equal(t, 1, Size("a"))
	})

	t.Run("whitespace", func(t *testing.T) {
		assert.Equal(t, 1, Size(" "))
		assert.Equal(t, 1, Size("\t"))
		assert.Equal(t, 1, Size("\n"))
	})
}
