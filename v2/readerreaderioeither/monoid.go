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

package readerreaderioeither

import (
	"github.com/IBM/fp-go/v2/monoid"
)

type (
	Monoid[R, C, E, A any] = monoid.Monoid[ReaderReaderIOEither[R, C, E, A]]
)

func ApplicativeMonoid[R, C, E, A any](m monoid.Monoid[A]) Monoid[R, C, E, A] {
	return monoid.ApplicativeMonoid(
		Of[R, C, E, A],
		MonadMap[R, C, E, A, func(A) A],
		MonadAp[R, C, E, A, A],
		m,
	)
}

func ApplicativeMonoidSeq[R, C, E, A any](m monoid.Monoid[A]) Monoid[R, C, E, A] {
	return monoid.ApplicativeMonoid(
		Of[R, C, E, A],
		MonadMap[R, C, E, A, func(A) A],
		MonadApSeq[R, C, E, A, A],
		m,
	)
}

func ApplicativeMonoidPar[R, C, E, A any](m monoid.Monoid[A]) Monoid[R, C, E, A] {
	return monoid.ApplicativeMonoid(
		Of[R, C, E, A],
		MonadMap[R, C, E, A, func(A) A],
		MonadApPar[R, C, E, A, A],
		m,
	)
}

func AlternativeMonoid[R, C, E, A any](m monoid.Monoid[A]) Monoid[R, C, E, A] {
	return monoid.AlternativeMonoid(
		Of[R, C, E, A],
		MonadMap[R, C, E, A, func(A) A],
		MonadAp[R, C, E, A, A],
		MonadAlt[R, C, E, A],
		m,
	)
}

func AltMonoid[R, C, E, A any](zero Lazy[ReaderReaderIOEither[R, C, E, A]]) Monoid[R, C, E, A] {
	return monoid.AltMonoid(
		zero,
		MonadAlt[R, C, E, A],
	)
}
