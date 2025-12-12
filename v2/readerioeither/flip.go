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

package readerioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/reader"
)

// Sequence swaps the order of nested environment parameters in a ReaderIOEither computation.
//
// This function takes a ReaderIOEither that produces another ReaderIOEither and returns a
// reader.Kleisli that reverses the order of the environment parameters. The result is
// a curried function that takes R2 first, then R1, and produces an IOEither[E, A].
//
// Type Parameters:
//   - R1: The first environment type (becomes inner after sequence)
//   - R2: The second environment type (becomes outer after sequence)
//   - E: The error type
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderIOEither that takes R2 and may produce a ReaderIOEither[R1, E, A]
//
// Returns:
//   - A Kleisli[R2, E, R1, A], which is func(R2) func(R1) IOEither[E, A]
//
// The function preserves error handling and IO effects at both levels.
func Sequence[R1, R2, E, A any](ma ReaderIOEither[R2, E, ReaderIOEither[R1, E, A]]) Kleisli[R2, E, R1, A] {
	return readert.Sequence(
		ioeither.Chain,
		ma,
	)
}

// SequenceReader swaps the order of environment parameters when the inner computation is a pure Reader.
//
// This function is similar to Sequence but specialized for the case where the innermost computation
// is a pure Reader (without IO or error handling) rather than another ReaderIOEither.
//
// Type Parameters:
//   - R1: The first environment type (becomes outer after sequence)
//   - R2: The second environment type (becomes inner after sequence)
//   - E: The error type (only present in the ReaderIOEither layer)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderIOEither that takes R2 and may produce a Reader[R1, A]
//
// Returns:
//   - A Kleisli[R2, E, R1, A], which is func(R2) func(R1) IOEither[E, A]
func SequenceReader[R1, R2, E, A any](ma ReaderIOEither[R2, E, Reader[R1, A]]) Kleisli[R2, E, R1, A] {
	return readert.SequenceReader(
		ioeither.Map,
		ma,
	)
}

// SequenceReaderIO swaps the order of environment parameters when the inner computation is a ReaderIO.
//
// This function is specialized for the case where the innermost computation is a ReaderIO
// (with IO effects but no error handling) rather than another ReaderIOEither.
//
// Type Parameters:
//   - R1: The first environment type (becomes outer after sequence)
//   - R2: The second environment type (becomes inner after sequence)
//   - E: The error type (only present in the outer ReaderIOEither layer)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderIOEither that takes R2 and may produce a ReaderIO[R1, A]
//
// Returns:
//   - A Kleisli[R2, E, R1, A], which is func(R2) func(R1) IOEither[E, A]
func SequenceReaderIO[R1, R2, E, A any](ma ReaderIOEither[R2, E, ReaderIO[R1, A]]) Kleisli[R2, E, R1, A] {
	return func(r1 R1) ReaderIOEither[R2, E, A] {
		return func(r2 R2) IOEither[E, A] {
			return func() Either[E, A] {
				return either.MonadMap(
					ma(r2)(),
					func(rr ReaderIO[R1, A]) A {
						return rr(r1)()
					},
				)
			}
		}
	}
}

// SequenceReaderEither swaps the order of environment parameters when the inner computation is a ReaderEither.
//
// This function is specialized for the case where the innermost computation is a ReaderEither
// (with error handling but no IO effects) rather than another ReaderIOEither.
//
// Type Parameters:
//   - R1: The first environment type (becomes outer after sequence)
//   - R2: The second environment type (becomes inner after sequence)
//   - E: The error type (present in both layers)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderIOEither that takes R2 and may produce a ReaderEither[R1, E, A]
//
// Returns:
//   - A Kleisli[R2, E, R1, A], which is func(R2) func(R1) IOEither[E, A]
func SequenceReaderEither[R1, R2, E, A any](ma ReaderIOEither[R2, E, ReaderEither[R1, E, A]]) Kleisli[R2, E, R1, A] {
	return func(r1 R1) ReaderIOEither[R2, E, A] {
		return func(r2 R2) IOEither[E, A] {
			return func() Either[E, A] {
				return either.MonadChain(
					ma(r2)(),
					identity.Ap[Either[E, A]](r1),
				)
			}
		}
	}
}

// Traverse transforms a ReaderIOEither computation by applying a function that produces
// another ReaderIOEither, effectively swapping the order of environment parameters.
//
// This function is useful when you have a computation that depends on environment R2 and
// produces a value of type A, and you want to transform it using a function that takes A
// and produces a computation depending on environment R1. The result is a curried function
// that takes R2 first, then R1, and produces an IOEither[E, B].
//
// Type Parameters:
//   - R2: The outer environment type (provided first)
//   - R1: The inner environment type (provided second)
//   - E: The error type
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Kleisli arrow that transforms A into a ReaderIOEither[R1, E, B]
//
// Returns:
//   - A function that takes a ReaderIOEither[R2, E, A] and returns a Kleisli[R2, E, R1, B],
//     which is func(R2) ReaderIOEither[R1, E, B]
//
// The function preserves error handling and IO effects while reordering the environment dependencies.
func Traverse[R2, R1, E, A, B any](
	f Kleisli[R1, E, A, B],
) func(ReaderIOEither[R2, E, A]) Kleisli[R2, E, R1, B] {
	return readert.Traverse[ReaderIOEither[R2, E, A]](
		ioeither.Map,
		ioeither.Chain,
		f,
	)
}

// TraverseReader transforms a ReaderIOEither computation by applying a Reader-based function,
// effectively introducing a new environment dependency.
//
// This function takes a Reader-based transformation (Kleisli arrow) and returns a function that
// can transform a ReaderIOEither. The result allows you to provide the Reader's environment (R1)
// first, which then produces a ReaderIOEither that depends on environment R2.
//
// Type Parameters:
//   - R2: The outer environment type (from the original ReaderIOEither)
//   - R1: The inner environment type (introduced by the Reader transformation)
//   - E: The error type
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Reader-based Kleisli arrow that transforms A to B using environment R1
//
// Returns:
//   - A function that takes a ReaderIOEither[R2, E, A] and returns a Kleisli[R2, E, R1, B],
//     which is func(R2) ReaderIOEither[R1, E, B]
//
// The function preserves error handling and IO effects while adding the Reader environment dependency.
//
// Example:
//
//	type Config struct {
//	    Multiplier int
//	}
//	type Database struct {
//	    ConnectionString string
//	}
//
//	// Original computation that depends on Database
//	original := func(db Database) IOEither[error, int] {
//	    return ioeither.Right[error](len(db.ConnectionString))
//	}
//
//	// Reader-based transformation that depends on Config
//	multiply := func(x int) func(Config) int {
//	    return func(cfg Config) int {
//	        return x * cfg.Multiplier
//	    }
//	}
//
//	// Apply TraverseReader to introduce Config dependency
//	traversed := TraverseReader[Database, Config, error, int, int](multiply)
//	result := traversed(original)
//
//	// Provide Config first, then Database
//	cfg := Config{Multiplier: 5}
//	db := Database{ConnectionString: "localhost:5432"}
//	finalResult := result(cfg)(db)() // Returns Right(80) = len("localhost:5432") * 5
func TraverseReader[R2, R1, E, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderIOEither[R2, E, A]) Kleisli[R2, E, R1, B] {
	return readert.TraverseReader[ReaderIOEither[R2, E, A]](
		ioeither.Map,
		ioeither.Map,
		f,
	)
}
