package array

import "github.com/IBM/fp-go/v2/option"

type (
	Kleisli[A, B any]  = func(A) []B
	Operator[A, B any] = Kleisli[[]A, B]
	Option[A any]      = option.Option[A]
)
