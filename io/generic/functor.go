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
	"github.com/IBM/fp-go/internal/functor"
)

type ioFunctor[A, B any, GA ~func() A, GB ~func() B] struct{}

func (o *ioFunctor[A, B, GA, GB]) Map(f func(A) B) func(GA) GB {
	return Map[GA, GB, A, B](f)
}

// Functor implements the functoric operations for [IO]
func Functor[A, B any, GA ~func() A, GB ~func() B]() functor.Functor[A, B, GA, GB] {
	return &ioFunctor[A, B, GA, GB]{}
}
