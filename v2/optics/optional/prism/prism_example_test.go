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

package prism_test

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/common"
	OPT "github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/optics/optional/prism"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
)

// ---- shared test types ----

type exResult interface{ isExResult() }
type exSuccess struct{ Value int }
type exFailure struct{ Error string }

func (exSuccess) isExResult() {}
func (exFailure) isExResult() {}

// ---- AsOptional examples ----

// ExampleAsOptional demonstrates converting a Prism into an Optional.
// The resulting Optional has the same GetOption, but its Set is a no-op when
// the prism does not match.
func ExampleAsOptional() {
	successPrism := P.MakePrism(
		func(r exResult) O.Option[int] {
			if s, ok := r.(exSuccess); ok {
				return O.Some(s.Value)
			}
			return O.None[int]()
		},
		func(v int) exResult { return exSuccess{Value: v} },
	)

	opt := prism.AsOptional(successPrism)

	// GetOption returns Some for the matching variant.
	fmt.Println(opt.GetOption(exSuccess{Value: 42}))

	// GetOption returns None for a non-matching variant.
	fmt.Println(opt.GetOption(exFailure{Error: "oops"}))

	// Set updates when the variant matches.
	updated := opt.Set(100)(exSuccess{Value: 1})
	fmt.Println(updated.(exSuccess).Value)

	// Set is a no-op when the variant does not match.
	unchanged := opt.Set(100)(exFailure{Error: "nope"})
	fmt.Println(unchanged.(exFailure).Error)

	// Output:
	// Some[int](42)
	// None[int]
	// 100
	// nope
}

// ExampleAsOptional_noOpOnMismatch demonstrates that Set is a no-op when the
// prism does not match, satisfying the Optional no-op law.
func ExampleAsOptional_noOpOnMismatch() {
	successPrism := P.MakePrism(
		func(r exResult) O.Option[int] {
			if s, ok := r.(exSuccess); ok {
				return O.Some(s.Value)
			}
			return O.None[int]()
		},
		func(v int) exResult { return exSuccess{Value: v} },
	)

	opt := prism.AsOptional(successPrism)
	failure := exResult(exFailure{Error: "original"})

	result := opt.Set(999)(failure)
	fmt.Println(result.(exFailure).Error)
	// Output:
	// original
}

// ---- PrismSome examples ----

// ExamplePrismSome demonstrates the PrismSome prism, which focuses on the Some
// variant of an Option type.
func ExamplePrismSome() {
	p := prism.PrismSome[int]()

	// GetOption on Some returns the same Some.
	fmt.Println(p.GetOption(O.Some(42)))

	// GetOption on None returns None.
	fmt.Println(p.GetOption(O.None[int]()))

	// ReverseGet wraps a value in Some.
	fmt.Println(p.ReverseGet(7))

	// Output:
	// Some[int](42)
	// None[int]
	// Some[int](7)
}

// ---- Some examples ----

// ExampleSome demonstrates using Some to drill through an Optional[S, Option[A]]
// and focus directly on the A inside the Option.
func ExampleSome() {
	type Settings struct {
		Timeout O.Option[int]
	}

	// Optional[Settings, Option[int]] – always focuses on the Timeout field.
	timeoutOpt := OPT.MakeOptional(
		func(s Settings) O.Option[O.Option[int]] { return O.Some(s.Timeout) },
		func(s Settings, opt O.Option[int]) Settings { s.Timeout = opt; return s },
	)

	// Drill through to the int inside Some.
	valueOpt := prism.Some[Settings, int](timeoutOpt)

	present := Settings{Timeout: O.Some(30)}
	absent := Settings{Timeout: O.None[int]()}

	fmt.Println(valueOpt.GetOption(present))
	fmt.Println(valueOpt.GetOption(absent))

	updated := valueOpt.Set(60)(present)
	fmt.Println(updated.Timeout)

	// Set is a no-op when Timeout is None.
	unchanged := valueOpt.Set(60)(absent)
	fmt.Println(unchanged.Timeout)

	// Output:
	// Some[int](30)
	// None[int]
	// Some[int](60)
	// None[int]
}

// ---- Compose examples ----

// ExampleCompose demonstrates composing an Optional with a Prism using F.Pipe1.
// The result is an Optional that is None when either the outer optional or the
// prism does not match.
func ExampleCompose() {
	type Notification struct {
		Payload C.Optional[exResult, int] // not stored — we build it via optional
	}

	// Optional[exResult, int] via AsOptional + sum-type prism.
	successPrism := P.MakePrism(
		func(r exResult) O.Option[int] {
			if s, ok := r.(exSuccess); ok {
				return O.Some(s.Value)
			}
			return O.None[int]()
		},
		func(v int) exResult { return exSuccess{Value: v} },
	)

	// Outer optional: focuses on exResult when it is non-nil (always here via
	// Id + Some, but we use a plain MakeOptional to keep the example clear).
	outerOpt := OPT.MakeOptional(
		func(r exResult) O.Option[exResult] {
			if r != nil {
				return O.Some(r)
			}
			return O.None[exResult]()
		},
		func(_ exResult, r exResult) exResult { return r },
	)

	// Compose outer optional with the success prism.
	composed := F.Pipe1(outerOpt, prism.Compose[exResult](successPrism))

	fmt.Println(composed.GetOption(exSuccess{Value: 42}))
	fmt.Println(composed.GetOption(exFailure{Error: "err"}))

	updated := composed.Set(99)(exSuccess{Value: 1})
	fmt.Println(updated.(exSuccess).Value)

	// Set is a no-op when the prism does not match.
	notUpdated := composed.Set(99)(exFailure{Error: "keep"})
	fmt.Println(notUpdated.(exFailure).Error)

	// Output:
	// Some[int](42)
	// None[int]
	// 99
	// keep
}
