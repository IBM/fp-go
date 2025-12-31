package ioresult

import (
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// IOResult represents a synchronous computation that may fail with an error.
	// It's an alias for ioresult.IOResult[T].
	IOResult[T any] = ioresult.IOResult[T]

	// Result represents a computation that may fail with an error.
	// It's an alias for result.Result[T].
	Result[T any] = result.Result[T]
)
