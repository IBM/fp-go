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

import F "github.com/IBM/fp-go/v2/function"

// PrismComposeLens composes a Prism[S, A] with a Lens[A, B] to produce an Optional[S, B].
//
// This composition lets you first select a variant within a sum type (via the Prism),
// and then focus on a field within that variant (via the Lens). The result is an
// Optional because the Prism may not match the source value.
//
// Internally the Prism is widened to an Optional[S, A] via PrismAsOptional, and that
// Optional is then composed with the Lens (also widened to an Optional[A, B]) using
// OptionalComposeLens.
//
// The resulting Optional satisfies the standard Optional laws:
//
//  1. GetSet (no-op on None):
//     GetOption(s) = None => Set(b)(s) = s
//
//  2. SetGet (get what you set):
//     GetOption(s) = Some(_) => GetOption(Set(b)(s)) = Some(b)
//
//  3. SetSet (last set wins):
//     Set(c)(Set(b)(s)) = Set(c)(s)
//
// Type Parameters:
//   - S: The source type (sum type)
//   - A: The intermediate type (the variant selected by the Prism)
//   - B: The target type (the field selected by the Lens within A)
//
// Parameters:
//   - l: Lens[A, B] that focuses on field B within variant A
//
// Returns:
//   - A function that takes a Prism[S, A] and returns an Optional[S, B]
//
// See Also:
//   - PrismAsOptional: widens a Prism to an Optional
//   - OptionalComposeLens: composes an Optional with a Lens
//   - LensComposePrism: the inverse composition (Lens then Prism)
func PrismComposeLens[S, A, B any](l Lens[A, B]) func(Prism[S, A]) Optional[S, B] {
	return F.Flow2(
		PrismAsOptional,
		OptionalComposeLens[S](l),
	)
}
