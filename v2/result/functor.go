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

import (
	"github.com/IBM/fp-go/v2/internal/functor"
)

type eitherFunctor[A, B any] struct{}

func (o *eitherFunctor[A, B]) Map(f func(A) B) Operator[A, B] {
	return Map(f)
}

// Functor implements the functoric operations for Either.
// A functor provides the Map operation that transforms values inside a context
// while preserving the structure.
//
// Example:
//
//	f := either.Functor[error, int, string]()
//	result := f.Map(strconv.Itoa)(either.Right[error](42))
//	// result is Right("42")
func Functor[A, B any]() functor.Functor[A, B, Result[A], Result[B]] {
	return &eitherFunctor[A, B]{}
}
