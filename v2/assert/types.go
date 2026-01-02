package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Result represents a computation that may fail with an error.
	Result[T any] = result.Result[T]

	// Reader represents a test assertion that depends on a testing.T context and returns a boolean.
	Reader = reader.Reader[*testing.T, bool]

	// Kleisli represents a function that produces a test assertion Reader from a value of type T.
	Kleisli[T any] = reader.Reader[T, Reader]

	// Predicate represents a function that tests a value of type T and returns a boolean.
	Predicate[T any] = predicate.Predicate[T]

	// Lens is a functional reference to a subpart of a data structure.
	Lens[S, T any] = lens.Lens[S, T]

	// Optional is an optic that focuses on a value that may or may not be present.
	Optional[S, T any] = optional.Optional[S, T]

	// Prism is an optic that focuses on a case of a sum type.
	Prism[S, T any] = prism.Prism[S, T]
)
