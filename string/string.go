// Copyright (c) 2023 IBM Corp.
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

	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/ord"
)

var (
	// ToUpperCase converts the string to uppercase
	ToUpperCase = strings.ToUpper

	// ToLowerCase converts the string to lowercase
	ToLowerCase = strings.ToLower

	// Ord implements the default ordering for strings
	Ord = ord.FromStrictCompare[string]()
)

func Eq(left string, right string) bool {
	return left == right
}

func ToBytes(s string) []byte {
	return []byte(s)
}

func ToRunes(s string) []rune {
	return []rune(s)
}

func IsEmpty(s string) bool {
	return len(s) == 0
}

func IsNonEmpty(s string) bool {
	return len(s) > 0
}

func Size(s string) int {
	return len(s)
}

// Includes returns a predicate that tests for the existence of the search string
func Includes(searchString string) func(string) bool {
	return F.Bind2nd(strings.Contains, searchString)
}

// Equals returns a predicate that tests if a string is equal
func Equals(other string) func(string) bool {
	return F.Bind2nd(Eq, other)
}

// Format applies a format string to an arbitrary value
func Format[T any](format string) func(t T) string {
	return func(t T) string {
		return fmt.Sprintf(format, t)
	}
}
