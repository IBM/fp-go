package generic

// these functions convert a golang function with the context as the first parameter into a reader with the context as the last parameter, which
// is a equivalent to a function returning a reader of that context
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[GA ~func(R) A, R, A any](f func(R) A) GA {
	return MakeReader[GA](f)
}

func From1[GA ~func(R) A, R, T1, A any](f func(R, T1) A) func(T1) GA {
	return func(t1 T1) GA {
		return MakeReader[GA](func(r R) A {
			return f(r, t1)
		})
	}
}

func From2[GA ~func(R) A, R, T1, T2, A any](f func(R, T1, T2) A) func(T1, T2) GA {
	return func(t1 T1, t2 T2) GA {
		return MakeReader[GA](func(r R) A {
			return f(r, t1, t2)
		})
	}
}

func From3[GA ~func(R) A, R, T1, T2, T3, A any](f func(R, T1, T2, T3) A) func(T1, T2, T3) GA {
	return func(t1 T1, t2 T2, t3 T3) GA {
		return MakeReader[GA](func(r R) A {
			return f(r, t1, t2, t3)
		})
	}
}

func From4[GA ~func(R) A, R, T1, T2, T3, T4, A any](f func(R, T1, T2, T3, T4) A) func(T1, T2, T3, T4) GA {
	return func(t1 T1, t2 T2, t3 T3, t4 T4) GA {
		return MakeReader[GA](func(r R) A {
			return f(r, t1, t2, t3, t4)
		})
	}
}
