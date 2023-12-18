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

// Iso is an optic which converts elements of type `S` into elements of type `A` without loss.
package iso

import (
	EM "github.com/IBM/fp-go/endomorphism"
	F "github.com/IBM/fp-go/function"
)

type Iso[S, A any] struct {
	Get        func(s S) A
	ReverseGet func(a A) S
}

func MakeIso[S, A any](get func(S) A, reverse func(A) S) Iso[S, A] {
	return Iso[S, A]{Get: get, ReverseGet: reverse}
}

// Id returns an iso implementing the identity operation
func Id[S any]() Iso[S, S] {
	return MakeIso(F.Identity[S], F.Identity[S])
}

// Compose combines an ISO with another ISO
func Compose[S, A, B any](ab Iso[A, B]) func(Iso[S, A]) Iso[S, B] {
	return func(sa Iso[S, A]) Iso[S, B] {
		return MakeIso(
			F.Flow2(sa.Get, ab.Get),
			F.Flow2(ab.ReverseGet, sa.ReverseGet),
		)
	}
}

// Reverse changes the order of parameters for an iso
func Reverse[S, A any](sa Iso[S, A]) Iso[A, S] {
	return MakeIso(
		sa.ReverseGet,
		sa.Get,
	)
}

func modify[FCT ~func(A) A, S, A any](f FCT, sa Iso[S, A], s S) S {
	return F.Pipe3(
		s,
		sa.Get,
		f,
		sa.ReverseGet,
	)
}

// Modify applies a transformation
func Modify[S any, FCT ~func(A) A, A any](f FCT) func(Iso[S, A]) EM.Endomorphism[S] {
	return EM.Curry3(modify[FCT, S, A])(f)
}

// Wrap wraps the value
func Unwrap[A, S any](s S) func(Iso[S, A]) A {
	return func(sa Iso[S, A]) A {
		return sa.Get(s)
	}
}

// Unwrap unwraps the value
func Wrap[S, A any](a A) func(Iso[S, A]) S {
	return func(sa Iso[S, A]) S {
		return sa.ReverseGet(a)
	}
}

// From wraps the value
func To[A, S any](s S) func(Iso[S, A]) A {
	return Unwrap[A, S](s)
}

// To unwraps the value
func From[S, A any](a A) func(Iso[S, A]) S {
	return Wrap[S](a)
}

func imap[S, A, B any](sa Iso[S, A], ab func(A) B, ba func(B) A) Iso[S, B] {
	return MakeIso(
		F.Flow2(sa.Get, ab),
		F.Flow2(ba, sa.ReverseGet),
	)
}

// IMap implements a bidirectional mapping of the transform
func IMap[S, A, B any](ab func(A) B, ba func(B) A) func(Iso[S, A]) Iso[S, B] {
	return func(sa Iso[S, A]) Iso[S, B] {
		return imap(sa, ab, ba)
	}
}
