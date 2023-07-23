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

package fromreader

import (
	F "github.com/IBM/fp-go/function"
	G "github.com/IBM/fp-go/reader/generic"
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
	mchain func(HKTRA, func(A) HKTRB) HKTRB,
	fromReader func(GB) HKTRB,
	f func(A) GB,
) func(HKTRA) HKTRB {
	return F.Bind2nd(mchain, FromReaderK(fromReader, f))
}
