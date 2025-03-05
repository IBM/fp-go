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

package predicate

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/semigroup"
)

type (
	Semigroup[A any] = semigroup.Semigroup[Predicate[A]]
	Monoid[A any]    = monoid.Monoid[Predicate[A]]
)

// SemigroupAny combines predicates via ||
func SemigroupAny[A any]() Semigroup[A] {
	return semigroup.MakeSemigroup(func(first Predicate[A], second Predicate[A]) Predicate[A] {
		return F.Pipe1(
			first,
			Or(second),
		)
	})
}

// SemigroupAll combines predicates via &&
func SemigroupAll[A any]() Semigroup[A] {
	return semigroup.MakeSemigroup(func(first Predicate[A], second Predicate[A]) Predicate[A] {
		return F.Pipe1(
			first,
			And(second),
		)
	})
}

// MonoidAny combines predicates via ||
func MonoidAny[A any]() Monoid[A] {
	return monoid.MakeMonoid(
		SemigroupAny[A]().Concat,
		F.Constant1[A](false),
	)
}

// MonoidAll combines predicates via &&
func MonoidAll[A any]() Monoid[A] {
	return monoid.MakeMonoid(
		SemigroupAll[A]().Concat,
		F.Constant1[A](true),
	)
}
