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

package readerioeither

import (
	"github.com/IBM/fp-go/v2/monoid"
)

type (
	Monoid[A any] = monoid.Monoid[ReaderIOEither[A]]
)

// ApplicativeMonoid returns a [Monoid] that concatenates [ReaderIOEither] instances via their applicative
func ApplicativeMonoid[A any](m monoid.Monoid[A]) Monoid[A] {
	return monoid.ApplicativeMonoid(
		Of[A],
		MonadMap[A, func(A) A],
		MonadAp[A, A],
		m,
	)
}

// ApplicativeMonoidSeq returns a [Monoid] that concatenates [ReaderIOEither] instances via their applicative
func ApplicativeMonoidSeq[A any](m monoid.Monoid[A]) Monoid[A] {
	return monoid.ApplicativeMonoid(
		Of[A],
		MonadMap[A, func(A) A],
		MonadApSeq[A, A],
		m,
	)
}

// ApplicativeMonoidPar returns a [Monoid] that concatenates [ReaderIOEither] instances via their applicative
func ApplicativeMonoidPar[A any](m monoid.Monoid[A]) Monoid[A] {
	return monoid.ApplicativeMonoid(
		Of[A],
		MonadMap[A, func(A) A],
		MonadApPar[A, A],
		m,
	)
}

// AlternativeMonoid is the alternative [Monoid] for [ReaderIOEither]
func AlternativeMonoid[A any](m monoid.Monoid[A]) Monoid[A] {
	return monoid.AlternativeMonoid(
		Of[A],
		MonadMap[A, func(A) A],
		MonadAp[A, A],
		MonadAlt[A],
		m,
	)
}

// AltMonoid is the alternative [Monoid] for an [ReaderIOEither]
func AltMonoid[A any](zero Lazy[ReaderIOEither[A]]) Monoid[A] {
	return monoid.AltMonoid(
		zero,
		MonadAlt[A],
	)
}
