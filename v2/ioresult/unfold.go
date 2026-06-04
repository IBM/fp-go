package ioresult

import (
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/iterator/iter"
)

// Unfold generates a lazy [SeqResult] of values by repeatedly applying f to a
// seed of type B. It is the dual of a fold (an anamorphism).
//
// On each iteration f is called with the current seed B:
//   - [option.None] terminates the sequence normally (no more values).
//   - [option.Some](Pair(nextSeed, value)) emits value (type A) as a [result.Right]
//     and advances the seed to nextSeed for the next iteration.
//   - A [result.Left] error terminates the sequence and propagates that error as
//     the final element.
//
// The consumer may stop early by breaking out of the range loop; the generator
// is not called again after that.
//
// Example — count down from n to 1:
//
//	countDown := ioresult.Unfold(
//	    func(n int) IOResult[Option[Pair[int, int]]] {
//	        if n == 0 {
//	            return Of[Option[Pair[int, int]]](option.None[Pair[int, int]]())
//	        }
//	        // Head = next seed, Tail = emitted value
//	        return Of(option.Some(pair.MakePair(n-1, n)))
//	    },
//	)(3)
//	// Iterating countDown yields: Right(3), Right(2), Right(1)
func Unfold[A, B any](
	f Kleisli[B, Option[Pair[B, A]]],
) iter.Kleisli[B, Result[A]] {
	return ioeither.Unfold(f)
}
