package common

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	// Endomorphism is a function from a type to itself (A → A).
	// It represents transformations that preserve the type.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Predicate[A any] = predicate.Predicate[A]

	Pair[L, R any] = pair.Pair[L, R]
)
