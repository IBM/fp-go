// Copyright (c) 2025 IBM Corp.
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
package readerreaderioresult

import (
	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/retry"
)

// Retrying executes an action with automatic retry logic based on a retry policy.
// It retries the action when it fails or when the check predicate returns false.
//
// This function is useful for handling transient failures in operations like:
//   - Network requests that may temporarily fail
//   - Database operations that may encounter locks
//   - External service calls that may be temporarily unavailable
//
// Parameters:
//   - policy: Defines the retry behavior (number of retries, delays, backoff strategy)
//   - action: The computation to retry, receives retry status information
//   - check: Predicate to determine if the result should trigger a retry (returns true to continue, false to retry)
//
// The action receives a retry.RetryStatus that contains:
//   - IterNumber: Current iteration number (0-based)
//   - CumulativeDelay: Total delay accumulated so far
//   - PreviousDelay: Delay from the previous iteration
//
// Returns:
//   - A ReaderReaderIOResult that executes the action with retry logic
//
// Example:
//
//	import (
//	    "errors"
//	    "time"
//	    "github.com/IBM/fp-go/v2/retry"
//	)
//
//	type Config struct {
//	    MaxRetries int
//	    BaseDelay  time.Duration
//	}
//
//	// Create a retry policy with exponential backoff
//	policy := retry.ExponentialBackoff(100*time.Millisecond, 5*time.Second)
//	policy = retry.LimitRetries(3, policy)
//
//	// Action that may fail transiently
//	action := func(status retry.RetryStatus) ReaderReaderIOResult[Config, string] {
//	    return func(cfg Config) ReaderIOResult[context.Context, string] {
//	        return func(ctx context.Context) IOResult[string] {
//	            return func() Either[error, string] {
//	                // Simulate transient failure
//	                if status.IterNumber < 2 {
//	                    return either.Left[string](errors.New("transient error"))
//	                }
//	                return either.Right[error]("success")
//	            }
//	        }
//	    }
//	}
//
//	// Check if we should retry (retry on any error)
//	check := func(result Result[string]) bool {
//	    return either.IsRight(result) // Continue only if successful
//	}
//
//	// Execute with retry logic
//	result := Retrying(policy, action, check)
//
//go:inline
func Retrying[R, A any](
	policy retry.RetryPolicy,
	action Kleisli[R, retry.RetryStatus, A],
	check Predicate[Result[A]],
) ReaderReaderIOResult[R, A] {
	// get an implementation for the types
	return F.Flow4(
		reader.Read[RIOE.ReaderIOResult[A]],
		reader.Map[retry.RetryStatus],
		reader.Read[RIOE.Kleisli[retry.RetryStatus, A]](action),
		F.Bind13of3(RIOE.Retrying[A])(policy, check),
	)
}
