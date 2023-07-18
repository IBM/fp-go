package ioeither

import (
	ET "github.com/ibm/fp-go/either"
	G "github.com/ibm/fp-go/ioeither/generic"
	R "github.com/ibm/fp-go/retry"
)

// Retrying will retry the actions according to the check policy
//
// policy - refers to the retry policy
// action - converts a status into an operation to be executed
// check  - checks if the result of the action needs to be retried
func Retrying[E, A any](
	policy R.RetryPolicy,
	action func(R.RetryStatus) IOEither[E, A],
	check func(ET.Either[E, A]) bool,
) IOEither[E, A] {
	return G.Retrying(policy, action, check)
}
