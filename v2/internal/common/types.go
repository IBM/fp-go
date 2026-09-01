package common

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/option"
)

type (
	// Endomorphism is a function from a type to itself (A → A).
	// It represents transformations that preserve the type.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Option[A any] = option.Option[A]
)
