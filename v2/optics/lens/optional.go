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

// ComposeOptional composes a Lens with an Optional to produce an Optional.
//
// A Lens always succeeds at focusing (S → A), while an Optional only succeeds
// when the inner value exists (A → B).  Their composition yields an Optional
// that can either find a B inside S or return nothing when the inner Optional
// does not match.
//
// The lens is first widened to an Optional[S, A] via LensAsOptional, then
// composed with ab using OptionalComposeOptional.
//
// The resulting Optional satisfies the three standard Optional laws:
//
//  1. GetSet (no-match is a no-op): if optional.GetOption(s) = None,
//     then optional.Set(b)(s) = s.
//
//  2. SetGet (match case): if optional.GetOption(s) = Some(_),
//     then optional.GetOption(optional.Set(b)(s)) = Some(b).
//
//  3. SetSet (last Set wins): optional.Set(b2)(optional.Set(b1)(s)) = optional.Set(b2)(s).
//
// Type Parameters:
//   - S: The outer structure type (focused by the Lens)
//   - A: The intermediate type (the Lens focus; the Optional source)
//   - B: The target type (the Optional focus)
//
// Parameters:
//   - ab: An Optional[A, B] that partially focuses on B within A
//
// Returns:
//   - A function that takes a Lens[S, A] and returns an Optional[S, B]
//
// See Also:
//   - ComposePrism: compose a Lens with a Prism instead of an Optional
//   - ComposeIso: compose a Lens with an Iso (always-succeeding)
//   - Compose: compose two Lenses together
func ComposeOptional[S, A, B any](ab common.Optional[A, B]) func(Lens[S, A]) common.Optional[S, B] {
	return common.LensComposeOptional[S](ab)
}
