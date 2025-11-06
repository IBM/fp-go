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
