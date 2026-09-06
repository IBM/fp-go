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

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/optics/prism"
)

// Compose creates a Kleisli arrow that composes a prism with an isomorphism.
//
// Deprecated: Use github.com/IBM/fp-go/v2/optics/iso.ComposePrism instead.
func Compose[S, A, B any](ab Prism[A, B]) P.Kleisli[S, Iso[S, A], B] {
	return func(ia Iso[S, A]) Prism[S, B] {
		return P.MakePrismWithName(
			F.Flow2(ia.Get, ab.GetOption),
			F.Flow2(ab.ReverseGet, ia.ReverseGet),
			fmt.Sprintf("IsoCompose[%s -> %s]", ia, ab),
		)
	}
}
