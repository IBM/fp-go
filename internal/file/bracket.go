package file

import (
	F "github.com/IBM/fp-go/function"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[
	GA, // IOEither[E, A]
	GB, // IOEither[E, A]
	GANY, // IOEither[E, ANY]

	EB, // Either[E, B]

	A, B, ANY any](

	ofeb func(EB) GB,

	chainab func(GA, func(A) GB) GB,
	chainebb func(GB, func(EB) GB) GB,
	chainany func(GANY, func(ANY) GB) GB,

	acquire GA,
	use func(A) GB,
	release func(A, EB) GANY,
) GB {
	return chainab(acquire,
		func(a A) GB {
			return chainebb(use(a), func(eb EB) GB {
				return chainany(
					release(a, eb),
					F.Constant1[ANY](ofeb(eb)),
				)
			})
		})
}
