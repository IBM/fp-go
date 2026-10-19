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
	"fmt"
	"time"

	"github.com/IBM/fp-go/v2/clock"
	"github.com/IBM/fp-go/v2/clock/testclock"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
)

// ExampleDelay waits for an hour on a test clock, without sleeping.
func ExampleDelay() {
	tc := testclock.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	ctx := pair.Second(clock.WithClock(tc)(context.Background()))

	done := make(chan Result[string])
	go func() { done <- F.Pipe1(Of("done"), Delay[string](time.Hour))(ctx)() }()

	tc.BlockUntil(1)
	tc.Advance(time.Hour)
	fmt.Println(<-done)
	// Output: Right[string](done)
}

// ExampleTimer produces the time of the test clock after the wait.
func ExampleTimer() {
	tc := testclock.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	ctx := pair.Second(clock.WithClock(tc)(context.Background()))

	done := make(chan Result[time.Time])
	go func() { done <- Timer(time.Minute)(ctx)() }()

	tc.BlockUntil(1)
	tc.Advance(time.Minute)
	fmt.Println(F.Pipe1(<-done, result.Map(F.Bind2nd(time.Time.Format, time.Kitchen))))
	// Output: Right[string](12:01AM)
}

// ExampleWithTimeout lets a timeout on a test clock expire.
func ExampleWithTimeout() {
	tc := testclock.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	ctx := pair.Second(clock.WithClock(tc)(context.Background()))

	slow := F.Pipe2(
		Of("slow"),
		Delay[string](time.Hour),
		WithTimeout[string](time.Minute),
	)

	done := make(chan Result[string])
	go func() { done <- slow(ctx)() }()

	// the timeout and the delay are pending
	tc.BlockUntil(2)
	tc.Advance(time.Minute)
	fmt.Println(<-done)
	// Output: Left[context.deadlineExceededError](context deadline exceeded)
}
