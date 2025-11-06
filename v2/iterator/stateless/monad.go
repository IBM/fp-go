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

package stateless

import (
	"github.com/IBM/fp-go/v2/internal/monad"
	G "github.com/IBM/fp-go/v2/iterator/stateless/generic"
)

// Monad returns the monadic operations for an [Iterator]
func Monad[A, B any]() monad.Monad[A, B, Iterator[A], Iterator[B], Iterator[func(A) B]] {
	return G.Monad[A, B, Iterator[A], Iterator[B], Iterator[func(A) B]]()
}
