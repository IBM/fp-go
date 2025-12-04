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

package optional

import (
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	OPT "github.com/IBM/fp-go/v2/optics/optional"
	O "github.com/IBM/fp-go/v2/option"
)

func lensAsOptional[S, A any](creator func(get O.Kleisli[S, A], set func(S, A) S) OPT.Optional[S, A], sa L.Lens[S, A]) OPT.Optional[S, A] {
	return creator(F.Flow2(sa.Get, O.Some[A]), func(s S, a A) S {
		return sa.Set(a)(s)
	})
}

// LensAsOptional converts a Lens into an Optional
func LensAsOptional[S, A any](sa L.Lens[S, A]) OPT.Optional[S, A] {
	return lensAsOptional(OPT.MakeOptional[S, A], sa)
}
