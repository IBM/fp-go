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

	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
)

// Curry0 converts a function that takes context.Context and returns (A, error) into a ReaderResult[A].
//
// This is useful for lifting existing functions that follow Go's context-first convention
// into the ReaderResult monad.
//
// Type Parameters:
//   - A: The return value type
//
// Parameters:
//   - f: A function that takes context.Context and returns (A, error)
//
// Returns:
//   - A ReaderResult[A] that wraps the function
//
// Example:
//
//	func getConfig(ctx context.Context) (Config, error) {
//	    // ... implementation
//	    return config, nil
//	}
//	rr := readerresult.Curry0(getConfig)
//	config, err := rr(ctx)
//
//go:inline
func Curry0[A any](f func(context.Context) (A, error)) ReaderResult[A] {
	return RR.Curry0(f)
}

// Curry1 converts a function with one parameter into a curried ReaderResult-returning function.
//
// The context.Context parameter is handled by the ReaderResult, allowing you to partially
// apply the business parameter before providing the context.
//
// Type Parameters:
//   - T1: The first parameter type
//   - A: The return value type
//
// Parameters:
//   - f: A function that takes (context.Context, T1) and returns (A, error)
//
// Returns:
//   - A curried function that takes T1 and returns ReaderResult[A]
//
// Example:
//
//	func getUser(ctx context.Context, id int) (User, error) {
//	    // ... implementation
//	    return user, nil
//	}
//	getUserRR := readerresult.Curry1(getUser)
//	rr := getUserRR(42)  // Partially applied
//	user, err := rr(ctx)  // Execute with context
//
//go:inline
func Curry1[T1, A any](f func(context.Context, T1) (A, error)) func(T1) ReaderResult[A] {
	return RR.Curry1(f)
}

// Curry2 converts a function with two parameters into a curried ReaderResult-returning function.
//
// The context.Context parameter is handled by the ReaderResult, allowing you to partially
// apply the business parameters before providing the context.
//
// Type Parameters:
//   - T1: The first parameter type
//   - T2: The second parameter type
//   - A: The return value type
//
// Parameters:
//   - f: A function that takes (context.Context, T1, T2) and returns (A, error)
//
// Returns:
//   - A curried function that takes T1, then T2, and returns ReaderResult[A]
//
// Example:
//
//	func updateUser(ctx context.Context, id int, name string) (User, error) {
//	    // ... implementation
//	    return user, nil
//	}
//	updateUserRR := readerresult.Curry2(updateUser)
//	rr := updateUserRR(42)("Alice")  // Partially applied
//	user, err := rr(ctx)  // Execute with context
//
//go:inline
func Curry2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1) func(T2) ReaderResult[A] {
	return RR.Curry2(f)
}

// Curry3 converts a function with three parameters into a curried ReaderResult-returning function.
//
// The context.Context parameter is handled by the ReaderResult, allowing you to partially
// apply the business parameters before providing the context.
//
// Type Parameters:
//   - T1: The first parameter type
//   - T2: The second parameter type
//   - T3: The third parameter type
//   - A: The return value type
//
// Parameters:
//   - f: A function that takes (context.Context, T1, T2, T3) and returns (A, error)
//
// Returns:
//   - A curried function that takes T1, then T2, then T3, and returns ReaderResult[A]
//
// Example:
//
//	func createPost(ctx context.Context, userID int, title string, body string) (Post, error) {
//	    // ... implementation
//	    return post, nil
//	}
//	createPostRR := readerresult.Curry3(createPost)
//	rr := createPostRR(42)("Title")("Body")  // Partially applied
//	post, err := rr(ctx)  // Execute with context
//
//go:inline
func Curry3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderResult[A] {
	return RR.Curry3(f)
}

// Uncurry1 converts a curried ReaderResult function back to a standard Go function.
//
// This is the inverse of Curry1, useful when you need to call curried functions
// in a traditional Go style.
//
// Type Parameters:
//   - T1: The parameter type
//   - A: The return value type
//
// Parameters:
//   - f: A curried function that takes T1 and returns ReaderResult[A]
//
// Returns:
//   - A function that takes (context.Context, T1) and returns (A, error)
//
// Example:
//
//	curriedFn := func(id int) readerresult.ReaderResult[User] { ... }
//	normalFn := readerresult.Uncurry1(curriedFn)
//	user, err := normalFn(ctx, 42)
//
//go:inline
func Uncurry1[T1, A any](f func(T1) ReaderResult[A]) func(context.Context, T1) (A, error) {
	return RR.Uncurry1(f)
}

// Uncurry2 converts a curried ReaderResult function with two parameters back to a standard Go function.
//
// This is the inverse of Curry2.
//
// Type Parameters:
//   - T1: The first parameter type
//   - T2: The second parameter type
//   - A: The return value type
//
// Parameters:
//   - f: A curried function that takes T1, then T2, and returns ReaderResult[A]
//
// Returns:
//   - A function that takes (context.Context, T1, T2) and returns (A, error)
//
//go:inline
func Uncurry2[T1, T2, A any](f func(T1) func(T2) ReaderResult[A]) func(context.Context, T1, T2) (A, error) {
	return RR.Uncurry2(f)
}

// Uncurry3 converts a curried ReaderResult function with three parameters back to a standard Go function.
//
// This is the inverse of Curry3.
//
// Type Parameters:
//   - T1: The first parameter type
//   - T2: The second parameter type
//   - T3: The third parameter type
//   - A: The return value type
//
// Parameters:
//   - f: A curried function that takes T1, then T2, then T3, and returns ReaderResult[A]
//
// Returns:
//   - A function that takes (context.Context, T1, T2, T3) and returns (A, error)
//
//go:inline
func Uncurry3[T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderResult[A]) func(context.Context, T1, T2, T3) (A, error) {
	return RR.Uncurry3(f)
}
