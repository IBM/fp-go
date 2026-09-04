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

// ComposeIso composes a Prism[S, A] with an Iso[A, B] to produce a Prism[S, B].
//
// Because an Iso is a total and lossless bidirectional transformation, composing a
// Prism with an Iso always yields another Prism — the result is never made partial
// by the Iso itself.  GetOption returns None only when the original Prism does not
// match the source value.
//
// The composition works as follows:
//   - GetOption: applies the prism's GetOption to get Option[A], then maps the
//     iso's forward function (A → B) over the Option, yielding Option[B].
//   - ReverseGet: converts B → A via the iso's ReverseGet, then A → S via the
//     prism's ReverseGet.
//
// The resulting Prism satisfies both prism laws whenever the original prism and the
// isomorphism individually satisfy their respective laws:
//
//  1. GetOption(ReverseGet(b)) = Some(b)
//  2. If GetOption(s) = Some(b) then ReverseGet(b) round-trips back to an
//     equivalent s
//
// Type Parameters:
//   - S: The source type (sum type)
//   - A: The intermediate type (focus of the original Prism)
//   - B: The target type (focus of the resulting Prism, target of the Iso)
//
// Parameters:
//   - ab: Iso[A, B] that converts A to B (forward) and B back to A (reverse)
//
// Returns:
//   - An Operator[S, A, B] — a function that takes a Prism[S, A] and returns a Prism[S, B]
//
// See Also:
//   - IMap: similar operation using a plain pair of functions instead of an Iso
//   - Compose: for composing two Prisms (when the inner step is also partial)
func ComposeIso[S, A, B any](ab common.Iso[A, B]) Operator[S, A, B] {
	return common.PrismComposeIso[S](ab)
}
