package ioeither

import (
	"iter"

	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Consumer represents a function that consumes a value of type A.
	// It's typically used for side effects like logging or updating state.
	Consumer[A any] = consumer.Consumer[A]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]

	// Trampoline represents a tail-recursive computation that can be evaluated safely
	// without stack overflow. It's used for implementing stack-safe recursive algorithms.
	Trampoline[B, L any] = tailrec.Trampoline[B, L]

	Void = function.Void

	Option[T any] = option.Option[T]

	Pair[A, B any] = pair.Pair[A, B]

	Seq[T any] = iter.Seq[T]
)
