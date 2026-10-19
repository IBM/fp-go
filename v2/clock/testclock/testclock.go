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

// Package testclock provides a clock.Clock for tests that advances only when told to.
//
// A TestClock starts at a fixed time and never moves on its own. Timers started with
// After fire when the test calls Advance or Set and moves the clock past their due time.
// BlockUntil waits until the code under test has started a given number of timers, so a
// test can advance the clock exactly when the code is waiting, without sleeping.
package testclock

import (
	"slices"
	"sync"
	"time"

	"github.com/IBM/fp-go/v2/clock"
)

type (
	// TestClock is a manually driven clock.Clock. It is safe for concurrent use.
	TestClock struct {
		mu      sync.Mutex
		changed *sync.Cond
		now     time.Time
		seq     uint64
		pending []*timer
	}

	timer struct {
		when time.Time
		seq  uint64
		ch   chan time.Time
	}
)

// compile time check that TestClock implements clock.Clock
var _ clock.Clock = (*TestClock)(nil)

// New creates a test clock that starts at the given time.
//
// Parameters:
//   - start: The initial time of the clock
//
// Returns:
//   - *TestClock: A clock that only moves with Advance and Set
func New(start time.Time) *TestClock {
	c := &TestClock{now: start}
	c.changed = sync.NewCond(&c.mu)
	return c
}

// Now returns the current time of the clock.
func (c *TestClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// After starts a timer that fires once the clock has advanced by d.
//
// A timer with a duration of zero or less fires immediately. The returned stop function
// removes a pending timer and reports whether it was still pending.
func (c *TestClock) After(d time.Duration) (<-chan time.Time, func() bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	t := &timer{when: c.now.Add(d), seq: c.seq, ch: make(chan time.Time, 1)}
	c.seq++

	if d <= 0 {
		t.ch <- c.now
		return t.ch, func() bool { return false }
	}

	c.pending = append(c.pending, t)
	c.changed.Broadcast()

	return t.ch, func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		idx := slices.Index(c.pending, t)
		if idx < 0 {
			return false
		}
		c.pending = slices.Delete(c.pending, idx, idx+1)
		c.changed.Broadcast()
		return true
	}
}

// Advance moves the clock forward by d and fires every timer that is due, in the order
// of their due times. Timers with the same due time fire in the order they were started.
//
// Parameters:
//   - d: The duration to advance the clock by
func (c *TestClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setLocked(c.now.Add(d))
}

// Set moves the clock to t and fires every timer that is due, in the order of their due
// times. Moving the clock backwards fires no timers.
//
// Parameters:
//   - t: The new time of the clock
func (c *TestClock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setLocked(t)
}

// BlockUntil waits until at least n timers are pending, that is started and neither
// fired nor stopped.
//
// Parameters:
//   - n: The number of pending timers to wait for
func (c *TestClock) BlockUntil(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for len(c.pending) < n {
		c.changed.Wait()
	}
}

func (c *TestClock) setLocked(now time.Time) {
	c.now = now

	due, rest := splitDue(c.pending, now)
	if len(due) == 0 {
		return
	}
	c.pending = rest

	slices.SortFunc(due, compareTimers)
	for _, t := range due {
		t.ch <- t.when
	}
	c.changed.Broadcast()
}

func splitDue(timers []*timer, now time.Time) (due, rest []*timer) {
	for _, t := range timers {
		if t.when.After(now) {
			rest = append(rest, t)
		} else {
			due = append(due, t)
		}
	}
	return due, rest
}

func compareTimers(a, b *timer) int {
	if c := a.when.Compare(b.when); c != 0 {
		return c
	}
	if a.seq < b.seq {
		return -1
	}
	return 1
}
