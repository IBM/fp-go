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
	"net/url"
	"regexp"
	"strconv"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validate"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

// validateFromParser creates a validation function from a parser that may fail.
// It wraps a parser function that returns (A, error) into a Validate[I, A] function
// that integrates with the validation framework.
//
// The returned validation function:
//   - Calls the parser with the input value
//   - On success: returns a successful validation containing the parsed value
//   - On failure: returns a validation failure with the error message and cause
//
// Type Parameters:
//   - A: The target type to parse into
//   - I: The input type to parse from
//
// Parameters:
//   - parser: A function that attempts to parse input I into type A, returning an error on failure
//
// Returns:
//   - A Validate[I, A] function that can be used in codec construction
//
// Example:
//
//	// Create a validator for parsing integers from strings
//	intValidator := validateFromParser(strconv.Atoi)
//	// Use in a codec
//	intCodec := MakeType("Int", Is[int](), intValidator, strconv.Itoa)
func validateFromParser[A, I any](parser func(I) (A, error)) Validate[I, A] {
	return func(i I) Decode[Context, A] {
		// Attempt to parse the input value
		a, err := parser(i)
		if err != nil {
			// On error, create a validation failure with the error details
			return validation.FailureWithError[A](i, err.Error())(err)
		}
		// On success, wrap the parsed value in a successful validation
		return reader.Of[Context](validation.Success(a))
	}
}

// URL creates a bidirectional codec for URL parsing and formatting.
// This codec can parse strings into *url.URL and encode *url.URL back to strings.
//
// The codec:
//   - Decodes: Parses a string using url.Parse, validating URL syntax
//   - Encodes: Converts a *url.URL to its string representation using String()
//   - Validates: Ensures the input string is a valid URL format
//
// Returns:
//   - A Type[*url.URL, string, string] codec that handles URL transformations
//
// Example:
//
//	urlCodec := URL()
//
//	// Decode a string to URL
//	validation := urlCodec.Decode("https://example.com/path?query=value")
//	// validation is Right(*url.URL{...})
//
//	// Encode a URL to string
//	u, _ := url.Parse("https://example.com")
//	str := urlCodec.Encode(u)
//	// str is "https://example.com"
//
//	// Invalid URL fails validation
//	validation := urlCodec.Decode("not a valid url")
//	// validation is Left(ValidationError{...})
func URL() Type[*url.URL, string, string] {
	return MakeType(
		"URL",
		Is[*url.URL](),
		validateFromParser(url.Parse),
		(*url.URL).String,
	)
}

// Date creates a bidirectional codec for date/time parsing and formatting with a specific layout.
// This codec uses Go's time.Parse and time.Format with the provided layout string.
//
// The codec:
//   - Decodes: Parses a string into time.Time using the specified layout
//   - Encodes: Formats a time.Time back to a string using the same layout
//   - Validates: Ensures the input string matches the expected date/time format
//
// Parameters:
//   - layout: The time layout string (e.g., "2006-01-02", time.RFC3339)
//     See time package documentation for layout format details
//
// Returns:
//   - A Type[time.Time, string, string] codec that handles date/time transformations
//
// Example:
//
//	// Create a codec for ISO 8601 dates
//	dateCodec := Date("2006-01-02")
//
//	// Decode a string to time.Time
//	validation := dateCodec.Decode("2024-03-15")
//	// validation is Right(time.Time{...})
//
//	// Encode a time.Time to string
//	t := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
//	str := dateCodec.Encode(t)
//	// str is "2024-03-15"
//
//	// Create a codec for RFC3339 timestamps
//	timestampCodec := Date(time.RFC3339)
//	validation := timestampCodec.Decode("2024-03-15T10:30:00Z")
//
//	// Invalid format fails validation
//	validation := dateCodec.Decode("15-03-2024")
//	// validation is Left(ValidationError{...})
func Date(layout string) Type[time.Time, string, string] {
	return MakeType(
		"Date",
		Is[time.Time](),
		validateFromParser(func(s string) (time.Time, error) { return time.Parse(layout, s) }),
		F.Bind2nd(time.Time.Format, layout),
	)
}

// Regex creates a bidirectional codec for regex pattern matching with capture groups.
// This codec can match strings against a regular expression pattern and extract capture groups,
// then reconstruct the original string from the match data.
//
// The codec uses prism.Match which contains:
//   - Before: Text before the match
//   - Groups: Capture groups (index 0 is the full match, 1+ are numbered capture groups)
//   - After: Text after the match
//
// The codec:
//   - Decodes: Attempts to match the regex against the input string
//   - Encodes: Reconstructs the original string from a Match structure
//   - Validates: Ensures the string matches the regex pattern
//
// Parameters:
//   - re: A compiled regular expression pattern
//
// Returns:
//   - A Type[prism.Match, string, string] codec that handles regex matching
//
// Example:
//
//	// Create a codec for matching numbers in text
//	numberRegex := regexp.MustCompile(`\d+`)
//	numberCodec := Regex(numberRegex)
//
//	// Decode a string with a number
//	validation := numberCodec.Decode("Price: 42 dollars")
//	// validation is Right(Match{Before: "Price: ", Groups: []string{"42"}, After: " dollars"})
//
//	// Encode a Match back to string
//	match := prism.Match{Before: "Price: ", Groups: []string{"42"}, After: " dollars"}
//	str := numberCodec.Encode(match)
//	// str is "Price: 42 dollars"
//
//	// Non-matching string fails validation
//	validation := numberCodec.Decode("no numbers here")
//	// validation is Left(ValidationError{...})
func Regex(re *regexp.Regexp) Type[prism.Match, string, string] {
	return FromRefinement(prism.RegexMatcher(re))
}

// RegexNamed creates a bidirectional codec for regex pattern matching with named capture groups.
// This codec can match strings against a regular expression with named groups and extract them
// by name, then reconstruct the original string from the match data.
//
// The codec uses prism.NamedMatch which contains:
//   - Before: Text before the match
//   - Groups: Map of named capture groups (name -> matched text)
//   - Full: The complete matched text
//   - After: Text after the match
//
// The codec:
//   - Decodes: Attempts to match the regex against the input string
//   - Encodes: Reconstructs the original string from a NamedMatch structure
//   - Validates: Ensures the string matches the regex pattern with named groups
//
// Parameters:
//   - re: A compiled regular expression with named capture groups (e.g., `(?P<name>pattern)`)
//
// Returns:
//   - A Type[prism.NamedMatch, string, string] codec that handles named regex matching
//
// Example:
//
//	// Create a codec for matching email addresses with named groups
//	emailRegex := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
//	emailCodec := RegexNamed(emailRegex)
//
//	// Decode an email string
//	validation := emailCodec.Decode("john@example.com")
//	// validation is Right(NamedMatch{
//	//     Before: "",
//	//     Groups: map[string]string{"user": "john", "domain": "example.com"},
//	//     Full: "john@example.com",
//	//     After: ""
//	// })
//
//	// Encode a NamedMatch back to string
//	match := prism.NamedMatch{
//	    Before: "",
//	    Groups: map[string]string{"user": "john", "domain": "example.com"},
//	    Full: "john@example.com",
//	    After: "",
//	}
//	str := emailCodec.Encode(match)
//	// str is "john@example.com"
//
//	// Non-matching string fails validation
//	validation := emailCodec.Decode("not-an-email")
//	// validation is Left(ValidationError{...})
func RegexNamed(re *regexp.Regexp) Type[prism.NamedMatch, string, string] {
	return FromRefinement(prism.RegexNamedMatcher(re))
}

// IntFromString creates a bidirectional codec for parsing integers from strings.
// This codec converts string representations of integers to int values and vice versa.
//
// The codec:
//   - Decodes: Parses a string to an int using strconv.Atoi
//   - Encodes: Converts an int to its string representation using strconv.Itoa
//   - Validates: Ensures the string contains a valid integer (base 10)
//
// The codec accepts integers in base 10 format, with optional leading sign (+/-).
// It does not accept hexadecimal, octal, or other number formats.
//
// Returns:
//   - A Type[int, string, string] codec that handles int/string conversions
//
// Example:
//
//	intCodec := IntFromString()
//
//	// Decode a valid integer string
//	validation := intCodec.Decode("42")
//	// validation is Right(42)
//
//	// Decode negative integer
//	validation := intCodec.Decode("-123")
//	// validation is Right(-123)
//
//	// Encode an integer to string
//	str := intCodec.Encode(42)
//	// str is "42"
//
//	// Invalid integer string fails validation
//	validation := intCodec.Decode("not a number")
//	// validation is Left(ValidationError{...})
//
//	// Floating point fails validation
//	validation := intCodec.Decode("3.14")
//	// validation is Left(ValidationError{...})
func IntFromString() Type[int, string, string] {
	return MakeType(
		"IntFromString",
		Is[int](),
		validateFromParser(strconv.Atoi),
		strconv.Itoa,
	)
}

// Int64FromString creates a bidirectional codec for parsing 64-bit integers from strings.
// This codec converts string representations of integers to int64 values and vice versa.
//
// The codec:
//   - Decodes: Parses a string to an int64 using strconv.ParseInt with base 10
//   - Encodes: Converts an int64 to its string representation
//   - Validates: Ensures the string contains a valid 64-bit integer (base 10)
//
// The codec accepts integers in base 10 format, with optional leading sign (+/-).
// It supports the full range of int64 values (-9223372036854775808 to 9223372036854775807).
//
// Returns:
//   - A Type[int64, string, string] codec that handles int64/string conversions
//
// Example:
//
//	int64Codec := Int64FromString()
//
//	// Decode a valid integer string
//	validation := int64Codec.Decode("9223372036854775807")
//	// validation is Right(9223372036854775807)
//
//	// Decode negative integer
//	validation := int64Codec.Decode("-9223372036854775808")
//	// validation is Right(-9223372036854775808)
//
//	// Encode an int64 to string
//	str := int64Codec.Encode(42)
//	// str is "42"
//
//	// Invalid integer string fails validation
//	validation := int64Codec.Decode("not a number")
//	// validation is Left(ValidationError{...})
//
//	// Out of range value fails validation
//	validation := int64Codec.Decode("9223372036854775808")
//	// validation is Left(ValidationError{...})
func Int64FromString() Type[int64, string, string] {
	return MakeType(
		"Int64FromString",
		Is[int64](),
		validateFromParser(func(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }),
		prism.ParseInt64().ReverseGet,
	)
}

// BoolFromString creates a bidirectional codec for parsing boolean values from strings.
// This codec converts string representations of booleans to bool values and vice versa.
//
// The codec:
//   - Decodes: Parses a string to a bool using strconv.ParseBool
//   - Encodes: Converts a bool to its string representation using strconv.FormatBool
//   - Validates: Ensures the string contains a valid boolean value
//
// The codec accepts the following string values (case-insensitive):
//   - true: "1", "t", "T", "true", "TRUE", "True"
//   - false: "0", "f", "F", "false", "FALSE", "False"
//
// Returns:
//   - A Type[bool, string, string] codec that handles bool/string conversions
//
// Example:
//
//	boolCodec := BoolFromString()
//
//	// Decode valid boolean strings
//	validation := boolCodec.Decode("true")
//	// validation is Right(true)
//
//	validation := boolCodec.Decode("1")
//	// validation is Right(true)
//
//	validation := boolCodec.Decode("false")
//	// validation is Right(false)
//
//	validation := boolCodec.Decode("0")
//	// validation is Right(false)
//
//	// Encode a boolean to string
//	str := boolCodec.Encode(true)
//	// str is "true"
//
//	str := boolCodec.Encode(false)
//	// str is "false"
//
//	// Invalid boolean string fails validation
//	validation := boolCodec.Decode("yes")
//	// validation is Left(ValidationError{...})
//
//	// Case variations are accepted
//	validation := boolCodec.Decode("TRUE")
//	// validation is Right(true)
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

// MarshalText creates a bidirectional codec for types that implement encoding.TextMarshaler
// and encoding.TextUnmarshaler. This codec handles binary text serialization formats.
//
// The codec:
//   - Decodes: Calls dec.UnmarshalText(b) to deserialize []byte into the target type T
//   - Encodes: Calls enc.MarshalText() to serialize the value to []byte
//   - Validates: Returns a failure if UnmarshalText returns an error
//
// Note: The enc and dec parameters are external marshaler/unmarshaler instances. The
// decoded value is the zero value of T after UnmarshalText has been called on dec
// (the caller is responsible for ensuring dec holds the decoded state).
//
// Type Parameters:
//   - T: The Go type to encode/decode
//
// Parameters:
//   - enc: An encoding.TextMarshaler used for encoding values to []byte
//   - dec: An encoding.TextUnmarshaler used for decoding []byte to the target type
//
// Returns:
//   - A Type[T, []byte, []byte] codec that handles text marshaling/unmarshaling
//
// Example:
//
//	type MyType struct{ Value string }
//
//	var instance MyType
//	codec := MarshalText[MyType](instance, &instance)
//
//	// Decode bytes to MyType
//	result := codec.Decode([]byte(`some text`))
//
//	// Encode MyType to bytes
//	encoded := codec.Encode(instance)
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

// MarshalJSON creates a bidirectional codec for types that implement encoding/json's
// json.Marshaler and json.Unmarshaler interfaces. This codec handles JSON serialization.
//
// The codec:
//   - Decodes: Calls dec.UnmarshalJSON(b) to deserialize []byte JSON into the target type T
//   - Encodes: Calls enc.MarshalJSON() to serialize the value to []byte JSON
//   - Validates: Returns a failure if UnmarshalJSON returns an error
//
// Note: The enc and dec parameters are external marshaler/unmarshaler instances. The
// decoded value is the zero value of T after UnmarshalJSON has been called on dec
// (the caller is responsible for ensuring dec holds the decoded state).
//
// Type Parameters:
//   - T: The Go type to encode/decode
//
// Parameters:
//   - enc: A json.Marshaler used for encoding values to JSON []byte
//   - dec: A json.Unmarshaler used for decoding JSON []byte to the target type
//
// Returns:
//   - A Type[T, []byte, []byte] codec that handles JSON marshaling/unmarshaling
//
// Example:
//
//	type MyData struct {
//	    Name  string `json:"name"`
//	    Value int    `json:"value"`
//	}
//
//	var instance MyData
//	codec := MarshalJSON[MyData](&instance, &instance)
//
//	// Decode JSON bytes to MyData
//	result := codec.Decode([]byte(`{"name":"test","value":42}`))
//
//	// Encode MyData to JSON bytes
//	encoded := codec.Encode(instance)
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

// FromNonZero creates a bidirectional codec for non-zero values of comparable types.
// This codec validates that values are not equal to their zero value (e.g., 0 for int,
// "" for string, false for bool, nil for pointers).
//
// The codec uses a refinement (prism) that:
//   - Decodes: Validates that the input is not the zero value of type T
//   - Encodes: Returns the value unchanged (identity function)
//   - Validates: Ensures the value is non-zero/non-default
//
// This is useful for enforcing that required fields have meaningful values rather than
// their default zero values, which often represent "not set" or "missing" states.
//
// Type Parameters:
//   - T: A comparable type (must support == and != operators)
//
// Returns:
//   - A Type[T, T, T] codec that validates non-zero values
//
// Example:
//
//	// Create a codec for non-zero integers
//	nonZeroInt := FromNonZero[int]()
//
//	// Decode non-zero value succeeds
//	result := nonZeroInt.Decode(42)
//	// result is Right(42)
//
//	// Decode zero value fails
//	result := nonZeroInt.Decode(0)
//	// result is Left(ValidationError{...})
//
//	// Encode is identity
//	encoded := nonZeroInt.Encode(42)
//	// encoded is 42
//
//	// Works with strings
//	nonEmptyStr := FromNonZero[string]()
//	result := nonEmptyStr.Decode("hello")  // Right("hello")
//	result = nonEmptyStr.Decode("")        // Left(ValidationError{...})
//
//	// Works with pointers
//	nonNilPtr := FromNonZero[*int]()
//	value := 42
//	result := nonNilPtr.Decode(&value)  // Right(&value)
//	result = nonNilPtr.Decode(nil)      // Left(ValidationError{...})
//
// Common use cases:
//   - Validating required numeric fields are not zero
//   - Ensuring string fields are not empty
//   - Checking pointers are not nil
//   - Validating boolean flags are explicitly set to true
//   - Composing with other codecs for multi-stage validation
//
// See Also:
//   - NonEmptyString: Specialized version for strings with clearer intent
//   - FromRefinement: General function for creating codecs from prisms
func FromNonZero[T comparable]() Type[T, T, T] {
	return FromRefinement(prism.FromNonZero[T]())
}

// NonEmptyString creates a bidirectional codec for non-empty strings.
// This codec validates that string values are not empty, providing a type-safe
// way to work with strings that must contain at least one character.
//
// This is a specialized version of FromNonZero[string]() that makes the intent
// clearer when working specifically with strings that must not be empty.
//
// The codec:
//   - Decodes: Validates that the input string is not empty ("")
//   - Encodes: Returns the string unchanged (identity function)
//   - Validates: Ensures the string has length > 0
//
// Note: This codec only checks for empty strings, not whitespace-only strings.
// A string containing only spaces, tabs, or newlines will pass validation.
//
// Returns:
//   - A Type[string, string, string] codec that validates non-empty strings
//
// Example:
//
//	nonEmpty := NonEmptyString()
//
//	// Decode non-empty string succeeds
//	result := nonEmpty.Decode("hello")
//	// result is Right("hello")
//
//	// Decode empty string fails
//	result := nonEmpty.Decode("")
//	// result is Left(ValidationError{...})
//
//	// Whitespace-only strings pass validation
//	result := nonEmpty.Decode("   ")
//	// result is Right("   ")
//
//	// Encode is identity
//	encoded := nonEmpty.Encode("world")
//	// encoded is "world"
//
//	// Compose with other codecs for validation pipelines
//	intFromNonEmptyString := Pipe(IntFromString())(nonEmpty)
//	result := intFromNonEmptyString.Decode("42")   // Right(42)
//	result = intFromNonEmptyString.Decode("")      // Left(ValidationError{...})
//	result = intFromNonEmptyString.Decode("abc")   // Left(ValidationError{...})
//
// Common use cases:
//   - Validating required string fields (usernames, names, IDs)
//   - Ensuring configuration values are provided
//   - Validating user input before processing
//   - Composing with parsing codecs to validate before parsing
//   - Building validation pipelines for string data
//
// See Also:
//   - FromNonZero: General version for any comparable type
//   - String: Basic string codec without validation
//   - IntFromString: Codec for parsing integers from strings
func NonEmptyString() Type[string, string, string] {
	return F.Pipe1(
		FromRefinement(prism.NonEmptyString()),
		WithName[string, string, string]("NonEmptyString"),
	)
}

// WithName creates an endomorphism that renames a codec without changing its behavior.
// This function returns a higher-order function that takes a codec and returns a new codec
// with the specified name, while preserving all validation, encoding, and type-checking logic.
//
// This is useful for:
//   - Providing more descriptive names for composed codecs
//   - Creating domain-specific codec names for better error messages
//   - Documenting the purpose of complex codec pipelines
//   - Improving debugging and logging output
//
// The renamed codec maintains the same:
//   - Type checking behavior (Is function)
//   - Validation logic (Validate function)
//   - Encoding behavior (Encode function)
//
// Only the name returned by the Name() method changes.
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
// Example:
//
//	// Create a codec with a generic name
//	positiveInt := Pipe[int, int, string, int](
//	    FromRefinement(prism.FromPredicate(func(n int) bool { return n > 0 })),
//	)(IntFromString())
//	// positiveInt.Name() returns something like "Pipe(FromRefinement(...), IntFromString)"
//
//	// Rename it for clarity
//	namedCodec := WithName[int, string, string]("PositiveIntFromString")(positiveInt)
//	// namedCodec.Name() returns "PositiveIntFromString"
//
//	// Use in a pipeline with F.Pipe
//	userAgeCodec := F.Pipe1(
//	    IntFromString(),
//	    WithName[int, string, string]("UserAge"),
//	)
//
//	// Validation errors will show the custom name
//	result := userAgeCodec.Decode("invalid")
//	// Error context will reference "UserAge" instead of "IntFromString"
//
// Common use cases:
//   - Naming composed codecs for better error messages
//   - Creating domain-specific codec names (e.g., "EmailAddress", "PhoneNumber")
//   - Documenting complex validation pipelines
//   - Improving debugging output in logs
//   - Making codec composition more readable
//
// Note: This function creates a new codec instance with the same behavior but a different
// name. The original codec is not modified.
//
// See Also:
//   - MakeType: For creating codecs with custom names from scratch
//   - Pipe: For composing codecs (which generates automatic names)
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
