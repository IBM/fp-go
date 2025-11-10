//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache LicensVersion 2.0 (the "License");
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
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/monoid"
)

// ApplicativeMonoid returns a [Monoid] that concatenates [IOResult] instances via their applicative
//
//go:inline
func ApplicativeMonoid[A any](
	m monoid.Monoid[A],
) Monoid[A] {
	return ioeither.ApplicativeMonoid[error](m)
}

// ApplicativeMonoid returns a [Monoid] that concatenates [IOResult] instances via their applicative
//
//go:inline
func ApplicativeMonoidSeq[A any](
	m monoid.Monoid[A],
) Monoid[A] {
	return ioeither.ApplicativeMonoidSeq[error](m)
}

// ApplicativeMonoid returns a [Monoid] that concatenates [IOResult] instances via their applicative
//
//go:inline
func ApplicativeMonoidPar[A any](
	m monoid.Monoid[A],
) Monoid[A] {
	return ioeither.ApplicativeMonoidPar[error](m)
}
