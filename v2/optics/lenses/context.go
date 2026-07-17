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

package lenses

import (
	"context"
	"fmt"

	CR "github.com/IBM/fp-go/v2/context/reader"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
)

// AtContext creates a Lens that focuses on a typed value stored in a context.Context
// under a given key.
//
// The getter reads the value associated with key from the context and performs a
// type assertion to V, returning Option[V]: Some(v) if the value exists and has
// type V, or None if the key is absent or the stored value has an incompatible type.
//
// The setter accepts Option[V]:
//   - Some(v): returns a new child context that carries key → v (via context.WithValue).
//   - None: returns the original context unchanged.
//
// The lens is named "AtContext[<key>]" where <key> is formatted with fmt.Sprintf.
//
// Type Parameters:
//   - V: The type of the value stored in the context under key.
//   - K: The type of the key used to look up and store the value.
//
// Parameters:
//   - key: The context key used to store and retrieve the value.
//
// Returns:
//   - Lens[context.Context, Option[V]]: A lens focusing on the optional V stored
//     in a context.Context under key.
//
// See Also:
//   - context.WithValue: The standard library function used by the setter.
//   - option.InstanceOf: The type-safe assertion used by the getter.
//   - context/reader.WithValue: The Kleisli arrow used to derive child contexts.
func AtContext[V, K any](key K) Lens[context.Context, Option[V]] {

	return lens.MakeLensCurriedWithName(
		F.Pipe1(
			F.Bind2nd((context.Context).Value, any(key)),
			reader.Map[context.Context](option.InstanceOf[V]),
		),
		F.Flow2(
			option.Map(CR.WithValue[V](key)),
			option.GetOrElse(reader.Ask[context.Context]),
		),
		fmt.Sprintf("AtContext[%v]", key),
	)
}
