package Retry

import (
	"math"
	"time"

	F "github.com/ibm/fp-go/function"
	M "github.com/ibm/fp-go/monoid"
	O "github.com/ibm/fp-go/option"
	"github.com/ibm/fp-go/ord"
)

type RetryStatus struct {
	// Iteration number, where `0` is the first try
	IterNumber uint
	// Delay incurred so far from retries
	CumulativeDelay time.Duration
	// Latest attempt's delay. Will always be `none` on first run.
	PreviousDelay O.Option[time.Duration]
}

// RetryPolicy is a function that takes an `RetryStatus` and
// possibly returns a delay in milliseconds. Iteration numbers start
// at zero and increase by one on each retry. A //None// return value from
// the function implies we have reached the retry limit.
type RetryPolicy = func(RetryStatus) O.Option[time.Duration]

const emptyDuration = time.Duration(0)

var ordDuration = ord.FromStrictCompare[time.Duration]()

// 'RetryPolicy' is a 'Monoid'. You can collapse multiple strategies into one using 'concat'.
// The semantics of this combination are as follows:
//
// 1. If either policy returns 'None', the combined policy returns
// 'None'. This can be used to inhibit after a number of retries,
// for example.
//
// 2. If both policies return a delay, the larger delay will be used.
// This is quite natural when combining multiple policies to achieve a
// certain effect.
var Monoid = M.FunctionMonoid[RetryStatus](O.ApplicativeMonoid(M.MakeMonoid(
	ord.MaxSemigroup(ordDuration).Concat, emptyDuration)))

// LimitRetries retries immediately, but only up to `i` times.
func LimitRetries(i uint) RetryPolicy {
	pred := func(value uint) bool {
		return value < i
	}
	empty := F.Constant1[uint](emptyDuration)
	return func(status RetryStatus) O.Option[time.Duration] {
		return F.Pipe2(
			status.IterNumber,
			O.FromPredicate(pred),
			O.Map(empty),
		)
	}
}

// ConstantDelay delays with unlimited retries
func ConstantDelay(delay time.Duration) RetryPolicy {
	return F.Constant1[RetryStatus](O.Of(delay))
}

// CapDelay sets a time-upperbound for any delays that may be directed by the
// given policy. This function does not terminate the retrying. The policy
// capDelay(maxDelay, exponentialBackoff(n))` will never stop retrying. It
// will reach a state where it retries forever with a delay of `maxDelay`
// between each one. To get termination you need to use one of the
// 'limitRetries' function variants.
func CapDelay(maxDelay time.Duration, policy RetryPolicy) RetryPolicy {
	return F.Flow2(
		policy,
		O.Map(F.Bind1st(ord.Min(ordDuration), maxDelay)),
	)
}

// ExponentialBackoff grows delay exponentially each iteration.
// Each delay will increase by a factor of two.
func ExponentialBackoff(delay time.Duration) RetryPolicy {
	return func(status RetryStatus) O.Option[time.Duration] {
		return O.Some(delay * time.Duration(math.Pow(2, float64(status.IterNumber))))
	}
}

/**
 * Initial, default retry status. Exported mostly to allow user code
 * to test their handlers and retry policies.
 */
var DefaultRetryStatus = RetryStatus{
	IterNumber:      0,
	CumulativeDelay: 0,
	PreviousDelay:   O.None[time.Duration](),
}

var getOrElseDelay = O.GetOrElse(F.Constant(emptyDuration))

/**
 * Apply policy on status to see what the decision would be.
 */
func ApplyPolicy(policy RetryPolicy, status RetryStatus) RetryStatus {
	previousDelay := policy(status)
	return RetryStatus{
		IterNumber:      status.IterNumber + 1,
		CumulativeDelay: status.CumulativeDelay + getOrElseDelay(previousDelay),
		PreviousDelay:   previousDelay,
	}
}
