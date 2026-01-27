// Copyright (c) 2025 IBM Corp.
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
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

// Promap is the profunctor map operation that transforms both the input and output of a ReaderIOResult.
// It applies f to the input environment (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Adapt the environment type before passing it to the ReaderIOResult (via f)
//   - Transform the success value after the IO effect completes (via g)
//
// The error type is fixed as error and remains unchanged through the transformation.
//
// Type Parameters:
//   - R: The original environment type expected by the ReaderIOResult
//   - A: The original success type produced by the ReaderIOResult
//   - D: The new input environment type
//   - B: The new output success type
//
// Parameters:
//   - f: Function to transform the input environment from D to R (contravariant)
//   - g: Function to transform the output success value from A to B (covariant)
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIOResult[R, A] and returns a ReaderIOResult[D, B]
//
//go:inline
func Promap[R, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, ReaderIOResult[R, A], B] {
	return RIOE.Promap[R, error](f, g)
}

// Contramap changes the value of the local environment during the execution of a ReaderIOResult.
// This is the contravariant functor operation that transforms the input environment.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Contramap is useful for adapting a ReaderIOResult to work with a different environment type
// by providing a function that converts the new environment to the expected one.
//
// Type Parameters:
//   - A: The success type (unchanged)
//   - R2: The new input environment type
//   - R1: The original environment type expected by the ReaderIOResult
//
// Parameters:
//   - f: Function to transform the environment from R2 to R1
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIOResult[R1, A] and returns a ReaderIOResult[R2, A]
//
//go:inline
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderIOResult[R1, A], A] {
	return RIOE.Contramap[error, A](f)
}

// LocalIOK transforms the environment of a ReaderIOResult using an IO-based Kleisli arrow.
// It allows you to modify the environment through an effectful computation before
// passing it to the ReaderIOResult.
//
// This is useful when the environment transformation itself requires IO effects,
// such as reading from a file, making a network call, or accessing system resources,
// but these effects cannot fail (or failures are not relevant).
//
// The transformation happens in two stages:
//  1. The IO effect f is executed with the R2 environment to produce an R1 value
//  2. The resulting R1 value is passed to the ReaderIOResult[R1, A] to produce the final result
//
// Type Parameters:
//   - A: The success type produced by the ReaderIOResult
//   - R1: The original environment type expected by the ReaderIOResult
//   - R2: The new input environment type
//
// Parameters:
//   - f: An IO Kleisli arrow that transforms R2 to R1 with IO effects
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIOResult[R1, A] and returns a ReaderIOResult[R2, A]
//
// Example:
//
//	// Transform a config path into a loaded config (infallible)
//	loadConfig := func(path string) IO[Config] {
//	    return func() Config {
//	        return getDefaultConfig() // Always succeeds
//	    }
//	}
//
//	// Use the config to perform an operation that might fail
//	useConfig := func(cfg Config) IOResult[string] {
//	    return func() Result[string] {
//	        if cfg.Valid {
//	            return Ok[string]("Success")
//	        }
//	        return Err[string](errors.New("invalid config"))
//	    }
//	}
//
//	// Compose them using LocalIOK
//	result := LocalIOK[string, Config, string](loadConfig)(useConfig)
//	output := result("config.json")() // Loads config and uses it
//
//go:inline
func LocalIOK[A, R1, R2 any](f io.Kleisli[R2, R1]) Kleisli[R2, ReaderIOResult[R1, A], A] {
	return RIOE.LocalIOK[error, A](f)
}

// LocalIOEitherK transforms the environment of a ReaderIOResult using an IOEither-based Kleisli arrow.
// It allows you to modify the environment through an effectful computation that can fail before
// passing it to the ReaderIOResult.
//
// This is useful when the environment transformation itself requires IO effects that can fail,
// such as reading from a file that might not exist, making a network call that might timeout,
// or parsing data that might be invalid.
//
// The transformation happens in two stages:
//  1. The IOEither effect f is executed with the R2 environment to produce Either[error, R1]
//  2. If successful (Right), the R1 value is passed to the ReaderIOResult[R1, A]
//  3. If failed (Left), the error is propagated without executing the ReaderIOResult
//
// Type Parameters:
//   - A: The success type produced by the ReaderIOResult
//   - R1: The original environment type expected by the ReaderIOResult
//   - R2: The new input environment type
//
// Parameters:
//   - f: An IOEither Kleisli arrow that transforms R2 to R1 with IO effects that can fail
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIOResult[R1, A] and returns a ReaderIOResult[R2, A]
//
// Example:
//
//	// Transform a config path into a loaded config (can fail)
//	loadConfig := func(path string) IOEither[error, Config] {
//	    return func() Either[error, Config] {
//	        cfg, err := readConfigFile(path)
//	        if err != nil {
//	            return Left[Config](err)
//	        }
//	        return Right[error](cfg)
//	    }
//	}
//
//	// Use the config to perform an operation that might fail
//	useConfig := func(cfg Config) IOResult[string] {
//	    return func() Result[string] {
//	        if cfg.Valid {
//	            return Ok[string]("Success: " + cfg.Name)
//	        }
//	        return Err[string](errors.New("invalid config"))
//	    }
//	}
//
//	// Compose them using LocalIOEitherK
//	result := LocalIOEitherK[string, Config, string](loadConfig)(useConfig)
//	output := result("config.json")() // Loads config (might fail) and uses it (might fail)
//
//go:inline
func LocalIOEitherK[A, R1, R2 any](f ioeither.Kleisli[error, R2, R1]) Kleisli[R2, ReaderIOResult[R1, A], A] {
	return RIOE.LocalIOEitherK[A](f)
}

// LocalIOResultK transforms the environment of a ReaderIOResult using an IOResult-based Kleisli arrow.
// It allows you to modify the environment through an effectful computation that can fail before
// passing it to the ReaderIOResult.
//
// This is a type-safe alias for LocalIOEitherK specialized for error type, providing a more
// idiomatic API when working with Result types (which use error as the error type).
//
// The transformation happens in two stages:
//  1. The IOResult effect f is executed with the R2 environment to produce Result[R1]
//  2. If successful (Ok), the R1 value is passed to the ReaderIOResult[R1, A]
//  3. If failed (Err), the error is propagated without executing the ReaderIOResult
//
// Type Parameters:
//   - A: The success type produced by the ReaderIOResult
//   - R1: The original environment type expected by the ReaderIOResult
//   - R2: The new input environment type
//
// Parameters:
//   - f: An IOResult Kleisli arrow that transforms R2 to R1 with IO effects that can fail
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIOResult[R1, A] and returns a ReaderIOResult[R2, A]
//
// Example:
//
//	// Transform a config path into a loaded config (can fail)
//	loadConfig := func(path string) IOResult[Config] {
//	    return func() Result[Config] {
//	        cfg, err := readConfigFile(path)
//	        if err != nil {
//	            return Err[Config](err)
//	        }
//	        return Ok(cfg)
//	    }
//	}
//
//	// Use the config to perform an operation that might fail
//	useConfig := func(cfg Config) IOResult[string] {
//	    return func() Result[string] {
//	        if cfg.Valid {
//	            return Ok("Success: " + cfg.Name)
//	        }
//	        return Err[string](errors.New("invalid config"))
//	    }
//	}
//
//	// Compose them using LocalIOResultK
//	result := LocalIOResultK[string, Config, string](loadConfig)(useConfig)
//	output := result("config.json")() // Loads config (might fail) and uses it (might fail)
//
//go:inline
func LocalIOResultK[A, R1, R2 any](f ioresult.Kleisli[R2, R1]) Kleisli[R2, ReaderIOResult[R1, A], A] {
	return RIOE.LocalIOEitherK[A](f)
}
