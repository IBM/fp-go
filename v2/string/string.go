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

package string

import (
	"fmt"
	"strings"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ord"
)

var (
	// ToUpperCase converts the string to uppercase
	ToUpperCase = strings.ToUpper

	// ToLowerCase converts the string to lowercase
	ToLowerCase = strings.ToLower

	// Ord implements the default ordering for strings
	Ord = ord.FromStrictCompare[string]()

	// Join joins strings
	Join = F.Curry2(F.Bind2nd[[]string, string, string])(strings.Join)

	// Equals returns a predicate that tests if a string is equal
	Equals = F.Curry2(Eq)

	// Includes returns a predicate that tests for the existence of the search string
	Includes = F.Bind2of2(strings.Contains)

	// HasPrefix returns a predicate that checks if the prefis is included in the string
	HasPrefix = F.Bind2of2(strings.HasPrefix)
)

// Eq tests if two strings are equal
func Eq(left, right string) bool {
	return left == right
}

// ToBytes converts a string to a byte slice
func ToBytes(s string) []byte {
	return []byte(s)
}

// ToRunes converts a string to a rune slice
func ToRunes(s string) []rune {
	return []rune(s)
}

// IsEmpty returns true if the string is empty
//
//go:inline
func IsEmpty(s string) bool {
	return s == ""
}

// IsNonEmpty returns true if the string is not empty
//
//go:inline
func IsNonEmpty(s string) bool {
	return s != ""
}

// Size returns the length of the string in bytes
//
//go:inline
func Size(s string) int {
	return len(s)
}

// Format applies a format string to an arbitrary value and returns a function
// that formats values of type T using the provided format string
func Format[T any](format string) func(T) string {
	return func(t T) string {
		return fmt.Sprintf(format, t)
	}
}

// Intersperse returns a function that concatenates two strings with a middle string in between.
// If either string is empty, the middle string is not added (to satisfy monoid identity laws).
func Intersperse(middle string) func(string, string) string {
	return func(l, r string) string {
		if l == "" {
			return r
		}
		if r == "" {
			return l
		}
		return l + middle + r
	}
}
