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

package prism

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	J "github.com/IBM/fp-go/v2/json"
	"github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
)

// FromEncoding creates a prism for base64 encoding/decoding operations.
// It provides a safe way to work with base64-encoded strings, handling
// encoding and decoding errors gracefully through the Option type.
//
// The prism's GetOption attempts to decode a base64 string into bytes.
// If decoding succeeds, it returns Some([]byte); if it fails (e.g., invalid
// base64 format), it returns None.
//
// The prism's ReverseGet always succeeds, encoding bytes into a base64 string.
//
// Parameters:
//   - enc: A base64.Encoding instance (e.g., base64.StdEncoding, base64.URLEncoding)
//
// Returns:
//   - A Prism[string, []byte] that safely handles base64 encoding/decoding
//
// Example:
//
//	// Create a prism for standard base64 encoding
//	b64Prism := FromEncoding(base64.StdEncoding)
//
//	// Decode valid base64 string
//	data := b64Prism.GetOption("SGVsbG8gV29ybGQ=")  // Some([]byte("Hello World"))
//
//	// Decode invalid base64 string
//	invalid := b64Prism.GetOption("not-valid-base64!!!")  // None[[]byte]()
//
//	// Encode bytes to base64
//	encoded := b64Prism.ReverseGet([]byte("Hello World"))  // "SGVsbG8gV29ybGQ="
//
//	// Use with Set to update encoded values
//	newData := []byte("New Data")
//	setter := Set[string, []byte](newData)
//	result := setter(b64Prism)("SGVsbG8gV29ybGQ=")  // Encodes newData to base64
//
// Common use cases:
//   - Safely decoding base64-encoded configuration values
//   - Working with base64-encoded API responses
//   - Validating and transforming base64 data in pipelines
//   - Using different encodings (Standard, URL-safe, RawStd, RawURL)
func FromEncoding(enc *base64.Encoding) Prism[string, []byte] {
	return MakePrismWithName(F.Flow2(
		either.Eitherize1(enc.DecodeString),
		either.Fold(F.Ignore1of1[error](option.None[[]byte]), option.Some),
	), enc.EncodeToString,
		"PrismFromEncoding",
	)
}

// ParseURL creates a prism for parsing and formatting URLs.
// It provides a safe way to work with URL strings, handling parsing
// errors gracefully through the Option type.
//
// The prism's GetOption attempts to parse a string into a *url.URL.
// If parsing succeeds, it returns Some(*url.URL); if it fails (e.g., invalid
// URL format), it returns None.
//
// The prism's ReverseGet always succeeds, converting a *url.URL back to its
// string representation.
//
// Returns:
//   - A Prism[string, *url.URL] that safely handles URL parsing/formatting
//
// Example:
//
//	// Create a URL parsing prism
//	urlPrism := ParseURL()
//
//	// Parse valid URL
//	parsed := urlPrism.GetOption("https://example.com/path?query=value")
//	// Some(*url.URL{Scheme: "https", Host: "example.com", ...})
//
//	// Parse invalid URL
//	invalid := urlPrism.GetOption("ht!tp://invalid url")  // None[*url.URL]()
//
//	// Convert URL back to string
//	u, _ := url.Parse("https://example.com")
//	str := urlPrism.ReverseGet(u)  // "https://example.com"
//
//	// Use with Set to update URLs
//	newURL, _ := url.Parse("https://newsite.com")
//	setter := Set[string, *url.URL](newURL)
//	result := setter(urlPrism)("https://oldsite.com")  // "https://newsite.com"
//
// Common use cases:
//   - Validating and parsing URL configuration values
//   - Working with API endpoints
//   - Transforming URL strings in data pipelines
//   - Extracting and modifying URL components safely
func ParseURL() Prism[string, *url.URL] {
	return MakePrismWithName(F.Flow2(
		either.Eitherize1(url.Parse),
		either.Fold(F.Ignore1of1[error](option.None[*url.URL]), option.Some),
	), (*url.URL).String,
		"PrismParseURL",
	)
}

// InstanceOf creates a prism for type assertions on interface{}/any values.
// It provides a safe way to extract values of a specific type from an any value,
// handling type mismatches gracefully through the Option type.
//
// The prism's GetOption attempts to assert that an any value is of type T.
// If the assertion succeeds, it returns Some(T); if it fails, it returns None.
//
// The prism's ReverseGet always succeeds, converting a value of type T back to any.
//
// Type Parameters:
//   - T: The target type to extract from any
//
// Returns:
//   - A Prism[any, T] that safely handles type assertions
//
// Example:
//
//	// Create a prism for extracting int values
//	intPrism := InstanceOf[int]()
//
//	// Extract int from any
//	var value any = 42
//	result := intPrism.GetOption(value)  // Some(42)
//
//	// Type mismatch returns None
//	var strValue any = "hello"
//	result = intPrism.GetOption(strValue)  // None[int]()
//
//	// Convert back to any
//	anyValue := intPrism.ReverseGet(42)  // any(42)
//
//	// Use with Set to update typed values
//	setter := Set[any, int](100)
//	result := setter(intPrism)(any(42))  // any(100)
//
// Common use cases:
//   - Safely extracting typed values from interface{} collections
//   - Working with heterogeneous data structures
//   - Type-safe deserialization and validation
//   - Pattern matching on interface{} values
func InstanceOf[T any]() Prism[any, T] {
	var t T
	return MakePrismWithName(option.InstanceOf[T], F.ToAny[T], fmt.Sprintf("PrismInstanceOf[%T]", t))
}

// ParseDate creates a prism for parsing and formatting dates with a specific layout.
// It provides a safe way to work with date strings, handling parsing errors
// gracefully through the Option type.
//
// The prism's GetOption attempts to parse a string into a time.Time using the
// specified layout. If parsing succeeds, it returns Some(time.Time); if it fails
// (e.g., invalid date format), it returns None.
//
// The prism's ReverseGet always succeeds, formatting a time.Time back to a string
// using the same layout.
//
// Parameters:
//   - layout: The time layout string (e.g., "2006-01-02", time.RFC3339)
//
// Returns:
//   - A Prism[string, time.Time] that safely handles date parsing/formatting
//
// Example:
//
//	// Create a prism for ISO date format
//	datePrism := ParseDate("2006-01-02")
//
//	// Parse valid date
//	parsed := datePrism.GetOption("2024-03-15")
//	// Some(time.Time{2024, 3, 15, ...})
//
//	// Parse invalid date
//	invalid := datePrism.GetOption("not-a-date")  // None[time.Time]()
//
//	// Format date back to string
//	date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
//	str := datePrism.ReverseGet(date)  // "2024-03-15"
//
//	// Use with Set to update dates
//	newDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
//	setter := Set[string, time.Time](newDate)
//	result := setter(datePrism)("2024-03-15")  // "2025-01-01"
//
//	// Different layouts for different formats
//	rfc3339Prism := ParseDate(time.RFC3339)
//	parsed = rfc3339Prism.GetOption("2024-03-15T10:30:00Z")
//
// Common use cases:
//   - Validating and parsing date configuration values
//   - Working with date strings in APIs
//   - Converting between date formats
//   - Safely handling user-provided date inputs
func ParseDate(layout string) Prism[string, time.Time] {
	return MakePrismWithName(F.Flow2(
		F.Bind1st(either.Eitherize2(time.Parse), layout),
		either.Fold(F.Ignore1of1[error](option.None[time.Time]), option.Some),
	), F.Bind2nd(time.Time.Format, layout),
		"PrismParseDate",
	)
}

// Deref creates a prism for safely dereferencing pointers.
// It provides a safe way to work with nullable pointers, handling nil values
// gracefully through the Option type.
//
// The prism's GetOption attempts to dereference a pointer.
// If the pointer is non-nil, it returns Some(*T); if it's nil, it returns None.
//
// The prism's ReverseGet is the identity function, returning the pointer unchanged.
//
// Type Parameters:
//   - T: The type being pointed to
//
// Returns:
//   - A Prism[*T, *T] that safely handles pointer dereferencing
//
// Example:
//
//	// Create a prism for dereferencing int pointers
//	derefPrism := Deref[int]()
//
//	// Dereference non-nil pointer
//	value := 42
//	ptr := &value
//	result := derefPrism.GetOption(ptr)  // Some(&42)
//
//	// Dereference nil pointer
//	var nilPtr *int
//	result = derefPrism.GetOption(nilPtr)  // None[*int]()
//
//	// ReverseGet returns the pointer unchanged
//	reconstructed := derefPrism.ReverseGet(ptr)  // &42
//
//	// Use with Set to update non-nil pointers
//	newValue := 100
//	newPtr := &newValue
//	setter := Set[*int, *int](newPtr)
//	result := setter(derefPrism)(ptr)  // &100
//	result = setter(derefPrism)(nilPtr) // nil (unchanged)
//
// Common use cases:
//   - Safely working with optional pointer fields
//   - Validating non-nil pointers before operations
//   - Filtering out nil values in data pipelines
//   - Working with database nullable columns
func Deref[T any]() Prism[*T, *T] {
	return MakePrismWithName(option.FromNillable[T], F.Identity[*T], "PrismDeref")
}

// FromEither creates a prism for extracting Right values from Either types.
// It provides a safe way to work with Either values, focusing on the success case
// and handling the error case gracefully through the Option type.
//
// The prism's GetOption attempts to extract the Right value from an Either.
// If the Either is Right(value), it returns Some(value); if it's Left(error), it returns None.
//
// The prism's ReverseGet always succeeds, wrapping a value into a Right.
//
// Type Parameters:
//   - E: The error/left type
//   - T: The value/right type
//
// Returns:
//   - A Prism[Either[E, T], T] that safely extracts Right values
//
// Example:
//
//	// Create a prism for extracting successful results
//	resultPrism := FromEither[error, int]()
//
//	// Extract from Right
//	success := either.Right[error](42)
//	result := resultPrism.GetOption(success)  // Some(42)
//
//	// Extract from Left
//	failure := either.Left[int](errors.New("failed"))
//	result = resultPrism.GetOption(failure)  // None[int]()
//
//	// Wrap value into Right
//	wrapped := resultPrism.ReverseGet(100)  // Right(100)
//
//	// Use with Set to update successful results
//	setter := Set[Either[error, int], int](200)
//	result := setter(resultPrism)(success)  // Right(200)
//	result = setter(resultPrism)(failure)   // Left(error) (unchanged)
//
// Common use cases:
//   - Extracting successful values from Either results
//   - Filtering out errors in data pipelines
//   - Working with fallible operations
//   - Composing with other prisms for complex error handling
func FromEither[E, T any]() Prism[Either[E, T], T] {
	return MakePrismWithName(either.ToOption[E, T], either.Of[E, T], "PrismFromEither")
}

// FromResult creates a prism for extracting values from Result types.
// It provides a safe way to work with Result values (which are Either[error, T]),
// focusing on the success case and handling errors gracefully through the Option type.
//
// This is a convenience function that is equivalent to FromEither[error, T]().
//
// The prism's GetOption attempts to extract the success value from a Result.
// If the Result is successful, it returns Some(value); if it's an error, it returns None.
//
// The prism's ReverseGet always succeeds, wrapping a value into a successful Result.
//
// Type Parameters:
//   - T: The value type contained in the Result
//
// Returns:
//   - A Prism[Result[T], T] that safely extracts success values
//
// Example:
//
//	// Create a prism for extracting successful results
//	resultPrism := FromResult[int]()
//
//	// Extract from successful result
//	success := result.Of[int](42)
//	value := resultPrism.GetOption(success)  // Some(42)
//
//	// Extract from error result
//	failure := result.Error[int](errors.New("failed"))
//	value = resultPrism.GetOption(failure)  // None[int]()
//
//	// Wrap value into successful Result
//	wrapped := resultPrism.ReverseGet(100)  // Result containing 100
//
//	// Use with Set to update successful results
//	setter := Set[Result[int], int](200)
//	result := setter(resultPrism)(success)  // Result containing 200
//	result = setter(resultPrism)(failure)   // Error result (unchanged)
//
// Common use cases:
//   - Extracting successful values from Result types
//   - Filtering out errors in data pipelines
//   - Working with fallible operations that return Result
//   - Composing with other prisms for complex error handling
//
//go:inline
func FromResult[T any]() Prism[Result[T], T] {
	return FromEither[error, T]()
}

// FromZero creates a prism that matches zero values of comparable types.
// It provides a safe way to work with zero values, handling non-zero values
// gracefully through the Option type.
//
// The prism's GetOption returns Some(t) if the value equals the zero value
// of type T; otherwise, it returns None.
//
// The prism's ReverseGet is the identity function, returning the value unchanged.
//
// Type Parameters:
//   - T: A comparable type (must support == and != operators)
//
// Returns:
//   - A Prism[T, T] that matches zero values
//
// Example:
//
//	// Create a prism for zero integers
//	zeroPrism := FromZero[int]()
//
//	// Match zero value
//	result := zeroPrism.GetOption(0)  // Some(0)
//
//	// Non-zero returns None
//	result = zeroPrism.GetOption(42)  // None[int]()
//
//	// ReverseGet is identity
//	value := zeroPrism.ReverseGet(0)  // 0
//
//	// Use with Set to update zero values
//	setter := Set[int, int](100)
//	result := setter(zeroPrism)(0)   // 100
//	result = setter(zeroPrism)(42)   // 42 (unchanged)
//
// Common use cases:
//   - Validating that values are zero/default
//   - Filtering zero values in data pipelines
//   - Working with optional fields that use zero as "not set"
//   - Replacing zero values with defaults
func FromZero[T comparable]() Prism[T, T] {
	return MakePrismWithName(option.FromZero[T](), F.Identity[T], "PrismFromZero")
}

// FromNonZero creates a prism that matches non-zero values of comparable types.
// It provides a safe way to work with non-zero values, handling zero values
// gracefully through the Option type.
//
// The prism's GetOption returns Some(t) if the value is not equal to the zero value
// of type T; otherwise, it returns None.
//
// The prism's ReverseGet is the identity function, returning the value unchanged.
//
// Type Parameters:
//   - T: A comparable type (must support == and != operators)
//
// Returns:
//   - A Prism[T, T] that matches non-zero values
//
// Example:
//
//	// Create a prism for non-zero integers
//	nonZeroPrism := FromNonZero[int]()
//
//	// Match non-zero value
//	result := nonZeroPrism.GetOption(42)  // Some(42)
//
//	// Zero returns None
//	result = nonZeroPrism.GetOption(0)  // None[int]()
//
//	// ReverseGet is identity
//	value := nonZeroPrism.ReverseGet(42)  // 42
//
//	// Use with Set to update non-zero values
//	setter := Set[int, int](100)
//	result := setter(nonZeroPrism)(42)   // 100
//	result = setter(nonZeroPrism)(0)     // 0 (unchanged)
//
// Common use cases:
//   - Validating that values are non-zero/non-default
//   - Filtering non-zero values in data pipelines
//   - Working with required fields that shouldn't be zero
//   - Replacing non-zero values with new values
func FromNonZero[T comparable]() Prism[T, T] {
	return MakePrismWithName(option.FromNonZero[T](), F.Identity[T], "PrismFromNonZero")
}

// Match represents a regex match result with full reconstruction capability.
// It contains everything needed to reconstruct the original string, making it
// suitable for use in a prism that maintains bidirectionality.
//
// Fields:
//   - Before: Text before the match
//   - Groups: Capture groups (index 0 is the full match, 1+ are capture groups)
//   - After: Text after the match
//
// Example:
//
//	// For string "hello world 123" with regex `\d+`:
//	// Match{
//	//     Before: "hello world ",
//	//     Groups: []string{"123"},
//	//     After: "",
//	// }
//
// fp-go:Lens
type Match struct {
	Before string   // Text before the match
	Groups []string // Capture groups (index 0 is full match)
	After  string   // Text after the match
}

// Reconstruct builds the original string from a Match.
// This is the inverse operation of regex matching, allowing full round-trip conversion.
//
// Returns:
//   - The original string that was matched
//
// Example:
//
//	match := Match{
//	    Before: "hello ",
//	    Groups: []string{"world"},
//	    After: "!",
//	}
//	original := match.Reconstruct()  // "hello world!"
func (m Match) Reconstruct() string {
	return m.Before + m.Groups[0] + m.After
}

// FullMatch returns the complete matched text (the entire regex match).
// This is equivalent to Groups[0] and represents what the regex matched.
//
// Returns:
//   - The full matched text
//
// Example:
//
//	match := Match{
//	    Before: "price: ",
//	    Groups: []string{"$99.99", "99.99"},
//	    After: " USD",
//	}
//	full := match.FullMatch()  // "$99.99"
func (m Match) FullMatch() string {
	return m.Groups[0]
}

// Group returns the nth capture group from the match (1-indexed).
// Capture group 0 is the full match, groups 1+ are the parenthesized captures.
// Returns an empty string if the group index is out of bounds.
//
// Parameters:
//   - n: The capture group index (1-indexed)
//
// Returns:
//   - The captured text, or empty string if index is invalid
//
// Example:
//
//	// Regex: `(\w+)@(\w+\.\w+)` matching "user@example.com"
//	match := Match{
//	    Groups: []string{"user@example.com", "user", "example.com"},
//	}
//	username := match.Group(1)  // "user"
//	domain := match.Group(2)    // "example.com"
//	invalid := match.Group(5)   // ""
func (m Match) Group(n int) string {
	if n < len(m.Groups) {
		return m.Groups[n]
	}
	return ""
}

// RegexMatcher creates a prism for regex pattern matching with full reconstruction.
// It provides a safe way to match strings against a regex pattern, extracting
// match information while maintaining the ability to reconstruct the original string.
//
// The prism's GetOption attempts to match the regex against the string.
// If a match is found, it returns Some(Match) with all capture groups and context;
// if no match is found, it returns None.
//
// The prism's ReverseGet reconstructs the original string from a Match.
//
// Parameters:
//   - re: A compiled regular expression
//
// Returns:
//   - A Prism[string, Match] that safely handles regex matching
//
// Example:
//
//	// Create a prism for matching numbers
//	numRegex := regexp.MustCompile(`\d+`)
//	numPrism := RegexMatcher(numRegex)
//
//	// Match a string
//	match := numPrism.GetOption("price: 42 dollars")
//	// Some(Match{Before: "price: ", Groups: ["42"], After: " dollars"})
//
//	// No match returns None
//	noMatch := numPrism.GetOption("no numbers here")  // None[Match]()
//
//	// Reconstruct original string
//	if m, ok := option.IsSome(match); ok {
//	    original := numPrism.ReverseGet(m)  // "price: 42 dollars"
//	}
//
//	// Extract capture groups
//	emailRegex := regexp.MustCompile(`(\w+)@(\w+\.\w+)`)
//	emailPrism := RegexMatcher(emailRegex)
//	match = emailPrism.GetOption("contact: user@example.com")
//	// Match.Group(1) = "user", Match.Group(2) = "example.com"
//
// Common use cases:
//   - Parsing structured text with regex patterns
//   - Extracting and validating data from strings
//   - Text transformation pipelines
//   - Pattern-based string manipulation with reconstruction
//
// Note: This prism is bijective - you can always reconstruct the original
// string from a Match, making it suitable for round-trip transformations.
func RegexMatcher(re *regexp.Regexp) Prism[string, Match] {
	noMatch := option.None[Match]()

	return MakePrismWithName(
		// String -> Option[Match]
		func(s string) Option[Match] {
			loc := re.FindStringSubmatchIndex(s)
			if loc == nil {
				return noMatch
			}

			// Extract all capture groups
			groups := make([]string, 0)
			for i := 0; i < len(loc); i += 2 {
				if loc[i] >= 0 {
					groups = append(groups, s[loc[i]:loc[i+1]])
				} else {
					groups = append(groups, "")
				}
			}

			match := Match{
				Before: s[:loc[0]],
				Groups: groups,
				After:  s[loc[1]:],
			}

			return option.Some(match)
		},
		Match.Reconstruct,
		fmt.Sprintf("PrismRegex[%s]", re),
	)
}

// NamedMatch represents a regex match result with named capture groups.
// It provides access to captured text by name rather than by index, making
// regex patterns more readable and maintainable.
//
// Fields:
//   - Before: Text before the match
//   - Groups: Map of capture group names to their matched text
//   - Full: The complete matched text
//   - After: Text after the match
//
// Example:
//
//	// For regex `(?P<user>\w+)@(?P<domain>\w+\.\w+)` matching "user@example.com":
//	// NamedMatch{
//	//     Before: "",
//	//     Groups: map[string]string{"user": "user", "domain": "example.com"},
//	//     Full: "user@example.com",
//	//     After: "",
//	// }
//
// fp-go:Lens
type NamedMatch struct {
	Before string
	Groups map[string]string
	Full   string // The full matched text
	After  string
}

// Reconstruct builds the original string from a NamedMatch.
// This is the inverse operation of regex matching, allowing full round-trip conversion.
//
// Returns:
//   - The original string that was matched
//
// Example:
//
//	match := NamedMatch{
//	    Before: "email: ",
//	    Full: "user@example.com",
//	    Groups: map[string]string{"user": "user", "domain": "example.com"},
//	    After: "",
//	}
//	original := match.Reconstruct()  // "email: user@example.com"
func (nm NamedMatch) Reconstruct() string {
	return nm.Before + nm.Full + nm.After
}

// RegexNamedMatcher creates a prism for regex pattern matching with named capture groups.
// It provides a safe way to match strings against a regex pattern with named groups,
// making it easier to extract specific parts of the match by name rather than index.
//
// The prism's GetOption attempts to match the regex against the string.
// If a match is found, it returns Some(NamedMatch) with all named capture groups;
// if no match is found, it returns None.
//
// The prism's ReverseGet reconstructs the original string from a NamedMatch.
//
// Parameters:
//   - re: A compiled regular expression with named capture groups
//
// Returns:
//   - A Prism[string, NamedMatch] that safely handles regex matching with named groups
//
// Example:
//
//	// Create a prism for matching email addresses with named groups
//	emailRegex := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
//	emailPrism := RegexNamedMatcher(emailRegex)
//
//	// Match a string
//	match := emailPrism.GetOption("contact: user@example.com")
//	// Some(NamedMatch{
//	//     Before: "contact: ",
//	//     Groups: {"user": "user", "domain": "example.com"},
//	//     Full: "user@example.com",
//	//     After: "",
//	// })
//
//	// Access named groups
//	if m, ok := option.IsSome(match); ok {
//	    username := m.Groups["user"]      // "user"
//	    domain := m.Groups["domain"]      // "example.com"
//	}
//
//	// No match returns None
//	noMatch := emailPrism.GetOption("invalid-email")  // None[NamedMatch]()
//
//	// Reconstruct original string
//	if m, ok := option.IsSome(match); ok {
//	    original := emailPrism.ReverseGet(m)  // "contact: user@example.com"
//	}
//
//	// More complex example with date parsing
//	dateRegex := regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})`)
//	datePrism := RegexNamedMatcher(dateRegex)
//	match = datePrism.GetOption("Date: 2024-03-15")
//	// Access: match.Groups["year"], match.Groups["month"], match.Groups["day"]
//
// Common use cases:
//   - Parsing structured text with meaningful field names
//   - Extracting and validating data from formatted strings
//   - Log parsing with named fields
//   - Configuration file parsing
//   - URL route parameter extraction
//
// Note: Only named capture groups appear in the Groups map. Unnamed groups
// are not included. The Full field always contains the complete matched text.
func RegexNamedMatcher(re *regexp.Regexp) Prism[string, NamedMatch] {
	names := re.SubexpNames()
	noMatch := option.None[NamedMatch]()

	return MakePrism(
		func(s string) Option[NamedMatch] {
			loc := re.FindStringSubmatchIndex(s)
			if loc == nil {
				return noMatch
			}

			groups := make(map[string]string)
			for i := 1; i < len(loc)/2; i++ {
				if S.IsNonEmpty(names[i]) && loc[2*i] >= 0 {
					groups[names[i]] = s[loc[2*i]:loc[2*i+1]]
				}
			}

			match := NamedMatch{
				Before: s[:loc[0]],
				Groups: groups,
				Full:   s[loc[0]:loc[1]],
				After:  s[loc[1]:],
			}

			return option.Some(match)
		},
		NamedMatch.Reconstruct,
	)
}

func getFromEither[A, B any](f func(A) (B, error)) func(A) Option[B] {
	return func(a A) Option[B] {
		b, err := f(a)
		if err != nil {
			return option.None[B]()
		}
		return option.Of(b)
	}
}

func atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}

// ParseInt creates a prism for parsing and formatting integers.
// It provides a safe way to convert between string and int, handling
// parsing errors gracefully through the Option type.
//
// The prism's GetOption attempts to parse a string into an int.
// If parsing succeeds, it returns Some(int); if it fails (e.g., invalid
// number format), it returns None.
//
// The prism's ReverseGet always succeeds, converting an int to its string representation.
//
// Returns:
//   - A Prism[string, int] that safely handles int parsing/formatting
//
// Example:
//
//	// Create an int parsing prism
//	intPrism := ParseInt()
//
//	// Parse valid integer
//	parsed := intPrism.GetOption("42")  // Some(42)
//
//	// Parse invalid integer
//	invalid := intPrism.GetOption("not-a-number")  // None[int]()
//
//	// Format int to string
//	str := intPrism.ReverseGet(42)  // "42"
//
//	// Use with Set to update integer values
//	setter := Set[string, int](100)
//	result := setter(intPrism)("42")  // "100"
//
// Common use cases:
//   - Parsing integer configuration values
//   - Validating numeric user input
//   - Converting between string and int in data pipelines
//   - Working with numeric API parameters
//
//go:inline
func ParseInt() Prism[string, int] {
	return MakePrismWithName(getFromEither(strconv.Atoi), strconv.Itoa, "PrismParseInt")
}

// ParseInt64 creates a prism for parsing and formatting 64-bit integers.
// It provides a safe way to convert between string and int64, handling
// parsing errors gracefully through the Option type.
//
// The prism's GetOption attempts to parse a string into an int64.
// If parsing succeeds, it returns Some(int64); if it fails (e.g., invalid
// number format or overflow), it returns None.
//
// The prism's ReverseGet always succeeds, converting an int64 to its string representation.
//
// Returns:
//   - A Prism[string, int64] that safely handles int64 parsing/formatting
//
// Example:
//
//	// Create an int64 parsing prism
//	int64Prism := ParseInt64()
//
//	// Parse valid 64-bit integer
//	parsed := int64Prism.GetOption("9223372036854775807")  // Some(9223372036854775807)
//
//	// Parse invalid integer
//	invalid := int64Prism.GetOption("not-a-number")  // None[int64]()
//
//	// Format int64 to string
//	str := int64Prism.ReverseGet(int64(42))  // "42"
//
//	// Use with Set to update int64 values
//	setter := Set[string, int64](int64(100))
//	result := setter(int64Prism)("42")  // "100"
//
// Common use cases:
//   - Parsing large integer values (timestamps, IDs)
//   - Working with database integer columns
//   - Handling 64-bit numeric API parameters
//   - Converting between string and int64 in data pipelines
//
//go:inline
func ParseInt64() Prism[string, int64] {
	return MakePrismWithName(getFromEither(atoi64), itoa64, "PrismParseInt64")
}

// ParseBool creates a prism for parsing and formatting boolean values.
// It provides a safe way to convert between string and bool, handling
// parsing errors gracefully through the Option type.
//
// The prism's GetOption attempts to parse a string into a bool.
// It accepts "1", "t", "T", "TRUE", "true", "True", "0", "f", "F", "FALSE", "false", "False".
// If parsing succeeds, it returns Some(bool); if it fails, it returns None.
//
// The prism's ReverseGet always succeeds, converting a bool to "true" or "false".
//
// Returns:
//   - A Prism[string, bool] that safely handles bool parsing/formatting
//
// Example:
//
//	// Create a bool parsing prism
//	boolPrism := ParseBool()
//
//	// Parse valid boolean strings
//	parsed := boolPrism.GetOption("true")   // Some(true)
//	parsed = boolPrism.GetOption("1")       // Some(true)
//	parsed = boolPrism.GetOption("false")   // Some(false)
//	parsed = boolPrism.GetOption("0")       // Some(false)
//
//	// Parse invalid boolean
//	invalid := boolPrism.GetOption("maybe")  // None[bool]()
//
//	// Format bool to string
//	str := boolPrism.ReverseGet(true)   // "true"
//	str = boolPrism.ReverseGet(false)   // "false"
//
//	// Use with Set to update boolean values
//	setter := Set[string, bool](true)
//	result := setter(boolPrism)("false")  // "true"
//
// Common use cases:
//   - Parsing boolean configuration values
//   - Validating boolean user input
//   - Converting between string and bool in data pipelines
//   - Working with boolean API parameters or flags
//
//go:inline
func ParseBool() Prism[string, bool] {
	return MakePrismWithName(getFromEither(strconv.ParseBool), strconv.FormatBool, "PrismParseBool")
}

func atof64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func atof32(s string) (float32, error) {
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f32), nil
}

func f32toa(f float32) string {
	return strconv.FormatFloat(float64(f), 'g', -1, 32)
}

func f64toa(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// ParseFloat32 creates a prism for parsing and formatting 32-bit floating-point numbers.
// It provides a safe way to convert between string and float32, handling
// parsing errors gracefully through the Option type.
//
// The prism's GetOption attempts to parse a string into a float32.
// If parsing succeeds, it returns Some(float32); if it fails (e.g., invalid
// number format or overflow), it returns None.
//
// The prism's ReverseGet always succeeds, converting a float32 to its string representation
// using the 'g' format (shortest representation).
//
// Returns:
//   - A Prism[string, float32] that safely handles float32 parsing/formatting
//
// Example:
//
//	// Create a float32 parsing prism
//	float32Prism := ParseFloat32()
//
//	// Parse valid float
//	parsed := float32Prism.GetOption("3.14")  // Some(3.14)
//	parsed = float32Prism.GetOption("1.5e10") // Some(1.5e10)
//
//	// Parse invalid float
//	invalid := float32Prism.GetOption("not-a-number")  // None[float32]()
//
//	// Format float32 to string
//	str := float32Prism.ReverseGet(float32(3.14))  // "3.14"
//
//	// Use with Set to update float32 values
//	setter := Set[string, float32](float32(2.71))
//	result := setter(float32Prism)("3.14")  // "2.71"
//
// Common use cases:
//   - Parsing floating-point configuration values
//   - Working with scientific notation
//   - Converting between string and float32 in data pipelines
//   - Handling numeric API parameters with decimal precision
//
//go:inline
func ParseFloat32() Prism[string, float32] {
	return MakePrismWithName(getFromEither(atof32), f32toa, "ParseFloat32")
}

// ParseFloat64 creates a prism for parsing and formatting 64-bit floating-point numbers.
// It provides a safe way to convert between string and float64, handling
// parsing errors gracefully through the Option type.
//
// The prism's GetOption attempts to parse a string into a float64.
// If parsing succeeds, it returns Some(float64); if it fails (e.g., invalid
// number format or overflow), it returns None.
//
// The prism's ReverseGet always succeeds, converting a float64 to its string representation
// using the 'g' format (shortest representation).
//
// Returns:
//   - A Prism[string, float64] that safely handles float64 parsing/formatting
//
// Example:
//
//	// Create a float64 parsing prism
//	float64Prism := ParseFloat64()
//
//	// Parse valid float
//	parsed := float64Prism.GetOption("3.141592653589793")  // Some(3.141592653589793)
//	parsed = float64Prism.GetOption("1.5e100")             // Some(1.5e100)
//
//	// Parse invalid float
//	invalid := float64Prism.GetOption("not-a-number")  // None[float64]()
//
//	// Format float64 to string
//	str := float64Prism.ReverseGet(3.141592653589793)  // "3.141592653589793"
//
//	// Use with Set to update float64 values
//	setter := Set[string, float64](2.718281828459045)
//	result := setter(float64Prism)("3.14")  // "2.718281828459045"
//
// Common use cases:
//   - Parsing high-precision floating-point values
//   - Working with scientific notation and large numbers
//   - Converting between string and float64 in data pipelines
//   - Handling precise numeric API parameters
//
//go:inline
func ParseFloat64() Prism[string, float64] {
	return MakePrismWithName(getFromEither(atof64), f64toa, "PrismParseFloat64")
}

// FromOption creates a prism for extracting values from Option types.
// It provides a safe way to work with Option values, focusing on the Some case
// and handling the None case gracefully through the prism's GetOption behavior.
//
// The prism's GetOption is the identity function - it returns the Option as-is.
// If the Option is Some(value), GetOption returns Some(value); if it's None, it returns None.
// This allows the prism to naturally handle the presence or absence of a value.
//
// The prism's ReverseGet wraps a value into Some, always succeeding.
//
// Type Parameters:
//   - T: The value type contained in the Option
//
// Returns:
//   - A Prism[Option[T], T] that safely extracts values from Options
//
// Example:
//
//	// Create a prism for extracting int values from Option[int]
//	optPrism := FromOption[int]()
//
//	// Extract from Some
//	someValue := option.Some(42)
//	result := optPrism.GetOption(someValue)  // Some(42)
//
//	// Extract from None
//	noneValue := option.None[int]()
//	result = optPrism.GetOption(noneValue)  // None[int]()
//
//	// Wrap value into Some
//	wrapped := optPrism.ReverseGet(100)  // Some(100)
//
//	// Use with Set to update Some values
//	setter := Set[Option[int], int](200)
//	result := setter(optPrism)(someValue)  // Some(200)
//	result = setter(optPrism)(noneValue)   // None[int]() (unchanged)
//
//	// Compose with other prisms for nested extraction
//	// Extract int from Option[Option[int]]
//	nestedPrism := Compose[Option[Option[int]], Option[int], int](
//	    FromOption[Option[int]](),
//	    FromOption[int](),
//	)
//	nested := option.Some(option.Some(42))
//	value := nestedPrism.GetOption(nested)  // Some(42)
//
// Common use cases:
//   - Extracting values from optional fields
//   - Working with nullable data in a type-safe way
//   - Composing with other prisms to handle nested Options
//   - Filtering and transforming optional values in pipelines
//   - Converting between Option and other optional representations
//
// Key insight: This prism treats Option[T] as a "container" that may or may not
// hold a value of type T. The prism focuses on the value inside, allowing you to
// work with it when present and gracefully handle its absence when not.
//
//go:inline
func FromOption[T any]() Prism[Option[T], T] {
	return MakePrismWithName(
		F.Identity[Option[T]],
		option.Some[T],
		"PrismFromOption",
	)
}

// NonEmptyString creates a prism that matches non-empty strings.
// It provides a safe way to work with non-empty string values, handling
// empty strings gracefully through the Option type.
//
// This is a specialized version of FromNonZero[string]() that makes the intent
// clearer when working specifically with strings that must not be empty.
//
// The prism's GetOption returns Some(s) if the string is not empty;
// otherwise, it returns None.
//
// The prism's ReverseGet is the identity function, returning the string unchanged.
//
// Returns:
//   - A Prism[string, string] that matches non-empty strings
//
// Example:
//
//	// Create a prism for non-empty strings
//	nonEmptyPrism := NonEmptyString()
//
//	// Match non-empty string
//	result := nonEmptyPrism.GetOption("hello")  // Some("hello")
//
//	// Empty string returns None
//	result = nonEmptyPrism.GetOption("")  // None[string]()
//
//	// ReverseGet is identity
//	value := nonEmptyPrism.ReverseGet("world")  // "world"
//
//	// Use with Set to update non-empty strings
//	setter := Set[string, string]("updated")
//	result := setter(nonEmptyPrism)("original")  // "updated"
//	result = setter(nonEmptyPrism)("")           // "" (unchanged)
//
//	// Compose with other prisms for validation pipelines
//	// Example: Parse a non-empty string as an integer
//	nonEmptyIntPrism := Compose[string, string, int](
//	    NonEmptyString(),
//	    ParseInt(),
//	)
//	value := nonEmptyIntPrism.GetOption("42")   // Some(42)
//	value = nonEmptyIntPrism.GetOption("")      // None[int]()
//	value = nonEmptyIntPrism.GetOption("abc")   // None[int]()
//
// Common use cases:
//   - Validating required string fields (usernames, names, IDs)
//   - Filtering empty strings from data pipelines
//   - Ensuring configuration values are non-empty
//   - Composing with parsing prisms to validate input before parsing
//   - Working with user input that must not be blank
//
// Key insight: This prism is particularly useful for validation scenarios where
// an empty string represents an invalid or missing value, allowing you to handle
// such cases gracefully through the Option type rather than with error handling.
//
//go:inline
func NonEmptyString() Prism[string, string] {
	return FromNonZero[string]()
}

// ErrorPrisms provides prisms for accessing fields of url.Error
type ErrorPrisms struct {
	Op  Prism[url.Error, string]
	URL Prism[url.Error, string]
	Err Prism[url.Error, error]
}

// MakeErrorPrisms creates a new ErrorPrisms with prisms for all fields
func MakeErrorPrisms() ErrorPrisms {
	_fromNonZeroOp := option.FromNonZero[string]()
	_prismOp := MakePrismWithName(
		func(s url.Error) Option[string] { return _fromNonZeroOp(s.Op) },
		func(v string) url.Error {
			return url.Error{Op: v}
		},
		"Error.Op",
	)
	_fromNonZeroURL := option.FromNonZero[string]()
	_prismURL := MakePrismWithName(
		func(s url.Error) Option[string] { return _fromNonZeroURL(s.URL) },
		func(v string) url.Error {
			return url.Error{URL: v}
		},
		"Error.URL",
	)
	_fromNonZeroErr := option.FromNonZero[error]()
	_prismErr := MakePrismWithName(
		func(s url.Error) Option[error] { return _fromNonZeroErr(s.Err) },
		func(v error) url.Error {
			return url.Error{Err: v}
		},
		"Error.Err",
	)
	return ErrorPrisms{
		Op:  _prismOp,
		URL: _prismURL,
		Err: _prismErr,
	}
}

// URLPrisms provides prisms for accessing fields of url.URL
type URLPrisms struct {
	Scheme      Prism[url.URL, string]
	Opaque      Prism[url.URL, string]
	User        Prism[url.URL, *url.Userinfo]
	Host        Prism[url.URL, string]
	Path        Prism[url.URL, string]
	RawPath     Prism[url.URL, string]
	OmitHost    Prism[url.URL, bool]
	ForceQuery  Prism[url.URL, bool]
	RawQuery    Prism[url.URL, string]
	Fragment    Prism[url.URL, string]
	RawFragment Prism[url.URL, string]
}

// MakeURLPrisms creates a new URLPrisms with prisms for all fields
func MakeURLPrisms() URLPrisms {
	_fromNonZeroScheme := option.FromNonZero[string]()
	_prismScheme := MakePrismWithName(
		func(s url.URL) Option[string] { return _fromNonZeroScheme(s.Scheme) },
		func(v string) url.URL {
			return url.URL{Scheme: v}
		},
		"URL.Scheme",
	)
	_fromNonZeroOpaque := option.FromNonZero[string]()
	_prismOpaque := MakePrismWithName(
		func(s url.URL) Option[string] { return _fromNonZeroOpaque(s.Opaque) },
		func(v string) url.URL {
			return url.URL{Opaque: v}
		},
		"URL.Opaque",
	)
	_fromNonZeroUser := option.FromNonZero[*url.Userinfo]()
	_prismUser := MakePrismWithName(
		func(s url.URL) Option[*url.Userinfo] { return _fromNonZeroUser(s.User) },
		func(v *url.Userinfo) url.URL {
			return url.URL{User: v}
		},
		"URL.User",
	)
	_fromNonZeroHost := option.FromNonZero[string]()
	_prismHost := MakePrismWithName(
		func(s url.URL) Option[string] { return _fromNonZeroHost(s.Host) },
		func(v string) url.URL {
			return url.URL{Host: v}
		},
		"URL.Host",
	)
	_fromNonZeroPath := option.FromNonZero[string]()
	_prismPath := MakePrismWithName(
		func(s url.URL) Option[string] { return _fromNonZeroPath(s.Path) },
		func(v string) url.URL {
			return url.URL{Path: v}
		},
		"URL.Path",
	)
	_fromNonZeroRawPath := option.FromNonZero[string]()
	_prismRawPath := MakePrismWithName(
		func(s url.URL) Option[string] { return _fromNonZeroRawPath(s.RawPath) },
		func(v string) url.URL {
			return url.URL{RawPath: v}
		},
		"URL.RawPath",
	)
	_fromNonZeroOmitHost := option.FromNonZero[bool]()
	_prismOmitHost := MakePrismWithName(
		func(s url.URL) Option[bool] { return _fromNonZeroOmitHost(s.OmitHost) },
		func(v bool) url.URL {
			return url.URL{OmitHost: v}
		},
		"URL.OmitHost",
	)
	_fromNonZeroForceQuery := option.FromNonZero[bool]()
	_prismForceQuery := MakePrismWithName(
		func(s url.URL) Option[bool] { return _fromNonZeroForceQuery(s.ForceQuery) },
		func(v bool) url.URL {
			return url.URL{ForceQuery: v}
		},
		"URL.ForceQuery",
	)
	_fromNonZeroRawQuery := option.FromNonZero[string]()
	_prismRawQuery := MakePrismWithName(
		func(s url.URL) Option[string] { return _fromNonZeroRawQuery(s.RawQuery) },
		func(v string) url.URL {
			return url.URL{RawQuery: v}
		},
		"URL.RawQuery",
	)
	_fromNonZeroFragment := option.FromNonZero[string]()
	_prismFragment := MakePrismWithName(
		func(s url.URL) Option[string] { return _fromNonZeroFragment(s.Fragment) },
		func(v string) url.URL {
			return url.URL{Fragment: v}
		},
		"URL.Fragment",
	)
	_fromNonZeroRawFragment := option.FromNonZero[string]()
	_prismRawFragment := MakePrismWithName(
		func(s url.URL) Option[string] { return _fromNonZeroRawFragment(s.RawFragment) },
		func(v string) url.URL {
			return url.URL{RawFragment: v}
		},
		"URL.RawFragment",
	)
	return URLPrisms{
		Scheme:      _prismScheme,
		Opaque:      _prismOpaque,
		User:        _prismUser,
		Host:        _prismHost,
		Path:        _prismPath,
		RawPath:     _prismRawPath,
		OmitHost:    _prismOmitHost,
		ForceQuery:  _prismForceQuery,
		RawQuery:    _prismRawQuery,
		Fragment:    _prismFragment,
		RawFragment: _prismRawFragment,
	}
}

// ParseJSON creates a prism for parsing and marshaling JSON data.
// It provides a safe way to convert between JSON bytes and Go types,
// handling parsing and marshaling errors gracefully through the Option type.
//
// The prism's GetOption attempts to unmarshal JSON bytes into type A.
// If unmarshaling succeeds, it returns Some(A); if it fails (e.g., invalid JSON
// or type mismatch), it returns None.
//
// The prism's ReverseGet marshals a value of type A into JSON bytes.
// If marshaling fails (which is rare), it returns an empty byte slice.
//
// Type Parameters:
//   - A: The Go type to unmarshal JSON into
//
// Returns:
//   - A Prism[[]byte, A] that safely handles JSON parsing/marshaling
//
// Example:
//
//	// Define a struct type
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	// Create a JSON parsing prism
//	jsonPrism := ParseJSON[Person]()
//
//	// Parse valid JSON
//	jsonData := []byte(`{"name":"Alice","age":30}`)
//	person := jsonPrism.GetOption(jsonData)
//	// Some(Person{Name: "Alice", Age: 30})
//
//	// Parse invalid JSON
//	invalidJSON := []byte(`{invalid json}`)
//	result := jsonPrism.GetOption(invalidJSON)  // None[Person]()
//
//	// Marshal to JSON
//	p := Person{Name: "Bob", Age: 25}
//	jsonBytes := jsonPrism.ReverseGet(p)
//	// []byte(`{"name":"Bob","age":25}`)
//
//	// Use with Set to update JSON data
//	newPerson := Person{Name: "Charlie", Age: 35}
//	setter := Set[[]byte, Person](newPerson)
//	updated := setter(jsonPrism)(jsonData)
//	// []byte(`{"name":"Charlie","age":35}`)
//
// Common use cases:
//   - Parsing JSON configuration files
//   - Working with JSON API responses
//   - Validating and transforming JSON data in pipelines
//   - Type-safe JSON deserialization
//   - Converting between JSON and Go structs
func ParseJSON[A any]() Prism[[]byte, A] {
	return MakePrismWithName(
		F.Flow2(
			J.Unmarshal[A],
			either.ToOption[error, A],
		),
		F.Flow2(
			J.Marshal[A],
			either.GetOrElse(F.Constant1[error](array.Empty[byte]())),
		),
		"JSON",
	)
}
