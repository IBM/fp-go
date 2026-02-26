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

package constant

import (
	"github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
)

// Monoid creates a monoid that always returns a constant value.
//
// This creates a trivial monoid where both the Concat operation and Empty
// always return the same constant value, regardless of inputs. This is useful
// for testing, placeholder implementations, or when you need a monoid instance
// but the actual combining behavior doesn't matter.
//
// # Monoid Laws
//
// The constant monoid satisfies all monoid laws trivially:
//   - Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z)) - always returns 'a'
//   - Left Identity: Concat(Empty(), x) = x - both return 'a'
//   - Right Identity: Concat(x, Empty()) = x - both return 'a'
//
// Type Parameters:
//   - A: The type of the constant value
//
// Parameters:
//   - a: The constant value to return in all operations
//
// Returns:
//   - A Monoid[A] that always returns the constant value
//
// Example:
//
//	// Create a monoid that always returns 42
//	m := Monoid(42)
//	result := m.Concat(1, 2)  // 42
//	empty := m.Empty()         // 42
//
//	// Useful for testing or placeholder implementations
//	type Config struct {
//	    Timeout int
//	}
//	defaultConfig := Monoid(Config{Timeout: 30})
//	config := defaultConfig.Concat(Config{Timeout: 10}, Config{Timeout: 20})
//	// config is Config{Timeout: 30}
//
// See also:
//   - function.Constant2: The underlying constant function
//   - M.MakeMonoid: The monoid constructor
func Monoid[A any](a A) M.Monoid[A] {
	return M.MakeMonoid(function.Constant2[A, A](a), a)
}
