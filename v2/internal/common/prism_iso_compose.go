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

package common

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

// PrismComposeIso composes a Prism[S, A] with an Iso[A, B] to produce a Prism[S, B].
//
// Because an Iso is a total bidirectional transformation, composing a prism with an
// iso always yields another prism — the result is never partial due to the iso.
// GetOption returns None only when the original prism does not match.
//
// Concretely:
//   - GetOption: applies the prism's GetOption to obtain Option[A], then maps the
//     iso's forward function (A → B) over the Option with OptionMap, yielding Option[B].
//   - ReverseGet: converts B → A via the iso's ReverseGet, then converts A → S via
//     the prism's ReverseGet.
//
// The resulting Prism satisfies both prism laws whenever the original prism
// and the isomorphism individually satisfy their respective laws:
//
//  1. GetOption(ReverseGet(b)) = Some(b)
//  2. If GetOption(s) = Some(b) then ReverseGet(b) round-trips correctly
//
// Type Parameters:
//   - S: The source type (sum type)
//   - A: The intermediate type (focus of the original prism)
//   - B: The target type (focus of the resulting prism, target of the iso)
//
// Parameters:
//   - ab: Iso[A, B] that converts A to B (forward) and B back to A (reverse)
//
// Returns:
//   - A function that takes a Prism[S, A] and returns a Prism[S, B]
//
// See Also:
//   - PrismIMap: similar operation using a pair of plain functions instead of an Iso
//   - PrismComposePrism: for composing two prisms (when the inner step is also partial)
func PrismComposeIso[S, A, B any](ab Iso[A, B]) PrismOperator[S, A, B] {
	return func(pa Prism[S, A]) Prism[S, B] {
		return MakePrismWithName(
			F.Flow2(
				pa.GetOption,
				OptionMap(ab.Get),
			),
			F.Flow2(
				ab.ReverseGet,
				pa.ReverseGet,
			),
			fmt.Sprintf("PrismComposeIso[%s -> %s]", pa, ab),
		)
	}
}
