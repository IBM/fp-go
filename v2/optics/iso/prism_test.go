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

package iso_test

import (
	"fmt"
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/iso"
	"github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// ---- test types -------------------------------------------------------

// tcToken is a sum type representing a parsed token.
type tcToken interface{ isToken() }

// tcNumber is the numeric variant.
type tcNumber struct{ Value int }

// tcWord is the string variant.
type tcWord struct{ Value string }

func (tcNumber) isToken() {}
func (tcWord) isToken()   {}

// ---- shared helpers ---------------------------------------------------

// numberPrism focuses on the tcNumber variant of tcToken.
var numberPrism = prism.MakePrism(
	F.Flow2(F.ToAny[tcToken], O.InstanceOf[tcNumber]),
	func(n tcNumber) tcToken { return n },
)

// utf8Iso is an Iso[[]byte, string].
var utf8Iso = iso.MakeIso(
	func(b []byte) string { return string(b) },
	func(s string) []byte { return []byte(s) },
)

// stringToTokenIso is an Iso[string, tcToken] that encodes a string as a tcWord token.
var stringToTokenIso = iso.MakeIso(
	func(s string) tcToken { return tcWord{Value: s} },
	func(t tcToken) string {
		if w, ok := t.(tcWord); ok {
			return w.Value
		}
		return ""
	},
)

// numberStringIso is an Iso[string, tcToken] that encodes a decimal integer string
// as a tcNumber token (best-effort; non-numeric strings produce tcNumber{0}).
var numberStringIso = iso.MakeIso(
	func(s string) tcToken {
		n, _ := strconv.Atoi(s)
		return tcNumber{Value: n}
	},
	func(t tcToken) string {
		if n, ok := t.(tcNumber); ok {
			return strconv.Itoa(n.Value)
		}
		return "0"
	},
)

// ---- TestComposePrism -------------------------------------------------

// TestComposePrism_GetOption_Match verifies Some is returned when the Prism
// matches after the Iso has transformed the source.
func TestComposePrism_GetOption_Match(t *testing.T) {
	p := iso.ComposePrism[string](numberPrism)(numberStringIso)

	got := p.GetOption("42")
	assert.Equal(t, O.Some(tcNumber{Value: 42}), got)
}

// TestComposePrism_GetOption_NoMatch verifies None is returned when the Prism
// does not match the value produced by the Iso.
func TestComposePrism_GetOption_NoMatch(t *testing.T) {
	// stringToTokenIso always produces a tcWord, so numberPrism never matches.
	p := iso.ComposePrism[string](numberPrism)(stringToTokenIso)

	assert.Equal(t, O.None[tcNumber](), p.GetOption("hello"))
	assert.Equal(t, O.None[tcNumber](), p.GetOption(""))
}

// TestComposePrism_ReverseGet verifies that ReverseGet constructs S from B by
// threading through prism.ReverseGet then iso.ReverseGet.
func TestComposePrism_ReverseGet(t *testing.T) {
	p := iso.ComposePrism[string](numberPrism)(numberStringIso)

	// prism.ReverseGet(tcNumber{7}) == tcNumber{7} (as tcToken)
	// iso.ReverseGet(tcNumber{7}) == "7"
	got := p.ReverseGet(tcNumber{Value: 7})
	assert.Equal(t, "7", got)
}

// TestComposePrism_PrismLaws verifies both prism laws on the composed prism.
func TestComposePrism_PrismLaws(t *testing.T) {
	p := iso.ComposePrism[string](numberPrism)(numberStringIso)

	t.Run("law 1: GetOption(ReverseGet(b)) == Some(b)", func(t *testing.T) {
		b := tcNumber{Value: 99}
		got := p.GetOption(p.ReverseGet(b))
		assert.Equal(t, O.Some(b), got)
	})

	t.Run("law 2: if GetOption(s)==Some(b) then GetOption(ReverseGet(b))==Some(b)", func(t *testing.T) {
		got := p.GetOption("100")
		if O.IsSome(got) {
			b := O.GetOrElse(F.Constant(tcNumber{}))(got)
			assert.Equal(t, O.Some(b), p.GetOption(p.ReverseGet(b)))
		}
	})
}

// TestComposePrism_IdentityIso verifies that composing with the identity Iso
// leaves the prism behaviour unchanged.
func TestComposePrism_IdentityIso(t *testing.T) {
	idIso := iso.Id[tcToken]()
	p := iso.ComposePrism[tcToken](numberPrism)(idIso)

	token := tcToken(tcNumber{Value: 5})

	t.Run("GetOption matches same as underlying prism", func(t *testing.T) {
		assert.Equal(t, numberPrism.GetOption(token), p.GetOption(token))
	})

	t.Run("ReverseGet matches same as underlying prism", func(t *testing.T) {
		assert.Equal(t, numberPrism.ReverseGet(tcNumber{Value: 5}), p.ReverseGet(tcNumber{Value: 5}))
	})
}

// TestComposePrism_PipelineUsage verifies the curried form composes cleanly
// inside F.Pipe1.
func TestComposePrism_PipelineUsage(t *testing.T) {
	piped := F.Pipe1(numberStringIso, iso.ComposePrism[string](numberPrism))
	direct := iso.ComposePrism[string](numberPrism)(numberStringIso)

	t.Run("GetOption parity", func(t *testing.T) {
		assert.Equal(t, direct.GetOption("7"), piped.GetOption("7"))
		assert.Equal(t, direct.GetOption("abc"), piped.GetOption("abc"))
	})
	t.Run("ReverseGet parity", func(t *testing.T) {
		b := tcNumber{Value: 3}
		assert.Equal(t, direct.ReverseGet(b), piped.ReverseGet(b))
	})
}

// TestComposePrism_Chaining verifies two sequential ComposePrism calls:
// []byte →(utf8Iso)→ string →(numberStringIso)→ tcToken, with numberPrism focusing tcNumber.
func TestComposePrism_Chaining(t *testing.T) {
	// step 1: string source — prism focuses on tcNumber inside a tcToken-bearing iso
	stringPrism := iso.ComposePrism[string](numberPrism)(numberStringIso)

	// step 2: widen source from string to []byte via utf8Iso
	bytesPrism := iso.ComposePrism[[]byte](stringPrism)(utf8Iso)

	// wordPrism uses the word iso so the numberPrism never matches.
	wordStringPrism := iso.ComposePrism[string](numberPrism)(stringToTokenIso)
	wordBytesPrism := iso.ComposePrism[[]byte](wordStringPrism)(utf8Iso)

	t.Run("GetOption navigates through both isos", func(t *testing.T) {
		got := bytesPrism.GetOption([]byte("55"))
		assert.Equal(t, O.Some(tcNumber{Value: 55}), got)
	})

	t.Run("GetOption returns None when prism misses", func(t *testing.T) {
		// stringToTokenIso always produces tcWord, so numberPrism never matches.
		assert.Equal(t, O.None[tcNumber](), wordBytesPrism.GetOption([]byte("hello")))
	})

	t.Run("ReverseGet threads through both iso.ReverseGet calls", func(t *testing.T) {
		got := bytesPrism.ReverseGet(tcNumber{Value: 8})
		assert.Equal(t, []byte("8"), got)
	})

	t.Run("law 1 still holds after chaining", func(t *testing.T) {
		b := tcNumber{Value: 42}
		assert.Equal(t, O.Some(b), bytesPrism.GetOption(bytesPrism.ReverseGet(b)))
	})
}

// ---- Examples ---------------------------------------------------------

// ExampleComposePrism demonstrates widening the source type of a Prism using
// an Iso.  The Prism focuses on the tcNumber variant of tcToken; the Iso
// encodes a decimal string as a tcNumber token, giving a Prism[string, tcNumber].
func ExampleComposePrism() {
	numIso := iso.MakeIso(
		func(s string) tcToken {
			n, _ := strconv.Atoi(s)
			return tcNumber{Value: n}
		},
		func(t tcToken) string {
			if n, ok := t.(tcNumber); ok {
				return strconv.Itoa(n.Value)
			}
			return "0"
		},
	)

	numPrism := prism.MakePrism(
		F.Flow2(F.ToAny[tcToken], O.InstanceOf[tcNumber]),
		func(n tcNumber) tcToken { return n },
	)

	p := F.Pipe1(numIso, iso.ComposePrism[string](numPrism))

	// GetOption: iso converts "42" to tcNumber{42}, prism extracts it.
	fmt.Println(O.IsSome(p.GetOption("42")))

	// GetOption returns None when composed with a non-matching iso.
	wordIso := iso.MakeIso(
		func(s string) tcToken { return tcWord{Value: s} },
		func(t tcToken) string {
			if w, ok := t.(tcWord); ok {
				return w.Value
			}
			return ""
		},
	)
	wordP := F.Pipe1(wordIso, iso.ComposePrism[string](numPrism))
	fmt.Println(O.IsNone(wordP.GetOption("hello")))

	// ReverseGet: prism.ReverseGet(tcNumber{7}) → tcToken, iso.ReverseGet → "7".
	fmt.Println(p.ReverseGet(tcNumber{Value: 7}))

	// Output:
	// true
	// true
	// 7
}

// ExampleComposePrism_eitherSource shows ComposePrism widening from
// Either[error, int] to []byte via two isos chained with F.Pipe1.
func ExampleComposePrism_eitherSource() {
	// Prism: extracts Right values from Either[error, int].
	rightPrism := prism.FromEither[error, int]()

	// Iso: string ↔ Either[error, int] (parse / format).
	stringEitherIso := iso.MakeIso(
		func(s string) E.Either[error, int] {
			n, err := strconv.Atoi(s)
			if err != nil {
				return E.Left[int](err)
			}
			return E.Right[error](n)
		},
		func(e E.Either[error, int]) string {
			return strconv.Itoa(E.GetOrElse(func(error) int { return 0 })(e))
		},
	)

	// Iso: []byte ↔ string.
	bytesIso := iso.MakeIso(
		func(b []byte) string { return string(b) },
		func(s string) []byte { return []byte(s) },
	)

	// Chain: []byte →(bytesIso)→ string →(stringEitherIso)→ Either → prism focus int.
	stringPrism := F.Pipe1(stringEitherIso, iso.ComposePrism[string](rightPrism))
	bytesPrism := F.Pipe1(bytesIso, iso.ComposePrism[[]byte](stringPrism))

	fmt.Println(O.IsSome(bytesPrism.GetOption([]byte("10"))))
	fmt.Println(O.IsNone(bytesPrism.GetOption([]byte("NaN"))))
	fmt.Println(string(bytesPrism.ReverseGet(99)))

	// Output:
	// true
	// true
	// 99
}
