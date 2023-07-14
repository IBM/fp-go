package fromreader

import (
	F "github.com/ibm/fp-go/function"
	G "github.com/ibm/fp-go/reader/generic"
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
