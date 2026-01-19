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

package readerreaderioresult

import (
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerioeither"
	RRIOE "github.com/IBM/fp-go/v2/readerreaderioeither"
)

// Sequence swaps the order of nested environment parameters in a ReaderReaderIOResult computation.
//
// This function takes a ReaderReaderIOResult that produces another ReaderReaderIOResult and returns a
// Kleisli arrow that reverses the order of the outer environment parameters (R1 and R2). The result is
// a curried function that takes R1 first, then R2, and produces a computation with context.Context and error handling.
//
// Type Parameters:
//   - R1: The first outer environment type (becomes the outermost after sequence)
//   - R2: The second outer environment type (becomes inner after sequence)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderReaderIOResult[R2, ReaderReaderIOResult[R1, A]]
//
// Returns:
//   - A Kleisli[R2, R1, A], which is func(R1) ReaderReaderIOResult[R2, A]
//
// The function preserves error handling and IO effects at all levels while reordering the
// outer environment dependencies. The inner context.Context layer remains unchanged.
//
// This is particularly useful when you need to change the order in which contexts are provided
// to a nested computation, such as when composing operations that have different dependency orders.
//
// Example:
//
//	type AppConfig struct {
//	    DatabaseURL string
//	}
//	type UserPrefs struct {
//	    Theme string
//	}
//
//	// Original: takes AppConfig, returns computation that may produce
//	// another computation depending on UserPrefs
//	original := func(cfg AppConfig) readerioresult.ReaderIOResult[context.Context,
//	    ReaderReaderIOResult[UserPrefs, string]] {
//	    return readerioresult.Of[context.Context](
//	        Of[UserPrefs]("result"),
//	    )
//	}
//
//	// Sequence swaps UserPrefs and AppConfig order
//	sequenced := Sequence[UserPrefs, AppConfig, string](original)
//
//	// Now provide UserPrefs first, then AppConfig
//	ctx := context.Background()
//	result := sequenced(UserPrefs{Theme: "dark"})(AppConfig{DatabaseURL: "db"})(ctx)()
func Sequence[R1, R2, A any](ma ReaderReaderIOResult[R2, ReaderReaderIOResult[R1, A]]) Kleisli[R2, R1, A] {
	return readert.Sequence(
		readerioeither.Chain,
		ma,
	)
}

// SequenceReader swaps the order of environment parameters when the inner computation is a pure Reader.
//
// This function is similar to Sequence but specialized for the case where the innermost computation
// is a pure Reader (without IO or error handling) rather than another ReaderReaderIOResult. It takes
// a ReaderReaderIOResult that produces a Reader and returns a Kleisli arrow that reverses the order
// of the outer environment parameters.
//
// Type Parameters:
//   - R1: The first environment type (becomes outermost after sequence)
//   - R2: The second environment type (becomes inner after sequence)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderReaderIOResult[R2, Reader[R1, A]]
//
// Returns:
//   - A Kleisli[R2, R1, A], which is func(R1) ReaderReaderIOResult[R2, A]
//
// The function lifts the pure Reader computation into the ReaderIOResult context (with context.Context
// and error handling) while reordering the environment dependencies.
//
// Example:
//
//	type AppConfig struct {
//	    Multiplier int
//	}
//	type Database struct {
//	    ConnectionString string
//	}
//
//	// Original: takes AppConfig, may produce a Reader[Database, int]
//	original := func(cfg AppConfig) readerioresult.ReaderIOResult[context.Context, reader.Reader[Database, int]] {
//	    return readerioresult.Of[context.Context](func(db Database) int {
//	        return len(db.ConnectionString) * cfg.Multiplier
//	    })
//	}
//
//	// Sequence to provide Database first, then AppConfig
//	sequenced := SequenceReader[Database, AppConfig, int](original)
//	ctx := context.Background()
//	result := sequenced(Database{ConnectionString: "localhost"})(AppConfig{Multiplier: 2})(ctx)()
func SequenceReader[R1, R2, A any](ma ReaderReaderIOResult[R2, Reader[R1, A]]) Kleisli[R2, R1, A] {
	return readert.SequenceReader(
		readerioeither.Map,
		ma,
	)
}

// SequenceReaderIO swaps the order of environment parameters when the inner computation is a ReaderIO.
//
// This function is specialized for the case where the innermost computation is a ReaderIO
// (with IO effects but no error handling) rather than another ReaderReaderIOResult. It takes
// a ReaderReaderIOResult that produces a ReaderIO and returns a Kleisli arrow that reverses
// the order of the outer environment parameters.
//
// Type Parameters:
//   - R1: The first environment type (becomes outermost after sequence)
//   - R2: The second environment type (becomes inner after sequence)
//   - A: The success value type
//
// Parameters:
//   - ma: A ReaderReaderIOResult[R2, ReaderIO[R1, A]]
//
// Returns:
//   - A Kleisli[R2, R1, A], which is func(R1) ReaderReaderIOResult[R2, A]
//
// The function lifts the ReaderIO computation (which has IO effects but no error handling)
// into the ReaderIOResult context (with context.Context and error handling) while reordering
// the environment dependencies.
//
// Example:
//
//	type AppConfig struct {
//	    FilePath string
//	}
//	type Logger struct {
//	    Level string
//	}
//
//	// Original: takes AppConfig, may produce a ReaderIO[Logger, string]
//	original := func(cfg AppConfig) readerioresult.ReaderIOResult[context.Context, readerio.ReaderIO[Logger, string]] {
//	    return readerioresult.Of[context.Context](func(logger Logger) io.IO[string] {
//	        return func() string {
//	            return fmt.Sprintf("[%s] Reading from %s", logger.Level, cfg.FilePath)
//	        }
//	    })
//	}
//
//	// Sequence to provide Logger first, then AppConfig
//	sequenced := SequenceReaderIO[Logger, AppConfig, string](original)
//	ctx := context.Background()
//	result := sequenced(Logger{Level: "INFO"})(AppConfig{FilePath: "/data"})(ctx)()
func SequenceReaderIO[R1, R2, A any](ma ReaderReaderIOResult[R2, ReaderIO[R1, A]]) Kleisli[R2, R1, A] {
	return RRIOE.SequenceReaderIO(ma)
}

// Traverse transforms a ReaderReaderIOResult computation by applying a function that produces
// another ReaderReaderIOResult, effectively swapping the order of outer environment parameters.
//
// This function is useful when you have a computation that depends on environment R2 and
// produces a value of type A, and you want to transform it using a function that takes A
// and produces a computation depending on environment R1. The result is a curried function
// that takes R1 first, then R2, and produces a computation with context.Context and error handling.
//
// Type Parameters:
//   - R2: The outer environment type from the original computation
//   - R1: The inner environment type introduced by the transformation
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Kleisli arrow that transforms A into a ReaderReaderIOResult[R1, B]
//
// Returns:
//   - A function that takes a ReaderReaderIOResult[R2, A] and returns a Kleisli[R2, R1, B],
//     which is func(R1) ReaderReaderIOResult[R2, B]
//
// The function preserves error handling and IO effects while reordering the environment dependencies.
// This is the generalized version of Sequence that also applies a transformation function.
//
// Example:
//
//	type AppConfig struct {
//	    SystemID string
//	}
//	type UserConfig struct {
//	    UserID int
//	}
//
//	// Original computation depending on AppConfig
//	original := Of[AppConfig](42)
//
//	// Transformation that introduces UserConfig dependency
//	transform := func(n int) ReaderReaderIOResult[UserConfig, string] {
//	    return func(userCfg UserConfig) readerioresult.ReaderIOResult[context.Context, string] {
//	        return readerioresult.Of[context.Context](fmt.Sprintf("User %d: %d", userCfg.UserID, n))
//	    }
//	}
//
//	// Apply traverse to swap order and transform
//	traversed := Traverse[AppConfig, UserConfig, int, string](transform)(original)
//
//	// Provide UserConfig first, then AppConfig
//	ctx := context.Background()
//	result := traversed(UserConfig{UserID: 1})(AppConfig{SystemID: "sys1"})(ctx)()
func Traverse[R2, R1, A, B any](
	f Kleisli[R1, A, B],
) func(ReaderReaderIOResult[R2, A]) Kleisli[R2, R1, B] {
	return readert.Traverse[ReaderReaderIOResult[R2, A]](
		readerioeither.Map,
		readerioeither.Chain,
		f,
	)
}

// TraverseReader transforms a ReaderReaderIOResult computation by applying a Reader-based function,
// effectively introducing a new environment dependency.
//
// This function takes a Reader-based transformation (Kleisli arrow) and returns a function that
// can transform a ReaderReaderIOResult. The result allows you to provide the Reader's environment (R1)
// first, which then produces a ReaderReaderIOResult that depends on environment R2.
//
// Type Parameters:
//   - R2: The outer environment type from the original ReaderReaderIOResult
//   - R1: The inner environment type introduced by the Reader transformation
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Reader-based Kleisli arrow that transforms A to B using environment R1
//
// Returns:
//   - A function that takes a ReaderReaderIOResult[R2, A] and returns a Kleisli[R2, R1, B],
//     which is func(R1) ReaderReaderIOResult[R2, B]
//
// The function preserves error handling and IO effects while adding the Reader environment dependency
// and reordering the environment parameters. This is useful when you want to introduce a pure
// (non-IO, non-error) environment dependency to an existing computation.
//
// Example:
//
//	type AppConfig struct {
//	    Timeout int
//	}
//	type UserPreferences struct {
//	    Theme string
//	}
//
//	// Original computation depending on AppConfig
//	original := Of[AppConfig](100)
//
//	// Pure Reader transformation that introduces UserPreferences dependency
//	formatWithTheme := func(value int) reader.Reader[UserPreferences, string] {
//	    return func(prefs UserPreferences) string {
//	        return fmt.Sprintf("[%s theme] Value: %d", prefs.Theme, value)
//	    }
//	}
//
//	// Apply traverse to introduce UserPreferences and swap order
//	traversed := TraverseReader[AppConfig, UserPreferences, int, string](formatWithTheme)(original)
//
//	// Provide UserPreferences first, then AppConfig
//	ctx := context.Background()
//	result := traversed(UserPreferences{Theme: "dark"})(AppConfig{Timeout: 30})(ctx)()
func TraverseReader[R2, R1, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderReaderIOResult[R2, A]) Kleisli[R2, R1, B] {
	return readert.TraverseReader[ReaderReaderIOResult[R2, A]](
		readerioeither.Map,
		readerioeither.Map,
		f,
	)
}
