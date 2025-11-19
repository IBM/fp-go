package http

import (
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/readerioeither"
)

type (
	Kleisli[E, A, B any]        = ioeither.Kleisli[E, A, B]
	ReaderIOEither[R, E, A any] = readerioeither.ReaderIOEither[R, E, A]
)
