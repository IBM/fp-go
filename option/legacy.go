package option

import "fmt"

type OptionTag int

const (
	NoneTag OptionTag = iota
	SomeTag
)

// Option defines a data structure that logically holds a value or not
type Option[A any] struct {
	Tag   OptionTag
	Value A
}

// String prints some debug info for the object
func (s Option[A]) String() string {
	switch s.Tag {
	case NoneTag:
		return fmt.Sprintf("None[%T]", s.Value)
	case SomeTag:
		return fmt.Sprintf("Some[%T](%v)", s.Value, s.Value)
	}
	return "Invalid"
}

func IsNone[T any](val Option[T]) bool {
	return val.Tag == NoneTag
}

func Some[T any](value T) Option[T] {
	return Option[T]{Tag: SomeTag, Value: value}
}

func Of[T any](value T) Option[T] {
	return Some(value)
}

func None[T any]() Option[T] {
	return Option[T]{Tag: NoneTag}
}

func IsSome[T any](val Option[T]) bool {
	return val.Tag == SomeTag
}

func MonadFold[A, B any](ma Option[A], onNone func() B, onSome func(A) B) B {
	if IsNone(ma) {
		return onNone()
	}
	return onSome(ma.Value)
}

func Unwrap[A any](ma Option[A]) (A, bool) {
	return ma.Value, ma.Tag == SomeTag
}
