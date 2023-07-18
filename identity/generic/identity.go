package generic

import (
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
)

func MonadAp[GAB ~func(A) B, B, A any](fab GAB, fa A) B {
	return fab(fa)
}

func Ap[GAB ~func(A) B, B, A any](fa A) func(GAB) B {
	return F.Bind2nd(MonadAp[GAB, B, A], fa)
}

func MonadMap[GAB ~func(A) B, A, B any](fa A, f GAB) B {
	return f(fa)
}

func Map[GAB ~func(A) B, A, B any](f GAB) func(A) B {
	return f
}

func MonadChain[GAB ~func(A) B, A, B any](ma A, f GAB) B {
	return f(ma)
}

func Chain[GAB ~func(A) B, A, B any](f GAB) func(A) B {
	return f
}

func MonadChainFirst[GAB ~func(A) B, A, B any](fa A, f GAB) A {
	return C.MonadChainFirst(MonadChain[func(A) A, A, A], MonadMap[func(B) A, B, A], fa, f)
}

func ChainFirst[GAB ~func(A) B, A, B any](f GAB) func(A) A {
	return C.ChainFirst(MonadChain[func(A) A, A, A], MonadMap[func(B) A, B, A], f)
}
