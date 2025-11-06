// Copyright (c) 2024 - 2025 IBM Corp.
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

package statet

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
)

func Of[
	HKTSA ~func(S) HKTA,
	HKTA,
	S, A any,
](
	fof func(pair.Pair[S, A]) HKTA,

	a A) HKTSA {

	return function.Flow2(
		function.Bind2nd(pair.MakePair[S, A], a),
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
	fmap func(HKTA, func(pair.Pair[S, A]) pair.Pair[S, B]) HKTB,

	fa HKTSA,
	f func(A) B,
) HKTSB {

	return function.Flow2(
		fa,
		function.Bind2nd(fmap, pair.Map[S](f)),
	)
}

func Map[
	HKTSA ~func(S) HKTA,
	HKTSB ~func(S) HKTB,
	HKTA,
	HKTB,
	S, A, B any,
](
	fmap func(func(pair.Pair[S, A]) pair.Pair[S, B]) func(HKTA) HKTB,

	f func(A) B,
) func(HKTSA) HKTSB {
	mp := fmap(pair.Map[S](f))

	return func(fa HKTSA) HKTSB {
		return function.Flow2(
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
	fchain func(HKTA, func(pair.Pair[S, A]) HKTB) HKTB,

	fa HKTSA,
	f func(A) HKTSB,
) HKTSB {
	return function.Flow2(
		fa,
		function.Bind2nd(fchain, func(a pair.Pair[S, A]) HKTB {
			return f(pair.Tail(a))(pair.Head(a))
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
	fchain func(func(pair.Pair[S, A]) HKTB) func(HKTA) HKTB,

	f func(A) HKTSB,
) func(HKTSA) HKTSB {
	mp := fchain(func(a pair.Pair[S, A]) HKTB {
		return f(pair.Tail(a))(pair.Head(a))
	})

	return func(fa HKTSA) HKTSB {
		return function.Flow2(
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
	fmap func(HKTA, func(pair.Pair[S, A]) pair.Pair[S, B]) HKTB,
	fchain func(HKTAB, func(pair.Pair[S, func(A) B]) HKTB) HKTB,

	fab HKTSAB,
	fa HKTSA,
) HKTSB {
	return func(s S) HKTB {
		return fchain(fab(s), func(ab pair.Pair[S, func(A) B]) HKTB {
			return fmap(fa(pair.Head(ab)), pair.Map[S](pair.Tail(ab)))
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
	fmap func(func(pair.Pair[S, A]) pair.Pair[S, B]) func(HKTA) HKTB,
	fchain func(func(pair.Pair[S, func(A) B]) HKTB) func(HKTAB) HKTB,

	fa HKTSA,
) func(HKTSAB) HKTSB {
	return func(fab HKTSAB) HKTSB {
		return function.Flow2(
			fab,
			fchain(func(ab pair.Pair[S, func(A) B]) HKTB {
				return fmap(pair.Map[S](pair.Tail(ab)))(fa(pair.Head(ab)))
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
	fmap func(HKTFA, func(A) pair.Pair[S, A]) HKTA,
	ma HKTFA) HKTSA {

	f1 := function.Bind1st(fmap, ma)

	return func(s S) HKTA {
		return f1(function.Bind1st(pair.MakePair[S, A], s))
	}
}

func FromState[
	HKTSA ~func(S) HKTA,
	ST ~func(S) pair.Pair[S, A],
	HKTA,

	S, A any,
](
	fof func(pair.Pair[S, A]) HKTA,
	sa ST,
) HKTSA {
	return function.Flow2(sa, fof)
}
