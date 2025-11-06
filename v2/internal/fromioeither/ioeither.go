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

package fromioeither

import (
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
)

func MonadChainFirstIOEitherK[GIOB ~func() ET.Either[E, B], E, A, B, HKTA, HKTB any](
	mchain func(HKTA, func(A) HKTA) HKTA,
	mmap func(HKTB, func(B) A) HKTA,
	fromio func(GIOB) HKTB,
	first HKTA, f func(A) GIOB) HKTA {
	// chain
	return C.MonadChainFirst(mchain, mmap, first, F.Flow2(f, fromio))
}

func ChainFirstIOEitherK[GIOB ~func() ET.Either[E, B], E, A, B, HKTA, HKTB any](
	mchain func(func(A) HKTA) func(HKTA) HKTA,
	mmap func(func(B) A) func(HKTB) HKTA,
	fromio func(GIOB) HKTB,
	f func(A) GIOB) func(HKTA) HKTA {
	// chain
	return C.ChainFirst(mchain, mmap, F.Flow2(f, fromio))
}

func MonadChainIOEitherK[GIOB ~func() ET.Either[E, B], E, A, B, HKTA, HKTB any](
	mchain func(HKTA, func(A) HKTB) HKTB,
	fromio func(GIOB) HKTB,
	first HKTA, f func(A) GIOB) HKTB {
	// chain
	return C.MonadChain[A, B](mchain, first, F.Flow2(f, fromio))
}

func ChainIOEitherK[GIOB ~func() ET.Either[E, B], E, A, B, HKTA, HKTB any](
	mchain func(func(A) HKTB) func(HKTA) HKTB,
	fromio func(GIOB) HKTB,
	f func(A) GIOB) func(HKTA) HKTB {
	// chain
	return C.Chain[A, B](mchain, F.Flow2(f, fromio))
}
