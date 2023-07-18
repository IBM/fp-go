package generic

import (
	R "github.com/IBM/fp-go/retry"
	G "github.com/IBM/fp-go/retry/generic"
)

type retryStatusIO = func() R.RetryStatus

// Retry combinator for actions that don't raise exceptions, but
// signal in their type the outcome has failed. Examples are the
// `Option`, `Either` and `EitherT` monads.
//
// policy - refers to the retry policy
// action - converts a status into an operation to be executed
// check  - checks if the result of the action needs to be retried
func Retrying[GA ~func() A, A any](
	policy R.RetryPolicy,
	action func(R.RetryStatus) GA,
	check func(A) bool,
) GA {
	// get an implementation for the types
	return G.Retrying(
		Chain[GA, GA, A, A],
		Chain[retryStatusIO, GA, R.RetryStatus, A],
		Of[GA, A],
		Of[retryStatusIO, R.RetryStatus],
		Delay[retryStatusIO, R.RetryStatus],

		policy,
		action,
		check,
	)
}
