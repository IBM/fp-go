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

package decode

import (
	"github.com/IBM/fp-go/v2/option"
)

// FromOption lifts an Option into the Decode monad.
//
// A present value is turned into a successful decoder via [Of]. An absent value
// falls back to the decoder produced by onNone, which decides how the absence is
// reported.
//
// The fallback is a whole Decode rather than a bare Errors value on purpose:
// a Decode has access to the validation context at decode time, so the failure it
// produces can carry the full path through the surrounding nested structure. This
// is the difference to readereither.FromOption, whose onNone yields a fixed error
// that knows nothing about where in the input the failure occurred.
//
// # Type Parameters
//
//   - I: The input type of the resulting decoder (the validation context in codec pipelines)
//   - A: The type held by the Option and produced by the decoder
//
// # Parameters
//
//   - onNone: A deferred Decode used when the Option is None. It is evaluated only
//     for the absent case, so building it may be as expensive as needed.
//
// # Returns
//
//   - A Kleisli arrow func(Option[A]) Decode[I, A]
//
// # Example
//
// Decoding a required map key, failing with the context aware message:
//
//	lookup := F.Flow2(
//	    record.Lookup[string]("name"),
//	    decode.FromOption[validation.Context](
//	        lazy.Of(validation.FailureWithMessage[string]("name", "key is unavailable")),
//	    ),
//	)
//
// Supplying a default instead of failing:
//
//	withDefault := decode.FromOption[validation.Context](
//	    lazy.Of(decode.Of[validation.Context](0)),
//	)
//
// # See Also
//
//   - Of: Lifts a plain value into Decode, used for the Some case
//   - Left: Builds a Decode that fails with fixed errors
//   - validation.FailureWithMessage: Builds a context aware failing Decode
func FromOption[I, A any](onNone Lazy[Decode[I, A]]) Kleisli[I, Option[A], A] {
	return option.Fold(onNone, Of[I, A])
}
