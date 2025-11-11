package ioresult

import (
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
)

type (
	IOResult[T any] = ioresult.IOResult[T]
	Result[T any]   = result.Result[T]
)
