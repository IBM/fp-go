package iso

import (
	"github.com/IBM/fp-go/v2/array/nonempty"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	// Number represents a numeric type constraint.
	Number = number.Number

	// Pair represents a tuple of two values of types A and B.
	Pair[A, B any] = pair.Pair[A, B]

	// Either represents a value of one of two possible types (a disjoint union).
	Either[E, A any] = either.Either[E, A]

	// NonEmptyArray represents an array that is guaranteed to have at least one element.
	NonEmptyArray[A any] = nonempty.NonEmptyArray[A]

	// Reader represents a computation that reads from a shared environment of type R
	// and produces a value of type T.  It is a plain function func(R) T, re-exported
	// here so that callers of the iso package do not need to import the reader package
	// directly.
	//
	// Type Parameters:
	//   - R: The environment/context type supplied to the computation
	//   - T: The result type produced by the computation
	Reader[R, T any] = reader.Reader[R, T]

	// Eq represents an equality type class for type T.
	// It provides a way to define custom equality semantics for any type,
	// not just those that are comparable with Go's == operator.
	//
	// Type Parameters:
	//   - T: The type for which equality is defined
	//
	// An Eq instance must satisfy the equivalence relation laws:
	//   - Reflexivity:  Equals(x, x) == true for all x
	//   - Symmetry:     Equals(x, y) == Equals(y, x) for all x, y
	//   - Transitivity: if Equals(x, y) && Equals(y, z) then Equals(x, z)
	//
	// See Also:
	//   - FromEquals: creates a Kleisli arrow that builds an Iso[T, bool] from an Eq[T]
	//   - FromStrictEquals: convenience variant using Go's == operator
	Eq[T any] = eq.Eq[T]

	// Kleisli represents a Kleisli arrow for the Reader monad specialised to
	// isomorphisms: a function that accepts an environment of type R and returns
	// an Iso[S, A].
	//
	// In the iso package the environment is typically an Eq[T] instance, making
	// the full type Kleisli[Eq[T], T, bool] = func(Eq[T]) Iso[T, bool].  This
	// lets callers defer the choice of equality predicate to call time.
	//
	// Type Parameters:
	//   - R: The environment type consumed to produce the isomorphism
	//   - S: The source type of the resulting Iso
	//   - A: The target type of the resulting Iso
	//
	// See Also:
	//   - FromEquals: the primary function that returns a Kleisli in this package
	Kleisli[R, S, A any] = Reader[R, Iso[S, A]]

	Lazy[T any] = lazy.Lazy[T]
)
