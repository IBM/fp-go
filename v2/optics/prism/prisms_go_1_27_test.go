//go:build go1.27

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
	"testing"
	"uuid"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// The example UUID of RFC 9562 in the four textual forms accepted by uuid.Parse,
// plus its upper case spelling. All of them denote the very same UUID.
const (
	canonicalUUID = "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	bracedUUID    = "{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}"
	urnUUID       = "urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	compactUUID   = "f81d4fae7dec11d0a76500a0c91e6bf6"
	upperUUID     = "F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6"

	nilUUID = "00000000-0000-0000-0000-000000000000"
	maxUUID = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	oneUUID = "00000000-0000-0000-0000-000000000001"
)

// getUUID extracts the focused UUID, falling back to the nil UUID.
var getUUID = O.GetOrElse(uuid.Nil)

// TestParseUUID_Success tests that the prism parses all accepted textual forms
func TestParseUUID_Success(t *testing.T) {
	prism := ParseUUID()
	expected := uuid.MustParse(canonicalUUID)

	t.Run("parse canonical form", func(t *testing.T) {
		parsed := prism.GetOption(canonicalUUID)
		assert.True(t, O.IsSome(parsed))
		assert.Equal(t, expected, getUUID(parsed))
	})

	t.Run("parse braced form", func(t *testing.T) {
		assert.Equal(t, O.Of(expected), prism.GetOption(bracedUUID))
	})

	t.Run("parse urn form", func(t *testing.T) {
		assert.Equal(t, O.Of(expected), prism.GetOption(urnUUID))
	})

	t.Run("parse unhyphenated form", func(t *testing.T) {
		assert.Equal(t, O.Of(expected), prism.GetOption(compactUUID))
	})

	t.Run("parse upper case digits", func(t *testing.T) {
		assert.Equal(t, O.Of(expected), prism.GetOption(upperUUID))
	})

	t.Run("parse nil UUID", func(t *testing.T) {
		assert.Equal(t, O.Of(uuid.Nil()), prism.GetOption(nilUUID))
	})

	t.Run("parse max UUID", func(t *testing.T) {
		assert.Equal(t, O.Of(uuid.Max()), prism.GetOption(maxUUID))
	})
}

// TestParseUUID_Failure tests that malformed input yields None
func TestParseUUID_Failure(t *testing.T) {
	prism := ParseUUID()

	invalid := []string{
		"",
		"not-a-uuid",
		"f81d4fae-7dec-11d0-a765", // truncated
		"f81d4fae-7dec-11d0-a765-00a0c91e6bf6-extra",   // trailing garbage
		"f81d4fae-7dec-11d0-a765-00a0c91e6bfg",         // non hexadecimal digit
		"f81d4fae7dec-11d0-a765-00a0c91e6bf6",          // hyphens in the wrong place
		"{f81d4fae-7dec-11d0-a765-00a0c91e6bf6",        // unbalanced brace
		"urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf", // truncated urn form
		" f81d4fae-7dec-11d0-a765-00a0c91e6bf6",        // leading whitespace
		"f81d4fae-7dec-11d0-a765-00a0c91e6bf6 ",        // trailing whitespace
		"00000000-0000-0000-0000-00000000000",          // one digit short
		"000000000-0000-0000-0000-000000000000",        // one digit too long
	}

	for _, s := range invalid {
		t.Run(fmt.Sprintf("reject %q", s), func(t *testing.T) {
			assert.True(t, O.IsNone(prism.GetOption(s)))
		})
	}
}

// TestParseUUID_ReverseGet tests that ReverseGet renders the canonical form
func TestParseUUID_ReverseGet(t *testing.T) {
	prism := ParseUUID()

	t.Run("render known UUID", func(t *testing.T) {
		assert.Equal(t, canonicalUUID, prism.ReverseGet(uuid.MustParse(canonicalUUID)))
	})

	t.Run("render nil UUID", func(t *testing.T) {
		assert.Equal(t, nilUUID, prism.ReverseGet(uuid.Nil()))
	})

	t.Run("render max UUID", func(t *testing.T) {
		assert.Equal(t, maxUUID, prism.ReverseGet(uuid.Max()))
	})
}

// TestParseUUID_EdgeCases tests normalization of the non canonical input forms
func TestParseUUID_EdgeCases(t *testing.T) {
	prism := ParseUUID()

	// normalize parses a string and renders it back in canonical form
	normalize := F.Flow2(
		prism.GetOption,
		O.Map(prism.ReverseGet),
	)

	t.Run("all accepted forms normalize to the canonical form", func(t *testing.T) {
		for _, s := range []string{canonicalUUID, bracedUUID, urnUUID, compactUUID, upperUUID} {
			assert.Equal(t, O.Of(canonicalUUID), normalize(s))
		}
	})

	t.Run("invalid input stays None", func(t *testing.T) {
		assert.Equal(t, O.None[string](), normalize("not-a-uuid"))
	})

	t.Run("prism carries its name", func(t *testing.T) {
		assert.Equal(t, "PrismParseUUID", prism.String())
	})
}

// TestParseUUID_PrismLaws tests the prism laws for ParseUUID
func TestParseUUID_PrismLaws(t *testing.T) {
	prism := ParseUUID()

	t.Run("GetOption after ReverseGet is Some", func(t *testing.T) {
		uuids := A.From(uuid.Nil(), uuid.Max(), uuid.MustParse(canonicalUUID), uuid.NewV4(), uuid.NewV7())
		for _, u := range uuids {
			assert.Equal(t, O.Of(u), prism.GetOption(prism.ReverseGet(u)))
		}
	})

	t.Run("ReverseGet after GetOption is the identity on canonical input", func(t *testing.T) {
		roundTrip := F.Flow2(
			prism.GetOption,
			O.Map(prism.ReverseGet),
		)
		for _, s := range []string{canonicalUUID, nilUUID, maxUUID} {
			assert.Equal(t, O.Of(s), roundTrip(s))
		}
	})
}

// TestParseUUID_WithSet tests using the prism as a setter
func TestParseUUID_WithSet(t *testing.T) {
	setter := Set[string](uuid.MustParse(oneUUID))(ParseUUID())

	t.Run("set replaces a matching value", func(t *testing.T) {
		assert.Equal(t, oneUUID, setter(canonicalUUID))
	})

	t.Run("set replaces a non canonical value", func(t *testing.T) {
		assert.Equal(t, oneUUID, setter(bracedUUID))
	})

	t.Run("set leaves a non matching value unchanged", func(t *testing.T) {
		assert.Equal(t, "not-a-uuid", setter("not-a-uuid"))
	})
}

// TestParseUUID_Integration tests composition with other prisms
func TestParseUUID_Integration(t *testing.T) {
	// nonNilUUID parses a string and additionally rejects the nil UUID
	nonNilUUID := F.Pipe1(
		ParseUUID(),
		Compose[string](FromNonZero[uuid.UUID]()),
	)

	t.Run("accepts a non nil UUID", func(t *testing.T) {
		assert.Equal(t, O.Of(uuid.MustParse(canonicalUUID)), nonNilUUID.GetOption(canonicalUUID))
	})

	t.Run("rejects the nil UUID", func(t *testing.T) {
		assert.True(t, O.IsNone(nonNilUUID.GetOption(nilUUID)))
	})

	t.Run("rejects malformed input", func(t *testing.T) {
		assert.True(t, O.IsNone(nonNilUUID.GetOption("not-a-uuid")))
	})

	t.Run("reverse get is unaffected by the composition", func(t *testing.T) {
		assert.Equal(t, canonicalUUID, nonNilUUID.ReverseGet(uuid.MustParse(canonicalUUID)))
	})
}

func BenchmarkParseUUID_GetOption(b *testing.B) {
	prism := ParseUUID()
	for range b.N {
		prism.GetOption(canonicalUUID)
	}
}

func BenchmarkParseUUID_ReverseGet(b *testing.B) {
	prism := ParseUUID()
	id := uuid.MustParse(canonicalUUID)
	for range b.N {
		prism.ReverseGet(id)
	}
}

// ExampleParseUUID demonstrates parsing strings into UUIDs and back.
func ExampleParseUUID() {
	uuidPrism := ParseUUID()

	fmt.Println(uuidPrism.GetOption("f81d4fae-7dec-11d0-a765-00a0c91e6bf6"))
	fmt.Println(uuidPrism.GetOption("not-a-uuid"))
	fmt.Println(uuidPrism.ReverseGet(uuid.Nil()))

	// Output:
	// Some[uuid.UUID](f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
	// None[uuid.UUID]
	// 00000000-0000-0000-0000-000000000000
}

// ExampleParseUUID_normalize demonstrates how the prism normalizes the
// alternative textual forms accepted by uuid.Parse.
func ExampleParseUUID_normalize() {
	uuidPrism := ParseUUID()

	normalize := F.Flow2(
		uuidPrism.GetOption,
		O.Map(uuidPrism.ReverseGet),
	)

	fmt.Println(normalize("{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}"))
	fmt.Println(normalize("urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6"))
	fmt.Println(normalize("F81D4FAE7DEC11D0A76500A0C91E6BF6"))
	fmt.Println(normalize("not-a-uuid"))

	// Output:
	// Some[string](f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
	// Some[string](f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
	// Some[string](f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
	// None[string]
}

// ExampleParseUUID_compose demonstrates composing the prism with FromNonZero so
// that the nil UUID is rejected in addition to malformed input.
func ExampleParseUUID_compose() {
	nonNilUUID := F.Pipe1(
		ParseUUID(),
		Compose[string](FromNonZero[uuid.UUID]()),
	)

	fmt.Println(nonNilUUID.GetOption("f81d4fae-7dec-11d0-a765-00a0c91e6bf6"))
	fmt.Println(nonNilUUID.GetOption("00000000-0000-0000-0000-000000000000"))
	fmt.Println(nonNilUUID.GetOption("not-a-uuid"))

	// Output:
	// Some[uuid.UUID](f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
	// None[uuid.UUID]
	// None[uuid.UUID]
}

// ExampleParseUUID_set demonstrates replacing a UUID in a string that is only
// modified when it actually holds a valid UUID.
func ExampleParseUUID_set() {
	setter := Set[string](uuid.MustParse("00000000-0000-0000-0000-000000000001"))(ParseUUID())

	fmt.Println(setter("f81d4fae-7dec-11d0-a765-00a0c91e6bf6"))
	fmt.Println(setter("not-a-uuid"))

	// Output:
	// 00000000-0000-0000-0000-000000000001
	// not-a-uuid
}
