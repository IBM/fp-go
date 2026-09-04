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

import "github.com/IBM/fp-go/v2/internal/common"

// ComposeOptional composes a Prism[S, A] with an Optional[A, B] to produce an Optional[S, B].
//
// This composition first tries to select a variant via the Prism (which may fail), then
// focuses on a sub-part of that variant via the Optional (which may also fail). The result
// is an Optional because either step may return None.
//
// Internally the Prism is widened to an Optional[S, A] via PrismAsOptional and the two
// optionals are then sequentially composed using OptionalComposeOptional.
//
// The resulting Optional satisfies the three standard Optional laws:
//
//  1. GetSet (no-op when Prism misses):
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
//   - B: The target type (the sub-part selected by the inner Optional)
//
// Parameters:
//   - ab: Optional[A, B] that focuses on a part of the variant A
//
// Returns:
//   - A function that takes a Prism[S, A] and returns an Optional[S, B]
//
// See Also:
//   - ComposeLens: for the simpler case where the inner step is a total Lens
//   - Compose: for composing two Prisms (when the inner step is also partial and reversible)
func ComposeOptional[S, A, B any](ab common.Optional[A, B]) func(Prism[S, A]) common.Optional[S, B] {
	return common.PrismComposeOptional[S](ab)
}
