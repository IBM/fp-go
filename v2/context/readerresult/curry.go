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

	"github.com/IBM/fp-go/v2/readereither"
)

// Curry and Uncurry functions convert between idiomatic Go functions (with context.Context as the first parameter)
// and functional ReaderResult/Kleisli compositions (with context.Context as the last parameter).
//
// This follows the Go convention from https://pkg.go.dev/context to put context as the first parameter,
// while enabling functional composition where context is typically the last parameter.
//
// The curry functions transform:
//   func(context.Context, T1, T2, ...) (A, error)  →  func(T1) func(T2) ... ReaderResult[A]
//
// The uncurry functions transform:
//   func(T1) func(T2) ... ReaderResult[A]  →  func(context.Context, T1, T2, ...) (A, error)

// Curry0 converts a Go function with context and no additional parameters into a ReaderResult.
// This is useful for adapting context-aware functions to the ReaderResult monad.
//
// Type Parameters:
//   - A: The return type of the function
//
// Parameters:
//   - f: A function that takes a context and returns a value and error
//
// Returns:
//   - A ReaderResult that wraps the function
//
// Example:
//
//	// Idiomatic Go function
//	getConfig := func(ctx context.Context) (Config, error) {
//	    // Check context cancellation
//	    if ctx.Err() != nil {
//	        return Config{}, ctx.Err()
//	    }
//	    return Config{Value: 42}, nil
//	}
//
//	// Convert to ReaderResult for functional composition
//	configRR := readerresult.Curry0(getConfig)
//	result := configRR(t.Context()) // Right(Config{Value: 42})
//
//go:inline
func Curry0[A any](f func(context.Context) (A, error)) ReaderResult[A] {
	return readereither.Curry0(f)
}

// Curry1 converts a Go function with context and one parameter into a Kleisli arrow.
// This enables functional composition of single-parameter functions.
//
// Type Parameters:
//   - T1: The type of the first parameter
//   - A: The return type of the function
//
// Parameters:
//   - f: A function that takes a context and one parameter, returning a value and error
//
// Returns:
//   - A Kleisli arrow that can be composed with other ReaderResult operations
//
// Example:
//
//	// Idiomatic Go function
//	getUserByID := func(ctx context.Context, id int) (User, error) {
//	    if ctx.Err() != nil {
//	        return User{}, ctx.Err()
//	    }
//	    return User{ID: id, Name: "Alice"}, nil
//	}
//
//	// Convert to Kleisli for functional composition
//	getUserKleisli := readerresult.Curry1(getUserByID)
//
//	// Use in a pipeline
//	pipeline := F.Pipe1(
//	    readerresult.Of(123),
//	    readerresult.Chain(getUserKleisli),
//	)
//	result := pipeline(t.Context()) // Right(User{ID: 123, Name: "Alice"})
//
//go:inline
func Curry1[T1, A any](f func(context.Context, T1) (A, error)) Kleisli[T1, A] {
	return readereither.Curry1(f)
}

// Curry2 converts a Go function with context and two parameters into a curried function.
// This enables partial application and functional composition of two-parameter functions.
//
// Type Parameters:
//   - T1: The type of the first parameter
//   - T2: The type of the second parameter
//   - A: The return type of the function
//
// Parameters:
//   - f: A function that takes a context and two parameters, returning a value and error
//
// Returns:
//   - A curried function that takes T1 and returns a Kleisli arrow for T2
//
// Example:
//
//	// Idiomatic Go function
//	updateUser := func(ctx context.Context, id int, name string) (User, error) {
//	    if ctx.Err() != nil {
//	        return User{}, ctx.Err()
//	    }
//	    return User{ID: id, Name: name}, nil
//	}
//
//	// Convert to curried form
//	updateUserCurried := readerresult.Curry2(updateUser)
//
//	// Partial application
//	updateUser123 := updateUserCurried(123)
//
//	// Use in a pipeline
//	pipeline := F.Pipe1(
//	    readerresult.Of("Bob"),
//	    readerresult.Chain(updateUser123),
//	)
//	result := pipeline(t.Context()) // Right(User{ID: 123, Name: "Bob"})
//
//go:inline
func Curry2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1) Kleisli[T2, A] {
	return readereither.Curry2(f)
}

// Curry3 converts a Go function with context and three parameters into a curried function.
// This enables partial application and functional composition of three-parameter functions.
//
// Type Parameters:
//   - T1: The type of the first parameter
//   - T2: The type of the second parameter
//   - T3: The type of the third parameter
//   - A: The return type of the function
//
// Parameters:
//   - f: A function that takes a context and three parameters, returning a value and error
//
// Returns:
//   - A curried function that takes T1, T2, and returns a Kleisli arrow for T3
//
// Example:
//
//	// Idiomatic Go function
//	createOrder := func(ctx context.Context, userID int, productID int, quantity int) (Order, error) {
//	    if ctx.Err() != nil {
//	        return Order{}, ctx.Err()
//	    }
//	    return Order{UserID: userID, ProductID: productID, Quantity: quantity}, nil
//	}
//
//	// Convert to curried form
//	createOrderCurried := readerresult.Curry3(createOrder)
//
//	// Partial application
//	createOrderForUser := createOrderCurried(123)
//	createOrderForProduct := createOrderForUser(456)
//
//	// Use in a pipeline
//	pipeline := F.Pipe1(
//	    readerresult.Of(2),
//	    readerresult.Chain(createOrderForProduct),
//	)
//	result := pipeline(t.Context()) // Right(Order{UserID: 123, ProductID: 456, Quantity: 2})
//
//go:inline
func Curry3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1) func(T2) Kleisli[T3, A] {
	return readereither.Curry3(f)
}

// Uncurry1 converts a Kleisli arrow back into an idiomatic Go function with context as the first parameter.
// This is useful for interfacing with code that expects standard Go function signatures.
//
// Type Parameters:
//   - T1: The type of the parameter
//   - A: The return type
//
// Parameters:
//   - f: A Kleisli arrow
//
// Returns:
//   - A Go function with context as the first parameter
//
// Example:
//
//	// Kleisli arrow
//	getUserKleisli := func(id int) readerresult.ReaderResult[User] {
//	    return func(ctx context.Context) result.Result[User] {
//	        if ctx.Err() != nil {
//	            return result.Error[User](ctx.Err())
//	        }
//	        return result.Of(User{ID: id, Name: "Alice"})
//	    }
//	}
//
//	// Convert back to idiomatic Go function
//	getUserByID := readerresult.Uncurry1(getUserKleisli)
//
//	// Use as a normal Go function
//	user, err := getUserByID(t.Context(), 123)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(user.Name) // "Alice"
//
//go:inline
func Uncurry1[T1, A any](f Kleisli[T1, A]) func(context.Context, T1) (A, error) {
	return readereither.Uncurry1(f)
}

// Uncurry2 converts a curried function back into an idiomatic Go function with context as the first parameter.
// This is useful for interfacing with code that expects standard Go function signatures.
//
// Type Parameters:
//   - T1: The type of the first parameter
//   - T2: The type of the second parameter
//   - A: The return type
//
// Parameters:
//   - f: A curried function
//
// Returns:
//   - A Go function with context as the first parameter
//
// Example:
//
//	// Curried function
//	updateUserCurried := func(id int) func(name string) readerresult.ReaderResult[User] {
//	    return func(name string) readerresult.ReaderResult[User] {
//	        return func(ctx context.Context) result.Result[User] {
//	            if ctx.Err() != nil {
//	                return result.Error[User](ctx.Err())
//	            }
//	            return result.Of(User{ID: id, Name: name})
//	        }
//	    }
//	}
//
//	// Convert back to idiomatic Go function
//	updateUser := readerresult.Uncurry2(updateUserCurried)
//
//	// Use as a normal Go function
//	user, err := updateUser(t.Context(), 123, "Bob")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(user.Name) // "Bob"
//
//go:inline
func Uncurry2[T1, T2, A any](f func(T1) Kleisli[T2, A]) func(context.Context, T1, T2) (A, error) {
	return readereither.Uncurry2(f)
}

// Uncurry3 converts a curried function back into an idiomatic Go function with context as the first parameter.
// This is useful for interfacing with code that expects standard Go function signatures.
//
// Type Parameters:
//   - T1: The type of the first parameter
//   - T2: The type of the second parameter
//   - T3: The type of the third parameter
//   - A: The return type
//
// Parameters:
//   - f: A curried function
//
// Returns:
//   - A Go function with context as the first parameter
//
// Example:
//
//	// Curried function
//	createOrderCurried := func(userID int) func(productID int) func(quantity int) readerresult.ReaderResult[Order] {
//	    return func(productID int) func(quantity int) readerresult.ReaderResult[Order] {
//	        return func(quantity int) readerresult.ReaderResult[Order] {
//	            return func(ctx context.Context) result.Result[Order] {
//	                if ctx.Err() != nil {
//	                    return result.Error[Order](ctx.Err())
//	                }
//	                return result.Of(Order{UserID: userID, ProductID: productID, Quantity: quantity})
//	            }
//	        }
//	    }
//	}
//
//	// Convert back to idiomatic Go function
//	createOrder := readerresult.Uncurry3(createOrderCurried)
//
//	// Use as a normal Go function
//	order, err := createOrder(t.Context(), 123, 456, 2)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Order: User=%d, Product=%d, Qty=%d\n", order.UserID, order.ProductID, order.Quantity)
//
//go:inline
func Uncurry3[T1, T2, T3, A any](f func(T1) func(T2) Kleisli[T3, A]) func(context.Context, T1, T2, T3) (A, error) {
	return readereither.Uncurry3(f)
}
