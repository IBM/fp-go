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

package readerioresult

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/clock"
	"github.com/IBM/fp-go/v2/clock/testclock"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	R "github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

var clockStart = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// withTestClock returns a context that carries a fresh test clock
func withTestClock(t *testing.T) (*testclock.TestClock, context.Context) {
	tc := testclock.New(clockStart)
	return tc, pair.Second(clock.WithClock(tc)(t.Context()))
}

// runAsync runs a computation in a goroutine and returns a channel with its result
func runAsync[A any](ma ReaderIOResult[A], ctx context.Context) <-chan Result[A] {
	done := make(chan Result[A], 1)
	go func() { done <- ma(ctx)() }()
	return done
}

// waitForCause blocks until the context is done and fails with its cause
func waitForCause(ctx context.Context) IOResult[string] {
	return func() Result[string] {
		<-ctx.Done()
		return result.Left[string](context.Cause(ctx))
	}
}

func TestDelay_TestClock(t *testing.T) {
	t.Run("runs the computation once the clock has advanced", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		done := runAsync(F.Pipe1(Of(42), Delay[int](time.Hour)), ctx)

		tc.BlockUntil(1)
		tc.Advance(time.Hour - time.Nanosecond)
		select {
		case <-done:
			t.Fatal("Delay completed before the clock was advanced far enough")
		default:
		}

		tc.Advance(time.Nanosecond)
		assert.Equal(t, result.Of(42), <-done)
	})

	t.Run("fails with the cause if cancelled while waiting", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		cause := errors.New("stop")
		ctx, cancel := context.WithCancelCause(ctx)
		done := runAsync(F.Pipe1(Of(42), Delay[int](time.Hour)), ctx)

		tc.BlockUntil(1)
		cancel(cause)
		assert.Equal(t, result.Left[int](cause), <-done)
		tc.BlockUntil(0)
	})

	t.Run("does not start a timer if already cancelled", func(t *testing.T) {
		_, ctx := withTestClock(t)
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		assert.Equal(t, result.Left[int](context.Canceled), F.Pipe1(Of(42), Delay[int](time.Hour))(ctx)())
	})
}

func TestTimer_TestClock(t *testing.T) {
	tc, ctx := withTestClock(t)
	done := runAsync(Timer(time.Minute), ctx)

	tc.BlockUntil(1)
	tc.Advance(time.Minute)
	assert.Equal(t, result.Of(clockStart.Add(time.Minute)), <-done)
}

func TestWithTimeout_TestClock(t *testing.T) {
	t.Run("fails with DeadlineExceeded once the timeout elapses", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		done := runAsync(F.Pipe1(waitForCause, WithTimeout[string](time.Minute)), ctx)

		tc.BlockUntil(1)
		tc.Advance(time.Minute)
		assert.Equal(t, result.Left[string](context.DeadlineExceeded), <-done)
	})

	t.Run("interrupts a Delay on the same clock", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		done := runAsync(F.Pipe2(Of(42), Delay[int](time.Hour), WithTimeout[int](time.Minute)), ctx)

		// the timeout and the delay
		tc.BlockUntil(2)
		tc.Advance(time.Minute)
		assert.Equal(t, result.Left[int](context.DeadlineExceeded), <-done)
		tc.BlockUntil(0)
	})

	t.Run("succeeds and stops its timer if the computation is quick", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		assert.Equal(t, result.Of(42), F.Pipe1(Of(42), WithTimeout[int](time.Minute))(ctx)())
		tc.BlockUntil(0)
	})

	t.Run("cause and error of the derived context", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		inspect := func(ctx context.Context) IOResult[[]error] {
			return func() Result[[]error] {
				<-ctx.Done()
				_, hasDeadline := ctx.Deadline()
				assert.False(t, hasDeadline)
				return result.Of([]error{ctx.Err(), context.Cause(ctx)})
			}
		}
		done := runAsync(F.Pipe1(inspect, WithTimeout[[]error](time.Second)), ctx)

		tc.BlockUntil(1)
		tc.Advance(time.Second)
		assert.Equal(t, result.Of([]error{context.Canceled, context.DeadlineExceeded}), <-done)
	})

	t.Run("a non-positive timeout expires immediately", func(t *testing.T) {
		_, ctx := withTestClock(t)
		assert.Equal(t, result.Left[string](context.DeadlineExceeded), F.Pipe1(waitForCause, WithTimeout[string](0))(ctx)())
	})

	t.Run("parent cancellation is propagated", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		cause := errors.New("parent")
		ctx, cancel := context.WithCancelCause(ctx)
		done := runAsync(F.Pipe1(waitForCause, WithTimeout[string](time.Hour)), ctx)

		tc.BlockUntil(1)
		cancel(cause)
		assert.Equal(t, result.Left[string](cause), <-done)
		tc.BlockUntil(0)
	})
}

func TestWithTimeout_RealClock(t *testing.T) {
	t.Run("keeps the deadline of context.WithTimeout", func(t *testing.T) {
		hasDeadline := func(ctx context.Context) IOResult[bool] {
			return func() Result[bool] {
				_, ok := ctx.Deadline()
				return result.Of(ok)
			}
		}
		assert.Equal(t, result.Of(true), F.Pipe1(hasDeadline, WithTimeout[bool](time.Hour))(t.Context())())
	})
}

func TestRetrying_TestClock(t *testing.T) {
	t.Run("exponential backoff runs without real sleeps", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		policy := R.Monoid.Concat(
			R.LimitRetries(5),
			R.ExponentialBackoff(time.Second),
		)
		failing := errors.New("transient")
		action := func(status R.RetryStatus) ReaderIOResult[uint] {
			if status.IterNumber < 4 {
				return Left[uint](failing)
			}
			return Of(status.IterNumber)
		}

		done := runAsync(Retrying(policy, action, result.IsLeft[uint]), ctx)

		// delays of 1s, 2s, 4s and 8s
		for _, d := range []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second} {
			tc.BlockUntil(1)
			tc.Advance(d)
		}

		assert.Equal(t, result.Of[uint](4), <-done)
		assert.Equal(t, clockStart.Add(15*time.Second), tc.Now())
	})

	t.Run("cancellation during a delay stops retrying", func(t *testing.T) {
		tc, ctx := withTestClock(t)
		ctx, cancel := context.WithCancel(ctx)
		attempts := 0
		action := func(R.RetryStatus) ReaderIOResult[int] {
			return func(context.Context) IOResult[int] {
				return func() Result[int] {
					attempts++
					return result.Left[int](errors.New("fail"))
				}
			}
		}
		policy := R.Monoid.Concat(R.LimitRetries(3), R.ConstantDelay(time.Hour))

		done := runAsync(Retrying(policy, action, result.IsLeft[int]), ctx)

		tc.BlockUntil(1)
		cancel()
		assert.Equal(t, result.Left[int](context.Canceled), <-done)
		assert.Equal(t, 1, attempts)
	})
}

func BenchmarkDelay_Execution(b *testing.B) {
	ctx := b.Context()
	rioe := Delay[int](0)(Right(42))
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(ctx)()
	}
}

func BenchmarkTimer_Execution(b *testing.B) {
	ctx := b.Context()
	timer := Timer(0)
	b.ReportAllocs()
	for b.Loop() {
		_ = timer(ctx)()
	}
}

func BenchmarkWithTimeout_Execution(b *testing.B) {
	ctx := b.Context()
	rioe := WithTimeout[int](time.Hour)(Right(42))
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(ctx)()
	}
}
