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

package result

import (
	EQ "github.com/IBM/fp-go/v2/eq"
)

// Eq constructs an equality predicate for Result values (A, error).
// Two Result values are equal if they are both Left (error) with equal error values,
// or both Right (success) with equal values according to the provided equality predicate.
//
// Parameters:
//   - eq: Equality predicate for the Right (success) type A
//
// Returns a curried comparison function that takes two Result values and returns true if equal.
//
// Example:
//
//	eq := result.Eq(eq.FromStrictEquals[int]())
//	result1 := eq(42, nil)(42, nil) // true
//	result2 := eq(42, nil)(43, nil) // false
func Eq[A any](eq EQ.Eq[A]) func(A, error) func(A, error) bool {
	return func(a A, aerr error) func(A, error) bool {
		return func(b A, berr error) bool {
			if aerr != nil {
				if berr != nil {
					return aerr == berr
				}
				return false
			}
			if berr != nil {
				return false
			}
			return eq.Equals(a, b)
		}
	}
}

// FromStrictEquals constructs an equality predicate using Go's == operator.
// The Right type must be comparable.
//
// Example:
//
//	eq := result.FromStrictEquals[int]()
//	result1 := eq(42, nil)(42, nil) // true
//	result2 := eq(42, nil)(43, nil) // false
func FromStrictEquals[A comparable]() func(A, error) func(A, error) bool {
	return Eq(EQ.FromStrictEquals[A]())
}
