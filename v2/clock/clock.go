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

package clock

import (
	"context"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/common"
	IC "github.com/IBM/fp-go/v2/internal/context"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
)

type (
	// Clock is a source of time.
	//
	// Code that reads the time or waits for a duration through a Clock can be tested
	// with a test clock, which advances only when the test says so, instead of sleeping.
	Clock interface {
		// Now returns the current time of the clock.
		Now() time.Time

		// After returns a channel that receives the time once d has elapsed on the clock,
		// and a function that stops the timer and releases its resources. The stop
		// function reports whether it stopped the timer before it fired, like
		// time.Timer.Stop.
		After(d time.Duration) (<-chan time.Time, func() bool)
	}

	realClock struct{}

	clockInContextType int
)

var (
	clockInContextKey clockInContextType

	// getClockFromContext looks up the clock stored under clockInContextKey
	getClockFromContext = IC.AskValue[Clock](clockInContextKey)

	// getOrReal falls back to the real clock
	getOrReal = common.OptionGetOrElse(Real)
)

// Now returns the current time, as reported by package time.
func (realClock) Now() time.Time {
	return time.Now()
}

// After starts a timer of package time.
func (realClock) After(d time.Duration) (<-chan time.Time, func() bool) {
	t := time.NewTimer(d)
	return t.C, t.Stop
}

// Real returns the clock backed by package time.
//
// It is the clock that FromContext returns when the context carries no clock.
//
// Returns:
//   - Clock: The real clock
//
// See Also:
//   - FromContext: Reads the clock from a context
func Real() Clock {
	return realClock{}
}

// IsReal reports whether c is the clock returned by Real.
//
// Combinators use it to keep the exact behaviour of package time and context, for
// example context.WithTimeout, when no other clock is installed.
//
// Parameters:
//   - c: The clock to check
//
// Returns:
//   - bool: true if c is the real clock
func IsReal(c Clock) bool {
	_, ok := c.(realClock)
	return ok
}

// WithClock returns a Kleisli arrow that stores a clock in a context.
//
// The resulting pair contains a no-op cancel function and the new context, so that it can
// be used with the Local combinators, in the same way as logging.WithLogger.
//
// Parameters:
//   - c: The clock to store in the context
//
// Returns:
//   - pair.Kleisli[context.CancelFunc, context.Context, context.Context]: A function that
//     adds the clock to a context
//
// See Also:
//   - FromContext: Reads the clock back
func WithClock(c Clock) pair.Kleisli[context.CancelFunc, context.Context, context.Context] {
	return F.Flow2(
		IC.WithValue[Clock](clockInContextKey)(c),
		IC.NopCancel,
	)
}

// FromContext returns the clock stored in a context, or the real clock if there is none.
//
// Parameters:
//   - ctx: The context to read the clock from
//
// Returns:
//   - Clock: The clock of the context, or Real()
//
// See Also:
//   - WithClock: Stores a clock in a context
//   - Real: The fallback clock
func FromContext(ctx context.Context) Clock {
	return getOrReal(getClockFromContext(ctx))
}

// Now returns the current time of the clock in a context.
//
// Parameters:
//   - ctx: The context to read the clock from
//
// Returns:
//   - time.Time: The current time of the clock
//
// See Also:
//   - FromContext: Reads the clock from a context
func Now(ctx context.Context) time.Time {
	return FromContext(ctx).Now()
}

// Sleep waits for d on the clock in a context.
//
// The wait honours the context: if the context is cancelled before or while waiting,
// Sleep returns immediately with the cancellation cause, context.Cause(ctx).
//
// Parameters:
//   - ctx: The context that provides the clock and the cancellation
//   - d: The duration to wait
//
// Returns:
//   - error: nil after d has elapsed, the cancellation cause otherwise
//
// See Also:
//   - FromContext: Reads the clock from a context
func Sleep(ctx context.Context, d time.Duration) error {
	// quick check for cancellation, select picks randomly among ready cases
	if ctx.Err() != nil {
		return context.Cause(ctx)
	}
	c := FromContext(ctx)
	if IsReal(c) {
		// fast path, avoids allocating the stop function
		t := time.NewTimer(d)
		defer t.Stop()
		return wait(ctx, t.C)
	}
	ch, stop := c.After(d)
	defer stop()
	return wait(ctx, ch)
}

// wait blocks until ch receives or ctx is done
func wait(ctx context.Context, ch <-chan time.Time) error {
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}

// NowIO returns an IO that reads the current time of a clock.
//
// It feeds constructors that expect an IO[time.Time], such as those of the circuit
// breaker, from a Clock.
//
// Parameters:
//   - c: The clock to read
//
// Returns:
//   - io.IO[time.Time]: An IO that returns c.Now() each time it runs
func NowIO(c Clock) io.IO[time.Time] {
	return c.Now
}
