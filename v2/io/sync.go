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

package io

import (
	"context"
)

// WithLock executes the provided IO operation in the scope of a lock.
// The lock parameter should be an IO that acquires a lock and returns a function to release it.
//
// This ensures that the operation is executed with exclusive access to a shared resource,
// and the lock is always released even if the operation panics.
//
// Example:
//
//	mutex := &sync.Mutex{}
//	lock := io.FromImpure(func() context.CancelFunc {
//	    mutex.Lock()
//	    return func() { mutex.Unlock() }
//	})
//
//	safeOperation := io.WithLock(lock)(dangerousOperation)
//	result := safeOperation()
func WithLock[A any](lock IO[context.CancelFunc]) Operator[A, A] {
	return func(fa IO[A]) IO[A] {
		return func() A {
			defer lock()()
			return fa()
		}
	}
}
