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

// Package bytes provides functional programming utilities for working with byte slices.
//
// This package offers algebraic structures (Monoid, Ord) and utility functions
// for byte slice operations in a functional style.
//
// # Monoid
//
// The Monoid instance for byte slices combines them through concatenation,
// with an empty byte slice as the identity element.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/bytes"
//
//	// Concatenate byte slices
//	result := bytes.Monoid.Concat([]byte("Hello"), []byte(" World"))
//	// result: []byte("Hello World")
//
//	// Identity element
//	empty := bytes.Empty() // []byte{}
//
// # ConcatAll
//
// Efficiently concatenates multiple byte slices into a single slice:
//
//	import "github.com/IBM/fp-go/v2/bytes"
//
//	result := bytes.ConcatAll(
//	    []byte("Hello"),
//	    []byte(" "),
//	    []byte("World"),
//	)
//	// result: []byte("Hello World")
//
// # Ordering
//
// The Ord instance provides lexicographic ordering for byte slices:
//
//	import "github.com/IBM/fp-go/v2/bytes"
//
//	cmp := bytes.Ord.Compare([]byte("abc"), []byte("abd"))
//	// cmp: -1 (abc < abd)
//
//	equal := bytes.Ord.Equals([]byte("test"), []byte("test"))
//	// equal: true
//
// # Utility Functions
//
// The package provides several utility functions:
//
//	// Get empty byte slice
//	empty := bytes.Empty() // []byte{}
//
//	// Convert to string
//	str := bytes.ToString([]byte("hello")) // "hello"
//
//	// Get size
//	size := bytes.Size([]byte("hello")) // 5
//
// # Use Cases
//
// The bytes package is particularly useful for:
//
//   - Building byte buffers functionally
//   - Combining multiple byte slices efficiently
//   - Comparing byte slices lexicographically
//   - Working with binary data in a functional style
//
// # Example - Building a Protocol Message
//
//	import (
//	    "github.com/IBM/fp-go/v2/bytes"
//	    "encoding/binary"
//	)
//
//	// Build a simple protocol message
//	header := []byte{0x01, 0x02}
//	length := make([]byte, 4)
//	binary.BigEndian.PutUint32(length, 100)
//	payload := []byte("data")
//
//	message := bytes.ConcatAll(header, length, payload)
//
// # Example - Sorting Byte Slices
//
//	import (
//	    "github.com/IBM/fp-go/v2/array"
//	    "github.com/IBM/fp-go/v2/bytes"
//	)
//
//	data := [][]byte{
//	    []byte("zebra"),
//	    []byte("apple"),
//	    []byte("mango"),
//	}
//
//	sorted := array.Sort(bytes.Ord)(data)
//	// sorted: [[]byte("apple"), []byte("mango"), []byte("zebra")]
package bytes
