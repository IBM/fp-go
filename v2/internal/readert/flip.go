package readert

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
)

func Sequence[
	HKTR2HKTR1A ~func(R2) HKTR1HKTA,
	R1, R2, HKTR1HKTA, HKTA any](
	mchain func(HKTR1HKTA, func(func(R1) HKTA) HKTA) HKTA,
	ma HKTR2HKTR1A,
) func(R1) func(R2) HKTA {
	return func(r1 R1) func(R2) HKTA {
		return func(r2 R2) HKTA {
			return mchain(
				ma(r2),
				identity.Ap[HKTA](r1),
			)
		}
	}
}

func SequenceReader[
	HKTR2HKTR1A ~func(R2) HKTR1HKTA,
	R1, R2, A, HKTR1HKTA, HKTA any](
	mmap func(HKTR1HKTA, func(func(R1) A) A) HKTA,
	ma HKTR2HKTR1A,
) func(R1) func(R2) HKTA {
	return func(r1 R1) func(R2) HKTA {
		return func(r2 R2) HKTA {
			return mmap(
				ma(r2),
				identity.Ap[A](r1),
			)
		}
	}
}

func Traverse[
	HKTR2A ~func(R2) HKTA,
	HKTR1B ~func(R1) HKTB,
	R1, R2, A, HKTR1HKTB, HKTA, HKTB any](
	mmap func(HKTA, func(A) HKTR1B) HKTR1HKTB,
	mchain func(HKTR1HKTB, func(func(R1) HKTB) HKTB) HKTB,
	f func(A) HKTR1B,
) func(HKTR2A) func(R1) func(R2) HKTB {
	return func(ma HKTR2A) func(R1) func(R2) HKTB {
		return func(r1 R1) func(R2) HKTB {
			return func(r2 R2) HKTB {
				return mchain(
					mmap(ma(r2), f),
					identity.Ap[HKTB](r1),
				)
			}
		}
	}
}

func TraverseReader[
	HKTR2A ~func(R2) HKTA,
	HKTR1B ~func(R1) B,
	R1, R2, A, B, HKTR1HKTB, HKTA, HKTB any](
	mmap1 func(func(A) HKTR1B) func(HKTA) HKTR1HKTB,
	mmap2 func(func(func(R1) B) B) func(HKTR1HKTB) HKTB,
	f func(A) HKTR1B,
) func(HKTR2A) func(R1) func(R2) HKTB {
	return function.Flow2(
		function.Bind1of2(function.Bind2of3(function.Flow3[HKTR2A, func(HKTA) HKTR1HKTB, func(HKTR1HKTB) HKTB])(mmap1(f))),
		function.Bind12of3(function.Flow3[func(fa R1) identity.Operator[func(R1) B, B], func(func(func(R1) B) B) func(HKTR1HKTB) HKTB, func(func(HKTR1HKTB) HKTB) func(R2) HKTB])(identity.Ap[B, R1], mmap2),
	)
}
