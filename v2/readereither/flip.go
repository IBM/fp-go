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
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/reader"
)

// Sequence swaps the order of nested environment parameters in a ReaderEither computation.
//
// This function takes a ReaderEither that produces another ReaderEither and returns a
// reader.Kleisli that reverses the order of the environment parameters. The result is
// a curried function that takes R2 first, then R1, and produces an Either[E, A].
//
// Type Parameters:
//   - R1: The first environment type (becomes inner after flip)
//   - R2: The second environment type (becomes outer after flip)
//   - E: The error type
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderEither that takes R2 and may produce a ReaderEither[R1, E, A]
//
// Returns:
//   - A reader.Kleisli[R2, R1, Either[E, A]], which is func(R2) func(R1) Either[E, A]
//
// The function preserves error handling at both levels. Errors from the outer computation
// become errors in the inner Either result.
//
// Example:
//
//	import S "github.com/IBM/fp-go/v2/string"
//
//	type Database struct {
//	    ConnectionString string
//	}
//	type Config struct {
//	    Timeout int
//	}
//
//	// Original: takes Config, may fail, produces ReaderEither[Database, error, string]
//	original := func(cfg Config) either.Either[error, ReaderEither[Database, error, string]] {
//	    if cfg.Timeout <= 0 {
//	        return either.Left[ReaderEither[Database, error, string]](errors.New("invalid timeout"))
//	    }
//	    return either.Right[error](func(db Database) either.Either[error, string] {
//	        if S.IsEmpty(db.ConnectionString) {
//	            return either.Left[string](errors.New("empty connection string"))
//	        }
//	        return either.Right[error](fmt.Sprintf("Query on %s with timeout %d",
//	            db.ConnectionString, cfg.Timeout))
//	    })
//	}
//
//	// Sequenced: takes Database first, then Config
//	sequenced := Sequence(original)
//
//	db := Database{ConnectionString: "localhost:5432"}
//	cfg := Config{Timeout: 30}
//
//	// Apply database first to get a function that takes config
//	configReader := sequenced(db)
//	// Then apply config to get the final result
//	result := configReader(cfg)
//	// result is Either[error, string]
func Sequence[R1, R2, E, A any](ma ReaderEither[R2, E, ReaderEither[R1, E, A]]) Kleisli[R2, E, R1, A] {
	return readert.Sequence(
		either.Chain,
		ma,
	)
}

// SequenceReader swaps the order of environment parameters when the inner computation is a Reader.
//
// This function is similar to Sequence but specialized for the case where the innermost computation
// is a pure Reader (without error handling) rather than another ReaderEither. It takes a
// ReaderEither that produces a Reader and returns a Reader that produces a ReaderEither.
//
// Type Parameters:
//   - R1: The first environment type (becomes outer after flip)
//   - R2: The second environment type (becomes inner after flip)
//   - E: The error type (only present in the ReaderEither layer)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderEither that takes R2 and may produce a Reader[R1, A]
//
// Returns:
//   - A reader.Kleisli[R2, R1, Either[E, A]], which is func(R2) func(R1) Either[E, A]
//
// The function preserves error handling from the outer ReaderEither layer. If the outer
// computation fails, the error is propagated to the inner ReaderEither result.
//
// Example:
//
//	type Database struct {
//	    ConnectionString string
//	}
//	type Config struct {
//	    Timeout int
//	}
//
//	// Original: takes Config, may fail, produces Reader[Database, string]
//	original := func(cfg Config) either.Either[error, Reader[Database, string]] {
//	    if cfg.Timeout <= 0 {
//	        return either.Left[Reader[Database, string]](errors.New("invalid timeout"))
//	    }
//	    return either.Right[error](func(db Database) string {
//	        return fmt.Sprintf("Query on %s with timeout %d",
//	            db.ConnectionString, cfg.Timeout)
//	    })
//	}
//
//	// Sequenced: takes Database first, then Config
//	sequenced := SequenceReader(original)
//
//	db := Database{ConnectionString: "localhost:5432"}
//	cfg := Config{Timeout: 30}
//
//	// Apply database first to get a function that takes config
//	configReader := sequenced(db)
//	// Then apply config to get the final result
//	result := configReader(cfg)
//	// result is Either[error, string]
func SequenceReader[R1, R2, E, A any](ma ReaderEither[R2, E, Reader[R1, A]]) Kleisli[R2, E, R1, A] {
	return readert.SequenceReader(
		either.Map,
		ma,
	)
}

// Traverse transforms a ReaderEither computation by applying a function that produces
// another ReaderEither, effectively swapping the order of environment parameters.
//
// This function is useful when you have a computation that depends on environment R2 and
// produces a value of type A, and you want to transform it using a function that takes A
// and produces a computation depending on environment R1. The result is a curried function
// that takes R2 first, then R1, and produces an Either[E, B].
//
// Type Parameters:
//   - R2: The outer environment type (provided first)
//   - R1: The inner environment type (provided second)
//   - E: The error type
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Kleisli arrow that transforms A into a ReaderEither[R1, E, B]
//
// Returns:
//   - A function that takes a ReaderEither[R2, E, A] and returns a Kleisli[R2, E, R1, B],
//     which is func(R2) ReaderEither[R1, E, B]
//
// The function preserves error handling at both levels while reordering the environment dependencies.
//
// Example:
//
//	type Database struct {
//	    ConnectionString string
//	}
//	type Config struct {
//	    Timeout int
//	}
//
//	// Original: ReaderEither[Config, error, int] - takes Config, may fail, produces int
//	original := func(cfg Config) either.Either[error, int] {
//	    if cfg.Timeout <= 0 {
//	        return either.Left[int](errors.New("invalid timeout"))
//	    }
//	    return either.Right[error](cfg.Timeout * 10)
//	}
//
//	// Kleisli function: transforms int to ReaderEither[Database, error, string]
//	kleisli := func(value int) ReaderEither[Database, error, string] {
//	    return func(db Database) either.Either[error, string] {
//	        if S.IsEmpty(db.ConnectionString) {
//	            return either.Left[string](errors.New("empty connection string"))
//	        }
//	        return either.Right[error](fmt.Sprintf("%s:%d", db.ConnectionString, value))
//	    }
//	}
//
//	// Apply Traverse to get: func(ReaderEither[Config, error, int]) func(Database) ReaderEither[Config, error, string]
//	traversed := Traverse[Config, Database, error, int, string](kleisli)
//	result := traversed(original)
//
//	db := Database{ConnectionString: "localhost:5432"}
//	cfg := Config{Timeout: 30}
//
//	// Apply database first to get a function that takes config
//	configReader := result(db)
//	// Then apply config to get the final result
//	finalResult := configReader(cfg)
//	// finalResult is Either[error, string] = Right("localhost:5432:300")
func Traverse[R2, R1, E, A, B any](
	f Kleisli[R1, E, A, B],
) func(ReaderEither[R2, E, A]) Kleisli[R2, E, R1, B] {
	return readert.Traverse[ReaderEither[R2, E, A]](
		either.Map,
		either.Chain,
		f,
	)
}

func TraverseReader[R2, R1, E, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderEither[R2, E, A]) Kleisli[R2, E, R1, B] {
	return readert.TraverseReader[ReaderEither[R2, E, A]](
		either.Map,
		either.Map,
		f,
	)
}
