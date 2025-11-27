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
	F "github.com/IBM/fp-go/v2/function"
	OPT "github.com/IBM/fp-go/v2/optics/optional"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
)

// AsOptional converts a prism into an optional
func AsOptional[S, A any](sa P.Prism[S, A]) OPT.Optional[S, A] {
	return OPT.MakeOptional(
		sa.GetOption,
		func(s S, a A) S {
			return P.Set[S](a)(sa)(s)
		},
	)
}

func PrismSome[A any]() P.Prism[O.Option[A], A] {
	return P.MakePrismWithName(F.Identity[O.Option[A]], O.Some[A], "PrismSome")
}

// Some returns a `Optional` from a `Optional` focused on the `Some` of a `Option` type.
func Some[S, A any](soa OPT.Optional[S, O.Option[A]]) OPT.Optional[S, A] {
	return OPT.Compose[S](AsOptional(PrismSome[A]()))(soa)
}
