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
	"errors"
	"testing"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

type RetryOuterCtx struct {
	maxRetries int
}

type RetryInnerCtx struct {
	endpoint string
}

func TestRetryingSuccess(t *testing.T) {
	outer := RetryOuterCtx{maxRetries: 3}
	inner := RetryInnerCtx{endpoint: "http://example.com"}

	// Action that succeeds immediately
	action := func(status retry.RetryStatus) ReaderReaderIOEither[RetryOuterCtx, RetryInnerCtx, error, string] {
		return func(o RetryOuterCtx) ReaderIOEither[RetryInnerCtx, error, string] {
			return func(i RetryInnerCtx) IOE.IOEither[error, string] {
				return IOE.Of[error]("success")
			}
		}
	}

	// Never retry on success
	check := func(result E.Either[error, string]) bool {
		return E.IsLeft(result)
	}

	policy := retry.LimitRetries(3)
	result := Retrying(policy, action, check)

	assert.Equal(t, E.Right[error]("success"), result(outer)(inner)())
}

func TestRetryingEventualSuccess(t *testing.T) {
	outer := RetryOuterCtx{maxRetries: 3}
	inner := RetryInnerCtx{endpoint: "http://example.com"}

	attempts := 0
	err := errors.New("temporary error")

	// Action that fails twice then succeeds
	action := func(status retry.RetryStatus) ReaderReaderIOEither[RetryOuterCtx, RetryInnerCtx, error, string] {
		return func(o RetryOuterCtx) ReaderIOEither[RetryInnerCtx, error, string] {
			return func(i RetryInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					attempts++
					if attempts < 3 {
						return E.Left[string](err)
					}
					return E.Right[error]("success after retries")
				}
			}
		}
	}

	// Retry on error
	check := func(result E.Either[error, string]) bool {
		return E.IsLeft(result)
	}

	policy := retry.LimitRetries(5)
	result := Retrying(policy, action, check)

	outcome := result(outer)(inner)()
	assert.Equal(t, E.Right[error]("success after retries"), outcome)
	assert.Equal(t, 3, attempts)
}

func TestRetryingMaxRetriesExceeded(t *testing.T) {
	outer := RetryOuterCtx{maxRetries: 2}
	inner := RetryInnerCtx{endpoint: "http://example.com"}

	err := errors.New("persistent error")
	attempts := 0

	// Action that always fails
	action := func(status retry.RetryStatus) ReaderReaderIOEither[RetryOuterCtx, RetryInnerCtx, error, string] {
		return func(o RetryOuterCtx) ReaderIOEither[RetryInnerCtx, error, string] {
			return func(i RetryInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					attempts++
					return E.Left[string](err)
				}
			}
		}
	}

	// Always retry on error
	check := func(result E.Either[error, string]) bool {
		return E.IsLeft(result)
	}

	policy := retry.LimitRetries(2)
	result := Retrying(policy, action, check)

	outcome := result(outer)(inner)()
	assert.Equal(t, E.Left[string](err), outcome)
	assert.Equal(t, 3, attempts) // Initial attempt + 2 retries
}

func TestRetryingWithRetryStatus(t *testing.T) {
	outer := RetryOuterCtx{maxRetries: 3}
	inner := RetryInnerCtx{endpoint: "http://example.com"}

	var statuses []int

	// Action that tracks retry attempts
	action := func(status retry.RetryStatus) ReaderReaderIOEither[RetryOuterCtx, RetryInnerCtx, error, string] {
		return func(o RetryOuterCtx) ReaderIOEither[RetryInnerCtx, error, string] {
			return func(i RetryInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					statuses = append(statuses, int(status.IterNumber))
					if status.IterNumber < 2 {
						return E.Left[string](errors.New("retry"))
					}
					return E.Right[error]("success")
				}
			}
		}
	}

	check := func(result E.Either[error, string]) bool {
		return E.IsLeft(result)
	}

	policy := retry.LimitRetries(5)
	result := Retrying(policy, action, check)

	outcome := result(outer)(inner)()
	assert.Equal(t, E.Right[error]("success"), outcome)
	assert.Equal(t, []int{0, 1, 2}, statuses)
}

func TestRetryingWithContextAccess(t *testing.T) {
	outer := RetryOuterCtx{maxRetries: 2}
	inner := RetryInnerCtx{endpoint: "http://api.example.com"}

	attempts := 0

	// Action that uses both contexts
	action := func(status retry.RetryStatus) ReaderReaderIOEither[RetryOuterCtx, RetryInnerCtx, error, string] {
		return func(o RetryOuterCtx) ReaderIOEither[RetryInnerCtx, error, string] {
			return func(i RetryInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					attempts++
					if attempts <= o.maxRetries {
						return E.Left[string](errors.New("retry"))
					}
					return E.Right[error](i.endpoint + " success")
				}
			}
		}
	}

	check := func(result E.Either[error, string]) bool {
		return E.IsLeft(result)
	}

	policy := retry.LimitRetries(5)
	result := Retrying(policy, action, check)

	outcome := result(outer)(inner)()
	assert.Equal(t, E.Right[error]("http://api.example.com success"), outcome)
	assert.Equal(t, 3, attempts) // Fails twice (attempts 1,2), succeeds on attempt 3
}

func TestRetryingWithExponentialBackoff(t *testing.T) {
	outer := RetryOuterCtx{maxRetries: 3}
	inner := RetryInnerCtx{endpoint: "http://example.com"}

	attempts := 0
	startTime := time.Now()

	// Action that fails twice
	action := func(status retry.RetryStatus) ReaderReaderIOEither[RetryOuterCtx, RetryInnerCtx, error, string] {
		return func(o RetryOuterCtx) ReaderIOEither[RetryInnerCtx, error, string] {
			return func(i RetryInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					attempts++
					if attempts < 3 {
						return E.Left[string](errors.New("retry"))
					}
					return E.Right[error]("success")
				}
			}
		}
	}

	check := func(result E.Either[error, string]) bool {
		return E.IsLeft(result)
	}

	// Exponential backoff: 10ms, 20ms, 40ms...
	policy := retry.ExponentialBackoff(10 * time.Millisecond)
	result := Retrying(policy, action, check)

	outcome := result(outer)(inner)()
	elapsed := time.Since(startTime)

	assert.Equal(t, E.Right[error]("success"), outcome)
	assert.Equal(t, 3, attempts)
	// Should have at least 30ms delay (10ms + 20ms)
	assert.True(t, elapsed >= 30*time.Millisecond, "Expected at least 30ms delay, got %v", elapsed)
}

func TestRetryingNoRetryOnSuccess(t *testing.T) {
	outer := RetryOuterCtx{maxRetries: 3}
	inner := RetryInnerCtx{endpoint: "http://example.com"}

	attempts := 0

	action := func(status retry.RetryStatus) ReaderReaderIOEither[RetryOuterCtx, RetryInnerCtx, error, string] {
		return func(o RetryOuterCtx) ReaderIOEither[RetryInnerCtx, error, string] {
			return func(i RetryInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					attempts++
					return E.Right[error]("immediate success")
				}
			}
		}
	}

	// Only retry on errors
	check := func(result E.Either[error, string]) bool {
		return E.IsLeft(result)
	}

	policy := retry.LimitRetries(5)
	result := Retrying(policy, action, check)

	outcome := result(outer)(inner)()
	assert.Equal(t, E.Right[error]("immediate success"), outcome)
	assert.Equal(t, 1, attempts) // Should only run once
}

func TestRetryingCustomCheckPredicate(t *testing.T) {
	outer := RetryOuterCtx{maxRetries: 3}
	inner := RetryInnerCtx{endpoint: "http://example.com"}

	attempts := 0
	retryableErr := errors.New("retryable")
	fatalErr := errors.New("fatal")

	action := func(status retry.RetryStatus) ReaderReaderIOEither[RetryOuterCtx, RetryInnerCtx, error, string] {
		return func(o RetryOuterCtx) ReaderIOEither[RetryInnerCtx, error, string] {
			return func(i RetryInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					attempts++
					if attempts == 1 {
						return E.Left[string](retryableErr)
					}
					if attempts == 2 {
						return E.Left[string](fatalErr)
					}
					return E.Right[error]("success")
				}
			}
		}
	}

	// Only retry on retryable errors
	check := func(result E.Either[error, string]) bool {
		return E.Fold(
			func(err error) bool {
				return err.Error() == "retryable"
			},
			func(string) bool {
				return false
			},
		)(result)
	}

	policy := retry.LimitRetries(5)
	result := Retrying(policy, action, check)

	outcome := result(outer)(inner)()
	// Should stop on fatal error without further retries
	assert.Equal(t, E.Left[string](fatalErr), outcome)
	assert.Equal(t, 2, attempts)
}
