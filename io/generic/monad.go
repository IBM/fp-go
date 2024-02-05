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

package generic

import (
	"github.com/IBM/fp-go/internal/monad"
)

type ioMonad[A, B any, GA ~func() A, GB ~func() B, GAB ~func() func(A) B] struct{}

func (o *ioMonad[A, B, GA, GB, GAB]) Of(a A) GA {
	return Of[GA, A](a)
}

func (o *ioMonad[A, B, GA, GB, GAB]) Map(f func(A) B) func(GA) GB {
	return Map[GA, GB, A, B](f)
}

func (o *ioMonad[A, B, GA, GB, GAB]) Chain(f func(A) GB) func(GA) GB {
	return Chain[GA, GB, A, B](f)
}

func (o *ioMonad[A, B, GA, GB, GAB]) Ap(fa GA) func(GAB) GB {
	return Ap[GB, GAB, GA, B, A](fa)
}

// Monad implements the monadic operations for [Option]
func Monad[A, B any, GA ~func() A, GB ~func() B, GAB ~func() func(A) B]() monad.Monad[A, B, GA, GB, GAB] {
	return &ioMonad[A, B, GA, GB, GAB]{}
}
