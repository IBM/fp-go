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

package codec

import (
	"fmt"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/record"
)

// fieldKey is a comparable+Stringer key type used in examples.
type fieldKey string

func (k fieldKey) String() string { return string(k) }

// ExampleAtKey demonstrates decoding a value from a record by key.
func ExampleAtKey_decode() {
	c := AtKey[int](fieldKey("age"))

	r := record.Record[fieldKey, int]{"age": 30}
	result := c.Decode(r)

	age, _ := either.Unwrap(result)
	fmt.Println(age)

	// Output:
	// 30
}

// ExampleAtKey_encode demonstrates encoding a value into a single-entry record.
func ExampleAtKey_encode() {
	c := AtKey[string](fieldKey("city"))

	encoded := c.Encode("Berlin")
	fmt.Println(encoded["city"])

	// Output:
	// Berlin
}

// ExampleAtKey_missingKey demonstrates the failure path when the key is absent.
func ExampleAtKey_missingKey() {
	c := AtKey[int](fieldKey("score"))

	r := record.Record[fieldKey, int]{}
	result := c.Decode(r)

	fmt.Println(either.IsLeft(result))

	// Output:
	// true
}

// ExampleAtKey_name demonstrates that the codec name includes the key.
func ExampleAtKey_name() {
	c := AtKey[int](fieldKey("count"))
	fmt.Println(c.Name())

	// Output:
	// AtKey[count]
}

// ExampleAtKey_roundTrip demonstrates that encoding and then decoding is an identity.
func ExampleAtKey_roundTrip() {
	c := AtKey[int](fieldKey("n"))

	encoded := c.Encode(42)
	result := c.Decode(encoded)

	v, _ := either.Unwrap(result)
	fmt.Println(v)

	// Output:
	// 42
}
