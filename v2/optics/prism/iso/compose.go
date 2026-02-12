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

package iso

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
)

// Compose creates an operator that composes an isomorphism with a prism.
//
// This function takes an isomorphism Iso[A, B] and returns an operator that can
// transform any Prism[S, A] into a Prism[S, B]. The resulting prism maintains
// the same source type S but changes the focus type from A to B using the
// bidirectional transformation provided by the isomorphism.
//
// The composition works as follows:
//   - GetOption: First extracts A from S using the prism, then transforms A to B using the iso's Get
//   - ReverseGet: First transforms B to A using the iso's ReverseGet, then constructs S using the prism's ReverseGet
//
// This is particularly useful when you have a prism that focuses on one type but
// you need to work with a different type that has a lossless bidirectional
// transformation to the original type.
//
// Haskell Equivalent:
// This corresponds to the (.) operator for composing optics in Haskell's lens library,
// specifically when composing a Prism with an Iso:
//
//	prism . iso :: Prism s a -> Iso a b -> Prism s b
//
// In Haskell's lens library, this is part of the general optic composition mechanism.
// See: https://hackage.haskell.org/package/lens/docs/Control-Lens-Prism.html
//
// Type Parameters:
//   - S: The source type (sum type) that the prism operates on
//   - A: The original focus type of the prism
//   - B: The new focus type after applying the isomorphism
//
// Parameters:
//   - ab: An isomorphism between types A and B that defines the bidirectional transformation
//
// Returns:
//   - An Operator[S, A, B] that transforms Prism[S, A] into Prism[S, B]
//
// Laws:
// The composed prism must satisfy the prism laws:
//  1. GetOption(ReverseGet(b)) == Some(b) for all b: B
//  2. If GetOption(s) == Some(a), then GetOption(ReverseGet(a)) == Some(a)
//
// These laws are preserved because:
//   - The isomorphism satisfies: ab.ReverseGet(ab.Get(a)) == a and ab.Get(ab.ReverseGet(b)) == b
//   - The original prism satisfies the prism laws
//
// Example - Composing string/bytes isomorphism with Either prism:
//
//	import (
//	    "github.com/IBM/fp-go/v2/either"
//	    "github.com/IBM/fp-go/v2/optics/iso"
//	    "github.com/IBM/fp-go/v2/optics/prism"
//	    PI "github.com/IBM/fp-go/v2/optics/prism/iso"
//	    O "github.com/IBM/fp-go/v2/option"
//	)
//
//	// Create an isomorphism between string and []byte
//	stringBytesIso := iso.MakeIso(
//	    func(s string) []byte { return []byte(s) },
//	    func(b []byte) string { return string(b) },
//	)
//
//	// Create a prism that extracts Right values from Either[error, string]
//	rightPrism := prism.FromEither[error, string]()
//
//	// Compose them to get a prism that works with []byte instead of string
//	bytesPrism := PI.Compose(stringBytesIso)(rightPrism)
//
//	// Extract bytes from a Right value
//	success := either.Right[error]("hello")
//	result := bytesPrism.GetOption(success)
//	// result is Some([]byte("hello"))
//
//	// Extract from a Left value returns None
//	failure := either.Left[string](errors.New("error"))
//	result = bytesPrism.GetOption(failure)
//	// result is None
//
//	// Construct an Either from bytes
//	constructed := bytesPrism.ReverseGet([]byte("world"))
//	// constructed is Right("world")
//
// Example - Composing with custom types:
//
//	type Celsius float64
//	type Fahrenheit float64
//
//	// Isomorphism between Celsius and Fahrenheit
//	tempIso := iso.MakeIso(
//	    func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
//	    func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
//	)
//
//	// Prism that extracts temperature from a weather report
//	type WeatherReport struct {
//	    Temperature Celsius
//	    Condition   string
//	}
//	tempPrism := prism.MakePrism(
//	    func(w WeatherReport) option.Option[Celsius] {
//	        return option.Some(w.Temperature)
//	    },
//	    func(c Celsius) WeatherReport {
//	        return WeatherReport{Temperature: c}
//	    },
//	)
//
//	// Compose to work with Fahrenheit instead
//	fahrenheitPrism := PI.Compose(tempIso)(tempPrism)
//
//	report := WeatherReport{Temperature: 20, Condition: "sunny"}
//	temp := fahrenheitPrism.GetOption(report)
//	// temp is Some(68.0) in Fahrenheit
//
// See also:
//   - github.com/IBM/fp-go/v2/optics/iso for isomorphism operations
//   - github.com/IBM/fp-go/v2/optics/prism for prism operations
//   - Operator for the type signature of the returned function
func Compose[S, A, B any](ab Iso[A, B]) Operator[S, A, B] {
	return func(pa Prism[S, A]) Prism[S, B] {
		return P.MakePrismWithName(
			F.Flow2(
				pa.GetOption,
				O.Map(ab.Get),
			),
			F.Flow2(
				ab.ReverseGet,
				pa.ReverseGet,
			),
			fmt.Sprintf("PrismCompose[%s -> %s]", pa, ab),
		)
	}
}
