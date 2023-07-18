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
