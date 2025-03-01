// Copyright (c) 2024 IBM Corp.
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

package either

import (
	"github.com/IBM/fp-go/v2/internal/functor"
)

type eitherFunctor[E, A, B any] struct{}

func (o *eitherFunctor[E, A, B]) Map(f func(A) B) func(Either[E, A]) Either[E, B] {
	return Map[E, A, B](f)
}

// Functor implements the functoric operations for [Either]
func Functor[E, A, B any]() functor.Functor[A, B, Either[E, A], Either[E, B]] {
	return &eitherFunctor[E, A, B]{}
}
