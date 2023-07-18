package either

import (
	EQ "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
)

// Constructs an equal predicate for an `Either`
func Eq[E, A any](e EQ.Eq[E], a EQ.Eq[A]) EQ.Eq[Either[E, A]] {
	// some convenient shortcuts
	eqa := F.Curry2(a.Equals)
	eqe := F.Curry2(e.Equals)

	fca := F.Bind2nd(Fold[E, A, bool], F.Constant1[A](false))
	fce := F.Bind1st(Fold[E, A, bool], F.Constant1[E](false))

	fld := Fold(F.Flow2(eqe, fca), F.Flow2(eqa, fce))

	return EQ.FromEquals(F.Uncurry2(fld))
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[E, A comparable]() EQ.Eq[Either[E, A]] {
	return Eq(EQ.FromStrictEquals[E](), EQ.FromStrictEquals[A]())
}
