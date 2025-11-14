// Copyright (c) 2023 - 2025 IBM Corp.
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
	"iter"

	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	Option[A any]    = option.Option[A]
	Lazy[A any]      = lazy.Lazy[A]
	Pair[L, R any]   = pair.Pair[L, R]
	Predicate[A any] = predicate.Predicate[A]
	IO[A any]        = io.IO[A]

	// Iterator represents a stateless, pure way to iterate over a sequence
	Iterator[U any] Lazy[Option[Pair[Iterator[U], U]]]

	Kleisli[A, B any]  = reader.Reader[A, Iterator[B]]
	Operator[A, B any] = Kleisli[Iterator[A], B]

	Seq[T any]     = iter.Seq[T]
	Seq2[K, V any] = iter.Seq2[K, V]
)
