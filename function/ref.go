package function

func Ref[A any](a A) *A {
	return &a
}

func Deref[A any](a *A) A {
	return *a
}
