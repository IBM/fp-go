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

package readerreaderioeither

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readerioeither"
)

// Sequence swaps the order of nested environment parameters in a ReaderReaderIOEither computation.
//
// This function takes a ReaderReaderIOEither that produces another ReaderReaderIOEither and returns a
// Kleisli arrow that reverses the order of the outer environment parameters (R1 and R2). The result is
// a curried function that takes R1 first, then R2, and produces a ReaderIOEither[C, E, A].
//
// Type Parameters:
//   - R1: The first outer environment type (becomes the outermost after sequence)
//   - R2: The second outer environment type (becomes inner after sequence)
//   - C: The inner context/environment type (for the ReaderIOEither layer)
//   - E: The error type
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderReaderIOEither[R2, C, E, ReaderReaderIOEither[R1, C, E, A]]
//
// Returns:
//   - A Kleisli[R2, C, E, R1, A], which is func(R1) ReaderReaderIOEither[R2, C, E, A]
//
// The function preserves error handling and IO effects at all levels while reordering the
// outer environment dependencies. This is particularly useful when you need to change the
// order in which contexts are provided to a nested computation.
//
// Example:
//
//	type OuterConfig struct {
//	    DatabaseURL string
//	}
//	type InnerConfig struct {
//	    APIKey string
//	}
//	type RequestContext struct {
//	    UserID int
//	}
//
//	// Original: takes OuterConfig, returns computation that may produce
//	// another computation depending on InnerConfig
//	original := func(outer OuterConfig) readerioeither.ReaderIOEither[RequestContext, error,
//	    readerreaderioeither.ReaderReaderIOEither[InnerConfig, RequestContext, error, string]] {
//	    return readerioeither.Of[RequestContext, error](
//	        readerreaderioeither.Of[InnerConfig, RequestContext, error]("result"),
//	    )
//	}
//
//	// Sequence swaps InnerConfig and OuterConfig order
//	sequenced := Sequence(original)
//
//	// Now provide InnerConfig first, then OuterConfig
//	result := sequenced(InnerConfig{APIKey: "key"})(OuterConfig{DatabaseURL: "db"})(RequestContext{UserID: 1})()
func Sequence[R1, R2, C, E, A any](ma ReaderReaderIOEither[R2, C, E, ReaderReaderIOEither[R1, C, E, A]]) Kleisli[R2, C, E, R1, A] {
	return readert.Sequence(
		readerioeither.Chain,
		ma,
	)
}

// SequenceReader swaps the order of environment parameters when the inner computation is a pure Reader.
//
// This function is similar to Sequence but specialized for the case where the innermost computation
// is a pure Reader (without IO or error handling) rather than another ReaderReaderIOEither. It takes
// a ReaderReaderIOEither that produces a Reader and returns a Kleisli arrow that reverses the order
// of the outer environment parameters.
//
// Type Parameters:
//   - R1: The first environment type (becomes outermost after sequence)
//   - R2: The second environment type (becomes inner after sequence)
//   - C: The inner context/environment type (for the ReaderIOEither layer)
//   - E: The error type (only present in the ReaderReaderIOEither layer)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderReaderIOEither[R2, C, E, Reader[R1, A]]
//
// Returns:
//   - A Kleisli[R2, C, E, R1, A], which is func(R1) ReaderReaderIOEither[R2, C, E, A]
//
// The function lifts the pure Reader computation into the ReaderIOEither context while
// reordering the environment dependencies.
//
// Example:
//
//	type Config struct {
//	    Multiplier int
//	}
//	type Database struct {
//	    ConnectionString string
//	}
//	type Context struct {
//	    RequestID string
//	}
//
//	// Original: takes Config, may produce a Reader[Database, int]
//	original := func(cfg Config) readerioeither.ReaderIOEither[Context, error, reader.Reader[Database, int]] {
//	    return readerioeither.Of[Context, error](func(db Database) int {
//	        return len(db.ConnectionString) * cfg.Multiplier
//	    })
//	}
//
//	// Sequence to provide Database first, then Config
//	sequenced := SequenceReader(original)
//	result := sequenced(Database{ConnectionString: "localhost"})(Config{Multiplier: 2})(Context{RequestID: "123"})()
func SequenceReader[R1, R2, C, E, A any](ma ReaderReaderIOEither[R2, C, E, Reader[R1, A]]) Kleisli[R2, C, E, R1, A] {
	return readert.SequenceReader(
		readerioeither.Map,
		ma,
	)
}

// SequenceReaderIO swaps the order of environment parameters when the inner computation is a ReaderIO.
//
// This function is specialized for the case where the innermost computation is a ReaderIO
// (with IO effects but no error handling) rather than another ReaderReaderIOEither. It takes
// a ReaderReaderIOEither that produces a ReaderIO and returns a Kleisli arrow that reverses
// the order of the outer environment parameters.
//
// Type Parameters:
//   - R1: The first environment type (becomes outermost after sequence)
//   - R2: The second environment type (becomes inner after sequence)
//   - C: The inner context/environment type (for the ReaderIOEither layer)
//   - E: The error type (only present in the outer ReaderReaderIOEither layer)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderReaderIOEither[R2, C, E, ReaderIO[R1, A]]
//
// Returns:
//   - A Kleisli[R2, C, E, R1, A], which is func(R1) ReaderReaderIOEither[R2, C, E, A]
//
// The function lifts the ReaderIO computation (which has IO effects but no error handling)
// into the ReaderIOEither context while reordering the environment dependencies.
//
// Example:
//
//	type Config struct {
//	    FilePath string
//	}
//	type Logger struct {
//	    Level string
//	}
//	type Context struct {
//	    TraceID string
//	}
//
//	// Original: takes Config, may produce a ReaderIO[Logger, string]
//	original := func(cfg Config) readerioeither.ReaderIOEither[Context, error, readerio.ReaderIO[Logger, string]] {
//	    return readerioeither.Of[Context, error](func(logger Logger) io.IO[string] {
//	        return func() string {
//	            return fmt.Sprintf("[%s] Reading from %s", logger.Level, cfg.FilePath)
//	        }
//	    })
//	}
//
//	// Sequence to provide Logger first, then Config
//	sequenced := SequenceReaderIO(original)
//	result := sequenced(Logger{Level: "INFO"})(Config{FilePath: "/data"})(Context{TraceID: "abc"})()
func SequenceReaderIO[R1, R2, C, E, A any](ma ReaderReaderIOEither[R2, C, E, ReaderIO[R1, A]]) Kleisli[R2, C, E, R1, A] {
	return func(r1 R1) ReaderReaderIOEither[R2, C, E, A] {
		rd := ioeither.ChainIOK[E](readerio.Read[A](r1))
		return func(r2 R2) ReaderIOEither[C, E, A] {
			return F.Pipe1(
				ma(r2),
				reader.Map[C](rd),
			)
		}
	}
}

// Traverse transforms a ReaderReaderIOEither computation by applying a function that produces
// another ReaderReaderIOEither, effectively swapping the order of outer environment parameters.
//
// This function is useful when you have a computation that depends on environment R2 and
// produces a value of type A, and you want to transform it using a function that takes A
// and produces a computation depending on environment R1. The result is a curried function
// that takes R1 first, then R2, and produces a ReaderIOEither[C, E, B].
//
// Type Parameters:
//   - R2: The outer environment type from the original computation
//   - R1: The inner environment type introduced by the transformation
//   - C: The inner context/environment type (for the ReaderIOEither layer)
//   - E: The error type
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Kleisli arrow that transforms A into a ReaderReaderIOEither[R1, C, E, B]
//
// Returns:
//   - A function that takes a ReaderReaderIOEither[R2, C, E, A] and returns a Kleisli[R2, C, E, R1, B],
//     which is func(R1) ReaderReaderIOEither[R2, C, E, B]
//
// The function preserves error handling and IO effects while reordering the environment dependencies.
// This is the generalized version of Sequence that also applies a transformation function.
//
// Example:
//
//	type UserConfig struct {
//	    UserID int
//	}
//	type SystemConfig struct {
//	    SystemID string
//	}
//	type Context struct {
//	    RequestID string
//	}
//
//	// Original computation depending on SystemConfig
//	original := readerreaderioeither.Of[SystemConfig, Context, error](42)
//
//	// Transformation that introduces UserConfig dependency
//	transform := func(n int) readerreaderioeither.ReaderReaderIOEither[UserConfig, Context, error, string] {
//	    return func(userCfg UserConfig) readerioeither.ReaderIOEither[Context, error, string] {
//	        return readerioeither.Of[Context, error](fmt.Sprintf("User %d: %d", userCfg.UserID, n))
//	    }
//	}
//
//	// Apply traverse to swap order and transform
//	traversed := Traverse[SystemConfig, UserConfig, Context, error, int, string](transform)(original)
//
//	// Provide UserConfig first, then SystemConfig
//	result := traversed(UserConfig{UserID: 1})(SystemConfig{SystemID: "sys1"})(Context{RequestID: "req1"})()
func Traverse[R2, R1, C, E, A, B any](
	f Kleisli[R1, C, E, A, B],
) func(ReaderReaderIOEither[R2, C, E, A]) Kleisli[R2, C, E, R1, B] {
	return readert.Traverse[ReaderReaderIOEither[R2, C, E, A]](
		readerioeither.Map,
		readerioeither.Chain,
		f,
	)
}

// TraverseReader transforms a ReaderReaderIOEither computation by applying a Reader-based function,
// effectively introducing a new environment dependency.
//
// This function takes a Reader-based transformation (Kleisli arrow) and returns a function that
// can transform a ReaderReaderIOEither. The result allows you to provide the Reader's environment (R1)
// first, which then produces a ReaderReaderIOEither that depends on environment R2.
//
// Type Parameters:
//   - R2: The outer environment type from the original ReaderReaderIOEither
//   - R1: The inner environment type introduced by the Reader transformation
//   - C: The inner context/environment type (for the ReaderIOEither layer)
//   - E: The error type
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Reader-based Kleisli arrow that transforms A to B using environment R1
//
// Returns:
//   - A function that takes a ReaderReaderIOEither[R2, C, E, A] and returns a Kleisli[R2, C, E, R1, B],
//     which is func(R1) ReaderReaderIOEither[R2, C, E, B]
//
// The function preserves error handling and IO effects while adding the Reader environment dependency
// and reordering the environment parameters. This is useful when you want to introduce a pure
// (non-IO, non-error) environment dependency to an existing computation.
//
// Example:
//
//	type SystemConfig struct {
//	    Timeout int
//	}
//	type UserPreferences struct {
//	    Theme string
//	}
//	type Context struct {
//	    SessionID string
//	}
//
//	// Original computation depending on SystemConfig
//	original := readerreaderioeither.Of[SystemConfig, Context, error](100)
//
//	// Pure Reader transformation that introduces UserPreferences dependency
//	formatWithTheme := func(value int) reader.Reader[UserPreferences, string] {
//	    return func(prefs UserPreferences) string {
//	        return fmt.Sprintf("[%s theme] Value: %d", prefs.Theme, value)
//	    }
//	}
//
//	// Apply traverse to introduce UserPreferences and swap order
//	traversed := TraverseReader[SystemConfig, UserPreferences, Context, error, int, string](formatWithTheme)(original)
//
//	// Provide UserPreferences first, then SystemConfig
//	result := traversed(UserPreferences{Theme: "dark"})(SystemConfig{Timeout: 30})(Context{SessionID: "sess1"})()
func TraverseReader[R2, R1, C, E, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderReaderIOEither[R2, C, E, A]) Kleisli[R2, C, E, R1, B] {
	return readert.TraverseReader[ReaderReaderIOEither[R2, C, E, A]](
		readerioeither.Map,
		readerioeither.Map,
		f,
	)
}
