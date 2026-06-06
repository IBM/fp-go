package ioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
)

// Unfold generates a lazy sequence of Either[E, A] values by repeatedly
// applying f to a seed of type B. f returns an IOEither that, on the right
// side, yields Some(Pair(nextSeed, value)) to continue or None to stop. The
// sequence terminates as soon as f returns None, the consumer stops iterating,
// or f returns a Left (error), in which case the error is emitted as the final
// element.
//
// Example:
//
//	// Count down from n, emitting each value
//	countDown := ioeither.Unfold(
//	    func(n int) IOEither[error, option.Option[pair.Pair[int, int]]] {
//	        return func() either.Either[error, option.Option[pair.Pair[int, int]]] {
//	            if n == 0 {
//	                return either.Of[error](option.None[pair.Pair[int, int]]())
//	            }
//	            return either.Of[error](option.Some(pair.MakePair(n-1, n)))
//	        }
//	    },
//	)(3)
//	// Iterating countDown yields: Right(3), Right(2), Right(1)
func Unfold[E, A, B any](
	f Kleisli[E, B, Option[Pair[B, A]]],
) iter.Kleisli[B, Either[E, A]] {
	return func(seed B) Seq[Either[E, A]] {
		return func(yield func(Either[E, A]) bool) {
			current := seed
			for {
				eoba := f(current)()
				oba, e := either.Unwrap(eoba)
				if either.IsLeft(eoba) {
					yield(either.Left[A](e))
					return
				}
				ba, o := option.Unwrap(oba)
				if !o || !yield(either.Of[E](pair.Tail(ba))) {
					return
				}
				current = pair.Head(ba)
			}
		}
	}
}
