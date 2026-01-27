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

package readerio

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
)

// Promap is the profunctor map operation that transforms both the input and output of a ReaderIO.
// It applies f to the input environment (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Adapt the environment type before passing it to the ReaderIO (via f)
//   - Transform the result value after the IO effect completes (via g)
//
// The transformation happens in two stages:
//  1. The input environment D is transformed to E using f (contravariant)
//  2. The ReaderIO[E, A] is executed with the transformed environment
//  3. The result value A is transformed to B using g (covariant) within the IO context
//
// Type Parameters:
//   - E: The original environment type expected by the ReaderIO
//   - A: The original result type produced by the ReaderIO
//   - D: The new input environment type
//   - B: The new output result type
//
// Parameters:
//   - f: Function to transform the input environment from D to E (contravariant)
//   - g: Function to transform the output value from A to B (covariant)
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIO[E, A] and returns a ReaderIO[D, B]
//
// Example - Adapting environment and transforming result:
//
//	type DetailedConfig struct {
//	    Host string
//	    Port int
//	    Debug bool
//	}
//
//	type SimpleConfig struct {
//	    Host string
//	    Port int
//	}
//
//	// ReaderIO that reads port from SimpleConfig
//	getPort := readerio.Asks(func(c SimpleConfig) io.IO[int] {
//	    return io.Of(c.Port)
//	})
//
//	// Adapt DetailedConfig to SimpleConfig and convert int to string
//	simplify := func(d DetailedConfig) SimpleConfig {
//	    return SimpleConfig{Host: d.Host, Port: d.Port}
//	}
//	toString := strconv.Itoa
//
//	adapted := readerio.Promap(simplify, toString)(getPort)
//	result := adapted(DetailedConfig{Host: "localhost", Port: 8080, Debug: true})()
//	// result = "8080"
//
// Example - Logging with environment transformation:
//
//	type AppEnv struct {
//	    Logger *log.Logger
//	    Config Config
//	}
//
//	type LoggerEnv struct {
//	    Logger *log.Logger
//	}
//
//	logMessage := func(msg string) readerio.ReaderIO[LoggerEnv, func()] {
//	    return readerio.Asks(func(env LoggerEnv) io.IO[func()] {
//	        return io.Of(func() { env.Logger.Println(msg) })
//	    })
//	}
//
//	extractLogger := func(app AppEnv) LoggerEnv {
//	    return LoggerEnv{Logger: app.Logger}
//	}
//	ignore := func(func()) string { return "logged" }
//
//	logAndReturn := readerio.Promap(extractLogger, ignore)(logMessage("Hello"))
//	// Now works with AppEnv and returns string instead of func()
//
//go:inline
func Promap[E, A, D, B any](f func(D) E, g func(A) B) Kleisli[D, ReaderIO[E, A], B] {
	return reader.Promap(f, io.Map(g))
}

// Local changes the value of the local environment during the execution of a ReaderIO.
// This allows you to modify or adapt the environment before passing it to a ReaderIO computation.
//
// Local is particularly useful for:
//   - Extracting a subset of a larger environment
//   - Transforming environment types
//   - Providing different views of the same environment to different computations
//
// The transformation is contravariant - it transforms the input environment before
// the ReaderIO computation sees it, but doesn't affect the output value.
//
// Type Parameters:
//   - A: The result type produced by the ReaderIO
//   - R2: The new input environment type
//   - R1: The original environment type expected by the ReaderIO
//
// Parameters:
//   - f: Function to transform the environment from R2 to R1
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIO[R1, A] and returns a ReaderIO[R2, A]
//
// Example - Extracting a subset of environment:
//
//	type AppConfig struct {
//	    Database DatabaseConfig
//	    Server   ServerConfig
//	    Logger   *log.Logger
//	}
//
//	type DatabaseConfig struct {
//	    Host string
//	    Port int
//	}
//
//	// ReaderIO that only needs DatabaseConfig
//	connectDB := readerio.Asks(func(cfg DatabaseConfig) io.IO[string] {
//	    return io.Of(fmt.Sprintf("Connected to %s:%d", cfg.Host, cfg.Port))
//	})
//
//	// Extract database config from full app config
//	extractDB := func(app AppConfig) DatabaseConfig {
//	    return app.Database
//	}
//
//	// Adapt to work with full AppConfig
//	connectWithAppConfig := readerio.Local(extractDB)(connectDB)
//	result := connectWithAppConfig(AppConfig{
//	    Database: DatabaseConfig{Host: "localhost", Port: 5432},
//	})()
//	// result = "Connected to localhost:5432"
//
// Example - Providing different views:
//
//	type FullEnv struct {
//	    UserID int
//	    Role   string
//	}
//
//	type UserEnv struct {
//	    UserID int
//	}
//
//	getUserData := readerio.Asks(func(env UserEnv) io.IO[string] {
//	    return io.Of(fmt.Sprintf("User: %d", env.UserID))
//	})
//
//	toUserEnv := func(full FullEnv) UserEnv {
//	    return UserEnv{UserID: full.UserID}
//	}
//
//	adapted := readerio.Local(toUserEnv)(getUserData)
//	result := adapted(FullEnv{UserID: 42, Role: "admin"})()
//	// result = "User: 42"
//
//go:inline
func Local[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderIO[R1, A], A] {
	return reader.Local[IO[A]](f)
}

// Contramap is an alias for Local.
// It changes the value of the local environment during the execution of a ReaderIO.
// This is the contravariant functor operation that transforms the input environment.
//
// Contramap is semantically identical to Local - both modify the environment before
// passing it to a ReaderIO. The name "Contramap" emphasizes the contravariant nature
// of the transformation (transforming the input rather than the output).
//
// Type Parameters:
//   - A: The result type produced by the ReaderIO
//   - R2: The new input environment type
//   - R1: The original environment type expected by the ReaderIO
//
// Parameters:
//   - f: Function to transform the environment from R2 to R1
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIO[R1, A] and returns a ReaderIO[R2, A]
//
// Example - Environment adaptation:
//
//	type DetailedEnv struct {
//	    Config   Config
//	    Logger   *log.Logger
//	    Metrics  Metrics
//	}
//
//	type SimpleEnv struct {
//	    Config Config
//	}
//
//	readConfig := readerio.Asks(func(env SimpleEnv) io.IO[string] {
//	    return io.Of(env.Config.Value)
//	})
//
//	simplify := func(detailed DetailedEnv) SimpleEnv {
//	    return SimpleEnv{Config: detailed.Config}
//	}
//
//	adapted := readerio.Contramap(simplify)(readConfig)
//	result := adapted(DetailedEnv{Config: Config{Value: "test"}})()
//	// result = "test"
//
// See also: Local
//
//go:inline
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderIO[R1, A], A] {
	return reader.Contramap[IO[A]](f)
}

// LocalIOK transforms the environment of a ReaderIO using an IO-based Kleisli arrow.
// It allows you to modify the environment through an effectful computation before
// passing it to the ReaderIO.
//
// This is useful when the environment transformation itself requires IO effects,
// such as reading from a file, making a network call, or accessing system resources.
//
// The transformation happens in two stages:
//  1. The IO effect f is executed with the R2 environment to produce an R1 value
//  2. The resulting R1 value is passed to the ReaderIO[R1, A] to produce the final result
//
// Type Parameters:
//   - A: The result type produced by the ReaderIO
//   - R1: The original environment type expected by the ReaderIO
//   - R2: The new input environment type
//
// Parameters:
//   - f: An IO Kleisli arrow that transforms R2 to R1 with IO effects
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIO[R1, A] and returns a ReaderIO[R2, A]
//
// Example:
//
//	// Transform a config path into a loaded config
//	loadConfig := func(path string) IO[Config] {
//	    return func() Config {
//	        // Load config from file
//	        return parseConfig(readFile(path))
//	    }
//	}
//
//	// Use the config to perform some operation
//	useConfig := func(cfg Config) IO[string] {
//	    return Of("Using: " + cfg.Name)
//	}
//
//	// Compose them using LocalIOK
//	result := LocalIOK[string, Config, string](loadConfig)(useConfig)
//	output := result("config.json")() // Loads config and uses it
//
//go:inline
func LocalIOK[A, R1, R2 any](f io.Kleisli[R2, R1]) Kleisli[R2, ReaderIO[R1, A], A] {
	return func(ri ReaderIO[R1, A]) ReaderIO[R2, A] {
		return F.Flow2(
			f,
			io.Chain(ri),
		)
	}
}
