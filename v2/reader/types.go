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

package reader

type (
	// Reader represents a computation that depends on a shared environment of type R and produces a value of type A.
	//
	// The purpose of the Reader monad is to avoid threading arguments through multiple functions
	// in order to only get them where they are needed. This enables dependency injection and
	// configuration management in a functional style.
	//
	// Type Parameters:
	//   - R: The environment/context type (read-only, shared across computations)
	//   - A: The result type produced by the computation
	//
	// A Reader[R, A] is simply a function from R to A: func(R) A
	//
	// Example:
	//
	//	type Config struct {
	//	    DatabaseURL string
	//	    APIKey      string
	//	}
	//
	//	// A Reader that extracts the database URL from config
	//	getDatabaseURL := func(c Config) string { return c.DatabaseURL }
	//
	//	// A Reader that extracts the API key from config
	//	getAPIKey := func(c Config) string { return c.APIKey }
	//
	//	// Use the readers with a config
	//	config := Config{DatabaseURL: "localhost:5432", APIKey: "secret"}
	//	dbURL := getDatabaseURL(config)  // "localhost:5432"
	//	apiKey := getAPIKey(config)      // "secret"
	Reader[R, A any] = func(R) A

	// Operator represents a transformation from one Reader to another.
	// It takes a Reader[R, A] and produces a Reader[R, B], where both readers
	// share the same environment type R.
	//
	// This type is commonly used for operations like Map, Chain, and other
	// transformations that convert readers while preserving the environment type.
	//
	// Type Parameters:
	//   - R: The shared environment/context type
	//   - A: The input Reader's result type
	//   - B: The output Reader's result type
	//
	// Example:
	//
	//	type Config struct { Multiplier int }
	//
	//	// An operator that transforms int readers to string readers
	//	intToString := reader.Map[Config, int, string](strconv.Itoa)
	//
	//	getNumber := reader.Asks(func(c Config) int { return c.Multiplier })
	//	getString := intToString(getNumber)
	//	result := getString(Config{Multiplier: 42}) // "42"
	Operator[R, A, B any] = func(Reader[R, A]) Reader[R, B]
)
