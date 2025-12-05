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
func Sequence[R1, R2, A any](ma Reader[R2, Reader[R1, A]]) Kleisli[R2, R1, A] {
	return function.Flip(ma)
}

func Traverse[R2, R1, A, B any](
	f Kleisli[R1, A, B],
) func(Reader[R2, A]) Kleisli[R2, R1, B] {
	return readert.Traverse[Reader[R2, A]](
		identity.MonadMap,
		identity.MonadChain,
		f,
	)
}
