package fromio

import (
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
)

func MonadChainFirstIOK[A, B, HKTA, HKTB any, GIOB ~func() B](
	mchain func(HKTA, func(A) HKTA) HKTA,
	mmap func(HKTB, func(B) A) HKTA,
	fromio func(GIOB) HKTB,
	first HKTA, f func(A) GIOB) HKTA {
	// chain
	return C.MonadChainFirst(mchain, mmap, first, F.Flow2(f, fromio))
}

func ChainFirstIOK[A, B, HKTA, HKTB any, GIOB ~func() B](
	mchain func(HKTA, func(A) HKTA) HKTA,
	mmap func(HKTB, func(B) A) HKTA,
	fromio func(GIOB) HKTB,
	f func(A) GIOB) func(HKTA) HKTA {
	// chain
	return C.ChainFirst(mchain, mmap, F.Flow2(f, fromio))
}

func MonadChainIOK[GR ~func() B, A, B, HKTA, HKTB any](
	mchain func(HKTA, func(A) HKTB) HKTB,
	fromio func(GR) HKTB,
	first HKTA, f func(A) GR) HKTB {
	// chain
	return C.MonadChain[A, B](mchain, first, F.Flow2(f, fromio))
}

func ChainIOK[GR ~func() B, A, B, HKTA, HKTB any](
	mchain func(HKTA, func(A) HKTB) HKTB,
	fromio func(GR) HKTB,
	f func(A) GR) func(HKTA) HKTB {
	// chain
	return C.Chain[A, B](mchain, F.Flow2(f, fromio))
}
