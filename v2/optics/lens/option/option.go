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

package option

import (
	LG "github.com/IBM/fp-go/v2/optics/lens/generic"
	T "github.com/IBM/fp-go/v2/optics/traversal/option"
	O "github.com/IBM/fp-go/v2/option"
)

// AsTraversal converts a Lens[S, A] to a Traversal[S, A] for optional values.
//
// A traversal is a generalization of a lens that can focus on zero or more values.
// This function converts a lens (which focuses on exactly one value) into a traversal,
// allowing it to be used with traversal operations like mapping over multiple values.
//
// This is particularly useful when you want to:
//   - Use lens operations in a traversal context
//   - Compose lenses with traversals
//   - Apply operations that work on collections of optional values
//
// The conversion uses the Option monad's map operation to handle the optional nature
// of the values being traversed.
//
// Type Parameters:
//   - S: The structure type containing the field
//   - A: The type of the field being focused on
//
// Returns:
//   - A function that takes a Lens[S, A] and returns a Traversal[S, A]
//
// Example:
//
//	type Config struct {
//	    Timeout Option[int]
//	}
//
//	timeoutLens := lens.MakeLens(
//	    func(c Config) Option[int] { return c.Timeout },
//	    func(c Config, t Option[int]) Config { c.Timeout = t; return c },
//	)
//
//	// Convert to traversal for use with traversal operations
//	timeoutTraversal := lens.AsTraversal[Config, int]()(timeoutLens)
//
//	// Now can use traversal operations
//	configs := []Config{{Timeout: O.Some(30)}, {Timeout: O.None[int]()}}
//	// Apply operations across all configs using the traversal
func AsTraversal[S, A any]() func(Lens[S, A]) T.Traversal[S, A] {
	return LG.AsTraversal[T.Traversal[S, A]](O.MonadMap[A, S])
}
