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

package iso

import (
	"strings"
	"time"

	"github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/array/nonempty"
	B "github.com/IBM/fp-go/v2/bytes"
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
)

// UTF8String creates an isomorphism between byte slices and UTF-8 strings.
// This isomorphism provides bidirectional conversion between []byte and string,
// treating the byte slice as UTF-8 encoded text.
//
// Returns:
//   - An Iso[[]byte, string] where:
//   - Get: Converts []byte to string using UTF-8 encoding
//   - ReverseGet: Converts string to []byte using UTF-8 encoding
//
// Behavior:
//   - Get direction: Interprets the byte slice as UTF-8 and returns the corresponding string
//   - ReverseGet direction: Encodes the string as UTF-8 bytes
//
// Example:
//
//	iso := UTF8String()
//
//	// Convert bytes to string
//	str := iso.Get([]byte("hello"))  // "hello"
//
//	// Convert string to bytes
//	bytes := iso.ReverseGet("world")  // []byte("world")
//
//	// Round-trip conversion
//	original := []byte("test")
//	result := iso.ReverseGet(iso.Get(original))  // []byte("test")
//
// Use cases:
//   - Converting between string and byte representations
//   - Working with APIs that use different text representations
//   - File I/O operations where you need to switch between strings and bytes
//   - Network protocols that work with byte streams
//
// Note: This isomorphism assumes valid UTF-8 encoding. Invalid UTF-8 sequences
// in the byte slice will be handled according to Go's string conversion rules
// (typically replaced with the Unicode replacement character U+FFFD).
func UTF8String() Iso[[]byte, string] {
	return MakeIso(B.ToString, S.ToBytes)
}

// lines creates an isomorphism between a slice of strings and a single string
// with lines separated by the specified separator.
// This is an internal helper function used by Lines.
//
// Parameters:
//   - sep: The separator string to use for joining/splitting lines
//
// Returns:
//   - An Iso[[]string, string] that joins/splits strings using the separator
//
// Behavior:
//   - Get direction: Joins the string slice into a single string with separators
//   - ReverseGet direction: Splits the string by the separator into a slice
func lines(sep string) Iso[[]string, string] {
	return MakeIso(S.Join(sep), F.Bind2nd(strings.Split, sep))
}

// Lines creates an isomorphism between a slice of strings and a single string
// with newline-separated lines.
// This is useful for working with multi-line text where you need to convert
// between a single string and individual lines.
//
// Returns:
//   - An Iso[[]string, string] where:
//   - Get: Joins string slice with newline characters ("\n")
//   - ReverseGet: Splits string by newline characters into a slice
//
// Behavior:
//   - Get direction: Joins each string in the slice with "\n" separator
//   - ReverseGet direction: Splits the string at each "\n" into a slice
//
// Example:
//
//	iso := Lines()
//
//	// Convert lines to single string
//	lines := []string{"line1", "line2", "line3"}
//	text := iso.Get(lines)  // "line1\nline2\nline3"
//
//	// Convert string to lines
//	text := "hello\nworld"
//	lines := iso.ReverseGet(text)  // []string{"hello", "world"}
//
//	// Round-trip conversion
//	original := []string{"a", "b", "c"}
//	result := iso.ReverseGet(iso.Get(original))  // []string{"a", "b", "c"}
//
// Use cases:
//   - Processing multi-line text files
//   - Converting between text editor representations (array of lines vs single string)
//   - Working with configuration files that have line-based structure
//   - Parsing or generating multi-line output
//
// Note: Empty strings in the slice will result in consecutive newlines in the output.
// Splitting a string with trailing newlines will include an empty string at the end.
//
// Example with edge cases:
//
//	iso := Lines()
//	lines := []string{"a", "", "b"}
//	text := iso.Get(lines)  // "a\n\nb"
//	result := iso.ReverseGet(text)  // []string{"a", "", "b"}
//
//	text := "a\nb\n"
//	lines := iso.ReverseGet(text)  // []string{"a", "b", ""}
func Lines() Iso[[]string, string] {
	return lines("\n")
}

// UnixMilli creates an isomorphism between Unix millisecond timestamps and time.Time values.
// This isomorphism provides bidirectional conversion between int64 milliseconds since
// the Unix epoch (January 1, 1970 UTC) and Go's time.Time type.
//
// Returns:
//   - An Iso[int64, time.Time] where:
//   - Get: Converts Unix milliseconds (int64) to time.Time
//   - ReverseGet: Converts time.Time to Unix milliseconds (int64)
//
// Behavior:
//   - Get direction: Creates a time.Time from milliseconds since Unix epoch
//   - ReverseGet direction: Extracts milliseconds since Unix epoch from time.Time
//
// Example:
//
//	iso := UnixMilli()
//
//	// Convert milliseconds to time.Time
//	millis := int64(1609459200000)  // 2021-01-01 00:00:00 UTC
//	t := iso.Get(millis)
//
//	// Convert time.Time to milliseconds
//	now := time.Now()
//	millis := iso.ReverseGet(now)
//
//	// Round-trip conversion
//	original := int64(1234567890000)
//	result := iso.ReverseGet(iso.Get(original))  // 1234567890000
//
// Use cases:
//   - Working with APIs that use Unix millisecond timestamps (e.g., JavaScript Date.now())
//   - Database storage where timestamps are stored as integers
//   - JSON serialization/deserialization of timestamps
//   - Converting between different time representations in distributed systems
//
// Precision notes:
//   - Millisecond precision is maintained in both directions
//   - Sub-millisecond precision in time.Time is lost when converting to int64
//   - The conversion is timezone-aware (time.Time includes location information)
//
// Example with precision:
//
//	iso := UnixMilli()
//	t := time.Date(2021, 1, 1, 12, 30, 45, 123456789, time.UTC)
//	millis := iso.ReverseGet(t)  // Nanoseconds are truncated to milliseconds
//	restored := iso.Get(millis)   // Nanoseconds will be 123000000
//
// Note: This isomorphism uses UTC for the time.Time values. If you need to preserve
// timezone information, consider storing it separately or using a different representation.
func UnixMilli() Iso[int64, time.Time] {
	return MakeIso(time.UnixMilli, time.Time.UnixMilli)
}

// Add creates an isomorphism that adds a constant value to a number.
// This isomorphism provides bidirectional conversion by adding a value in one direction
// and subtracting it in the reverse direction, effectively shifting numbers by a constant offset.
//
// Type Parameters:
//   - T: A numeric type (integer, float, or complex number)
//
// Parameters:
//   - n: The constant value to add in the Get direction (and subtract in ReverseGet)
//
// Returns:
//   - An Iso[T, T] where:
//   - Get: Adds n to the input value (x + n)
//   - ReverseGet: Subtracts n from the input value (x - n)
//
// Behavior:
//   - Get direction: Shifts the value up by n
//   - ReverseGet direction: Shifts the value down by n (inverse operation)
//
// Example:
//
//	// Create an isomorphism that adds 10
//	addTen := Add(10)
//
//	// Add 10 to a value
//	result := addTen.Get(5)  // 15
//
//	// Subtract 10 from a value (reverse)
//	original := addTen.ReverseGet(15)  // 5
//
//	// Round-trip conversion
//	value := 42
//	roundTrip := addTen.ReverseGet(addTen.Get(value))  // 42
//
// Use cases:
//   - Converting between different numeric scales or offsets
//   - Temperature conversions (e.g., Celsius offset adjustments)
//   - Index transformations (e.g., 0-based to 1-based indexing)
//   - Coordinate system translations
//   - Time zone offset adjustments (when working with numeric timestamps)
//
// Example with different numeric types:
//
//	// Integer addition
//	intIso := Add(5)
//	intResult := intIso.Get(10)  // 15
//
//	// Float addition
//	floatIso := Add(2.5)
//	floatResult := floatIso.Get(7.5)  // 10.0
//
//	// Complex number addition
//	complexIso := Add(complex(1, 2))
//	complexResult := complexIso.Get(complex(3, 4))  // (4+6i)
//
// Example with coordinate translation:
//
//	// Translate x-coordinates by 100 units
//	translateX := Add(100)
//	newX := translateX.Get(50)  // 150
//	originalX := translateX.ReverseGet(150)  // 50
//
// Note: This isomorphism satisfies the round-trip laws:
//   - ReverseGet(Get(x)) == x (because (x + n) - n == x)
//   - Get(ReverseGet(x)) == x (because (x - n) + n == x)
func Add[T Number](n T) Iso[T, T] {
	return MakeIso(
		N.Add(n),
		N.Sub(n),
	)
}

// Sub creates an isomorphism that subtracts a constant value from a number.
// This isomorphism provides bidirectional conversion by subtracting a value in one direction
// and adding it in the reverse direction, effectively shifting numbers by a constant offset in the opposite direction of Add.
//
// Type Parameters:
//   - T: A numeric type (integer, float, or complex number)
//
// Parameters:
//   - n: The constant value to subtract in the Get direction (and add in ReverseGet)
//
// Returns:
//   - An Iso[T, T] where:
//   - Get: Subtracts n from the input value (x - n)
//   - ReverseGet: Adds n to the input value (x + n)
//
// Behavior:
//   - Get direction: Shifts the value down by n
//   - ReverseGet direction: Shifts the value up by n (inverse operation)
//
// Example:
//
//	// Create an isomorphism that subtracts 10
//	subTen := Sub(10)
//
//	// Subtract 10 from a value
//	result := subTen.Get(15)  // 5
//
//	// Add 10 to a value (reverse)
//	original := subTen.ReverseGet(5)  // 15
//
//	// Round-trip conversion
//	value := 42
//	roundTrip := subTen.ReverseGet(subTen.Get(value))  // 42
//
// Use cases:
//   - Converting between different numeric scales or offsets
//   - Discount or reduction calculations
//   - Index transformations (e.g., 1-based to 0-based indexing)
//   - Coordinate system translations in the negative direction
//   - Time calculations (e.g., going back in time by a fixed amount)
//
// Example with different numeric types:
//
//	// Integer subtraction
//	intIso := Sub(5)
//	intResult := intIso.Get(10)  // 5
//
//	// Float subtraction
//	floatIso := Sub(2.5)
//	floatResult := floatIso.Get(10.0)  // 7.5
//
//	// Complex number subtraction
//	complexIso := Sub(complex(1, 2))
//	complexResult := complexIso.Get(complex(4, 6))  // (3+4i)
//
// Example with discount calculation:
//
//	// Apply a discount of 10 units
//	discount := Sub(10)
//	discountedPrice := discount.Get(50)  // 40
//	originalPrice := discount.ReverseGet(40)  // 50
//
// Relationship with Add:
// Sub(n) is equivalent to Add(-n). The following are equivalent:
//
//	sub5 := Sub(5)
//	addNeg5 := Add(-5)
//	// Both produce the same results:
//	sub5.Get(10)      // 5
//	addNeg5.Get(10)   // 5
//
// Note: This isomorphism satisfies the round-trip laws:
//   - ReverseGet(Get(x)) == x (because (x - n) + n == x)
//   - Get(ReverseGet(x)) == x (because (x + n) - n == x)
func Sub[T Number](n T) Iso[T, T] {
	return MakeIso(
		N.Sub(n),
		N.Add(n),
	)
}

// SwapPair creates an isomorphism that swaps the elements of a Pair.
// This isomorphism provides bidirectional conversion between Pair[A, B] and Pair[B, A],
// effectively exchanging the first and second elements of the pair.
//
// Type Parameters:
//   - A: The type of the first element in the source pair (becomes second in target)
//   - B: The type of the second element in the source pair (becomes first in target)
//
// Returns:
//   - An Iso[Pair[A, B], Pair[B, A]] where:
//   - Get: Swaps the pair elements from (A, B) to (B, A)
//   - ReverseGet: Swaps the pair elements from (B, A) back to (A, B)
//
// Behavior:
//   - Get direction: Transforms Pair[A, B] to Pair[B, A]
//   - ReverseGet direction: Transforms Pair[B, A] to Pair[A, B] (inverse operation)
//
// Example:
//
//	// Create a swap isomorphism for pairs of string and int
//	swapIso := SwapPair[string, int]()
//
//	// Swap a pair
//	original := pair.MakePair("hello", 42)  // Pair[string, int]
//	swapped := swapIso.Get(original)        // Pair[int, string] = (42, "hello")
//
//	// Swap back
//	restored := swapIso.ReverseGet(swapped) // Pair[string, int] = ("hello", 42)
//
//	// Round-trip conversion
//	p := pair.MakePair(1, "a")
//	roundTrip := swapIso.ReverseGet(swapIso.Get(p))  // (1, "a")
//
// Use cases:
//   - Reordering pair elements for API compatibility
//   - Converting between different pair representations
//   - Normalizing data structures with swapped element order
//   - Working with functions that expect arguments in different order
//   - Adapting between coordinate systems (e.g., (x, y) to (y, x))
//
// Example with coordinates:
//
//	// Swap between (x, y) and (y, x) coordinates
//	swapCoords := SwapPair[float64, float64]()
//
//	point := pair.MakePair(3.0, 4.0)  // (x=3, y=4)
//	swapped := swapCoords.Get(point)   // (y=4, x=3)
//
// Example with heterogeneous types:
//
//	// Swap between (name, age) and (age, name)
//	swapPerson := SwapPair[string, int]()
//
//	person := pair.MakePair("Alice", 30)
//	swapped := swapPerson.Get(person)  // (30, "Alice")
//
// Example with function composition:
//
//	// Use with other isomorphisms
//	swapIso := SwapPair[int, string]()
//	p := pair.MakePair(1, "test")
//
//	// Apply swap twice to get back to original
//	result := F.Pipe2(p, swapIso.Get, swapIso.ReverseGet)  // (1, "test")
//
// Note: This isomorphism satisfies the round-trip laws:
//   - ReverseGet(Get(pair)) == pair (swapping twice returns to original)
//   - Get(ReverseGet(pair)) == pair (swapping twice returns to original)
//
// Note: SwapPair is self-inverse, meaning applying it twice returns the original value.
// This makes it particularly useful for symmetric transformations.
func SwapPair[A, B any]() Iso[Pair[A, B], Pair[B, A]] {
	return MakeIso(
		pair.Swap[A, B],
		pair.Swap[B, A],
	)
}

// SwapEither creates an isomorphism that swaps the type parameters of an Either.
// This isomorphism provides bidirectional conversion between Either[E, A] and Either[A, E],
// effectively exchanging the Left and Right type positions while preserving which side
// the value is on.
//
// Type Parameters:
//   - E: The type of the Left value in the source Either (becomes Right in target)
//   - A: The type of the Right value in the source Either (becomes Left in target)
//
// Returns:
//   - An Iso[Either[E, A], Either[A, E]] where:
//   - Get: Swaps Either[E, A] to Either[A, E] (Left[E] becomes Right[E], Right[A] becomes Left[A])
//   - ReverseGet: Swaps Either[A, E] back to Either[E, A] (inverse operation)
//
// Behavior:
//   - Get direction: Transforms Either[E, A] to Either[A, E]
//   - Left[E] becomes Right[E]
//   - Right[A] becomes Left[A]
//   - ReverseGet direction: Transforms Either[A, E] to Either[E, A] (inverse)
//   - Right[E] becomes Left[E]
//   - Left[A] becomes Right[A]
//
// Example:
//
//	// Create a swap isomorphism for Either[string, int]
//	swapIso := SwapEither[string, int]()
//
//	// Swap a Left value
//	leftVal := either.Left[int]("error")     // Either[string, int] with Left
//	swapped := swapIso.Get(leftVal)          // Either[int, string] with Right("error")
//
//	// Swap a Right value
//	rightVal := either.Right[string](42)     // Either[string, int] with Right
//	swapped2 := swapIso.Get(rightVal)        // Either[int, string] with Left(42)
//
//	// Round-trip conversion
//	original := either.Left[int]("test")
//	roundTrip := swapIso.ReverseGet(swapIso.Get(original))  // Left("test")
//
// Use cases:
//   - Converting between different Either conventions (error-left vs error-right)
//   - Adapting between APIs with different Either type parameter orders
//   - Normalizing error handling patterns
//   - Working with libraries that use opposite Either conventions
//   - Transforming result types to match expected signatures
//
// Example with error handling:
//
//	// Convert from Either[Error, Value] to Either[Value, Error]
//	swapError := SwapEither[error, string]()
//
//	result := either.Right[error]("success")  // Either[error, string]
//	swapped := swapError.Get(result)          // Either[string, error] with Left("success")
//
// Example with validation:
//
//	// Swap validation result types
//	swapValidation := SwapEither[[]string, User]()
//
//	// Valid user
//	valid := either.Right[[]string](User{Name: "Alice"})
//	swapped := swapValidation.Get(valid)  // Either[User, []string] with Left(User)
//
//	// Invalid user
//	invalid := either.Left[User]([]string{"error1", "error2"})
//	swapped2 := swapValidation.Get(invalid)  // Either[User, []string] with Right(errors)
//
// Example demonstrating self-inverse property:
//
//	swapIso := SwapEither[string, int]()
//	value := either.Left[int]("error")
//
//	// Apply swap twice to get back to original
//	result := F.Pipe2(value, swapIso.Get, swapIso.ReverseGet)  // Left("error")
//
// Note: This isomorphism satisfies the round-trip laws:
//   - ReverseGet(Get(either)) == either (swapping twice returns to original)
//   - Get(ReverseGet(either)) == either (swapping twice returns to original)
//
// Note: SwapEither is self-inverse, meaning applying it twice returns the original value.
// The swap operation preserves which side (Left/Right) the value is on, only changing
// the type parameter positions.
func SwapEither[E, A any]() Iso[Either[E, A], Either[A, E]] {
	return MakeIso(
		either.Swap[E, A],
		either.Swap[A, E],
	)
}

// ReverseArray creates an isomorphism that reverses the order of elements in a slice.
// This isomorphism is self-inverse, meaning applying it twice returns the original slice.
// It provides bidirectional conversion where both Get and ReverseGet perform the same
// reversal operation.
//
// Type Parameters:
//   - A: The type of elements in the slice
//
// Returns:
//   - An Iso[[]A, []A] where:
//   - Get: Reverses the slice order
//   - ReverseGet: Reverses the slice order (same as Get, since reverse is self-inverse)
//
// Behavior:
//   - Get direction: Returns a new slice with elements in reverse order
//   - ReverseGet direction: Returns a new slice with elements in reverse order
//   - Both directions create new slices without modifying the original
//   - Self-inverse property: Get(Get(x)) == x and ReverseGet(ReverseGet(x)) == x
//
// Example:
//
//	iso := ReverseArray[int]()
//
//	// Reverse a slice
//	numbers := []int{1, 2, 3, 4, 5}
//	reversed := iso.Get(numbers)  // []int{5, 4, 3, 2, 1}
//
//	// Reverse back (using ReverseGet)
//	original := iso.ReverseGet(reversed)  // []int{1, 2, 3, 4, 5}
//
//	// Round-trip conversion
//	data := []int{10, 20, 30}
//	roundTrip := iso.ReverseGet(iso.Get(data))  // []int{10, 20, 30}
//
// Use cases:
//   - Processing data in reverse order within an isomorphism pipeline
//   - Implementing reversible transformations on collections
//   - Converting between forward and backward iteration orders
//   - Working with optics that need to reverse array order
//   - Composing with other isomorphisms for complex transformations
//
// Example with strings:
//
//	iso := ReverseArray[string]()
//	words := []string{"hello", "world", "foo"}
//	reversed := iso.Get(words)  // []string{"foo", "world", "hello"}
//
// Example with composition:
//
//	// Reverse, then map, then reverse back
//	iso := ReverseArray[int]()
//	numbers := []int{1, 2, 3, 4, 5}
//
//	result := F.Pipe3(
//	    numbers,
//	    iso.Get,                                    // Reverse
//	    array.Map(N.Mul(2)), // Map
//	    iso.ReverseGet,                             // Reverse back
//	)
//	// result: []int{2, 4, 6, 8, 10}
//
// Example with lens composition:
//
//	// Use ReverseArray as part of a lens pipeline
//	type Container struct {
//	    Items []int
//	}
//
//	// Create an isomorphism that reverses items in a container
//	reverseItems := ReverseArray[int]()
//
// Note: This isomorphism is self-inverse, which means:
//   - Get and ReverseGet perform the same operation
//   - Applying the isomorphism twice returns the original value
//   - ReverseArray().Get == ReverseArray().ReverseGet
//
// Performance:
//   - Time complexity: O(n) where n is the length of the slice
//   - Space complexity: O(n) for the new slice
//   - Both Get and ReverseGet have the same performance characteristics
func ReverseArray[A any]() Iso[[]A, []A] {
	return MakeIso(
		array.Reverse[A],
		array.Reverse[A],
	)
}

// Head creates an isomorphism between a single element and a non-empty array containing that element.
// This isomorphism provides bidirectional conversion between a value and a singleton non-empty array,
// where the value becomes the head (first element) of the array.
//
// Type Parameters:
//   - A: The type of the element
//
// Returns:
//   - An Iso[A, NonEmptyArray[A]] where:
//   - Get: Wraps a single element into a non-empty array (singleton array)
//   - ReverseGet: Extracts the head (first element) from a non-empty array
//
// Behavior:
//   - Get direction: Creates a non-empty array with the single element as its head
//   - ReverseGet direction: Extracts the first element from a non-empty array
//   - The non-empty array type guarantees at least one element exists
//
// Example:
//
//	iso := Head[int]()
//
//	// Wrap a value into a non-empty array
//	value := 42
//	arr := iso.Get(value)  // NonEmptyArray[int] with head=42
//
//	// Extract the head from a non-empty array
//	head := iso.ReverseGet(arr)  // 42
//
//	// Round-trip conversion
//	original := 100
//	roundTrip := iso.ReverseGet(iso.Get(original))  // 100
//
// Use cases:
//   - Converting between single values and non-empty collections
//   - Working with APIs that require non-empty arrays
//   - Ensuring type safety when a value must be in a non-empty context
//   - Composing with other optics that work on non-empty arrays
//   - Lifting single values into collection contexts
//
// Example with strings:
//
//	iso := Head[string]()
//	name := "Alice"
//	arr := iso.Get(name)  // NonEmptyArray[string] with head="Alice"
//	extracted := iso.ReverseGet(arr)  // "Alice"
//
// Example with structs:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	iso := Head[User]()
//	user := User{"Bob", 30}
//	arr := iso.Get(user)  // NonEmptyArray[User] with head=User{"Bob", 30}
//
// Example with composition:
//
//	// Lift a value into a non-empty array, then process it
//	iso := Head[int]()
//	value := 5
//
//	result := F.Pipe2(
//	    value,
//	    iso.Get,                                    // Wrap in non-empty array
//	    nonempty.Map(N.Mul(2)), // Map over array
//	)
//	// result: NonEmptyArray[int] with head=10
//
// Example with lens usage:
//
//	// Use Head to focus on a single value as a non-empty array
//	type Config struct {
//	    DefaultValue int
//	}
//
//	headIso := Head[int]()
//	config := Config{DefaultValue: 42}
//
//	// Convert default value to non-empty array for processing
//	arr := headIso.Get(config.DefaultValue)
//
// Note: This isomorphism satisfies the round-trip laws:
//   - ReverseGet(Get(x)) == x (extracting head from singleton returns original)
//   - Get(ReverseGet(arr)) creates a singleton with the same head
//
// Important: When using ReverseGet on a non-empty array with multiple elements,
// only the head (first element) is extracted. Other elements are discarded.
//
// Example with multi-element array:
//
//	iso := Head[int]()
//	// If you have a non-empty array with multiple elements
//	arr := nonempty.From(1, 2, 3, 4, 5)
//	head := iso.ReverseGet(arr)  // 1 (only the head is extracted)
func Head[A any]() Iso[A, NonEmptyArray[A]] {
	return MakeIso(
		nonempty.Of[A],
		nonempty.Head[A],
	)
}
