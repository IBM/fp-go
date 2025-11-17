// Copyright (c) 2025 IBM Corp.
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

package option

import "fmt"

type (
	Operator[A, B any] = func(A, bool) (B, bool)
	Kleisli[A, B any]  = func(A) (B, bool)
)

// IsSome checks if an Option contains a value.
//
// Parameters:
//   - t: The value of the Option
//   - tok: Whether the Option contains a value (true for Some, false for None)
//
// Example:
//
//	opt := Some(42)
//	IsSome(opt) // true
//	opt := None[int]()
//	IsSome(opt) // false
//
//go:inline
func IsSome[T any](t T, tok bool) bool {
	return tok
}

// IsNone checks if an Option is None (contains no value).
//
// Parameters:
//   - t: The value of the Option
//   - tok: Whether the Option contains a value (true for Some, false for None)
//
// Example:
//
//	opt := None[int]()
//	IsNone(opt) // true
//	opt := Some(42)
//	IsNone(opt) // false
//
//go:inline
func IsNone[T any](t T, tok bool) bool {
	return !tok
}

// Some creates an Option that contains a value.
//
// Parameters:
//   - value: The value to wrap in Some
//
// Example:
//
//	opt := Some(42) // Option containing 42
//	opt := Some("hello") // Option containing "hello"
//
//go:inline
func Some[T any](value T) (T, bool) {
	return value, true
}

// Of creates an Option that contains a value.
// This is an alias for Some and is used in monadic contexts.
//
// Parameters:
//   - value: The value to wrap in Some
//
// Example:
//
//	opt := Of(42) // Option containing 42
//
//go:inline
func Of[T any](value T) (T, bool) {
	return Some(value)
}

// None creates an Option that contains no value.
//
// Example:
//
//	opt := None[int]() // Empty Option of type int
//	opt := None[string]() // Empty Option of type string
//
//go:inline
func None[T any]() (t T, tok bool) {
	return
}

// ToString converts an Option to a string representation for debugging.
//
// Parameters:
//   - t: The value of the Option
//   - tok: Whether the Option contains a value (true for Some, false for None)
func ToString[T any](t T, tok bool) string {
	if tok {
		return fmt.Sprintf("Some[%T](%v)", t, t)
	}
	return fmt.Sprintf("None[%T]", t)
}
