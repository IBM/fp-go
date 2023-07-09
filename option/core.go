package option

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var (
	jsonNull = []byte("null")
)

// Option defines a data structure that logically holds a value or not
type Option[A any] struct {
	isSome bool
	some   A
}

// String prints some debug info for the object
func (s Option[A]) String() string {
	if s.isSome {
		return fmt.Sprintf("Some[%T](%v)", s.some, s.some)
	}
	return fmt.Sprintf("None[%T]", s.some)
}

// Format prints some debug info for the object
func (s Option[A]) Format(f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, s.String())
	default:
		fmt.Fprint(f, s.String())
	}
}

func (s Option[A]) MarshalJSON() ([]byte, error) {
	if IsSome(s) {
		return json.Marshal(s.some)
	}
	return jsonNull, nil
}

func (s *Option[A]) UnmarshalJSON(data []byte) error {
	// decode the value
	if bytes.Equal(data, jsonNull) {
		s.isSome = false
		s.some = *new(A)
		return nil
	}
	s.isSome = true
	return json.Unmarshal(data, &s.some)
}

func IsNone[T any](val Option[T]) bool {
	return !val.isSome
}

func Some[T any](value T) Option[T] {
	return Option[T]{isSome: true, some: value}
}

func Of[T any](value T) Option[T] {
	return Some(value)
}

func None[T any]() Option[T] {
	return Option[T]{isSome: false}
}

func IsSome[T any](val Option[T]) bool {
	return val.isSome
}

func MonadFold[A, B any](ma Option[A], onNone func() B, onSome func(A) B) B {
	if IsSome(ma) {
		return onSome(ma.some)
	}
	return onNone()
}

func Unwrap[A any](ma Option[A]) (A, bool) {
	return ma.some, ma.isSome
}
