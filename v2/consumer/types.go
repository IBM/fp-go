package consumer

type (
	Consumer[A any] = func(A)
)
