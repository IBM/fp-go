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
	"net/url"
	"regexp"
	"time"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
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
	return MakePrism(F.Flow2(
		either.Eitherize1(enc.DecodeString),
		either.Fold(F.Ignore1of1[error](option.None[[]byte]), option.Some),
	), enc.EncodeToString)
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
	return MakePrism(F.Flow2(
		either.Eitherize1(url.Parse),
		either.Fold(F.Ignore1of1[error](option.None[*url.URL]), option.Some),
	), (*url.URL).String)
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
	return MakePrism(option.ToType[T], F.ToAny[T])
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
	return MakePrism(F.Flow2(
		F.Bind1st(either.Eitherize2(time.Parse), layout),
		either.Fold(F.Ignore1of1[error](option.None[time.Time]), option.Some),
	), F.Bind2nd(time.Time.Format, layout))
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
	return MakePrism(option.FromNillable[T], F.Identity[*T])
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
	return MakePrism(either.ToOption[E, T], either.Of[E, T])
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
	return MakePrism(option.FromZero[T](), F.Identity[T])
}

func FromNonZero[T comparable]() Prism[T, T] {
	return MakePrism(option.FromNonZero[T](), F.Identity[T])
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

	return MakePrism(
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
				if names[i] != "" && loc[2*i] >= 0 {
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
