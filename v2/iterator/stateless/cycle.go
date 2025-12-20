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

package stateless

import (
	G "github.com/IBM/fp-go/v2/iterator/stateless/generic"
)

// Cycle creates an [Iterator] that repeats the elements of the input [Iterator] indefinitely.
// The iterator cycles through all elements of the input, and when it reaches the end, it starts over from the beginning.
// This creates an infinite iterator, so it should be used with caution and typically combined with operations that limit the output.
//
// Type Parameters:
//   - U: The type of elements in the iterator
//
// Parameters:
//   - ma: The input iterator to cycle through
//
// Returns:
//   - An iterator that infinitely repeats the elements of the input iterator
//
// Example:
//
//	iter := stateless.FromArray([]int{1, 2, 3})
//	cycled := stateless.Cycle(iter)
//	// Produces: 1, 2, 3, 1, 2, 3, 1, 2, 3, ... (infinitely)
//
//	// Typically used with Take to limit output:
//	limited := stateless.Take(7)(cycled)
//	// Produces: 1, 2, 3, 1, 2, 3, 1
func Cycle[U any](ma Iterator[U]) Iterator[U] {
	return G.Cycle(ma)
}
