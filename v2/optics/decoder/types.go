package decoder

import (
	"github.com/IBM/fp-go/v2/result"
)

type (
	Result[A any] = result.Result[A]

	Decoder[I, A any] = result.Kleisli[I, A]
)
