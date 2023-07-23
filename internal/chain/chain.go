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

package chain

import (
	F "github.com/IBM/fp-go/function"
)

// HKTA=HKT[A]
// HKTB=HKT[B]
func MonadChainFirst[A, B, HKTA, HKTB any](
	mchain func(HKTA, func(A) HKTA) HKTA,
	mmap func(HKTB, func(B) A) HKTA,
	first HKTA,
	f func(A) HKTB,
) HKTA {
	return mchain(first, func(a A) HKTA {
		return mmap(f(a), F.Constant1[B](a))
	})
}

func MonadChain[A, B, HKTA, HKTB any](
	mchain func(HKTA, func(A) HKTB) HKTB,
	first HKTA,
	f func(A) HKTB,
) HKTB {
	return mchain(first, f)
}

// HKTA=HKT[A]
// HKTB=HKT[B]
func ChainFirst[A, B, HKTA, HKTB any](
	mchain func(HKTA, func(A) HKTA) HKTA,
	mmap func(HKTB, func(B) A) HKTA,
	f func(A) HKTB) func(HKTA) HKTA {
	return F.Bind2nd(mchain, func(a A) HKTA {
		return mmap(f(a), F.Constant1[B](a))
	})
}

func Chain[A, B, HKTA, HKTB any](
	mchain func(HKTA, func(A) HKTB) HKTB,
	f func(A) HKTB,
) func(HKTA) HKTB {
	return func(first HKTA) HKTB {
		return MonadChain[A, B](mchain, first, f)
	}
}
