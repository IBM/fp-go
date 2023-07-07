package function

// what a mess, golang does not have proper types ...

func Pipe1[A, R any](a A, f1 func(a A) R) R {
	// return f1(a)
	r1 := f1(a)
	return r1
}

func Pipe2[A, T1, R any](a A, f1 func(a A) T1, f2 func(t1 T1) R) R {
	// return f2(f1(a))
	r1 := f1(a)
	r2 := f2(r1)
	return r2
}

func Pipe3[A, T1, T2, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) R) R {
	// return f3(f2(f1(a)))
	r1 := f1(a)
	r2 := f2(r1)
	r3 := f3(r2)
	return r3
}

func Pipe4[A, T1, T2, T3, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) R) R {
	// return f4(f3(f2(f1(a))))
	r1 := f1(a)
	r2 := f2(r1)
	r3 := f3(r2)
	r4 := f4(r3)
	return r4
}

func Pipe5[A, T1, T2, T3, T4, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) R) R {
	// return f5(f4(f3(f2(f1(a)))))
	r1 := f1(a)
	r2 := f2(r1)
	r3 := f3(r2)
	r4 := f4(r3)
	r5 := f5(r4)
	return r5
}

func Pipe6[A, T1, T2, T3, T4, T5, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) R) R {
	// return f6(f5(f4(f3(f2(f1(a))))))
	r1 := f1(a)
	r2 := f2(r1)
	r3 := f3(r2)
	r4 := f4(r3)
	r5 := f5(r4)
	r6 := f6(r5)
	return r6
}

func Pipe7[A, T1, T2, T3, T4, T5, T6, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) R) R {
	// return f7(f6(f5(f4(f3(f2(f1(a)))))))
	r1 := f1(a)
	r2 := f2(r1)
	r3 := f3(r2)
	r4 := f4(r3)
	r5 := f5(r4)
	r6 := f6(r5)
	r7 := f7(r6)
	return r7
}

func Pipe8[A, T1, T2, T3, T4, T5, T6, T7, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) R) R {
	// return f8(f7(f6(f5(f4(f3(f2(f1(a))))))))
	r1 := f1(a)
	r2 := f2(r1)
	r3 := f3(r2)
	r4 := f4(r3)
	r5 := f5(r4)
	r6 := f6(r5)
	r7 := f7(r6)
	r8 := f8(r7)
	return r8
}

func Pipe9[A, T1, T2, T3, T4, T5, T6, T7, T8, R any](a A, f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) T8, f9 func(t8 T8) R) R {
	// return f9(f8(f7(f6(f5(f4(f3(f2(f1(a)))))))))
	r1 := f1(a)
	r2 := f2(r1)
	r3 := f3(r2)
	r4 := f4(r3)
	r5 := f5(r4)
	r6 := f6(r5)
	r7 := f7(r6)
	r8 := f8(r7)
	r9 := f9(r8)
	return r9
}

func Flow1[A, R any](f1 func(a A) R) func(a A) R {
	return f1
}

func Flow2[A, T1, R any](f1 func(a A) T1, f2 func(t1 T1) R) func(a A) R {
	return func(a A) R {
		return Pipe2(a, f1, f2)
	}
}

func Flow3[A, T1, T2, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) R) func(a A) R {
	return func(a A) R {
		return Pipe3(a, f1, f2, f3)
	}
}

func Flow4[A, T1, T2, T3, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) R) func(a A) R {
	return func(a A) R {
		return Pipe4(a, f1, f2, f3, f4)
	}
}

func Flow5[A, T1, T2, T3, T4, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) R) func(a A) R {
	return func(a A) R {
		return Pipe5(a, f1, f2, f3, f4, f5)
	}
}

func Flow6[A, T1, T2, T3, T4, T5, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) R) func(a A) R {
	return func(a A) R {
		return Pipe6(a, f1, f2, f3, f4, f5, f6)
	}
}

func Flow7[A, T1, T2, T3, T4, T5, T6, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) R) func(a A) R {
	return func(a A) R {
		return Pipe7(a, f1, f2, f3, f4, f5, f6, f7)
	}
}

func Flow8[A, T1, T2, T3, T4, T5, T6, T7, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) R) func(a A) R {
	return func(a A) R {
		return Pipe8(a, f1, f2, f3, f4, f5, f6, f7, f8)
	}
}

func Flow9[A, T1, T2, T3, T4, T5, T6, T7, T8, R any](f1 func(a A) T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) T4, f5 func(t4 T4) T5, f6 func(t5 T5) T6, f7 func(t6 T6) T7, f8 func(t7 T7) T8, f9 func(t8 T8) R) func(a A) R {
	return func(a A) R {
		return Pipe9(a, f1, f2, f3, f4, f5, f6, f7, f8, f9)
	}
}

// Identity returns the value 'a'
func Identity[A any](a A) A {
	return a
}

// Constant creates a nullary function that returns the constant value 'a'
func Constant[A any](a A) func() A {
	return func() A {
		return a
	}
}

// Constant1 creates a unary function that returns the constant value 'a' and ignores its input
func Constant1[B, A any](a A) func(B) A {
	return func(_ B) A {
		return a
	}
}

// Constant2 creates a binary function that returns the constant value 'a' and ignores its inputs
func Constant2[B, C, A any](a A) func(B, C) A {
	return func(_ B, _ C) A {
		return a
	}
}

func IsNil[A any](a *A) bool {
	return a == nil
}

func IsNonNil[A any](a *A) bool {
	return a != nil
}

// Swap returns a new binary function that changes the order of input parameters
func Swap[T1, T2, R any](f func(T1, T2) R) func(T2, T1) R {
	return func(t2 T2, t1 T1) R {
		return f(t1, t2)
	}
}

// First returns the first out of two input values
func First[T1, T2 any](t1 T1, _ T2) T1 {
	return t1
}

// Second returns the second out of two input values
func Second[T1, T2 any](_ T1, t2 T2) T2 {
	return t2
}

func Nullary1[R any](f1 func() R) func() R {
	return f1
}

func Nullary2[T1, R any](f1 func() T1, f2 func(t1 T1) R) func() R {
	return func() R {
		return Pipe1(f1(), f2)
	}
}

func Nullary3[T1, T2, R any](f1 func() T1, f2 func(t1 T1) T2, f3 func(t2 T2) R) func() R {
	return func() R {
		return Pipe2(f1(), f2, f3)
	}
}

func Nullary4[T1, T2, T3, R any](f1 func() T1, f2 func(t1 T1) T2, f3 func(t2 T2) T3, f4 func(t3 T3) R) func() R {
	return func() R {
		return Pipe3(f1(), f2, f3, f4)
	}
}
