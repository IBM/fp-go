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

package bytes

// Empty returns an empty byte slice.
//
// This function returns the identity element for the byte slice Monoid,
// which is an empty byte slice. It's useful as a starting point for
// building byte slices or as a default value.
//
// Returns:
//   - An empty byte slice ([]byte{})
//
// Properties:
//   - Empty() is the identity element for Monoid.Concat
//   - Monoid.Concat(Empty(), x) == x
//   - Monoid.Concat(x, Empty()) == x
//
// Example - Basic usage:
//
//	empty := Empty()
//	fmt.Println(len(empty)) // 0
//
// Example - As identity element:
//
//	data := []byte("hello")
//	result1 := Monoid.Concat(Empty(), data) // []byte("hello")
//	result2 := Monoid.Concat(data, Empty()) // []byte("hello")
//
// Example - Building byte slices:
//
//	// Start with empty and build up
//	buffer := Empty()
//	buffer = Monoid.Concat(buffer, []byte("Hello"))
//	buffer = Monoid.Concat(buffer, []byte(" "))
//	buffer = Monoid.Concat(buffer, []byte("World"))
//	// buffer: []byte("Hello World")
//
// See also:
//   - Monoid.Empty(): Alternative way to get empty byte slice
//   - ConcatAll(): For concatenating multiple byte slices
func Empty() []byte {
	return Monoid.Empty()
}

// ToString converts a byte slice to a string.
//
// This function performs a direct conversion from []byte to string.
// The conversion creates a new string with a copy of the byte data.
//
// Parameters:
//   - a: The byte slice to convert
//
// Returns:
//   - A string containing the same data as the byte slice
//
// Performance Note:
//
// This conversion allocates a new string. For performance-critical code
// that needs to avoid allocations, consider using unsafe.String (Go 1.20+)
// or working directly with byte slices.
//
// Example - Basic conversion:
//
//	bytes := []byte("hello")
//	str := ToString(bytes)
//	fmt.Println(str) // "hello"
//
// Example - Converting binary data:
//
//	// ASCII codes for "Hello"
//	data := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
//	str := ToString(data)
//	fmt.Println(str) // "Hello"
//
// Example - Empty byte slice:
//
//	empty := Empty()
//	str := ToString(empty)
//	fmt.Println(str == "") // true
//
// Example - UTF-8 encoded text:
//
//	utf8Bytes := []byte("Hello, 世界")
//	str := ToString(utf8Bytes)
//	fmt.Println(str) // "Hello, 世界"
//
// Example - Round-trip conversion:
//
//	original := "test string"
//	bytes := []byte(original)
//	result := ToString(bytes)
//	fmt.Println(original == result) // true
//
// See also:
//   - []byte(string): For converting string to byte slice
//   - Size(): For getting the length of a byte slice
func ToString(a []byte) string {
	return string(a)
}

// Size returns the number of bytes in a byte slice.
//
// This function returns the length of the byte slice, which is the number
// of bytes it contains. This is equivalent to len(as) but provided as a
// named function for use in functional composition.
//
// Parameters:
//   - as: The byte slice to measure
//
// Returns:
//   - The number of bytes in the slice
//
// Example - Basic usage:
//
//	data := []byte("hello")
//	size := Size(data)
//	fmt.Println(size) // 5
//
// Example - Empty slice:
//
//	empty := Empty()
//	size := Size(empty)
//	fmt.Println(size) // 0
//
// Example - Binary data:
//
//	binary := []byte{0x01, 0x02, 0x03, 0x04}
//	size := Size(binary)
//	fmt.Println(size) // 4
//
// Example - UTF-8 encoded text:
//
//	// Note: Size returns byte count, not character count
//	utf8 := []byte("Hello, 世界")
//	byteCount := Size(utf8)
//	fmt.Println(byteCount) // 13 (not 9 characters)
//
// Example - Using in functional composition:
//
//	import "github.com/IBM/fp-go/v2/array"
//
//	slices := [][]byte{
//	    []byte("a"),
//	    []byte("bb"),
//	    []byte("ccc"),
//	}
//
//	// Map to get sizes
//	sizes := array.Map(Size)(slices)
//	// sizes: []int{1, 2, 3}
//
// Example - Checking if slice is empty:
//
//	data := []byte("test")
//	isEmpty := Size(data) == 0
//	fmt.Println(isEmpty) // false
//
// See also:
//   - len(): Built-in function for getting slice length
//   - ToString(): For converting byte slice to string
func Size(as []byte) int {
	return len(as)
}
