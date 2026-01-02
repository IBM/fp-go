package builder

import (
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Result represents a computation that may fail with an error.
	// It's an alias for Either[error, T].
	Result[T any] = result.Result[T]

	// Prism is an optic that focuses on a case of a sum type.
	// It provides a way to extract and construct values of a specific variant.
	Prism[S, A any] = prism.Prism[S, A]

	// Option represents an optional value that may or may not be present.
	Option[T any] = option.Option[T]
)
