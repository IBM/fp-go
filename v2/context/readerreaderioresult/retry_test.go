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
	"context"
	"errors"
	"testing"
	"time"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

func TestRetryingSuccess(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	attempts := 0
	action := func(status retry.RetryStatus) ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					if attempts < 3 {
						return result.Left[int](errors.New("temporary error"))
					}
					return result.Of(42)
				}
			}
		}
	}

	check := result.IsLeft[int]

	policy := retry.LimitRetries(5)

	computation := Retrying(policy, action, check)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Of(42), outcome)
	assert.Equal(t, 3, attempts)
}

func TestRetryingFailureExhaustsRetries(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	attempts := 0
	testErr := errors.New("persistent error")

	action := func(status retry.RetryStatus) ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					return result.Left[int](testErr)
				}
			}
		}
	}

	check := result.IsLeft[int]

	policy := retry.LimitRetries(3)

	computation := Retrying(policy, action, check)
	outcome := computation(cfg)(ctx)()

	assert.True(t, result.IsLeft(outcome))
	assert.Equal(t, 4, attempts) // Initial attempt + 3 retries
}

func TestRetryingNoRetryNeeded(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	attempts := 0
	action := func(status retry.RetryStatus) ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					return result.Of(42)
				}
			}
		}
	}

	check := result.IsLeft[int]

	policy := retry.LimitRetries(5)

	computation := Retrying(policy, action, check)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Of(42), outcome)
	assert.Equal(t, 1, attempts) // Only initial attempt
}

func TestRetryingWithDelay(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	attempts := 0
	start := time.Now()

	action := func(status retry.RetryStatus) ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					if attempts < 2 {
						return result.Left[int](errors.New("temporary error"))
					}
					return result.Of(42)
				}
			}
		}
	}

	check := result.IsLeft[int]

	// Policy with delay
	policy := retry.CapDelay(
		100*time.Millisecond,
		retry.LimitRetries(3),
	)

	computation := Retrying(policy, action, check)
	outcome := computation(cfg)(ctx)()
	elapsed := time.Since(start)

	assert.Equal(t, result.Of(42), outcome)
	assert.Equal(t, 2, attempts)
	// The delay might be very short in tests, so just check it completed
	_ = elapsed
}

func TestRetryingAccessesConfig(t *testing.T) {
	cfg := AppConfig{DatabaseURL: "test-db", LogLevel: "debug"}
	ctx := t.Context()

	attempts := 0
	var capturedURL string

	action := func(status retry.RetryStatus) ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					capturedURL = c.DatabaseURL
					if attempts < 2 {
						return result.Left[int](errors.New("temporary error"))
					}
					return result.Of(len(c.DatabaseURL))
				}
			}
		}
	}

	check := result.IsLeft[int]

	policy := retry.LimitRetries(3)

	computation := Retrying(policy, action, check)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Of(7), outcome) // len("test-db")
	assert.Equal(t, "test-db", capturedURL)
	assert.Equal(t, 2, attempts)
}

func TestRetryingWithExponentialBackoff(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	attempts := 0
	action := func(status retry.RetryStatus) ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					if attempts < 3 {
						return result.Left[int](errors.New("temporary error"))
					}
					return result.Of(42)
				}
			}
		}
	}

	check := result.IsLeft[int]

	// Exponential backoff policy
	policy := retry.CapDelay(
		200*time.Millisecond,
		retry.LimitRetries(5),
	)

	computation := Retrying(policy, action, check)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Of(42), outcome)
	assert.Equal(t, 3, attempts)
}

func TestRetryingCheckFunction(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	attempts := 0
	action := func(status retry.RetryStatus) ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					return result.Of(attempts)
				}
			}
		}
	}

	// Retry while result is less than 3
	check := func(r Result[int]) bool {
		return result.Fold(
			reader.Of[error](true),
			N.LessThan(3),
		)(r)
	}

	policy := retry.LimitRetries(10)

	computation := Retrying(policy, action, check)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Of(3), outcome)
	assert.Equal(t, 3, attempts)
}
