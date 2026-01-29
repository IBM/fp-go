package ioresult

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// IO represents a computation that performs side effects and returns a value of type A.
	IO[A any] = io.IO[A]

	// Lazy represents a deferred computation that produces a value of type A when evaluated.
	Lazy[A any] = lazy.Lazy[A]

	// Result represents an Either with error as the left type, compatible with Go's (value, error) tuple.
	Result[A any] = result.Result[A]

	// Reader represents a computation that depends on a read-only environment of type R and produces a value of type A.
	Reader[R, A any] = reader.Reader[R, A]

	// Endomorphism represents a function from type A to type A.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// IOResult represents a computation that performs IO and may fail with an error.
	// It follows Go's idiomatic pattern of returning (value, error) tuples.
	// A successful computation returns (value, nil), while a failed one returns (zero, error).
	IOResult[A any] = func() (A, error)

	// Kleisli represents a function from A to an IOResult of B.
	// It is used for chaining computations that may fail.
	Kleisli[A, B any] = Reader[A, IOResult[B]]

	// Operator represents a transformation from IOResult[A] to IOResult[B].
	// It is commonly used in function composition pipelines.
	Operator[A, B any] = Kleisli[IOResult[A], B]

	Predicate[A any] = predicate.Predicate[A]

	Void = function.Void
)
