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

package readerioresult

import (
	"context"
	"strconv"
	"testing"

	"github.com/IBM/fp-go/v2/context/ioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestPromapBasic tests basic Promap functionality
func TestPromapBasic(t *testing.T) {
	t.Run("transform both context and output", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[int] {
			return func() R.Result[int] {
				if v := ctx.Value("key"); v != nil {
					return R.Of(v.(int))
				}
				return R.Of(0)
			}
		}

		addKey := func(ctx context.Context) pair.Pair[context.CancelFunc, context.Context] {
			newCtx := context.WithValue(ctx, "key", 42)
			return pair.MakePair(context.CancelFunc(func() {}), newCtx)
		}
		toString := strconv.Itoa

		adapted := Promap(addKey, toString)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, R.Of("42"), result)
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("context transformation", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[int] {
			return func() R.Result[int] {
				if v := ctx.Value("key"); v != nil {
					return R.Of(v.(int))
				}
				return R.Of(0)
			}
		}

		addKey := func(ctx context.Context) pair.Pair[context.CancelFunc, context.Context] {
			newCtx := context.WithValue(ctx, "key", 100)
			return pair.MakePair(context.CancelFunc(func() {}), newCtx)
		}

		adapted := Contramap[int](addKey)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, R.Of(100), result)
	})
}

// TestLocalBasic tests basic Local functionality
func TestLocalBasic(t *testing.T) {
	t.Run("adds value to context", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				if v := ctx.Value("user"); v != nil {
					return R.Of(v.(string))
				}
				return R.Of("unknown")
			}
		}

		addUser := func(ctx context.Context) pair.Pair[context.CancelFunc, context.Context] {
			newCtx := context.WithValue(ctx, "user", "Alice")
			return pair.MakePair(context.CancelFunc(func() {}), newCtx)
		}

		adapted := Local[string](addUser)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, R.Of("Alice"), result)
	})
}

// TestLocalIOK_Success tests LocalIOK with successful context transformation
func TestLocalIOK_Success(t *testing.T) {
	t.Run("transforms context with IO effect", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				if v := ctx.Value("user"); v != nil {
					return R.Of(v.(string))
				}
				return R.Of("unknown")
			}
		}

		addUser := func(ctx context.Context) io.IO[ContextCancel] {
			return func() ContextCancel {
				newCtx := context.WithValue(ctx, "user", "Bob")
				return pair.MakePair(context.CancelFunc(func() {}), newCtx)
			}
		}

		adapted := LocalIOK[string](addUser)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, R.Of("Bob"), result)
	})

	t.Run("preserves original value type", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[int] {
			return func() R.Result[int] {
				if v := ctx.Value("count"); v != nil {
					return R.Of(v.(int))
				}
				return R.Of(0)
			}
		}

		addCount := func(ctx context.Context) io.IO[ContextCancel] {
			return func() ContextCancel {
				newCtx := context.WithValue(ctx, "count", 42)
				return pair.MakePair(context.CancelFunc(func() {}), newCtx)
			}
		}

		adapted := LocalIOK[int](addCount)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, R.Of(42), result)
	})
}

// TestLocalIOK_CancelledContext tests LocalIOK with cancelled context
func TestLocalIOK_CancelledContext(t *testing.T) {
	t.Run("returns error when context is cancelled", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				return R.Of("should not reach here")
			}
		}

		addUser := func(ctx context.Context) io.IO[ContextCancel] {
			return func() ContextCancel {
				newCtx := context.WithValue(ctx, "user", "Charlie")
				return pair.MakePair(context.CancelFunc(func() {}), newCtx)
			}
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		adapted := LocalIOK[string](addUser)(getValue)
		result := adapted(ctx)()

		assert.True(t, R.IsLeft(result))
	})
}

// TestLocalIOK_CancelFuncCalled tests that CancelFunc is properly called
func TestLocalIOK_CancelFuncCalled(t *testing.T) {
	t.Run("calls cancel function after execution", func(t *testing.T) {
		cancelCalled := false

		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				return R.Of("test")
			}
		}

		addUser := func(ctx context.Context) io.IO[ContextCancel] {
			return func() ContextCancel {
				newCtx := context.WithValue(ctx, "user", "Dave")
				cancelFunc := context.CancelFunc(func() {
					cancelCalled = true
				})
				return pair.MakePair(cancelFunc, newCtx)
			}
		}

		adapted := LocalIOK[string](addUser)(getValue)
		_ = adapted(t.Context())()

		assert.True(t, cancelCalled, "cancel function should be called")
	})
}

// TestLocalIOResultK_Success tests LocalIOResultK with successful context transformation
func TestLocalIOResultK_Success(t *testing.T) {
	t.Run("transforms context with IOResult effect", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				if v := ctx.Value("role"); v != nil {
					return R.Of(v.(string))
				}
				return R.Of("guest")
			}
		}

		addRole := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
			return func() R.Result[ContextCancel] {
				newCtx := context.WithValue(ctx, "role", "admin")
				return R.Of(pair.MakePair(context.CancelFunc(func() {}), newCtx))
			}
		}

		adapted := LocalIOResultK[string](addRole)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, R.Of("admin"), result)
	})

	t.Run("preserves original value type", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[int] {
			return func() R.Result[int] {
				if v := ctx.Value("score"); v != nil {
					return R.Of(v.(int))
				}
				return R.Of(0)
			}
		}

		addScore := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
			return func() R.Result[ContextCancel] {
				newCtx := context.WithValue(ctx, "score", 100)
				return R.Of(pair.MakePair(context.CancelFunc(func() {}), newCtx))
			}
		}

		adapted := LocalIOResultK[int](addScore)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, R.Of(100), result)
	})
}

// TestLocalIOResultK_Failure tests LocalIOResultK with failed context transformation
func TestLocalIOResultK_Failure(t *testing.T) {
	t.Run("propagates transformation error", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				return R.Of("should not reach here")
			}
		}

		failTransform := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
			return func() R.Result[ContextCancel] {
				return R.Left[ContextCancel](assert.AnError)
			}
		}

		adapted := LocalIOResultK[string](failTransform)(getValue)
		result := adapted(t.Context())()

		assert.True(t, R.IsLeft(result))
		_, err := R.UnwrapError(result)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("does not execute original computation on transformation failure", func(t *testing.T) {
		executed := false

		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				executed = true
				return R.Of("should not execute")
			}
		}

		failTransform := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
			return func() R.Result[ContextCancel] {
				return R.Left[ContextCancel](assert.AnError)
			}
		}

		adapted := LocalIOResultK[string](failTransform)(getValue)
		_ = adapted(t.Context())()

		assert.False(t, executed, "original computation should not execute")
	})
}

// TestLocalIOResultK_CancelledContext tests LocalIOResultK with cancelled context
func TestLocalIOResultK_CancelledContext(t *testing.T) {
	t.Run("returns error when context is cancelled", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				return R.Of("should not reach here")
			}
		}

		addRole := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
			return func() R.Result[ContextCancel] {
				newCtx := context.WithValue(ctx, "role", "user")
				return R.Of(pair.MakePair(context.CancelFunc(func() {}), newCtx))
			}
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		adapted := LocalIOResultK[string](addRole)(getValue)
		result := adapted(ctx)()

		assert.True(t, R.IsLeft(result))
	})
}

// TestLocalIOResultK_CancelFuncCalled tests that CancelFunc is properly called
func TestLocalIOResultK_CancelFuncCalled(t *testing.T) {
	t.Run("calls cancel function after successful execution", func(t *testing.T) {
		cancelCalled := false

		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				return R.Of("test")
			}
		}

		addRole := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
			return func() R.Result[ContextCancel] {
				newCtx := context.WithValue(ctx, "role", "user")
				cancelFunc := context.CancelFunc(func() {
					cancelCalled = true
				})
				return R.Of(pair.MakePair(cancelFunc, newCtx))
			}
		}

		adapted := LocalIOResultK[string](addRole)(getValue)
		_ = adapted(t.Context())()

		assert.True(t, cancelCalled, "cancel function should be called")
	})

	t.Run("does not call cancel function on transformation failure", func(t *testing.T) {
		cancelCalled := false

		getValue := func(ctx context.Context) IOResult[string] {
			return func() R.Result[string] {
				return R.Of("test")
			}
		}

		failTransform := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
			return func() R.Result[ContextCancel] {
				cancelFunc := context.CancelFunc(func() {
					cancelCalled = true
				})
				_ = cancelFunc // avoid unused warning
				return R.Left[ContextCancel](assert.AnError)
			}
		}

		adapted := LocalIOResultK[string](failTransform)(getValue)
		_ = adapted(t.Context())()

		assert.False(t, cancelCalled, "cancel function should not be called on failure")
	})
}

// TestLocalIOResultK_Integration tests integration with other operations
func TestLocalIOResultK_Integration(t *testing.T) {
	t.Run("composes with Map", func(t *testing.T) {
		getValue := func(ctx context.Context) IOResult[int] {
			return func() R.Result[int] {
				if v := ctx.Value("value"); v != nil {
					return R.Of(v.(int))
				}
				return R.Of(0)
			}
		}

		addValue := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
			return func() R.Result[ContextCancel] {
				newCtx := context.WithValue(ctx, "value", 10)
				return R.Of(pair.MakePair(context.CancelFunc(func() {}), newCtx))
			}
		}

		double := func(x int) int { return x * 2 }

		adapted := F.Flow2(
			LocalIOResultK[int](addValue),
			Map(double),
		)(getValue)
		result := adapted(t.Context())()

		assert.Equal(t, R.Of(20), result)
	})
}
