package option

import (
	"iter"

	"github.com/IBM/fp-go/v2/endomorphism"
)

type (
	Seq[T any]          = iter.Seq[T]
	Endomorphism[T any] = endomorphism.Endomorphism[T]
)
