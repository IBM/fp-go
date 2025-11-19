package file

import (
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
)

type (
	IOResult[T any]    = ioresult.IOResult[T]
	Kleisli[A, B any]  = ioresult.Kleisli[A, B]
	Operator[A, B any] = ioresult.Operator[A, B]
)
