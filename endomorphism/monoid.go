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
	F "github.com/IBM/fp-go/function"
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

func concat[A any](first, second func(A) A) func(A) A {
	return F.Flow2(first, second)
}

// Semigroup for the Endomorphism where the `concat` operation is the usual function composition.
func Semigroup[A any]() S.Semigroup[func(A) A] {
	return S.MakeSemigroup(concat[A])
}

// Monoid for the Endomorphism where the `concat` operation is the usual function composition.
func Monoid[A any]() M.Monoid[func(A) A] {
	return M.MakeMonoid(concat[A], F.Identity[A])
}
