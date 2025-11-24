package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Result[T any]      = result.Result[T]
	Reader             = reader.Reader[*testing.T, bool]
	Kleisli[T any]     = reader.Reader[T, Reader]
	Predicate[T any]   = predicate.Predicate[T]
	Lens[S, T any]     = lens.Lens[S, T]
	Optional[S, T any] = optional.Optional[S, T]
	Prism[S, T any]    = prism.Prism[S, T]
)
