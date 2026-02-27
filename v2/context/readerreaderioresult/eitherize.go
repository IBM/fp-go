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

package readerreaderioresult

import (
	"context"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
)

// Eitherize converts a function that returns a value and error into a ReaderReaderIOResult.
//
// This function takes a function that accepts an outer context R and context.Context,
// returning a value T and an error, and converts it into a ReaderReaderIOResult[R, T].
// The error is automatically converted into the Left case of the Result, while successful
// values become the Right case.
//
// This is particularly useful for integrating standard Go error-handling patterns into
// the functional programming style of ReaderReaderIOResult. It is especially helpful
// for adapting interface member functions that accept a context. When you have an
// interface method with signature (receiver, context.Context) (T, error), you can
// use Eitherize to convert it into a ReaderReaderIOResult where the receiver becomes
// the outer reader context R.
//
// # Type Parameters
//
//   - R: The outer reader context type (e.g., application configuration)
//   - T: The success value type
//
// # Parameters
//
//   - f: A function that takes R and context.Context and returns (T, error)
//
// # Returns
//
//   - ReaderReaderIOResult[R, T]: A computation that depends on R and context.Context,
//     performs IO, and produces a Result[T]
//
// # Example Usage
//
//	type AppConfig struct {
//	    DatabaseURL string
//	}
//
//	// A function using standard Go error handling
//	func fetchUser(cfg AppConfig, ctx context.Context) (*User, error) {
//	    // Implementation that may return an error
//	    return &User{ID: 1, Name: "Alice"}, nil
//	}
//
//	// Convert to ReaderReaderIOResult
//	fetchUserRR := Eitherize(fetchUser)
//
//	// Use in functional composition
//	result := F.Pipe1(
//	    fetchUserRR,
//	    Map[AppConfig](func(u *User) string { return u.Name }),
//	)
//
//	// Execute with config and context
//	cfg := AppConfig{DatabaseURL: "postgres://localhost"}
//	outcome := result(cfg)(context.Background())()
//
// # Adapting Interface Methods
//
// Eitherize is particularly useful for adapting interface member functions:
//
//	type UserRepository interface {
//	    GetUser(ctx context.Context, id int) (*User, error)
//	}
//
//	type UserRepo struct {
//	    db *sql.DB
//	}
//
//	func (r *UserRepo) GetUser(ctx context.Context, id int) (*User, error) {
//	    // Implementation
//	    return &User{ID: id}, nil
//	}
//
//	// Adapt the method by binding the first parameter (receiver)
//	repo := &UserRepo{db: db}
//	getUserRR := Eitherize(func(id int, ctx context.Context) (*User, error) {
//	    return repo.GetUser(ctx, id)
//	})
//
//	// Now getUserRR has type: ReaderReaderIOResult[int, *User]
//	// The receiver (repo) is captured in the closure
//	// The id becomes the outer reader context R
//
// # See Also
//
//   - Eitherize1: For functions that take an additional parameter
//   - ioresult.Eitherize2: The underlying conversion function
func Eitherize[R, T any](f func(R, context.Context) (T, error)) ReaderReaderIOResult[R, T] {
	return F.Pipe1(
		ioresult.Eitherize2(f),
		F.Curry2,
	)
}

// Eitherize1 converts a function that takes an additional parameter and returns a value
// and error into a Kleisli arrow.
//
// This function takes a function that accepts an outer context R, context.Context, and
// an additional parameter A, returning a value T and an error, and converts it into a
// Kleisli arrow (A -> ReaderReaderIOResult[R, T]). The error is automatically converted
// into the Left case of the Result, while successful values become the Right case.
//
// This is useful for creating composable operations that depend on both contexts and
// an input value, following standard Go error-handling patterns. It is especially helpful
// for adapting interface member functions that accept a context and additional parameters.
// When you have an interface method with signature (receiver, context.Context, A) (T, error),
// you can use Eitherize1 to convert it into a Kleisli arrow where the receiver becomes
// the outer reader context R and A becomes the input parameter.
//
// # Type Parameters
//
//   - R: The outer reader context type (e.g., application configuration)
//   - A: The input parameter type
//   - T: The success value type
//
// # Parameters
//
//   - f: A function that takes R, context.Context, and A, returning (T, error)
//
// # Returns
//
//   - Kleisli[R, A, T]: A function from A to ReaderReaderIOResult[R, T]
//
// # Example Usage
//
//	type AppConfig struct {
//	    DatabaseURL string
//	}
//
//	// A function using standard Go error handling
//	func fetchUserByID(cfg AppConfig, ctx context.Context, id int) (*User, error) {
//	    // Implementation that may return an error
//	    return &User{ID: id, Name: "Alice"}, nil
//	}
//
//	// Convert to Kleisli arrow
//	fetchUserKleisli := Eitherize1(fetchUserByID)
//
//	// Use in functional composition with Chain
//	pipeline := F.Pipe1(
//	    Of[AppConfig](123),
//	    Chain[AppConfig](fetchUserKleisli),
//	)
//
//	// Execute with config and context
//	cfg := AppConfig{DatabaseURL: "postgres://localhost"}
//	outcome := pipeline(cfg)(context.Background())()
//
// # Adapting Interface Methods
//
// Eitherize1 is particularly useful for adapting interface member functions with parameters:
//
//	type UserRepository interface {
//	    GetUserByID(ctx context.Context, id int) (*User, error)
//	    UpdateUser(ctx context.Context, user *User) error
//	}
//
//	type UserRepo struct {
//	    db *sql.DB
//	}
//
//	func (r *UserRepo) GetUserByID(ctx context.Context, id int) (*User, error) {
//	    // Implementation
//	    return &User{ID: id}, nil
//	}
//
//	// Adapt the method - receiver becomes R, id becomes A
//	repo := &UserRepo{db: db}
//	getUserKleisli := Eitherize1(func(r *UserRepo, ctx context.Context, id int) (*User, error) {
//	    return r.GetUserByID(ctx, id)
//	})
//
//	// Now getUserKleisli has type: Kleisli[*UserRepo, int, *User]
//	// Which is: func(int) ReaderReaderIOResult[*UserRepo, *User]
//	// Use it in composition:
//	pipeline := F.Pipe1(
//	    Of[*UserRepo](123),
//	    Chain[*UserRepo](getUserKleisli),
//	)
//	result := pipeline(repo)(context.Background())()
//
// # See Also
//
//   - Eitherize: For functions without an additional parameter
//   - Chain: For composing Kleisli arrows
//   - ioresult.Eitherize3: The underlying conversion function
func Eitherize1[R, A, T any](f func(R, context.Context, A) (T, error)) Kleisli[R, A, T] {
	return F.Flow2(
		F.Bind3of3(ioresult.Eitherize3(f)),
		F.Curry2,
	)
}
