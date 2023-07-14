// Package readerioeither implements a specialization of the Reader monad assuming a golang context as the context of the monad and a standard golang error
package readerioeither

import (
	"context"

	RE "github.com/ibm/fp-go/readerioeither"
)

// ReaderIOEither is a specialization of the Reader monad for the typical golang scenario
type ReaderIOEither[A any] RE.ReaderIOEither[context.Context, error, A]
