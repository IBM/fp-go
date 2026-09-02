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
	F "github.com/IBM/fp-go/v2/function"
)

func lensAsOptional[S, A any](creator func(get OptionKleisli[S, A], set func(S, A) S) Optional[S, A], sa Lens[S, A]) Optional[S, A] {
	return creator(F.Flow2(sa.Get, OptionSome[A]), func(s S, a A) S {
		return sa.Set(a)(s)
	})
}

// LensAsOptional converts a Lens into an Optional
func LensAsOptional[S, A any](sa Lens[S, A]) Optional[S, A] {
	return lensAsOptional(MakeOptional[S, A], sa)
}

// LensComposeOptional composes a lens with an optional
func OptionalComposeLens[S, A, B any](ab Lens[A, B]) OptionalOperator[S, A, B] {
	return F.Pipe2(
		ab,
		LensAsOptional[A, B],
		OptionalComposeOptional[S, A, B],
	)
}
