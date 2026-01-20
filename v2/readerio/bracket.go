package readerio

import (
	G "github.com/IBM/fp-go/v2/internal/bracket"
)

//go:inline
func Bracket[
	R, A, B, ANY any](

	acquire ReaderIO[R, A],
	use Kleisli[R, A, B],
	release func(A, B) ReaderIO[R, ANY],
) ReaderIO[R, B] {
	return G.MonadBracket[
		ReaderIO[R, A],
		ReaderIO[R, B],
		ReaderIO[R, ANY],
		B,
		A,
		B,
	](
		Of[R, B],
		MonadChain[R, A, B],
		MonadChain[R, B, B],
		MonadChain[R, ANY, B],

		acquire,
		use,
		release,
	)
}
