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

import (
	C "github.com/IBM/fp-go/v2/internal/common"
)

// Compose returns a function that composes an Optional[S, A] with a Lens[A, B]
// to produce an Optional[S, B].
//
// Because a Lens always succeeds, the resulting Optional inherits its
// partiality exclusively from the outer optional: GetOption returns None only
// when the outer optional produces None, and Set is a no-op only in that same
// case.
//
// This is a curried free-function form designed for use with F.Pipe1.  On
// Go 1.27+ the equivalent method form is also available:
//
//	// pipe form (all versions)
//	sb := F.Pipe1(sa, Compose[S, A, B](ab))
//
//	// method form (go1.27+)
//	sb := sa.ComposeLens(ab)
//
// GetOption of the composed optional applies the outer optional's GetOption to
// obtain an Option[A]; when Some, it applies the lens's Get to retrieve B.
// None is returned whenever the outer optional produces None.
//
// Set of the composed optional uses the lens's Set to embed the new B into A,
// then propagates the updated A upward through the outer optional's Set.  When
// the outer optional produces None, Set is a no-op and the original S is
// returned unchanged, consistent with the Optional no-op law.
//
// The composed optional satisfies all three optional laws (GetSet, SetGet,
// SetSet) whenever the outer optional individually satisfies them.
//
// Type Parameters:
//   - S: the outer structure type
//   - A: the intermediate type focused by the optional
//   - B: the focus type of the lens and of the resulting optional
//
// Parameters:
//   - ab: the lens from A to B
//
// Returns:
//   - a function that takes an Optional[S, A] and returns an Optional[S, B]
//
// See Also:
//   - OptionalComposeLens: the underlying implementation in internal/common
//   - github.com/IBM/fp-go/v2/optics/optional.Compose: analogous composition
//     with an Optional instead of a Lens
func Compose[S, A, B any](ab C.Lens[A, B]) func(C.Optional[S, A]) C.Optional[S, B] {
	return C.OptionalComposeLens[S](ab)
}
