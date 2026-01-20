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

import (
	"github.com/IBM/fp-go/v2/function"
)

// Promap is the profunctor map operation that transforms both the input and output of a Reader.
// It applies f to the input (contravariantly) and g to the output (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Example:
//
//	type Config struct { Port int }
//	type Env struct { Config Config }
//	getPort := func(c Config) int { return c.Port }
//	extractConfig := func(e Env) Config { return e.Config }
//	toString := strconv.Itoa
//	r := reader.Promap(extractConfig, toString)(getPort)
//	result := r(Env{Config: Config{Port: 8080}}) // "8080"
func Promap[E, A, D, B any](f func(D) E, g func(A) B) Kleisli[D, Reader[E, A], B] {
	return function.Bind13of3(function.Flow3[func(D) E, func(E) A, func(A) B])(f, g)
}

// Local changes the value of the local context during the execution of the action `ma`.
// This is similar to Contravariant's contramap and allows you to modify the environment
// before passing it to a Reader.
//
// Example:
//
//	type DetailedConfig struct { Host string; Port int }
//	type SimpleConfig struct { Host string }
//	getHost := func(c SimpleConfig) string { return c.Host }
//	simplify := func(d DetailedConfig) SimpleConfig { return SimpleConfig{Host: d.Host} }
//	r := reader.Local(simplify)(getHost)
//	result := r(DetailedConfig{Host: "localhost", Port: 8080}) // "localhost"
//
//go:inline
func Local[A, R1, R2 any](f func(R2) R1) Kleisli[R2, Reader[R1, A], A] {
	return Compose[A](f)
}

// WithLocal applies a local environment transformation to an existing Reader.
// It takes a Reader that expects environment R1 and a transformation function
// that converts R2 to R1, returning a new Reader that expects environment R2.
//
// This is the "applied" version of Local - while Local returns a Kleisli arrow
// that can transform any Reader, WithLocal directly applies the transformation
// to a specific Reader instance.
//
// Type Parameters:
//   - A: The output type of the Reader
//   - R1: The original environment type that the Reader expects
//   - R2: The new environment type that the returned Reader will expect
//
// Parameters:
//   - fa: The original Reader[R1, A] to transform
//   - f: A function that converts the new environment R2 to the original environment R1
//
// Returns:
//   - Reader[R2, A]: A new Reader that accepts R2 and produces A
//
// The transformation is contravariant - it modifies the input environment before
// passing it to the original Reader. This allows you to adapt a Reader to work
// with a different (typically broader) environment type.
//
// Example:
//
//	type DetailedConfig struct { Host string; Port int; Timeout time.Duration }
//	type SimpleConfig struct { Host string; Port int }
//
//	// A Reader that works with SimpleConfig
//	getEndpoint := func(c SimpleConfig) string {
//	    return fmt.Sprintf("%s:%d", c.Host, c.Port)
//	}
//
//	// Transform it to work with DetailedConfig
//	simplify := func(d DetailedConfig) SimpleConfig {
//	    return SimpleConfig{Host: d.Host, Port: d.Port}
//	}
//	adaptedReader := reader.WithLocal(getEndpoint, simplify)
//
//	// Now it can be called with DetailedConfig
//	detailed := DetailedConfig{Host: "localhost", Port: 8080, Timeout: 30*time.Second}
//	endpoint := adaptedReader(detailed) // "localhost:8080"
//
// See also: Local, Contramap
//
//go:inline
func WithLocal[A, R1, R2 any](fa Reader[R1, A], f func(R2) R1) Reader[R2, A] {
	return function.Flow2(f, fa)
}

// Contramap is an alias for Local.
// It changes the value of the local context during the execution of a Reader.
// This is the contravariant functor operation that transforms the input environment.
//
// Contramap is semantically identical to Local - both modify the environment before
// passing it to a Reader. The name "Contramap" emphasizes the contravariant nature
// of the transformation (transforming the input rather than the output).
//
// Example:
//
//	type DetailedConfig struct { Host string; Port int }
//	type SimpleConfig struct { Host string }
//	getHost := func(c SimpleConfig) string { return c.Host }
//	simplify := func(d DetailedConfig) SimpleConfig { return SimpleConfig{Host: d.Host} }
//	r := reader.Contramap(simplify)(getHost)
//	result := r(DetailedConfig{Host: "localhost", Port: 8080}) // "localhost"
//
// See also: Local
//
//go:inline
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, Reader[R1, A], A] {
	return Compose[A](f)
}
