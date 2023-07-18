package generic

import (
	ET "github.com/IBM/fp-go/either"
	GIO "github.com/IBM/fp-go/io/generic"
	R "github.com/IBM/fp-go/retry"
)

// Retry combinator for actions that don't raise exceptions, but
// signal in their type the outcome has failed. Examples are the
// `Option`, `Either` and `EitherT` monads.
//
// policy - refers to the retry policy
// action - converts a status into an operation to be executed
// check  - checks if the result of the action needs to be retried
func Retrying[GA ~func() ET.Either[E, A], E, A any](
	policy R.RetryPolicy,
	action func(R.RetryStatus) GA,
	check func(ET.Either[E, A]) bool,
) GA {
	// get an implementation for the types
	return GIO.Retrying(policy, action, check)
}
