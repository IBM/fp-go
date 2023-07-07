package either

import "fmt"

type EitherTag int

const (
	LeftTag EitherTag = iota
	RightTag
)

// Either defines a data structure that logically holds either an E or an A. The tag discriminates the cases
type Either[E, A any] struct {
	Tag   EitherTag
	Left  E
	Right A
}

// String prints some debug info for the object
func (s Either[E, A]) String() string {
	switch s.Tag {
	case LeftTag:
		return fmt.Sprintf("Left[%T, %T](%v)", s.Left, s.Right, s.Left)
	case RightTag:
		return fmt.Sprintf("Right[%T, %T](%v)", s.Left, s.Right, s.Right)
	}
	return "Invalid"
}

func IsLeft[E, A any](val Either[E, A]) bool {
	return val.Tag == LeftTag
}

func IsRight[E, A any](val Either[E, A]) bool {
	return val.Tag == RightTag
}

func Left[E, A any](value E) Either[E, A] {
	return Either[E, A]{Tag: LeftTag, Left: value}
}

func Right[E, A any](value A) Either[E, A] {
	return Either[E, A]{Tag: RightTag, Right: value}
}

func fold[E, A, B any](ma Either[E, A], onLeft func(e E) B, onRight func(a A) B) B {
	if IsLeft(ma) {
		return onLeft(ma.Left)
	}
	return onRight(ma.Right)
}

func Unwrap[E, A any](ma Either[E, A]) (A, E) {
	return ma.Right, ma.Left
}
