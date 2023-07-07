// Package tuple contains type definitions and functions for data structures for tuples of heterogenous types. For homogeneous types
// consider to use arrays for simplicity
package tuple

// Tuple1 is a structure carrying one element
type Tuple1[T1 any] struct {
	F1 T1
}

// Tuple2 is a structure carrying two elements
type Tuple2[T1, T2 any] struct {
	F1 T1
	F2 T2
}

// Tuple3 is a structure carrying three elements
type Tuple3[T1, T2, T3 any] struct {
	F1 T1
	F2 T2
	F3 T3
}

// Tuple4 is a structure carrying four elements
type Tuple4[T1, T2, T3, T4 any] struct {
	F1 T1
	F2 T2
	F3 T3
	F4 T4
}

func MakeTuple1[T1 any](t1 T1) Tuple1[T1] {
	return Tuple1[T1]{F1: t1}
}

func MakeTuple2[T1, T2 any](t1 T1, t2 T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2]{F1: t1, F2: t2}
}

func MakeTuple3[T1, T2, T3 any](t1 T1, t2 T2, t3 T3) Tuple3[T1, T2, T3] {
	return Tuple3[T1, T2, T3]{F1: t1, F2: t2, F3: t3}
}

func MakeTuple4[T1, T2, T3, T4 any](t1 T1, t2 T2, t3 T3, t4 T4) Tuple4[T1, T2, T3, T4] {
	return Tuple4[T1, T2, T3, T4]{F1: t1, F2: t2, F3: t3, F4: t4}
}

// Tupled2 converts a function that accepts two parameters into a function that accepts a tuple
func Tupled2[T1, T2, R any](f func(t1 T1, t2 T2) R) func(Tuple2[T1, T2]) R {
	return func(t Tuple2[T1, T2]) R {
		return f(t.F1, t.F2)
	}
}

// Tupled3 converts a function that accepts three parameters into a function that accepts a tuple
func Tupled3[T1, T2, T3, R any](f func(t1 T1, t2 T2, t3 T3) R) func(Tuple3[T1, T2, T3]) R {
	return func(t Tuple3[T1, T2, T3]) R {
		return f(t.F1, t.F2, t.F3)
	}
}

// Untupled2 converts a function that accepts a tuple into a function that accepts two parameters
func Untupled2[T1, T2, R any](f func(Tuple2[T1, T2]) R) func(T1, T2) R {
	return func(t1 T1, t2 T2) R {
		return f(MakeTuple2(t1, t2))
	}
}

// Untupled3 converts a function that accepts a tuple into a function that accepts three parameters
func Untupled3[T1, T2, T3, R any](f func(Tuple3[T1, T2, T3]) R) func(T1, T2, T3) R {
	return func(t1 T1, t2 T2, t3 T3) R {
		return f(MakeTuple3(t1, t2, t3))
	}
}

func First[T1, T2 any](t Tuple2[T1, T2]) T1 {
	return t.F1
}

func Second[T1, T2 any](t Tuple2[T1, T2]) T2 {
	return t.F2
}

func Swap[T1, T2 any](t Tuple2[T1, T2]) Tuple2[T2, T1] {
	return MakeTuple2(t.F2, t.F1)
}

func Of[T1, T2 any](e T2) func(T1) Tuple2[T1, T2] {
	return func(t T1) Tuple2[T1, T2] {
		return MakeTuple2(t, e)
	}
}

func BiMap[E, G, A, B any](mapSnd func(E) G, mapFst func(A) B) func(Tuple2[A, E]) Tuple2[B, G] {
	return func(t Tuple2[A, E]) Tuple2[B, G] {
		return MakeTuple2(mapFst(First(t)), mapSnd(Second(t)))
	}
}
