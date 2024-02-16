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

package ioeither

import (
	"github.com/IBM/fp-go/internal/monad"
	"github.com/IBM/fp-go/internal/pointed"
	G "github.com/IBM/fp-go/ioeither/generic"
)

// Pointed returns the pointed operations for [IOEither]
func Pointed[E, A any]() pointed.Pointed[A, IOEither[E, A]] {
	return G.Pointed[E, A, IOEither[E, A]]()
}

// Monad returns the monadic operations for [IOEither]
func Monad[E, A, B any]() monad.Monad[A, B, IOEither[E, A], IOEither[E, B], IOEither[E, func(A) B]] {
	return G.Monad[E, A, B, IOEither[E, A], IOEither[E, B], IOEither[E, func(A) B]]()
}
