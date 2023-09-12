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
	"context"

	ET "github.com/IBM/fp-go/either"
	M "github.com/IBM/fp-go/monoid"
)

func ApplicativeMonoid[GRA ~func(context.Context) GIOA, GRFA ~func(context.Context) GIOFA, GIOA ~func() ET.Either[error, A], GIOFA ~func() ET.Either[error, func(A) A], A any](
	m M.Monoid[A],
) M.Monoid[GRA] {
	return M.ApplicativeMonoid(
		Of[GRA],
		MonadMap[GRA, GRFA],
		MonadAp[GRA, GRA, GRFA],
		m,
	)
}

func ApplicativeMonoidSeq[GRA ~func(context.Context) GIOA, GRFA ~func(context.Context) GIOFA, GIOA ~func() ET.Either[error, A], GIOFA ~func() ET.Either[error, func(A) A], A any](
	m M.Monoid[A],
) M.Monoid[GRA] {
	return M.ApplicativeMonoid(
		Of[GRA],
		MonadMap[GRA, GRFA],
		MonadApSeq[GRA, GRA, GRFA],
		m,
	)
}

func ApplicativeMonoidPar[GRA ~func(context.Context) GIOA, GRFA ~func(context.Context) GIOFA, GIOA ~func() ET.Either[error, A], GIOFA ~func() ET.Either[error, func(A) A], A any](
	m M.Monoid[A],
) M.Monoid[GRA] {
	return M.ApplicativeMonoid(
		Of[GRA],
		MonadMap[GRA, GRFA],
		MonadApPar[GRA, GRA, GRFA],
		m,
	)
}
