package predicate

func Not[A any](predicate func(A) bool) func(A) bool {
	return func(a A) bool {
		return !predicate((a))
	}
}

// And creates a predicate that combines other predicates via &&
func And[A any](second func(A) bool) func(func(A) bool) func(A) bool {
	return func(first func(A) bool) func(A) bool {
		return func(a A) bool {
			return first(a) && second(a)
		}
	}
}

// Or creates a predicate that combines other predicates via ||
func Or[A any](second func(A) bool) func(func(A) bool) func(A) bool {
	return func(first func(A) bool) func(A) bool {
		return func(a A) bool {
			return first(a) || second(a)
		}
	}
}
