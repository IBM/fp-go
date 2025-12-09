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

// From0 converts a context-taking function into a thunk that returns a ReaderResult.
//
// Unlike Curry0 which returns a ReaderResult directly, From0 returns a function
// that when called produces a ReaderResult. This is useful for lazy evaluation.
//
// Type Parameters:
//   - A: The return value type
//
// Parameters:
//   - f: A function that takes context.Context and returns (A, error)
//
// Returns:
//   - A thunk (function with no parameters) that returns ReaderResult[A]
//
// Example:
//
//	func getConfig(ctx context.Context) (Config, error) {
//	    return Config{Port: 8080}, nil
//	}
//	thunk := readerresult.From0(getConfig)
//	rr := thunk()  // Create the ReaderResult
//	config, err := rr(ctx)  // Execute it
//
//go:inline
func From0[A any](f func(context.Context) (A, error)) func() ReaderResult[A] {
	return RR.From0(f)
}

// From1 converts a function with one parameter into an uncurried ReaderResult-returning function.
//
// Unlike Curry1 which returns a curried function, From1 returns a function that takes
// all parameters at once (except context). This is more convenient for direct calls.
//
// Type Parameters:
//   - T1: The parameter type
//   - A: The return value type
//
// Parameters:
//   - f: A function that takes (context.Context, T1) and returns (A, error)
//
// Returns:
//   - A function that takes T1 and returns ReaderResult[A]
//
// Example:
//
//	func getUser(ctx context.Context, id int) (User, error) {
//	    return User{ID: id}, nil
//	}
//	getUserRR := readerresult.From1(getUser)
//	rr := getUserRR(42)
//	user, err := rr(ctx)
//
//go:inline
func From1[T1, A any](f func(context.Context, T1) (A, error)) func(T1) ReaderResult[A] {
	return RR.From1(f)
}

// From2 converts a function with two parameters into an uncurried ReaderResult-returning function.
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
//   - A function that takes (T1, T2) and returns ReaderResult[A]
//
// Example:
//
//	func updateUser(ctx context.Context, id int, name string) (User, error) {
//	    return User{ID: id, Name: name}, nil
//	}
//	updateUserRR := readerresult.From2(updateUser)
//	rr := updateUserRR(42, "Alice")
//	user, err := rr(ctx)
//
//go:inline
func From2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1, T2) ReaderResult[A] {
	return RR.From2(f)
}

// From3 converts a function with three parameters into an uncurried ReaderResult-returning function.
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
//   - A function that takes (T1, T2, T3) and returns ReaderResult[A]
//
// Example:
//
//	func createPost(ctx context.Context, userID int, title, body string) (Post, error) {
//	    return Post{UserID: userID, Title: title, Body: body}, nil
//	}
//	createPostRR := readerresult.From3(createPost)
//	rr := createPostRR(42, "Title", "Body")
//	post, err := rr(ctx)
//
//go:inline
func From3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderResult[A] {
	return RR.From3(f)
}
