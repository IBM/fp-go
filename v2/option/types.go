package option

import (
	"iter"

	"github.com/IBM/fp-go/v2/endomorphism"
)

type (
	// Seq represents an iterator sequence over values of type T.
	// It's an alias for Go's standard iter.Seq[T] type.
	Seq[T any] = iter.Seq[T]

	// Endomorphism represents a function from a type to itself (T -> T).
	Endomorphism[T any] = endomorphism.Endomorphism[T]
)
