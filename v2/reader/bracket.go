package reader

import (
	G "github.com/IBM/fp-go/v2/internal/bracket"
)

//go:inline
func Bracket[
	R, A, B, ANY any](

	acquire Reader[R, A],
	use Kleisli[R, A, B],
	release func(A, B) Reader[R, ANY],
) Reader[R, B] {
	return G.MonadBracket[
		Reader[R, A],
		Reader[R, B],
		Reader[R, ANY],
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
