//   Copyright (c) 2025 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// Package number provides Option-based utilities for number conversions.
package number

import (
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

func atoi(value string) (int, bool) {
	data, err := strconv.Atoi(value)
	return data, err == nil
}

var (
	// Atoi converts a string to an integer, returning Some(int) on success or None on failure.
	//
	// Example:
	//
	//	result := Atoi("42") // Some(42)
	//	result := Atoi("abc") // None
	//	result := Atoi("") // None
	Atoi = O.Optionize1(atoi)

	// Itoa converts an integer to a string, always returning Some(string).
	//
	// Example:
	//
	//	result := Itoa(42) // Some("42")
	//	result := Itoa(-10) // Some("-10")
	//	result := Itoa(0) // Some("0")
	Itoa = F.Flow2(strconv.Itoa, O.Of[string])
)
