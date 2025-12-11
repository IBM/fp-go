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

// Package option defines the Option data structure and its monadic operations.
//
// Option represents an optional value: every Option is either Some and contains a value,
// or None, and does not contain a value. This is a type-safe alternative to using nil
// pointers or special sentinel values to represent the absence of a value.
//
// # Fantasy Land Specification
//
// This implementation corresponds to the Fantasy Land Maybe type:
// https://github.com/fantasyland/fantasy-land#maybe
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//   - Alt: https://github.com/fantasyland/fantasy-land#alt
//   - Plus: https://github.com/fantasyland/fantasy-land#plus
//   - Alternative: https://github.com/fantasyland/fantasy-land#alternative
//   - Foldable: https://github.com/fantasyland/fantasy-land#foldable
//
// # Basic Usage
//
// Create an Option with Some or None:
//
//	opt := Some(42)           // Option containing 42
//	opt := None[int]()        // Empty Option
//	opt := Of(42)             // Alternative to Some
//
// Check if an Option contains a value:
//
//	if IsSome(opt) {
//	    // opt contains a value
//	}
//	if IsNone(opt) {
//	    // opt is empty
//	}
//
// Extract values:
//
//	value, ok := Unwrap(opt)  // Returns (value, true) or (zero, false)
//	value := GetOrElse(func() int { return 0 })(opt)  // Returns value or default
//
// # Transformations
//
// Map transforms the contained value:
//
//	result := Map(func(x int) string {
//	    return fmt.Sprintf("%d", x)
//	})(Some(42))  // Some("42")
//
// Chain sequences operations that may fail:
//
//	result := Chain(func(x int) Option[int] {
//	    if x > 0 { return Some(x * 2) }
//	    return None[int]()
//	})(Some(5))  // Some(10)
//
// Filter keeps values that satisfy a predicate:
//
//	result := Filter(func(x int) bool {
//	    return x > 0
//	})(Some(5))  // Some(5)
//
// # Working with Collections
//
// Transform arrays:
//
//	result := TraverseArray(func(x int) Option[int] {
//	    if x > 0 { return Some(x * 2) }
//	    return None[int]()
//	})([]int{1, 2, 3})  // Some([2, 4, 6])
//
// Sequence arrays of Options:
//
//	result := SequenceArray([]Option[int]{
//	    Some(1), Some(2), Some(3),
//	})  // Some([1, 2, 3])
//
// Compact arrays (remove None values):
//
//	result := CompactArray([]Option[int]{
//	    Some(1), None[int](), Some(3),
//	})  // [1, 3]
//
// # Algebraic Operations
//
// Option supports various algebraic structures:
//
//   - Functor: Map operations
//   - Applicative: Ap operations for applying wrapped functions
//   - Monad: Chain operations for sequencing computations
//   - Eq: Equality comparison
//   - Ord: Ordering comparison
//   - Semigroup/Monoid: Combining Options
//
// # Error Handling
//
// Convert error-returning functions:
//
//	result := TryCatch(func() (int, error) {
//	    return strconv.Atoi("42")
//	})  // Some(42)
//
// Convert validation functions:
//
//	parse := FromValidation(func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	})
//	result := parse("42")  // Some(42)
//
// # Subpackages
//
//   - option/number: Number conversion utilities (Atoi, Itoa)
//   - option/testing: Testing utilities for verifying monad laws
package option

//go:generate go run .. option --count 10 --filename gen.go
