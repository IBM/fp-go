package result

type Applicative[A, B any] interface {
	Apply[A, B]
	Pointed[A]
}
