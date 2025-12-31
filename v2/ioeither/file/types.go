package file

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/ioeither"
)

type (
	Either[E, T any]      = either.Either[E, T]
	IOEither[E, T any]    = ioeither.IOEither[E, T]
	Kleisli[E, A, B any]  = ioeither.Kleisli[E, A, B]
	Operator[E, A, B any] = ioeither.Operator[E, A, B]
)
