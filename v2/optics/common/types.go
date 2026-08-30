package common

import "github.com/IBM/fp-go/v2/endomorphism"

type (
	// Endomorphism is a function from a type to itself (A → A).
	// It represents transformations that preserve the type.
	Endomorphism[A any] = endomorphism.Endomorphism[A]
)
