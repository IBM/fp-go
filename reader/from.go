package reader

import (
	G "github.com/ibm/fp-go/reader/generic"
)

func From0[R, A any](f func(R) A) Reader[R, A] {
	return G.From0[Reader[R, A]](f)
}

func From1[R, T1, A any](f func(R, T1) A) func(T1) Reader[R, A] {
	return G.From1[Reader[R, A]](f)
}

func From2[R, T1, T2, A any](f func(R, T1, T2) A) func(T1, T2) Reader[R, A] {
	return G.From2[Reader[R, A]](f)
}

func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) A) func(T1, T2, T3) Reader[R, A] {
	return G.From3[Reader[R, A]](f)
}

func From4[R, T1, T2, T3, T4, A any](f func(R, T1, T2, T3, T4) A) func(T1, T2, T3, T4) Reader[R, A] {
	return G.From4[Reader[R, A]](f)
}
