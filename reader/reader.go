package reader

import (
	G "github.com/IBM/fp-go/reader/generic"
	T "github.com/IBM/fp-go/tuple"
)

// The purpose of the `Reader` monad is to avoid threading arguments through multiple functions in order to only get them where they are needed.
// The first template argument `R` is the the context to read from, the second argument `A` is the return value of the monad
type Reader[R, A any] func(R) A

// MakeReader creates a reader, i.e. a method that accepts a context and that returns a value
func MakeReader[R, A any](r Reader[R, A]) Reader[R, A] {
	return G.MakeReader(r)
}

// Ask reads the current context
func Ask[R any]() Reader[R, R] {
	return G.Ask[Reader[R, R]]()
}

// Asks projects a value from the global context in a Reader
func Asks[R, A any](f Reader[R, A]) Reader[R, A] {
	return G.Asks(f)
}

func AsksReader[R, A any](f func(R) Reader[R, A]) Reader[R, A] {
	return G.AsksReader(f)
}

func MonadMap[E, A, B any](fa Reader[E, A], f func(A) B) Reader[E, B] {
	return G.MonadMap[Reader[E, A], Reader[E, B]](fa, f)
}

// Map can be used to turn functions `func(A)B` into functions `(fa F[A])F[B]` whose argument and return types
// use the type constructor `F` to represent some computational context.
func Map[E, A, B any](f func(A) B) func(Reader[E, A]) Reader[E, B] {
	return G.Map[Reader[E, A], Reader[E, B]](f)
}

func MonadAp[B, R, A any](fab Reader[R, func(A) B], fa Reader[R, A]) Reader[R, B] {
	return G.MonadAp[Reader[R, A], Reader[R, B]](fab, fa)
}

// Ap applies a function to an argument under a type constructor.
func Ap[B, R, A any](fa Reader[R, A]) func(Reader[R, func(A) B]) Reader[R, B] {
	return G.Ap[Reader[R, A], Reader[R, B], Reader[R, func(A) B]](fa)
}

func Of[R, A any](a A) Reader[R, A] {
	return G.Of[Reader[R, A]](a)
}

func MonadChain[R, A, B any](ma Reader[R, A], f func(A) Reader[R, B]) Reader[R, B] {
	return G.MonadChain(ma, f)
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[R, A, B any](f func(A) Reader[R, B]) func(Reader[R, A]) Reader[R, B] {
	return G.Chain[Reader[R, A]](f)
}

func Flatten[R, A any](mma func(R) Reader[R, A]) Reader[R, A] {
	return G.Flatten(mma)
}

func Compose[R, B, C any](ab Reader[R, B]) func(Reader[B, C]) Reader[R, C] {
	return G.Compose[Reader[R, B], Reader[B, C], Reader[R, C]](ab)
}

func First[A, B, C any](pab Reader[A, B]) Reader[T.Tuple2[A, C], T.Tuple2[B, C]] {
	return G.First[Reader[A, B], Reader[T.Tuple2[A, C], T.Tuple2[B, C]]](pab)
}

func Second[A, B, C any](pbc Reader[B, C]) Reader[T.Tuple2[A, B], T.Tuple2[A, C]] {
	return G.Second[Reader[B, C], Reader[T.Tuple2[A, B], T.Tuple2[A, C]]](pbc)
}

func Promap[E, A, D, B any](f func(D) E, g func(A) B) func(Reader[E, A]) Reader[D, B] {
	return G.Promap[Reader[E, A], Reader[D, B]](f, g)
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[R2, R1, A any](f func(R2) R1) func(Reader[R1, A]) Reader[R2, A] {
	return G.Local[Reader[R1, A], Reader[R2, A]](f)
}

// Read applies a context to a reader to obtain its value
func Read[E, A any](e E) func(Reader[E, A]) A {
	return G.Read[Reader[E, A]](e)
}
