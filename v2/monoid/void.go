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

package monoid

import (
	"github.com/IBM/fp-go/v2/function"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// VoidMonoid creates a Monoid for the Void (unit) type.
//
// The Void type has exactly one value (function.VOID), making it trivial to define
// a monoid. This monoid uses the Last semigroup, which always returns the second
// argument, though since all Void values are identical, the choice of semigroup
// doesn't affect the result.
//
// This monoid is useful in contexts where:
//   - A monoid instance is required but no meaningful data needs to be combined
//   - You need to track that an operation occurred without caring about its result
//   - Building generic abstractions that work with any monoid, including the trivial case
//
// # Monoid Laws
//
// The VoidMonoid satisfies all monoid laws trivially:
//   - Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z)) - always VOID
//   - Left Identity: Concat(Empty(), x) = x - always VOID
//   - Right Identity: Concat(x, Empty()) = x - always VOID
//
// Returns:
//   - A Monoid[Void] instance
//
// Example:
//
//	m := VoidMonoid()
//	result := m.Concat(function.VOID, function.VOID)  // function.VOID
//	empty := m.Empty()                                 // function.VOID
//
//	// Useful for tracking operations without data
//	type Action = func() Void
//	actions := []Action{
//	    func() Void { fmt.Println("Action 1"); return function.VOID },
//	    func() Void { fmt.Println("Action 2"); return function.VOID },
//	}
//	// Execute all actions and combine results
//	results := A.Map(func(a Action) Void { return a() })(actions)
//	_ = ConcatAll(m)(results)  // All actions executed, result is VOID
func VoidMonoid() Monoid[Void] {
	return MakeMonoid(
		S.Last[Void]().Concat,
		function.VOID,
	)
}
