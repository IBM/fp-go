package function

func Curry1[T1, R any](f func(T1) R) func(T1) R {
	return f
}

func Curry2[T1, T2, R any](f func(T1, T2) R) func(T1) func(T2) R {
	return func(t1 T1) func(T2) R {
		return func(t2 T2) R {
			return f(t1, t2)
		}
	}
}

func Curry3[T1, T2, T3, R any](f func(T1, T2, T3) R) func(T1) func(T2) func(T3) R {
	return func(t1 T1) func(T2) func(T3) R {
		return func(t2 T2) func(T3) R {
			return func(t3 T3) R {
				return f(t1, t2, t3)
			}
		}
	}
}

func Curry4[T1, T2, T3, T4, R any](f func(T1, T2, T3, T4) R) func(T1) func(T2) func(T3) func(T4) R {
	return func(t1 T1) func(T2) func(T3) func(T4) R {
		return func(t2 T2) func(T3) func(T4) R {
			return func(t3 T3) func(T4) R {
				return func(t4 T4) R {
					return f(t1, t2, t3, t4)
				}
			}
		}
	}
}

func Curry5[T1, T2, T3, T4, T5, R any](f func(T1, T2, T3, T4, T5) R) func(T1) func(T2) func(T3) func(T4) func(T5) R {
	return func(t1 T1) func(T2) func(T3) func(T4) func(T5) R {
		return func(t2 T2) func(T3) func(T4) func(T5) R {
			return func(t3 T3) func(T4) func(T5) R {
				return func(t4 T4) func(T5) R {
					return func(t5 T5) R {
						return f(t1, t2, t3, t4, t5)
					}
				}
			}
		}
	}
}

func Uncurry1[T1, R any](f func(T1) R) func(T1) R {
	return f
}

func Uncurry2[T1, T2, R any](f func(T1) func(T2) R) func(T1, T2) R {
	return func(t1 T1, t2 T2) R {
		return f(t1)(t2)
	}
}

func Uncurry3[T1, T2, T3, R any](f func(T1) func(T2) func(T3) R) func(T1, T2, T3) R {
	return func(t1 T1, t2 T2, t3 T3) R {
		return f(t1)(t2)(t3)
	}
}

func Uncurry4[T1, T2, T3, T4, R any](f func(T1) func(T2) func(T3) func(T4) R) func(T1, T2, T3, T4) R {
	return func(t1 T1, t2 T2, t3 T3, t4 T4) R {
		return f(t1)(t2)(t3)(t4)
	}
}

func Uncurry5[T1, T2, T3, T4, T5, R any](f func(T1) func(T2) func(T3) func(T4) func(T5) R) func(T1, T2, T3, T4, T5) R {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) R {
		return f(t1)(t2)(t3)(t4)(t5)
	}
}
