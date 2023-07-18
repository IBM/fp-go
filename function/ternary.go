package function

func Ternary[A, B any](pred func(A) bool, onTrue func(A) B, onFalse func(A) B) func(A) B {
	return func(a A) B {
		if pred(a) {
			return onTrue(a)
		}
		return onFalse(a)
	}
}
