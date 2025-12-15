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
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
)

// Sequence swaps the order of nested environment parameters in a ReaderResult computation.
//
// This function takes a ReaderResult that produces another ReaderResult and returns a
// reader.Kleisli that reverses the order of the environment parameters. The result is
// a curried function that takes R2 first, then R1, and produces a Result[A].
//
// Type Parameters:
//   - R1: The first environment type (becomes inner after flip)
//   - R2: The second environment type (becomes outer after flip)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderResult that takes R2 and may produce a ReaderResult[R1, A]
//
// Returns:
//   - A reader.Kleisli[R2, R1, Result[A]], which is func(R2) func(R1) Result[A]
//
// The function preserves error handling at both levels. Errors from the outer computation
// become errors in the inner Result.
//
// Note: This is an inline wrapper around readereither.Sequence since ReaderResult is an alias
// for ReaderEither with error type fixed to error.
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
//	// Original: takes Config, may fail, produces ReaderResult[Database, string]
//	original := func(cfg Config) result.Result[ReaderResult[Database, string]] {
//	    if cfg.Timeout <= 0 {
//	        return result.Error[ReaderResult[Database, string]](errors.New("invalid timeout"))
//	    }
//	    return result.Ok[error](func(db Database) result.Result[string] {
//	        if S.IsEmpty(db.ConnectionString) {
//	            return result.Error[string](errors.New("empty connection string"))
//	        }
//	        return result.Ok[error](fmt.Sprintf("Query on %s with timeout %d",
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
//	// result is Result[string]
//
//go:inline
func Sequence[R1, R2, A any](ma ReaderResult[R2, ReaderResult[R1, A]]) reader.Kleisli[R2, R1, Result[A]] {
	return readereither.Sequence(ma)
}

// SequenceReader swaps the order of environment parameters when the inner computation is a Reader.
//
// This function is similar to Sequence but specialized for the case where the innermost computation
// is a pure Reader (without error handling) rather than another ReaderResult. It takes a
// ReaderResult that produces a Reader and returns a reader.Kleisli that produces Results.
//
// Type Parameters:
//   - R1: The first environment type (becomes outer after flip)
//   - R2: The second environment type (becomes inner after flip)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderResult that takes R2 and may produce a Reader[R1, A]
//
// Returns:
//   - A reader.Kleisli[R2, R1, Result[A]], which is func(R2) func(R1) Result[A]
//
// The function preserves error handling from the outer ReaderResult layer. If the outer
// computation fails, the error is propagated to the inner Result.
//
// Note: This is an inline wrapper around readereither.SequenceReader since ReaderResult is an alias
// for ReaderEither with error type fixed to error.
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
//	original := func(cfg Config) result.Result[Reader[Database, string]] {
//	    if cfg.Timeout <= 0 {
//	        return result.Error[Reader[Database, string]](errors.New("invalid timeout"))
//	    }
//	    return result.Ok[error](func(db Database) string {
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
//	// result is Result[string]
//
//go:inline
func SequenceReader[R1, R2, A any](ma ReaderResult[R2, Reader[R1, A]]) reader.Kleisli[R2, R1, Result[A]] {
	return readereither.SequenceReader(ma)
}

func Traverse[R2, R1, A, B any](
	f Kleisli[R1, A, B],
) func(ReaderResult[R2, A]) Kleisli[R2, R1, B] {
	return readereither.Traverse[R2](f)
}

func TraverseReader[R2, R1, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderResult[R2, A]) Kleisli[R2, R1, B] {
	return readereither.TraverseReader[R2, R1, error](f)
}
