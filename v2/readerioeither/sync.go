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

package readerioeither

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/readerio"
)

// WithLock executes a ReaderIOEither operation within the scope of a lock.
// The lock is acquired before the operation executes and released after it completes,
// regardless of whether the operation succeeds or fails.
//
// This is useful for ensuring thread-safe access to shared resources or for
// implementing critical sections in concurrent code.
//
// Type parameters:
//   - R: The context type
//   - E: The error type
//   - A: The value type
//
// Parameters:
//   - lock: A function that acquires a lock and returns a CancelFunc to release it
//
// Returns:
//
//	An Operator that wraps the computation with lock acquisition and release
//
// Example:
//
//	var mu sync.Mutex
//	safeFetch := F.Pipe1(
//	    fetchData(),
//	    WithLock[Config, error, Data](func() context.CancelFunc {
//	        mu.Lock()
//	        return func() { mu.Unlock() }
//	    }),
//	)
//
//go:inline
func WithLock[R, E, A any](lock func() context.CancelFunc) Operator[R, E, A, A] {
	return readerio.WithLock[R, either.Either[E, A]](lock)
}
