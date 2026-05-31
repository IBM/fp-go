package io

import (
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
)

// Unfold generates a lazy sequence of A values by repeatedly applying f to a
// seed of type B. f returns an IO that yields Some(Pair(nextSeed, value)) to
// continue or None to stop. The sequence terminates as soon as f returns None
// or the consumer stops iterating.
func Unfold[A, B any](
	f Kleisli[B, Option[Pair[B, A]]],
	seed B,
) Seq[A] {
	return func(yield func(A) bool) {
		current := seed
		for {
			pba, ok := option.Unwrap(f(current)())
			if !ok || !yield(pair.Tail(pba)) {
				return
			}
			current = pair.Head(pba)
		}
	}
}
