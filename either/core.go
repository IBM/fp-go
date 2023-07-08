package either

import "fmt"

type EitherTag string

const (
	leftTag  = "Left"
	rightTag = "Right"
)

// Either defines a data structure that logically holds either an E or an A. The tag discriminates the cases
type Either[E, A any] struct {
	tag   EitherTag `default:"Left"`
	left  E
	right A
}

// String prints some debug info for the object
func (s Either[E, A]) String() string {
	switch s.tag {
	case leftTag:
		return fmt.Sprintf("%s[%T, %T](%v)", s.tag, s.left, s.right, s.left)
	case rightTag:
		return fmt.Sprintf("%s[%T, %T](%v)", s.tag, s.left, s.right, s.right)
	}
	return "Invalid"
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
	return val.tag == leftTag
}

func IsRight[E, A any](val Either[E, A]) bool {
	return val.tag == rightTag
}

func Left[E, A any](value E) Either[E, A] {
	return Either[E, A]{tag: leftTag, left: value}
}

func Right[E, A any](value A) Either[E, A] {
	return Either[E, A]{tag: rightTag, right: value}
}

func MonadFold[E, A, B any](ma Either[E, A], onLeft func(e E) B, onRight func(a A) B) B {
	if IsLeft(ma) {
		return onLeft(ma.left)
	}
	return onRight(ma.right)
}

func Unwrap[E, A any](ma Either[E, A]) (A, E) {
	return ma.right, ma.left
}
