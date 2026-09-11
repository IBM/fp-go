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

// Package traversalendo provides traversal optics that return Endomorphism[S] instead of S.
//
// # Overview
//
// A traversal endomorphism is a specialized form of traversal that works with Endomorphism[S]
// as its effect type rather than directly with S. This design choice enables powerful
// composition capabilities that are not available with regular traversals under an
// applicative-only constraint.
//
// The key insight is that while regular traversals can be composed sequentially (nested
// access), they cannot be easily combined in parallel (accessing multiple fields at the
// same level) without a monad constraint. Traversal endomorphisms solve this by forming
// a monoid, allowing multiple traversals to be combined using standard monoid operations.
//
// # Core Concepts
//
// Type Signature:
//
//	Traversal[S, A, HKTS, HKTA] = func(func(A) HKTA) func(S) HKTS
//
// For traversal endomorphisms, HKTS is instantiated as HKT[Endomorphism[S]], meaning
// the traversal returns an endomorphism (a function from S to S) wrapped in an effect
// context, rather than returning S directly.
//
// This representation has several important properties:
//
//  1. Monoid Structure: Traversal endomorphisms form a monoid with:
//     - Empty: A traversal that focuses on no values (identity element)
//     - Concat: Combines two traversals to focus on all values from both
//
//  2. Composability: Multiple traversal endomorphisms can be combined using monoid
//     operations like Fold, allowing you to build complex traversals from simpler ones.
//
//  3. Conversion: A traversal endomorphism can be converted to a regular traversal
//     using ToTraversal, but the reverse conversion is not possible. This makes
//     traversal endomorphisms more powerful for construction, while regular traversals
//     are more convenient for use.
//
// # Why Endomorphisms?
//
// Regular traversals under an applicative constraint cannot be easily combined in
// parallel because applicative functors don't provide a way to sequence effects that
// depend on previous results. However, endomorphisms (functions from S to S) can be
// composed using function composition, which forms a monoid.
//
// By working with Endomorphism[S] as the effect type, we can:
//
//  1. Use the endomorphism monoid to combine multiple traversals
//  2. Build complex traversals incrementally using monoid operations
//  3. Convert the final result to a regular traversal for actual use
//
// This approach separates the construction phase (using endomorphisms) from the
// usage phase (using regular traversals), providing both power and convenience.
//
// # Basic Usage
//
// Creating and combining traversal endomorphisms:
//
//	import (
//	    A "github.com/IBM/fp-go/v2/array"
//	    thunk "github.com/IBM/fp-go/v2/context/readerioresult"
//	    "github.com/IBM/fp-go/v2/endomorphism"
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/monoid"
//	    TA "github.com/IBM/fp-go/v2/optics/traversalendo/array/generic"
//	    TE "github.com/IBM/fp-go/v2/optics/traversalendo/generic"
//	    TL "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
//	)
//
//	type Address struct {
//	    Street   string
//	    Name     string
//	    Keywords []string
//	}
//
//	// Create a monoid for combining traversal endomorphisms
//	m := TE.MakeMonoid[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Of,
//	    thunk.Map,
//	    thunk.Ap,
//	)
//
//	// Create individual field traversals
//	lenses := MakeAddressRefLenses()
//	fromString := TL.FromLens[*Address, string](thunk.Map)
//	fromArray := TA.FromArrayLens[[]string, *Address, string](
//	    thunk.Of[[]string],
//	    thunk.Map[[]string, func(string) []string],
//	    thunk.Map[[]string, endomorphism.Endomorphism[*Address]],
//	    thunk.Ap[[]string, string],
//	)
//
//	streetTrav := fromString(lenses.Street)
//	nameTrav := fromString(lenses.Name)
//	keywordsTrav := fromArray(lenses.Keywords)
//
//	// Combine all traversals using monoid fold
//	allFieldsTrav := F.Pipe1(
//	    A.From(streetTrav, nameTrav, keywordsTrav),
//	    monoid.Fold(m),
//	)
//
//	// Convert to a regular traversal for use
//	regularTrav := TE.ToTraversal[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Map,
//	)(allFieldsTrav)
//
//	// Use the traversal to modify all fields
//	addr := &Address{
//	    Street:   "main st",
//	    Name:     "john",
//	    Keywords: []string{"a", "b"},
//	}
//	res := regularTrav(func(s string) thunk.ReaderIOResult[string] {
//	    return thunk.Of(strings.ToUpper(s))
//	})(addr)(ctx)()
//	// res == result.Of(&Address{
//	//     Street:   "MAIN ST",
//	//     Name:     "JOHN",
//	//     Keywords: []string{"A", "B"},
//	// })
//
// # Monoid Operations
//
// The monoid structure enables powerful composition patterns:
//
//	// Empty traversal (identity element)
//	emptyTrav := TE.Empty[string, thunk.ReaderIOResult[string], *Address](thunk.Of)
//
//	// Concat two traversals
//	combined := TE.Concat[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Map,
//	    thunk.Ap,
//	)(streetTrav, nameTrav)
//
//	// Fold multiple traversals
//	allFields := monoid.Fold(m)(A.From(streetTrav, nameTrav, keywordsTrav))
//
// Concat runs the effects of the left traversal before the effects of the right traversal.
// The resulting endomorphisms compose as left ∘ right, so if both traversals focus on the
// same value, the update produced by the left traversal wins.
//
// # Conversion to Regular Traversals
//
// Once you've built a traversal endomorphism, convert it to a regular traversal:
//
//	regularTrav := TE.ToTraversal[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Map,
//	)(traversalEndo)
//
// This conversion is one-way: you can convert a traversal endomorphism to a regular
// traversal, but not vice versa. This asymmetry reflects the fact that traversal
// endomorphisms are more powerful for construction but less convenient for direct use.
//
// # Comparison with Regular Traversals
//
// Regular Traversal:
//
//	Traversal[S, A, HKTS, HKTA] = func(func(A) HKTA) func(S) HKTS
//	where HKTS = HKT[S]
//
// Traversal Endomorphism:
//
//	Traversal[S, A, HKTES, HKTA] = func(func(A) HKTA) func(S) HKTES
//	where HKTES = HKT[Endomorphism[S]]
//
// Key Differences:
//
//  1. Effect Type: Regular traversals work with any effect type, while traversal
//     endomorphisms specifically use Endomorphism[S] as the effect type.
//
//  2. Composition: Traversal endomorphisms can be combined using monoid operations
//     (Empty, Concat, Fold), while regular traversals require monad constraints for
//     similar parallel composition.
//
//  3. Conversion: Traversal endomorphisms can be converted to regular traversals
//     using ToTraversal, but the reverse is not possible.
//
//  4. Use Case: Traversal endomorphisms are ideal for building complex traversals
//     incrementally, while regular traversals are more convenient for direct use.
//
// # Advanced Patterns
//
// Building traversals from lenses and optionals:
//
//	// From a simple lens
//	fieldTrav := TL.FromLens[*Address, string](thunk.Map)(fieldLens)
//
//	// From a lens focusing on an array
//	arrayTrav := TA.FromArrayLens[[]string, *Address, string](
//	    thunk.Of[[]string],
//	    thunk.Map[[]string, func(string) []string],
//	    thunk.Map[[]string, endomorphism.Endomorphism[*Address]],
//	    thunk.Ap[[]string, string],
//	)(arrayLens)
//
//	// From a lens focusing on any traversable structure, e.g. an option
//	optionTrav := TT.FromTraversableLens[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Map[option.Option[string], endomorphism.Endomorphism[*Address]],
//	)(option.MakeTraversable[string, string](
//	    thunk.Of[option.Option[string]],
//	    thunk.Map[string, option.Option[string]],
//	))(optionLens)
//
//	// From an optional that focuses on zero or one value
//	optionalTrav := TO.FromOptional[*Address, string](
//	    thunk.Of[endomorphism.Endomorphism[*Address]],
//	    thunk.Map[string, endomorphism.Endomorphism[*Address]],
//	)(fieldOptional)
//
// Composing traversal endomorphisms:
//
//	// Compose for nested access (sequential): outerTrav focuses on A within S,
//	// innerTrav focuses on B within A
//	nested := F.Pipe1(
//	    outerTrav,
//	    TE.Compose[S, B, thunk.ReaderIOResult[endomorphism.Endomorphism[S]], thunk.ReaderIOResult[B]](
//	        thunk.Map[endomorphism.Endomorphism[A], A],
//	    )(innerTrav),
//	)
//
//	// Combine for parallel access (using monoid)
//	parallel := monoid.Fold(m)(A.From(trav1, trav2, trav3))
//
// # Type Parameters
//
// Throughout these packages, you'll encounter these type parameters:
//
//   - S: The source type (the structure being traversed)
//   - A: The focus type (the values being accessed or modified)
//   - HKTES: Higher-kinded type for Endomorphism[S] in the effect context (named HKTS
//     in the array and traversable sub-packages)
//   - HKTA: Higher-kinded type for A in the effect context
//   - HKTS: In ToTraversal, the higher-kinded type for S of the resulting regular traversal
//
// # Sub-packages
//
//   - generic: Empty, Concat and MakeMonoid (the monoid), Compose and ToTraversal
//   - lens/generic: FromLens builds a traversal endomorphism from a lens
//   - optional/generic: FromOptional and FromArrayOptional build traversal endomorphisms from optionals
//   - array/generic: FromArrayLens traverses the elements of an array focused by a lens
//   - traversable/generic: FromTraversableLens traverses any traversable structure focused by a lens
//
// # Related Packages
//
//   - optics/traversal: Regular traversal optics
//   - optics/lens: For focusing on a single field
//   - optics/optional: For focusing on optional values
//   - endomorphism: Endomorphism type and operations
//   - monoid: Monoid algebraic structure
package traversalendo
