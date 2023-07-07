package errors

import (
	"errors"

	O "github.com/ibm/fp-go/option"
)

// As tries to extract the error of desired type from the given error
func As[A error]() func(error) O.Option[A] {
	return O.FromValidation(func(err error) (A, bool) {
		var a A
		ok := errors.As(err, &a)
		return a, ok
	})
}
