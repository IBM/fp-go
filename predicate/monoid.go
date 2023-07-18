package predicate

import (
	F "github.com/IBM/fp-go/function"
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

// SemigroupAny combines predicates via ||
func SemigroupAny[A any](predicate func(A) bool) S.Semigroup[func(A) bool] {
	return S.MakeSemigroup(func(first func(A) bool, second func(A) bool) func(A) bool {
		return F.Pipe1(
			first,
			Or(second),
		)
	})
}

// SemigroupAll combines predicates via &&
func SemigroupAll[A any](predicate func(A) bool) S.Semigroup[func(A) bool] {
	return S.MakeSemigroup(func(first func(A) bool, second func(A) bool) func(A) bool {
		return F.Pipe1(
			first,
			And(second),
		)
	})
}

// MonoidAny combines predicates via ||
func MonoidAny[A any](predicate func(A) bool) S.Semigroup[func(A) bool] {
	return M.MakeMonoid(
		SemigroupAny(predicate).Concat,
		F.Constant1[A](false),
	)
}

// MonoidAll combines predicates via &&
func MonoidAll[A any](predicate func(A) bool) S.Semigroup[func(A) bool] {
	return M.MakeMonoid(
		SemigroupAll(predicate).Concat,
		F.Constant1[A](true),
	)
}
