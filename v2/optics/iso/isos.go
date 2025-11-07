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
	"strings"
	"time"

	B "github.com/IBM/fp-go/v2/bytes"
	F "github.com/IBM/fp-go/v2/function"
	S "github.com/IBM/fp-go/v2/string"
)

// UTF8String creates an isomorphism between byte slices and UTF-8 strings.
// This isomorphism provides bidirectional conversion between []byte and string,
// treating the byte slice as UTF-8 encoded text.
//
// Returns:
//   - An Iso[[]byte, string] where:
//   - Get: Converts []byte to string using UTF-8 encoding
//   - ReverseGet: Converts string to []byte using UTF-8 encoding
//
// Behavior:
//   - Get direction: Interprets the byte slice as UTF-8 and returns the corresponding string
//   - ReverseGet direction: Encodes the string as UTF-8 bytes
//
// Example:
//
//	iso := UTF8String()
//
//	// Convert bytes to string
//	str := iso.Get([]byte("hello"))  // "hello"
//
//	// Convert string to bytes
//	bytes := iso.ReverseGet("world")  // []byte("world")
//
//	// Round-trip conversion
//	original := []byte("test")
//	result := iso.ReverseGet(iso.Get(original))  // []byte("test")
//
// Use cases:
//   - Converting between string and byte representations
//   - Working with APIs that use different text representations
//   - File I/O operations where you need to switch between strings and bytes
//   - Network protocols that work with byte streams
//
// Note: This isomorphism assumes valid UTF-8 encoding. Invalid UTF-8 sequences
// in the byte slice will be handled according to Go's string conversion rules
// (typically replaced with the Unicode replacement character U+FFFD).
func UTF8String() Iso[[]byte, string] {
	return MakeIso(B.ToString, S.ToBytes)
}

// lines creates an isomorphism between a slice of strings and a single string
// with lines separated by the specified separator.
// This is an internal helper function used by Lines.
//
// Parameters:
//   - sep: The separator string to use for joining/splitting lines
//
// Returns:
//   - An Iso[[]string, string] that joins/splits strings using the separator
//
// Behavior:
//   - Get direction: Joins the string slice into a single string with separators
//   - ReverseGet direction: Splits the string by the separator into a slice
func lines(sep string) Iso[[]string, string] {
	return MakeIso(S.Join(sep), F.Bind2nd(strings.Split, sep))
}

// Lines creates an isomorphism between a slice of strings and a single string
// with newline-separated lines.
// This is useful for working with multi-line text where you need to convert
// between a single string and individual lines.
//
// Returns:
//   - An Iso[[]string, string] where:
//   - Get: Joins string slice with newline characters ("\n")
//   - ReverseGet: Splits string by newline characters into a slice
//
// Behavior:
//   - Get direction: Joins each string in the slice with "\n" separator
//   - ReverseGet direction: Splits the string at each "\n" into a slice
//
// Example:
//
//	iso := Lines()
//
//	// Convert lines to single string
//	lines := []string{"line1", "line2", "line3"}
//	text := iso.Get(lines)  // "line1\nline2\nline3"
//
//	// Convert string to lines
//	text := "hello\nworld"
//	lines := iso.ReverseGet(text)  // []string{"hello", "world"}
//
//	// Round-trip conversion
//	original := []string{"a", "b", "c"}
//	result := iso.ReverseGet(iso.Get(original))  // []string{"a", "b", "c"}
//
// Use cases:
//   - Processing multi-line text files
//   - Converting between text editor representations (array of lines vs single string)
//   - Working with configuration files that have line-based structure
//   - Parsing or generating multi-line output
//
// Note: Empty strings in the slice will result in consecutive newlines in the output.
// Splitting a string with trailing newlines will include an empty string at the end.
//
// Example with edge cases:
//
//	iso := Lines()
//	lines := []string{"a", "", "b"}
//	text := iso.Get(lines)  // "a\n\nb"
//	result := iso.ReverseGet(text)  // []string{"a", "", "b"}
//
//	text := "a\nb\n"
//	lines := iso.ReverseGet(text)  // []string{"a", "b", ""}
func Lines() Iso[[]string, string] {
	return lines("\n")
}

// UnixMilli creates an isomorphism between Unix millisecond timestamps and time.Time values.
// This isomorphism provides bidirectional conversion between int64 milliseconds since
// the Unix epoch (January 1, 1970 UTC) and Go's time.Time type.
//
// Returns:
//   - An Iso[int64, time.Time] where:
//   - Get: Converts Unix milliseconds (int64) to time.Time
//   - ReverseGet: Converts time.Time to Unix milliseconds (int64)
//
// Behavior:
//   - Get direction: Creates a time.Time from milliseconds since Unix epoch
//   - ReverseGet direction: Extracts milliseconds since Unix epoch from time.Time
//
// Example:
//
//	iso := UnixMilli()
//
//	// Convert milliseconds to time.Time
//	millis := int64(1609459200000)  // 2021-01-01 00:00:00 UTC
//	t := iso.Get(millis)
//
//	// Convert time.Time to milliseconds
//	now := time.Now()
//	millis := iso.ReverseGet(now)
//
//	// Round-trip conversion
//	original := int64(1234567890000)
//	result := iso.ReverseGet(iso.Get(original))  // 1234567890000
//
// Use cases:
//   - Working with APIs that use Unix millisecond timestamps (e.g., JavaScript Date.now())
//   - Database storage where timestamps are stored as integers
//   - JSON serialization/deserialization of timestamps
//   - Converting between different time representations in distributed systems
//
// Precision notes:
//   - Millisecond precision is maintained in both directions
//   - Sub-millisecond precision in time.Time is lost when converting to int64
//   - The conversion is timezone-aware (time.Time includes location information)
//
// Example with precision:
//
//	iso := UnixMilli()
//	t := time.Date(2021, 1, 1, 12, 30, 45, 123456789, time.UTC)
//	millis := iso.ReverseGet(t)  // Nanoseconds are truncated to milliseconds
//	restored := iso.Get(millis)   // Nanoseconds will be 123000000
//
// Note: This isomorphism uses UTC for the time.Time values. If you need to preserve
// timezone information, consider storing it separately or using a different representation.
func UnixMilli() Iso[int64, time.Time] {
	return MakeIso(time.UnixMilli, time.Time.UnixMilli)
}
