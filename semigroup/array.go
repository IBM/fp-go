package semigroup

import (
	M "github.com/ibm/fp-go/magma"
)

func GenericMonadConcatAll[GA ~[]A, A any](s Semigroup[A]) func(GA, A) A {
	return M.GenericMonadConcatAll[GA](M.MakeMagma(s.Concat))
}

func GenericConcatAll[GA ~[]A, A any](s Semigroup[A]) func(A) func(GA) A {
	return M.GenericConcatAll[GA](M.MakeMagma(s.Concat))
}

func MonadConcatAll[A any](s Semigroup[A]) func([]A, A) A {
	return GenericMonadConcatAll[[]A](s)
}

func ConcatAll[A any](s Semigroup[A]) func(A) func([]A) A {
	return GenericConcatAll[[]A](s)
}
