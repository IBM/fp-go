// Copyright (c) 2024 IBM Corp.
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

// Package codec provides pre-built codec implementations for common types.
// This package includes codecs for URL parsing, date/time formatting, and other
// standard data transformations that require bidirectional encoding/decoding.
//
// The codecs in this package follow functional programming principles and integrate
// with the validation framework to provide type-safe, composable transformations.
package codec

import (
	"encoding"
	"encoding/json"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validate"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/readerresult"
	"github.com/IBM/fp-go/v2/result"
)

// validateFromParser creates a validation function from a parser that may fail.
// It wraps a parser function that returns (A, error) into a Validate[I, A] function
// that integrates with the validation framework.
//
// On success the parsed value is returned as a successful validation.
// On failure a validation error carrying the underlying error is returned.
//
// Type Parameters:
//   - A: The target type to parse into
//   - I: The input type to parse from
//
// Parameters:
//   - parser: A function that attempts to parse input I into type A,
//     returning an error on failure
//
// Returns:
//   - A Validate[I, A] function that can be used in codec construction
func validateFromParser[A, I any](parser func(I) (A, error)) Validate[I, A] {
	return F.Pipe2(
		parser,
		readerresult.FromIdiomatic,
		validate.FromReaderResult,
	)
}

// URL creates a bidirectional codec for URL parsing and formatting.
//
// The codec decodes by calling url.Parse on the input string and validates URL
// syntax.  Encoding converts a *url.URL back to its string representation via
// (*url.URL).String.
//
// Returns:
//   - A Type[*url.URL, string, string] codec
func URL() Type[*url.URL, string, string] {
	return MakeType(
		"URL",
		Is[*url.URL](),
		validateFromParser(url.Parse),
		(*url.URL).String,
	)
}

// FileInfoWithPath extends fs.FileInfo with an absolute filesystem path.
// It combines all standard file metadata from fs.FileInfo with the resolved
// absolute path of the file, as returned by filepath.Abs.
//
// This interface is the decoded value type produced by the Stat codec.
//
// See Also:
//   - Stat: the codec that decodes a path string into a FileInfoWithPath
type FileInfoWithPath interface {
	fs.FileInfo
	// AbsPath returns the absolute path of the file as resolved by filepath.Abs.
	AbsPath() string
}

type fileInfoWithPath struct {
	fs.FileInfo
	absPath string
}

func (f *fileInfoWithPath) AbsPath() string {
	return f.absPath
}

func statFileInfoWithPath(path string) (FileInfoWithPath, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	return &fileInfoWithPath{FileInfo: stat, absPath: absPath}, nil
}

// Stat creates a bidirectional codec for file-system path resolution and stat
// retrieval.
//
// The codec resolves the input string to an absolute path via filepath.Abs,
// calls os.Stat on it, and wraps the resulting fs.FileInfo together with the
// resolved absolute path in a FileInfoWithPath value.
//
// Encoding converts a FileInfoWithPath back to its absolute path string via
// AbsPath.
//
// Decoding fails if filepath.Abs fails or if the path does not exist or is not
// accessible (i.e. os.Stat returns an error).
//
// Returns:
//   - A Type[FileInfoWithPath, string, string] codec
//
// See Also:
//   - FileInfoWithPath: the decoded value type
func Stat() Type[FileInfoWithPath, string, string] {
	return MakeType(
		"Stat",
		Is[FileInfoWithPath](),
		validateFromParser(statFileInfoWithPath),
		FileInfoWithPath.AbsPath,
	)
}

// Date creates a bidirectional codec for date/time parsing and formatting with
// a specific layout.
//
// The codec decodes by calling time.Parse(layout, s) and encodes by calling
// time.Time.Format(layout).  The codec name is always "Date" regardless of the
// layout.
//
// Parameters:
//   - layout: The time layout string (e.g., "2006-01-02", time.RFC3339).
//     See the time package documentation for layout format details.
//
// Returns:
//   - A Type[time.Time, string, string] codec
func Date(layout string) Type[time.Time, string, string] {
	return MakeType(
		"Date",
		Is[time.Time](),
		validateFromParser(func(s string) (time.Time, error) { return time.Parse(layout, s) }),
		F.Bind2nd(time.Time.Format, layout),
	)
}

// Regex creates a bidirectional codec for regex pattern matching with capture
// groups.
//
// The codec decodes by attempting to match the compiled regular expression
// against the input string.  On success it returns a prism.Match containing:
//   - Before: text before the match
//   - Groups: capture groups (index 0 is the full match, 1+ are numbered groups)
//   - After: text after the match
//
// Encoding reconstructs the original string from a prism.Match value.
// Decoding fails if the regex does not match.
//
// Parameters:
//   - re: A compiled regular expression
//
// Returns:
//   - A Type[prism.Match, string, string] codec
//
// See Also:
//   - RegexNamed: variant that exposes named capture groups as a map
func Regex(re *regexp.Regexp) Type[prism.Match, string, string] {
	return FromRefinement(prism.RegexMatcher(re))
}

// RegexNamed creates a bidirectional codec for regex pattern matching with
// named capture groups.
//
// The codec decodes by attempting to match the compiled regular expression
// against the input string.  On success it returns a prism.NamedMatch
// containing:
//   - Before: text before the match
//   - Groups: map of named capture groups (name → matched text)
//   - Full: the complete matched text
//   - After: text after the match
//
// Encoding reconstructs the original string from a prism.NamedMatch value.
// Decoding fails if the regex does not match.
//
// Parameters:
//   - re: A compiled regular expression with named capture groups
//     (e.g., `(?P<name>pattern)`)
//
// Returns:
//   - A Type[prism.NamedMatch, string, string] codec
//
// See Also:
//   - Regex: variant that exposes capture groups as an ordered slice
func RegexNamed(re *regexp.Regexp) Type[prism.NamedMatch, string, string] {
	return FromRefinement(prism.RegexNamedMatcher(re))
}

// IntFromString creates a bidirectional codec for parsing integers from
// strings.
//
// The codec decodes by calling strconv.Atoi and encodes by calling
// strconv.Itoa.  Only base-10 integers with an optional leading sign are
// accepted; hexadecimal, octal, and floating-point strings are rejected.
//
// Returns:
//   - A Type[int, string, string] codec
//
// See Also:
//   - Int64FromString: 64-bit variant with explicit range
func IntFromString() Type[int, string, string] {
	return MakeType(
		"IntFromString",
		Is[int](),
		validateFromParser(strconv.Atoi),
		strconv.Itoa,
	)
}

// Int64FromString creates a bidirectional codec for parsing 64-bit integers
// from strings.
//
// The codec decodes by calling strconv.ParseInt(s, 10, 64) and encodes by
// calling strconv.FormatInt.  Only base-10 integers are accepted.  Values
// outside the int64 range (-9223372036854775808 to 9223372036854775807) are
// rejected.
//
// Returns:
//   - A Type[int64, string, string] codec
//
// See Also:
//   - IntFromString: platform-width int variant
func Int64FromString() Type[int64, string, string] {
	return MakeType(
		"Int64FromString",
		Is[int64](),
		validateFromParser(func(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }),
		prism.ParseInt64().ReverseGet,
	)
}

// BoolFromString creates a bidirectional codec for parsing boolean values from
// strings.
//
// The codec decodes by calling strconv.ParseBool and encodes by calling
// strconv.FormatBool.  The accepted string values (per strconv.ParseBool) are:
//   - true:  "1", "t", "T", "TRUE", "true", "True"
//   - false: "0", "f", "F", "FALSE", "false", "False"
//
// Note that encoding always produces "true" or "false" regardless of which
// accepted input form was decoded (e.g. "1" decodes to true but re-encodes as
// "true").
//
// Returns:
//   - A Type[bool, string, string] codec
func BoolFromString() Type[bool, string, string] {
	return MakeType(
		"BoolFromString",
		Is[bool](),
		validateFromParser(strconv.ParseBool),
		strconv.FormatBool,
	)
}

func decodeJSON[T any](dec json.Unmarshaler) ReaderResult[[]byte, T] {
	return func(b []byte) Result[T] {
		var t T
		err := dec.UnmarshalJSON(b)
		return result.TryCatchError(t, err)
	}
}

func decodeText[T any](dec encoding.TextUnmarshaler) ReaderResult[[]byte, T] {
	return func(b []byte) Result[T] {
		var t T
		err := dec.UnmarshalText(b)
		return result.TryCatchError(t, err)
	}
}

// MarshalText creates a bidirectional codec for types that implement
// encoding.TextMarshaler and encoding.TextUnmarshaler.
//
// The codec decodes by calling dec.UnmarshalText(b) and encodes by calling
// enc.MarshalText().  Both enc and dec are caller-supplied instances; the
// caller is responsible for ensuring they share the same underlying state when
// a round-trip is required.
//
// Note: The codec name is "UnmarshalText" to reflect the primary decode
// operation.
//
// Type Parameters:
//   - T: The Go type to encode/decode
//
// Parameters:
//   - enc: An encoding.TextMarshaler used to encode values to []byte
//   - dec: An encoding.TextUnmarshaler used to decode []byte into type T
//
// Returns:
//   - A Type[T, []byte, []byte] codec
func MarshalText[T any](
	enc encoding.TextMarshaler,
	dec encoding.TextUnmarshaler,
) Type[T, []byte, []byte] {
	return MakeType(
		"UnmarshalText",
		Is[T](),
		F.Pipe2(
			dec,
			decodeText[T],
			validate.FromReaderResult,
		),
		func(t T) []byte {
			b, _ := enc.MarshalText()
			return b
		},
	)
}

// MarshalJSON creates a bidirectional codec for types that implement
// json.Marshaler and json.Unmarshaler.
//
// The codec decodes by calling dec.UnmarshalJSON(b) and encodes by calling
// enc.MarshalJSON().  Both enc and dec are caller-supplied instances; the
// caller is responsible for ensuring they share the same underlying state when
// a round-trip is required.
//
// Note: The codec name is "UnmarshalJSON" to reflect the primary decode
// operation.
//
// Type Parameters:
//   - T: The Go type to encode/decode
//
// Parameters:
//   - enc: A json.Marshaler used to encode values to JSON []byte
//   - dec: A json.Unmarshaler used to decode JSON []byte into type T
//
// Returns:
//   - A Type[T, []byte, []byte] codec
func MarshalJSON[T any](
	enc json.Marshaler,
	dec json.Unmarshaler,
) Type[T, []byte, []byte] {
	return MakeType(
		"UnmarshalJSON",
		Is[T](),
		F.Pipe2(
			dec,
			decodeJSON[T],
			validate.FromReaderResult,
		),
		func(t T) []byte {
			b, _ := enc.MarshalJSON()
			return b
		},
	)
}

// FromNonZero creates a bidirectional codec that validates the input is not
// equal to the zero value of T.
//
// The codec decodes by asserting the input is non-zero (using a prism built
// from prism.FromNonZero) and encodes by returning the value unchanged.
// Decoding fails with a validation error when the input equals the zero value
// of T (0 for numeric types, "" for string, false for bool, nil for pointers).
//
// Type Parameters:
//   - T: A comparable type
//
// Returns:
//   - A Type[T, T, T] codec that validates non-zero values
//
// See Also:
//   - NonEmptyString: specialised version for strings with a descriptive name
//   - FromRefinement: general function for creating codecs from prisms
func FromNonZero[T comparable]() Type[T, T, T] {
	return FromRefinement(prism.FromNonZero[T]())
}

// NonEmptyString creates a bidirectional codec for non-empty strings.
//
// The codec decodes by asserting the input string is not empty and encodes by
// returning the string unchanged.  A string containing only whitespace passes
// validation; only the empty string "" is rejected.
//
// This is a specialised version of FromNonZero[string]() that carries the name
// "NonEmptyString" for clearer error messages.
//
// Returns:
//   - A Type[string, string, string] codec
//
// See Also:
//   - FromNonZero: general non-zero codec for any comparable type
func NonEmptyString() Type[string, string, string] {
	return F.Pipe1(
		FromRefinement(prism.NonEmptyString()),
		WithName[string, string, string]("NonEmptyString"),
	)
}

// WithName returns an endomorphism that renames a codec without changing any
// of its behaviour.
//
// The returned function accepts a Type[A, O, I] and returns a new codec with
// the specified name while preserving all validation, encoding, and
// type-checking logic.  Only the value returned by Name() changes.
//
// Type Parameters:
//   - A: The target type (what we decode to and encode from)
//   - O: The output type (what we encode to)
//   - I: The input type (what we decode from)
//
// Parameters:
//   - name: The new name for the codec
//
// Returns:
//   - An Endomorphism[Type[A, O, I]] that renames the codec
//
// See Also:
//   - MakeType: for creating codecs with custom names from scratch
func WithName[A, O, I any](name string) Endomorphism[Type[A, O, I]] {
	return func(codec Type[A, O, I]) Type[A, O, I] {
		return MakeType(
			name,
			codec.Is,
			codec.Validate,
			codec.Encode,
		)
	}
}
