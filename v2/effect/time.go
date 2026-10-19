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

package effect

import (
	"time"

	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/reader"
)

// Delay waits for the given duration before running an effect.
//
// The wait honours the context: if the context is cancelled while waiting, the effect does
// not run and fails with the cancellation cause.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//   - A: The type of the value produced by the effect
//
// Parameters:
//   - delay: The duration to wait
//
// Returns:
//   - Operator[C, A, A]: A function that delays an effect
//
// See Also:
//   - Timer: An effect that waits and produces the current time
func Delay[C, A any](delay time.Duration) Operator[C, A, A] {
	return readerreaderioresult.Delay[C, A](delay)
}

// Timer returns an effect that waits for the given duration and produces the time at which
// the wait ended.
//
// This corresponds to Effect.sleep in effect-ts, with the current time as the result.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//
// Parameters:
//   - delay: The duration to wait
//
// Returns:
//   - Effect[C, time.Time]: An effect that produces the current time after the wait
//
// See Also:
//   - Delay: Delays an existing effect
func Timer[C any](delay time.Duration) Effect[C, time.Time] {
	return FromThunk[C](thunk.Timer(delay))
}

// WithTimeout runs an effect with a deadline.
//
// The context passed to the effect is cancelled after the timeout. If the effect does not
// complete in time, it fails with context.DeadlineExceeded, provided it honours the
// context, as all effects built with this package do.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//   - A: The type of the value produced by the effect
//
// Parameters:
//   - timeout: The maximum duration of the effect
//
// Returns:
//   - Operator[C, A, A]: A function that applies the timeout to an effect
//
// See Also:
//   - Retrying: Retries an effect that failed, for example because of a timeout
func WithTimeout[C, A any](timeout time.Duration) Operator[C, A, A] {
	return reader.Map[C](thunk.WithTimeout[A](timeout))
}
