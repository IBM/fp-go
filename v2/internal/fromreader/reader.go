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

package fromreader

import (
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
	G "github.com/IBM/fp-go/v2/reader/generic"
)

func Ask[GR ~func(R) R, R, HKTRA any](fromReader func(GR) HKTRA) func() HKTRA {
	return func() HKTRA {
		return fromReader(G.Ask[GR]())
	}
}

func Asks[GA ~func(R) A, R, A, HKTRA any](fromReader func(GA) HKTRA) func(GA) HKTRA {
	return fromReader
}

func FromReaderK[GB ~func(R) B, R, A, B, HKTRB any](
	fromReader func(GB) HKTRB,
	f func(A) GB) func(A) HKTRB {
	return F.Flow2(f, fromReader)
}

func MonadChainReaderK[GB ~func(R) B, R, A, B, HKTRA, HKTRB any](
	mchain func(HKTRA, func(A) HKTRB) HKTRB,
	fromReader func(GB) HKTRB,
	ma HKTRA,
	f func(A) GB,
) HKTRB {
	return mchain(ma, FromReaderK(fromReader, f))
}

func ChainReaderK[GB ~func(R) B, R, A, B, HKTRA, HKTRB any](
	mchain func(func(A) HKTRB) func(HKTRA) HKTRB,
	fromReader func(GB) HKTRB,
	f func(A) GB,
) func(HKTRA) HKTRB {
	return mchain(FromReaderK(fromReader, f))
}

func MonadChainFirstReaderK[GB ~func(R) B, R, A, B, HKTRA, HKTRB any](
	mchain func(HKTRA, func(A) HKTRB) HKTRA,
	fromReader func(GB) HKTRB,
	ma HKTRA,
	f func(A) GB,
) HKTRA {
	return mchain(ma, FromReaderK(fromReader, f))
}

func ChainFirstReaderK[GB ~func(R) B, R, A, B, HKTRA, HKTRB any](
	mchain func(func(A) HKTRB) func(HKTRA) HKTRA,
	fromReader func(GB) HKTRB,
	f func(A) GB,
) func(HKTRA) HKTRA {
	return mchain(FromReaderK(fromReader, f))
}

//go:inline
func BindReaderK[
	GT ~func(R) T,
	R, S1, S2, T,
	HKTET,
	HKTES1,
	HKTES2 any](
	mchain func(func(S1) HKTES2) func(HKTES1) HKTES2,
	mmap func(func(T) S2) func(HKTET) HKTES2,
	fromReader func(GT) HKTET,
	setter func(T) func(S1) S2,
	f func(S1) GT,
) func(HKTES1) HKTES2 {
	return C.Bind(
		mchain,
		mmap,
		setter,
		FromReaderK(fromReader, f),
	)
}
