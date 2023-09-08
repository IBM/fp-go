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

package bindt

import (
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	T "github.com/IBM/fp-go/tuple"
)

func Bind[SET ~func(B) func(S1) S2, FCT ~func(S1) HKTB, S1, S2, B, HKTS1, HKTS2, HKTB any](
	mchain func(func(S1) HKTS2) func(HKTS1) HKTS2,
	mmap func(func(B) S2) func(HKTB) HKTS2,
	s SET,
	f FCT,
) func(HKTS1) HKTS2 {
	return mchain(F.Flow3(
		T.Replicate2[S1],
		T.Map2(F.Flow2(
			I.Ap[S2, S1],
			F.Flow2(
				F.Bind1st(F.Flow2[SET, func(func(S1) S2) S2], s),
				mmap,
			)), f),
		T.Tupled2(I.MonadAp[HKTS2, HKTB]),
	))
}

func BindTo[SET ~func(B) S2, S2, B, HKTS2, HKTB any](
	mmap func(func(B) S2) func(HKTB) HKTS2,
	s SET,
) func(HKTB) HKTS2 {
	return mmap(s)
}

func ApS[
	SET ~func(B) func(S1) S2,
	S1, S2, B, HKTS1S2, HKTS1, HKTS2, HKTB any,
](
	ap func(HKTS1) func(HKTS1S2) HKTS2,
	mmap func(func(B) func(S1) S2) func(HKTB) HKTS1S2,
	s SET, fb HKTB) func(HKTS1) HKTS2 {

	return F.Flow2(
		ap,
		I.Ap[HKTS2, HKTS1S2](mmap(s)(fb)),
	)
}

func Let[SET ~func(B) func(S1) S2, FCT ~func(S1) B, S1, S2, B, HKTS1, HKTS2 any](
	mmap func(func(S1) S2) func(HKTS1) HKTS2,
	s SET,
	f FCT,
) func(HKTS1) HKTS2 {
	return mmap(F.Flow3(
		T.Replicate2[S1],
		T.Map2(F.Flow2(
			I.Ap[S2, S1],
			F.Bind1st(F.Flow2[SET, func(func(S1) S2) S2], s)), f),
		T.Tupled2(I.MonadAp[S2, B]),
	))
}
