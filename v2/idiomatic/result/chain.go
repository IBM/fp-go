package result

type Chainable[A, B any] interface {
	Apply[A, B]
	Chain(Kleisli[A, B]) Operator[A, B]
}
