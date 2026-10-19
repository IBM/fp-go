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

/*
Package clock provides an injectable time source.

A Clock reads the current time and starts timers. The clock in use travels in the
context.Context, like the logger of package logging: WithClock stores a clock in a
context and FromContext reads it back. A context without a clock uses the real clock,
backed by package time.

Time-based combinators of fp-go read the clock from the context, among them Delay,
Timer, WithTimeout and Retrying of context/readerioresult and package effect. Tests
install the clock of package clock/testclock and advance it explicitly, so that code
which waits for minutes runs in microseconds and without sleeping.

	ctx := pair.Second(clock.WithClock(testclock.New(start))(context.Background()))
	now := clock.Now(ctx)
*/
package clock
