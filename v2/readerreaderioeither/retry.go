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
package readerreaderioeither

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/reader"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	"github.com/IBM/fp-go/v2/retry"
)

// Retrying retries an action according to a retry policy until it succeeds or the policy gives up.
// The action receives a RetryStatus that tracks the retry attempt number and cumulative delay.
// The check predicate determines whether a result should trigger a retry.
//
// This is useful for operations that may fail transiently and need to be retried with backoff,
// while having access to both outer (R) and inner (C) reader contexts.
//
// Parameters:
//   - policy: The retry policy that determines delays and when to give up
//   - action: A Kleisli function that takes RetryStatus and returns a ReaderReaderIOEither
//   - check: A predicate that returns true if the result should trigger a retry
//
// Returns:
//   - A ReaderReaderIOEither that will retry according to the policy
//
// Example:
//
//	type OuterConfig struct {
//	    MaxRetries int
//	}
//	type InnerConfig struct {
//	    Endpoint string
//	}
//
//	// Retry a network call with exponential backoff
//	policy := retry.ExponentialBackoff(100*time.Millisecond, 2.0)
//
//	action := func(status retry.RetryStatus) readerreaderioeither.ReaderReaderIOEither[OuterConfig, InnerConfig, error, Response] {
//	    return func(outer OuterConfig) readerioeither.ReaderIOEither[InnerConfig, error, Response] {
//	        return func(inner InnerConfig) ioeither.IOEither[error, Response] {
//	            return ioeither.TryCatch(
//	                func() (Response, error) {
//	                    return makeRequest(inner.Endpoint)
//	                },
//	                func(err error) error { return err },
//	            )
//	        }
//	    }
//	}
//
//	// Retry on network errors
//	check := func(result either.Either[error, Response]) bool {
//	    return either.IsLeft(result) && isNetworkError(either.GetLeft(result))
//	}
//
//	result := readerreaderioeither.Retrying(policy, action, check)
//
//go:inline
func Retrying[R, C, E, A any](
	policy retry.RetryPolicy,
	action Kleisli[R, C, E, retry.RetryStatus, A],
	check Predicate[Either[E, A]],
) ReaderReaderIOEither[R, C, E, A] {
	// get an implementation for the types
	return func(r R) ReaderIOEither[C, E, A] {
		return RIOE.Retrying(policy, F.Pipe1(action, reader.Map[retry.RetryStatus](reader.Read[ReaderIOEither[C, E, A]](r))), check)
	}
}
