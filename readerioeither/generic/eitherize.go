package generic

import (
	ET "github.com/ibm/fp-go/either"
)

// these functions From a golang function with the context as the first parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Eitherize0[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, A any](f func(R) (A, error)) GEA {
	return From0[GEA](func(r R) func() (A, error) {
		return func() (A, error) {
			return f(r)
		}
	})
}

func Eitherize1[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, T1, A any](f func(R, T1) (A, error)) func(T1) GEA {
	return From1[GEA](func(r R, t1 T1) func() (A, error) {
		return func() (A, error) {
			return f(r, t1)
		}
	})
}

func Eitherize2[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1, T2) GEA {
	return From2[GEA](func(r R, t1 T1, t2 T2) func() (A, error) {
		return func() (A, error) {
			return f(r, t1, t2)
		}
	})
}

func Eitherize3[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1, T2, T3) GEA {
	return From3[GEA](func(r R, t1 T1, t2 T2, t3 T3) func() (A, error) {
		return func() (A, error) {
			return f(r, t1, t2, t3)
		}
	})
}

func Eitherize4[GEA ~func(R) GIOA, GIOA ~func() ET.Either[error, A], R, T1, T2, T3, T4, A any](f func(R, T1, T2, T3, T4) (A, error)) func(T1, T2, T3, T4) GEA {
	return From4[GEA](func(r R, t1 T1, t2 T2, t3 T3, t4 T4) func() (A, error) {
		return func() (A, error) {
			return f(r, t1, t2, t3, t4)
		}
	})
}
