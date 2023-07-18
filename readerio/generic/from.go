package generic

import (
	G "github.com/IBM/fp-go/reader/generic"
)

// these functions From a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[GEA ~func(R) GIOA, GIOA ~func() A, R, A any](f func(R) GIOA) func() GEA {
	return G.From0[GEA](f)
}

func From1[GEA ~func(R) GIOA, GIOA ~func() A, R, T1, A any](f func(R, T1) GIOA) func(T1) GEA {
	return G.From1[GEA](f)
}

func From2[GEA ~func(R) GIOA, GIOA ~func() A, R, T1, T2, A any](f func(R, T1, T2) GIOA) func(T1, T2) GEA {
	return G.From2[GEA](f)
}

func From3[GEA ~func(R) GIOA, GIOA ~func() A, R, T1, T2, T3, A any](f func(R, T1, T2, T3) GIOA) func(T1, T2, T3) GEA {
	return G.From3[GEA](f)
}
