package fromioeither

import (
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
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
	mchain func(HKTA, func(A) HKTA) HKTA,
	mmap func(HKTB, func(B) A) HKTA,
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
	mchain func(HKTA, func(A) HKTB) HKTB,
	fromio func(GIOB) HKTB,
	f func(A) GIOB) func(HKTA) HKTB {
	// chain
	return C.Chain[A, B](mchain, F.Flow2(f, fromio))
}
