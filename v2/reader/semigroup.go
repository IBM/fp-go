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

package reader

import "github.com/IBM/fp-go/v2/monoid"

// ApplicativeMonoid returns a [Monoid] that concatenates [Reader] instances via their applicative.
// This combines two Reader values by applying the underlying monoid's combine operation
// to their results using applicative application.
//
// The applicative behavior means that both Reader computations are executed with the same
// environment, and their results are combined using the underlying monoid. This is useful
// for accumulating values from multiple Reader computations that all depend on the same
// environment.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for Reader[R, A].
//
// Example:
//
//	type Config struct { Port int; Timeout int }
//	intMonoid := number.MonoidSum[int]()
//	readerMonoid := ApplicativeMonoid[Config](intMonoid)
//
//	getPort := func(c Config) int { return c.Port }
//	getTimeout := func(c Config) int { return c.Timeout }
//	combined := readerMonoid.Concat(getPort, getTimeout)
//	// Result: func(c Config) int { return c.Port + c.Timeout }
//
//	config := Config{Port: 8080, Timeout: 30}
//	result := combined(config) // 8110
//
//go:inline
func ApplicativeMonoid[R, A any](m monoid.Monoid[A]) monoid.Monoid[Reader[R, A]] {
	return monoid.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[A, R, A],
		m,
	)
}
