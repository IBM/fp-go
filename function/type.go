package function

func ToAny[A any](a A) any {
	return any(a)
}
