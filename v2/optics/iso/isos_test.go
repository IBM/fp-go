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
	"testing"
	"time"

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
		lines := []string{"line1", "line2", "line3"}
		result := iso.Get(lines)
		assert.Equal(t, "line1\nline2\nline3", result)
	})

	t.Run("Get handles single line", func(t *testing.T) {
		lines := []string{"single line"}
		result := iso.Get(lines)
		assert.Equal(t, "single line", result)
	})

	t.Run("Get handles empty slice", func(t *testing.T) {
		lines := []string{}
		result := iso.Get(lines)
		assert.Equal(t, "", result)
	})

	t.Run("Get handles empty strings in slice", func(t *testing.T) {
		lines := []string{"a", "", "b"}
		result := iso.Get(lines)
		assert.Equal(t, "a\n\nb", result)
	})

	t.Run("Get handles slice with only empty strings", func(t *testing.T) {
		lines := []string{"", "", ""}
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
		original := []string{}
		text := iso.Get(original)      // ""
		result := iso.ReverseGet(text) // [""]
		assert.Equal(t, []string{""}, result)
	})

	t.Run("Law 2: Get(ReverseGet(str)) == str", func(t *testing.T) {
		testCases := []string{
			"line1\nline2",
			"single",
			"",
			"a\n\nb",
			"\n\n",
		}

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
		testCases := []int64{
			0,
			1609459200000,
			-86400000,
			1234567890000,
			time.Now().UnixMilli(),
		}

		for _, original := range testCases {
			result := iso.ReverseGet(iso.Get(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Law 2: Get(ReverseGet(time)) == time (with millisecond precision)", func(t *testing.T) {
		testCases := []time.Time{
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Unix(0, 0).UTC(),
			time.Date(1969, 12, 31, 0, 0, 0, 0, time.UTC),
			time.Now().Truncate(time.Millisecond),
		}

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
		originalLines := []string{"a", "b", "c"}
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
