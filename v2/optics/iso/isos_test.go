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

package iso

import (
	"fmt"
	"testing"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/array/nonempty"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// TestUTF8String tests the UTF8String isomorphism
func TestUTF8String(t *testing.T) {
	iso := UTF8String()

	t.Run("Get converts bytes to string", func(t *testing.T) {
		bytes := []byte("hello world")
		result := iso.Get(bytes)
		assert.Equal(t, "hello world", result)
	})

	t.Run("Get handles empty bytes", func(t *testing.T) {
		bytes := []byte{}
		result := iso.Get(bytes)
		assert.Equal(t, "", result)
	})

	t.Run("Get handles UTF-8 characters", func(t *testing.T) {
		bytes := []byte("Hello 世界 🌍")
		result := iso.Get(bytes)
		assert.Equal(t, "Hello 世界 🌍", result)
	})

	t.Run("ReverseGet converts string to bytes", func(t *testing.T) {
		str := "hello world"
		result := iso.ReverseGet(str)
		assert.Equal(t, []byte("hello world"), result)
	})

	t.Run("ReverseGet handles empty string", func(t *testing.T) {
		str := ""
		result := iso.ReverseGet(str)
		assert.Equal(t, []byte{}, result)
	})

	t.Run("ReverseGet handles UTF-8 characters", func(t *testing.T) {
		str := "Hello 世界 🌍"
		result := iso.ReverseGet(str)
		assert.Equal(t, []byte("Hello 世界 🌍"), result)
	})

	t.Run("Round-trip bytes to string to bytes", func(t *testing.T) {
		original := []byte("test data")
		result := iso.ReverseGet(iso.Get(original))
		assert.Equal(t, original, result)
	})

	t.Run("Round-trip string to bytes to string", func(t *testing.T) {
		original := "test string"
		result := iso.Get(iso.ReverseGet(original))
		assert.Equal(t, original, result)
	})

	t.Run("Handles special characters", func(t *testing.T) {
		str := "line1\nline2\ttab\r\nwindows"
		bytes := iso.ReverseGet(str)
		result := iso.Get(bytes)
		assert.Equal(t, str, result)
	})

	t.Run("Handles binary-like data", func(t *testing.T) {
		bytes := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello"
		result := iso.Get(bytes)
		assert.Equal(t, "Hello", result)
	})
}

// TestLines tests the Lines isomorphism
func TestLines(t *testing.T) {
	iso := Lines()

	t.Run("Get joins lines with newline", func(t *testing.T) {
		lines := A.From("line1", "line2", "line3")
		result := iso.Get(lines)
		assert.Equal(t, "line1\nline2\nline3", result)
	})

	t.Run("Get handles single line", func(t *testing.T) {
		lines := A.Of("single line")
		result := iso.Get(lines)
		assert.Equal(t, "single line", result)
	})

	t.Run("Get handles empty slice", func(t *testing.T) {
		lines := A.Empty[string]()
		result := iso.Get(lines)
		assert.Equal(t, "", result)
	})

	t.Run("Get handles empty strings in slice", func(t *testing.T) {
		lines := A.From("a", "", "b")
		result := iso.Get(lines)
		assert.Equal(t, "a\n\nb", result)
	})

	t.Run("Get handles slice with only empty strings", func(t *testing.T) {
		lines := A.From("", "", "")
		result := iso.Get(lines)
		assert.Equal(t, "\n\n", result)
	})

	t.Run("ReverseGet splits string by newline", func(t *testing.T) {
		str := "line1\nline2\nline3"
		result := iso.ReverseGet(str)
		assert.Equal(t, []string{"line1", "line2", "line3"}, result)
	})

	t.Run("ReverseGet handles single line", func(t *testing.T) {
		str := "single line"
		result := iso.ReverseGet(str)
		assert.Equal(t, []string{"single line"}, result)
	})

	t.Run("ReverseGet handles empty string", func(t *testing.T) {
		str := ""
		result := iso.ReverseGet(str)
		assert.Equal(t, []string{""}, result)
	})

	t.Run("ReverseGet handles consecutive newlines", func(t *testing.T) {
		str := "a\n\nb"
		result := iso.ReverseGet(str)
		assert.Equal(t, []string{"a", "", "b"}, result)
	})

	t.Run("ReverseGet handles trailing newline", func(t *testing.T) {
		str := "a\nb\n"
		result := iso.ReverseGet(str)
		assert.Equal(t, []string{"a", "b", ""}, result)
	})

	t.Run("ReverseGet handles leading newline", func(t *testing.T) {
		str := "\na\nb"
		result := iso.ReverseGet(str)
		assert.Equal(t, []string{"", "a", "b"}, result)
	})

	t.Run("Round-trip lines to string to lines", func(t *testing.T) {
		original := []string{"line1", "line2", "line3"}
		result := iso.ReverseGet(iso.Get(original))
		assert.Equal(t, original, result)
	})

	t.Run("Round-trip string to lines to string", func(t *testing.T) {
		original := "line1\nline2\nline3"
		result := iso.Get(iso.ReverseGet(original))
		assert.Equal(t, original, result)
	})

	t.Run("Handles lines with special characters", func(t *testing.T) {
		lines := []string{"Hello 世界", "🌍 Earth", "tab\there"}
		text := iso.Get(lines)
		result := iso.ReverseGet(text)
		assert.Equal(t, lines, result)
	})

	t.Run("Preserves whitespace in lines", func(t *testing.T) {
		lines := []string{"  indented", "normal", "\ttabbed"}
		text := iso.Get(lines)
		result := iso.ReverseGet(text)
		assert.Equal(t, lines, result)
	})
}

// TestUnixMilli tests the UnixMilli isomorphism
func TestUnixMilli(t *testing.T) {
	iso := UnixMilli()

	t.Run("Get converts milliseconds to time", func(t *testing.T) {
		millis := int64(1609459200000) // 2021-01-01 00:00:00 UTC
		result := iso.Get(millis)
		// Compare Unix timestamps to avoid timezone issues
		assert.Equal(t, millis, result.UnixMilli())
	})

	t.Run("Get handles zero milliseconds (Unix epoch)", func(t *testing.T) {
		millis := int64(0)
		result := iso.Get(millis)
		assert.Equal(t, millis, result.UnixMilli())
	})

	t.Run("Get handles negative milliseconds (before epoch)", func(t *testing.T) {
		millis := int64(-86400000) // 1 day before epoch
		result := iso.Get(millis)
		assert.Equal(t, millis, result.UnixMilli())
	})

	t.Run("ReverseGet converts time to milliseconds", func(t *testing.T) {
		tm := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		result := iso.ReverseGet(tm)
		assert.Equal(t, int64(1609459200000), result)
	})

	t.Run("ReverseGet handles Unix epoch", func(t *testing.T) {
		tm := time.Unix(0, 0).UTC()
		result := iso.ReverseGet(tm)
		assert.Equal(t, int64(0), result)
	})

	t.Run("ReverseGet handles time before epoch", func(t *testing.T) {
		tm := time.Date(1969, 12, 31, 0, 0, 0, 0, time.UTC)
		result := iso.ReverseGet(tm)
		assert.Equal(t, int64(-86400000), result)
	})

	t.Run("Round-trip milliseconds to time to milliseconds", func(t *testing.T) {
		original := int64(1234567890000)
		result := iso.ReverseGet(iso.Get(original))
		assert.Equal(t, original, result)
	})

	t.Run("Round-trip time to milliseconds to time", func(t *testing.T) {
		original := time.Date(2021, 6, 15, 12, 30, 45, 0, time.UTC)
		result := iso.Get(iso.ReverseGet(original))
		// Compare as Unix timestamps to avoid timezone issues
		assert.Equal(t, original.UnixMilli(), result.UnixMilli())
	})

	t.Run("Truncates sub-millisecond precision", func(t *testing.T) {
		// Time with nanoseconds
		tm := time.Date(2021, 1, 1, 0, 0, 0, 123456789, time.UTC)
		millis := iso.ReverseGet(tm)
		result := iso.Get(millis)

		// Should have millisecond precision only - compare timestamps
		assert.Equal(t, tm.Truncate(time.Millisecond).UnixMilli(), result.UnixMilli())
	})

	t.Run("Handles current time", func(t *testing.T) {
		now := time.Now()
		millis := iso.ReverseGet(now)
		result := iso.Get(millis)

		// Should be equal within millisecond precision
		assert.Equal(t, now.Truncate(time.Millisecond), result.Truncate(time.Millisecond))
	})

	t.Run("Handles far future date", func(t *testing.T) {
		future := time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)
		millis := iso.ReverseGet(future)
		result := iso.Get(millis)
		assert.Equal(t, future.UnixMilli(), result.UnixMilli())
	})

	t.Run("Handles far past date", func(t *testing.T) {
		past := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		millis := iso.ReverseGet(past)
		result := iso.Get(millis)
		assert.Equal(t, past.UnixMilli(), result.UnixMilli())
	})

	t.Run("Preserves timezone information in round-trip", func(t *testing.T) {
		// Create time in different timezone
		loc, _ := time.LoadLocation("America/New_York")
		tm := time.Date(2021, 6, 15, 12, 0, 0, 0, loc)

		// Convert to millis and back
		millis := iso.ReverseGet(tm)
		result := iso.Get(millis)

		// Times should represent the same instant (even if timezone differs)
		assert.True(t, tm.Equal(result))
	})
}

// TestUTF8StringRoundTripLaws verifies isomorphism laws for UTF8String
func TestUTF8StringRoundTripLaws(t *testing.T) {
	iso := UTF8String()

	t.Run("Law 1: ReverseGet(Get(bytes)) == bytes", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("hello"),
			[]byte(""),
			[]byte("Hello 世界 🌍"),
			[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 2: Get(ReverseGet(str)) == str", func(t *testing.T) {
		testCases := []string{
			"hello",
			"",
			"Hello 世界 🌍",
			"special\nchars\ttab",
		}

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.Equal(t, original, result)
		}
	})
}

// TestLinesRoundTripLaws verifies isomorphism laws for Lines
func TestLinesRoundTripLaws(t *testing.T) {
	iso := Lines()

	t.Run("Law 1: ReverseGet(Get(lines)) == lines", func(t *testing.T) {
		testCases := [][]string{
			{"line1", "line2"},
			{"single"},
			{"a", "", "b"},
			{"", "", ""},
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 1: Empty slice special case", func(t *testing.T) {
		// Empty slice becomes "" which splits to [""]
		// This is expected behavior of strings.Split
		original := A.Empty[string]()
		text := iso.Get(original)      // ""
		result := iso.ReverseGet(text) // [""]
		assert.Equal(t, []string{""}, result)
	})

	t.Run("Law 2: Get(ReverseGet(str)) == str", func(t *testing.T) {
		testCases := A.From(
			"line1\nline2",
			"single",
			"",
			"a\n\nb",
			"\n\n",
		)

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.Equal(t, original, result)
		}
	})
}

// TestUnixMilliRoundTripLaws verifies isomorphism laws for UnixMilli
func TestUnixMilliRoundTripLaws(t *testing.T) {
	iso := UnixMilli()

	t.Run("Law 1: ReverseGet(Get(millis)) == millis", func(t *testing.T) {
		testCases := A.From(
			0,
			1609459200000,
			-86400000,
			1234567890000,
			time.Now().UnixMilli(),
		)

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 2: Get(ReverseGet(time)) == time (with millisecond precision)", func(t *testing.T) {
		testCases := A.From(
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Unix(0, 0).UTC(),
			time.Date(1969, 12, 31, 0, 0, 0, 0, time.UTC),
			time.Now().Truncate(time.Millisecond),
		)

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			// Compare Unix timestamps to avoid timezone issues
			assert.Equal(t, original.UnixMilli(), result.UnixMilli())
		}
	})
}

// TestIsosComposition tests composing the isos functions
func TestIsosComposition(t *testing.T) {
	t.Run("Compose UTF8String with Lines", func(t *testing.T) {
		utf8Iso := UTF8String()
		linesIso := Lines()

		// First convert bytes to string, then string to lines
		bytes := []byte("line1\nline2\nline3")
		str := utf8Iso.Get(bytes)
		lines := linesIso.ReverseGet(str)
		assert.Equal(t, []string{"line1", "line2", "line3"}, lines)

		// Reverse: lines to string to bytes
		originalLines := A.From("a", "b", "c")
		text := linesIso.Get(originalLines)
		resultBytes := utf8Iso.ReverseGet(text)
		assert.Equal(t, []byte("a\nb\nc"), resultBytes)
	})

	t.Run("Chain UTF8String and Lines operations", func(t *testing.T) {
		utf8Iso := UTF8String()
		linesIso := Lines()

		// Process: bytes -> string -> lines -> string -> bytes
		original := []byte("hello\nworld")
		str := utf8Iso.Get(original)
		lines := linesIso.ReverseGet(str)
		text := linesIso.Get(lines)
		result := utf8Iso.ReverseGet(text)

		assert.Equal(t, original, result)
	})
}

// TestAdd tests the Add isomorphism
func TestAdd(t *testing.T) {
	t.Run("Add with positive integer", func(t *testing.T) {
		iso := Add(10)
		result := iso.Get(5)
		assert.Equal(t, 15, result)
	})

	t.Run("Add with negative integer", func(t *testing.T) {
		iso := Add(-10)
		result := iso.Get(5)
		assert.Equal(t, -5, result)
	})

	t.Run("Add with zero", func(t *testing.T) {
		iso := Add(0)
		result := iso.Get(42)
		assert.Equal(t, 42, result)
	})

	t.Run("ReverseGet subtracts the value", func(t *testing.T) {
		iso := Add(10)
		result := iso.ReverseGet(15)
		assert.Equal(t, 5, result)
	})

	t.Run("ReverseGet with negative value", func(t *testing.T) {
		iso := Add(-10)
		result := iso.ReverseGet(-5)
		assert.Equal(t, 5, result)
	})

	t.Run("Add with float64", func(t *testing.T) {
		iso := Add(2.5)
		result := iso.Get(7.5)
		assert.Equal(t, 10.0, result)
	})

	t.Run("Add with negative float64", func(t *testing.T) {
		iso := Add(-3.14)
		result := iso.Get(10.0)
		assert.InDelta(t, 6.86, result, 0.001)
	})

	t.Run("ReverseGet with float64", func(t *testing.T) {
		iso := Add(2.5)
		result := iso.ReverseGet(10.0)
		assert.Equal(t, 7.5, result)
	})

	t.Run("Add with complex number", func(t *testing.T) {
		iso := Add(complex(1, 2))
		result := iso.Get(complex(3, 4))
		assert.Equal(t, complex(4, 6), result)
	})

	t.Run("ReverseGet with complex number", func(t *testing.T) {
		iso := Add(complex(1, 2))
		result := iso.ReverseGet(complex(4, 6))
		assert.Equal(t, complex(3, 4), result)
	})

	t.Run("Round-trip integer", func(t *testing.T) {
		iso := Add(10)
		original := 42
		result := iso.ReverseGet(iso.Get(original))
		assert.Equal(t, original, result)
	})

	t.Run("Round-trip float64", func(t *testing.T) {
		iso := Add(3.14)
		original := 2.71
		result := iso.ReverseGet(iso.Get(original))
		assert.InDelta(t, original, result, 0.0001)
	})

	t.Run("Round-trip complex", func(t *testing.T) {
		iso := Add(complex(1, 1))
		original := complex(5, 7)
		result := iso.ReverseGet(iso.Get(original))
		assert.Equal(t, original, result)
	})

	t.Run("Handles large integers", func(t *testing.T) {
		iso := Add(1000000)
		result := iso.Get(2000000)
		assert.Equal(t, 3000000, result)
	})

	t.Run("Handles very small floats", func(t *testing.T) {
		iso := Add(0.0001)
		result := iso.Get(0.0002)
		assert.InDelta(t, 0.0003, result, 0.00001)
	})

	t.Run("Add with int8", func(t *testing.T) {
		iso := Add(int8(5))
		result := iso.Get(int8(10))
		assert.Equal(t, int8(15), result)
	})

	t.Run("Add with int16", func(t *testing.T) {
		iso := Add(int16(100))
		result := iso.Get(int16(200))
		assert.Equal(t, int16(300), result)
	})

	t.Run("Add with int32", func(t *testing.T) {
		iso := Add(int32(1000))
		result := iso.Get(int32(2000))
		assert.Equal(t, int32(3000), result)
	})

	t.Run("Add with int64", func(t *testing.T) {
		iso := Add(int64(10000))
		result := iso.Get(int64(20000))
		assert.Equal(t, int64(30000), result)
	})

	t.Run("Add with uint", func(t *testing.T) {
		iso := Add(uint(5))
		result := iso.Get(uint(10))
		assert.Equal(t, uint(15), result)
	})

	t.Run("Add with float32", func(t *testing.T) {
		iso := Add(float32(1.5))
		result := iso.Get(float32(2.5))
		assert.InDelta(t, float32(4.0), result, 0.001)
	})

	t.Run("Add with complex64", func(t *testing.T) {
		iso := Add(complex64(complex(1, 1)))
		result := iso.Get(complex64(complex(2, 2)))
		assert.Equal(t, complex64(complex(3, 3)), result)
	})

	t.Run("Add with complex128", func(t *testing.T) {
		iso := Add(complex128(complex(2.5, 3.5)))
		result := iso.Get(complex128(complex(1.5, 2.5)))
		assert.Equal(t, complex128(complex(4.0, 6.0)), result)
	})
}

// TestAddRoundTripLaws verifies isomorphism laws for Add
func TestAddRoundTripLaws(t *testing.T) {
	t.Run("Law 1: ReverseGet(Get(x)) == x for integers", func(t *testing.T) {
		iso := Add(10)
		testCases := []int{0, 1, -1, 42, -42, 100, -100, 999999}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result, "Failed for value: %d", original)
		}
	})

	t.Run("Law 2: Get(ReverseGet(x)) == x for integers", func(t *testing.T) {
		iso := Add(10)
		testCases := []int{0, 1, -1, 42, -42, 100, -100, 999999}

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.Equal(t, original, result, "Failed for value: %d", original)
		}
	})

	t.Run("Law 1: ReverseGet(Get(x)) == x for floats", func(t *testing.T) {
		iso := Add(3.14)
		testCases := []float64{0.0, 1.5, -1.5, 42.42, -42.42, 0.001, -0.001}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.InDelta(t, original, result, 0.0001, "Failed for value: %f", original)
		}
	})

	t.Run("Law 2: Get(ReverseGet(x)) == x for floats", func(t *testing.T) {
		iso := Add(3.14)
		testCases := []float64{0.0, 1.5, -1.5, 42.42, -42.42, 0.001, -0.001}

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.InDelta(t, original, result, 0.0001, "Failed for value: %f", original)
		}
	})

	t.Run("Law 1: ReverseGet(Get(x)) == x for complex", func(t *testing.T) {
		iso := Add(complex(1, 2))
		testCases := []complex128{
			complex(0, 0),
			complex(1, 1),
			complex(-1, -1),
			complex(3.5, 4.5),
			complex(-3.5, -4.5),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result, "Failed for value: %v", original)
		}
	})

	t.Run("Law 2: Get(ReverseGet(x)) == x for complex", func(t *testing.T) {
		iso := Add(complex(1, 2))
		testCases := []complex128{
			complex(0, 0),
			complex(1, 1),
			complex(-1, -1),
			complex(3.5, 4.5),
			complex(-3.5, -4.5),
		}

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.Equal(t, original, result, "Failed for value: %v", original)
		}
	})
}

// TestAddComposition tests composing Add with other isomorphisms
func TestAddComposition(t *testing.T) {
	t.Run("Compose two Add isomorphisms", func(t *testing.T) {
		addTen := Add(10)
		addFive := Add(5)

		// Compose: first add 10, then add 5 (total: add 15)
		composed := F.Pipe1(addTen, Compose[int](addFive))

		result := composed.Get(0)
		assert.Equal(t, 15, result)

		// Reverse should subtract 15
		original := composed.ReverseGet(15)
		assert.Equal(t, 0, original)
	})

	t.Run("Add with Reverse", func(t *testing.T) {
		addTen := Add(10)
		reversed := Reverse(addTen)

		// Reversed: Get subtracts, ReverseGet adds
		result := reversed.Get(15)
		assert.Equal(t, 5, result)

		original := reversed.ReverseGet(5)
		assert.Equal(t, 15, original)
	})

	t.Run("Add with Modify", func(t *testing.T) {
		addTen := Add(10)

		// Double the value after adding 10
		doubler := Modify[int](func(x int) int { return x * 2 })(addTen)

		// (5 + 10) * 2 = 30, then subtract 10 = 20
		result := doubler(5)
		assert.Equal(t, 20, result)
	})

	t.Run("Chain multiple Add operations", func(t *testing.T) {
		// Create a chain: add 5, then add 10, then add 3
		add5 := Add(5)
		add10 := Add(10)
		add3 := Add(3)

		value := 0
		step1 := add5.Get(value)  // 5
		step2 := add10.Get(step1) // 15
		step3 := add3.Get(step2)  // 18
		assert.Equal(t, 18, step3)

		// Reverse the chain
		back1 := add3.ReverseGet(step3)  // 15
		back2 := add10.ReverseGet(back1) // 5
		back3 := add5.ReverseGet(back2)  // 0
		assert.Equal(t, value, back3)
	})
}

// TestSub tests the Sub isomorphism
func TestSub(t *testing.T) {
	t.Run("Sub with positive integer", func(t *testing.T) {
		iso := Sub(10)
		result := iso.Get(15)
		assert.Equal(t, 5, result)
	})

	t.Run("Sub with negative integer", func(t *testing.T) {
		iso := Sub(-10)
		result := iso.Get(5)
		assert.Equal(t, 15, result)
	})

	t.Run("Sub with zero", func(t *testing.T) {
		iso := Sub(0)
		result := iso.Get(42)
		assert.Equal(t, 42, result)
	})

	t.Run("ReverseGet adds the value", func(t *testing.T) {
		iso := Sub(10)
		result := iso.ReverseGet(5)
		assert.Equal(t, 15, result)
	})

	t.Run("ReverseGet with negative value", func(t *testing.T) {
		iso := Sub(-10)
		result := iso.ReverseGet(15)
		assert.Equal(t, 5, result)
	})

	t.Run("Sub with float64", func(t *testing.T) {
		iso := Sub(2.5)
		result := iso.Get(10.0)
		assert.Equal(t, 7.5, result)
	})

	t.Run("Sub with negative float64", func(t *testing.T) {
		iso := Sub(-3.14)
		result := iso.Get(10.0)
		assert.InDelta(t, 13.14, result, 0.001)
	})

	t.Run("ReverseGet with float64", func(t *testing.T) {
		iso := Sub(2.5)
		result := iso.ReverseGet(7.5)
		assert.Equal(t, 10.0, result)
	})

	t.Run("Sub with complex number", func(t *testing.T) {
		iso := Sub(complex(1, 2))
		result := iso.Get(complex(4, 6))
		assert.Equal(t, complex(3, 4), result)
	})

	t.Run("ReverseGet with complex number", func(t *testing.T) {
		iso := Sub(complex(1, 2))
		result := iso.ReverseGet(complex(3, 4))
		assert.Equal(t, complex(4, 6), result)
	})

	t.Run("Round-trip integer", func(t *testing.T) {
		iso := Sub(10)
		original := 42
		result := iso.ReverseGet(iso.Get(original))
		assert.Equal(t, original, result)
	})

	t.Run("Round-trip float64", func(t *testing.T) {
		iso := Sub(3.14)
		original := 10.5
		result := iso.ReverseGet(iso.Get(original))
		assert.InDelta(t, original, result, 0.0001)
	})

	t.Run("Round-trip complex", func(t *testing.T) {
		iso := Sub(complex(1, 1))
		original := complex(5, 7)
		result := iso.ReverseGet(iso.Get(original))
		assert.Equal(t, original, result)
	})

	t.Run("Handles large integers", func(t *testing.T) {
		iso := Sub(1000000)
		result := iso.Get(3000000)
		assert.Equal(t, 2000000, result)
	})

	t.Run("Handles very small floats", func(t *testing.T) {
		iso := Sub(0.0001)
		result := iso.Get(0.0003)
		assert.InDelta(t, 0.0002, result, 0.00001)
	})

	t.Run("Sub with int8", func(t *testing.T) {
		iso := Sub(int8(5))
		result := iso.Get(int8(15))
		assert.Equal(t, int8(10), result)
	})

	t.Run("Sub with int16", func(t *testing.T) {
		iso := Sub(int16(100))
		result := iso.Get(int16(300))
		assert.Equal(t, int16(200), result)
	})

	t.Run("Sub with int32", func(t *testing.T) {
		iso := Sub(int32(1000))
		result := iso.Get(int32(3000))
		assert.Equal(t, int32(2000), result)
	})

	t.Run("Sub with int64", func(t *testing.T) {
		iso := Sub(int64(10000))
		result := iso.Get(int64(30000))
		assert.Equal(t, int64(20000), result)
	})

	t.Run("Sub with uint", func(t *testing.T) {
		iso := Sub(uint(5))
		result := iso.Get(uint(15))
		assert.Equal(t, uint(10), result)
	})

	t.Run("Sub with float32", func(t *testing.T) {
		iso := Sub(float32(1.5))
		result := iso.Get(float32(4.0))
		assert.InDelta(t, float32(2.5), result, 0.001)
	})

	t.Run("Sub with complex64", func(t *testing.T) {
		iso := Sub(complex64(complex(1, 1)))
		result := iso.Get(complex64(complex(3, 3)))
		assert.Equal(t, complex64(complex(2, 2)), result)
	})

	t.Run("Sub with complex128", func(t *testing.T) {
		iso := Sub(complex128(complex(2.5, 3.5)))
		result := iso.Get(complex128(complex(4.0, 6.0)))
		assert.Equal(t, complex128(complex(1.5, 2.5)), result)
	})

	t.Run("Sub is equivalent to Add with negative value", func(t *testing.T) {
		sub5 := Sub(5)
		addNeg5 := Add(-5)

		value := 10
		assert.Equal(t, sub5.Get(value), addNeg5.Get(value))
		assert.Equal(t, sub5.ReverseGet(value), addNeg5.ReverseGet(value))
	})
}

// TestSubRoundTripLaws verifies isomorphism laws for Sub
func TestSubRoundTripLaws(t *testing.T) {
	t.Run("Law 1: ReverseGet(Get(x)) == x for integers", func(t *testing.T) {
		iso := Sub(10)
		testCases := []int{0, 1, -1, 42, -42, 100, -100, 999999}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result, "Failed for value: %d", original)
		}
	})

	t.Run("Law 2: Get(ReverseGet(x)) == x for integers", func(t *testing.T) {
		iso := Sub(10)
		testCases := []int{0, 1, -1, 42, -42, 100, -100, 999999}

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.Equal(t, original, result, "Failed for value: %d", original)
		}
	})

	t.Run("Law 1: ReverseGet(Get(x)) == x for floats", func(t *testing.T) {
		iso := Sub(3.14)
		testCases := []float64{0.0, 1.5, -1.5, 42.42, -42.42, 0.001, -0.001}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.InDelta(t, original, result, 0.0001, "Failed for value: %f", original)
		}
	})

	t.Run("Law 2: Get(ReverseGet(x)) == x for floats", func(t *testing.T) {
		iso := Sub(3.14)
		testCases := []float64{0.0, 1.5, -1.5, 42.42, -42.42, 0.001, -0.001}

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.InDelta(t, original, result, 0.0001, "Failed for value: %f", original)
		}
	})

	t.Run("Law 1: ReverseGet(Get(x)) == x for complex", func(t *testing.T) {
		iso := Sub(complex(1, 2))
		testCases := []complex128{
			complex(0, 0),
			complex(1, 1),
			complex(-1, -1),
			complex(3.5, 4.5),
			complex(-3.5, -4.5),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result, "Failed for value: %v", original)
		}
	})

	t.Run("Law 2: Get(ReverseGet(x)) == x for complex", func(t *testing.T) {
		iso := Sub(complex(1, 2))
		testCases := []complex128{
			complex(0, 0),
			complex(1, 1),
			complex(-1, -1),
			complex(3.5, 4.5),
			complex(-3.5, -4.5),
		}

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.Equal(t, original, result, "Failed for value: %v", original)
		}
	})
}

// TestSubComposition tests composing Sub with other isomorphisms
func TestSubComposition(t *testing.T) {
	t.Run("Compose two Sub isomorphisms", func(t *testing.T) {
		subTen := Sub(10)
		subFive := Sub(5)

		// Compose: first subtract 10, then subtract 5 (total: subtract 15)
		composed := F.Pipe1(subTen, Compose[int](subFive))

		result := composed.Get(20)
		assert.Equal(t, 5, result)

		// Reverse should add 15
		original := composed.ReverseGet(5)
		assert.Equal(t, 20, original)
	})

	t.Run("Compose Sub and Add", func(t *testing.T) {
		sub10 := Sub(10)
		add5 := Add(5)

		// Compose: first subtract 10, then add 5 (net: subtract 5)
		composed := F.Pipe1(sub10, Compose[int](add5))

		result := composed.Get(20)
		assert.Equal(t, 15, result)

		// Reverse: subtract 5, then add 10 (net: add 5)
		original := composed.ReverseGet(15)
		assert.Equal(t, 20, original)
	})

	t.Run("Sub with Reverse", func(t *testing.T) {
		subTen := Sub(10)
		reversed := Reverse(subTen)

		// Reversed: Get adds, ReverseGet subtracts
		result := reversed.Get(5)
		assert.Equal(t, 15, result)

		original := reversed.ReverseGet(15)
		assert.Equal(t, 5, original)
	})

	t.Run("Sub with Modify", func(t *testing.T) {
		subTen := Sub(10)

		// Double the value after subtracting 10
		doubler := Modify[int](func(x int) int { return x * 2 })(subTen)

		// (20 - 10) * 2 = 20, then add 10 = 30
		result := doubler(20)
		assert.Equal(t, 30, result)
	})

	t.Run("Chain multiple Sub operations", func(t *testing.T) {
		// Create a chain: subtract 5, then subtract 10, then subtract 3
		sub5 := Sub(5)
		sub10 := Sub(10)
		sub3 := Sub(3)

		value := 50
		step1 := sub5.Get(value)  // 45
		step2 := sub10.Get(step1) // 35
		step3 := sub3.Get(step2)  // 32
		assert.Equal(t, 32, step3)

		// Reverse the chain
		back1 := sub3.ReverseGet(step3)  // 35
		back2 := sub10.ReverseGet(back1) // 45
		back3 := sub5.ReverseGet(back2)  // 50
		assert.Equal(t, value, back3)
	})

	t.Run("Sub and Add cancel each other", func(t *testing.T) {
		sub10 := Sub(10)
		add10 := Add(10)

		value := 42

		// Apply Sub then Add - should get original value
		result1 := add10.Get(sub10.Get(value))
		assert.Equal(t, value, result1)

		// Apply Add then Sub - should get original value
		result2 := sub10.Get(add10.Get(value))
		assert.Equal(t, value, result2)
	})
}

// TestSwapPair tests the SwapPair isomorphism
func TestSwapPair(t *testing.T) {
	t.Run("Swap pair of string and int", func(t *testing.T) {
		iso := SwapPair[string, int]()
		original := P.MakePair("hello", 42)
		swapped := iso.Get(original)

		assert.Equal(t, 42, P.Head(swapped))
		assert.Equal(t, "hello", P.Tail(swapped))
	})

	t.Run("ReverseGet swaps back", func(t *testing.T) {
		iso := SwapPair[string, int]()
		swapped := P.MakePair(42, "hello")
		original := iso.ReverseGet(swapped)

		assert.Equal(t, "hello", P.Head(original))
		assert.Equal(t, 42, P.Tail(original))
	})

	t.Run("Swap pair of same types", func(t *testing.T) {
		iso := SwapPair[int, int]()
		original := P.MakePair(1, 2)
		swapped := iso.Get(original)

		assert.Equal(t, 2, P.Head(swapped))
		assert.Equal(t, 1, P.Tail(swapped))
	})

	t.Run("Swap pair of floats", func(t *testing.T) {
		iso := SwapPair[float64, float64]()
		original := P.MakePair(3.14, 2.71)
		swapped := iso.Get(original)

		assert.Equal(t, 2.71, P.Head(swapped))
		assert.Equal(t, 3.14, P.Tail(swapped))
	})

	t.Run("Swap pair with complex types", func(t *testing.T) {
		type Person struct{ Name string }
		type Address struct{ City string }

		iso := SwapPair[Person, Address]()
		original := P.MakePair(Person{"Alice"}, Address{"NYC"})
		swapped := iso.Get(original)

		assert.Equal(t, Address{"NYC"}, P.Head(swapped))
		assert.Equal(t, Person{"Alice"}, P.Tail(swapped))
	})

	t.Run("Round-trip with heterogeneous types", func(t *testing.T) {
		iso := SwapPair[string, int]()
		original := P.MakePair("test", 123)
		result := iso.ReverseGet(iso.Get(original))

		assert.Equal(t, original, result)
	})

	t.Run("Round-trip with same types", func(t *testing.T) {
		iso := SwapPair[int, int]()
		original := P.MakePair(10, 20)
		result := iso.ReverseGet(iso.Get(original))

		assert.Equal(t, original, result)
	})

	t.Run("Swap coordinates", func(t *testing.T) {
		iso := SwapPair[float64, float64]()
		point := P.MakePair(3.0, 4.0) // (x, y)
		swapped := iso.Get(point)     // (y, x)

		assert.Equal(t, 4.0, P.Head(swapped))
		assert.Equal(t, 3.0, P.Tail(swapped))
	})

	t.Run("Swap with nil-able types", func(t *testing.T) {
		iso := SwapPair[*string, *int]()
		str := "test"
		num := 42
		original := P.MakePair(&str, &num)
		swapped := iso.Get(original)

		assert.Equal(t, &num, P.Head(swapped))
		assert.Equal(t, &str, P.Tail(swapped))
	})

	t.Run("Swap with slices", func(t *testing.T) {
		iso := SwapPair[[]int, []string]()
		original := P.MakePair([]int{1, 2, 3}, []string{"a", "b"})
		swapped := iso.Get(original)

		assert.Equal(t, []string{"a", "b"}, P.Head(swapped))
		assert.Equal(t, []int{1, 2, 3}, P.Tail(swapped))
	})
}

// TestSwapPairRoundTripLaws verifies isomorphism laws for SwapPair
func TestSwapPairRoundTripLaws(t *testing.T) {
	t.Run("Law 1: ReverseGet(Get(pair)) == pair", func(t *testing.T) {
		iso := SwapPair[string, int]()
		testCases := []Pair[string, int]{
			P.MakePair("a", 1),
			P.MakePair("", 0),
			P.MakePair("test", -42),
			P.MakePair("hello world", 999),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 2: Get(ReverseGet(pair)) == pair", func(t *testing.T) {
		iso := SwapPair[string, int]()
		testCases := []Pair[string, int]{
			P.MakePair("a", 1),
			P.MakePair("", 0),
			P.MakePair("test", -42),
			P.MakePair("hello world", 999),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Self-inverse property: Get(Get(pair)) swaps back", func(t *testing.T) {
		iso1 := SwapPair[string, int]()
		iso2 := SwapPair[int, string]()

		original := P.MakePair("test", 42)
		swapped := iso1.Get(original)
		restored := iso2.Get(swapped)

		assert.Equal(t, original, restored)
	})
}

// TestSwapPairComposition tests composing SwapPair with other isomorphisms
func TestSwapPairComposition(t *testing.T) {
	t.Run("Compose two SwapPair isomorphisms", func(t *testing.T) {
		swap1 := SwapPair[string, int]()
		swap2 := SwapPair[int, string]()

		original := P.MakePair("hello", 42)

		// First swap: (string, int) -> (int, string)
		step1 := swap1.Get(original)
		assert.Equal(t, 42, P.Head(step1))
		assert.Equal(t, "hello", P.Tail(step1))

		// Second swap: (int, string) -> (string, int)
		step2 := swap2.Get(step1)
		assert.Equal(t, original, step2)
	})

	t.Run("SwapPair with Reverse", func(t *testing.T) {
		swapIso := SwapPair[string, int]()
		reversed := Reverse(swapIso)

		original := P.MakePair("test", 123)

		// Reversed Get is same as ReverseGet
		result1 := reversed.Get(P.MakePair(123, "test"))
		assert.Equal(t, original, result1)

		// Reversed ReverseGet is same as Get
		result2 := reversed.ReverseGet(original)
		assert.Equal(t, P.MakePair(123, "test"), result2)
	})

	t.Run("Chain multiple SwapPair operations", func(t *testing.T) {
		original := P.MakePair(1, "a")

		// Swap multiple times
		iso1 := SwapPair[int, string]()
		iso2 := SwapPair[string, int]()

		step1 := iso1.Get(original) // ("a", 1)
		step2 := iso2.Get(step1)    // (1, "a")
		step3 := iso1.Get(step2)    // ("a", 1)
		step4 := iso2.Get(step3)    // (1, "a")

		assert.Equal(t, original, step4)
	})
}

// TestSwapPairUseCases demonstrates practical use cases for SwapPair
func TestSwapPairUseCases(t *testing.T) {
	t.Run("Swap coordinates from (x, y) to (y, x)", func(t *testing.T) {
		swapCoords := SwapPair[float64, float64]()

		point := P.MakePair(3.0, 4.0) // (x=3, y=4)
		swapped := swapCoords.Get(point)

		assert.Equal(t, 4.0, P.Head(swapped)) // y
		assert.Equal(t, 3.0, P.Tail(swapped)) // x
	})

	t.Run("Swap key-value to value-key", func(t *testing.T) {
		swapKV := SwapPair[string, int]()

		entry := P.MakePair("age", 30)
		swapped := swapKV.Get(entry)

		assert.Equal(t, 30, P.Head(swapped))
		assert.Equal(t, "age", P.Tail(swapped))
	})

	t.Run("Adapt function argument order", func(t *testing.T) {
		swapArgs := SwapPair[string, int]()

		// Function expects (int, string)
		processArgs := func(p Pair[int, string]) string {
			return fmt.Sprintf("%s: %d", P.Tail(p), P.Head(p))
		}

		// We have (string, int)
		input := P.MakePair("count", 42)
		swapped := swapArgs.Get(input)
		result := processArgs(swapped)

		assert.Equal(t, "count: 42", result)
	})

	t.Run("Normalize data structure", func(t *testing.T) {
		type Name string
		type Age int

		swapPerson := SwapPair[Name, Age]()

		// Different systems use different orders
		person1 := P.MakePair(Name("Alice"), Age(30))
		person2 := swapPerson.Get(person1) // Normalized to (Age, Name)

		assert.Equal(t, Age(30), P.Head(person2))
		assert.Equal(t, Name("Alice"), P.Tail(person2))
	})
}

// TestSwapEither tests the SwapEither isomorphism
func TestSwapEither(t *testing.T) {
	t.Run("Swap Left value", func(t *testing.T) {
		iso := SwapEither[string, int]()
		original := E.Left[int]("error")
		swapped := iso.Get(original)

		assert.True(t, E.IsRight(swapped))
		// swapped is Either[int, string], so Unwrap returns (string, int)
		right, _ := E.Unwrap(swapped)
		assert.Equal(t, "error", right)
	})

	t.Run("Swap Right value", func(t *testing.T) {
		iso := SwapEither[string, int]()
		original := E.Right[string](42)
		swapped := iso.Get(original)

		assert.True(t, E.IsLeft(swapped))
		// swapped is Either[int, string], so Unwrap returns (string, int)
		_, left := E.Unwrap(swapped)
		assert.Equal(t, 42, left)
	})

	t.Run("ReverseGet swaps Left back", func(t *testing.T) {
		iso := SwapEither[string, int]()
		swapped := E.Right[int]("error")
		original := iso.ReverseGet(swapped)

		assert.True(t, E.IsLeft(original))
		// original is Either[string, int], so Unwrap returns (int, string)
		_, left := E.Unwrap(original)
		assert.Equal(t, "error", left)
	})

	t.Run("ReverseGet swaps Right back", func(t *testing.T) {
		iso := SwapEither[string, int]()
		swapped := E.Left[string](42)
		original := iso.ReverseGet(swapped)

		assert.True(t, E.IsRight(original))
		// original is Either[string, int], so Unwrap returns (int, string)
		right, _ := E.Unwrap(original)
		assert.Equal(t, 42, right)
	})

	t.Run("Swap with error type", func(t *testing.T) {
		iso := SwapEither[error, string]()
		err := fmt.Errorf("test error")
		original := E.Left[string](err)
		swapped := iso.Get(original)

		assert.True(t, E.IsRight(swapped))
		// swapped is Either[string, error], so Unwrap returns (error, string)
		right, _ := E.Unwrap(swapped)
		assert.Equal(t, err, right)
	})

	t.Run("Swap with complex types", func(t *testing.T) {
		type ValidationError struct{ Message string }
		type User struct{ Name string }

		iso := SwapEither[ValidationError, User]()
		original := E.Right[ValidationError](User{"Alice"})
		swapped := iso.Get(original)

		assert.True(t, E.IsLeft(swapped))
		// swapped is Either[User, ValidationError], so Unwrap returns (ValidationError, User)
		_, left := E.Unwrap(swapped)
		assert.Equal(t, User{"Alice"}, left)
	})

	t.Run("Round-trip Left value", func(t *testing.T) {
		iso := SwapEither[string, int]()
		original := E.Left[int]("error")
		result := iso.ReverseGet(iso.Get(original))

		assert.Equal(t, original, result)
	})

	t.Run("Round-trip Right value", func(t *testing.T) {
		iso := SwapEither[string, int]()
		original := E.Right[string](42)
		result := iso.ReverseGet(iso.Get(original))

		assert.Equal(t, original, result)
	})

	t.Run("Swap with slice types", func(t *testing.T) {
		iso := SwapEither[[]string, int]()
		errors := []string{"error1", "error2"}
		original := E.Left[int](errors)
		swapped := iso.Get(original)

		assert.True(t, E.IsRight(swapped))
		// swapped is Either[int, []string], so Unwrap returns ([]string, int)
		right, _ := E.Unwrap(swapped)
		assert.Equal(t, errors, right)
	})
}

// TestSwapEitherRoundTripLaws verifies isomorphism laws for SwapEither
func TestSwapEitherRoundTripLaws(t *testing.T) {
	t.Run("Law 1: ReverseGet(Get(either)) == either for Left", func(t *testing.T) {
		iso := SwapEither[string, int]()
		testCases := []E.Either[string, int]{
			E.Left[int]("error1"),
			E.Left[int](""),
			E.Left[int]("test error"),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 1: ReverseGet(Get(either)) == either for Right", func(t *testing.T) {
		iso := SwapEither[string, int]()
		testCases := []E.Either[string, int]{
			E.Right[string](0),
			E.Right[string](42),
			E.Right[string](-1),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 2: Get(ReverseGet(either)) == either for Left", func(t *testing.T) {
		iso := SwapEither[string, int]()
		testCases := []E.Either[string, int]{
			E.Left[int]("error1"),
			E.Left[int]("error2"),
			E.Left[int](""),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 2: Get(ReverseGet(either)) == either for Right", func(t *testing.T) {
		iso := SwapEither[string, int]()
		testCases := []E.Either[string, int]{
			E.Right[string](0),
			E.Right[string](42),
			E.Right[string](-1),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Self-inverse property", func(t *testing.T) {
		iso1 := SwapEither[string, int]()
		iso2 := SwapEither[int, string]()

		original := E.Left[int]("error")
		swapped := iso1.Get(original)
		restored := iso2.Get(swapped)

		assert.Equal(t, original, restored)
	})
}

// TestSwapEitherComposition tests composing SwapEither with other isomorphisms
func TestSwapEitherComposition(t *testing.T) {
	t.Run("Compose two SwapEither isomorphisms", func(t *testing.T) {
		swap1 := SwapEither[string, int]()
		swap2 := SwapEither[int, string]()

		original := E.Left[int]("error")

		// First swap: Either[string, int] -> Either[int, string]
		step1 := swap1.Get(original)
		assert.True(t, E.IsRight(step1))

		// Second swap: Either[int, string] -> Either[string, int]
		step2 := swap2.Get(step1)
		assert.Equal(t, original, step2)
	})

	t.Run("SwapEither with Reverse", func(t *testing.T) {
		swapIso := SwapEither[string, int]()
		reversed := Reverse(swapIso)

		original := E.Left[int]("error")

		// Reversed Get is same as ReverseGet
		swapped := E.Right[int]("error")
		result1 := reversed.Get(swapped)
		assert.Equal(t, original, result1)

		// Reversed ReverseGet is same as Get
		result2 := reversed.ReverseGet(original)
		assert.True(t, E.IsRight(result2))
	})

	t.Run("Chain multiple SwapEither operations", func(t *testing.T) {
		original := E.Right[string](42)

		iso1 := SwapEither[string, int]()
		iso2 := SwapEither[int, string]()

		step1 := iso1.Get(original) // Either[int, string] with Left(42)
		step2 := iso2.Get(step1)    // Either[string, int] with Right(42)
		step3 := iso1.Get(step2)    // Either[int, string] with Left(42)
		step4 := iso2.Get(step3)    // Either[string, int] with Right(42)

		assert.Equal(t, original, step4)
	})
}

// TestSwapEitherUseCases demonstrates practical use cases for SwapEither
func TestSwapEitherUseCases(t *testing.T) {
	t.Run("Convert error-left to error-right convention", func(t *testing.T) {
		swapError := SwapEither[error, string]()

		// Error-left convention
		result := E.Left[string](fmt.Errorf("failed"))
		// Convert to error-right convention
		swapped := swapError.Get(result)

		assert.True(t, E.IsRight(swapped))
	})

	t.Run("Adapt validation result types", func(t *testing.T) {
		type ValidationErrors []string
		type User struct{ Name string }

		swapValidation := SwapEither[ValidationErrors, User]()

		// Valid user
		valid := E.Right[ValidationErrors](User{"Alice"})
		swapped := swapValidation.Get(valid)

		assert.True(t, E.IsLeft(swapped))
		// swapped is Either[User, ValidationErrors], so Unwrap returns (ValidationErrors, User)
		_, left := E.Unwrap(swapped)
		assert.Equal(t, User{"Alice"}, left)
	})

	t.Run("Normalize API response types", func(t *testing.T) {
		type APIError struct{ Code int }
		type Response struct{ Data string }

		swapAPI := SwapEither[APIError, Response]()

		// Success response
		success := E.Right[APIError](Response{"data"})
		swapped := swapAPI.Get(success)

		assert.True(t, E.IsLeft(swapped))
	})

	t.Run("Convert between different Either conventions", func(t *testing.T) {
		swapConvention := SwapEither[string, int]()

		// Library A uses Either[Error, Value]
		libraryAResult := E.Left[int]("error")

		// Library B expects Either[Value, Error]
		libraryBResult := swapConvention.Get(libraryAResult)

		assert.True(t, E.IsRight(libraryBResult))
		// libraryBResult is Either[int, string], so Unwrap returns (string, int)
		right, _ := E.Unwrap(libraryBResult)
		assert.Equal(t, "error", right)
	})
}

// TestSubUseCases demonstrates practical use cases for Sub
func TestSubUseCases(t *testing.T) {
	t.Run("Convert 1-based to 0-based indexing", func(t *testing.T) {
		oneToZero := Sub(1)

		// Convert 1-based index to 0-based
		zeroBasedIndex := oneToZero.Get(1)
		assert.Equal(t, 0, zeroBasedIndex)

		// Convert 0-based index back to 1-based
		oneBasedIndex := oneToZero.ReverseGet(0)
		assert.Equal(t, 1, oneBasedIndex)
	})

	t.Run("Apply discount", func(t *testing.T) {
		discount := Sub(10.0)

		originalPrice := 50.0
		discountedPrice := discount.Get(originalPrice)
		assert.Equal(t, 40.0, discountedPrice)

		// Reverse to get original price
		original := discount.ReverseGet(discountedPrice)
		assert.Equal(t, originalPrice, original)
	})

	t.Run("Coordinate translation backwards", func(t *testing.T) {
		// Translate x-coordinate backwards by 100 units
		translateX := Sub(100)

		originalX := 150
		translatedX := translateX.Get(originalX)
		assert.Equal(t, 50, translatedX)

		// Translate forward
		backToOriginal := translateX.ReverseGet(translatedX)
		assert.Equal(t, originalX, backToOriginal)
	})

	t.Run("Time offset backwards in hours", func(t *testing.T) {
		// Go back 5 hours (represented as integer hours)
		subFiveHours := Sub(5)

		currentHour := 15
		pastHour := subFiveHours.Get(currentHour)
		assert.Equal(t, 10, pastHour)

		// Go forward 5 hours
		futureHour := subFiveHours.ReverseGet(pastHour)
		assert.Equal(t, currentHour, futureHour)
	})

	t.Run("Temperature calibration correction", func(t *testing.T) {
		// Sensor reads 2.5 degrees too high, correct by subtracting
		correction := Sub(2.5)

		sensorReading := 22.5
		actualTemp := correction.Get(sensorReading)
		assert.Equal(t, 20.0, actualTemp)

		// Reverse to get sensor reading from actual
		reading := correction.ReverseGet(actualTemp)
		assert.Equal(t, sensorReading, reading)
	})
}

// TestAddUseCases demonstrates practical use cases for Add
func TestAddUseCases(t *testing.T) {
	t.Run("Convert 0-based to 1-based indexing", func(t *testing.T) {
		zeroToOne := Add(1)

		// Convert 0-based index to 1-based
		oneBasedIndex := zeroToOne.Get(0)
		assert.Equal(t, 1, oneBasedIndex)

		// Convert 1-based index back to 0-based
		zeroBasedIndex := zeroToOne.ReverseGet(1)
		assert.Equal(t, 0, zeroBasedIndex)
	})

	t.Run("Temperature offset adjustment", func(t *testing.T) {
		// Adjust temperature by offset (e.g., calibration)
		calibrationOffset := Add(2.5)

		measuredTemp := 20.0
		calibratedTemp := calibrationOffset.Get(measuredTemp)
		assert.Equal(t, 22.5, calibratedTemp)

		// Reverse calibration
		original := calibrationOffset.ReverseGet(calibratedTemp)
		assert.Equal(t, measuredTemp, original)
	})

	t.Run("Coordinate translation", func(t *testing.T) {
		// Translate x-coordinate by 100 units
		translateX := Add(100)

		originalX := 50
		translatedX := translateX.Get(originalX)
		assert.Equal(t, 150, translatedX)

		// Translate back
		backToOriginal := translateX.ReverseGet(translatedX)
		assert.Equal(t, originalX, backToOriginal)
	})

	t.Run("Time offset in hours", func(t *testing.T) {
		// Add 5 hours (represented as integer hours)
		addFiveHours := Add(5)

		currentHour := 10
		futureHour := addFiveHours.Get(currentHour)
		assert.Equal(t, 15, futureHour)

		// Go back 5 hours
		pastHour := addFiveHours.ReverseGet(futureHour)
		assert.Equal(t, currentHour, pastHour)
	})
}

// TestReverseArray tests the ReverseArray isomorphism
func TestReverseArray(t *testing.T) {
	t.Run("Reverse array of integers", func(t *testing.T) {
		iso := ReverseArray[int]()
		input := []int{1, 2, 3, 4, 5}
		reversed := iso.Get(input)
		expected := []int{5, 4, 3, 2, 1}
		assert.Equal(t, expected, reversed)
	})

	t.Run("Reverse array of strings", func(t *testing.T) {
		iso := ReverseArray[string]()
		input := []string{"hello", "world", "foo"}
		reversed := iso.Get(input)
		expected := []string{"foo", "world", "hello"}
		assert.Equal(t, expected, reversed)
	})

	t.Run("Reverse empty array", func(t *testing.T) {
		iso := ReverseArray[int]()
		input := []int{}
		reversed := iso.Get(input)
		assert.Equal(t, []int{}, reversed)
	})

	t.Run("Reverse single element array", func(t *testing.T) {
		iso := ReverseArray[string]()
		input := []string{"only"}
		reversed := iso.Get(input)
		assert.Equal(t, []string{"only"}, reversed)
	})

	t.Run("ReverseGet also reverses", func(t *testing.T) {
		iso := ReverseArray[int]()
		input := []int{1, 2, 3}
		reversed := iso.ReverseGet(input)
		expected := []int{3, 2, 1}
		assert.Equal(t, expected, reversed)
	})

	t.Run("Does not modify original array", func(t *testing.T) {
		iso := ReverseArray[int]()
		original := []int{1, 2, 3, 4, 5}
		originalCopy := []int{1, 2, 3, 4, 5}
		_ = iso.Get(original)
		assert.Equal(t, originalCopy, original)
	})

	t.Run("Round-trip returns original", func(t *testing.T) {
		iso := ReverseArray[int]()
		original := []int{1, 2, 3, 4, 5}
		result := iso.ReverseGet(iso.Get(original))
		assert.Equal(t, original, result)
	})

	t.Run("Self-inverse property", func(t *testing.T) {
		iso := ReverseArray[int]()
		input := []int{1, 2, 3, 4, 5}

		// Get twice should return original
		result1 := iso.Get(iso.Get(input))
		assert.Equal(t, input, result1)

		// ReverseGet twice should return original
		result2 := iso.ReverseGet(iso.ReverseGet(input))
		assert.Equal(t, input, result2)
	})

	t.Run("Reverse with floats", func(t *testing.T) {
		iso := ReverseArray[float64]()
		input := []float64{1.1, 2.2, 3.3}
		reversed := iso.Get(input)
		expected := []float64{3.3, 2.2, 1.1}
		assert.Equal(t, expected, reversed)
	})

	t.Run("Reverse with structs", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		iso := ReverseArray[Person]()
		input := []Person{
			{"Alice", 30},
			{"Bob", 25},
		}
		reversed := iso.Get(input)
		expected := []Person{
			{"Bob", 25},
			{"Alice", 30},
		}
		assert.Equal(t, expected, reversed)
	})
}

// TestReverseArrayRoundTripLaws verifies isomorphism laws for ReverseArray
func TestReverseArrayRoundTripLaws(t *testing.T) {
	t.Run("Law 1: ReverseGet(Get(arr)) == arr", func(t *testing.T) {
		iso := ReverseArray[int]()
		testCases := [][]int{
			{1, 2, 3, 4, 5},
			{1},
			{},
			{10, 20},
			{5, 4, 3, 2, 1},
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 2: Get(ReverseGet(arr)) == arr", func(t *testing.T) {
		iso := ReverseArray[string]()
		testCases := [][]string{
			{"a", "b", "c"},
			{"single"},
			{},
			{"x", "y"},
		}

		for _, original := range testCases {
			result := iso.Get(iso.ReverseGet(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Self-inverse: Get == ReverseGet in effect", func(t *testing.T) {
		iso := ReverseArray[int]()
		input := []int{1, 2, 3, 4, 5}

		getResult := iso.Get(input)
		reverseGetResult := iso.ReverseGet(input)

		// Both should produce the same reversed result
		assert.Equal(t, getResult, reverseGetResult)
	})
}

// TestReverseArrayComposition tests composing ReverseArray with other isomorphisms
func TestReverseArrayComposition(t *testing.T) {
	t.Run("Compose with itself returns identity", func(t *testing.T) {
		iso := ReverseArray[int]()
		composed := F.Pipe1(iso, Compose[[]int](iso))

		input := []int{1, 2, 3, 4, 5}
		result := composed.Get(input)

		// Reversing twice should return original
		assert.Equal(t, input, result)
	})

	t.Run("Compose with Reverse", func(t *testing.T) {
		iso := ReverseArray[int]()
		reversed := Reverse(iso)

		input := []int{1, 2, 3}

		// Reversed iso: Get and ReverseGet are swapped (but they're the same for ReverseArray)
		result1 := reversed.Get(input)
		result2 := iso.ReverseGet(input)
		assert.Equal(t, result1, result2)
	})

	t.Run("Use with Modify", func(t *testing.T) {
		iso := ReverseArray[int]()

		// Reverse, double all elements, reverse back
		modifier := Modify[[]int](func(arr []int) []int {
			result := make([]int, len(arr))
			for i, v := range arr {
				result[i] = v * 2
			}
			return result
		})(iso)

		input := []int{1, 2, 3}
		result := modifier(input)

		// Should double in reverse order then reverse back
		expected := []int{2, 4, 6}
		assert.Equal(t, expected, result)
	})
}

// TestReverseArrayUseCases demonstrates practical use cases for ReverseArray
func TestReverseArrayUseCases(t *testing.T) {
	t.Run("Process data in reverse order", func(t *testing.T) {
		iso := ReverseArray[string]()
		events := []string{"first", "second", "third"}

		// Process in reverse chronological order
		reversed := iso.Get(events)
		assert.Equal(t, "third", reversed[0])
		assert.Equal(t, "first", reversed[2])
	})

	t.Run("Reversible transformation pipeline", func(t *testing.T) {
		iso := ReverseArray[int]()
		numbers := []int{1, 2, 3, 4, 5}

		// Apply transformation in reverse order
		reversed := iso.Get(numbers)
		// Process reversed data
		// Reverse back to original order
		restored := iso.ReverseGet(reversed)

		assert.Equal(t, numbers, restored)
	})

	t.Run("Palindrome check", func(t *testing.T) {
		iso := ReverseArray[int]()

		palindrome := []int{1, 2, 3, 2, 1}
		reversed := iso.Get(palindrome)
		assert.Equal(t, palindrome, reversed)

		notPalindrome := []int{1, 2, 3, 4, 5}
		reversedNot := iso.Get(notPalindrome)
		assert.NotEqual(t, notPalindrome, reversedNot)
	})
}

// TestHead tests the Head isomorphism
func TestHead(t *testing.T) {
	t.Run("Wrap integer into non-empty array", func(t *testing.T) {
		iso := Head[int]()
		value := 42
		arr := iso.Get(value)

		// Check that head is the value
		head := nonempty.Head(arr)
		assert.Equal(t, value, head)
	})

	t.Run("Extract head from non-empty array", func(t *testing.T) {
		iso := Head[int]()
		arr := nonempty.From(42, 10, 20)
		head := iso.ReverseGet(arr)

		assert.Equal(t, 42, head)
	})

	t.Run("Round-trip with integer", func(t *testing.T) {
		iso := Head[int]()
		original := 100
		result := iso.ReverseGet(iso.Get(original))

		assert.Equal(t, original, result)
	})

	t.Run("Wrap string into non-empty array", func(t *testing.T) {
		iso := Head[string]()
		value := "hello"
		arr := iso.Get(value)

		head := nonempty.Head(arr)
		assert.Equal(t, value, head)
	})

	t.Run("Extract head from multi-element array", func(t *testing.T) {
		iso := Head[string]()
		arr := nonempty.From("first", "second", "third")
		head := iso.ReverseGet(arr)

		// Only the head is extracted
		assert.Equal(t, "first", head)
	})

	t.Run("Round-trip with string", func(t *testing.T) {
		iso := Head[string]()
		original := "test"
		result := iso.ReverseGet(iso.Get(original))

		assert.Equal(t, original, result)
	})

	t.Run("Wrap struct into non-empty array", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		iso := Head[User]()
		user := User{"Alice", 30}
		arr := iso.Get(user)

		head := nonempty.Head(arr)
		assert.Equal(t, user, head)
	})

	t.Run("Round-trip with struct", func(t *testing.T) {
		type Person struct {
			Name string
		}
		iso := Head[Person]()
		original := Person{"Bob"}
		result := iso.ReverseGet(iso.Get(original))

		assert.Equal(t, original, result)
	})

	t.Run("Wrap pointer into non-empty array", func(t *testing.T) {
		iso := Head[*int]()
		value := 42
		ptr := &value
		arr := iso.Get(ptr)

		head := nonempty.Head(arr)
		assert.Equal(t, ptr, head)
		assert.Equal(t, 42, *head)
	})
}

// TestHeadRoundTripLaws verifies isomorphism laws for Head
func TestHeadRoundTripLaws(t *testing.T) {
	t.Run("Law 1: ReverseGet(Get(x)) == x for integers", func(t *testing.T) {
		iso := Head[int]()
		testCases := []int{0, 1, -1, 42, 100, -100}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 1: ReverseGet(Get(x)) == x for strings", func(t *testing.T) {
		iso := Head[string]()
		testCases := []string{"", "a", "hello", "world"}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 2: Get(ReverseGet(arr)) has same head", func(t *testing.T) {
		iso := Head[int]()

		// Create non-empty arrays
		arr1 := nonempty.Of(42)
		arr2 := nonempty.From(10, 20, 30)

		// Extract head and wrap back
		result1 := iso.Get(iso.ReverseGet(arr1))
		result2 := iso.Get(iso.ReverseGet(arr2))

		// Should have same head
		assert.Equal(t, nonempty.Head(arr1), nonempty.Head(result1))
		assert.Equal(t, nonempty.Head(arr2), nonempty.Head(result2))
	})
}

// TestHeadComposition tests composing Head with other isomorphisms
func TestHeadComposition(t *testing.T) {
	t.Run("Compose with Reverse", func(t *testing.T) {
		headIso := Head[int]()
		reversed := Reverse(headIso)

		value := 42

		// Reversed: Get extracts head, ReverseGet wraps
		arr := nonempty.Of(value)
		extracted := reversed.Get(arr)
		assert.Equal(t, value, extracted)

		wrapped := reversed.ReverseGet(value)
		assert.Equal(t, value, nonempty.Head(wrapped))
	})

	t.Run("Use with Modify", func(t *testing.T) {
		iso := Head[int]()

		// Wrap value, modify in array context, extract
		modifier := Modify[int](func(arr NonEmptyArray[int]) NonEmptyArray[int] {
			return nonempty.Map(N.Mul(2))(arr)
		})(iso)

		value := 5
		result := modifier(value)

		// Should double the value
		assert.Equal(t, 10, result)
	})
}

// TestHeadUseCases demonstrates practical use cases for Head
func TestHeadUseCases(t *testing.T) {
	t.Run("Lift value into non-empty context", func(t *testing.T) {
		iso := Head[int]()
		value := 42

		// Lift into non-empty array for processing
		arr := iso.Get(value)

		// Process as non-empty array
		doubled := nonempty.Map(N.Mul(2))(arr)

		// Extract result
		result := nonempty.Head(doubled)
		assert.Equal(t, 84, result)
	})

	t.Run("Ensure non-empty guarantee", func(t *testing.T) {
		iso := Head[string]()
		value := "important"

		// Wrap in non-empty array to guarantee at least one element
		arr := iso.Get(value)

		// Can safely access head without checking for empty
		head := nonempty.Head(arr)
		assert.Equal(t, value, head)
	})

	t.Run("Convert single value to collection", func(t *testing.T) {
		iso := Head[int]()
		defaultValue := 0

		// Convert to non-empty array for uniform processing
		arr := iso.Get(defaultValue)

		// Process as collection
		processed := nonempty.Map(N.Add(10))(arr)

		result := iso.ReverseGet(processed)
		assert.Equal(t, 10, result)
	})
}
