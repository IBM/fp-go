package either

import "fmt"

// Either defines a data structure that logically holds either an E or an A. The flag discriminates the cases
type (
	Either[E, A any] struct {
		isLeft bool
		left   E
		right  A
	}
)

// String prints some debug info for the object
func (s Either[E, A]) String() string {
	if s.isLeft {
		return fmt.Sprintf("Left[%T, %T](%v)", s.left, s.right, s.left)
	}
	return fmt.Sprintf("Right[%T, %T](%v)", s.left, s.right, s.right)
}

// Format prints some debug info for the object
func (s Either[E, A]) Format(f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, s.String())
	default:
		fmt.Fprint(f, s.String())
	}
}

func IsLeft[E, A any](val Either[E, A]) bool {
	return val.isLeft
}

func IsRight[E, A any](val Either[E, A]) bool {
	return !val.isLeft
}

func Left[E, A any](value E) Either[E, A] {
	return Either[E, A]{isLeft: true, left: value}
}

func Right[E, A any](value A) Either[E, A] {
	return Either[E, A]{isLeft: false, right: value}
}

func MonadFold[E, A, B any](ma Either[E, A], onLeft func(e E) B, onRight func(a A) B) B {
	if IsLeft(ma) {
		return onLeft(ma.left)
	}
	return onRight(ma.right)
}

// Unwrap converts an Either into the idiomatic tuple
func Unwrap[E, A any](ma Either[E, A]) (A, E) {
	return ma.right, ma.left
}
