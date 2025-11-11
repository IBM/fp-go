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

package generic

import (
	"github.com/IBM/fp-go/v2/internal/monad"
)

type arrayMonad[A, B any, GA ~[]A, GB ~[]B, GAB ~[]func(A) B] struct{}

func (o *arrayMonad[A, B, GA, GB, GAB]) Of(a A) GA {
	return Of[GA](a)
}

func (o *arrayMonad[A, B, GA, GB, GAB]) Map(f func(A) B) func(GA) GB {
	return Map[GA, GB](f)
}

func (o *arrayMonad[A, B, GA, GB, GAB]) Chain(f func(A) GB) func(GA) GB {
	return Chain[GA](f)
}

func (o *arrayMonad[A, B, GA, GB, GAB]) Ap(fa GA) func(GAB) GB {
	return Ap[GB, GAB](fa)
}

// Monad implements the monadic operations for an array
func Monad[A, B any, GA ~[]A, GB ~[]B, GAB ~[]func(A) B]() monad.Monad[A, B, GA, GB, GAB] {
	return &arrayMonad[A, B, GA, GB, GAB]{}
}
