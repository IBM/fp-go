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

package endomorphism

import (
	"github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Of converts any function to an [Endomorphism]
func Of[F ~func(A) A, A any](f F) Endomorphism[A] {
	return f
}

// Wrap converts any function to an [Endomorphism]
// Deprecated: no need to use it
func Wrap[F ~func(A) A, A any](f F) Endomorphism[A] {
	return f
}

// Unwrap converts any [Endomorphism] to a function
// Deprecated: no need to use it
func Unwrap[F ~func(A) A, A any](f Endomorphism[A]) F {
	return f
}

// Identity returns the identity [Endomorphism]
func Identity[A any]() Endomorphism[A] {
	return function.Identity[A]
}

// Semigroup for the Endomorphism where the `concat` operation is the usual function composition.
func Semigroup[A any]() S.Semigroup[Endomorphism[A]] {
	return S.MakeSemigroup(Compose[A])
}

// Monoid for the Endomorphism where the `concat` operation is the usual function composition.
func Monoid[A any]() M.Monoid[Endomorphism[A]] {
	return M.MakeMonoid(Compose[A], Identity[A]())
}
