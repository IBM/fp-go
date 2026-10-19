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

package clock_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/clock"
	"github.com/IBM/fp-go/v2/clock/testclock"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

var start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func withClock(ctx context.Context, c clock.Clock) context.Context {
	return pair.Second(clock.WithClock(c)(ctx))
}

func TestFromContext_Success(t *testing.T) {
	t.Run("returns the clock stored in the context", func(t *testing.T) {
		tc := testclock.New(start)
		ctx := withClock(t.Context(), tc)
		assert.Same(t, tc, clock.FromContext(ctx))
		assert.Equal(t, start, clock.Now(ctx))
	})

	t.Run("the cancel function of WithClock is a no-op", func(t *testing.T) {
		cancel, ctx := pair.Unpack(clock.WithClock(testclock.New(start))(t.Context()))
		cancel()
		assert.NoError(t, ctx.Err())
	})
}

func TestFromContext_EdgeCases(t *testing.T) {
	t.Run("falls back to the real clock", func(t *testing.T) {
		c := clock.FromContext(t.Context())
		assert.True(t, clock.IsReal(c))
	})

	t.Run("a test clock is not real", func(t *testing.T) {
		assert.False(t, clock.IsReal(testclock.New(start)))
	})

	t.Run("an inner clock shadows an outer one", func(t *testing.T) {
		outer := testclock.New(start)
		inner := testclock.New(start.Add(time.Hour))
		ctx := withClock(withClock(t.Context(), outer), inner)
		assert.Equal(t, start.Add(time.Hour), clock.Now(ctx))
	})
}

func TestReal_Success(t *testing.T) {
	t.Run("Now is close to time.Now", func(t *testing.T) {
		before := time.Now()
		now := clock.Real().Now()
		assert.False(t, now.Before(before))
	})

	t.Run("After fires", func(t *testing.T) {
		ch, stop := clock.Real().After(0)
		defer stop()
		<-ch
	})

	t.Run("stop reports a pending timer", func(t *testing.T) {
		_, stop := clock.Real().After(time.Hour)
		assert.True(t, stop())
		assert.False(t, stop())
	})
}

func TestSleep_Success(t *testing.T) {
	t.Run("returns after the clock has advanced", func(t *testing.T) {
		tc := testclock.New(start)
		ctx := withClock(t.Context(), tc)

		done := make(chan error, 1)
		go func() { done <- clock.Sleep(ctx, time.Minute) }()

		tc.BlockUntil(1)
		tc.Advance(time.Minute)
		assert.NoError(t, <-done)
	})

	t.Run("zero duration returns immediately", func(t *testing.T) {
		ctx := withClock(t.Context(), testclock.New(start))
		assert.NoError(t, clock.Sleep(ctx, 0))
	})
}

func TestSleep_Failure(t *testing.T) {
	t.Run("already cancelled context", func(t *testing.T) {
		cause := errors.New("stop")
		ctx, cancel := context.WithCancelCause(t.Context())
		cancel(cause)
		assert.ErrorIs(t, clock.Sleep(ctx, time.Hour), cause)
	})

	t.Run("cancelled while sleeping, the timer is stopped", func(t *testing.T) {
		tc := testclock.New(start)
		cause := errors.New("stop")
		ctx, cancel := context.WithCancelCause(withClock(t.Context(), tc))

		done := make(chan error, 1)
		go func() { done <- clock.Sleep(ctx, time.Hour) }()

		tc.BlockUntil(1)
		cancel(cause)
		assert.ErrorIs(t, <-done, cause)
		// the stopped timer no longer counts as pending
		tc.BlockUntil(0)
	})
}

func TestNowIO_Success(t *testing.T) {
	tc := testclock.New(start)
	now := clock.NowIO(tc)
	assert.Equal(t, start, now())
	tc.Advance(time.Second)
	assert.Equal(t, start.Add(time.Second), now())
}
