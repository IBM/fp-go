package filterable

import (
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
)

type (
	Option[A any]       = option.Option[A]
	Separated[A, B any] = pair.Pair[A, B]

	FilterType[A, HKTA any] = func(func(A) bool) func(HKTA) HKTA

	FilterMapType[A, B, HKTA, HKTB any] = func(func(A) Option[B]) func(HKTA) HKTB
)
