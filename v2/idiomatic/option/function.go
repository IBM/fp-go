package option

// Flow1 creates a function that takes an initial value t0 and successively applies 1 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//go:inline
func Flow1[F1 ~func(T0, bool) (T1, bool), T0, T1 any](f1 F1) func(T0, bool) (T1, bool) {
	return f1
}

// Flow2 creates a function that takes an initial value t0 and successively applies 2 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//go:inline
func Flow2[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), T0, T1, T2 any](f1 F1, f2 F2) func(T0, bool) (T2, bool) {
	return func(t0 T0, t0ok bool) (T2, bool) {
		return f2(f1(t0, t0ok))
	}
}

// Flow3 creates a function that takes an initial value t0 and successively applies 3 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//go:inline
func Flow3[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), T0, T1, T2, T3 any](f1 F1, f2 F2, f3 F3) func(T0, bool) (T3, bool) {
	return func(t0 T0, t0ok bool) (T3, bool) {
		return f3(f2(f1(t0, t0ok)))
	}
}

// Flow4 creates a function that takes an initial value t0 and successively applies 4 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//go:inline
func Flow4[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), T0, T1, T2, T3, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T0, bool) (T4, bool) {
	return func(t0 T0, t0ok bool) (T4, bool) {
		return f4(f3(f2(f1(t0, t0ok))))
	}
}
