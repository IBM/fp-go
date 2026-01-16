package codec

import (
	"fmt"

	"github.com/IBM/fp-go/v2/result"
)

func onTypeError(expType string) func(any) error {
	return func(u any) error {
		return fmt.Errorf("expecting type [%s] but got [%T]", expType, u)
	}
}

// Is checks if a value can be converted to type T.
// Returns Some(value) if the conversion succeeds, None otherwise.
// This is a type-safe cast operation.
func Is[T any]() func(any) Result[T] {
	var zero T
	return result.ToType[T](onTypeError(fmt.Sprintf("%T", zero)))
}
