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

// PrismComposeOptional composes a Prism[S, A] with an Optional[A, B] to produce
// an Optional[S, B].
//
// This composition lets you first select a variant within a sum type (via the
// Prism, which may return None), and then optionally focus on a sub-part of that
// variant (via the inner Optional, which may also return None). The result is an
// Optional because either step may fail.
//
// Internally, the Prism is widened to an Optional[S, A] via PrismAsOptional, and
// the two optionals are then sequentially composed using OptionalComposeOptional.
//
// The resulting Optional satisfies the three standard Optional laws whenever the
// constituent prism and inner optional each satisfy their own laws:
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
//   - l: Optional[A, B] that focuses on a sub-part of the intermediate type A
//
// Returns:
//   - A function that takes a Prism[S, A] and returns an Optional[S, B]
//
// See Also:
//   - PrismComposeLens: simpler variant where the inner step is a total Lens
//   - PrismComposePrism: variant where the inner step is a Prism (yields a Prism, not Optional)
//   - OptionalComposeOptional: the underlying combinator used in the implementation
func PrismComposeOptional[S, A, B any](l Optional[A, B]) func(Prism[S, A]) Optional[S, B] {
	return F.Flow2(
		PrismAsOptional,
		OptionalComposeOptional[S](l),
	)
}
