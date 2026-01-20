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

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/readert"
)

// Sequence swaps the order of the two environment parameters in a Kleisli arrow.
// It transforms a function that takes A and returns Reader[R2, B] into a function
// that takes R2 and returns Reader[R1, A].
//
// This is useful when you need to change the order of dependencies or when composing
// functions that expect their parameters in a different order.
//
// Type Parameters:
//   - R1: The first environment type (becomes second after flip)
//   - R2: The second environment type (becomes first after flip)
//   - A: The result type
//
// Parameters:
//   - ma: A Kleisli arrow from R2 to Reader[R1, A]
//
// Returns:
//   - A Kleisli arrow from R1 to Reader[R2, A] with swapped environment parameters
//
// Example:
//
//	type Config struct { Host string }
//	type Port int
//
//	// Original: takes Port, returns Reader[Config, string]
//	makeURL := func(port Port) reader.Reader[Config, string] {
//	    return func(c Config) string {
//	        return fmt.Sprintf("%s:%d", c.Host, port)
//	    }
//	}
//
//	// Sequenced: takes Config, returns Reader[Port, string]
//	sequenced := reader.Sequence(makeURL)
//	result := sequenced(Config{Host: "localhost"})(Port(8080))
//	// result: "localhost:8080"
//
// The Sequence operation is particularly useful when:
//   - You need to partially apply environments in a different order
//   - You're composing functions that expect parameters in reverse order
//   - You want to curry multi-parameter functions differently
//
//go:inline
func Sequence[R1, R2, A any](ma Kleisli[R1, R2, A]) Kleisli[R2, R1, A] {
	return function.Flip(ma)
}

// Traverse applies a Kleisli arrow to a value wrapped in a Reader, then sequences the result.
// It transforms a Reader[R2, A] into a function that takes R1 and returns Reader[R2, B],
// where the transformation from A to B is defined by a Kleisli arrow that depends on R1.
//
// This is useful when you have a Reader computation that produces a value, and you want to
// apply another Reader computation to that value, but with a different environment type.
// The result is a function that takes the second environment and returns a Reader that
// takes the first environment.
//
// Type Parameters:
//   - R2: The first environment type (outer Reader)
//   - R1: The second environment type (inner Reader/Kleisli)
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Kleisli arrow from A to B that depends on environment R1
//
// Returns:
//   - A function that takes a Reader[R2, A] and returns a Kleisli[R2, R1, B]
//
// The signature can be understood as:
//   - Input: Reader[R2, A] (a computation that produces A given R2)
//   - Output: func(R1) Reader[R2, B] (a function that takes R1 and produces a computation that produces B given R2)
//
// Example:
//
//	type Database struct { ConnectionString string }
//	type Config struct { TableName string }
//
//	// A Reader that gets a user ID from the database
//	getUserID := func(db Database) int {
//	    // Simulate database query
//	    return 42
//	}
//
//	// A Kleisli arrow that takes a user ID and returns a Reader that formats it with config
//	formatUser := func(id int) reader.Reader[Config, string] {
//	    return func(c Config) string {
//	        return fmt.Sprintf("User %d from table %s", id, c.TableName)
//	    }
//	}
//
//	// Traverse applies formatUser to the result of getUserID
//	traversed := reader.Traverse(formatUser)(getUserID)
//
//	// Now we can apply both environments
//	config := Config{TableName: "users"}
//	db := Database{ConnectionString: "localhost:5432"}
//	result := traversed(config)(db) // "User 42 from table users"
//
// The Traverse operation is particularly useful when:
//   - You need to compose computations that depend on different environments
//   - You want to apply a transformation that itself requires environmental context
//   - You're building pipelines where each stage has its own configuration
func Traverse[R2, R1, A, B any](
	f Kleisli[R1, A, B],
) func(Reader[R2, A]) Kleisli[R2, R1, B] {
	return readert.Traverse[Reader[R2, A]](
		identity.Map,
		identity.Chain,
		f,
	)
}
