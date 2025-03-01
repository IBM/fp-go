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

package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Of converts any function to an [Endomorphism]
func Of[ENDO ~func(A) A, F ~func(A) A, A any](f F) ENDO {
	return func(a A) A {
		return f(a)
	}
}

// Wrap converts any function to an [Endomorphism]
func Wrap[ENDO ~func(A) A, F ~func(A) A, A any](f F) ENDO {
	return Of[ENDO](f)
}

// Unwrap converts any [Endomorphism] to a normal function
func Unwrap[F ~func(A) A, ENDO ~func(A) A, A any](f ENDO) F {
	return Of[F](f)
}

func Identity[ENDO ~func(A) A, A any]() ENDO {
	return Of[ENDO](F.Identity[A])
}

func Compose[ENDO ~func(A) A, A any](f1, f2 ENDO) ENDO {
	return Of[ENDO](F.Flow2(f1, f2))
}

// Semigroup for the Endomorphism where the `concat` operation is the usual function composition.
func Semigroup[ENDO ~func(A) A, A any]() S.Semigroup[ENDO] {
	return S.MakeSemigroup(Compose[ENDO])
}

// Monoid for the Endomorphism where the `concat` operation is the usual function composition.
func Monoid[ENDO ~func(A) A, A any]() M.Monoid[ENDO] {
	return M.MakeMonoid(Compose[ENDO], Identity[ENDO]())
}
