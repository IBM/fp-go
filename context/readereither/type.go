// Package readereither implements a specialization of the Reader monad assuming a golang context as the context of the monad and a standard golang error
package readereither

import (
	"context"

	RE "github.com/ibm/fp-go/readereither"
)

// ReaderEither is a specialization of the Reader monad for the typical golang scenario
type ReaderEither[A any] RE.ReaderEither[context.Context, error, A]
