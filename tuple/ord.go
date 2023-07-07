package tuple

import (
	O "github.com/ibm/fp-go/ord"
)

// Ord1 implements ordering on a 1-tuple
func Ord1[T1 any](o1 O.Ord[T1]) O.Ord[Tuple1[T1]] {
	return O.MakeOrd(func(l, r Tuple1[T1]) int {
		return o1.Compare(l.F1, r.F1)
	}, func(l, r Tuple1[T1]) bool {
		return o1.Equals(l.F1, r.F1)
	})
}

// Ord2 implements ordering on a 2-tuple
func Ord2[T1, T2 any](o1 O.Ord[T1], o2 O.Ord[T2]) O.Ord[Tuple2[T1, T2]] {
	return O.MakeOrd(func(l, r Tuple2[T1, T2]) int {
		c := o1.Compare(l.F1, r.F1)
		if c != 0 {
			return c
		}
		c = o2.Compare(l.F2, r.F2)
		return c
	}, func(l, r Tuple2[T1, T2]) bool {
		return o1.Equals(l.F1, r.F1) && o2.Equals(l.F2, r.F2)
	})
}

// Ord3 implements ordering on a 3-tuple
func Ord3[T1, T2, T3 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3]) O.Ord[Tuple3[T1, T2, T3]] {
	return O.MakeOrd(func(l, r Tuple3[T1, T2, T3]) int {
		c := o1.Compare(l.F1, r.F1)
		if c != 0 {
			return c
		}
		c = o2.Compare(l.F2, r.F2)
		if c != 0 {
			return c
		}
		c = o3.Compare(l.F3, r.F3)
		return c
	}, func(l, r Tuple3[T1, T2, T3]) bool {
		return o1.Equals(l.F1, r.F1) && o2.Equals(l.F2, r.F2) && o3.Equals(l.F3, r.F3)
	})
}

// Ord4 implements ordering on a 4-tuple
func Ord4[T1, T2, T3, T4 any](o1 O.Ord[T1], o2 O.Ord[T2], o3 O.Ord[T3], o4 O.Ord[T4]) O.Ord[Tuple4[T1, T2, T3, T4]] {
	return O.MakeOrd(func(l, r Tuple4[T1, T2, T3, T4]) int {
		c := o1.Compare(l.F1, r.F1)
		if c != 0 {
			return c
		}
		c = o2.Compare(l.F2, r.F2)
		if c != 0 {
			return c
		}
		c = o3.Compare(l.F3, r.F3)
		if c != 0 {
			return c
		}
		c = o4.Compare(l.F4, r.F4)
		return c
	}, func(l, r Tuple4[T1, T2, T3, T4]) bool {
		return o1.Equals(l.F1, r.F1) && o2.Equals(l.F2, r.F2) && o3.Equals(l.F3, r.F3) && o4.Equals(l.F4, r.F4)
	})
}
