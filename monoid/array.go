package monoid

import (
	S "github.com/ibm/fp-go/semigroup"
)

func GenericConcatAll[GA ~[]A, A any](m Monoid[A]) func(GA) A {
	return S.GenericConcatAll[GA](S.MakeSemigroup(m.Concat))(m.Empty())
}

// ConcatAll concatenates all values using the monoid and the default empty value
func ConcatAll[A any](m Monoid[A]) func([]A) A {
	return GenericConcatAll[[]A](m)
}

// Fold concatenates all values using the monoid and the default empty value
func Fold[A any](m Monoid[A]) func([]A) A {
	return GenericConcatAll[[]A](m)
}
