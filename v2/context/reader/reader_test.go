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

package reader

import (
	"context"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
)

// TestWithValue_BasicUsage tests basic usage of WithValue
func TestWithValue_BasicUsage(t *testing.T) {
	t.Run("adds string value to context", func(t *testing.T) {
		setUserID := WithValue[string]("userID")
		ctx := context.Background()

		newCtx := setUserID("user123")(ctx)

		value := newCtx.Value("userID")
		assert.Equal(t, "user123", value)
	})

	t.Run("adds int value to context", func(t *testing.T) {
		setTimeout := WithValue[int]("timeout")
		ctx := context.Background()

		newCtx := setTimeout(30)(ctx)

		value := newCtx.Value("timeout")
		assert.Equal(t, 30, value)
	})

	t.Run("adds struct value to context", func(t *testing.T) {
		type User struct {
			ID   string
			Name string
		}

		setUser := WithValue[User]("user")
		ctx := context.Background()
		user := User{ID: "123", Name: "Alice"}

		newCtx := setUser(user)(ctx)

		value := newCtx.Value("user")
		assert.Equal(t, user, value)
	})
}

// TestWithValue_CustomKeyType tests WithValue with custom key types
func TestWithValue_CustomKeyType(t *testing.T) {
	type contextKey string

	const (
		userKey    contextKey = "user"
		sessionKey contextKey = "session"
	)

	t.Run("uses custom key type", func(t *testing.T) {
		setUser := WithValue[string](userKey)
		ctx := context.Background()

		newCtx := setUser("Alice")(ctx)

		value := newCtx.Value(userKey)
		assert.Equal(t, "Alice", value)
	})

	t.Run("different custom keys don't conflict", func(t *testing.T) {
		setUser := WithValue[string](userKey)
		setSession := WithValue[string](sessionKey)
		ctx := context.Background()

		ctx = setUser("Alice")(ctx)
		ctx = setSession("session-token")(ctx)

		assert.Equal(t, "Alice", ctx.Value(userKey))
		assert.Equal(t, "session-token", ctx.Value(sessionKey))
	})
}

// TestWithValue_Chaining tests chaining multiple WithValue operations
func TestWithValue_Chaining(t *testing.T) {
	t.Run("chains multiple values using composition", func(t *testing.T) {
		// Compose multiple WithValue operations
		enrichContext := func(ctx context.Context) context.Context {
			ctx = WithValue[string]("userID")("user123")(ctx)
			ctx = WithValue[string]("requestID")("req456")(ctx)
			ctx = WithValue[int]("timeout")(30)(ctx)
			return ctx
		}

		ctx := context.Background()
		enrichedCtx := enrichContext(ctx)

		assert.Equal(t, "user123", enrichedCtx.Value("userID"))
		assert.Equal(t, "req456", enrichedCtx.Value("requestID"))
		assert.Equal(t, 30, enrichedCtx.Value("timeout"))
	})

	t.Run("chains values sequentially", func(t *testing.T) {
		ctx := context.Background()

		ctx = WithValue[string]("key1")("value1")(ctx)
		ctx = WithValue[string]("key2")("value2")(ctx)
		ctx = WithValue[string]("key3")("value3")(ctx)

		assert.Equal(t, "value1", ctx.Value("key1"))
		assert.Equal(t, "value2", ctx.Value("key2"))
		assert.Equal(t, "value3", ctx.Value("key3"))
	})
}

// TestWithValue_ContextImmutability tests that original context is not modified
func TestWithValue_ContextImmutability(t *testing.T) {
	t.Run("original context is not modified", func(t *testing.T) {
		originalCtx := context.Background()
		setUserID := WithValue[string]("userID")

		newCtx := setUserID("user123")(originalCtx)

		// Original context should not have the value
		assert.Nil(t, originalCtx.Value("userID"))
		// New context should have the value
		assert.Equal(t, "user123", newCtx.Value("userID"))
	})

	t.Run("parent context values are preserved", func(t *testing.T) {
		parentCtx := context.WithValue(context.Background(), "parent", "parentValue")
		setChild := WithValue[string]("child")

		childCtx := setChild("childValue")(parentCtx)

		// Both parent and child values should be accessible
		assert.Equal(t, "parentValue", childCtx.Value("parent"))
		assert.Equal(t, "childValue", childCtx.Value("child"))
		// Parent context should not have child value
		assert.Nil(t, parentCtx.Value("child"))
	})
}

// TestWithValue_OverwritingValues tests overwriting existing values
func TestWithValue_OverwritingValues(t *testing.T) {
	t.Run("overwrites existing value with same key", func(t *testing.T) {
		ctx := context.Background()
		setUserID := WithValue[string]("userID")

		ctx = setUserID("user123")(ctx)
		assert.Equal(t, "user123", ctx.Value("userID"))

		ctx = setUserID("user456")(ctx)
		assert.Equal(t, "user456", ctx.Value("userID"))
	})

	t.Run("child context shadows parent value", func(t *testing.T) {
		parentCtx := context.WithValue(context.Background(), "key", "parent")
		setKey := WithValue[string]("key")

		childCtx := setKey("child")(parentCtx)

		// Child context should have the new value
		assert.Equal(t, "child", childCtx.Value("key"))
		// Parent context should still have the old value
		assert.Equal(t, "parent", parentCtx.Value("key"))
	})
}

// TestWithValue_NilValues tests handling of nil values
func TestWithValue_NilValues(t *testing.T) {
	t.Run("stores nil pointer value", func(t *testing.T) {
		type User struct {
			Name string
		}

		setUser := WithValue[*User]("user")
		ctx := context.Background()

		newCtx := setUser(nil)(ctx)

		value := newCtx.Value("user")
		assert.Nil(t, value)
	})

	t.Run("stores nil interface value", func(t *testing.T) {
		setData := WithValue[interface{}]("data")
		ctx := context.Background()

		newCtx := setData(nil)(ctx)

		value := newCtx.Value("data")
		assert.Nil(t, value)
	})
}

// TestWithValue_ComplexTypes tests WithValue with complex types
func TestWithValue_ComplexTypes(t *testing.T) {
	t.Run("stores slice value", func(t *testing.T) {
		setTags := WithValue[[]string]("tags")
		ctx := context.Background()
		tags := []string{"go", "functional", "programming"}

		newCtx := setTags(tags)(ctx)

		value := newCtx.Value("tags")
		assert.Equal(t, tags, value)
	})

	t.Run("stores map value", func(t *testing.T) {
		setMetadata := WithValue[map[string]int]("metadata")
		ctx := context.Background()
		metadata := map[string]int{"count": 42, "limit": 100}

		newCtx := setMetadata(metadata)(ctx)

		value := newCtx.Value("metadata")
		assert.Equal(t, metadata, value)
	})

	t.Run("stores function value", func(t *testing.T) {
		type Handler func(string) string

		setHandler := WithValue[Handler]("handler")
		ctx := context.Background()
		handler := func(s string) string { return "handled: " + s }

		newCtx := setHandler(handler)(ctx)

		value := newCtx.Value("handler")
		assert.NotNil(t, value)
		// Verify it's a function by calling it
		if h, ok := value.(Handler); ok {
			assert.Equal(t, "handled: test", h("test"))
		} else {
			t.Fatal("Expected handler function")
		}
	})
}

// TestWithValue_Integration tests integration with other Reader operations
func TestWithValue_Integration(t *testing.T) {
	t.Run("integrates with Reader Map", func(t *testing.T) {
		// Create a reader that adds a value and then extracts it
		pipeline := F.Pipe1(
			WithValue[string]("userID")("user123"),
			R.Map[context.Context](func(ctx context.Context) string {
				return ctx.Value("userID").(string)
			}),
		)

		ctx := context.Background()
		result := pipeline(ctx)

		assert.Equal(t, "user123", result)
	})

	t.Run("integrates with Reader Chain", func(t *testing.T) {
		// Chain multiple context enrichments
		pipeline := F.Pipe2(
			WithValue[string]("step")("1"),
			R.Chain(func(ctx context.Context) Reader[context.Context] {
				step := ctx.Value("step").(string)
				return WithValue[string]("result")("step " + step + " complete")
			}),
			R.Map[context.Context](func(ctx context.Context) string {
				return ctx.Value("result").(string)
			}),
		)

		ctx := context.Background()
		result := pipeline(ctx)

		assert.Equal(t, "step 1 complete", result)
	})
}

// TestWithValue_RealWorldScenario tests a realistic use case
func TestWithValue_RealWorldScenario(t *testing.T) {
	t.Run("HTTP request context enrichment", func(t *testing.T) {
		type RequestContext struct {
			UserID    string
			RequestID string
			TraceID   string
		}

		// Simulate enriching a context with request metadata
		enrichRequestContext := func(userID, requestID, traceID string) Reader[context.Context] {
			return func(ctx context.Context) context.Context {
				ctx = WithValue[string]("userID")(userID)(ctx)
				ctx = WithValue[string]("requestID")(requestID)(ctx)
				ctx = WithValue[string]("traceID")(traceID)(ctx)
				return ctx
			}
		}

		// Extract request context from enriched context
		getRequestContext := func(ctx context.Context) RequestContext {
			return RequestContext{
				UserID:    ctx.Value("userID").(string),
				RequestID: ctx.Value("requestID").(string),
				TraceID:   ctx.Value("traceID").(string),
			}
		}

		// Use the enrichment
		ctx := context.Background()
		enrichedCtx := enrichRequestContext("user123", "req456", "trace789")(ctx)
		reqCtx := getRequestContext(enrichedCtx)

		assert.Equal(t, "user123", reqCtx.UserID)
		assert.Equal(t, "req456", reqCtx.RequestID)
		assert.Equal(t, "trace789", reqCtx.TraceID)
	})
}

// ExampleNopCancel demonstrates wrapping a plain context in a no-op ContextCancel.
func ExampleNopCancel() {
	ctx := context.WithValue(context.Background(), "key", "value")
	cc := NopCancel(ctx)

	// The second element of the pair is the original context, unchanged.
	wrappedCtx := pair.Tail(cc)
	fmt.Println(wrappedCtx.Value("key"))

	// The first element is a no-op cancel function; calling it is safe.
	cancel := pair.Head(cc)
	cancel()

	// The context is still active after the no-op cancel.
	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
	default:
		fmt.Println("still active")
	}
	// Output:
	// value
	// still active
}

// TestNopCancel verifies the no-op cancellation semantics of NopCancel.
func TestNopCancel(t *testing.T) {
	t.Run("returns the same context unchanged", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", "value")
		cc := NopCancel(ctx)

		assert.Equal(t, ctx, pair.Tail(cc))
	})

	t.Run("cancel func is a no-op and does not cancel the context", func(t *testing.T) {
		ctx := context.Background()
		cc := NopCancel(ctx)

		cancel := pair.Head(cc)
		cancel() // must not panic and must not cancel ctx

		select {
		case <-ctx.Done():
			t.Fatal("context should not be cancelled after calling the no-op cancel func")
		default:
		}
	})

	t.Run("calling cancel multiple times does not panic", func(t *testing.T) {
		cc := NopCancel(context.Background())
		cancel := pair.Head(cc)

		assert.NotPanics(t, func() {
			cancel()
			cancel()
			cancel()
		})
	})
}
