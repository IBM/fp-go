package tuple

import (
	M "github.com/ibm/fp-go/monoid"
)

// Monoid1 implements a monoid for a 1-tuple
func Monoid1[T1 any](m1 M.Monoid[T1]) M.Monoid[Tuple1[T1]] {
	return M.MakeMonoid(func(l, r Tuple1[T1]) Tuple1[T1] {
		return MakeTuple1(m1.Concat(l.F1, l.F1))
	}, MakeTuple1(m1.Empty()))
}

// Monoid2 implements a monoid for a 2-tuple
func Monoid2[T1, T2 any](m1 M.Monoid[T1], m2 M.Monoid[T2]) M.Monoid[Tuple2[T1, T2]] {
	return M.MakeMonoid(func(l, r Tuple2[T1, T2]) Tuple2[T1, T2] {
		return MakeTuple2(m1.Concat(l.F1, l.F1), m2.Concat(l.F2, l.F2))
	}, MakeTuple2(m1.Empty(), m2.Empty()))
}

// Monoid3 implements a monoid for a 3-tuple
func Monoid3[T1, T2, T3 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3]) M.Monoid[Tuple3[T1, T2, T3]] {
	return M.MakeMonoid(func(l, r Tuple3[T1, T2, T3]) Tuple3[T1, T2, T3] {
		return MakeTuple3(m1.Concat(l.F1, l.F1), m2.Concat(l.F2, l.F2), m3.Concat(l.F3, l.F3))
	}, MakeTuple3(m1.Empty(), m2.Empty(), m3.Empty()))
}

// Monoid4 implements a monoid for a 4-tuple
func Monoid4[T1, T2, T3, T4 any](m1 M.Monoid[T1], m2 M.Monoid[T2], m3 M.Monoid[T3], m4 M.Monoid[T4]) M.Monoid[Tuple4[T1, T2, T3, T4]] {
	return M.MakeMonoid(func(l, r Tuple4[T1, T2, T3, T4]) Tuple4[T1, T2, T3, T4] {
		return MakeTuple4(m1.Concat(l.F1, l.F1), m2.Concat(l.F2, l.F2), m3.Concat(l.F3, l.F3), m4.Concat(l.F4, l.F4))
	}, MakeTuple4(m1.Empty(), m2.Empty(), m3.Empty(), m4.Empty()))
}
