// Copyright (c) 2024 IBM Corp.
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

package state

import (
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/pair"
)

func Of[
	HKTSA ~func(S) HKTA,
	HKTA,
	S, A any,
](
	fof func(P.Pair[A, S]) HKTA,

	a A) HKTSA {

	return F.Flow2(
		F.Bind1st(P.MakePair[A, S], a),
		fof,
	)
}

func MonadMap[
	HKTSA ~func(S) HKTA,
	HKTSB ~func(S) HKTB,
	HKTA,
	HKTB,
	S, A, B any,
](
	fmap func(HKTA, func(P.Pair[A, S]) P.Pair[B, S]) HKTB,

	fa HKTSA,
	f func(A) B,
) HKTSB {

	return F.Flow2(
		fa,
		F.Bind2nd(fmap, P.Map[S](f)),
	)
}

func Map[
	HKTSA ~func(S) HKTA,
	HKTSB ~func(S) HKTB,
	HKTA,
	HKTB,
	S, A, B any,
](
	fmap func(func(P.Pair[A, S]) P.Pair[B, S]) func(HKTA) HKTB,

	f func(A) B,
) func(HKTSA) HKTSB {
	mp := fmap(P.Map[S](f))

	return func(fa HKTSA) HKTSB {
		return F.Flow2(
			fa,
			mp,
		)
	}
}

func MonadChain[
	HKTSA ~func(S) HKTA,
	HKTSB ~func(S) HKTB,
	HKTA,
	HKTB,
	S, A any,
](
	fchain func(HKTA, func(P.Pair[A, S]) HKTB) HKTB,

	fa HKTSA,
	f func(A) HKTSB,
) HKTSB {
	return F.Flow2(
		fa,
		F.Bind2nd(fchain, func(a P.Pair[A, S]) HKTB {
			return f(P.Head(a))(P.Tail(a))
		}),
	)
}

func Chain[
	HKTSA ~func(S) HKTA,
	HKTSB ~func(S) HKTB,
	HKTA,
	HKTB,
	S, A any,
](
	fchain func(func(P.Pair[A, S]) HKTB) func(HKTA) HKTB,

	f func(A) HKTSB,
) func(HKTSA) HKTSB {
	mp := fchain(func(a P.Pair[A, S]) HKTB {
		return f(P.Head(a))(P.Tail(a))
	})

	return func(fa HKTSA) HKTSB {
		return F.Flow2(
			fa,
			mp,
		)
	}
}

func MonadAp[
	HKTSA ~func(S) HKTA,
	HKTSB ~func(S) HKTB,
	HKTSAB ~func(S) HKTAB,
	HKTA,
	HKTB,
	HKTAB,

	S, A, B any,
](
	fmap func(HKTA, func(P.Pair[A, S]) P.Pair[B, S]) HKTB,
	fchain func(HKTAB, func(P.Pair[func(A) B, S]) HKTB) HKTB,

	fab HKTSAB,
	fa HKTSA,
) HKTSB {
	return func(s S) HKTB {
		return fchain(fab(s), func(ab P.Pair[func(A) B, S]) HKTB {
			return fmap(fa(P.Tail(ab)), P.Map[S](P.Head(ab)))
		})
	}
}

func Ap[
	HKTSA ~func(S) HKTA,
	HKTSB ~func(S) HKTB,
	HKTSAB ~func(S) HKTAB,
	HKTA,
	HKTB,
	HKTAB,

	S, A, B any,
](
	fmap func(func(P.Pair[A, S]) P.Pair[B, S]) func(HKTA) HKTB,
	fchain func(func(P.Pair[func(A) B, S]) HKTB) func(HKTAB) HKTB,

	fa HKTSA,
) func(HKTSAB) HKTSB {
	return func(fab HKTSAB) HKTSB {
		return F.Flow2(
			fab,
			fchain(func(ab P.Pair[func(A) B, S]) HKTB {
				return fmap(P.Map[S](P.Head(ab)))(fa(P.Tail(ab)))
			}),
		)
	}
}

func FromF[
	HKTSA ~func(S) HKTA,
	HKTA,

	HKTFA,

	S, A any,
](
	fmap func(HKTFA, func(A) P.Pair[A, S]) HKTA,
	ma HKTFA) HKTSA {

	f1 := F.Bind1st(fmap, ma)

	return func(s S) HKTA {
		return f1(F.Bind2nd(P.MakePair[A, S], s))
	}
}

func FromState[
	HKTSA ~func(S) HKTA,
	ST ~func(S) P.Pair[A, S],
	HKTA,

	S, A any,
](
	fof func(P.Pair[A, S]) HKTA,
	sa ST,
) HKTSA {
	return F.Flow2(sa, fof)
}
