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

package chain

import (
	F "github.com/IBM/fp-go/v2/function"
)

// HKTA=HKT[A]
// HKTB=HKT[B]
func MonadChainFirst[A, B, HKTA, HKTB any](
	mchain func(HKTA, Kleisli[A, HKTA]) HKTA,
	mmap func(HKTB, func(B) A) HKTA,
	first HKTA,
	f Kleisli[A, HKTB],
) HKTA {
	return mchain(first, func(a A) HKTA {
		return mmap(f(a), F.Constant1[B](a))
	})
}

func MonadChain[A, B, HKTA, HKTB any](
	mchain func(HKTA, Kleisli[A, HKTB]) HKTB,
	first HKTA,
	f Kleisli[A, HKTB],
) HKTB {
	return mchain(first, f)
}

// HKTA=HKT[A]
// HKTB=HKT[B]
func ChainFirst[A, B, HKTA, HKTB any](
	mchain func(Kleisli[A, HKTA]) Operator[HKTA, HKTA],
	mmap func(func(B) A) func(HKTB) HKTA,
	f Kleisli[A, HKTB]) Operator[HKTA, HKTA] {

	return mchain(func(a A) HKTA {
		return mmap(F.Constant1[B](a))(f(a))
	})
}

func Chain[A, B, HKTA, HKTB any](
	mchain func(Kleisli[A, HKTB]) Operator[HKTA, HKTB],
	f Kleisli[A, HKTB],
) Operator[HKTA, HKTB] {
	return mchain(f)
}

func MonadBind[S1, S2, B, HKTS1, HKTS2, HKTB any](
	mchain func(HKTS1, func(S1) HKTS2) HKTS2,
	mmap func(HKTB, func(B) S2) HKTS2,
	first HKTS1,
	key func(B) func(S1) S2,
	f func(S1) HKTB,
) HKTS2 {
	return mchain(first, func(s1 S1) HKTS2 {
		return mmap(f(s1), func(b B) S2 {
			return key(b)(s1)
		})
	})
}

func Bind[S1, S2, B, HKTS1, HKTS2, HKTB any](
	mchain func(func(S1) HKTS2) func(HKTS1) HKTS2,
	mmap func(func(B) S2) func(HKTB) HKTS2,
	key func(B) func(S1) S2,
	f func(S1) HKTB,
) func(HKTS1) HKTS2 {
	mapb := F.Flow2(
		F.Flip(key),
		mmap,
	)
	return mchain(func(s1 S1) HKTS2 {
		return F.Pipe2(
			s1,
			f,
			F.Pipe1(
				s1,
				mapb,
			),
		)
	})
}

func BindTo[S1, B, HKTS1, HKTB any](
	mmap func(func(B) S1) func(HKTB) HKTS1,
	key func(B) S1,
) func(fa HKTB) HKTS1 {
	return mmap(key)
}

func MonadBindTo[S1, B, HKTS1, HKTB any](
	mmap func(HKTB, func(B) S1) HKTS1,
	first HKTB,
	key func(B) S1,
) HKTS1 {
	return mmap(first, key)
}
