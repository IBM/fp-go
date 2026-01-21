// Copyright (c) 2023 - 2025 IBM Corp.
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

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

// TestContextCancellationInMap tests that context cancellation is properly handled in Map operations
func TestContextCancellationInMap(t *testing.T) {
	cfg := defaultConfig
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	computation := F.Pipe1(
		Of[AppConfig](42),
		Map[AppConfig](func(n int) int {
			// This should still execute as Map doesn't check context
			return n * 2
		}),
	)

	outcome := computation(cfg)(ctx)()
	// Map operations don't inherently check context, so they succeed
	assert.Equal(t, result.Of(84), outcome)
}

// TestContextCancellationInChain tests context cancellation in Chain operations
func TestContextCancellationInChain(t *testing.T) {
	cfg := defaultConfig
	ctx, cancel := context.WithCancel(context.Background())

	executed := false
	computation := F.Pipe1(
		Of[AppConfig](42),
		Chain(func(n int) ReaderReaderIOResult[AppConfig, int] {
			return func(c AppConfig) ReaderIOResult[context.Context, int] {
				return func(ctx context.Context) IOResult[int] {
					return func() Result[int] {
						// Check if context is cancelled
						select {
						case <-ctx.Done():
							return result.Left[int](ctx.Err())
						default:
							executed = true
							return result.Of(n * 2)
						}
					}
				}
			}
		}),
	)

	cancel() // Cancel before execution
	outcome := computation(cfg)(ctx)()

	assert.True(t, result.IsLeft(outcome))
	assert.False(t, executed, "Chained operation should not execute when context is cancelled")
}

// TestContextCancellationWithTimeout tests timeout-based cancellation
func TestContextCancellationWithTimeout(t *testing.T) {
	cfg := defaultConfig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	computation := func(c AppConfig) ReaderIOResult[context.Context, int] {
		return func(ctx context.Context) IOResult[int] {
			return func() Result[int] {
				// Simulate long-running operation
				select {
				case <-time.After(100 * time.Millisecond):
					return result.Of(42)
				case <-ctx.Done():
					return result.Left[int](ctx.Err())
				}
			}
		}
	}

	outcome := computation(cfg)(ctx)()
	assert.True(t, result.IsLeft(outcome))

	result.Fold(
		func(err error) any {
			assert.ErrorIs(t, err, context.DeadlineExceeded)
			return nil
		},
		func(v int) any {
			t.Fatal("Should have timed out")
			return nil
		},
	)(outcome)
}

// TestContextCancellationInBracket tests that bracket properly handles context cancellation
func TestContextCancellationInBracket(t *testing.T) {
	cfg := defaultConfig
	ctx, cancel := context.WithCancel(context.Background())

	resource := &Resource{id: "res1"}
	useCalled := false

	acquire := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
		return func(ctx context.Context) IOResult[*Resource] {
			return func() Result[*Resource] {
				resource.acquired = true
				return result.Of(resource)
			}
		}
	}

	use := func(r *Resource) ReaderReaderIOResult[AppConfig, string] {
		return func(c AppConfig) ReaderIOResult[context.Context, string] {
			return func(ctx context.Context) IOResult[string] {
				return func() Result[string] {
					select {
					case <-ctx.Done():
						return result.Left[string](ctx.Err())
					default:
						useCalled = true
						return result.Of("success")
					}
				}
			}
		}
	}

	release := func(r *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
		return func(c AppConfig) ReaderIOResult[context.Context, any] {
			return func(ctx context.Context) IOResult[any] {
				return func() Result[any] {
					r.released = true
					return result.Of[any](nil)
				}
			}
		}
	}

	cancel() // Cancel before use
	computation := Bracket(acquire, use, release)
	outcome := computation(cfg)(ctx)()

	assert.True(t, resource.acquired, "Resource should be acquired")
	assert.True(t, resource.released, "Resource should be released even with cancellation")
	assert.False(t, useCalled, "Use should not execute when context is cancelled")
	assert.True(t, result.IsLeft(outcome))
}

// TestContextCancellationInRetry tests context cancellation during retry operations
func TestContextCancellationInRetry(t *testing.T) {
	cfg := defaultConfig
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	attempts := 0
	action := func(status retry.RetryStatus) ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					select {
					case <-ctx.Done():
						return result.Left[int](ctx.Err())
					case <-time.After(30 * time.Millisecond):
						return result.Left[int](errors.New("temporary error"))
					}
				}
			}
		}
	}

	check := func(r Result[int]) bool {
		return result.IsLeft(r)
	}

	policy := retry.LimitRetries(10)
	computation := Retrying(policy, action, check)
	outcome := computation(cfg)(ctx)()

	assert.True(t, result.IsLeft(outcome))
	// Should stop retrying when context is cancelled
	assert.Less(t, attempts, 10, "Should stop retrying when context is cancelled")
}

// TestContextPropagationThroughMonadTransforms tests that context is properly propagated
func TestContextPropagationThroughMonadTransforms(t *testing.T) {
	cfg := defaultConfig

	t.Run("context propagates through Map", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", "value")

		var capturedCtx context.Context
		computation := func(c AppConfig) ReaderIOResult[context.Context, string] {
			return func(ctx context.Context) IOResult[string] {
				return func() Result[string] {
					capturedCtx = ctx
					return result.Of("test")
				}
			}
		}

		_ = computation(cfg)(ctx)()
		assert.Equal(t, "value", capturedCtx.Value("key"))
	})

	t.Run("context propagates through Chain", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", "value")

		var capturedCtx context.Context
		computation := F.Pipe1(
			Of[AppConfig](42),
			Chain(func(n int) ReaderReaderIOResult[AppConfig, int] {
				return func(c AppConfig) ReaderIOResult[context.Context, int] {
					return func(ctx context.Context) IOResult[int] {
						return func() Result[int] {
							capturedCtx = ctx
							return result.Of(n * 2)
						}
					}
				}
			}),
		)

		_ = computation(cfg)(ctx)()
		assert.Equal(t, "value", capturedCtx.Value("key"))
	})

	t.Run("context propagates through Ap", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", "value")

		var capturedCtx context.Context
		fab := func(c AppConfig) ReaderIOResult[context.Context, func(int) int] {
			return func(ctx context.Context) IOResult[func(int) int] {
				return func() Result[func(int) int] {
					capturedCtx = ctx
					return result.Of(N.Mul(2))
				}
			}
		}

		fa := Of[AppConfig](21)
		computation := MonadAp(fab, fa)

		_ = computation(cfg)(ctx)()
		assert.Equal(t, "value", capturedCtx.Value("key"))
	})
}

// TestContextCancellationInAlt tests Alt operation with context cancellation
func TestContextCancellationInAlt(t *testing.T) {
	cfg := defaultConfig
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	firstCalled := false
	secondCalled := false

	first := func(c AppConfig) ReaderIOResult[context.Context, int] {
		return func(ctx context.Context) IOResult[int] {
			return func() Result[int] {
				select {
				case <-ctx.Done():
					return result.Left[int](ctx.Err())
				default:
					firstCalled = true
					return result.Left[int](errors.New("first error"))
				}
			}
		}
	}

	second := func() ReaderReaderIOResult[AppConfig, int] {
		return func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					select {
					case <-ctx.Done():
						return result.Left[int](ctx.Err())
					default:
						secondCalled = true
						return result.Of(42)
					}
				}
			}
		}
	}

	computation := MonadAlt(first, second)
	outcome := computation(cfg)(ctx)()

	assert.True(t, result.IsLeft(outcome))
	assert.False(t, firstCalled, "First should not execute when context is cancelled")
	assert.False(t, secondCalled, "Second should not execute when context is cancelled")
}

// TestContextCancellationInDoNotation tests context cancellation in do-notation
func TestContextCancellationInDoNotation(t *testing.T) {
	cfg := defaultConfig
	ctx, cancel := context.WithCancel(context.Background())

	type State struct {
		Value1 int
		Value2 int
	}

	step1Executed := false
	step2Executed := false

	computation := F.Pipe2(
		Do[AppConfig](State{}),
		Bind(
			func(v int) func(State) State {
				return func(s State) State {
					s.Value1 = v
					return s
				}
			},
			func(s State) ReaderReaderIOResult[AppConfig, int] {
				return func(c AppConfig) ReaderIOResult[context.Context, int] {
					return func(ctx context.Context) IOResult[int] {
						return func() Result[int] {
							select {
							case <-ctx.Done():
								return result.Left[int](ctx.Err())
							default:
								step1Executed = true
								return result.Of(10)
							}
						}
					}
				}
			},
		),
		Bind(
			func(v int) func(State) State {
				return func(s State) State {
					s.Value2 = v
					return s
				}
			},
			func(s State) ReaderReaderIOResult[AppConfig, int] {
				return func(c AppConfig) ReaderIOResult[context.Context, int] {
					return func(ctx context.Context) IOResult[int] {
						return func() Result[int] {
							select {
							case <-ctx.Done():
								return result.Left[int](ctx.Err())
							default:
								step2Executed = true
								return result.Of(20)
							}
						}
					}
				}
			},
		),
	)

	cancel() // Cancel before execution
	outcome := computation(cfg)(ctx)()

	assert.True(t, result.IsLeft(outcome))
	assert.False(t, step1Executed, "Step 1 should not execute when context is cancelled")
	assert.False(t, step2Executed, "Step 2 should not execute when context is cancelled")
}

// TestContextCancellationBetweenSteps tests cancellation between sequential steps
func TestContextCancellationBetweenSteps(t *testing.T) {
	cfg := defaultConfig
	ctx, cancel := context.WithCancel(context.Background())

	step1Executed := false
	step2Executed := false

	computation := F.Pipe1(
		func(c AppConfig) ReaderIOResult[context.Context, int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					step1Executed = true
					cancel() // Cancel after first step
					return result.Of(42)
				}
			}
		},
		Chain(func(n int) ReaderReaderIOResult[AppConfig, int] {
			return func(c AppConfig) ReaderIOResult[context.Context, int] {
				return func(ctx context.Context) IOResult[int] {
					return func() Result[int] {
						select {
						case <-ctx.Done():
							return result.Left[int](ctx.Err())
						default:
							step2Executed = true
							return result.Of(n * 2)
						}
					}
				}
			}
		}),
	)

	outcome := computation(cfg)(ctx)()

	assert.True(t, step1Executed, "Step 1 should execute")
	assert.False(t, step2Executed, "Step 2 should not execute after cancellation")
	assert.True(t, result.IsLeft(outcome))
}
