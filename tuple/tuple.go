// Package tuple contains type definitions and functions for data structures for tuples of heterogenous types. For homogeneous types
// consider to use arrays for simplicity
package tuple

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
