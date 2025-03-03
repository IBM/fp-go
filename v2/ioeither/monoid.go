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

package ioeither

import (
	"github.com/IBM/fp-go/v2/monoid"
)

type (
	Monoid[E, A any] = monoid.Monoid[IOEither[E, A]]
)

// ApplicativeMonoid returns a [Monoid] that concatenates [IOEither] instances via their applicative
func ApplicativeMonoid[E, A any](
	m monoid.Monoid[A],
) Monoid[E, A] {
	return monoid.ApplicativeMonoid(
		MonadOf[E, A],
		MonadMap[E, A, func(A) A],
		MonadAp[A, E, A],
		m,
	)
}

// ApplicativeMonoid returns a [Monoid] that concatenates [IOEither] instances via their applicative
func ApplicativeMonoidSeq[E, A any](
	m monoid.Monoid[A],
) Monoid[E, A] {
	return monoid.ApplicativeMonoid(
		MonadOf[E, A],
		MonadMap[E, A, func(A) A],
		MonadApSeq[A, E, A],
		m,
	)
}

// ApplicativeMonoid returns a [Monoid] that concatenates [IOEither] instances via their applicative
func ApplicativeMonoidPar[E, A any](
	m monoid.Monoid[A],
) Monoid[E, A] {
	return monoid.ApplicativeMonoid(
		MonadOf[E, A],
		MonadMap[E, A, func(A) A],
		MonadApPar[A, E, A],
		m,
	)
}
