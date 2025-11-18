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

package state

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	Endomorphism[A any] = endomorphism.Endomorphism[A]
	Lens[S, A any]      = lens.Lens[S, A]
	// some type aliases
	Reader[R, A any] = reader.Reader[R, A]
	Pair[L, R any]   = pair.Pair[L, R]

	// State represents an operation on top of a current [State] that produces a value and a new [State]
	State[S, A any] = Reader[S, pair.Pair[S, A]]

	Kleisli[S, A, B any]  = Reader[A, State[S, B]]
	Operator[S, A, B any] = Kleisli[S, State[S, A], B]
)
