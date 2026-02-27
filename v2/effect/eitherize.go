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

package effect

import (
	"context"

	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
)

// Eitherize converts a function that returns a value and error into an Effect.
//
// This function takes a function that accepts a context C and context.Context,
// returning a value T and an error, and converts it into an Effect[C, T].
// The error is automatically converted into a failure, while successful
// values become successes.
//
// This is particularly useful for integrating standard Go error-handling patterns into
// the effect system. It is especially helpful for adapting interface member functions
// that accept a context. When you have an interface method with signature
// (receiver, context.Context) (T, error), you can use Eitherize to convert it into
// an Effect where the receiver becomes the context C.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - T: The success value type
//
// # Parameters
//
//   - f: A function that takes C and context.Context and returns (T, error)
//
// # Returns
//
//   - Effect[C, T]: An effect that depends on C, performs IO, and produces T
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
//	// Convert to Effect
//	fetchUserEffect := effect.Eitherize(fetchUser)
//
//	// Use in functional composition
//	pipeline := F.Pipe1(
//	    fetchUserEffect,
//	    effect.Map[AppConfig](func(u *User) string { return u.Name }),
//	)
//
//	// Execute with config
//	cfg := AppConfig{DatabaseURL: "postgres://localhost"}
//	result, err := effect.RunSync(effect.Provide[*User](cfg)(pipeline))(context.Background())
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
//	getUserEffect := effect.Eitherize(func(id int, ctx context.Context) (*User, error) {
//	    return repo.GetUser(ctx, id)
//	})
//
//	// Now getUserEffect has type: Effect[int, *User]
//	// The receiver (repo) is captured in the closure
//	// The id becomes the context C
//
// # See Also
//
//   - Eitherize1: For functions that take an additional parameter
//   - readerreaderioresult.Eitherize: The underlying implementation
//
//go:inline
func Eitherize[C, T any](f func(C, context.Context) (T, error)) Effect[C, T] {
	return readerreaderioresult.Eitherize(f)
}

// Eitherize1 converts a function that takes an additional parameter and returns a value
// and error into a Kleisli arrow.
//
// This function takes a function that accepts a context C, context.Context, and
// an additional parameter A, returning a value T and an error, and converts it into a
// Kleisli arrow (A -> Effect[C, T]). The error is automatically converted into a failure,
// while successful values become successes.
//
// This is useful for creating composable operations that depend on context and
// an input value, following standard Go error-handling patterns. It is especially helpful
// for adapting interface member functions that accept a context and additional parameters.
// When you have an interface method with signature (receiver, context.Context, A) (T, error),
// you can use Eitherize1 to convert it into a Kleisli arrow where the receiver becomes
// the context C and A becomes the input parameter.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input parameter type
//   - T: The success value type
//
// # Parameters
//
//   - f: A function that takes C, context.Context, and A, returning (T, error)
//
// # Returns
//
//   - Kleisli[C, A, T]: A function from A to Effect[C, T]
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
//	fetchUserKleisli := effect.Eitherize1(fetchUserByID)
//
//	// Use in functional composition with Chain
//	pipeline := F.Pipe1(
//	    effect.Succeed[AppConfig](123),
//	    effect.Chain[AppConfig](fetchUserKleisli),
//	)
//
//	// Execute with config
//	cfg := AppConfig{DatabaseURL: "postgres://localhost"}
//	result, err := effect.RunSync(effect.Provide[*User](cfg)(pipeline))(context.Background())
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
//	// Adapt the method - receiver becomes C, id becomes A
//	repo := &UserRepo{db: db}
//	getUserKleisli := effect.Eitherize1(func(r *UserRepo, ctx context.Context, id int) (*User, error) {
//	    return r.GetUserByID(ctx, id)
//	})
//
//	// Now getUserKleisli has type: Kleisli[*UserRepo, int, *User]
//	// Which is: func(int) Effect[*UserRepo, *User]
//	// Use it in composition:
//	pipeline := F.Pipe1(
//	    effect.Succeed[*UserRepo](123),
//	    effect.Chain[*UserRepo](getUserKleisli),
//	)
//	result, err := effect.RunSync(effect.Provide[*User](repo)(pipeline))(context.Background())
//
// # See Also
//
//   - Eitherize: For functions without an additional parameter
//   - Chain: For composing Kleisli arrows
//   - readerreaderioresult.Eitherize1: The underlying implementation
//
//go:inline
func Eitherize1[C, A, T any](f func(C, context.Context, A) (T, error)) Kleisli[C, A, T] {
	return readerreaderioresult.Eitherize1(f)
}
