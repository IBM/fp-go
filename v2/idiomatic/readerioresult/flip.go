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
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/reader"
)

// Sequence swaps the order of nested environment parameters in a ReaderIOResult computation.
//
// This function transforms a computation that takes environment R2 and produces a ReaderIOResult[R1, A]
// into a Kleisli arrow that takes R1 first and returns a ReaderIOResult[R2, A].
//
// Type Parameters:
//   - R1: The type of the inner environment (becomes the outer parameter after sequencing)
//   - R2: The type of the outer environment (becomes the inner environment after sequencing)
//   - A: The type of the value produced by the computation
//
// Parameters:
//   - ma: A ReaderIOResult that depends on R2 and produces a ReaderIOResult[R1, A]
//
// Returns:
//   - A Kleisli arrow (func(R1) func(R2) func() (A, error)) that reverses the environment order
//
// The transformation preserves error handling - if the outer computation fails, the error
// is propagated; if the inner computation fails, that error is also propagated.
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
//	// Original: takes Config, produces ReaderIOResult[Database, string]
//	original := func(cfg Config) func() (func(Database) func() (string, error), error) {
//	    return func() (func(Database) func() (string, error), error) {
//	        if cfg.Timeout <= 0 {
//	            return nil, errors.New("invalid timeout")
//	        }
//	        return func(db Database) func() (string, error) {
//	            return func() (string, error) {
//	                if S.IsEmpty(db.ConnectionString) {
//	                    return "", errors.New("empty connection")
//	                }
//	                return fmt.Sprintf("Query on %s with timeout %d",
//	                    db.ConnectionString, cfg.Timeout), nil
//	            }
//	        }, nil
//	    }
//	}
//
//	// Sequenced: takes Database first, then Config
//	sequenced := Sequence(original)
//	db := Database{ConnectionString: "localhost:5432"}
//	cfg := Config{Timeout: 30}
//	result, err := sequenced(db)(cfg)()
//	// result: "Query on localhost:5432 with timeout 30"
func Sequence[R1, R2, A any](ma ReaderIOResult[R2, ReaderIOResult[R1, A]]) reader.Kleisli[R2, R1, IOResult[A]] {
	return readert.Sequence(
		ioresult.Chain,
		ma,
	)
}

// SequenceReader swaps the order of environment parameters when the inner computation is a pure Reader.
//
// This function is similar to Sequence but specialized for cases where the inner computation
// is a Reader (pure function) rather than a ReaderIOResult. It transforms a ReaderIOResult that
// produces a Reader into a Kleisli arrow with swapped environment order.
//
// Type Parameters:
//   - R1: The type of the Reader's environment (becomes the outer parameter after sequencing)
//   - R2: The type of the ReaderIOResult's environment (becomes the inner environment after sequencing)
//   - A: The type of the value produced by the computation
//
// Parameters:
//   - ma: A ReaderIOResult[R2, Reader[R1, A]] - depends on R2 and produces a pure Reader
//
// Returns:
//   - A Kleisli arrow (func(R1) func(R2) func() (A, error)) that reverses the environment order
//
// The inner Reader computation is automatically lifted into the IOResult context (cannot fail).
// Only the outer ReaderIOResult can fail with an error.
//
// Example:
//
//	type Config struct {
//	    Multiplier int
//	}
//
//	// Original: takes int, produces Reader[Config, int]
//	original := func(x int) func() (func(Config) int, error) {
//	    return func() (func(Config) int, error) {
//	        if x < 0 {
//	            return nil, errors.New("negative value")
//	        }
//	        return func(cfg Config) int {
//	            return x * cfg.Multiplier
//	        }, nil
//	    }
//	}
//
//	// Sequenced: takes Config first, then int
//	sequenced := SequenceReader(original)
//	cfg := Config{Multiplier: 5}
//	result, err := sequenced(cfg)(10)()
//	// result: 50, err: nil
func SequenceReader[R1, R2, A any](ma ReaderIOResult[R2, Reader[R1, A]]) reader.Kleisli[R2, R1, IOResult[A]] {
	return readert.SequenceReader(
		ioresult.Map,
		ma,
	)
}

// Traverse transforms a ReaderIOResult computation by applying a Kleisli arrow that introduces
// a new environment dependency, effectively swapping the environment order.
//
// This is a higher-order function that takes a Kleisli arrow and returns a function that
// can transform ReaderIOResult computations. It's useful for introducing environment-dependent
// transformations into existing computations while reordering the environment parameters.
//
// Type Parameters:
//   - R2: The type of the original computation's environment
//   - R1: The type of the new environment introduced by the Kleisli arrow
//   - A: The input type to the Kleisli arrow
//   - B: The output type of the transformation
//
// Parameters:
//   - f: A Kleisli arrow (func(A) ReaderIOResult[R1, B]) that transforms A to B with R1 dependency
//
// Returns:
//   - A function that transforms ReaderIOResult[R2, A] into a Kleisli arrow with swapped environments
//
// The transformation preserves error handling from both the original computation and the
// Kleisli arrow. The resulting computation takes R1 first, then R2.
//
// Example:
//
//	type Database struct {
//	    Prefix string
//	}
//
//	// Original computation: depends on int environment
//	original := func(x int) func() (int, error) {
//	    return func() (int, error) {
//	        if x < 0 {
//	            return 0, errors.New("negative value")
//	        }
//	        return x * 2, nil
//	    }
//	}
//
//	// Kleisli arrow: transforms int to string with Database dependency
//	format := func(value int) func(Database) func() (string, error) {
//	    return func(db Database) func() (string, error) {
//	        return func() (string, error) {
//	            return fmt.Sprintf("%s:%d", db.Prefix, value), nil
//	        }
//	    }
//	}
//
//	// Apply Traverse
//	traversed := Traverse[int](format)
//	result := traversed(original)
//
//	// Use with Database first, then int
//	db := Database{Prefix: "ID"}
//	output, err := result(db)(10)()
//	// output: "ID:20", err: nil
func Traverse[R2, R1, A, B any](
	f Kleisli[R1, A, B],
) func(ReaderIOResult[R2, A]) Kleisli[R2, R1, B] {
	return readert.Traverse[ReaderIOResult[R2, A]](
		ioresult.Map,
		ioresult.Chain,
		f,
	)
}

// TraverseReader transforms a ReaderIOResult computation by applying a Reader-based Kleisli arrow,
// introducing a new environment dependency while swapping the environment order.
//
// This function is similar to Traverse but specialized for pure Reader transformations that
// cannot fail. It's useful when you want to introduce environment-dependent logic without
// adding error handling complexity.
//
// Type Parameters:
//   - R2: The type of the original computation's environment
//   - R1: The type of the new environment introduced by the Reader Kleisli arrow
//   - A: The input type to the Kleisli arrow
//   - B: The output type of the transformation
//
// Parameters:
//   - f: A Reader Kleisli arrow (func(A) func(R1) B) that transforms A to B with R1 dependency
//
// Returns:
//   - A function that transforms ReaderIOResult[R2, A] into a Kleisli arrow with swapped environments
//
// The Reader transformation is automatically lifted into the IOResult context. Only the original
// ReaderIOResult computation can fail; the Reader transformation itself is pure and cannot fail.
//
// Example:
//
//	type Config struct {
//	    Multiplier int
//	}
//
//	// Original computation: depends on int environment, may fail
//	original := func(x int) func() (int, error) {
//	    return func() (int, error) {
//	        if x < 0 {
//	            return 0, errors.New("negative value")
//	        }
//	        return x * 2, nil
//	    }
//	}
//
//	// Pure Reader transformation: multiplies by config value
//	multiply := func(value int) func(Config) int {
//	    return func(cfg Config) int {
//	        return value * cfg.Multiplier
//	    }
//	}
//
//	// Apply TraverseReader
//	traversed := TraverseReader[int, Config](multiply)
//	result := traversed(original)
//
//	// Use with Config first, then int
//	cfg := Config{Multiplier: 5}
//	output, err := result(cfg)(10)()
//	// output: 100 (10 * 2 * 5), err: nil
func TraverseReader[R2, R1, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderIOResult[R2, A]) Kleisli[R2, R1, B] {
	return readert.TraverseReader[ReaderIOResult[R2, A]](
		ioresult.Map,
		ioresult.Map,
		f,
	)
}
