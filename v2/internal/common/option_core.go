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

package common

import (
	"bytes"
	"encoding/json"
	"reflect"
)

var (
	// jsonNull is the cached representation of the `null` serialization in JSON
	jsonNull = []byte("null")
)

// Option defines a data structure that logically holds a value or not.
// It represents an optional value: every Option is either Some and contains a value,
// or None, and does not contain a value.
//
// Option is commonly used to represent the result of operations that may fail,
// as an alternative to returning nil pointers or using error values.
//
// Example:
//
//	var opt Option[int] = Some(42)  // Contains a value
//	var opt Option[int] = None[int]() // Contains no value
type Option[A any] struct {
	a A
	s bool
}

type (
	OptionKleisli[A, B any]  = func(A) Option[B]
	OptionOperator[A, B any] = OptionKleisli[Option[A], B]
	OptionKleisliI[A, B any] = func(A) (B, bool)

	// OptionTraversable represents a data structure that can be traversed from left to right,
	// applying an effectful function to each element and collecting the results.
	//
	// A Traversable takes a Kleisli arrow (a function that returns an Option) and
	// produces another Kleisli arrow that operates on a container of values.
	//
	// Type Parameters:
	//   - A: The input element type
	//   - B: The output element type
	//   - GA: The input container type (e.g., []A, map[K]A)
	//   - GB: The output container type (e.g., []B, map[K]B)
	//
	// The Traversable signature:
	//
	//	func(Kleisli[A, B]) Kleisli[GA, GB]
	//
	// expands to:
	//
	//	func(func(A) Option[B]) func(GA) Option[GB]
	//
	// This means: given a function that transforms A to Option[B], produce a function
	// that transforms a container of A values into an Option of a container of B values.
	//
	// Behavior:
	//   - If all transformations succeed (return Some), the result is Some containing
	//     the container of all transformed values
	//   - If any transformation fails (returns None), the entire result is None
	//
	// Common Use Cases:
	//   - Validating and transforming collections where any failure should fail the whole operation
	//   - Parsing collections of strings where all must parse successfully
	//   - Applying optional transformations across data structures
	//
	// Example:
	//
	//	// Array traversable
	//	traversable := TraversableArray[string, int]()
	//	parse := func(s string) Option[int] {
	//	    n, err := strconv.Atoi(s)
	//	    if err != nil { return None[int]() }
	//	    return Some(n)
	//	}
	//	result := traversable(parse)([]string{"1", "2", "3"}) // Some([1, 2, 3])
	//	result := traversable(parse)([]string{"1", "x", "3"}) // None
	//
	// See Also:
	//   - TraversableArray: Traversable instance for arrays
	//   - TraverseArray: Direct array traversal function
	//   - TraverseRecord: Traversal for maps/records
	OptionTraversable[A, B, GA, GB any] = func(OptionKleisli[A, B]) OptionKleisli[GA, GB]
)

// String implements fmt.Stringer for Option.
// Returns a human-readable string representation.
//
// Example:
//
//	Some(42).String() // "Some[int](42)"
//	None[int]().String() // "None[int]"
func (s Option[A]) String() string {
	return optString(s.s, s.a)
}

func optMarshalJSON(isSome bool, value any) ([]byte, error) {
	if isSome {
		return json.Marshal(value)
	}
	return jsonNull, nil
}

func (s Option[A]) MarshalJSON() ([]byte, error) {
	return optMarshalJSON(s.s, s.a)
}

// optUnmarshalJSON unmarshals the [Option] from a JSON string
//

func optUnmarshalJSON(isSome *bool, value any, data []byte) error {
	// decode the value
	if bytes.Equal(data, jsonNull) {
		*isSome = false
		reflect.ValueOf(value).Elem().SetZero()
		return nil
	}
	*isSome = true
	return json.Unmarshal(data, value)
}

func (s *Option[A]) UnmarshalJSON(data []byte) error {
	return optUnmarshalJSON(&s.s, &s.a, data)
}

// OptionIsNone checks if an Option is None (contains no value).
//
// Example:
//
//	opt := OptionNone[int]()
//	OptionIsNone(opt) // true
//	opt := OptionSome(42)
//	OptionIsNone(opt) // false
//
//go:inline
func OptionIsNone[T any](val Option[T]) bool {
	return !val.s
}

// OptionSome creates an Option that contains a value.
//
// Example:
//
//	opt := OptionSome(42) // Option containing 42
//	opt := OptionSome("hello") // Option containing "hello"
//
//go:inline
func OptionSome[T any](value T) Option[T] {
	return Option[T]{s: true, a: value}
}

// OptionOf creates an Option that contains a value.
// This is an alias for OptionSome and is used in monadic contexts.
//
// Example:
//
//	opt := OptionOf(42) // Option containing 42
//
//go:inline
func OptionOf[T any](value T) Option[T] {
	return OptionSome(value)
}

// OptionNone creates an Option that contains no value.
//
// Example:
//
//	opt := OptionNone[int]() // Empty Option of type int
//	opt := OptionNone[string]() // Empty Option of type string
//
//go:inline
func OptionNone[T any]() Option[T] {
	return Option[T]{}
}

// OptionIsSome checks if an Option contains a value.
//
// Example:
//
//	opt := OptionSome(42)
//	OptionIsSome(opt) // true
//	opt := OptionNone[int]()
//	OptionIsSome(opt) // false
//
//go:inline
func OptionIsSome[T any](val Option[T]) bool {
	return val.s
}

// OptionMonadFold performs a fold operation on an Option.
// If the Option is Some, applies onSome to the value.
// If the Option is None, calls onNone.
//
// Example:
//
//	opt := OptionSome(42)
//	result := OptionMonadFold(opt,
//	    func() string { return "no value" },
//	    func(x int) string { return fmt.Sprintf("value: %d", x) },
//	) // "value: 42"
func OptionMonadFold[A, B any](ma Option[A], onNone func() B, onSome func(A) B) B {
	if OptionIsSome(ma) {
		return onSome(ma.a)
	}
	return onNone()
}

// OptionUnwrap extracts the value and presence flag from an Option.
// Returns the value and true if Some, or zero value and false if None.
//
// Example:
//
//	opt := OptionSome(42)
//	val, ok := OptionUnwrap(opt) // val = 42, ok = true
//	opt := OptionNone[int]()
//	val, ok := OptionUnwrap(opt) // val = 0, ok = false
//
//go:inline
func OptionUnwrap[A any](ma Option[A]) (A, bool) {
	return ma.a, ma.s
}
