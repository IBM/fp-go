package option

func toType[T any](a any) (T, bool) {
	b, ok := a.(T)
	return b, ok
}

func ToType[T any](src any) Option[T] {
	return fromValidation(src, toType[T])
}

func ToAny[T any](src T) Option[any] {
	return Of(any(src))
}
