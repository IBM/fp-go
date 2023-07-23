package writer

import (
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/writer/generic"
)

// Constructs an equal predicate for a [Writer]
func Eq[W, A any](w EQ.Eq[W], a EQ.Eq[A]) EQ.Eq[Writer[W, A]] {
	return G.Eq[Writer[W, A]](w, a)
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[W, A comparable]() EQ.Eq[Writer[W, A]] {
	return G.FromStrictEquals[Writer[W, A]]()
}
