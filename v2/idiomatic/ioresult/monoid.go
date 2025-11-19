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

package ioresult

import (
	"github.com/IBM/fp-go/v2/monoid"
)

type (
	Monoid[A any] = monoid.Monoid[IOResult[A]]
)

// ApplicativeMonoid returns a [Monoid] that concatenates [IOEither] instances via their applicative
// ApplicativeMonoid returns a Monoid that concatenates IOResult instances via their applicative.
// Uses parallel execution (default Ap behavior).
func ApplicativeMonoid[A any](
	m monoid.Monoid[A],
) Monoid[A] {
	return monoid.ApplicativeMonoid(
		MonadOf[A],
		MonadMap[A, func(A) A],
		MonadAp[A, A],
		m,
	)
}

// ApplicativeMonoid returns a [Monoid] that concatenates [IOEither] instances via their applicative
// ApplicativeMonoidSeq returns a Monoid that concatenates IOResult instances sequentially.
// Uses sequential execution (ApSeq).
func ApplicativeMonoidSeq[A any](
	m monoid.Monoid[A],
) Monoid[A] {
	return monoid.ApplicativeMonoid(
		MonadOf[A],
		MonadMap[A, func(A) A],
		MonadApSeq[A, A],
		m,
	)
}

// ApplicativeMonoid returns a [Monoid] that concatenates [IOEither] instances via their applicative
// ApplicativeMonoidPar returns a Monoid that concatenates IOResult instances in parallel.
// Uses parallel execution (ApPar) explicitly.
func ApplicativeMonoidPar[A any](
	m monoid.Monoid[A],
) Monoid[A] {
	return monoid.ApplicativeMonoid(
		MonadOf[A],
		MonadMap[A, func(A) A],
		MonadApPar[A, A],
		m,
	)
}
