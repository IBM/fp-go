package result

import (
	T "github.com/IBM/fp-go/v2/optics/traversal/generic"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Traversal represents an optic that focuses on zero or more values of type A within a structure S.
	// It's specialized for Result types, allowing traversal over successful values.
	Traversal[S, A any] = T.Traversal[S, A, Result[S], Result[A]]

	// Result represents a computation that may fail with an error.
	Result[T any] = result.Result[T]

	// Operator represents a function that transforms one Traversal into another.
	// It takes a Traversal[S, A] and produces a Traversal[S, B].
	Operator[S, A, B any] = func(Traversal[S, A]) Traversal[S, B]
)
