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

package readeroption

import (
	"context"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestCurry0(t *testing.T) {
	// Function that returns a value from context
	getConfig := func(ctx context.Context) (string, bool) {
		if val := ctx.Value("config"); val != nil {
			return val.(string), true
		}
		return "", false
	}

	ro := Curry0(getConfig)

	// Test with value in context
	ctx1 := context.WithValue(context.Background(), "config", "test-config")
	result1 := ro(ctx1)
	assert.Equal(t, O.Of("test-config"), result1)

	// Test without value in context
	ctx2 := context.Background()
	result2 := ro(ctx2)
	assert.Equal(t, O.None[string](), result2)
}

func TestCurry1(t *testing.T) {
	// Function that looks up a value by key
	lookup := func(ctx context.Context, key string) (int, bool) {
		if val := ctx.Value(key); val != nil {
			return val.(int), true
		}
		return 0, false
	}

	ro := Curry1(lookup)

	// Test with value in context
	ctx1 := context.WithValue(context.Background(), "count", 42)
	result1 := ro("count")(ctx1)
	assert.Equal(t, O.Of(42), result1)

	// Test without value in context
	ctx2 := context.Background()
	result2 := ro("count")(ctx2)
	assert.Equal(t, O.None[int](), result2)
}

func TestCurry2(t *testing.T) {
	// Function that combines two parameters with context
	combine := func(ctx context.Context, a string, b int) (string, bool) {
		if ctx.Value("enabled") == true {
			return a + ":" + string(rune('0'+b)), true
		}
		return "", false
	}

	ro := Curry2(combine)

	// Test with enabled context
	ctx1 := context.WithValue(context.Background(), "enabled", true)
	result1 := ro("test")(5)(ctx1)
	assert.Equal(t, O.Of("test:5"), result1)

	// Test with disabled context
	ctx2 := context.Background()
	result2 := ro("test")(5)(ctx2)
	assert.Equal(t, O.None[string](), result2)
}

func TestCurry3(t *testing.T) {
	// Function that combines three parameters with context
	combine := func(ctx context.Context, a string, b int, c bool) (string, bool) {
		if ctx.Value("enabled") == true && c {
			return a + ":" + string(rune('0'+b)), true
		}
		return "", false
	}

	ro := Curry3(combine)

	// Test with enabled context and true flag
	ctx1 := context.WithValue(context.Background(), "enabled", true)
	result1 := ro("test")(5)(true)(ctx1)
	assert.Equal(t, O.Of("test:5"), result1)

	// Test with false flag
	result2 := ro("test")(5)(false)(ctx1)
	assert.Equal(t, O.None[string](), result2)
}

func TestUncurry1(t *testing.T) {
	// Create a curried function
	curried := func(x int) ReaderOption[context.Context, int] {
		return Of[context.Context](x * 2)
	}

	// Uncurry it
	uncurried := Uncurry1(curried)

	// Test the uncurried function
	result, ok := uncurried(context.Background(), 21)
	assert.True(t, ok)
	assert.Equal(t, 42, result)
}

func TestUncurry2(t *testing.T) {
	// Create a curried function
	curried := func(x int) func(y int) ReaderOption[context.Context, int] {
		return func(y int) ReaderOption[context.Context, int] {
			return Of[context.Context](x + y)
		}
	}

	// Uncurry it
	uncurried := Uncurry2(curried)

	// Test the uncurried function
	result, ok := uncurried(context.Background(), 10, 32)
	assert.True(t, ok)
	assert.Equal(t, 42, result)
}

func TestUncurry3(t *testing.T) {
	// Create a curried function
	curried := func(x int) func(y int) func(z int) ReaderOption[context.Context, int] {
		return func(y int) func(z int) ReaderOption[context.Context, int] {
			return func(z int) ReaderOption[context.Context, int] {
				return Of[context.Context](x + y + z)
			}
		}
	}

	// Uncurry it
	uncurried := Uncurry3(curried)

	// Test the uncurried function
	result, ok := uncurried(context.Background(), 10, 20, 12)
	assert.True(t, ok)
	assert.Equal(t, 42, result)
}
