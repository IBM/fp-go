package lazy

type (
	// Lazy represents a synchronous computation without side effects
	Lazy[A any] = func() A

	Kleisli[A, B any]  = func(A) Lazy[B]
	Operator[A, B any] = Kleisli[Lazy[A], B]
)
