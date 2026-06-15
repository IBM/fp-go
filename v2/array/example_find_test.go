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

// Package array provides example tests demonstrating array search and filtering operations.
package array

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

// ExampleFindFirst demonstrates finding the first element in an array that satisfies a predicate.
//
// This example shows how to use FindFirst to locate the first odd number in an array.
// The predicate checks if a number is odd by testing if the bitwise AND with 2 equals 0.
//
// Returns:
//   - Option[int]: Some containing the first matching element, or None if no match found
func ExampleFindFirst() {

	pred := func(val int) bool {
		return val&2 == 0
	}

	data1 := From(1, 2, 3)

	fmt.Println(FindFirst(pred)(data1))

	// Output:
	// Some[int](1)
}

// ExampleFilter demonstrates filtering an array and retrieving the first matching element.
//
// This example shows how to compose Filter with Head to achieve similar functionality
// to FindFirst. The Filter operation creates a new array containing only elements that
// satisfy the predicate, and Head retrieves the first element from that filtered array.
//
// The composition using F.Flow2 creates a reusable pipeline that:
//  1. Filters the array to keep only odd numbers
//  2. Takes the first element from the filtered result
//
// Returns:
//   - Option[int]: Some containing the first matching element, or None if no match found
func ExampleFilter() {

	pred := func(val int) bool {
		return val&2 == 0
	}

	data1 := From(1, 2, 3)

	Find := F.Flow2(
		Filter(pred),
		Head[int],
	)

	fmt.Println(Find(data1))

	// Output:
	// Some[int](1)
}
