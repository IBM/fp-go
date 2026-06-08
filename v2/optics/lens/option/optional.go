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
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/optional"
	O "github.com/IBM/fp-go/v2/option"
)

// AsOptional converts a LensO[S, A] (a lens focusing on Option[A]) into an Optional[S, A].
//
// This conversion bridges two different optics paradigms:
//   - LensO[S, A] is a Lens[S, Option[A]] that always focuses on an Option[A] field
//   - Optional[S, A] is an optic that may or may not find a value of type A
//
// The conversion works by:
//   - Using the lens getter directly as the optional's GetOption
//   - Wrapping values in Some before passing them to the lens setter
//
// This is useful when you have a lens that focuses on an optional field and want to
// use it with optional operations that expect the focus to be directly on the value
// rather than on the Option wrapper.
//
// The resulting Optional satisfies the three optional laws:
//
//  1. GetSet Law (No-op on None):
//     If GetOption(s) returns None, then Set(a)(s) returns s unchanged (no-op).
//     This is enforced by checking GetOption before applying the set operation.
//     The implementation uses MapTo which only maps when the Option is Some,
//     ensuring that setting when GetOption returns None has no effect.
//
//     Formally: GetOption(s) = None => Set(a)(s) = s
//
//  2. SetGet Law (Get what you Set):
//     If GetOption(s) returns Some(_), then GetOption(Set(a)(s)) returns Some(a).
//     This is satisfied because Set only updates when GetOption returns Some,
//     wrapping the new value in Some before passing to the lens setter.
//
//     Formally: GetOption(s) = Some(_) => GetOption(Set(a)(s)) = Some(a)
//
//  3. SetSet Law (Last Set Wins):
//     Set(b)(Set(a)(s)) equals Set(b)(s).
//     This is satisfied because both operations check GetOption and only update
//     when it returns Some, with the lens SetSet law ensuring the last set wins.
//
//     Formally: Set(b)(Set(a)(s)) = Set(b)(s)
//
// Type Parameters:
//   - S: The structure type containing the optional field
//   - A: The type of the value that may be present
//
// Parameters:
//   - l: A lens focusing on an Option[A] field within structure S
//
// Returns:
//   - An Optional[S, A] that focuses directly on values of type A
//
// Example:
//
//	type Config struct {
//	    Timeout Option[int]
//	}
//
//	timeoutLens := lens.MakeLens(
//	    func(c Config) Option[int] { return c.Timeout },
//	    func(c Config, t Option[int]) Config { c.Timeout = t; return c },
//	)
//
//	// Convert to optional for direct value operations
//	timeoutOptional := AsOptional(timeoutLens)
//
//	config := Config{Timeout: O.Some(30)}
//
//	// Get the value directly (not wrapped in Option)
//	value := timeoutOptional.GetOption(config)
//	// value is Some(30)
//
//	// Set a value when Some exists (automatically wrapped in Some)
//	updated := timeoutOptional.Set(60)(config)
//	// updated.Timeout is Some(60)
//
//	// Set is a no-op when GetOption returns None (Law 1)
//	emptyConfig := Config{Timeout: O.None[int]()}
//	stillEmpty := timeoutOptional.Set(60)(emptyConfig)
//	// stillEmpty.Timeout is still None - Set is a no-op
func AsOptional[S, A any](l LensO[S, A]) Optional[S, A] {
	return optional.MakeOptionalCurriedWithName(
		l.Get,
		func(a A) func(S) S {
			return endomorphism.Join(F.Flow3(
				l.Get,
				O.MapTo[A](a),
				l.Set,
			))
		},
		l.String(),
	)
}
