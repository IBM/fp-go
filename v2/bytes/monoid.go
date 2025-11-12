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

import (
	"bytes"

	A "github.com/IBM/fp-go/v2/array"
	O "github.com/IBM/fp-go/v2/ord"
)

var (
	// Monoid is the Monoid instance for byte slices.
	//
	// This Monoid combines byte slices through concatenation, with an empty
	// byte slice as the identity element. It satisfies the monoid laws:
	//
	// Identity laws:
	//   - Monoid.Concat(Monoid.Empty(), x) == x (left identity)
	//   - Monoid.Concat(x, Monoid.Empty()) == x (right identity)
	//
	// Associativity law:
	//   - Monoid.Concat(Monoid.Concat(a, b), c) == Monoid.Concat(a, Monoid.Concat(b, c))
	//
	// Operations:
	//   - Empty(): Returns an empty byte slice []byte{}
	//   - Concat(a, b []byte): Concatenates two byte slices
	//
	// Example - Basic concatenation:
	//
	//	result := Monoid.Concat([]byte("Hello"), []byte(" World"))
	//	// result: []byte("Hello World")
	//
	// Example - Identity element:
	//
	//	empty := Monoid.Empty()
	//	data := []byte("test")
	//	result1 := Monoid.Concat(empty, data) // []byte("test")
	//	result2 := Monoid.Concat(data, empty) // []byte("test")
	//
	// Example - Building byte buffers:
	//
	//	buffer := Monoid.Empty()
	//	buffer = Monoid.Concat(buffer, []byte("Line 1\n"))
	//	buffer = Monoid.Concat(buffer, []byte("Line 2\n"))
	//	buffer = Monoid.Concat(buffer, []byte("Line 3\n"))
	//
	// Example - Associativity:
	//
	//	a := []byte("a")
	//	b := []byte("b")
	//	c := []byte("c")
	//	left := Monoid.Concat(Monoid.Concat(a, b), c)   // []byte("abc")
	//	right := Monoid.Concat(a, Monoid.Concat(b, c))  // []byte("abc")
	//	// left == right
	//
	// See also:
	//   - ConcatAll: For concatenating multiple byte slices at once
	//   - Empty(): Convenience function for getting empty byte slice
	Monoid = A.Monoid[byte]()

	// ConcatAll efficiently concatenates multiple byte slices into a single slice.
	//
	// This function takes a variadic number of byte slices and combines them
	// into a single byte slice. It pre-allocates the exact amount of memory
	// needed, making it more efficient than repeated concatenation.
	//
	// Parameters:
	//   - slices: Zero or more byte slices to concatenate
	//
	// Returns:
	//   - A new byte slice containing all input slices concatenated in order
	//
	// Performance:
	//
	// ConcatAll is more efficient than using Monoid.Concat repeatedly because
	// it calculates the total size upfront and allocates memory once, avoiding
	// multiple allocations and copies.
	//
	// Example - Basic usage:
	//
	//	result := ConcatAll(
	//	    []byte("Hello"),
	//	    []byte(" "),
	//	    []byte("World"),
	//	)
	//	// result: []byte("Hello World")
	//
	// Example - Empty input:
	//
	//	result := ConcatAll()
	//	// result: []byte{}
	//
	// Example - Single slice:
	//
	//	result := ConcatAll([]byte("test"))
	//	// result: []byte("test")
	//
	// Example - Building protocol messages:
	//
	//	import "encoding/binary"
	//
	//	header := []byte{0x01, 0x02}
	//	length := make([]byte, 4)
	//	binary.BigEndian.PutUint32(length, 100)
	//	payload := []byte("data")
	//	footer := []byte{0xFF}
	//
	//	message := ConcatAll(header, length, payload, footer)
	//
	// Example - With empty slices:
	//
	//	result := ConcatAll(
	//	    []byte("a"),
	//	    []byte{},
	//	    []byte("b"),
	//	    []byte{},
	//	    []byte("c"),
	//	)
	//	// result: []byte("abc")
	//
	// Example - Building CSV line:
	//
	//	fields := [][]byte{
	//	    []byte("John"),
	//	    []byte("Doe"),
	//	    []byte("30"),
	//	}
	//	separator := []byte(",")
	//
	//	// Interleave fields with separators
	//	parts := [][]byte{
	//	    fields[0], separator,
	//	    fields[1], separator,
	//	    fields[2],
	//	}
	//	line := ConcatAll(parts...)
	//	// line: []byte("John,Doe,30")
	//
	// See also:
	//   - Monoid.Concat: For concatenating exactly two byte slices
	//   - bytes.Join: Standard library function for joining with separator
	ConcatAll = A.ArrayConcatAll[byte]

	// Ord is the Ord instance for byte slices providing lexicographic ordering.
	//
	// This Ord instance compares byte slices lexicographically (dictionary order),
	// comparing bytes from left to right until a difference is found or one slice
	// ends. It uses the standard library's bytes.Compare and bytes.Equal functions.
	//
	// Comparison rules:
	//   - Compares byte-by-byte from left to right
	//   - First differing byte determines the order
	//   - Shorter slice is less than longer slice if all bytes match
	//   - Empty slice is less than any non-empty slice
	//
	// Operations:
	//   - Compare(a, b []byte) int: Returns -1 if a < b, 0 if a == b, 1 if a > b
	//   - Equals(a, b []byte) bool: Returns true if slices are equal
	//
	// Example - Basic comparison:
	//
	//	cmp := Ord.Compare([]byte("abc"), []byte("abd"))
	//	// cmp: -1 (abc < abd)
	//
	//	cmp = Ord.Compare([]byte("xyz"), []byte("abc"))
	//	// cmp: 1 (xyz > abc)
	//
	//	cmp = Ord.Compare([]byte("test"), []byte("test"))
	//	// cmp: 0 (equal)
	//
	// Example - Length differences:
	//
	//	cmp := Ord.Compare([]byte("ab"), []byte("abc"))
	//	// cmp: -1 (shorter is less)
	//
	//	cmp = Ord.Compare([]byte("abc"), []byte("ab"))
	//	// cmp: 1 (longer is greater)
	//
	// Example - Empty slices:
	//
	//	cmp := Ord.Compare([]byte{}, []byte("a"))
	//	// cmp: -1 (empty is less)
	//
	//	cmp = Ord.Compare([]byte{}, []byte{})
	//	// cmp: 0 (both empty)
	//
	// Example - Equality check:
	//
	//	equal := Ord.Equals([]byte("test"), []byte("test"))
	//	// equal: true
	//
	//	equal = Ord.Equals([]byte("test"), []byte("Test"))
	//	// equal: false (case-sensitive)
	//
	// Example - Sorting byte slices:
	//
	//	import "github.com/IBM/fp-go/v2/array"
	//
	//	data := [][]byte{
	//	    []byte("zebra"),
	//	    []byte("apple"),
	//	    []byte("mango"),
	//	}
	//
	//	sorted := array.Sort(Ord)(data)
	//	// sorted: [[]byte("apple"), []byte("mango"), []byte("zebra")]
	//
	// Example - Binary data comparison:
	//
	//	cmp := Ord.Compare([]byte{0x01, 0x02}, []byte{0x01, 0x03})
	//	// cmp: -1 (0x02 < 0x03)
	//
	// Example - Finding minimum:
	//
	//	import O "github.com/IBM/fp-go/v2/ord"
	//
	//	a := []byte("xyz")
	//	b := []byte("abc")
	//	min := O.Min(Ord)(a, b)
	//	// min: []byte("abc")
	//
	// See also:
	//   - bytes.Compare: Standard library comparison function
	//   - bytes.Equal: Standard library equality function
	//   - array.Sort: For sorting slices using an Ord instance
	Ord = O.MakeOrd(bytes.Compare, bytes.Equal)
)
