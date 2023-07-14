package generic

import (
	ET "github.com/ibm/fp-go/either"
	G "github.com/ibm/fp-go/internal/file"
	I "github.com/ibm/fp-go/io/generic"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[
	GA ~func() ET.Either[E, A],
	GB ~func() ET.Either[E, B],
	GANY ~func() ET.Either[E, ANY],
	E, A, B, ANY any](

	acquire GA,
	use func(A) GB,
	release func(A, ET.Either[E, B]) GANY,
) GB {
	return G.Bracket[GA, GB, GANY, ET.Either[E, B], A, B](
		I.Of[GB, ET.Either[E, B]],
		MonadChain[GA, GB, E, A, B],
		I.MonadChain[GB, GB, ET.Either[E, B], ET.Either[E, B]],
		MonadChain[GANY, GB, E, ANY, B],

		acquire,
		use,
		release,
	)
}
