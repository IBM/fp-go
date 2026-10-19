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

package testclock_test

import (
	"fmt"
	"time"

	"github.com/IBM/fp-go/v2/clock/testclock"
)

var start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// ExampleNew creates a test clock at a fixed time.
func ExampleNew() {
	c := testclock.New(start)
	fmt.Println(c.Now().Format(time.RFC3339))
	// Output: 2024-01-01T00:00:00Z
}

// ExampleTestClock_Now shows that the clock does not move on its own.
func ExampleTestClock_Now() {
	c := testclock.New(start)
	fmt.Println(c.Now().Equal(c.Now()))
	// Output: true
}

// ExampleTestClock_After starts a timer that fires when the clock is advanced.
func ExampleTestClock_After() {
	c := testclock.New(start)
	ch, _ := c.After(time.Minute)
	c.Advance(time.Minute)
	fmt.Println((<-ch).Format(time.Kitchen))
	// Output: 12:01AM
}

// ExampleTestClock_Advance fires the timers that become due.
func ExampleTestClock_Advance() {
	c := testclock.New(start)
	short, _ := c.After(time.Second)
	long, _ := c.After(time.Hour)

	c.Advance(time.Minute)
	fmt.Println((<-short).Sub(start))

	c.Advance(time.Hour)
	fmt.Println((<-long).Sub(start))
	// Output:
	// 1s
	// 1h0m0s
}

// ExampleTestClock_Set moves the clock to an absolute time.
func ExampleTestClock_Set() {
	c := testclock.New(start)
	c.Set(start.AddDate(0, 0, 1))
	fmt.Println(c.Now().Format(time.DateOnly))
	// Output: 2024-01-02
}

// ExampleTestClock_BlockUntil waits until another goroutine is sleeping before advancing.
func ExampleTestClock_BlockUntil() {
	c := testclock.New(start)
	done := make(chan time.Time)
	go func() {
		ch, stop := c.After(time.Hour)
		defer stop()
		done <- <-ch
	}()

	c.BlockUntil(1)
	c.Advance(time.Hour)
	fmt.Println((<-done).Format(time.Kitchen))
	// Output: 1:00AM
}
