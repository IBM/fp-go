package decode

import (
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
)

// Of creates a Decode that always succeeds with the given value.
// This is the pointed functor operation that lifts a pure value into the Decode context.
//
// Example:
//
//	decoder := decode.Of[string](42)
//	result := decoder("any input") // Always returns validation.Success(42)
func Of[I, A any](a A) Decode[I, A] {
	return reader.Of[I](validation.Of(a))
}

// MonadChain sequences two decode operations, passing the result of the first to the second.
// This is the monadic bind operation that enables sequential composition of decoders.
//
// Example:
//
//	decoder1 := decode.Of[string](42)
//	decoder2 := decode.MonadChain(decoder1, func(n int) Decode[string, string] {
//	    return decode.Of[string](fmt.Sprintf("Number: %d", n))
//	})
func MonadChain[I, A, B any](fa Decode[I, A], f Kleisli[I, A, B]) Decode[I, B] {
	return readert.MonadChain(
		validation.MonadChain,
		fa,
		f,
	)
}

// Chain creates an operator that sequences decode operations.
// This is the curried version of MonadChain, useful for composition pipelines.
//
// Example:
//
//	chainOp := decode.Chain(func(n int) Decode[string, string] {
//	    return decode.Of[string](fmt.Sprintf("Number: %d", n))
//	})
//	decoder := chainOp(decode.Of[string](42))
func Chain[I, A, B any](f Kleisli[I, A, B]) Operator[I, A, B] {
	return readert.Chain[Decode[I, A]](
		validation.Chain,
		f,
	)
}

func ChainLeft[I, A any](f Kleisli[I, Errors, A]) Operator[I, A, A] {
	return readert.Chain[Decode[I, A]](
		validation.ChainLeft,
		f,
	)
}

func OrElse[I, A any](f Kleisli[I, Errors, A]) Operator[I, A, A] {
	return ChainLeft(f)
}

// MonadMap transforms the decoded value using the provided function.
// This is the functor map operation that applies a transformation to successful decode results.
//
// Example:
//
//	decoder := decode.Of[string](42)
//	mapped := decode.MonadMap(decoder, func(n int) string {
//	    return fmt.Sprintf("Number: %d", n)
//	})
func MonadMap[I, A, B any](fa Decode[I, A], f func(A) B) Decode[I, B] {
	return readert.MonadMap[
		Decode[I, A],
		Decode[I, B]](
		validation.MonadMap,
		fa,
		f,
	)
}

// Map creates an operator that transforms decoded values.
// This is the curried version of MonadMap, useful for composition pipelines.
//
// Example:
//
//	mapOp := decode.Map(func(n int) string {
//	    return fmt.Sprintf("Number: %d", n)
//	})
//	decoder := mapOp(decode.Of[string](42))
func Map[I, A, B any](f func(A) B) Operator[I, A, B] {
	return readert.Map[
		Decode[I, A],
		Decode[I, B]](
		validation.Map,
		f,
	)
}

// MonadAp applies a decoder containing a function to a decoder containing a value.
// This is the applicative apply operation that enables parallel composition of decoders.
//
// Example:
//
//	decoderFn := decode.Of[string](func(n int) string {
//	    return fmt.Sprintf("Number: %d", n)
//	})
//	decoderVal := decode.Of[string](42)
//	result := decode.MonadAp(decoderFn, decoderVal)
func MonadAp[B, I, A any](fab Decode[I, func(A) B], fa Decode[I, A]) Decode[I, B] {
	return readert.MonadAp[
		Decode[I, A],
		Decode[I, B],
		Decode[I, func(A) B], I, A](
		validation.MonadAp[B, A],
		fab,
		fa,
	)
}

// Ap creates an operator that applies a function decoder to a value decoder.
// This is the curried version of MonadAp, useful for composition pipelines.
//
// Example:
//
//	apOp := decode.Ap[string](decode.Of[string](42))
//	decoderFn := decode.Of[string](func(n int) string {
//	    return fmt.Sprintf("Number: %d", n)
//	})
//	result := apOp(decoderFn)
func Ap[B, I, A any](fa Decode[I, A]) Operator[I, func(A) B, B] {
	return readert.Ap[
		Decode[I, A],
		Decode[I, B],
		Decode[I, func(A) B], I, A](
		validation.Ap[B, A],
		fa,
	)
}
