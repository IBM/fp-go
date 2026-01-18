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

	"github.com/IBM/fp-go/v2/reader"
	RR "github.com/IBM/fp-go/v2/readerresult"
)

// SequenceReader swaps the order of environment parameters when the inner computation is a Reader.
//
// This function is specialized for the context.Context-based ReaderResult monad. It takes a
// ReaderResult that produces a Reader and returns a reader.Kleisli that produces Results.
// The context.Context is implicitly used as the outer environment type.
//
// Type Parameters:
//   - R: The inner environment type (becomes outer after flip)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderResult that takes context.Context and may produce a Reader[R, A]
//
// Returns:
//   - A reader.Kleisli[context.Context, R, Result[A]], which is func(context.Context) func(R) Result[A]
//
// The function preserves error handling from the outer ReaderResult layer. If the outer
// computation fails, the error is propagated to the inner Result.
//
// Note: This is an inline wrapper around readerresult.SequenceReader, specialized for
// context.Context as the outer environment type.
//
// Example:
//
//	type Database struct {
//	    ConnectionString string
//	}
//
//	// Original: takes context, may fail, produces Reader[Database, string]
//	original := func(ctx context.Context) result.Result[reader.Reader[Database, string]] {
//	    if ctx.Err() != nil {
//	        return result.Error[reader.Reader[Database, string]](ctx.Err())
//	    }
//	    return result.Ok[error](func(db Database) string {
//	        return fmt.Sprintf("Query on %s", db.ConnectionString)
//	    })
//	}
//
//	// Sequenced: takes context first, then Database
//	sequenced := SequenceReader(original)
//
//	ctx := t.Context()
//	db := Database{ConnectionString: "localhost:5432"}
//
//	// Apply context first to get a function that takes database
//	dbReader := sequenced(ctx)
//	// Then apply database to get the final result
//	result := dbReader(db)
//	// result is Result[string]
//
// Use Cases:
//   - Dependency injection: Flip parameter order to inject context first, then dependencies
//   - Testing: Separate context handling from business logic for easier testing
//   - Composition: Enable point-free style by fixing the context parameter first
//
//go:inline
func SequenceReader[R, A any](ma ReaderResult[Reader[R, A]]) reader.Kleisli[context.Context, R, Result[A]] {
	return RR.SequenceReader(ma)
}

// TraverseReader transforms a value using a Reader function and swaps environment parameter order.
//
// This function combines mapping and parameter flipping in a single operation. It takes a
// Reader function (pure computation without error handling) and returns a function that:
// 1. Maps a ReaderResult[A] to ReaderResult[B] using the provided Reader function
// 2. Flips the parameter order so R comes before context.Context
//
// Type Parameters:
//   - R: The inner environment type (becomes outer after flip)
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A reader.Kleisli[R, A, B], which is func(R) func(A) B - a pure Reader function
//
// Returns:
//   - A function that takes ReaderResult[A] and returns Kleisli[R, B]
//   - Kleisli[R, B] is func(R) ReaderResult[B], which is func(R) func(context.Context) Result[B]
//
// The function preserves error handling from the input ReaderResult. If the input computation
// fails, the error is propagated without applying the transformation function.
//
// Note: This is a wrapper around readerresult.TraverseReader, specialized for context.Context.
//
// Example:
//
//	type Config struct {
//	    MaxRetries int
//	}
//
//	// A pure Reader function that depends on Config
//	formatMessage := func(cfg Config) func(int) string {
//	    return func(value int) string {
//	        return fmt.Sprintf("Value: %d, MaxRetries: %d", value, cfg.MaxRetries)
//	    }
//	}
//
//	// Original computation that may fail
//	computation := func(ctx context.Context) result.Result[int] {
//	    if ctx.Err() != nil {
//	        return result.Error[int](ctx.Err())
//	    }
//	    return result.Ok[error](42)
//	}
//
//	// Create a traversal that applies formatMessage and flips parameters
//	traverse := TraverseReader[Config, int, string](formatMessage)
//
//	// Apply to the computation
//	flipped := traverse(computation)
//
//	// Now we can provide Config first, then context
//	cfg := Config{MaxRetries: 3}
//	ctx := t.Context()
//
//	result := flipped(cfg)(ctx)
//	// result is Result[string] containing "Value: 42, MaxRetries: 3"
//
// Use Cases:
//   - Dependency injection: Inject configuration/dependencies before context
//   - Testing: Separate pure business logic from context handling
//   - Composition: Build pipelines where dependencies are fixed before execution
//   - Point-free style: Enable partial application by fixing dependencies first
//
//go:inline
func TraverseReader[R, A, B any](
	f reader.Kleisli[R, A, B],
) func(ReaderResult[A]) Kleisli[R, B] {
	return RR.TraverseReader[context.Context](f)
}
