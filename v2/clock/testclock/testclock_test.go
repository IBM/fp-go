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

package testclock

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fired(ch <-chan time.Time) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func TestTestClock_Success(t *testing.T) {
	t.Run("Now starts at the start time and follows Advance", func(t *testing.T) {
		c := New(start)
		assert.Equal(t, start, c.Now())
		c.Advance(time.Minute)
		assert.Equal(t, start.Add(time.Minute), c.Now())
	})

	t.Run("a timer fires once it is due and receives its due time", func(t *testing.T) {
		c := New(start)
		ch, _ := c.After(time.Minute)

		c.Advance(59 * time.Second)
		assert.False(t, fired(ch))

		c.Advance(time.Second)
		assert.Equal(t, start.Add(time.Minute), <-ch)
	})

	t.Run("Set fires the due timers", func(t *testing.T) {
		c := New(start)
		ch, _ := c.After(time.Hour)
		c.Set(start.Add(2 * time.Hour))
		assert.Equal(t, start.Add(time.Hour), <-ch)
		assert.Equal(t, start.Add(2*time.Hour), c.Now())
	})

	t.Run("one advance fires all due timers with their own due times", func(t *testing.T) {
		c := New(start)
		ch3, _ := c.After(3 * time.Second)
		ch1, _ := c.After(time.Second)
		ch2, _ := c.After(2 * time.Second)
		chLater, _ := c.After(time.Minute)

		c.Advance(3 * time.Second)
		assert.Equal(t, start.Add(time.Second), <-ch1)
		assert.Equal(t, start.Add(2*time.Second), <-ch2)
		assert.Equal(t, start.Add(3*time.Second), <-ch3)
		assert.False(t, fired(chLater))
		assert.Equal(t, 1, len(c.pending))
	})

	t.Run("due timers are ordered by due time, then start order", func(t *testing.T) {
		a := &timer{when: start.Add(time.Second), seq: 2}
		b := &timer{when: start.Add(time.Second), seq: 1}
		e := &timer{when: start, seq: 3}
		due, rest := splitDue([]*timer{a, b, e}, start.Add(time.Second))
		assert.Empty(t, rest)
		assert.Equal(t, -1, compareTimers(e, b))
		assert.Equal(t, -1, compareTimers(b, a))
		assert.Equal(t, 1, compareTimers(a, b))
		assert.Len(t, due, 3)
	})
}

func TestTestClock_Stop(t *testing.T) {
	t.Run("stop removes a pending timer", func(t *testing.T) {
		c := New(start)
		ch, stop := c.After(time.Minute)
		assert.True(t, stop())
		assert.False(t, stop())
		c.Advance(time.Hour)
		assert.False(t, fired(ch))
	})

	t.Run("stop after firing reports false", func(t *testing.T) {
		c := New(start)
		_, stop := c.After(time.Minute)
		c.Advance(time.Minute)
		assert.False(t, stop())
	})
}

func TestTestClock_EdgeCases(t *testing.T) {
	t.Run("non-positive durations fire immediately", func(t *testing.T) {
		c := New(start)
		ch0, stop0 := c.After(0)
		chn, _ := c.After(-time.Second)
		assert.Equal(t, start, <-ch0)
		assert.Equal(t, start, <-chn)
		assert.False(t, stop0())
		assert.Equal(t, 0, len(c.pending))
	})

	t.Run("moving backwards fires nothing", func(t *testing.T) {
		c := New(start)
		ch, _ := c.After(time.Minute)
		c.Set(start.Add(-time.Hour))
		assert.False(t, fired(ch))
		c.Advance(time.Hour + time.Minute)
		assert.True(t, fired(ch))
	})

	t.Run("BlockUntil(0) returns immediately", func(t *testing.T) {
		New(start).BlockUntil(0)
	})
}

func TestTestClock_Concurrency(t *testing.T) {
	t.Run("BlockUntil waits for the given number of timers", func(t *testing.T) {
		c := New(start)
		const n = 10

		var wg sync.WaitGroup
		wg.Add(n)
		for range n {
			go func() {
				defer wg.Done()
				ch, stop := c.After(time.Second)
				defer stop()
				<-ch
			}()
		}

		c.BlockUntil(n)
		c.Advance(time.Second)
		wg.Wait()
		assert.Equal(t, 0, len(c.pending))
	})
}
