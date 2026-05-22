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

package generic

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
)

// TraversalFromOptional converts an Optional into a Traversal that works with endomorphisms.
//
// This function creates a traversal that focuses on zero or one value within a structure,
// using endomorphisms as the effect type. When the optional value exists, the traversal
// applies the provided function and returns an endomorphism that updates the structure.
// When the optional value doesn't exist, it returns the identity endomorphism (no-op).
//
// The key difference from the standard optional traversal is that this version uses
// Endomorphism[S] as the effect type (HKTES), which means the traversal produces
// functions that transform the entire structure S rather than wrapping it in another
// effect type.
//
// Type Parameters:
//   - S: The source structure type
//   - A: The focus type (the value the optional may contain)
//   - HKTES: Higher-kinded type for Endomorphism[S] (the effect wrapping S transformations)
//   - HKTA: Higher-kinded type for A (the effect wrapping the focus value)
//
// Parameters:
//   - fof: Pointed instance that lifts an Endomorphism[S] into HKTES
//   - fmap: Functor instance that maps over HKTA to produce HKTES
//
// Returns:
//
//	A function that converts an Optional[S, A] into a Traversal[S, A, HKTES, HKTA]
//
// The returned traversal has the signature:
//
//	func(f func(A) HKTA) func(S) HKTES
//
// When applied:
//   - If the optional value exists: applies f to the value and maps the result to update S
//   - If the optional value doesn't exist: returns the identity endomorphism (no change to S)
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	// Optional that focuses on non-empty names
//	nameOpt := optional.MakeOptional(
//	    func(p Person) option.Option[string] {
//	        if p.Name != "" { return option.Some(p.Name) }
//	        return option.None[string]()
//	    },
//	    func(name string) endomorphism.Endomorphism[Person] {
//	        return func(p Person) Person {
//	            p.Name = name
//	            return p
//	        }
//	    },
//	)
//
//	// Convert to traversal with identity functor
//	traversal := TraversalFromOptional(
//	    identity.Of[endomorphism.Endomorphism[Person]],
//	    identity.Map[string, endomorphism.Endomorphism[Person]],
//	)(nameOpt)
//
//	// Use the traversal to modify names
//	toUpper := func(s string) identity.Identity[string] {
//	    return identity.Of(strings.ToUpper(s))
//	}
//	modify := traversal(toUpper)
//	person := Person{Name: "alice", Age: 30}
//	updated := modify(person) // Person{Name: "ALICE", Age: 30}
//
// See Also:
//   - Optional: The optional optic type
//   - Traversal: The traversal optic type
//   - TraversalFromLens: For converting lenses to traversals
func FromOptional[S, A, HKTES, HKTA any](
	fof pointed.OfType[Endomorphism[S], HKTES],
	fmap functor.MapType[A, Endomorphism[S], HKTA, HKTES],
) func(Optional[S, A]) Traversal[S, A, HKTES, HKTA] {
	onNone := lazy.Of(fof(endomorphism.Identity[S]()))
	return func(sa Optional[S, A]) Traversal[S, A, HKTES, HKTA] {
		saGet := sa.GetOption
		saSet := fmap(sa.Set)
		return func(f func(A) HKTA) func(S) HKTES {
			return F.Flow2(
				saGet,
				option.Fold(
					onNone,
					F.Flow2(
						f,
						saSet,
					),
				),
			)
		}
	}
}
