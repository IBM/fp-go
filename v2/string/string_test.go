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

package string

import (
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	assert.True(t, IsEmpty(""))
	assert.False(t, IsEmpty("Carsten"))
}

func TestJoin(t *testing.T) {

	x := Join(",")(A.From("a", "b", "c"))
	assert.Equal(t, x, x)

	assert.Equal(t, "a,b,c", Join(",")(A.From("a", "b", "c")))
	assert.Equal(t, "a", Join(",")(A.From("a")))
	assert.Equal(t, "", Join(",")(A.Empty[string]()))
}

func TestEquals(t *testing.T) {
	assert.True(t, Equals("a")("a"))
	assert.False(t, Equals("a")("b"))
	assert.False(t, Equals("b")("a"))
}

func TestIncludes(t *testing.T) {
	assert.True(t, Includes("a")("bab"))
	assert.False(t, Includes("bab")("a"))
	assert.False(t, Includes("b")("a"))
}

func TestHasPrefix(t *testing.T) {
	assert.True(t, HasPrefix("prefix")("prefixbab"))
	assert.False(t, HasPrefix("bab")("a"))
	assert.False(t, HasPrefix("b")("a"))
}

func TestEq(t *testing.T) {
	assert.True(t, Eq("hello", "hello"))
	assert.False(t, Eq("hello", "world"))
	assert.True(t, Eq("", ""))
	assert.False(t, Eq("a", ""))
}

func TestToBytes(t *testing.T) {
	result := ToBytes("hello")
	expected := []byte{'h', 'e', 'l', 'l', 'o'}
	assert.Equal(t, expected, result)

	empty := ToBytes("")
	assert.Equal(t, []byte{}, empty)
}

func TestToRunes(t *testing.T) {
	result := ToRunes("hello")
	expected := []rune{'h', 'e', 'l', 'l', 'o'}
	assert.Equal(t, expected, result)

	// Test with unicode characters
	unicode := ToRunes("你好")
	assert.Equal(t, 2, len(unicode))

	empty := ToRunes("")
	assert.Equal(t, []rune{}, empty)
}

func TestIsNonEmpty(t *testing.T) {
	assert.False(t, IsNonEmpty(""))
	assert.True(t, IsNonEmpty("Carsten"))
	assert.True(t, IsNonEmpty(" "))
}

func TestSize(t *testing.T) {
	assert.Equal(t, 0, Size(""))
	assert.Equal(t, 5, Size("hello"))
	assert.Equal(t, 7, Size("Carsten"))
	// Note: Size returns byte length, not rune count
	assert.Equal(t, 6, Size("你好")) // 2 Chinese characters = 6 bytes in UTF-8
}

func TestFormat(t *testing.T) {
	formatInt := Format[int]("Number: %d")
	assert.Equal(t, "Number: 42", formatInt(42))

	formatString := Format[string]("Hello, %s!")
	assert.Equal(t, "Hello, World!", formatString("World"))

	formatFloat := Format[float64]("Value: %.2f")
	assert.Equal(t, "Value: 3.14", formatFloat(3.14159))
}

func TestIntersperse(t *testing.T) {
	join := Intersperse(", ")
	assert.Equal(t, "a, b", join("a", "b"))
	assert.Equal(t, "hello, world", join("hello", "world"))

	joinDash := Intersperse("-")
	assert.Equal(t, "foo-bar", joinDash("foo", "bar"))

	// Test with empty strings - should not add separator (monoid identity)
	assert.Equal(t, "b", join("", "b"))
	assert.Equal(t, "a", join("a", ""))
	assert.Equal(t, "", join("", ""))
}

func TestToUpperCase(t *testing.T) {
	assert.Equal(t, "HELLO", ToUpperCase("hello"))
	assert.Equal(t, "WORLD", ToUpperCase("WoRlD"))
	assert.Equal(t, "", ToUpperCase(""))
}

func TestToLowerCase(t *testing.T) {
	assert.Equal(t, "hello", ToLowerCase("HELLO"))
	assert.Equal(t, "world", ToLowerCase("WoRlD"))
	assert.Equal(t, "", ToLowerCase(""))
}

func TestOrd(t *testing.T) {
	assert.True(t, Ord.Compare("a", "b") < 0)
	assert.True(t, Ord.Compare("b", "a") > 0)
	assert.True(t, Ord.Compare("a", "a") == 0)
	assert.True(t, Ord.Compare("abc", "abd") < 0)
}

func TestPrepend(t *testing.T) {
	t.Run("prepend to non-empty string", func(t *testing.T) {
		addHello := Prepend("Hello, ")
		result := addHello("World")
		assert.Equal(t, "Hello, World", result)
	})

	t.Run("prepend to empty string", func(t *testing.T) {
		addPrefix := Prepend("prefix")
		result := addPrefix("")
		assert.Equal(t, "prefix", result)
	})

	t.Run("prepend empty string", func(t *testing.T) {
		addNothing := Prepend("")
		result := addNothing("test")
		assert.Equal(t, "test", result)
	})

	t.Run("prepend with special characters", func(t *testing.T) {
		addSymbols := Prepend(">>> ")
		result := addSymbols("message")
		assert.Equal(t, ">>> message", result)
	})

	t.Run("prepend with unicode", func(t *testing.T) {
		addEmoji := Prepend("🎉 ")
		result := addEmoji("Party!")
		assert.Equal(t, "🎉 Party!", result)
	})

	t.Run("multiple prepends", func(t *testing.T) {
		addA := Prepend("A")
		addB := Prepend("B")
		result := addB(addA("C"))
		assert.Equal(t, "BAC", result)
	})
}

func TestAppend(t *testing.T) {
	t.Run("append to non-empty string", func(t *testing.T) {
		addExclamation := Append("!")
		result := addExclamation("Hello")
		assert.Equal(t, "Hello!", result)
	})

	t.Run("append to empty string", func(t *testing.T) {
		addSuffix := Append("suffix")
		result := addSuffix("")
		assert.Equal(t, "suffix", result)
	})

	t.Run("append empty string", func(t *testing.T) {
		addNothing := Append("")
		result := addNothing("test")
		assert.Equal(t, "test", result)
	})

	t.Run("append with special characters", func(t *testing.T) {
		addEllipsis := Append("...")
		result := addEllipsis("To be continued")
		assert.Equal(t, "To be continued...", result)
	})

	t.Run("append with unicode", func(t *testing.T) {
		addEmoji := Append(" 🎉")
		result := addEmoji("Party")
		assert.Equal(t, "Party 🎉", result)
	})

	t.Run("multiple appends", func(t *testing.T) {
		addA := Append("A")
		addB := Append("B")
		result := addB(addA("C"))
		assert.Equal(t, "CAB", result)
	})
}

func TestPrependAndAppend(t *testing.T) {
	t.Run("combine prepend and append", func(t *testing.T) {
		addPrefix := Prepend("[ ")
		addSuffix := Append(" ]")
		result := addSuffix(addPrefix("content"))
		assert.Equal(t, "[ content ]", result)
	})

	t.Run("chain multiple operations", func(t *testing.T) {
		addQuotes := Prepend("\"")
		closeQuotes := Append("\"")
		addLabel := Prepend("Value: ")

		result := addLabel(addQuotes(closeQuotes("test")))
		assert.Equal(t, "Value: \"test\"", result)
	})
}
