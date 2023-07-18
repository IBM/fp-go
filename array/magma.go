package array

import (
	F "github.com/IBM/fp-go/function"
	M "github.com/IBM/fp-go/magma"
)

func ConcatAll[A any](m M.Magma[A]) func(A) func([]A) A {
	return F.Bind1st(Reduce[A, A], m.Concat)
}
