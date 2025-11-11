package builder

import (
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Result[T any] = result.Result[T]

	Prism[S, A any] = prism.Prism[S, A]

	Option[T any] = option.Option[T]
)
