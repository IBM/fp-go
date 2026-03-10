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
	F "github.com/IBM/fp-go/v2/function"
)

// Promap is the profunctor map operation that transforms both the input and output of an Effect.
// It applies f to the input context (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Modify the context before passing it to the Effect (via f)
//   - Transform the success value after the computation completes (via g)
//
// Promap is particularly useful for adapting effects to work with different context types
// while simultaneously transforming their output values.
//
// # Type Parameters
//
//   - E: The original context type expected by the Effect
//   - A: The original success type produced by the Effect
//   - D: The new input context type
//   - B: The new output success type
//
// # Parameters
//
//   - f: Function to transform the input context from D to E (contravariant)
//   - g: Function to transform the output success value from A to B (covariant)
//
// # Returns
//
//   - A Kleisli arrow that takes an Effect[E, A] and returns a function from D to B
//
// # Example Usage
//
//	type AppConfig struct {
//	    DatabaseURL string
//	    APIKey      string
//	}
//
//	type DBConfig struct {
//	    URL string
//	}
//
//	// Effect that uses DBConfig and returns an int
//	getUserCount := func(cfg DBConfig) effect.Effect[context.Context, int] {
//	    return effect.Succeed[context.Context](42)
//	}
//
//	// Transform AppConfig to DBConfig
//	extractDBConfig := func(app AppConfig) DBConfig {
//	    return DBConfig{URL: app.DatabaseURL}
//	}
//
//	// Transform int to string
//	formatCount := func(count int) string {
//	    return fmt.Sprintf("Users: %d", count)
//	}
//
//	// Adapt the effect to work with AppConfig and return string
//	adapted := effect.Promap(extractDBConfig, formatCount)(getUserCount)
//	result := adapted(AppConfig{DatabaseURL: "localhost:5432", APIKey: "secret"})
//
//go:inline
func Promap[E, A, D, B any](f Reader[D, E], g Reader[A, B]) Kleisli[D, Effect[E, A], B] {
	return F.Flow2(
		Local[A](f),
		Map[D](g),
	)
}
