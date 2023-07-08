package option

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type OptionTag string

const (
	noneTag = "None"
	someTag = "Some"
)

var (
	jsonNull = []byte("null")
)

// Option defines a data structure that logically holds a value or not
type Option[A any] struct {
	tag   OptionTag `default:"None"`
	value A
}

// String prints some debug info for the object
func (s Option[A]) String() string {
	switch s.tag {
	case noneTag:
		return fmt.Sprintf("%s[%T]", s.tag, s.value)
	case someTag:
		return fmt.Sprintf("%s[%T](%v)", s.tag, s.value, s.value)
	}
	return "Invalid"
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
	if s.tag == noneTag {
		return jsonNull, nil
	}
	return json.Marshal(s.value)
}

func (s Option[A]) UnmarshalJSON(data []byte) error {
	// decode the value
	if bytes.Equal(data, jsonNull) {
		s.tag = noneTag
		return nil
	}
	s.tag = someTag
	return json.Unmarshal(data, &s.value)
}

func IsNone[T any](val Option[T]) bool {
	return val.tag == noneTag
}

func Some[T any](value T) Option[T] {
	return Option[T]{tag: someTag, value: value}
}

func Of[T any](value T) Option[T] {
	return Some(value)
}

func None[T any]() Option[T] {
	return Option[T]{tag: noneTag}
}

func IsSome[T any](val Option[T]) bool {
	return val.tag == someTag
}

func MonadFold[A, B any](ma Option[A], onNone func() B, onSome func(A) B) B {
	if IsNone(ma) {
		return onNone()
	}
	return onSome(ma.value)
}

func Unwrap[A any](ma Option[A]) (A, bool) {
	return ma.value, ma.tag == someTag
}
