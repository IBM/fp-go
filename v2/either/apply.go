// Copyright (c) 2023 IBM Corp.
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

package either

import (
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

func ApplySemigroup[E, A any](s S.Semigroup[A]) S.Semigroup[Either[E, A]] {
	return S.ApplySemigroup(MonadMap[E, A, func(A) A], MonadAp[A, E, A], s)
}

// ApplicativeMonoid returns a [Monoid] that concatenates [Either] instances via their applicative
func ApplicativeMonoid[E, A any](m M.Monoid[A]) M.Monoid[Either[E, A]] {
	return M.ApplicativeMonoid(Of[E, A], MonadMap[E, A, func(A) A], MonadAp[A, E, A], m)
}
