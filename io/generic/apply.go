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

package generic

import (
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

func ApplySemigroup[GA ~func() A, A any](s S.Semigroup[A]) S.Semigroup[GA] {
	return S.ApplySemigroup(MonadMap[GA, func() func(A) A, A, func(A) A], MonadAp[GA, GA, func() func(A) A, A, A], s)
}

func ApplicativeMonoid[GA ~func() A, A any](m M.Monoid[A]) M.Monoid[GA] {
	return M.ApplicativeMonoid(Of[GA, A], MonadMap[GA, func() func(A) A, A, func(A) A], MonadAp[GA, GA, func() func(A) A, A, A], m)
}
