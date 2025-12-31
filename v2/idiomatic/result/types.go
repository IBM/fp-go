// Copyright (c) 2025 IBM Corp.
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
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	// Option is a type alias for option.Option, provided for convenience
	// when working with Result and Option together.
	Option[A any] = option.Option[A]

	// Lens is an optic that focuses on a field of type T within a structure of type S.
	Lens[S, T any] = lens.Lens[S, T]

	// Endomorphism represents a function from a type to itself (T -> T).
	Endomorphism[T any] = endomorphism.Endomorphism[T]

	// Kleisli represents a Kleisli arrow for the idiomatic Result pattern.
	// It's a function from A to (B, error), following Go's idiomatic error handling.
	Kleisli[A, B any] = func(A) (B, error)

	// Operator represents a function that transforms one Result into another.
	// It takes (A, error) and produces (B, error), following Go's idiomatic pattern.
	Operator[A, B any] = func(A, error) (B, error)

	Predicate[A any] = predicate.Predicate[A]
)
