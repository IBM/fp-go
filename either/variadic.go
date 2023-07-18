package either

func Variadic0[V, R any](f func([]V) (R, error)) func(...V) Either[error, R] {
	return func(v ...V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(v)
		})
	}
}

func Variadic1[T1, V, R any](f func(T1, []V) (R, error)) func(T1, ...V) Either[error, R] {
	return func(t1 T1, v ...V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, v)
		})
	}
}

func Variadic2[T1, T2, V, R any](f func(T1, T2, []V) (R, error)) func(T1, T2, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, v ...V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2, v)
		})
	}
}

func Variadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, []V) (R, error)) func(T1, T2, T3, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, v ...V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2, t3, v)
		})
	}
}

func Variadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, []V) (R, error)) func(T1, T2, T3, T4, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, v ...V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2, t3, t4, v)
		})
	}
}

func Unvariadic0[V, R any](f func(...V) (R, error)) func([]V) Either[error, R] {
	return func(v []V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(v...)
		})
	}
}

func Unvariadic1[T1, V, R any](f func(T1, ...V) (R, error)) func(T1, []V) Either[error, R] {
	return func(t1 T1, v []V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, v...)
		})
	}
}

func Unvariadic2[T1, T2, V, R any](f func(T1, T2, ...V) (R, error)) func(T1, T2, []V) Either[error, R] {
	return func(t1 T1, t2 T2, v []V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2, v...)
		})
	}
}

func Unvariadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, ...V) (R, error)) func(T1, T2, T3, []V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, v []V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2, t3, v...)
		})
	}
}

func Unvariadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, ...V) (R, error)) func(T1, T2, T3, T4, []V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, v []V) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2, t3, t4, v...)
		})
	}
}
