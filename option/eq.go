package option

import (
	EQ "github.com/ibm/fp-go/eq"
	F "github.com/ibm/fp-go/function"
)

// Constructs an equal predicate for an `Option`
func Eq[A any](a EQ.Eq[A]) EQ.Eq[Option[A]] {
	// some convenient shortcuts
	fld := Fold(
		F.Constant(Fold(F.ConstTrue, F.Constant1[A](false))),
		F.Flow2(F.Curry2(a.Equals), F.Bind1st(Fold[A, bool], F.ConstFalse)),
	)
	// convert to an equals predicate
	return EQ.FromEquals(F.Uncurry2(fld))
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[A comparable]() EQ.Eq[Option[A]] {
	return Eq(EQ.FromStrictEquals[A]())
}
