package generic

import (
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	Option[A any]    = option.Option[A]
	Lazy[A any]      = lazy.Lazy[A]
	Pair[L, R any]   = pair.Pair[L, R]
	Predicate[A any] = predicate.Predicate[A]
)
