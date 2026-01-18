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

package readerresult

import (
	"context"
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// TestCurry0 tests the Curry0 function
func TestCurry0(t *testing.T) {
	t.Run("converts Go function to ReaderResult on success", func(t *testing.T) {
		// Idiomatic Go function
		getConfig := func(ctx context.Context) (int, error) {
			return 42, nil
		}

		// Convert to ReaderResult
		configRR := Curry0(getConfig)
		result := configRR(t.Context())

		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("converts Go function to ReaderResult on error", func(t *testing.T) {
		testErr := errors.New("config error")
		getConfig := func(ctx context.Context) (int, error) {
			return 0, testErr
		}

		configRR := Curry0(getConfig)
		result := configRR(t.Context())

		assert.Equal(t, E.Left[int](testErr), result)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		getConfig := func(ctx context.Context) (int, error) {
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}
			return 42, nil
		}

		configRR := Curry0(getConfig)

		// Test with cancelled context
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		result := configRR(ctx)

		assert.True(t, E.IsLeft(result))
	})

	t.Run("can be used in functional composition", func(t *testing.T) {
		getConfig := func(ctx context.Context) (int, error) {
			return 42, nil
		}

		pipeline := F.Pipe1(
			Curry0(getConfig),
			Map(func(x int) string {
				return "value"
			}),
		)

		result := pipeline(t.Context())
		assert.True(t, E.IsRight(result))
	})
}

// TestCurry1 tests the Curry1 function
func TestCurry1(t *testing.T) {
	t.Run("converts Go function to Kleisli on success", func(t *testing.T) {
		getUserByID := func(ctx context.Context, id int) (string, error) {
			return "Alice", nil
		}

		getUserKleisli := Curry1(getUserByID)

		// Use in a pipeline
		pipeline := F.Pipe1(
			Of(123),
			Chain(getUserKleisli),
		)

		result := pipeline(t.Context())
		assert.Equal(t, E.Of[error]("Alice"), result)
	})

	t.Run("converts Go function to Kleisli on error", func(t *testing.T) {
		testErr := errors.New("user not found")
		getUserByID := func(ctx context.Context, id int) (string, error) {
			return "", testErr
		}

		getUserKleisli := Curry1(getUserByID)
		result := getUserKleisli(123)(t.Context())

		assert.Equal(t, E.Left[string](testErr), result)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		getUserByID := func(ctx context.Context, id int) (string, error) {
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			return "Alice", nil
		}

		getUserKleisli := Curry1(getUserByID)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		result := getUserKleisli(123)(ctx)

		assert.True(t, E.IsLeft(result))
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		getUserByID := func(ctx context.Context, id int) (string, error) {
			return "Alice", nil
		}

		pipeline := F.Pipe2(
			Of(123),
			Chain(Curry1(getUserByID)),
			Map(func(name string) int {
				return len(name)
			}),
		)

		result := pipeline(t.Context())
		assert.Equal(t, E.Of[error](5), result) // len("Alice") = 5
	})
}

// TestCurry2 tests the Curry2 function
func TestCurry2(t *testing.T) {
	t.Run("converts Go function to curried form on success", func(t *testing.T) {
		updateUser := func(ctx context.Context, id int, name string) (string, error) {
			return name, nil
		}

		updateUserCurried := Curry2(updateUser)

		// Partial application
		updateUser123 := updateUserCurried(123)

		// Use in a pipeline
		pipeline := F.Pipe1(
			Of("Bob"),
			Chain(updateUser123),
		)

		result := pipeline(t.Context())
		assert.Equal(t, E.Of[error]("Bob"), result)
	})

	t.Run("converts Go function to curried form on error", func(t *testing.T) {
		testErr := errors.New("update failed")
		updateUser := func(ctx context.Context, id int, name string) (string, error) {
			return "", testErr
		}

		updateUserCurried := Curry2(updateUser)
		result := updateUserCurried(123)("Bob")(t.Context())

		assert.Equal(t, E.Left[string](testErr), result)
	})

	t.Run("supports partial application", func(t *testing.T) {
		concat := func(ctx context.Context, a string, b string) (string, error) {
			return a + b, nil
		}

		concatCurried := Curry2(concat)

		// Partial application
		prependHello := concatCurried("Hello, ")

		result1 := prependHello("World")(t.Context())
		result2 := prependHello("Alice")(t.Context())

		assert.Equal(t, E.Of[error]("Hello, World"), result1)
		assert.Equal(t, E.Of[error]("Hello, Alice"), result2)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		updateUser := func(ctx context.Context, id int, name string) (string, error) {
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			return name, nil
		}

		updateUserCurried := Curry2(updateUser)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		result := updateUserCurried(123)("Bob")(ctx)

		assert.True(t, E.IsLeft(result))
	})
}

// TestCurry3 tests the Curry3 function
func TestCurry3(t *testing.T) {
	t.Run("converts Go function to curried form on success", func(t *testing.T) {
		createOrder := func(ctx context.Context, userID int, productID int, quantity int) (int, error) {
			return userID + productID + quantity, nil
		}

		createOrderCurried := Curry3(createOrder)

		// Partial application
		createOrderForUser := createOrderCurried(100)
		createOrderForProduct := createOrderForUser(200)

		// Use in a pipeline
		pipeline := F.Pipe1(
			Of(3),
			Chain(createOrderForProduct),
		)

		result := pipeline(t.Context())
		assert.Equal(t, E.Of[error](303), result) // 100 + 200 + 3
	})

	t.Run("converts Go function to curried form on error", func(t *testing.T) {
		testErr := errors.New("order creation failed")
		createOrder := func(ctx context.Context, userID int, productID int, quantity int) (int, error) {
			return 0, testErr
		}

		createOrderCurried := Curry3(createOrder)
		result := createOrderCurried(100)(200)(3)(t.Context())

		assert.Equal(t, E.Left[int](testErr), result)
	})

	t.Run("supports multiple levels of partial application", func(t *testing.T) {
		sum3 := func(ctx context.Context, a int, b int, c int) (int, error) {
			return a + b + c, nil
		}

		sum3Curried := Curry3(sum3)

		// First level partial application
		add10 := sum3Curried(10)

		// Second level partial application
		add10And20 := add10(20)

		result1 := add10And20(5)(t.Context())
		result2 := add10And20(15)(t.Context())

		assert.Equal(t, E.Of[error](35), result1) // 10 + 20 + 5
		assert.Equal(t, E.Of[error](45), result2) // 10 + 20 + 15
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		createOrder := func(ctx context.Context, userID int, productID int, quantity int) (int, error) {
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}
			return userID + productID + quantity, nil
		}

		createOrderCurried := Curry3(createOrder)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		result := createOrderCurried(100)(200)(3)(ctx)

		assert.True(t, E.IsLeft(result))
	})
}

// TestUncurry1 tests the Uncurry1 function
func TestUncurry1(t *testing.T) {
	t.Run("converts Kleisli back to Go function on success", func(t *testing.T) {
		getUserKleisli := func(id int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				return E.Of[error]("Alice")
			}
		}

		getUserByID := Uncurry1(getUserKleisli)

		user, err := getUserByID(t.Context(), 123)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", user)
	})

	t.Run("converts Kleisli back to Go function on error", func(t *testing.T) {
		testErr := errors.New("user not found")
		getUserKleisli := func(id int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				return E.Left[string](testErr)
			}
		}

		getUserByID := Uncurry1(getUserKleisli)

		user, err := getUserByID(t.Context(), 123)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
		assert.Equal(t, "", user)
	})

	t.Run("respects context in uncurried function", func(t *testing.T) {
		getUserKleisli := func(id int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				if ctx.Err() != nil {
					return E.Left[string](ctx.Err())
				}
				return E.Of[error]("Alice")
			}
		}

		getUserByID := Uncurry1(getUserKleisli)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		user, err := getUserByID(ctx, 123)
		assert.Error(t, err)
		assert.Equal(t, "", user)
	})

	t.Run("round-trip with Curry1", func(t *testing.T) {
		// Original Go function
		original := func(ctx context.Context, id int) (string, error) {
			return "Alice", nil
		}

		// Curry then uncurry
		roundTrip := Uncurry1(Curry1(original))

		user, err := roundTrip(t.Context(), 123)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", user)
	})
}

// TestUncurry2 tests the Uncurry2 function
func TestUncurry2(t *testing.T) {
	t.Run("converts curried function back to Go function on success", func(t *testing.T) {
		updateUserCurried := func(id int) func(name string) ReaderResult[string] {
			return func(name string) ReaderResult[string] {
				return func(ctx context.Context) E.Either[error, string] {
					return E.Of[error](name)
				}
			}
		}

		updateUser := Uncurry2(updateUserCurried)

		result, err := updateUser(t.Context(), 123, "Bob")
		assert.NoError(t, err)
		assert.Equal(t, "Bob", result)
	})

	t.Run("converts curried function back to Go function on error", func(t *testing.T) {
		testErr := errors.New("update failed")
		updateUserCurried := func(id int) func(name string) ReaderResult[string] {
			return func(name string) ReaderResult[string] {
				return func(ctx context.Context) E.Either[error, string] {
					return E.Left[string](testErr)
				}
			}
		}

		updateUser := Uncurry2(updateUserCurried)

		result, err := updateUser(t.Context(), 123, "Bob")
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
		assert.Equal(t, "", result)
	})

	t.Run("respects context in uncurried function", func(t *testing.T) {
		updateUserCurried := func(id int) func(name string) ReaderResult[string] {
			return func(name string) ReaderResult[string] {
				return func(ctx context.Context) E.Either[error, string] {
					if ctx.Err() != nil {
						return E.Left[string](ctx.Err())
					}
					return E.Of[error](name)
				}
			}
		}

		updateUser := Uncurry2(updateUserCurried)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result, err := updateUser(ctx, 123, "Bob")
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("round-trip with Curry2", func(t *testing.T) {
		// Original Go function
		original := func(ctx context.Context, a string, b string) (string, error) {
			return a + b, nil
		}

		// Curry then uncurry
		roundTrip := Uncurry2(Curry2(original))

		result, err := roundTrip(t.Context(), "Hello, ", "World")
		assert.NoError(t, err)
		assert.Equal(t, "Hello, World", result)
	})
}

// TestUncurry3 tests the Uncurry3 function
func TestUncurry3(t *testing.T) {
	t.Run("converts curried function back to Go function on success", func(t *testing.T) {
		createOrderCurried := func(userID int) func(productID int) func(quantity int) ReaderResult[int] {
			return func(productID int) func(quantity int) ReaderResult[int] {
				return func(quantity int) ReaderResult[int] {
					return func(ctx context.Context) E.Either[error, int] {
						return E.Of[error](userID + productID + quantity)
					}
				}
			}
		}

		createOrder := Uncurry3(createOrderCurried)

		result, err := createOrder(t.Context(), 100, 200, 3)
		assert.NoError(t, err)
		assert.Equal(t, 303, result) // 100 + 200 + 3
	})

	t.Run("converts curried function back to Go function on error", func(t *testing.T) {
		testErr := errors.New("order creation failed")
		createOrderCurried := func(userID int) func(productID int) func(quantity int) ReaderResult[int] {
			return func(productID int) func(quantity int) ReaderResult[int] {
				return func(quantity int) ReaderResult[int] {
					return func(ctx context.Context) E.Either[error, int] {
						return E.Left[int](testErr)
					}
				}
			}
		}

		createOrder := Uncurry3(createOrderCurried)

		result, err := createOrder(t.Context(), 100, 200, 3)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
		assert.Equal(t, 0, result)
	})

	t.Run("respects context in uncurried function", func(t *testing.T) {
		createOrderCurried := func(userID int) func(productID int) func(quantity int) ReaderResult[int] {
			return func(productID int) func(quantity int) ReaderResult[int] {
				return func(quantity int) ReaderResult[int] {
					return func(ctx context.Context) E.Either[error, int] {
						if ctx.Err() != nil {
							return E.Left[int](ctx.Err())
						}
						return E.Of[error](userID + productID + quantity)
					}
				}
			}
		}

		createOrder := Uncurry3(createOrderCurried)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result, err := createOrder(ctx, 100, 200, 3)
		assert.Error(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("round-trip with Curry3", func(t *testing.T) {
		// Original Go function
		original := func(ctx context.Context, a int, b int, c int) (int, error) {
			return a + b + c, nil
		}

		// Curry then uncurry
		roundTrip := Uncurry3(Curry3(original))

		result, err := roundTrip(t.Context(), 10, 20, 5)
		assert.NoError(t, err)
		assert.Equal(t, 35, result) // 10 + 20 + 5
	})
}

// TestCurryUncurryIntegration tests integration between curry and uncurry functions
func TestCurryUncurryIntegration(t *testing.T) {
	t.Run("Curry1 and Uncurry1 are inverses", func(t *testing.T) {
		original := func(ctx context.Context, x int) (int, error) {
			return x * 2, nil
		}

		// Curry then uncurry should give back equivalent function
		roundTrip := Uncurry1(Curry1(original))

		result1, err1 := original(t.Context(), 21)
		result2, err2 := roundTrip(t.Context(), 21)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, result1, result2)
	})

	t.Run("Curry2 and Uncurry2 are inverses", func(t *testing.T) {
		original := func(ctx context.Context, x int, y int) (int, error) {
			return x + y, nil
		}

		roundTrip := Uncurry2(Curry2(original))

		result1, err1 := original(t.Context(), 10, 20)
		result2, err2 := roundTrip(t.Context(), 10, 20)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, result1, result2)
	})

	t.Run("Curry3 and Uncurry3 are inverses", func(t *testing.T) {
		original := func(ctx context.Context, x int, y int, z int) (int, error) {
			return x * y * z, nil
		}

		roundTrip := Uncurry3(Curry3(original))

		result1, err1 := original(t.Context(), 2, 3, 4)
		result2, err2 := roundTrip(t.Context(), 2, 3, 4)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, result1, result2)
	})
}
