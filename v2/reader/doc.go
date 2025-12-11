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

// Package reader provides the Reader monad implementation for functional programming in Go.
//
// The Reader monad is used to pass a shared environment or configuration through a computation
// without explicitly threading it through every function call. It represents a computation that
// depends on some external context of type R and produces a value of type A.
//
// # Fantasy Land Specification
//
// This implementation corresponds to the Fantasy Land Reader type:
// https://github.com/fantasyland/fantasy-land
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//
// # Core Concept
//
// A Reader[R, A] is simply a function from R to A: func(R) A
// - R is the environment/context type (read-only)
// - A is the result type
//
// # Key Benefits
//
//   - Dependency Injection: Pass configuration or dependencies implicitly
//   - Composition: Combine readers that share the same environment
//   - Testability: Easy to test by providing different environments
//   - Avoid Threading: No need to pass context through every function
//
// # Basic Usage
//
//	// Define a configuration type
//	type Config struct {
//	    Host string
//	    Port int
//	}
//
//	// Create readers that depend on Config
//	getHost := reader.Asks(func(c Config) string { return c.Host })
//	getPort := reader.Asks(func(c Config) int { return c.Port })
//
//	// Compose readers
//	getURL := reader.Map(func(host string) string {
//	    return "http://" + host
//	})(getHost)
//
//	// Run the reader with a config
//	config := Config{Host: "localhost", Port: 8080}
//	url := getURL(config) // "http://localhost"
//
// # Common Operations
//
//   - Ask: Get the current environment
//   - Asks: Project a value from the environment
//   - Map: Transform the result value
//   - Chain: Sequence computations that depend on previous results
//   - Local: Modify the environment for a sub-computation
//
// # Monadic Operations
//
// The Reader type implements the Functor, Applicative, and Monad type classes:
//
//   - Functor: Map over the result value
//   - Applicative: Combine multiple readers with independent computations
//   - Monad: Chain readers where later computations depend on earlier results
//
// # Related Packages
//
//   - reader/generic: Generic implementations for custom reader types
//   - readerio: Reader combined with IO effects
//   - readereither: Reader combined with Either for error handling
//   - readerioeither: Reader combined with IO and Either
package reader

//go:generate go run .. reader --count 10 --filename gen.go
