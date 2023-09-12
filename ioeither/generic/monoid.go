//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package generic

import (
	ET "github.com/IBM/fp-go/either"
	M "github.com/IBM/fp-go/monoid"
)

func ApplicativeMonoid[GEA ~func() ET.Either[E, A], GEFA ~func() ET.Either[E, func(A) A], E, A any](
	m M.Monoid[A],
) M.Monoid[GEA] {
	return M.ApplicativeMonoid(
		MonadOf[GEA],
		MonadMap[GEA, GEFA],
		MonadAp[GEA, GEFA, GEA],
		m,
	)
}

func ApplicativeMonoidSeq[GEA ~func() ET.Either[E, A], GEFA ~func() ET.Either[E, func(A) A], E, A any](
	m M.Monoid[A],
) M.Monoid[GEA] {
	return M.ApplicativeMonoid(
		MonadOf[GEA],
		MonadMap[GEA, GEFA],
		MonadApSeq[GEA, GEFA, GEA],
		m,
	)
}

func ApplicativeMonoidPar[GEA ~func() ET.Either[E, A], GEFA ~func() ET.Either[E, func(A) A], E, A any](
	m M.Monoid[A],
) M.Monoid[GEA] {
	return M.ApplicativeMonoid(
		MonadOf[GEA],
		MonadMap[GEA, GEFA],
		MonadApPar[GEA, GEFA, GEA],
		m,
	)
}
