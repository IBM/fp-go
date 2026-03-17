package readerreaderioresult

import (
	"github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/internal/witherable"
	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/IBM/fp-go/v2/option"
)

//go:inline
func Filter[C, HKTA, A any](
	filter func(Predicate[A]) Endomorphism[HKTA],
) func(Predicate[A]) Operator[C, HKTA, HKTA] {
	return witherable.Filter(
		Map[C],
		filter,
	)
}

//go:inline
func FilterArray[C, A any](p Predicate[A]) Operator[C, []A, []A] {
	return Filter[C](array.Filter[A])(p)
}

//go:inline
func FilterIter[C, A any](p Predicate[A]) Operator[C, Seq[A], Seq[A]] {
	return Filter[C](iter.Filter[A])(p)
}

//go:inline
func FilterMap[C, HKTA, HKTB, A, B any](
	filter func(option.Kleisli[A, B]) Reader[HKTA, HKTB],
) func(option.Kleisli[A, B]) Operator[C, HKTA, HKTB] {
	return witherable.FilterMap(
		Map[C],
		filter,
	)
}

//go:inline
func FilterMapArray[C, A, B any](p option.Kleisli[A, B]) Operator[C, []A, []B] {
	return FilterMap[C](array.FilterMap[A, B])(p)
}

//go:inline
func FilterMapIter[C, A, B any](p option.Kleisli[A, B]) Operator[C, Seq[A], Seq[B]] {
	return FilterMap[C](iter.FilterMap[A, B])(p)
}
