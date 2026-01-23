package codec

import (
	"fmt"

	"github.com/IBM/fp-go/v2/errors"
	"github.com/IBM/fp-go/v2/internal/formatting"
	"github.com/IBM/fp-go/v2/result"
)

func onTypeError(expType string) func(any) error {
	return errors.OnSome[any](fmt.Sprintf("expecting type [%s] but got [%%T]", expType))
}

// Is checks if a value can be converted to type T.
// Returns Some(value) if the conversion succeeds, None otherwise.
// This is a type-safe cast operation.
func Is[T any]() ReaderResult[any, T] {
	return result.ToType[T](onTypeError(formatting.TypeInfo(*new(T))))
}
