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

package lens

import "github.com/IBM/fp-go/v2/internal/common"

// ComposePrism composes a Lens with a Prism to produce an Optional.
//
// A Lens always succeeds at focusing (S → A), while a Prism only succeeds for
// one variant of a sum type (A → B).  Their composition yields an Optional that
// can either find a B inside S or return nothing when the Prism does not match.
//
// The resulting Optional satisfies two laws:
//
//  1. SetGet (match case): if optional.GetOption(s) = Some(b),
//     then optional.GetOption(optional.Set(b)(s)) = Some(b).
//
//  2. GetSet (no-match case): if optional.GetOption(s) = None,
//     then optional.Set(b)(s) = s   (Set is a no-op).
//
// Type Parameters:
//   - S: The outer structure type (focused by the Lens)
//   - A: The intermediate type (the Lens focus; the Prism source)
//   - B: The target type (the Prism focus)
//
// Parameters:
//   - p: A Prism[A, B] that optionally extracts B from A
//
// Returns:
//   - A function that takes a Lens[S, A] and returns an Optional[S, B]
//
// See Also:
//   - ComposeIso: compose a Lens with an Iso instead of a Prism
//   - Compose: compose two Lenses together
func ComposePrism[S, A, B any](p common.Prism[A, B]) func(Lens[S, A]) common.Optional[S, B] {
	return common.LensComposePrism[S](p)
}
