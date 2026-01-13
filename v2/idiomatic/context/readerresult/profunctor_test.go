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
	"fmt"
	"strconv"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

// TestPromapBasic tests basic Promap functionality
func TestPromapBasic(t *testing.T) {
	t.Run("transform both context and output", func(t *testing.T) {
		// ReaderResult that reads a value from context
		getValue := func(ctx context.Context) (int, error) {
			if val := ctx.Value("port"); val != nil {
				return val.(int), nil
			}
			return 0, fmt.Errorf("port not found")
		}

		// Transform context to add a value and int to string
		addPort := func(ctx context.Context) (context.Context, context.CancelFunc) {
			return context.WithValue(ctx, "port", 8080), func() {}
		}
		toString := strconv.Itoa

		adapted := Promap(addPort, toString)(getValue)
		result, err := adapted(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "8080", result)
	})

	t.Run("handles error case", func(t *testing.T) {
		// ReaderResult that returns an error
		getError := func(ctx context.Context) (int, error) {
			return 0, fmt.Errorf("error occurred")
		}

		addPort := func(ctx context.Context) (context.Context, context.CancelFunc) {
			return context.WithValue(ctx, "port", 8080), func() {}
		}
		toString := strconv.Itoa

		adapted := Promap(addPort, toString)(getError)
		_, err := adapted(context.Background())

		assert.Error(t, err)
		assert.Equal(t, "error occurred", err.Error())
	})

	t.Run("context transformation with cancellation", func(t *testing.T) {
		getValue := func(ctx context.Context) (string, error) {
			if val := ctx.Value("key"); val != nil {
				return val.(string), nil
			}
			return "", fmt.Errorf("key not found")
		}

		addValue := func(ctx context.Context) (context.Context, context.CancelFunc) {
			ctx, cancel := context.WithCancel(ctx)
			return context.WithValue(ctx, "key", "value"), cancel
		}
		toUpper := func(s string) string {
			return "UPPER_" + s
		}

		adapted := Promap(addValue, toUpper)(getValue)
		result, err := adapted(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "UPPER_value", result)
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("context adaptation", func(t *testing.T) {
		// ReaderResult that reads from context
		getPort := func(ctx context.Context) (int, error) {
			if val := ctx.Value("port"); val != nil {
				return val.(int), nil
			}
			return 0, fmt.Errorf("port not found")
		}

		// Adapt context to add port value
		addPort := func(ctx context.Context) (context.Context, context.CancelFunc) {
			return context.WithValue(ctx, "port", 9000), func() {}
		}

		adapted := Contramap[int](addPort)(getPort)
		result, err := adapted(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 9000, result)
	})

	t.Run("preserves error", func(t *testing.T) {
		getError := func(ctx context.Context) (int, error) {
			return 0, fmt.Errorf("config error")
		}

		addPort := func(ctx context.Context) (context.Context, context.CancelFunc) {
			return context.WithValue(ctx, "port", 9000), func() {}
		}

		adapted := Contramap[int](addPort)(getError)
		_, err := adapted(context.Background())

		assert.Error(t, err)
		assert.Equal(t, "config error", err.Error())
	})

	t.Run("multiple context values", func(t *testing.T) {
		getValues := func(ctx context.Context) (string, error) {
			host := ctx.Value("host")
			port := ctx.Value("port")
			if host != nil && port != nil {
				return fmt.Sprintf("%s:%d", host, port), nil
			}
			return "", fmt.Errorf("missing values")
		}

		addValues := func(ctx context.Context) (context.Context, context.CancelFunc) {
			ctx = context.WithValue(ctx, "host", "localhost")
			ctx = context.WithValue(ctx, "port", 8080)
			return ctx, func() {}
		}

		adapted := Contramap[string](addValues)(getValues)
		result, err := adapted(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", result)
	})
}

// TestPromapComposition tests that Promap can be composed
func TestPromapComposition(t *testing.T) {
	t.Run("compose two Promap transformations", func(t *testing.T) {
		reader := func(ctx context.Context) (int, error) {
			if val := ctx.Value("value"); val != nil {
				return val.(int), nil
			}
			return 0, fmt.Errorf("value not found")
		}

		f1 := func(ctx context.Context) (context.Context, context.CancelFunc) {
			return context.WithValue(ctx, "value", 5), func() {}
		}
		g1 := N.Mul(2)

		f2 := func(ctx context.Context) (context.Context, context.CancelFunc) {
			return ctx, func() {}
		}
		g2 := N.Add(10)

		// Apply two Promap transformations
		step1 := Promap(f1, g1)(reader)
		step2 := Promap(f2, g2)(step1)

		result, err := step2(context.Background())

		// (5 * 2) + 10 = 20
		assert.NoError(t, err)
		assert.Equal(t, 20, result)
	})
}
