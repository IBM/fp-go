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

package readerreaderioresult

import (
	"github.com/IBM/fp-go/v2/monoid"
)

type (
	Monoid[R, A any] = monoid.Monoid[ReaderReaderIOResult[R, A]]
)

func ApplicativeMonoid[R, A any](m monoid.Monoid[A]) Monoid[R, A] {
	return monoid.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[R, A, A],
		m,
	)
}

func ApplicativeMonoidSeq[R, A any](m monoid.Monoid[A]) Monoid[R, A] {
	return monoid.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadApSeq[R, A, A],
		m,
	)
}

func ApplicativeMonoidPar[R, A any](m monoid.Monoid[A]) Monoid[R, A] {
	return monoid.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadApPar[R, A, A],
		m,
	)
}

func AlternativeMonoid[R, A any](m monoid.Monoid[A]) Monoid[R, A] {
	return monoid.AlternativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[R, A, A],
		MonadAlt[R, A],
		m,
	)
}

func AltMonoid[R, A any](zero Lazy[ReaderReaderIOResult[R, A]]) Monoid[R, A] {
	return monoid.AltMonoid(
		zero,
		MonadAlt[R, A],
	)
}
