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

package identity

// MakeTraversable creates a traversal function for Identity types.
// Since Identity has no computational context, traversing is equivalent to mapping.
//
// This function enables traversing an Identity value by applying a transformation that produces
// a higher-kinded type. Because Identity is just a value with no wrapper, the traversal
// simply applies the transformation function directly.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type after transformation
//   - HKTB: The higher-kinded type containing B (e.g., IO[B], Option[B])
//
// Returns:
//   - A function that takes a transformation function and returns a function that applies it
//
// Behavior:
//   - The value is transformed by f, producing HKTB directly
//   - No sequencing or lifting is needed since Identity has no context
//
// Example:
//
//	import (
//	    I "github.com/IBM/fp-go/v2/identity"
//	    IO "github.com/IBM/fp-go/v2/io"
//	    "strconv"
//	)
//
//	// Create a traversable for Identity that works with IO
//	traverseIO := I.MakeTraversable[int, string, IO.IO[string]]()
//
//	// Use it to transform int to IO[string]
//	fetchUser := func(id int) IO.IO[string] {
//	    return IO.Of(strconv.Itoa(id))
//	}
//
//	result := traverseIO(fetchUser)(42)  // IO["42"]
//
// See Also:
//   - Map: The underlying implementation (Identity traversal is just mapping)
func MakeTraversable[A, B, HKTB any]() func(func(A) HKTB) func(A) HKTB {
	return Map
}
