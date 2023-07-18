package option

import (
	F "github.com/ibm/fp-go/function"
)

func toType[T any](a any) (T, bool) {
	b, ok := a.(T)
	return b, ok
}

func ToType[T any](src any) Option[T] {
	return F.Pipe1(
		src,
		Optionize1(toType[T]),
	)
}

func ToAny[T any](src T) Option[any] {
	return Of(any(src))
}
