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

	"github.com/IBM/fp-go/v2/reader"
	RIO "github.com/IBM/fp-go/v2/readerio"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
	RR "github.com/IBM/fp-go/v2/readerresult"
)

// SequenceReader transforms a ReaderIOResult containing a Reader into a function that
// takes the Reader's environment first, then returns a ReaderIOResult.
//
// This function "flips" or "sequences" the nested structure, changing the order in which
// parameters are applied. It's particularly useful for point-free style programming where
// you want to partially apply the inner Reader's environment before dealing with the
// outer context.
//
// Type transformation:
//
//	From: ReaderIOResult[Reader[R, A]]
//	      = func(context.Context) func() Either[error, func(R) A]
//
//	To:   func(context.Context) func(R) IOResult[A]
//	      = func(context.Context) func(R) func() Either[error, A]
//
// This allows you to:
//  1. Provide the context.Context first
//  2. Then provide the Reader's environment R
//  3. Finally execute the IO effect to get Either[error, A]
//
// Point-free style benefits:
//   - Enables partial application of the Reader environment
//   - Facilitates composition of Reader-based computations
//   - Allows building reusable computation pipelines
//   - Supports dependency injection patterns where R represents dependencies
//
// Example:
//
//	type Config struct {
//	    Timeout int
//	}
//
//	// A computation that produces a Reader based on context
//	func getMultiplier(ctx context.Context) func() Either[error, func(Config) int] {
//	    return func() Either[error, func(Config) int] {
//	        return Right[error](func(cfg Config) int {
//	            return cfg.Timeout * 2
//	        })
//	    }
//	}
//
//	// Sequence it to apply Config first
//	sequenced := SequenceReader[Config, int](getMultiplier)
//
//	// Now we can partially apply the Config
//	cfg := Config{Timeout: 30}
//	ctx := t.Context()
//	result := sequenced(ctx)(cfg)() // Returns Right(60)
//
// This is especially useful in point-free style when building computation pipelines:
//
//	var pipeline = F.Flow3(
//	    loadConfig,           // ReaderIOResult[Reader[Database, Config]]
//	    SequenceReader,       // func(context.Context) func(Database) IOResult[Config]
//	    applyToDatabase(db),  // IOResult[Config]
//	)
//
//go:inline
func SequenceReader[R, A any](ma ReaderIOResult[Reader[R, A]]) Kleisli[R, A] {
	return RIOR.SequenceReader(ma)
}

// SequenceReaderIO transforms a ReaderIOResult containing a ReaderIO into a function that
// takes the ReaderIO's environment first, then returns a ReaderIOResult.
//
// This is similar to SequenceReader but works with ReaderIO, which represents a computation
// that depends on an environment R and performs IO effects.
//
// Type transformation:
//
//	From: ReaderIOResult[ReaderIO[R, A]]
//	      = func(context.Context) func() Either[error, func(R) func() A]
//
//	To:   func(context.Context) func(R) IOResult[A]
//	      = func(context.Context) func(R) func() Either[error, A]
//
// The key difference from SequenceReader is that the inner computation (ReaderIO) already
// performs IO effects, so the sequencing combines these effects properly.
//
// Point-free style benefits:
//   - Enables composition of ReaderIO-based computations
//   - Allows partial application of environment before IO execution
//   - Facilitates building effect pipelines with dependency injection
//   - Supports layered architecture where R represents service dependencies
//
// Example:
//
//	type Database struct {
//	    ConnectionString string
//	}
//
//	// A computation that produces a ReaderIO based on context
//	func getQuery(ctx context.Context) func() Either[error, func(Database) func() string] {
//	    return func() Either[error, func(Database) func() string] {
//	        return Right[error](func(db Database) func() string {
//	            return func() string {
//	                // Perform actual IO here
//	                return "Query result from " + db.ConnectionString
//	            }
//	        })
//	    }
//	}
//
//	// Sequence it to apply Database first
//	sequenced := SequenceReaderIO[Database, string](getQuery)
//
//	// Partially apply the Database
//	db := Database{ConnectionString: "localhost:5432"}
//	ctx := t.Context()
//	result := sequenced(ctx)(db)() // Executes IO and returns Right("Query result...")
//
// In point-free style, this enables clean composition:
//
//	var executeQuery = F.Flow3(
//	    prepareQuery,         // ReaderIOResult[ReaderIO[Database, QueryResult]]
//	    SequenceReaderIO,     // func(context.Context) func(Database) IOResult[QueryResult]
//	    withDatabase(db),     // IOResult[QueryResult]
//	)
//
//go:inline
func SequenceReaderIO[R, A any](ma ReaderIOResult[RIO.ReaderIO[R, A]]) Kleisli[R, A] {
	return RIOR.SequenceReaderIO(ma)
}

// SequenceReaderResult transforms a ReaderIOResult containing a ReaderResult into a function
// that takes the ReaderResult's environment first, then returns a ReaderIOResult.
//
// This is similar to SequenceReader but works with ReaderResult, which represents a computation
// that depends on an environment R and can fail with an error.
//
// Type transformation:
//
//	From: ReaderIOResult[ReaderResult[R, A]]
//	      = func(context.Context) func() Either[error, func(R) Either[error, A]]
//
//	To:   func(context.Context) func(R) IOResult[A]
//	      = func(context.Context) func(R) func() Either[error, A]
//
// The sequencing properly combines the error handling from both the outer ReaderIOResult
// and the inner ReaderResult, ensuring that errors from either level are propagated correctly.
//
// Point-free style benefits:
//   - Enables composition of error-handling computations with dependency injection
//   - Allows partial application of dependencies before error handling
//   - Facilitates building validation pipelines with environment dependencies
//   - Supports service-oriented architectures with proper error propagation
//
// Example:
//
//	type Config struct {
//	    MaxRetries int
//	}
//
//	// A computation that produces a ReaderResult based on context
//	func validateRetries(ctx context.Context) func() Either[error, func(Config) Either[error, int]] {
//	    return func() Either[error, func(Config) Either[error, int]] {
//	        return Right[error](func(cfg Config) Either[error, int] {
//	            if cfg.MaxRetries < 0 {
//	                return Left[int](errors.New("negative retries"))
//	            }
//	            return Right[error](cfg.MaxRetries)
//	        })
//	    }
//	}
//
//	// Sequence it to apply Config first
//	sequenced := SequenceReaderResult[Config, int](validateRetries)
//
//	// Partially apply the Config
//	cfg := Config{MaxRetries: 3}
//	ctx := t.Context()
//	result := sequenced(ctx)(cfg)() // Returns Right(3)
//
//	// With invalid config
//	badCfg := Config{MaxRetries: -1}
//	badResult := sequenced(ctx)(badCfg)() // Returns Left(error("negative retries"))
//
// In point-free style, this enables validation pipelines:
//
//	var validateAndProcess = F.Flow4(
//	    loadConfig,              // ReaderIOResult[ReaderResult[Config, Settings]]
//	    SequenceReaderResult,    // func(context.Context) func(Config) IOResult[Settings]
//	    applyConfig(cfg),        // IOResult[Settings]
//	    Chain(processSettings),  // IOResult[Result]
//	)
//
//go:inline
func SequenceReaderResult[R, A any](ma ReaderIOResult[RR.ReaderResult[R, A]]) Kleisli[R, A] {
	return RIOR.SequenceReaderEither(ma)
}

// TraverseReader transforms a ReaderIOResult computation by applying a Reader-based function,
// effectively introducing a new environment dependency.
//
// This function takes a Reader-based transformation (Kleisli arrow) and returns a function that
// can transform a ReaderIOResult. The result allows you to provide the Reader's environment (R)
// first, which then produces a ReaderIOResult that depends on the context.
//
// Type transformation:
//
//	From: ReaderIOResult[A]
//	      = func(context.Context) func() Either[error, A]
//
//	With: reader.Kleisli[R, A, B]
//	      = func(A) func(R) B
//
//	To:   func(ReaderIOResult[A]) func(R) ReaderIOResult[B]
//	      = func(ReaderIOResult[A]) func(R) func(context.Context) func() Either[error, B]
//
// This enables:
//  1. Transforming values within a ReaderIOResult using environment-dependent logic
//  2. Introducing new environment dependencies into existing computations
//  3. Building composable pipelines where transformations depend on configuration or dependencies
//  4. Point-free style composition with Reader-based transformations
//
// Type Parameters:
//   - R: The environment type that the Reader depends on
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Reader-based Kleisli arrow that transforms A to B using environment R
//
// Returns:
//   - A function that takes a ReaderIOResult[A] and returns a Kleisli[R, B],
//     which is func(R) ReaderIOResult[B]
//
// The function preserves error handling and IO effects while adding the Reader environment dependency.
//
// Example:
//
//	type Config struct {
//	    Multiplier int
//	}
//
//	// A Reader-based transformation that depends on Config
//	multiply := func(x int) func(Config) int {
//	    return func(cfg Config) int {
//	        return x * cfg.Multiplier
//	    }
//	}
//
//	// Original computation that produces an int
//	computation := Right[int](10)
//
//	// Apply TraverseReader to introduce Config dependency
//	traversed := TraverseReader[Config, int, int](multiply)
//	result := traversed(computation)
//
//	// Now we can provide the Config to get the final result
//	cfg := Config{Multiplier: 5}
//	ctx := t.Context()
//	finalResult := result(cfg)(ctx)() // Returns Right(50)
//
// In point-free style, this enables clean composition:
//
//	var pipeline = F.Flow3(
//	    loadValue,                    // ReaderIOResult[int]
//	    TraverseReader(multiplyByConfig), // func(Config) ReaderIOResult[int]
//	    applyConfig(cfg),             // ReaderIOResult[int]
//	)
//
//go:inline
func TraverseReader[R, A, B any](
	f reader.Kleisli[R, A, B],
) func(ReaderIOResult[A]) Kleisli[R, B] {
	return RIOR.TraverseReader[context.Context](f)
}
