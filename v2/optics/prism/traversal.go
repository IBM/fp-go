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

package prism

import (
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

// AsTraversal converts a Prism into a Traversal.
//
// A Traversal is a more general optic that can focus on zero or more values,
// while a Prism focuses on zero or one value. This function lifts a Prism
// into the Traversal abstraction, allowing it to be used in contexts that
// expect traversals.
//
// The conversion works by:
//   - If the prism matches (GetOption returns Some), the traversal focuses on that value
//   - If the prism doesn't match (GetOption returns None), the traversal focuses on zero values
//
// Type Parameters:
//   - R: The traversal function type ~func(func(A) HKTA) func(S) HKTS
//   - S: The source type
//   - A: The focus type
//   - HKTS: Higher-kinded type for S (e.g., functor/applicative context)
//   - HKTA: Higher-kinded type for A (e.g., functor/applicative context)
//
// Parameters:
//   - fof: Function to lift S into the higher-kinded type HKTS (pure/of operation)
//   - fmap: Function to map over HKTA and produce HKTS (functor map operation)
//
// Returns:
//   - A function that converts a Prism[S, A] into a Traversal R
//
// Example:
//
//	// Convert a prism to a traversal for use with applicative functors
//	prism := MakePrism(...)
//	traversal := AsTraversal(
//	    func(s S) HKTS { return pure(s) },
//	    func(hkta HKTA, f func(A) S) HKTS { return fmap(hkta, f) },
//	)(prism)
//
// Note: This function is typically used in advanced scenarios involving
// higher-kinded types and applicative functors. Most users will work
// directly with prisms rather than converting them to traversals.
func AsTraversal[R ~func(func(A) HKTA) func(S) HKTS, S, A, HKTS, HKTA any](
	fof func(S) HKTS,
	fmap func(HKTA, func(A) S) HKTS,
) func(Prism[S, A]) R {
	return func(sa Prism[S, A]) R {
		return func(f func(a A) HKTA) func(S) HKTS {
			return func(s S) HKTS {
				return F.Pipe2(
					s,
					sa.GetOption,
					O.Fold(
						// If prism doesn't match, return the original value lifted into HKTS
						F.Nullary2(F.Constant(s), fof),
						// If prism matches, apply f to the extracted value and map back
						func(a A) HKTS {
							return fmap(f(a), func(a A) S {
								return prismModify(F.Constant1[A](a), sa, s)
							})
						},
					),
				)
			}
		}
	}
}
