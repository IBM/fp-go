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
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestChainConsumer_Success tests that ChainConsumer executes the consumer
// and returns Void when the computation succeeds
func TestChainConsumer_Success(t *testing.T) {
	// Track if consumer was called
	var consumed int
	consumer := func(x int) {
		consumed = x
	}

	// Create a successful computation and chain the consumer
	computation := F.Pipe1(
		Of[context.Context](42),
		ChainConsumer[context.Context](consumer),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer was called with correct value
	assert.Equal(t, 42, consumed)

	// Verify result is successful with Void
	assert.True(t, result.IsRight(res))
	if result.IsRight(res) {
		val := result.GetOrElse(func(error) Void { return Void{} })(res)
		assert.Equal(t, Void{}, val)
	}
}

// TestChainConsumer_Failure tests that ChainConsumer does not execute
// the consumer when the computation fails
func TestChainConsumer_Failure(t *testing.T) {
	// Track if consumer was called
	consumerCalled := false
	consumer := func(x int) {
		consumerCalled = true
	}

	// Create a failing computation
	expectedErr := errors.New("test error")
	computation := F.Pipe1(
		Left[context.Context, int](expectedErr),
		ChainConsumer[context.Context](consumer),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer was NOT called
	assert.False(t, consumerCalled)

	// Verify result is an error
	assert.True(t, result.IsLeft(res))
}

// TestChainConsumer_MultipleOperations tests chaining multiple operations
// with ChainConsumer in a pipeline
func TestChainConsumer_MultipleOperations(t *testing.T) {
	// Track consumer calls
	var values []int
	consumer := func(x int) {
		values = append(values, x)
	}

	// Create a pipeline with multiple operations
	computation := F.Pipe2(
		Of[context.Context](10),
		Map[context.Context](N.Mul(2)),
		ChainConsumer[context.Context](consumer),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer was called with transformed value
	assert.Equal(t, []int{20}, values)

	// Verify result is successful
	assert.True(t, result.IsRight(res))
}

// TestChainFirstConsumer_Success tests that ChainFirstConsumer executes
// the consumer and preserves the original value
func TestChainFirstConsumer_Success(t *testing.T) {
	// Track if consumer was called
	var consumed int
	consumer := func(x int) {
		consumed = x
	}

	// Create a successful computation and chain the consumer
	computation := F.Pipe1(
		Of[context.Context](42),
		ChainFirstConsumer[context.Context](consumer),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer was called with correct value
	assert.Equal(t, 42, consumed)

	// Verify result is successful and preserves original value
	assert.True(t, result.IsRight(res))
	if result.IsRight(res) {
		val := result.GetOrElse(func(error) int { return 0 })(res)
		assert.Equal(t, 42, val)
	}
}

// TestChainFirstConsumer_Failure tests that ChainFirstConsumer does not
// execute the consumer when the computation fails
func TestChainFirstConsumer_Failure(t *testing.T) {
	// Track if consumer was called
	consumerCalled := false
	consumer := func(x int) {
		consumerCalled = true
	}

	// Create a failing computation
	expectedErr := errors.New("test error")
	computation := F.Pipe1(
		Left[context.Context, int](expectedErr),
		ChainFirstConsumer[context.Context](consumer),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer was NOT called
	assert.False(t, consumerCalled)

	// Verify result is an error
	assert.True(t, result.IsLeft(res))
}

// TestChainFirstConsumer_PreservesValue tests that ChainFirstConsumer
// preserves the value for further processing
func TestChainFirstConsumer_PreservesValue(t *testing.T) {
	// Track consumer calls
	var logged []int
	logger := func(x int) {
		logged = append(logged, x)
	}

	// Create a pipeline that logs intermediate values
	computation := F.Pipe3(
		Of[context.Context](10),
		ChainFirstConsumer[context.Context](logger),
		Map[context.Context](N.Mul(2)),
		ChainFirstConsumer[context.Context](logger),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer was called at each step
	assert.Equal(t, []int{10, 20}, logged)

	// Verify final result
	assert.True(t, result.IsRight(res))
	if result.IsRight(res) {
		val := result.GetOrElse(func(error) int { return 0 })(res)
		assert.Equal(t, 20, val)
	}
}

// TestChainFirstConsumer_WithMap tests combining ChainFirstConsumer with Map
func TestChainFirstConsumer_WithMap(t *testing.T) {
	// Track intermediate values
	var intermediate int
	consumer := func(x int) {
		intermediate = x
	}

	// Create a pipeline with logging and transformation
	computation := F.Pipe2(
		Of[context.Context](5),
		ChainFirstConsumer[context.Context](consumer),
		Map[context.Context](N.Mul(3)),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer saw original value
	assert.Equal(t, 5, intermediate)

	// Verify final result is transformed
	assert.True(t, result.IsRight(res))
	if result.IsRight(res) {
		val := result.GetOrElse(func(error) int { return 0 })(res)
		assert.Equal(t, 15, val)
	}
}

// TestChainConsumer_WithContext tests that consumers work with context
func TestChainConsumer_WithContext(t *testing.T) {
	type Config struct {
		Multiplier int
	}

	// Track consumer calls
	var consumed int
	consumer := func(x int) {
		consumed = x
	}

	// Create a computation that uses context
	computation := F.Pipe2(
		Of[Config](10),
		Map[Config](N.Mul(2)),
		ChainConsumer[Config](consumer),
	)

	// Execute with context
	cfg := Config{Multiplier: 3}
	res := computation(cfg)()

	// Verify consumer was called
	assert.Equal(t, 20, consumed)

	// Verify result is successful
	assert.True(t, result.IsRight(res))
}

// TestChainFirstConsumer_SideEffects tests that ChainFirstConsumer
// can be used for side effects like logging
func TestChainFirstConsumer_SideEffects(t *testing.T) {
	// Simulate a logging side effect
	var logs []string
	logValue := func(x string) {
		logs = append(logs, "Processing: "+x)
	}

	// Create a pipeline with logging
	computation := F.Pipe3(
		Of[context.Context]("hello"),
		ChainFirstConsumer[context.Context](logValue),
		Map[context.Context](S.Append(" world")),
		ChainFirstConsumer[context.Context](logValue),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify logs were created
	assert.Equal(t, []string{
		"Processing: hello",
		"Processing: hello world",
	}, logs)

	// Verify final result
	assert.True(t, result.IsRight(res))
	if result.IsRight(res) {
		val := result.GetOrElse(func(error) string { return "" })(res)
		assert.Equal(t, "hello world", val)
	}
}

// TestChainConsumer_ComplexType tests consumers with complex types
func TestChainConsumer_ComplexType(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	// Track consumed user
	var consumedUser *User
	consumer := func(u User) {
		consumedUser = &u
	}

	// Create a computation with a complex type
	user := User{Name: "Alice", Age: 30}
	computation := F.Pipe1(
		Of[context.Context](user),
		ChainConsumer[context.Context](consumer),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer received the user
	assert.NotNil(t, consumedUser)
	assert.Equal(t, "Alice", consumedUser.Name)
	assert.Equal(t, 30, consumedUser.Age)

	// Verify result is successful
	assert.True(t, result.IsRight(res))
}

// TestChainFirstConsumer_ComplexType tests ChainFirstConsumer with complex types
func TestChainFirstConsumer_ComplexType(t *testing.T) {
	type Product struct {
		ID    int
		Name  string
		Price float64
	}

	// Track consumed products
	var consumedProducts []Product
	consumer := func(p Product) {
		consumedProducts = append(consumedProducts, p)
	}

	// Create a pipeline with complex type
	product := Product{ID: 1, Name: "Widget", Price: 9.99}
	computation := F.Pipe2(
		Of[context.Context](product),
		ChainFirstConsumer[context.Context](consumer),
		Map[context.Context](func(p Product) Product {
			p.Price = p.Price * 1.1 // Apply 10% markup
			return p
		}),
	)

	// Execute the computation
	res := computation(context.Background())()

	// Verify consumer saw original product
	assert.Len(t, consumedProducts, 1)
	assert.Equal(t, 9.99, consumedProducts[0].Price)

	// Verify final result has updated price
	assert.True(t, result.IsRight(res))
	if result.IsRight(res) {
		finalProduct := result.GetOrElse(func(error) Product { return Product{} })(res)
		assert.InDelta(t, 10.989, finalProduct.Price, 0.001)
	}
}
