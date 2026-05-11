package itereither

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]

	// Trampoline represents a tail-recursive computation that can be evaluated safely
	// without stack overflow. It's used for implementing stack-safe recursive algorithms.
	Trampoline[B, L any] = tailrec.Trampoline[B, L]

	Void = function.Void
)
