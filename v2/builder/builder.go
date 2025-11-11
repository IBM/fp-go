package builder

type (
	Builder[T any] interface {
		Build() Result[T]
	}
)
