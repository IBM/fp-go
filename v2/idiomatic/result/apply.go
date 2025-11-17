package result

type Apply[A, B any] interface {
	Functor[A, B]
	Ap(A, error) Operator[func(A) B, B]
}
