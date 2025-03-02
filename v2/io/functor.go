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

package io

import (
	"github.com/IBM/fp-go/v2/internal/functor"
)

type (
	ioFunctor[A, B any] struct{}

	IOFunctor[A, B any] = functor.Functor[A, B, IO[A], IO[B]]
)

func (o *ioFunctor[A, B]) Map(f func(A) B) func(IO[A]) IO[B] {
	return Map(f)
}

// Functor implements the functoric operations for [IO]
func Functor[A, B any]() IOFunctor[A, B] {
	return &ioFunctor[A, B]{}
}
