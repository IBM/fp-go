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

package effect

import (
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
)

// Local transforms the context required by an effect using a pure function.
// This allows you to adapt an effect that requires one context type to work
// with a different context type by providing a transformation function.
//
// # Type Parameters
//
//   - C1: The outer context type (what you have)
//   - C2: The inner context type (what the effect needs)
//   - A: The value type produced by the effect
//
// # Parameters
//
//   - acc: A pure function that transforms C1 to C2
//
// # Returns
//
//   - Kleisli[C1, Effect[C2, A], A]: A function that adapts the effect to use C1
//
// # Example
//
//	type AppConfig struct { DB DatabaseConfig }
//	type DatabaseConfig struct { Host string }
//	dbEffect := effect.Of[DatabaseConfig]("connected")
//	appEffect := effect.Local[AppConfig, DatabaseConfig, string](
//		func(app AppConfig) DatabaseConfig { return app.DB },
//	)(dbEffect)
//
//go:inline
func Local[C1, C2, A any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A] {
	return readerreaderioresult.Local[A](acc)
}

// Contramap is an alias for Local, following the contravariant functor naming convention.
// It transforms the context required by an effect using a pure function.
//
// # Type Parameters
//
//   - C1: The outer context type (what you have)
//   - C2: The inner context type (what the effect needs)
//   - A: The value type produced by the effect
//
// # Parameters
//
//   - acc: A pure function that transforms C1 to C2
//
// # Returns
//
//   - Kleisli[C1, Effect[C2, A], A]: A function that adapts the effect to use C1
//
//go:inline
func Contramap[C1, C2, A any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A] {
	return readerreaderioresult.Local[A](acc)
}

// LocalIOK transforms the context using an IO-based function.
// This allows the context transformation itself to perform I/O operations.
//
// # Type Parameters
//
//   - A: The value type produced by the effect
//   - C1: The inner context type (what the effect needs)
//   - C2: The outer context type (what you have)
//
// # Parameters
//
//   - f: An IO function that transforms C2 to C1
//
// # Returns
//
//   - func(Effect[C1, A]) Effect[C2, A]: A function that adapts the effect
//
// # Example
//
//	loadConfig := func(path string) io.IO[Config] {
//		return func() Config { /* load from file */ }
//	}
//	transform := effect.LocalIOK[string](loadConfig)
//	adapted := transform(configEffect)
//
//go:inline
func LocalIOK[A, C1, C2 any](f io.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalIOK[A](f)
}

// LocalIOResultK transforms the context using an IOResult-based function.
// This allows the context transformation to perform I/O and handle errors.
//
// # Type Parameters
//
//   - A: The value type produced by the effect
//   - C1: The inner context type (what the effect needs)
//   - C2: The outer context type (what you have)
//
// # Parameters
//
//   - f: An IOResult function that transforms C2 to C1
//
// # Returns
//
//   - func(Effect[C1, A]) Effect[C2, A]: A function that adapts the effect
//
// # Example
//
//	loadConfig := func(path string) ioresult.IOResult[Config] {
//		return func() result.Result[Config] {
//			// load from file, may fail
//		}
//	}
//	transform := effect.LocalIOResultK[string](loadConfig)
//	adapted := transform(configEffect)
//
//go:inline
func LocalIOResultK[A, C1, C2 any](f ioresult.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalIOResultK[A](f)
}

// LocalResultK transforms the context using a Result-based function.
// This allows the context transformation to fail with an error.
//
// # Type Parameters
//
//   - A: The value type produced by the effect
//   - C1: The inner context type (what the effect needs)
//   - C2: The outer context type (what you have)
//
// # Parameters
//
//   - f: A Result function that transforms C2 to C1
//
// # Returns
//
//   - func(Effect[C1, A]) Effect[C2, A]: A function that adapts the effect
//
// # Example
//
//	validateConfig := func(raw RawConfig) result.Result[Config] {
//		if raw.IsValid() {
//			return result.Of(raw.ToConfig())
//		}
//		return result.Left[Config](errors.New("invalid"))
//	}
//	transform := effect.LocalResultK[string](validateConfig)
//	adapted := transform(configEffect)
//
//go:inline
func LocalResultK[A, C1, C2 any](f result.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalResultK[A](f)
}

// LocalThunkK transforms the context using a Thunk (ReaderIOResult) function.
// This allows the context transformation to depend on context.Context, perform I/O, and handle errors.
//
// # Type Parameters
//
//   - A: The value type produced by the effect
//   - C1: The inner context type (what the effect needs)
//   - C2: The outer context type (what you have)
//
// # Parameters
//
//   - f: A Thunk function that transforms C2 to C1
//
// # Returns
//
//   - func(Effect[C1, A]) Effect[C2, A]: A function that adapts the effect
//
// # Example
//
//	loadConfig := func(path string) readerioresult.ReaderIOResult[Config] {
//		return func(ctx context.Context) ioresult.IOResult[Config] {
//			// load from file with context, may fail
//		}
//	}
//	transform := effect.LocalThunkK[string](loadConfig)
//	adapted := transform(configEffect)
//
//go:inline
func LocalThunkK[A, C1, C2 any](f thunk.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalReaderIOResultK[A](f)
}

// LocalEffectK transforms the context of an Effect using an Effect-returning function.
// This is the most powerful context transformation function, allowing the transformation
// itself to be effectful (can fail, perform I/O, and access the outer context).
//
// LocalEffectK takes a Kleisli arrow that:
//   - Accepts the outer context C2
//   - Returns an Effect that produces the inner context C1
//   - Can fail with an error during context transformation
//   - Can perform I/O operations during transformation
//
// This is useful when:
//   - Context transformation requires I/O (e.g., loading config from a file)
//   - Context transformation can fail (e.g., validating or parsing context)
//   - Context transformation needs to access the outer context
//
// Type Parameters:
//   - A: The value type produced by the effect
//   - C1: The inner context type (required by the original effect)
//   - C2: The outer context type (provided to the transformed effect)
//
// Parameters:
//   - f: A Kleisli arrow (C2 -> Effect[C2, C1]) that transforms C2 to C1 effectfully
//
// Returns:
//   - A function that transforms Effect[C1, A] to Effect[C2, A]
//
// Example:
//
//	type DatabaseConfig struct {
//		ConnectionString string
//	}
//
//	type AppConfig struct {
//		ConfigPath string
//	}
//
//	// Effect that needs DatabaseConfig
//	dbEffect := effect.Of[DatabaseConfig, string]("query result")
//
//	// Transform AppConfig to DatabaseConfig effectfully
//	// (e.g., load config from file, which can fail)
//	loadConfig := func(app AppConfig) Effect[AppConfig, DatabaseConfig] {
//		return effect.Chain[AppConfig](func(_ AppConfig) Effect[AppConfig, DatabaseConfig] {
//			// Simulate loading config from file (can fail)
//			return effect.Of[AppConfig, DatabaseConfig](DatabaseConfig{
//				ConnectionString: "loaded from " + app.ConfigPath,
//			})
//		})(effect.Of[AppConfig, AppConfig](app))
//	}
//
//	// Apply the transformation
//	transform := effect.LocalEffectK[string, DatabaseConfig, AppConfig](loadConfig)
//	appEffect := transform(dbEffect)
//
//	// Run with AppConfig
//	ioResult := effect.Provide(AppConfig{ConfigPath: "/etc/app.conf"})(appEffect)
//	readerResult := effect.RunSync(ioResult)
//	result, err := readerResult(context.Background())
//
// Comparison with other Local functions:
//   - Local/Contramap: Pure context transformation (C2 -> C1)
//   - LocalIOK: IO-based transformation (C2 -> IO[C1])
//   - LocalIOResultK: IO with error handling (C2 -> IOResult[C1])
//   - LocalReaderIOResultK: Reader-based with IO and errors (C2 -> ReaderIOResult[C1])
//   - LocalEffectK: Full Effect transformation (C2 -> Effect[C2, C1])
//
//go:inline
func LocalEffectK[A, C1, C2 any](f Kleisli[C2, C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalReaderReaderIOEitherK[A](f)
}
