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

package reader

import (
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

func ApplySemigroup[R, A any](
	_map func(func(R) A, func(A) func(A) A) func(R, func(A) A),
	_ap func(func(R, func(A) A), func(R) A) func(R) A,

	s S.Semigroup[A],
) S.Semigroup[func(R) A] {
	return S.ApplySemigroup(_map, _ap, s)
}

func ApplicativeMonoid[R, A any](
	_of func(A) func(R) A,
	_map func(func(R) A, func(A) func(A) A) func(R, func(A) A),
	_ap func(func(R, func(A) A), func(R) A) func(R) A,

	m M.Monoid[A],
) M.Monoid[func(R) A] {
	return M.ApplicativeMonoid(_of, _map, _ap, m)
}
