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
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/reader"
)

// FromLens converts a Lens into a Traversal that works with endomorphisms.
//
// The traversal is the reader profunctor mapping of the transformation function f:
// the lens getter is pre-composed and the lifted setter is post-composed, giving
// fmap(Set) ∘ f ∘ Get, a function from S to HKTES. The resulting traversal focuses on
// exactly one value; its Endomorphism[S] sets the focus to the value computed from the
// source that was passed to the traversal.
//
// Type Parameters:
//   - S: The source structure type
//   - A: The focus type
//   - HKTES: Higher-kinded type for Endomorphism[S]
//   - HKTA: Higher-kinded type for A
//
// Parameters:
//   - fmap: Functor map operation that lifts the lens setter func(A) Endomorphism[S] into the effect
//
// Returns:
//   - A function that converts a Lens[S, A] into a Traversal[S, A, HKTES, HKTA]
//
// Example:
//
//	import (
//	    "github.com/IBM/fp-go/v2/endomorphism"
//	    "github.com/IBM/fp-go/v2/identity"
//	    "github.com/IBM/fp-go/v2/optics/lens"
//	)
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLens(
//	    func(p Person) string { return p.Name },
//	    func(p Person, name string) Person {
//	        p.Name = name
//	        return p
//	    },
//	)
//
//	nameTrav := FromLens[Person, string](identity.Map[string, endomorphism.Endomorphism[Person]])(nameLens)
//
//	person := Person{Name: "alice", Age: 30}
//	updated := nameTrav(strings.ToUpper)(person)(person) // Person{Name: "ALICE", Age: 30}
//
// See Also:
//   - optional/generic.FromOptional: For optics that focus on zero or one value
//   - traversable/generic.FromTraversableLens: For lenses that focus on a traversable structure
func FromLens[S, A, HKTES, HKTA any](
	fmap functor.MapType[A, Endomorphism[S], HKTA, HKTES],
) func(Lens[S, A]) Traversal[S, A, HKTES, HKTA] {
	return func(sa Lens[S, A]) Traversal[S, A, HKTES, HKTA] {
		return reader.ProMap(sa.Get, fmap(sa.Set))
	}
}
