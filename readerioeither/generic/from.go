package generic

import (
	ET "github.com/ibm/fp-go/either"
	G "github.com/ibm/fp-go/reader/generic"
)

// these functions From a golang function with the context as the first parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, A any](f func(R) func() (A, error)) GEA {
	return G.From0[GEA](func(r R) GIOA {
		return ET.Eitherize0(f(r))
	})
}

func From1[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, T1, A any](f func(R, T1) func() (A, error)) func(T1) GEA {
	return G.From1[GEA](func(r R, t1 T1) GIOA {
		return ET.Eitherize0(f(r, t1))
	})
}

func From2[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, T1, T2, A any](f func(R, T1, T2) func() (A, error)) func(T1, T2) GEA {
	return G.From2[GEA](func(r R, t1 T1, t2 T2) GIOA {
		return ET.Eitherize0(f(r, t1, t2))
	})
}

func From3[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, T1, T2, T3, A any](f func(R, T1, T2, T3) func() (A, error)) func(T1, T2, T3) GEA {
	return G.From3[GEA](func(r R, t1 T1, t2 T2, t3 T3) GIOA {
		return ET.Eitherize0(f(r, t1, t2, t3))
	})
}

func From4[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, T1, T2, T3, T4, A any](f func(R, T1, T2, T3, T4) func() (A, error)) func(T1, T2, T3, T4) GEA {
	return G.From4[GEA](func(r R, t1 T1, t2 T2, t3 T3, t4 T4) GIOA {
		return ET.Eitherize0(f(r, t1, t2, t3, t4))
	})
}
