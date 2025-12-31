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

package option

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/iso"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	// Endomorphism is a function from a type to itself (A → A).
	// It represents transformations that preserve the type.
	//
	// This is commonly used in lens setters to transform a structure
	// by applying a function that takes and returns the same type.
	//
	// Example:
	//   increment := N.Add(1)
	//   // increment is an Endomorphism[int]
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Lens represents a functional reference to a field within a structure.
	//
	// A Lens[S, A] provides a way to get and set a value of type A within
	// a structure of type S in an immutable way. It consists of:
	//   - Get: S → A (retrieve the value)
	//   - Set: A → S → S (update the value, returning a new structure)
	//
	// Lenses satisfy three laws:
	//   1. Get-Put: lens.Set(lens.Get(s))(s) == s
	//   2. Put-Get: lens.Get(lens.Set(a)(s)) == a
	//   3. Put-Put: lens.Set(b)(lens.Set(a)(s)) == lens.Set(b)(s)
	//
	// Type Parameters:
	//   - S: The structure type containing the field
	//   - A: The type of the field being focused on
	Lens[S, A any] = lens.Lens[S, A]

	// Option represents a value that may or may not be present.
	//
	// It is either Some[T] containing a value of type T, or None[T]
	// representing the absence of a value. This is a type-safe alternative
	// to using nil pointers.
	//
	// Type Parameters:
	//   - T: The type of the value that may be present
	Option[T any] = option.Option[T]

	// LensO is a lens that focuses on an optional value.
	//
	// A LensO[S, A] is equivalent to Lens[S, Option[A]], representing
	// a lens that focuses on a value of type A that may or may not be
	// present within a structure S.
	//
	// This is particularly useful for:
	//   - Nullable pointer fields
	//   - Optional configuration values
	//   - Fields that may be conditionally present
	//
	// The getter returns Option[A] (Some if present, None if absent).
	// The setter takes Option[A] (Some to set, None to remove).
	//
	// Type Parameters:
	//   - S: The structure type containing the optional field
	//   - A: The type of the optional value being focused on
	//
	// Example:
	//   type Config struct {
	//       Timeout *int
	//   }
	//
	//   timeoutLens := lens.MakeLensRef(
	//       func(c *Config) *int { return c.Timeout },
	//       func(c *Config, t *int) *Config { c.Timeout = t; return c },
	//   )
	//
	//   optLens := lens.FromNillableRef(timeoutLens)
	//   // optLens is a LensO[*Config, *int]
	LensO[S, A any] = Lens[S, Option[A]]

	// Kleisli represents a Kleisli arrow for optional lenses.
	// It's a function from A to LensO[S, B], used for composing optional lens operations.
	Kleisli[S, A, B any] = reader.Reader[A, LensO[S, B]]

	// Operator represents a function that transforms one optional lens into another.
	// It takes a LensO[S, A] and produces a LensO[S, B].
	Operator[S, A, B any] = Kleisli[S, LensO[S, A], B]

	// Iso represents an isomorphism between types S and A.
	// An isomorphism is a bidirectional transformation that preserves structure.
	Iso[S, A any] = iso.Iso[S, A]
)
