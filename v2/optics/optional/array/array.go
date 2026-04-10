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

package array

import (
	OP "github.com/IBM/fp-go/v2/optics/optional"
	G "github.com/IBM/fp-go/v2/optics/optional/array/generic"
)

// At creates an Optional that focuses on the element at a specific index in an array.
//
// This function returns an Optional that can get and set the element at the given index.
// If the index is out of bounds, GetOption returns None and Set operations are no-ops
// (the array is returned unchanged). This follows the Optional laws where operations
// on non-existent values have no effect.
//
// The Optional provides safe array access without panicking on invalid indices, making
// it ideal for functional transformations where you want to modify array elements only
// when they exist.
//
// Type Parameters:
//   - A: The type of elements in the array
//
// Parameters:
//   - idx: The zero-based index to focus on
//
// Returns:
//   - An Optional that focuses on the element at the specified index
//
// Example:
//
//	import (
//	    AR "github.com/IBM/fp-go/v2/array"
//	    OP "github.com/IBM/fp-go/v2/optics/optional"
//	    OA "github.com/IBM/fp-go/v2/optics/optional/array"
//	)
//
//	numbers := []int{10, 20, 30, 40}
//
//	// Create an optional focusing on index 1
//	second := OA.At[int](1)
//
//	// Get the element at index 1
//	value := second.GetOption(numbers)
//	// value: option.Some(20)
//
//	// Set the element at index 1
//	updated := second.Set(25)(numbers)
//	// updated: []int{10, 25, 30, 40}
//
//	// Out of bounds access returns None
//	outOfBounds := OA.At[int](10)
//	value = outOfBounds.GetOption(numbers)
//	// value: option.None[int]()
//
//	// Out of bounds set is a no-op
//	unchanged := outOfBounds.Set(99)(numbers)
//	// unchanged: []int{10, 20, 30, 40} (original array)
//
// See Also:
//   - AR.Lookup: Gets an element at an index, returning an Option
//   - AR.UpdateAt: Updates an element at an index, returning an Option
//   - OP.Optional: The Optional optic type
func At[A any](idx int) OP.Optional[[]A, A] {
	return G.At[[]A](idx)
}
