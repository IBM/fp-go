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

// Package generic provides generic string utility functions that work with any type
// that has string as its underlying type (using the ~string constraint).
// This allows these functions to work with custom string types while maintaining type safety.
package generic

// ToBytes converts the string to bytes
func ToBytes[T ~string](s T) []byte {
	return []byte(s)
}

// ToRunes converts the string to runes
func ToRunes[T ~string](s T) []rune {
	return []rune(s)
}

// IsEmpty tests if the string is empty
func IsEmpty[T ~string](s T) bool {
	return len(s) == 0
}

// IsNonEmpty tests if the string is not empty
func IsNonEmpty[T ~string](s T) bool {
	return len(s) > 0
}

// Size returns the size of the string
func Size[T ~string](s T) int {
	return len(s)
}
