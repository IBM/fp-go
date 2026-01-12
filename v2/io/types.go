package io

import (
	"iter"

	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	S "github.com/IBM/fp-go/v2/semigroup"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// IO represents a synchronous computation that cannot fail.
	// It's a function that takes no arguments and returns a value of type A.
	// Refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioltagt] for more details.
	IO[A any] = func() A

	// Pair represents a tuple of two values of types L and R.
	Pair[L, R any] = pair.Pair[L, R]

	// Kleisli represents a Kleisli arrow for the IO monad.
	// It's a function from A to IO[B], used for composing IO operations.
	Kleisli[A, B any] = reader.Reader[A, IO[B]]

	// Operator represents a function that transforms one IO into another.
	// It takes an IO[A] and produces an IO[B].
	Operator[A, B any] = Kleisli[IO[A], B]

	// Monoid represents a monoid structure for IO values.
	Monoid[A any] = M.Monoid[IO[A]]

	// Semigroup represents a semigroup structure for IO values.
	Semigroup[A any] = S.Semigroup[IO[A]]

	// Consumer represents a function that consumes a value of type A.
	Consumer[A any] = consumer.Consumer[A]

	// Seq represents an iterator sequence over values of type T.
	Seq[T any] = iter.Seq[T]

	// Trampoline represents a tail-recursive computation that can be evaluated safely
	// without stack overflow. It's used for implementing stack-safe recursive algorithms.
	Trampoline[B, L any] = tailrec.Trampoline[B, L]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]

	Void = function.Void
)
