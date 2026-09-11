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
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/reader"
	SG "github.com/IBM/fp-go/v2/semigroup"
)

// Empty creates an empty traversal endomorphism that focuses on no values.
//
// An empty traversal endomorphism is the identity element of the traversal endomorphism
// monoid. When composed with another traversal using Concat, it has no effect. This is
// useful as a base case when building traversals dynamically or when you need a traversal
// that performs no operations.
//
// The empty traversal simply lifts the source value S into an endomorphism context
// without applying any transformation to focused values. It wraps the identity
// endomorphism in the appropriate reader and pointed structure.
//
// Type Parameters:
//   - A: The focus type (not used, as no values are focused)
//   - HKTA: Higher-kinded type for A (not used)
//   - S: The source type
//   - HKTES: Higher-kinded type for Endomorphism[S]
//
// Parameters:
//   - fof: Function to lift Endomorphism[S] into HKTES (pure/of operation)
//
// Returns:
//   - A traversal endomorphism that focuses on no values
//
// Example:
//
//	import (
//	    thunk "github.com/IBM/fp-go/v2/context/readerioresult"
//	)
//
//	type Address struct {
//	    Street string
//	    Name   string
//	}
//
//	// Create an empty traversal endomorphism
//	emptyTrav := Empty[string, thunk.ReaderIOResult[string], *Address](thunk.Of)
//
//	// Empty is the identity for Concat
//	concat := Concat[string, thunk.ReaderIOResult[string], *Address](thunk.Map, thunk.Ap)
//	concat(emptyTrav, trav) // behaves like trav
//	concat(trav, emptyTrav) // behaves like trav
//
// See Also:
//   - Concat: Combine two traversal endomorphisms
//   - MakeMonoid: Create a monoid instance for traversal endomorphisms
func Empty[A, HKTA, S, HKTES any](
	fof pointed.OfType[Endomorphism[S], HKTES],
) Traversal[S, A, HKTES, HKTA] {
	return reader.Of[func(A) HKTA](reader.Of[S](fof(endomorphism.Identity[S]())))
}

// Concat combines two traversal endomorphisms into a single traversal that focuses on all values from both.
//
// This function creates a new traversal endomorphism that applies both input traversals in
// sequence, combining their effects using applicative operations. The resulting traversal
// will focus on all values that either the left or right traversal focuses on.
//
// Concat is associative, meaning:
//
//	Concat(fmap, fap)(a, Concat(fmap, fap)(b, c)) == Concat(fmap, fap)(Concat(fmap, fap)(a, b), c)
//
// Together with Empty, Concat forms a monoid for traversal endomorphisms, allowing you to
// combine multiple traversals using standard monoid operations like fold or reduce.
//
// The implementation is derived from standard semigroups: the endomorphism semigroup is lifted
// into the effect via ApplySemigroup, and then lifted pointwise over the source type S and over
// the transformation function via FunctionSemigroup. The effects of the left traversal are
// sequenced before the effects of the right traversal.
//
// The resulting endomorphisms compose as l ∘ r: the update produced by the right traversal is
// applied first, the update produced by the left traversal last. For traversals with disjoint
// foci this order is irrelevant; if both traversals focus on the same value, the update of the
// left traversal wins.
//
// Type Parameters:
//   - A: The focus type
//   - HKTA: Higher-kinded type for A
//   - S: The source type
//   - HKTES: Higher-kinded type for Endomorphism[S]
//   - HKTESES: Higher-kinded type for Endomorphism[Endomorphism[S]]
//
// Parameters:
//   - fmap: Functor map operation for endomorphisms
//   - fap: Applicative apply operation for sequencing effects
//
// Returns:
//   - A function that takes two traversals and returns their concatenation
//
// Example:
//
//	import (
//	    thunk "github.com/IBM/fp-go/v2/context/readerioresult"
//	    F "github.com/IBM/fp-go/v2/function"
//	    TLG "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
//	)
//
//	type Address struct {
//	    Street string
//	    Name   string
//	}
//
//	lenses := MakeAddressRefLenses()
//	fromString := TLG.FromLens[*Address, string](thunk.Map)
//
//	streetTrav := fromString(lenses.Street)
//	nameTrav := fromString(lenses.Name)
//
//	// Combine to traverse both fields
//	bothFieldsTrav := Concat[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Map,
//	    thunk.Ap,
//	)(streetTrav, nameTrav)
//
//	// Modify both fields; the result is an effectful endomorphism of *Address,
//	// use ToTraversal to obtain the modified structure directly
//	addr := &Address{Street: "main st", Name: "john"}
//	endo := bothFieldsTrav(func(s string) thunk.ReaderIOResult[string] {
//	    return thunk.Of(strings.ToUpper(s))
//	})(addr)
//
// See Also:
//   - Empty: Create an empty traversal endomorphism (identity element)
//   - MakeMonoid: Create a monoid instance for traversal endomorphisms
func Concat[A, HKTA, S, HKTES, HKTESES any](
	fmap functor.MapType[Endomorphism[S], Endomorphism[Endomorphism[S]], HKTES, HKTESES],
	fap apply.ApType[HKTES, HKTES, HKTESES],
) func(l, r Traversal[S, A, HKTES, HKTA]) Traversal[S, A, HKTES, HKTA] {
	return F.Pipe2(
		SG.ApplySemigroup(fmap, fap, endomorphism.Semigroup[S]()),
		SG.FunctionSemigroup[S, HKTES],
		SG.FunctionSemigroup[func(A) HKTA, func(S) HKTES],
	).Concat
}

// MakeMonoid creates a monoid instance for traversal endomorphisms.
//
// A monoid provides two operations: an identity element (Empty) and an associative
// binary operation (Concat). This allows traversal endomorphisms to be combined using
// standard monoid operations, making it easy to build complex traversals from simpler ones.
//
// The monoid laws are satisfied:
//  1. Identity: Concat(Empty, x) == x and Concat(x, Empty) == x
//  2. Associativity: Concat(Concat(a, b), c) == Concat(a, Concat(b, c))
//
// This monoid instance enables you to use traversal endomorphisms with generic monoid
// operations like fold, reduce, or mconcat to combine multiple traversals into one.
//
// Type Parameters:
//   - A: The focus type
//   - HKTA: Higher-kinded type for A
//   - S: The source type
//   - HKTES: Higher-kinded type for Endomorphism[S]
//   - HKTESES: Higher-kinded type for Endomorphism[Endomorphism[S]]
//
// Parameters:
//   - fof: Function to lift Endomorphism[S] into HKTES (for Empty)
//   - fmap: Functor map operation for endomorphisms (for Concat)
//   - fap: Applicative apply operation for sequencing effects (for Concat)
//
// Returns:
//   - A monoid instance for traversal endomorphisms
//
// Example:
//
//	import (
//	    A "github.com/IBM/fp-go/v2/array"
//	    thunk "github.com/IBM/fp-go/v2/context/readerioresult"
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/monoid"
//	    TLG "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
//	)
//
//	type Address struct {
//	    Street string
//	    Name   string
//	}
//
//	// Create monoid for traversal endomorphisms
//	m := MakeMonoid[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Of,
//	    thunk.Map,
//	    thunk.Ap,
//	)
//
//	// Create individual field traversals
//	lenses := MakeAddressRefLenses()
//	fromString := TLG.FromLens[*Address, string](thunk.Map)
//	streetTrav := fromString(lenses.Street)
//	nameTrav := fromString(lenses.Name)
//
//	// Combine all field traversals using monoid fold
//	allFieldsTrav := F.Pipe1(
//	    A.From(streetTrav, nameTrav),
//	    monoid.Fold(m),
//	)
//
//	// Convert to a regular traversal and modify all fields at once
//	addrTrav := ToTraversal[string, thunk.ReaderIOResult[string], *Address](thunk.Map)(allFieldsTrav)
//	addr := &Address{Street: "street", Name: "name"}
//	res := addrTrav(func(s string) thunk.ReaderIOResult[string] {
//	    return thunk.Of(strings.ToUpper(s))
//	})(addr)(ctx)()
//	// res == result.Of(&Address{Street: "STREET", Name: "NAME"})
//
// See Also:
//   - Empty: The identity element of the monoid
//   - Concat: The binary operation of the monoid
//   - monoid.Fold: Combine multiple monoid values
//   - monoid.MakeMonoid: Generic monoid constructor
func MakeMonoid[A, HKTA, S, HKTES, HKTESES any](
	fof pointed.OfType[Endomorphism[S], HKTES],
	fmap functor.MapType[Endomorphism[S], Endomorphism[Endomorphism[S]], HKTES, HKTESES],
	fap apply.ApType[HKTES, HKTES, HKTESES],
) Monoid[Traversal[S, A, HKTES, HKTA]] {
	return M.MakeMonoid(
		Concat[A, HKTA](fmap, fap),
		Empty[A, HKTA](fof),
	)
}
