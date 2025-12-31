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

package eithert

import (
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	FC "github.com/IBM/fp-go/v2/internal/functor"
)

func MonadAlt[LAZY ~func() HKTFA, E, A, HKTFA any](
	fof func(ET.Either[E, A]) HKTFA,
	fchain func(HKTFA, func(ET.Either[E, A]) HKTFA) HKTFA,

	first HKTFA,
	second LAZY) HKTFA {

	return fchain(first, ET.Fold(F.Ignore1of1[E](second), F.Flow2(ET.Of[E, A], fof)))
}

func Alt[LAZY ~func() HKTFA, E, A, HKTFA any](
	fof func(ET.Either[E, A]) HKTFA,
	fchain func(func(ET.Either[E, A]) HKTFA) func(HKTFA) HKTFA,

	second LAZY) func(HKTFA) HKTFA {

	return fchain(ET.Fold(F.Ignore1of1[E](second), F.Flow2(ET.Of[E, A], fof)))
}

// HKTFA = HKT<F, Either<E, A>>
// HKTFB = HKT<F, Either<E, B>>
func MonadMap[E, A, B, HKTFA, HKTFB any](fmap func(HKTFA, func(ET.Either[E, A]) ET.Either[E, B]) HKTFB, fa HKTFA, f func(A) B) HKTFB {
	// HKTGA = Either[E, A]
	// HKTGB = Either[E, B]
	return FC.MonadMap(fmap, ET.MonadMap[E, A, B], fa, f)
}

func Map[E, A, B, HKTFA, HKTFB any](
	fmap func(func(ET.Either[E, A]) ET.Either[E, B]) func(HKTFA) HKTFB,
	f func(A) B) func(HKTFA) HKTFB {
	// HKTGA = Either[E, A]
	// HKTGB = Either[E, B]
	return FC.Map(fmap, ET.Map[E, A, B], f)
}

// HKTFA = HKT<F, Either<E, A>>
// HKTFB = HKT<F, Either<E, B>>
func MonadBiMap[E1, E2, A, B, HKTFA, HKTFB any](fmap func(HKTFA, func(ET.Either[E1, A]) ET.Either[E2, B]) HKTFB, fa HKTFA, f func(E1) E2, g func(A) B) HKTFB {
	// HKTGA = Either[E, A]
	// HKTGB = Either[E, B]
	return fmap(fa, ET.BiMap(f, g))
}

// HKTFA = HKT<F, Either<E, A>>
// HKTFB = HKT<F, Either<E, B>>
func BiMap[E1, E2, A, B, HKTFA, HKTFB any](
	fmap func(func(ET.Either[E1, A]) ET.Either[E2, B]) func(HKTFA) HKTFB,
	f func(E1) E2, g func(A) B) func(HKTFA) HKTFB {
	// HKTGA = Either[E, A]
	// HKTGB = Either[E, B]
	return fmap(ET.BiMap(f, g))
}

// HKTFA = HKT<F, Either<E, A>>
// HKTFB = HKT<F, Either<E, B>>
func MonadChain[E, A, B, HKTFA, HKTFB any](
	fchain func(HKTFA, func(ET.Either[E, A]) HKTFB) HKTFB,
	fof func(ET.Either[E, B]) HKTFB,
	ma HKTFA,
	f func(A) HKTFB) HKTFB {
	// dispatch to the even more generic implementation
	return fchain(ma, ET.Fold(F.Flow2(ET.Left[B, E], fof), f))
}

func Chain[E, A, B, HKTFA, HKTFB any](
	fchain func(func(ET.Either[E, A]) HKTFB) func(HKTFA) HKTFB,
	fof func(ET.Either[E, B]) HKTFB,
	f func(A) HKTFB) func(HKTFA) HKTFB {
	// dispatch to the even more generic implementation
	return fchain(ET.Fold(F.Flow2(ET.Left[B, E], fof), f))
}

func MonadAp[E, A, B, HKTFAB, HKTFGAB, HKTFA, HKTFB any](
	fap func(HKTFGAB, HKTFA) HKTFB,
	fmap func(HKTFAB, func(ET.Either[E, func(A) B]) func(ET.Either[E, A]) ET.Either[E, B]) HKTFGAB,
	fab HKTFAB,
	fa HKTFA) HKTFB {
	return apply.MonadAp(fap, fmap, ET.MonadAp[B, E, A], fab, fa)
}

func Ap[E, A, B, HKTFAB, HKTFGAB, HKTFA, HKTFB any](
	fap func(HKTFA) func(HKTFGAB) HKTFB,
	fmap func(func(ET.Either[E, func(A) B]) func(ET.Either[E, A]) ET.Either[E, B]) func(HKTFAB) HKTFGAB,
	fa HKTFA) func(HKTFAB) HKTFB {
	return apply.Ap(fap, fmap, ET.Ap[B, E, A], fa)
}

func Right[E, A, HKTA any](fof func(ET.Either[E, A]) HKTA, a A) HKTA {
	return F.Pipe2(a, ET.Right[E, A], fof)
}

func Left[E, A, HKTA any](fof func(ET.Either[E, A]) HKTA, e E) HKTA {
	return F.Pipe2(e, ET.Left[A, E], fof)
}

// HKTA  = HKT[A]
// HKTEA = HKT[Either[E, A]]
func RightF[E, A, HKTA, HKTEA any](fmap func(HKTA, func(A) ET.Either[E, A]) HKTEA, fa HKTA) HKTEA {
	return fmap(fa, ET.Right[E, A])
}

// HKTE  = HKT[E]
// HKTEA = HKT[Either[E, A]]
func LeftF[E, A, HKTE, HKTEA any](fmap func(HKTE, func(E) ET.Either[E, A]) HKTEA, fe HKTE) HKTEA {
	return fmap(fe, ET.Left[A, E])
}

func FoldE[E, A, HKTEA, HKTB any](mchain func(HKTEA, func(ET.Either[E, A]) HKTB) HKTB, ma HKTEA, onLeft func(E) HKTB, onRight func(A) HKTB) HKTB {
	return mchain(ma, ET.Fold(onLeft, onRight))
}

func MatchE[E, A, HKTEA, HKTB any](mchain func(HKTEA, func(ET.Either[E, A]) HKTB) HKTB, onLeft func(E) HKTB, onRight func(A) HKTB) func(HKTEA) HKTB {
	return F.Bind2nd(mchain, ET.Fold(onLeft, onRight))
}

func GetOrElse[E, A, HKTEA, HKTA any](mchain func(HKTEA, func(ET.Either[E, A]) HKTA) HKTA, mof func(A) HKTA, onLeft func(E) HKTA) func(HKTEA) HKTA {
	return MatchE(mchain, onLeft, mof)
}

func GetOrElseOf[E, A, HKTEA, HKTA any](mchain func(HKTEA, func(ET.Either[E, A]) HKTA) HKTA, mof func(A) HKTA, onLeft func(E) A) func(HKTEA) HKTA {
	return MatchE(mchain, F.Flow2(onLeft, mof), mof)
}

func OrElse[E1, E2, A, HKTE1A, HKTE2A any](mchain func(HKTE1A, func(ET.Either[E1, A]) HKTE2A) HKTE2A, mof func(ET.Either[E2, A]) HKTE2A, onLeft func(E1) HKTE2A) func(HKTE1A) HKTE2A {
	return MatchE(mchain, onLeft, F.Flow2(ET.Right[E2, A], mof))
}

func OrLeft[E1, E2, A, HKTE1A, HKTE2, HKTE2A any](
	mchain func(HKTE1A, func(ET.Either[E1, A]) HKTE2A) HKTE2A,
	mmap func(HKTE2, func(E2) ET.Either[E2, A]) HKTE2A,
	mof func(ET.Either[E2, A]) HKTE2A,
	onLeft func(E1) HKTE2) func(HKTE1A) HKTE2A {

	return F.Bind2nd(mchain, ET.Fold(F.Flow2(onLeft, F.Bind2nd(mmap, ET.Left[A, E2])), F.Flow2(ET.Right[E2, A], mof)))
}

func MonadMapLeft[E, A, B, HKTFA, HKTFB any](fmap func(HKTFA, func(ET.Either[E, A]) ET.Either[B, A]) HKTFB, fa HKTFA, f func(E) B) HKTFB {
	return FC.MonadMap(fmap, ET.MonadMapLeft[E, A, B], fa, f)
}

func MapLeft[E, A, B, HKTFA, HKTFB any](fmap func(func(ET.Either[E, A]) ET.Either[B, A]) func(HKTFA) HKTFB, f func(E) B) func(HKTFA) HKTFB {
	return FC.Map(fmap, ET.MapLeft[A, E, B], f)
}

func MonadChainLeft[EA, A, EB, HKTFA, HKTFB any](
	fchain func(HKTFA, func(ET.Either[EA, A]) HKTFB) HKTFB,
	fof func(ET.Either[EB, A]) HKTFB,
	fa HKTFA,
	f func(EA) HKTFB) HKTFB {
	return fchain(fa, ET.Fold(f, F.Flow2(ET.Right[EB, A], fof)))
}

func ChainLeft[EA, A, EB, HKTFA, HKTFB any](
	fchain func(func(ET.Either[EA, A]) HKTFB) func(HKTFA) HKTFB,
	fof func(ET.Either[EB, A]) HKTFB,
	f func(EA) HKTFB) func(HKTFA) HKTFB {
	return fchain(ET.Fold(f, F.Flow2(ET.Right[EB, A], fof)))
}

func MonadChainFirstLeft[EA, A, EB, B, HKTFA, HKTFB any](
	fchain func(HKTFA, func(ET.Either[EA, A]) HKTFA) HKTFA,
	fmap func(HKTFB, func(ET.Either[EB, B]) ET.Either[EA, A]) HKTFA,
	fof func(ET.Either[EA, A]) HKTFA,
	fa HKTFA,
	f func(EA) HKTFB) HKTFA {

	return fchain(fa, func(e ET.Either[EA, A]) HKTFA {
		return ET.Fold(func(ea EA) HKTFA {
			return fmap(f(ea), F.Constant1[ET.Either[EB, B]](e))
		}, F.Flow2(ET.Right[EA, A], fof))(e)
	})
}

func ChainFirstLeft[EA, A, EB, B, HKTFA, HKTFB any](
	fchain func(func(ET.Either[EA, A]) HKTFA) func(HKTFA) HKTFA,
	fmap func(func(ET.Either[EB, B]) ET.Either[EA, A]) func(HKTFB) HKTFA,
	fof func(ET.Either[EA, A]) HKTFA,
	f func(EA) HKTFB) func(HKTFA) HKTFA {

	return fchain(func(e ET.Either[EA, A]) HKTFA {
		return ET.Fold(F.Flow2(f, fmap(F.Constant1[ET.Either[EB, B]](e))), F.Flow2(ET.Right[EA, A], fof))(e)
	})
}
