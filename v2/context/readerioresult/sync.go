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

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
)

// WithLock executes the provided IO operation in the scope of a lock.
// This ensures that the operation is executed with mutual exclusion, preventing concurrent access.
//
// The lock is acquired before the operation and released after it completes (or fails).
// The lock parameter should return a CancelFunc that releases the lock when called.
//
// Parameters:
//   - lock: ReaderIOResult that acquires a lock and returns a CancelFunc to release it
//
// Returns a function that wraps a ReaderIOResult with lock protection.
//
// Example:
//
//	mutex := &sync.Mutex{}
//	lock := TryCatch(func(ctx context.Context) func() (context.CancelFunc, error) {
//	    return func() (context.CancelFunc, error) {
//	        mutex.Lock()
//	        return func() { mutex.Unlock() }, nil
//	    }
//	})
//	protectedOp := WithLock(lock)(myOperation)
func WithLock[A any](lock ReaderIOResult[context.CancelFunc]) Operator[A, A] {
	return function.Flow2(
		function.Constant1[context.CancelFunc, ReaderIOResult[A]],
		WithResource[A](lock, function.Flow2(
			io.FromImpure[context.CancelFunc],
			FromIO[Void],
		)),
	)
}
