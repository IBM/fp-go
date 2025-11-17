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
	eitherMonad[A, B any] struct{}

	Monad[A, B any] interface {
		Applicative[A, B]
		Chainable[A, B]
	}
)

func (o eitherMonad[A, B]) Of(a A) (A, error) {
	return Of(a)
}

func (o eitherMonad[A, B]) Map(f func(A) B) Operator[A, B] {
	return Map(f)
}

func (o eitherMonad[A, B]) Chain(f func(A) (B, error)) Operator[A, B] {
	return Chain(f)
}

func (o eitherMonad[A, B]) Ap(a A, err error) Operator[func(A) B, B] {
	return Ap[B](a, err)
}

// MakeMonad creates a Monad instance for Result operations.
// A monad combines the capabilities of Functor (Map), Applicative (Ap), and Chain (flatMap/bind).
// This allows for sequential composition of computations that may fail.
//
// Example:
//
//	m := result.MakeMonad[int, string]()
//	val, err := m.Chain(func(x int) (string, error) {
//	    if x > 0 {
//	        return result.Right[error](strconv.Itoa(x))
//	    }
//	    return result.Left[string](errors.New("negative"))
//	})(result.Right[error](42))
//	// val is "42", err is nil
func MakeMonad[A, B any]() Monad[A, B] {
	return eitherMonad[A, B]{}
}
