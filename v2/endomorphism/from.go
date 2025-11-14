package endomorphism

import (
	"github.com/IBM/fp-go/v2/function"
	S "github.com/IBM/fp-go/v2/semigroup"
)

func FromSemigroup[A any](s S.Semigroup[A]) Kleisli[A] {
	return function.Bind2of2(s.Concat)
}
