// Package readerio implements a specialization of the ReaderIO monad assuming a golang context as the context of the monad
package readerio

import (
	"context"

	R "github.com/ibm/fp-go/readerio"
)

// ReaderIO is a specialization of the ReaderIO monad assuming a golang context as the context of the monad
type ReaderIO[A any] R.ReaderIO[context.Context, A]
