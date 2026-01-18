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

package readerioresult_test

import (
	"context"
	"fmt"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

// Example_sequenceReader_basicUsage demonstrates the basic usage of SequenceReader
// to flip the parameter order, enabling point-free style programming.
func Example_sequenceReader_basicUsage() {
	type Config struct {
		Multiplier int
	}

	// A computation that produces a Reader based on context
	getComputation := func(ctx context.Context) func() either.Either[error, func(Config) int] {
		return func() either.Either[error, func(Config) int] {
			// This could check context for cancellation, deadlines, etc.
			return either.Right[error](func(cfg Config) int {
				return cfg.Multiplier * 10
			})
		}
	}

	// Sequence it to flip the parameter order
	// Now Config comes first, then context
	sequenced := RIOE.SequenceReader(getComputation)

	// Partially apply the Config - this is the key benefit for point-free style
	cfg := Config{Multiplier: 5}
	withConfig := sequenced(cfg)

	// Now we have a ReaderIOResult[int] that can be used with any context
	ctx := context.Background()
	result := withConfig(ctx)()

	if value, err := either.Unwrap(result); err == nil {
		fmt.Println(value)
	}
	// Output: 50
}

// Example_sequenceReader_dependencyInjection demonstrates how SequenceReader
// enables clean dependency injection patterns in point-free style.
func Example_sequenceReader_dependencyInjection() {
	// Define our dependencies
	type Database struct {
		ConnectionString string
	}

	type UserService struct {
		db Database
	}

	// A function that creates a computation requiring a Database
	makeQuery := func(ctx context.Context) func() either.Either[error, func(Database) string] {
		return func() either.Either[error, func(Database) string] {
			return either.Right[error](func(db Database) string {
				return fmt.Sprintf("Querying %s", db.ConnectionString)
			})
		}
	}

	// Sequence to enable dependency injection
	queryWithDB := RIOE.SequenceReader(makeQuery)

	// Inject the database dependency
	db := Database{ConnectionString: "localhost:5432"}
	query := queryWithDB(db)

	// Execute with context
	ctx := context.Background()
	result := query(ctx)()

	if value, err := either.Unwrap(result); err == nil {
		fmt.Println(value)
	}
	// Output: Querying localhost:5432
}

// Example_sequenceReader_pointFreeComposition demonstrates how SequenceReader
// enables point-free style composition of computations.
func Example_sequenceReader_pointFreeComposition() {
	type Config struct {
		BaseValue int
	}

	// Step 1: Create a computation that produces a Reader
	step1 := func(ctx context.Context) func() either.Either[error, func(Config) int] {
		return func() either.Either[error, func(Config) int] {
			return either.Right[error](func(cfg Config) int {
				return cfg.BaseValue * 2
			})
		}
	}

	// Step 2: Sequence it to enable partial application
	sequenced := RIOE.SequenceReader(step1)

	// Step 3: Build a pipeline using point-free style
	// Partially apply the config
	cfg := Config{BaseValue: 10}

	// Create a reusable computation with the config baked in
	computation := F.Pipe1(
		sequenced(cfg),
		RIOE.Map(func(x int) int { return x + 5 }),
	)

	// Execute the pipeline
	ctx := context.Background()
	result := computation(ctx)()

	if value, err := either.Unwrap(result); err == nil {
		fmt.Println(value)
	}
	// Output: 25
}

// Example_sequenceReader_multipleEnvironments demonstrates using SequenceReader
// to work with multiple environment types in a clean, composable way.
func Example_sequenceReader_multipleEnvironments() {
	type DatabaseConfig struct {
		Host string
		Port int
	}

	type APIConfig struct {
		Endpoint string
		APIKey   string
	}

	// Function that needs DatabaseConfig
	getDatabaseURL := func(ctx context.Context) func() either.Either[error, func(DatabaseConfig) string] {
		return func() either.Either[error, func(DatabaseConfig) string] {
			return either.Right[error](func(cfg DatabaseConfig) string {
				return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
			})
		}
	}

	// Function that needs APIConfig
	getAPIURL := func(ctx context.Context) func() either.Either[error, func(APIConfig) string] {
		return func() either.Either[error, func(APIConfig) string] {
			return either.Right[error](func(cfg APIConfig) string {
				return cfg.Endpoint
			})
		}
	}

	// Sequence both to enable partial application
	withDBConfig := RIOE.SequenceReader(getDatabaseURL)
	withAPIConfig := RIOE.SequenceReader(getAPIURL)

	// Partially apply different configs
	dbCfg := DatabaseConfig{Host: "localhost", Port: 5432}
	apiCfg := APIConfig{Endpoint: "https://api.example.com", APIKey: "secret"}

	dbQuery := withDBConfig(dbCfg)
	apiQuery := withAPIConfig(apiCfg)

	// Execute both with the same context
	ctx := context.Background()

	dbResult := dbQuery(ctx)()
	apiResult := apiQuery(ctx)()

	if dbURL, err := either.Unwrap(dbResult); err == nil {
		fmt.Println("Database:", dbURL)
	}
	if apiURL, err := either.Unwrap(apiResult); err == nil {
		fmt.Println("API:", apiURL)
	}
	// Output:
	// Database: localhost:5432
	// API: https://api.example.com
}

// Example_sequenceReaderResult_errorHandling demonstrates how SequenceReaderResult
// enables point-free style with proper error handling at multiple levels.
func Example_sequenceReaderResult_errorHandling() {
	type ValidationConfig struct {
		MinValue int
		MaxValue int
	}

	// A computation that can fail at both outer and inner levels
	makeValidator := func(ctx context.Context) func() either.Either[error, func(context.Context) either.Either[error, int]] {
		return func() either.Either[error, func(context.Context) either.Either[error, int]] {
			// Outer level: check context
			if ctx.Err() != nil {
				return either.Left[func(context.Context) either.Either[error, int]](ctx.Err())
			}

			// Return inner computation
			return either.Right[error](func(innerCtx context.Context) either.Either[error, int] {
				// Inner level: perform validation
				value := 42
				if value < 0 {
					return either.Left[int](fmt.Errorf("value too small: %d", value))
				}
				if value > 100 {
					return either.Left[int](fmt.Errorf("value too large: %d", value))
				}
				return either.Right[error](value)
			})
		}
	}

	// Sequence to enable point-free composition
	sequenced := RIOE.SequenceReaderResult(makeValidator)

	// Build a pipeline with error handling
	ctx := context.Background()
	pipeline := F.Pipe2(
		sequenced(ctx),
		RIOE.Map(N.Mul(2)),
		RIOE.Chain(func(x int) RIOE.ReaderIOResult[string] {
			return RIOE.Of(fmt.Sprintf("Result: %d", x))
		}),
	)

	result := pipeline(ctx)()

	if value, err := either.Unwrap(result); err == nil {
		fmt.Println(value)
	}
	// Output: Result: 84
}

// Example_sequenceReader_partialApplication demonstrates the power of partial
// application enabled by SequenceReader for building reusable computations.
func Example_sequenceReader_partialApplication() {
	type ServiceConfig struct {
		ServiceName string
		Version     string
	}

	// Create a computation factory
	makeServiceInfo := func(ctx context.Context) func() either.Either[error, func(ServiceConfig) string] {
		return func() either.Either[error, func(ServiceConfig) string] {
			return either.Right[error](func(cfg ServiceConfig) string {
				return fmt.Sprintf("%s v%s", cfg.ServiceName, cfg.Version)
			})
		}
	}

	// Sequence it
	sequenced := RIOE.SequenceReader(makeServiceInfo)

	// Create multiple service configurations
	authConfig := ServiceConfig{ServiceName: "AuthService", Version: "1.0.0"}
	userConfig := ServiceConfig{ServiceName: "UserService", Version: "2.1.0"}

	// Partially apply each config to create specialized computations
	getAuthInfo := sequenced(authConfig)
	getUserInfo := sequenced(userConfig)

	// These can now be reused across different contexts
	ctx := context.Background()

	authResult := getAuthInfo(ctx)()
	userResult := getUserInfo(ctx)()

	if auth, err := either.Unwrap(authResult); err == nil {
		fmt.Println(auth)
	}
	if user, err := either.Unwrap(userResult); err == nil {
		fmt.Println(user)
	}
	// Output:
	// AuthService v1.0.0
	// UserService v2.1.0
}

// Example_sequenceReader_testingBenefits demonstrates how SequenceReader
// makes testing easier by allowing you to inject test dependencies.
func Example_sequenceReader_testingBenefits() {
	// Simple logger that collects messages
	type SimpleLogger struct {
		Messages []string
	}

	// A computation that depends on a logger (using the struct directly)
	makeLoggingOperation := func(ctx context.Context) func() either.Either[error, func(*SimpleLogger) string] {
		return func() either.Either[error, func(*SimpleLogger) string] {
			return either.Right[error](func(logger *SimpleLogger) string {
				logger.Messages = append(logger.Messages, "Operation started")
				result := "Success"
				logger.Messages = append(logger.Messages, fmt.Sprintf("Operation completed: %s", result))
				return result
			})
		}
	}

	// Sequence to enable dependency injection
	sequenced := RIOE.SequenceReader(makeLoggingOperation)

	// Inject a test logger
	testLogger := &SimpleLogger{Messages: []string{}}
	operation := sequenced(testLogger)

	// Execute
	ctx := context.Background()
	result := operation(ctx)()

	if value, err := either.Unwrap(result); err == nil {
		fmt.Println("Result:", value)
		fmt.Println("Logs:", len(testLogger.Messages))
	}
	// Output:
	// Result: Success
	// Logs: 2
}
