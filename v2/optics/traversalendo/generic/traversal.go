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
	"github.com/IBM/fp-go/v2/reader"
)

// ToTraversal converts a traversal endomorphism into a regular traversal.
//
// A traversal endomorphism works with Endomorphism[S] as its effect type, which represents
// functions from S to S. This is useful for building traversals using monoid operations
// like Concat and Empty. However, to actually use the traversal for modification or other
// operations, you need to convert it to a regular traversal that works with the desired
// effect type (like Identity, Const, Option, etc.).
//
// ToTraversal performs this conversion by:
//  1. Lifting endomorphism.Read with the functor map, giving func(S) func(HKTES) HKTS
//  2. Applying that reader to the result of the traversal endomorphism via reader.MonadAp
//  3. Post-composing the traversal endomorphism with this fixed step via reader.Map
//
// This allows you to build complex traversals using the monoid structure of traversal
// endomorphisms, then convert the result to a usable traversal for your specific effect type.
//
// Type Parameters:
//   - A: The focus type (the values being accessed or modified)
//   - HKTA: Higher-kinded type for A in the target traversal
//   - S: The source type (the structure being traversed)
//   - HKTES: Higher-kinded type for Endomorphism[S] (used in the endomorphism traversal)
//   - HKTS: Higher-kinded type for S in the target traversal
//
// Parameters:
//   - fmap: Functor map operation that transforms Endomorphism[S] into S within the effect context
//
// Returns:
//   - A function that converts a traversal endomorphism into a regular traversal
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
//	// Build a traversal endomorphism using monoid operations
//	m := MakeMonoid[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Of,
//	    thunk.Map,
//	    thunk.Ap,
//	)
//
//	lenses := MakeAddressRefLenses()
//	fromString := TLG.FromLens[*Address, string](thunk.Map)
//
//	// Combine multiple field traversals
//	addrTravEndo := F.Pipe1(
//	    A.From(fromString(lenses.Street), fromString(lenses.Name)),
//	    monoid.Fold(m),
//	)
//
//	// Convert to a regular traversal for use
//	addrTrav := ToTraversal[string, thunk.ReaderIOResult[string], *Address](
//	    thunk.Map,
//	)(addrTravEndo)
//
//	// Now use the traversal to modify all fields
//	addr := &Address{Street: "main st", Name: "john"}
//	res := addrTrav(func(s string) thunk.ReaderIOResult[string] {
//	    return thunk.Of(strings.ToUpper(s))
//	})(addr)(ctx)()
//	// res == result.Of(&Address{Street: "MAIN ST", Name: "JOHN"})
//
// See Also:
//   - MakeMonoid: Create a monoid for combining traversal endomorphisms
//   - Concat: Combine two traversal endomorphisms
//   - Empty: Create an empty traversal endomorphism
func ToTraversal[A, HKTA, S, HKTES, HKTS any](
	fmap functor.MapType[Endomorphism[S], S, HKTES, HKTS],
) func(Traversal[S, A, HKTES, HKTA]) Traversal[S, A, HKTS, HKTA] {
	return reader.Map[func(A) HKTA](
		F.Bind1st(
			reader.MonadAp[HKTS, S, HKTES],
			F.Flow2(endomorphism.Read[S], fmap),
		),
	)
}
