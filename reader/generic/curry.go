package generic

import (
	F "github.com/ibm/fp-go/function"
)

// these functions curry a golang function with the context as the firsr parameter into a reader with the context as the last parameter, which
// is a equivalent to a function returning a reader of that context
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Curry0[GA ~func(R) A, R, A any](f func(R) A) GA {
	return MakeReader[GA](f)
}

func Curry1[GA ~func(R) A, R, T1, A any](f func(R, T1) A) func(T1) GA {
	return F.Curry1(From1[GA](f))
}

func Curry2[GA ~func(R) A, R, T1, T2, A any](f func(R, T1, T2) A) func(T1) func(T2) GA {
	return F.Curry2(From2[GA](f))
}

func Curry3[GA ~func(R) A, R, T1, T2, T3, A any](f func(R, T1, T2, T3) A) func(T1) func(T2) func(T3) GA {
	return F.Curry3(From3[GA](f))
}

func Curry4[GA ~func(R) A, R, T1, T2, T3, T4, A any](f func(R, T1, T2, T3, T4) A) func(T1) func(T2) func(T3) func(T4) GA {
	return F.Curry4(From4[GA](f))
}

func Uncurry0[GA ~func(R) A, R, A any](f GA) func(R) A {
	return f
}

func Uncurry1[GA ~func(R) A, R, T1, A any](f func(T1) GA) func(R, T1) A {
	uc := F.Uncurry1(f)
	return func(r R, t1 T1) A {
		return uc(t1)(r)
	}
}

func Uncurry2[GA ~func(R) A, R, T1, T2, A any](f func(T1) func(T2) GA) func(R, T1, T2) A {
	uc := F.Uncurry2(f)
	return func(r R, t1 T1, t2 T2) A {
		return uc(t1, t2)(r)
	}
}

func Uncurry3[GA ~func(R) A, R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) GA) func(R, T1, T2, T3) A {
	uc := F.Uncurry3(f)
	return func(r R, t1 T1, t2 T2, t3 T3) A {
		return uc(t1, t2, t3)(r)
	}
}

func Uncurry4[GA ~func(R) A, R, T1, T2, T3, T4, A any](f func(T1) func(T2) func(T3) func(T4) GA) func(R, T1, T2, T3, T4) A {
	uc := F.Uncurry4(f)
	return func(r R, t1 T1, t2 T2, t3 T3, t4 T4) A {
		return uc(t1, t2, t3, t4)(r)
	}
}
