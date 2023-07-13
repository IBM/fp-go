package function

func Bind1st[T1, T2, R any](f func(T1, T2) R, t1 T1) func(T2) R {
	return func(t2 T2) R {
		return f(t1, t2)
	}
}
func Bind2nd[T1, T2, R any](f func(T1, T2) R, t2 T2) func(T1) R {
	return func(t1 T1) R {
		return f(t1, t2)
	}
}

func SK[T1, T2 any](_ T1, t2 T2) T2 {
	return t2
}
