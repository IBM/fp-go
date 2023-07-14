package ioeither

import (
	ET "github.com/ibm/fp-go/either"
	G "github.com/ibm/fp-go/ioeither/generic"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[E, A, B, ANY any](
	acquire IOEither[E, A],
	use func(A) IOEither[E, B],
	release func(A, ET.Either[E, B]) IOEither[E, ANY],
) IOEither[E, B] {
	return G.Bracket(acquire, use, release)
}
