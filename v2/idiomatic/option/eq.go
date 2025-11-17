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

import (
	EQ "github.com/IBM/fp-go/v2/eq"
)

// Eq constructs an equality predicate for Option[A] given an equality predicate for A.
// Two Options are equal if:
//   - Both are None, or
//   - Both are Some and their contained values are equal according to the provided Eq
//
// Parameters:
//   - eq: An equality predicate for the contained type A
//
// Returns a curried function that takes two Options (as tuples) and returns true if they are equal.
//
// Example:
//
//	intEq := eq.FromStrictEquals[int]()
//	optEq := Eq(intEq)
//
//	opt1 := Some(42)  // (42, true)
//	opt2 := Some(42)  // (42, true)
//	optEq(opt1)(opt2) // true
//
//	opt3 := Some(43)  // (43, true)
//	optEq(opt1)(opt3) // false
//
//	none1 := None[int]()  // (0, false)
//	none2 := None[int]()  // (0, false)
//	optEq(none1)(none2) // true
//
//	optEq(opt1)(none1) // false
func Eq[A any](eq EQ.Eq[A]) func(A, bool) func(A, bool) bool {
	return func(a1 A, a1ok bool) func(A, bool) bool {
		return func(a2 A, a2ok bool) bool {
			if a1ok {
				if a2ok {
					return eq.Equals(a1, a2)
				}
				return false
			}
			return !a2ok
		}
	}
}

// FromStrictEquals constructs an Eq for Option[A] using Go's built-in equality (==) for type A.
// This is a convenience function for comparable types.
//
// Returns a curried function that takes two Options (as tuples) and returns true if they are equal.
//
// Example:
//
//	optEq := FromStrictEquals[int]()
//
//	opt1 := Some(42)  // (42, true)
//	opt2 := Some(42)  // (42, true)
//	optEq(opt1)(opt2) // true
//
//	none1 := None[int]()  // (0, false)
//	none2 := None[int]()  // (0, false)
//	optEq(none1)(none2) // true
//
//	opt3 := Some(43)  // (43, true)
//	optEq(opt1)(opt3) // false
func FromStrictEquals[A comparable]() func(A, bool) func(A, bool) bool {
	return Eq(EQ.FromStrictEquals[A]())
}
