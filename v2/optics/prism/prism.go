// Copyright (c) 2023 - 2025 IBM Corp.
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

// Prism is an optic used to select part of a sum type.
package prism

import (
	EM "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

type (
	// Prism is an optic used to select part of a sum type.
	Prism[S, A any] interface {
		GetOption(s S) O.Option[A]
		ReverseGet(a A) S
	}

	prismImpl[S, A any] struct {
		get func(S) O.Option[A]
		rev func(A) S
	}
)

func (prism prismImpl[S, A]) GetOption(s S) O.Option[A] {
	return prism.get(s)
}

func (prism prismImpl[S, A]) ReverseGet(a A) S {
	return prism.rev(a)
}

func MakePrism[S, A any](get func(S) O.Option[A], rev func(A) S) Prism[S, A] {
	return prismImpl[S, A]{get, rev}
}

// Id returns a prism implementing the identity operation
func Id[S any]() Prism[S, S] {
	return MakePrism(O.Some[S], F.Identity[S])
}

func FromPredicate[S any](pred func(S) bool) Prism[S, S] {
	return MakePrism(O.FromPredicate(pred), F.Identity[S])
}

// Compose composes a `Prism` with a `Prism`.
func Compose[S, A, B any](ab Prism[A, B]) func(Prism[S, A]) Prism[S, B] {
	return func(sa Prism[S, A]) Prism[S, B] {
		return MakePrism(F.Flow2(
			sa.GetOption,
			O.Chain(ab.GetOption),
		), F.Flow2(
			ab.ReverseGet,
			sa.ReverseGet,
		))
	}
}

func prismModifyOption[S, A any](f func(A) A, sa Prism[S, A], s S) O.Option[S] {
	return F.Pipe2(
		s,
		sa.GetOption,
		O.Map(F.Flow2(
			f,
			sa.ReverseGet,
		)),
	)
}

func prismModify[S, A any](f func(A) A, sa Prism[S, A], s S) S {
	return F.Pipe1(
		prismModifyOption(f, sa, s),
		O.GetOrElse(F.Constant(s)),
	)
}

func prismSet[S, A any](a A) func(Prism[S, A]) EM.Endomorphism[S] {
	return EM.Curry3(prismModify[S, A])(F.Constant1[A](a))
}

func Set[S, A any](a A) func(Prism[S, A]) EM.Endomorphism[S] {
	return EM.Curry3(prismModify[S, A])(F.Constant1[A](a))
}

func prismSome[A any]() Prism[O.Option[A], A] {
	return MakePrism(F.Identity[O.Option[A]], O.Some[A])
}

// Some returns a `Prism` from a `Prism` focused on the `Some` of a `Option` type.
func Some[S, A any](soa Prism[S, O.Option[A]]) Prism[S, A] {
	return Compose[S](prismSome[A]())(soa)
}

func imap[S any, AB ~func(A) B, BA ~func(B) A, A, B any](sa Prism[S, A], ab AB, ba BA) Prism[S, B] {
	return MakePrism(
		F.Flow2(sa.GetOption, O.Map(ab)),
		F.Flow2(ba, sa.ReverseGet),
	)
}

func IMap[S any, AB ~func(A) B, BA ~func(B) A, A, B any](ab AB, ba BA) func(Prism[S, A]) Prism[S, B] {
	return func(sa Prism[S, A]) Prism[S, B] {
		return imap(sa, ab, ba)
	}
}
