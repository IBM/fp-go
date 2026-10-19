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
	"fmt"
	"time"

	"github.com/IBM/fp-go/v2/clock"
	"github.com/IBM/fp-go/v2/clock/testclock"
	"github.com/IBM/fp-go/v2/pair"
)

// ExampleReal shows that a context without a clock uses the real clock.
func ExampleReal() {
	fmt.Println(clock.IsReal(clock.Real()))
	fmt.Println(clock.IsReal(clock.FromContext(context.Background())))
	// Output:
	// true
	// true
}

// ExampleIsReal distinguishes the real clock from a test clock.
func ExampleIsReal() {
	fmt.Println(clock.IsReal(testclock.New(start)))
	// Output: false
}

// ExampleWithClock installs a test clock in a context.
func ExampleWithClock() {
	ctx := pair.Second(clock.WithClock(testclock.New(start))(context.Background()))
	fmt.Println(clock.Now(ctx).Format(time.RFC3339))
	// Output: 2024-01-01T00:00:00Z
}

// ExampleFromContext reads the clock back from a context.
func ExampleFromContext() {
	tc := testclock.New(start)
	ctx := pair.Second(clock.WithClock(tc)(context.Background()))
	fmt.Println(clock.FromContext(ctx) == clock.Clock(tc))
	// Output: true
}

// ExampleNow reads the time of the clock in a context.
func ExampleNow() {
	tc := testclock.New(start)
	ctx := pair.Second(clock.WithClock(tc)(context.Background()))
	tc.Advance(90 * time.Minute)
	fmt.Println(clock.Now(ctx).Format(time.Kitchen))
	// Output: 1:30AM
}

// ExampleSleep waits for an hour on a test clock without sleeping.
func ExampleSleep() {
	tc := testclock.New(start)
	ctx := pair.Second(clock.WithClock(tc)(context.Background()))

	done := make(chan error)
	go func() { done <- clock.Sleep(ctx, time.Hour) }()

	tc.BlockUntil(1)
	tc.Advance(time.Hour)
	fmt.Println(<-done, clock.Now(ctx).Format(time.Kitchen))
	// Output: <nil> 1:00AM
}

// ExampleNowIO feeds an IO[time.Time] from a clock.
func ExampleNowIO() {
	now := clock.NowIO(testclock.New(start))
	fmt.Println(now().Year())
	// Output: 2024
}
