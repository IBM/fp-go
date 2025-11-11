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

// Package generic provides generic array operations for custom Reader types.
package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
)

// MonadTraverseArray transforms each element of an array using a function that returns a generic Reader,
// then collects the results into a single generic Reader containing an array.
// This is the monadic version that takes the array as the first parameter.
//
// This generic version works with custom reader types that match the pattern ~func(R) B.
//
// Type Parameters:
//   - GB: The generic Reader type for individual elements (~func(R) B)
//   - GBS: The generic Reader type for the result array (~func(R) BBS)
//   - AAS: The input array type (~[]A)
//   - BBS: The output array type (~[]B)
//   - R: The environment/context type
//   - A: The input element type
//   - B: The output element type
func MonadTraverseArray[GB ~func(R) B, GBS ~func(R) BBS, AAS ~[]A, BBS ~[]B, R, A, B any](tas AAS, f func(A) GB) GBS {
	return RA.MonadTraverse(
		Of[GBS, R, BBS],
		Map[GBS, func(R) func(B) BBS, R, BBS, func(B) BBS],
		Ap[GB, GBS, func(R) func(B) BBS, R, B, BBS],
		tas, f,
	)
}

// TraverseArray transforms each element of an array using a function that returns a generic Reader,
// then collects the results into a single generic Reader containing an array.
//
// This generic version works with custom reader types that match the pattern ~func(R) B.
//
// Type Parameters:
//   - GB: The generic Reader type for individual elements (~func(R) B)
//   - GBS: The generic Reader type for the result array (~func(R) BBS)
//   - AAS: The input array type (~[]A)
//   - BBS: The output array type (~[]B)
//   - R: The environment/context type
//   - A: The input element type
//   - B: The output element type
func TraverseArray[GB ~func(R) B, GBS ~func(R) BBS, AAS ~[]A, BBS ~[]B, R, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, R, BBS],
		Map[GBS, func(R) func(B) BBS, R, BBS, func(B) BBS],
		Ap[GB, GBS, func(R) func(B) BBS, R, B, BBS],
		f,
	)
}

// TraverseArrayWithIndex transforms each element of an array using a function that takes
// both the index and the element, returning a generic Reader. The results are collected into
// a single generic Reader containing an array.
//
// This generic version works with custom reader types that match the pattern ~func(R) B.
//
// Type Parameters:
//   - GB: The generic Reader type for individual elements (~func(R) B)
//   - GBS: The generic Reader type for the result array (~func(R) BBS)
//   - AAS: The input array type (~[]A)
//   - BBS: The output array type (~[]B)
//   - R: The environment/context type
//   - A: The input element type
//   - B: The output element type
func TraverseArrayWithIndex[GB ~func(R) B, GBS ~func(R) BBS, AAS ~[]A, BBS ~[]B, R, A, B any](f func(int, A) GB) func(AAS) GBS {
	return RA.TraverseWithIndex[AAS](
		Of[GBS, R, BBS],
		Map[GBS, func(R) func(B) BBS, R, BBS, func(B) BBS],
		Ap[GB, GBS, func(R) func(B) BBS, R, B, BBS],
		f,
	)
}

// SequenceArray converts an array of generic Readers into a single generic Reader containing an array.
// All Readers in the input array share the same environment and are evaluated with it.
//
// This generic version works with custom reader types that match the pattern ~func(R) A.
//
// Type Parameters:
//   - GA: The generic Reader type for individual elements (~func(R) A)
//   - GAS: The generic Reader type for the result array (~func(R) AAS)
//   - AAS: The array type (~[]A)
//   - GAAS: The input array of Readers type (~[]GA)
//   - R: The environment/context type
//   - A: The element type
func SequenceArray[GA ~func(R) A, GAS ~func(R) AAS, AAS ~[]A, GAAS ~[]GA, R, A any](ma GAAS) GAS {
	return MonadTraverseArray[GA, GAS](ma, F.Identity[GA])
}
