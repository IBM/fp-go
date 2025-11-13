package builder

import (
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Option[T any]  = option.Option[T]
	Result[T any]  = result.Result[T]
	Lens[S, T any] = lens.Lens[S, T]
)
