package array

import (
	F "github.com/ibm/fp-go/function"
	M "github.com/ibm/fp-go/magma"
)

func ConcatAll[A any](m M.Magma[A]) func(A) func([]A) A {
	return F.Bind1st(Reduce[A, A], m.Concat)
}
