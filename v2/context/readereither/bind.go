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

package readereither

import (
	"context"

	G "github.com/IBM/fp-go/v2/readereither/generic"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    UserID   string
//	    TenantID string
//	}
//	result := readereither.Do(State{})
func Do[S any](
	empty S,
) ReaderEither[S] {
	return G.Do[ReaderEither[S], context.Context, error, S](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps
// and access the context.Context from the environment.
//
// The setter function takes the result of the computation and returns a function that
// updates the context from S1 to S2.
//
// Example:
//
//	type State struct {
//	    UserID   string
//	    TenantID string
//	}
//
//	result := F.Pipe2(
//	    readereither.Do(State{}),
//	    readereither.Bind(
//	        func(uid string) func(State) State {
//	            return func(s State) State { s.UserID = uid; return s }
//	        },
//	        func(s State) readereither.ReaderEither[string] {
//	            return func(ctx context.Context) either.Either[error, string] {
//	                if uid, ok := ctx.Value("userID").(string); ok {
//	                    return either.Right[error](uid)
//	                }
//	                return either.Left[string](errors.New("no userID"))
//	            }
//	        },
//	    ),
//	    readereither.Bind(
//	        func(tid string) func(State) State {
//	            return func(s State) State { s.TenantID = tid; return s }
//	        },
//	        func(s State) readereither.ReaderEither[string] {
//	            // This can access s.UserID from the previous step
//	            return func(ctx context.Context) either.Either[error, string] {
//	                return either.Right[error]("tenant-" + s.UserID)
//	            }
//	        },
//	    ),
//	)
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) ReaderEither[T],
) func(ReaderEither[S1]) ReaderEither[S2] {
	return G.Bind[ReaderEither[S1], ReaderEither[S2], ReaderEither[T], context.Context, error, S1, S2, T](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(ReaderEither[S1]) ReaderEither[S2] {
	return G.Let[ReaderEither[S1], ReaderEither[S2], context.Context, error, S1, S2, T](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(ReaderEither[S1]) ReaderEither[S2] {
	return G.LetTo[ReaderEither[S1], ReaderEither[S2], context.Context, error, S1, S2, T](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) func(ReaderEither[T]) ReaderEither[S1] {
	return G.BindTo[ReaderEither[S1], ReaderEither[T], context.Context, error, S1, T](setter)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
// This allows independent computations to be combined without one depending on the result of the other.
//
// Unlike Bind, which sequences operations, ApS can be used when operations are independent
// and can conceptually run in parallel.
//
// Example:
//
//	type State struct {
//	    UserID   string
//	    TenantID string
//	}
//
//	// These operations are independent and can be combined with ApS
//	getUserID := func(ctx context.Context) either.Either[error, string] {
//	    return either.Right[error](ctx.Value("userID").(string))
//	}
//	getTenantID := func(ctx context.Context) either.Either[error, string] {
//	    return either.Right[error](ctx.Value("tenantID").(string))
//	}
//
//	result := F.Pipe2(
//	    readereither.Do(State{}),
//	    readereither.ApS(
//	        func(uid string) func(State) State {
//	            return func(s State) State { s.UserID = uid; return s }
//	        },
//	        getUserID,
//	    ),
//	    readereither.ApS(
//	        func(tid string) func(State) State {
//	            return func(s State) State { s.TenantID = tid; return s }
//	        },
//	        getTenantID,
//	    ),
//	)
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderEither[T],
) func(ReaderEither[S1]) ReaderEither[S2] {
	return G.ApS[ReaderEither[S1], ReaderEither[S2], ReaderEither[T], context.Context, error, S1, S2, T](setter, fa)
}
