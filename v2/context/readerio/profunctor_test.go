// Copyright (c) 2025 IBM Corp.
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

package readerio

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// TestPromapBasic tests basic Promap functionality
func TestPromapBasic(t *testing.T) {
	t.Run("transform both context and output", func(t *testing.T) {
		// ReaderIO that reads a value from context
		getValue := func(ctx context.Context) IO[int] {
			return func() int {
				if v := ctx.Value("key"); v != nil {
					return v.(int)
				}
				return 0
			}
		}

		// Transform context and result
		addKey := func(ctx context.Context) ContextCancel {
			newCtx := context.WithValue(ctx, "key", 42)
			return pair.MakePair[context.CancelFunc](func() {}, newCtx)
		}
		toString := strconv.Itoa

		adapted := Promap(addKey, toString)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, "42", result)
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("context transformation", func(t *testing.T) {
		getValue := func(ctx context.Context) IO[int] {
			return func() int {
				if v := ctx.Value("key"); v != nil {
					return v.(int)
				}
				return 0
			}
		}

		addKey := func(ctx context.Context) ContextCancel {
			newCtx := context.WithValue(ctx, "key", 100)
			return pair.MakePair[context.CancelFunc](func() {}, newCtx)
		}

		adapted := Contramap[int](addKey)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, 100, result)
	})
}

// TestLocalBasic tests basic Local functionality
func TestLocalBasic(t *testing.T) {
	t.Run("adds timeout to context", func(t *testing.T) {
		getValue := func(ctx context.Context) IO[bool] {
			return func() bool {
				_, hasDeadline := ctx.Deadline()
				return hasDeadline
			}
		}

		addTimeout := func(ctx context.Context) ContextCancel {
			newCtx, cancelFct := context.WithTimeout(ctx, time.Second)
			return pair.MakePair(cancelFct, newCtx)
		}

		adapted := Local[bool](addTimeout)(getValue)
		result := adapted(t.Context())()

		assert.True(t, result)
	})
}

// TestLocalIOKBasic tests basic LocalIOK functionality
func TestLocalIOKBasic(t *testing.T) {
	t.Run("context transformation with IO effect", func(t *testing.T) {
		getValue := func(ctx context.Context) IO[string] {
			return func() string {
				if v := ctx.Value("key"); v != nil {
					return v.(string)
				}
				return "default"
			}
		}

		// Context transformation wrapped in IO effect
		addKeyIO := func(ctx context.Context) IO[ContextCancel] {
			return func() ContextCancel {
				// Simulate side effect (e.g., loading config)
				newCtx := context.WithValue(ctx, "key", "loaded-value")
				return pair.MakePair[context.CancelFunc](func() {}, newCtx)
			}
		}

		adapted := LocalIOK[string](addKeyIO)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, "loaded-value", result)
	})

	t.Run("cleanup function is called", func(t *testing.T) {
		cleanupCalled := false

		getValue := func(ctx context.Context) IO[int] {
			return func() int {
				if v := ctx.Value("value"); v != nil {
					return v.(int)
				}
				return 0
			}
		}

		addValueIO := func(ctx context.Context) IO[ContextCancel] {
			return func() ContextCancel {
				newCtx := context.WithValue(ctx, "value", 42)
				cleanup := context.CancelFunc(func() {
					cleanupCalled = true
				})
				return pair.MakePair(cleanup, newCtx)
			}
		}

		adapted := LocalIOK[int](addValueIO)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, 42, result)
		assert.True(t, cleanupCalled, "cleanup function should be called")
	})

	t.Run("works with timeout context", func(t *testing.T) {
		getValue := func(ctx context.Context) IO[bool] {
			return func() bool {
				_, hasDeadline := ctx.Deadline()
				return hasDeadline
			}
		}

		addTimeoutIO := func(ctx context.Context) IO[ContextCancel] {
			return func() ContextCancel {
				newCtx, cancelFct := context.WithTimeout(ctx, time.Second)
				return pair.MakePair(cancelFct, newCtx)
			}
		}

		adapted := LocalIOK[bool](addTimeoutIO)(getValue)
		result := adapted(t.Context())()

		assert.True(t, result, "context should have deadline")
	})
}
