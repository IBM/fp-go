package readerioeither

import (
	ET "github.com/ibm/fp-go/either"
	G "github.com/ibm/fp-go/readerioeither/generic"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[
	R, E, A, B, ANY any](

	acquire ReaderIOEither[R, E, A],
	use func(A) ReaderIOEither[R, E, B],
	release func(A, ET.Either[E, B]) ReaderIOEither[R, E, ANY],
) ReaderIOEither[R, E, B] {
	return G.Bracket(acquire, use, release)
}
