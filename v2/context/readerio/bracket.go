package readerio

import (
	RIO "github.com/IBM/fp-go/v2/readerio"
)

//go:inline
func Bracket[
	A, B, ANY any](

	acquire ReaderIO[A],
	use Kleisli[A, B],
	release func(A, B) ReaderIO[ANY],
) ReaderIO[B] {
	return RIO.Bracket(acquire, use, release)
}
