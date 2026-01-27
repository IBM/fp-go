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
	"net/url"
	"regexp"
	"strconv"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/reader"
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
