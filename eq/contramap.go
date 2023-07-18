package eq

// Contramap implements an Equals predicate based on a mapping
func Contramap[A, B any](f func(b B) A) func(Eq[A]) Eq[B] {
	return func(fa Eq[A]) Eq[B] {
		equals := fa.Equals
		return FromEquals(func(x, y B) bool {
			return equals(f(x), f(y))
		})
	}
}
