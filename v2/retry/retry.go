// Copyright (c) 2023 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import (
	"math"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
)

type RetryStatus struct {
	// Iteration number, where `0` is the first try
	IterNumber uint
	// Delay incurred so far from retries
	CumulativeDelay time.Duration
	// Latest attempt's delay. Will always be `none` on first run.
	PreviousDelay Option[time.Duration]
}

// RetryPolicy is a function that takes an `RetryStatus` and
// possibly returns a delay in milliseconds. Iteration numbers start
// at zero and increase by one on each retry. A //None// return value from
// the function implies we have reached the retry limit.
type RetryPolicy = func(RetryStatus) Option[time.Duration]

const emptyDuration = time.Duration(0)

var ordDuration = ord.FromStrictCompare[time.Duration]()

// Monoid 'RetryPolicy' is a 'Monoid'. You can collapse multiple strategies into one using 'concat'.
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
	return func(status RetryStatus) Option[time.Duration] {
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
	return func(status RetryStatus) Option[time.Duration] {
		return O.Some(delay * time.Duration(math.Pow(2, float64(status.IterNumber))))
	}
}

// DefaultRetryStatus is the default retry status. Exported mostly to allow user code
// to test their handlers and retry policies.
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
