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

package prism

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/optics/prism"
)

// Compose creates a Kleisli arrow that composes a prism with an isomorphism.
//
// This function takes a Prism[A, B] and returns a Kleisli arrow that can transform
// any Iso[S, A] into a Prism[S, B]. The resulting prism changes the source type from
// A to S using the bidirectional transformation provided by the isomorphism, while
// maintaining the same focus type B.
//
// The composition works as follows:
//   - GetOption: First transforms S to A using the iso's Get, then extracts B from A using the prism's GetOption
//   - ReverseGet: First constructs A from B using the prism's ReverseGet, then transforms A to S using the iso's ReverseGet
//
// This is the dual operation of optics/prism/iso.Compose:
//   - optics/prism/iso.Compose: Transforms the focus type (A → B) while keeping source type (S) constant
//   - optics/iso/prism.Compose: Transforms the source type (A → S) while keeping focus type (B) constant
//
// This is particularly useful when you have a prism that works with one type but you
// need to adapt it to work with a different source type that has a lossless bidirectional
// transformation to the original type.
//
// Type Parameters:
//   - S: The new source type after applying the isomorphism
//   - A: The original source type of the prism
//   - B: The focus type (remains constant through composition)
//
// Parameters:
//   - ab: A prism that extracts B from A
//
// Returns:
//   - A Kleisli arrow (function) that takes an Iso[S, A] and returns a Prism[S, B]
//
// Laws:
// The composed prism must satisfy the prism laws:
//  1. GetOption(ReverseGet(b)) == Some(b) for all b: B
//  2. If GetOption(s) == Some(a), then GetOption(ReverseGet(a)) == Some(a)
//
// These laws are preserved because:
//   - The isomorphism satisfies: ia.ReverseGet(ia.Get(s)) == s and ia.Get(ia.ReverseGet(a)) == a
//   - The original prism satisfies the prism laws
//
// Haskell Equivalent:
// This corresponds to the (.) operator for composing optics in Haskell's lens library,
// specifically when composing an Iso with a Prism:
//
//	iso . prism :: Iso s a -> Prism a b -> Prism s b
//
// In Haskell's lens library, this is part of the general optic composition mechanism.
// See: https://hackage.haskell.org/package/lens/docs/Control-Lens-Iso.html
//
// Example - Composing with Either prism:
//
//	import (
//	    "github.com/IBM/fp-go/v2/either"
//	    "github.com/IBM/fp-go/v2/optics/iso"
//	    "github.com/IBM/fp-go/v2/optics/prism"
//	    IP "github.com/IBM/fp-go/v2/optics/iso/prism"
//	    O "github.com/IBM/fp-go/v2/option"
//	)
//
//	// Create a prism that extracts Right values from Either[error, string]
//	rightPrism := prism.FromEither[error, string]()
//
//	// Create an isomorphism between []byte and string
//	bytesStringIso := iso.MakeIso(
//	    func(b []byte) string { return string(b) },
//	    func(s string) []byte { return []byte(s) },
//	)
//
//	// Compose them to get a prism that works with []byte as source
//	bytesPrism := IP.Compose(rightPrism)(bytesStringIso)
//
//	// Use the composed prism
//	// First converts []byte to string via iso, then extracts Right value
//	bytes := []byte("hello")
//	either := either.Right[error](string(bytes))
//	result := bytesPrism.GetOption(bytes)  // Extracts "hello" if Right
//
//	// Construct []byte from string
//	constructed := bytesPrism.ReverseGet("world")
//	// Returns []byte("world") wrapped in Right
//
// Example - Composing with custom types:
//
//	type JSON []byte
//	type Config struct {
//	    Host string
//	    Port int
//	}
//
//	// Isomorphism between JSON and []byte
//	jsonIso := iso.MakeIso(
//	    func(j JSON) []byte { return []byte(j) },
//	    func(b []byte) JSON { return JSON(b) },
//	)
//
//	// Prism that extracts Config from []byte (via JSON parsing)
//	configPrism := prism.MakePrism(
//	    func(b []byte) option.Option[Config] {
//	        var cfg Config
//	        if err := json.Unmarshal(b, &cfg); err != nil {
//	            return option.None[Config]()
//	        }
//	        return option.Some(cfg)
//	    },
//	    func(cfg Config) []byte {
//	        b, _ := json.Marshal(cfg)
//	        return b
//	    },
//	)
//
//	// Compose to work with JSON type instead of []byte
//	jsonConfigPrism := IP.Compose(configPrism)(jsonIso)
//
//	jsonData := JSON(`{"host":"localhost","port":8080}`)
//	config := jsonConfigPrism.GetOption(jsonData)
//	// config is Some(Config{Host: "localhost", Port: 8080})
//
// See also:
//   - github.com/IBM/fp-go/v2/optics/iso for isomorphism operations
//   - github.com/IBM/fp-go/v2/optics/prism for prism operations
//   - github.com/IBM/fp-go/v2/optics/prism/iso for the dual composition (transforming focus type)
func Compose[S, A, B any](ab Prism[A, B]) P.Kleisli[S, Iso[S, A], B] {
	return func(ia Iso[S, A]) Prism[S, B] {
		return P.MakePrismWithName(
			F.Flow2(ia.Get, ab.GetOption),
			F.Flow2(ab.ReverseGet, ia.ReverseGet),
			fmt.Sprintf("IsoCompose[%s -> %s]", ia, ab),
		)
	}
}
