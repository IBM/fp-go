package io

import (
	G "github.com/IBM/fp-go/io/generic"
	R "github.com/IBM/fp-go/retry"
)

// Retrying will retry the actions according to the check policy
//
// policy - refers to the retry policy
// action - converts a status into an operation to be executed
// check  - checks if the result of the action needs to be retried
func Retrying[A any](
	policy R.RetryPolicy,
	action func(R.RetryStatus) IO[A],
	check func(A) bool,
) IO[A] {
	return G.Retrying(policy, action, check)
}
