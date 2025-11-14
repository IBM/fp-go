package result

import (
	T "github.com/IBM/fp-go/v2/optics/traversal/generic"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Traversal[S, A any]   = T.Traversal[S, A, Result[S], Result[A]]
	Result[T any]         = result.Result[T]
	Operator[S, A, B any] = func(Traversal[S, A]) Traversal[S, B]
)
