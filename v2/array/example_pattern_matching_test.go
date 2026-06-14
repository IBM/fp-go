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

package array_test

import (
	"fmt"
	"strconv"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	RO "github.com/IBM/fp-go/v2/readeroption"
)

// Example_pattern_matching demonstrates using FindFirstMap with option.Alt for
// multi-branch pattern matching.
//
// FindFirstMap searches an array and applies a selector function to each element,
// returning the first Some result. When combined with option.Alt, you can create
// sophisticated pattern matchers that try multiple conditions in sequence.
//
// For pattern matching on a single value (similar to switch/case), combine
// multiple matcher functions using option.Alt via function.Flow.
func Example_pattern_matching() {
	// Define a type to classify
	type Request struct {
		Method string
		Path   string
		Body   string
	}

	// Define matchers as functions that return Some on match, None otherwise
	matchGET := func(r Request) O.Option[string] {
		if r.Method == "GET" {
			return O.Some(fmt.Sprintf("Fetching: %s", r.Path))
		}
		return O.None[string]()
	}

	matchPOST := func(r Request) O.Option[string] {
		if r.Method == "POST" {
			return O.Some(fmt.Sprintf("Creating: %s with body: %s", r.Path, r.Body))
		}
		return O.None[string]()
	}

	matchDELETE := func(r Request) O.Option[string] {
		if r.Method == "DELETE" {
			return O.Some(fmt.Sprintf("Deleting: %s", r.Path))
		}
		return O.None[string]()
	}

	defaultCase := func(r Request) string {
		return fmt.Sprintf("Unsupported method: %s", r.Method)
	}

	matchers := A.From(
		matchGET,
		matchPOST,
		matchDELETE,
	)

	altMonoid := RO.AltMonoid[Request, string]()

	handleRequest := F.Pipe2(
		matchers,
		A.Fold(altMonoid),
		RO.GetOrElse(defaultCase),
	)

	// Test various requests
	requests := []Request{
		{Method: "GET", Path: "/users"},
		{Method: "POST", Path: "/users", Body: `{"name":"Alice"}`},
		{Method: "DELETE", Path: "/users/123"},
		{Method: "PATCH", Path: "/users/123"},
	}

	for _, req := range requests {
		fmt.Println(handleRequest(req))
	}

	// Output:
	// Fetching: /users
	// Creating: /users with body: {"name":"Alice"}
	// Deleting: /users/123
	// Unsupported method: PATCH
}

// Example_pattern_matching_array demonstrates using FindFirstMap to find and
// transform the first matching element in an array.
func Example_pattern_matching_array() {
	// Parse different string formats into integers
	parseDecimal := func(s string) O.Option[int] {
		if len(s) > 0 && s[0] != '0' {
			n, err := strconv.Atoi(s)
			if err == nil {
				return O.Some(n)
			}
		}
		return O.None[int]()
	}

	parseHex := func(s string) O.Option[int] {
		if len(s) > 2 && s[:2] == "0x" {
			n, err := strconv.ParseInt(s[2:], 16, 64)
			if err == nil {
				return O.Some(int(n))
			}
		}
		return O.None[int]()
	}

	parseOctal := func(s string) O.Option[int] {
		if len(s) > 1 && s[0] == '0' && s[1] != 'x' {
			n, err := strconv.ParseInt(s[1:], 8, 64)
			if err == nil {
				return O.Some(int(n))
			}
		}
		return O.None[int]()
	}

	parseBinary := func(s string) O.Option[int] {
		if len(s) > 2 && s[:2] == "0b" {
			n, err := strconv.ParseInt(s[2:], 2, 64)
			if err == nil {
				return O.Some(int(n))
			}
		}
		return O.None[int]()
	}

	// Combine parsers using AltAllArray - tries each format in sequence
	parseNumber := func(s string) O.Option[int] {
		parsers := []O.Option[int]{
			parseDecimal(s),
			parseHex(s),
			parseOctal(s),
			parseBinary(s),
		}
		return O.AltAllArray(O.None[int]())(parsers)
	}

	// Use FindFirstMap to find the first parseable string in an array
	inputs := []string{"invalid", "also bad", "42", "0x2A", "052"}
	result := A.FindFirstMap(parseNumber)(inputs)

	fmt.Printf("First parseable number: %d\n",
		F.Pipe1(result, O.GetOrElse(F.Constant(-1))))

	// Output:
	// First parseable number: 42
}

// Example_pattern_matching_numeric demonstrates pattern matching on numeric values
// with range checks and special cases.
func Example_pattern_matching_numeric() {
	// Classify numbers into categories
	isZero := func(n int) O.Option[string] {
		if n == 0 {
			return O.Some("zero")
		}
		return O.None[string]()
	}

	isNegative := func(n int) O.Option[string] {
		if n < 0 {
			return O.Some("negative")
		}
		return O.None[string]()
	}

	isSmallPositive := func(n int) O.Option[string] {
		if n > 0 && n <= 10 {
			return O.Some("small positive")
		}
		return O.None[string]()
	}

	isLargePositive := func(n int) O.Option[string] {
		if n > 10 {
			return O.Some("large positive")
		}
		return O.None[string]()
	}

	// Combine classifiers using AltAllArray
	classify := func(n int) O.Option[string] {
		classifiers := []O.Option[string]{
			isZero(n),
			isNegative(n),
			isSmallPositive(n),
			isLargePositive(n),
		}
		return O.AltAllArray(O.None[string]())(classifiers)
	}

	numbers := []int{0, -5, 3, 15, -100, 10}
	for _, n := range numbers {
		category := F.Pipe1(classify(n), O.GetOrElse(F.Constant("unknown")))
		fmt.Printf("%d: %s\n", n, category)
	}

	// Output:
	// 0: zero
	// -5: negative
	// 3: small positive
	// 15: large positive
	// -100: negative
	// 10: small positive
}

// Example_pattern_matching_with_guards demonstrates using guards (additional conditions)
// within pattern matchers for more precise matching.
func Example_pattern_matching_with_guards() {
	type Event struct {
		Type     string
		Priority int
		Message  string
	}

	// Match critical errors
	matchCriticalError := func(e Event) O.Option[string] {
		if e.Type == "error" && e.Priority >= 9 {
			return O.Some(fmt.Sprintf("CRITICAL: %s", e.Message))
		}
		return O.None[string]()
	}

	// Match regular errors
	matchError := func(e Event) O.Option[string] {
		if e.Type == "error" {
			return O.Some(fmt.Sprintf("ERROR: %s", e.Message))
		}
		return O.None[string]()
	}

	// Match warnings
	matchWarning := func(e Event) O.Option[string] {
		if e.Type == "warning" {
			return O.Some(fmt.Sprintf("WARNING: %s", e.Message))
		}
		return O.None[string]()
	}

	// Match info
	matchInfo := func(e Event) O.Option[string] {
		if e.Type == "info" {
			return O.Some(fmt.Sprintf("INFO: %s", e.Message))
		}
		return O.None[string]()
	}

	// Combine matchers - most specific first
	formatEvent := func(e Event) O.Option[string] {
		matchers := []O.Option[string]{
			matchCriticalError(e),
			matchError(e),
			matchWarning(e),
			matchInfo(e),
		}
		return O.AltAllArray(O.None[string]())(matchers)
	}

	events := []Event{
		{Type: "error", Priority: 10, Message: "System failure"},
		{Type: "error", Priority: 5, Message: "Connection lost"},
		{Type: "warning", Priority: 3, Message: "High memory usage"},
		{Type: "info", Priority: 1, Message: "User logged in"},
	}

	for _, event := range events {
		formatted := F.Pipe1(formatEvent(event), O.GetOrElse(F.Constant("UNKNOWN")))
		fmt.Println(formatted)
	}

	// Output:
	// CRITICAL: System failure
	// ERROR: Connection lost
	// WARNING: High memory usage
	// INFO: User logged in
}

// Made with Bob
