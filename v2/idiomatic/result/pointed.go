// Copyright (c) 2024 - 2025 IBM Corp.
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

package result

type (
	eitherPointed[A any] struct{}

	Pointed[A any] interface {
		Of(a A) (A, error)
	}
)

func (o eitherPointed[A]) Of(a A) (A, error) {
	return Of(a)
}

// Pointed implements the pointed functor operations for Either.
// A pointed functor provides the Of operation to lift a value into the Either context.
//
// Example:
//
//	p := either.Pointed[error, int]()
//	result := p.Of(42) // Right(42)
func MakePointed[A any]() Pointed[A] {
	return eitherPointed[A]{}
}
