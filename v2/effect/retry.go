package effect

import (
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/retry"
)

func Retrying[C, A any](
	policy retry.RetryPolicy,
	action Kleisli[C, retry.RetryStatus, A],
	check Predicate[Result[A]],
) Effect[C, A] {
	return readerreaderioresult.Retrying(policy, action, check)
}
