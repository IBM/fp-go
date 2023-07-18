// Package reader implements a specialization of the Reader monad assuming a golang context as the context of the monad
package reader

import (
	"context"

	R "github.com/IBM/fp-go/reader"
)

// Reader is a specialization of the Reader monad assuming a golang context as the context of the monad
type Reader[A any] R.Reader[context.Context, A]
