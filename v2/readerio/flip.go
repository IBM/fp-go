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

package readerio

import (
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
)

// Sequence swaps the order of nested environment parameters in a ReaderIO computation.
//
// This function takes a ReaderIO that produces another ReaderIO and returns a
// reader.Kleisli that reverses the order of the environment parameters. The result is
// a curried function that takes R2 first, then R1, and produces an IO[A].
//
// Type Parameters:
//   - R1: The first environment type (becomes inner after flip)
//   - R2: The second environment type (becomes outer after flip)
//   - A: The result type
//
// Parameters:
//   - ma: A ReaderIO that takes R2 and produces a ReaderIO[R1, A]
//
// Returns:
//   - A reader.Kleisli[R2, R1, IO[A]], which is func(R2) func(R1) IO[A]
//
// The function preserves IO effects at both levels. The transformation can be visualized as:
//
//	Before: R2 -> IO[R1 -> IO[A]]
//	After:  R2 -> (R1 -> IO[A])
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
//	// Original: takes Config with IO, produces ReaderIO[Database, string]
//	original := func(cfg Config) io.IO[ReaderIO[Database, string]] {
//	    return io.Of(func(db Database) io.IO[string] {
//	        return io.Of(fmt.Sprintf("Query on %s with timeout %d",
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
//	// Then apply config to get the final IO result
//	result := configReader(cfg)
//	// result is IO[string]
func Sequence[R1, R2, A any](ma ReaderIO[R2, ReaderIO[R1, A]]) Kleisli[R2, R1, A] {
	return readert.Sequence(
		io.Chain,
		ma,
	)
}

// SequenceReader swaps the order of environment parameters when the inner computation is a Reader.
//
// This function is similar to Sequence but specialized for the case where the innermost computation
// is a pure Reader (without IO effects) rather than another ReaderIO. It takes a
// ReaderIO that produces a Reader and returns a reader.Kleisli that produces IO effects.
//
// This function is useful when you have a computation that depends on two environments where
// the outer environment is wrapped in IO (ReaderIO) and the inner is pure (Reader), and you need
// to change the order in which they are applied. The result moves the IO effect to the inner level.
//
// Type Parameters:
//   - R1: The first environment type (from inner Reader, becomes outer after flip)
//   - R2: The second environment type (from outer ReaderIO, becomes inner after flip)
//   - A: The result type
//
// Parameters:
//   - ma: A ReaderIO[R2, Reader[R1, A]] - a computation that takes R2 with IO effects and produces a pure Reader[R1, A]
//
// Returns:
//   - A reader.Kleisli[R2, R1, IO[A]], which is func(R2) func(R1) IO[A]
//
// The transformation can be visualized as:
//
//	Before: R2 -> IO[R1 -> A]
//	After:  R2 -> (R1 -> IO[A])
//
// Note the key difference from Sequence:
//   - Sequence: Both levels have IO effects (ReaderIO[R2, ReaderIO[R1, A]] -> ReaderIO[R1, ReaderIO[R2, A]])
//   - SequenceReader: IO moves from outer to inner (ReaderIO[R2, Reader[R1, A]] -> Reader[R1, ReaderIO[R2, A]])
//
// Example:
//
//	type Database struct { ConnectionString string }
//	type Config struct { Timeout int }
//
//	// Original: takes Database with IO, returns pure Reader[Config, string]
//	query := func(db Database) io.IO[reader.Reader[Config, string]] {
//	    return io.Of(func(cfg Config) string {
//	        return fmt.Sprintf("Query on %s with timeout %d", db.ConnectionString, cfg.Timeout)
//	    })
//	}
//
//	// Sequenced: takes Config first, then Database
//	sequenced := readerio.SequenceReader(query)
//
//	db := Database{ConnectionString: "localhost:5432"}
//	cfg := Config{Timeout: 30}
//
//	// Apply config first to get a function that takes database
//	dbReader := sequenced(cfg)
//	// Then apply database to get the final IO result
//	result := dbReader(db)
//	// result is IO[string]
//
// Use cases:
//   - Reordering dependencies when you want to defer IO effects
//   - Adapting functions where the IO effect should be associated with a different parameter
//   - Building pipelines that need pure outer layers with effectful inner computations
//   - Optimizing by controlling which environment access triggers IO effects
func SequenceReader[R1, R2, A any](ma ReaderIO[R2, Reader[R1, A]]) Kleisli[R2, R1, A] {
	return readert.SequenceReader(
		io.Map,
		ma,
	)
}

func Traverse[R2, R1, A, B any](
	f Kleisli[R1, A, B],
) func(ReaderIO[R2, A]) Kleisli[R2, R1, B] {
	return readert.Traverse[ReaderIO[R2, A]](
		io.Map,
		io.Chain,
		f,
	)
}

func TraverseReader[R2, R1, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderIO[R2, A]) Kleisli[R2, R1, B] {
	return readert.TraverseReader[ReaderIO[R2, A]](
		io.Map,
		io.Map,
		f,
	)
}
