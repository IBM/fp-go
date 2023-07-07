package Apply

import (
	F "github.com/ibm/fp-go/function"
	T "github.com/ibm/fp-go/tuple"
)

func tupleConstructor1[A any]() func(A) T.Tuple1[A] {
	return F.Curry1(T.MakeTuple1[A])
}

func tupleConstructor2[A, B any]() func(A) func(B) T.Tuple2[A, B] {
	return F.Curry2(T.MakeTuple2[A, B])
}

func tupleConstructor3[A, B, C any]() func(A) func(B) func(C) T.Tuple3[A, B, C] {
	return F.Curry3(T.MakeTuple3[A, B, C])
}

func tupleConstructor4[A, B, C, D any]() func(A) func(B) func(C) func(D) T.Tuple4[A, B, C, D] {
	return F.Curry4(T.MakeTuple4[A, B, C, D])
}

func SequenceT1[A, HKTA, HKT1A any](
	fmap func(HKTA, func(A) T.Tuple1[A]) HKT1A,
	a HKTA) HKT1A {
	return fmap(a, tupleConstructor1[A]())
}

// HKTA = HKT[A]
// HKTB = HKT[B]
// HKT2AB = HKT[Tuple[A, B]]
// HKTFB2AB = HKT[func(B)Tuple[A, B]]
func SequenceT2[A, B, HKTA, HKTB, HKTFB2AB, HKT2AB any](
	fmap func(HKTA, func(A) func(B) T.Tuple2[A, B]) HKTFB2AB,
	fap1 func(HKTFB2AB, HKTB) HKT2AB,
	a HKTA, b HKTB,
) HKT2AB {
	return fap1(fmap(a, tupleConstructor2[A, B]()), b)
}

// HKTA = HKT[A]
// HKTB = HKT[B]
// HKTC = HKT[C]
// HKT3ABC = HKT[Tuple[A, B, C]]
// HKTFB3ABC = HKT[func(B)func(C)Tuple[A, B, C]]
// HKTFC3ABC = HKT[func(C)Tuple[A, B, C]]
func SequenceT3[A, B, C, HKTA, HKTB, HKTC, HKTFB3ABC, HKTFC3ABC, HKT3ABC any](
	fmap func(HKTA, func(A) func(B) func(C) T.Tuple3[A, B, C]) HKTFB3ABC,
	fap1 func(HKTFB3ABC, HKTB) HKTFC3ABC,
	fap2 func(HKTFC3ABC, HKTC) HKT3ABC,

	a HKTA, b HKTB, c HKTC) HKT3ABC {
	return fap2(fap1(fmap(a, tupleConstructor3[A, B, C]()), b), c)
}

// HKTA = HKT[A]
// HKTB = HKT[B]
// HKTC = HKT[C]
// HKT3ABCD = HKT[Tuple[A, B, C, D]]
// HKTFB3ABCD = HKT[func(B)func(C)func(D)Tuple[A, B, C, D]]
// HKTFC3ABCD = HKT[func(C)func(D)Tuple[A, B, C, D]]
// HKTFD3ABCD = HKT[func(D)Tuple[A, B, C, D]]
func SequenceT4[A, B, C, D, HKTA, HKTB, HKTC, HKTD, HKTFB4ABCD, HKTFC4ABCD, HKTFD4ABCD, HKT4ABCD any](
	fmap func(HKTA, func(A) func(B) func(C) func(D) T.Tuple4[A, B, C, D]) HKTFB4ABCD,
	fap1 func(HKTFB4ABCD, HKTB) HKTFC4ABCD,
	fap2 func(HKTFC4ABCD, HKTC) HKTFD4ABCD,
	fap3 func(HKTFD4ABCD, HKTD) HKT4ABCD,

	a HKTA, b HKTB, c HKTC, d HKTD) HKT4ABCD {
	return fap3(fap2(fap1(fmap(a, tupleConstructor4[A, B, C, D]()), b), c), d)
}
