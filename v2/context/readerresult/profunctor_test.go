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

package readerresult

import (
	"context"
	"strconv"
	"testing"

	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestPromapBasic tests basic Promap functionality
func TestPromapBasic(t *testing.T) {
	t.Run("transform both context and output", func(t *testing.T) {
		getValue := func(ctx context.Context) Result[int] {
			if v := ctx.Value("key"); v != nil {
				return R.Of(v.(int))
			}
			return R.Of(0)
		}

		addKey := func(ctx context.Context) (context.Context, context.CancelFunc) {
			newCtx := context.WithValue(ctx, "key", 42)
			return newCtx, func() {}
		}
		toString := strconv.Itoa

		adapted := Promap(addKey, toString)(getValue)
		result := adapted(t.Context())

		assert.Equal(t, R.Of("42"), result)
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("context transformation", func(t *testing.T) {
		getValue := func(ctx context.Context) Result[int] {
			if v := ctx.Value("key"); v != nil {
				return R.Of(v.(int))
			}
			return R.Of(0)
		}

		addKey := func(ctx context.Context) (context.Context, context.CancelFunc) {
			newCtx := context.WithValue(ctx, "key", 100)
			return newCtx, func() {}
		}

		adapted := Contramap[int](addKey)(getValue)
		result := adapted(t.Context())

		assert.Equal(t, R.Of(100), result)
	})
}

// TestLocalBasic tests basic Local functionality
func TestLocalBasic(t *testing.T) {
	t.Run("adds value to context", func(t *testing.T) {
		getValue := func(ctx context.Context) Result[string] {
			if v := ctx.Value("user"); v != nil {
				return R.Of(v.(string))
			}
			return R.Of("unknown")
		}

		addUser := func(ctx context.Context) (context.Context, context.CancelFunc) {
			newCtx := context.WithValue(ctx, "user", "Alice")
			return newCtx, func() {}
		}

		adapted := Local[string](addUser)(getValue)
		result := adapted(t.Context())

		assert.Equal(t, R.Of("Alice"), result)
	})
}
