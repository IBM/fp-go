package ioresult

import (
	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/semigroup"
)

type (
	// IO represents a synchronous computation that cannot fail.
	IO[A any] = io.IO[A]

	// Lazy represents a deferred computation that produces a value of type A.
	Lazy[A any] = lazy.Lazy[A]

	// Either represents a value of one of two possible types (a disjoint union).
	Either[E, A any] = either.Either[E, A]

	// Result represents a computation that may fail with an error.
	// It's an alias for Either[error, A].
	Result[A any] = result.Result[A]

	// Endomorphism represents a function from a type to itself (A -> A).
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// IOResult represents a synchronous computation that may fail with an error.
	// It combines IO (side effects) with Result (error handling).
	// Refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioeitherlte-agt] for more details.
	IOResult[A any] = IO[Result[A]]

	// Monoid represents a monoid structure for IOResult values.
	Monoid[A any] = monoid.Monoid[IOResult[A]]

	// Semigroup represents a semigroup structure for IOResult values.
	Semigroup[A any] = semigroup.Semigroup[IOResult[A]]

	// Kleisli represents a Kleisli arrow for the IOResult monad.
	// It's a function from A to IOResult[B], used for composing operations that may fail.
	Kleisli[A, B any] = reader.Reader[A, IOResult[B]]

	// Operator represents a function that transforms one IOResult into another.
	// It takes an IOResult[A] and produces an IOResult[B].
	Operator[A, B any] = Kleisli[IOResult[A], B]

	// Consumer represents a function that consumes a value of type A.
	// It's typically used for side effects like logging or updating state.
	Consumer[A any] = consumer.Consumer[A]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]

	Void = function.Void
)
