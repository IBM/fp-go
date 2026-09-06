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

package iso

import "github.com/IBM/fp-go/v2/internal/common"

// ComposePrism composes an Iso[S, A] with a Prism[A, B] to produce a Prism[S, B].
//
// Because the Iso is a total bidirectional transformation, the source type of the
// Prism is widened from A to S.  The resulting Prism:
//   - GetOption: converts S → A via the Iso's Get, then applies the Prism's GetOption
//     to obtain Option[B].
//   - ReverseGet: constructs A from B via the Prism's ReverseGet, then converts A → S
//     via the Iso's ReverseGet.
//
// The resulting Prism satisfies both prism laws whenever the original Prism and the
// Iso individually satisfy their respective laws:
//
//  1. GetOption(ReverseGet(b)) = Some(b)
//  2. If GetOption(s) = Some(b) then ReverseGet(b) round-trips correctly
//
// Type Parameters:
//   - S: The new source type (introduced by the Iso)
//   - A: The intermediate type (target of the Iso; source of the Prism)
//   - B: The focus type of the resulting Prism
//
// Parameters:
//   - ab: Prism[A, B] that optionally extracts B from A
//
// Returns:
//   - A function that takes an Iso[S, A] and returns a Prism[S, B]
//
// See Also:
//   - Compose: for composing two Isos (when the inner step is also total)
func ComposePrism[S, A, B any](ab common.Prism[A, B]) func(Iso[S, A]) common.Prism[S, B] {
	return common.IsoComposePrism[S](ab)
}
