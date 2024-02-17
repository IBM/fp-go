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
	F "github.com/IBM/fp-go/function"
	P "github.com/IBM/fp-go/pair"
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
		F.Bind2nd(fmap, func(a P.Pair[A, S]) P.Pair[B, S] {
			return P.MakePair(f(P.Head(a)), P.Tail(a))
		}),
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
	mp := fmap(func(a P.Pair[A, S]) P.Pair[B, S] {
		return P.MakePair(f(P.Head(a)), P.Tail(a))
	})

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
