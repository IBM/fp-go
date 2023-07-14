package identity

import (
	F "github.com/ibm/fp-go/function"
	G "github.com/ibm/fp-go/identity/generic"
)

func MonadAp[B, A any](fab func(A) B, fa A) B {
	return G.MonadAp(fab, fa)
}

func Ap[B, A any](fa A) func(func(A) B) B {
	return G.Ap[func(A) B](fa)
}

func MonadMap[A, B any](fa A, f func(A) B) B {
	return G.MonadMap(fa, f)
}

func Map[A, B any](f func(A) B) func(A) B {
	return G.Map(f)
}

func MonadMapTo[A, B any](fa A, b B) B {
	return b
}

func MapTo[A, B any](b B) func(A) B {
	return F.Constant1[A](b)
}

func Of[A any](a A) A {
	return a
}

func MonadChain[A, B any](ma A, f func(A) B) B {
	return G.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) B) func(A) B {
	return G.Chain(f)
}

func MonadChainFirst[A, B any](fa A, f func(A) B) A {
	return G.MonadChainFirst(fa, f)
}

func ChainFirst[A, B any](f func(A) B) func(A) A {
	return G.ChainFirst(f)
}
