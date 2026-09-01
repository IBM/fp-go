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

package common

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// optionals used across the method-compose tests.
// outer: Person → Name  (None when Name is empty)
// inner: string → first rune  (None when string is empty)

func makeNameOptional() Optional[Person, string] {
	return MakeOptional(
		func(p Person) Option[string] {
			if p.Name != "" {
				return OptionSome(p.Name)
			}
			return OptionNone[string]()
		},
		func(p Person, name string) Person {
			p.Name = name
			return p
		},
	)
}

func makeFirstCharOptional() Optional[string, rune] {
	return MakeOptional(
		func(s string) Option[rune] {
			if len(s) > 0 {
				return OptionSome(rune(s[0]))
			}
			return OptionNone[rune]()
		},
		func(s string, r rune) string {
			if len(s) > 0 {
				return string(r) + s[1:]
			}
			return string(r)
		},
	)
}

// TestOptionalCompose_MethodEquivalentToFreeFunction verifies that the method
// form o.Compose(ab) returns an optional identical in behaviour to the
// free-function Compose[S](ab)(o).
func TestOptionalCompose_MethodEquivalentToFreeFunction(t *testing.T) {
	outer := makeNameOptional()
	inner := makeFirstCharOptional()

	method := outer.Compose(inner)
	free := OptionalComposeOptional[Person](inner)(outer)

	persons := []Person{
		{Name: "Alice", Age: 30},
		{Name: "", Age: 20},
	}
	for _, p := range persons {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(p), method.GetOption(p))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set('Z')(p), method.Set('Z')(p))
		})
	}
}

// TestOptionalCompose_Method_GetOption_Match verifies Some is returned when
// both optionals produce a value.
func TestOptionalCompose_Method_GetOption_Match(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	got := composed.GetOption(Person{Name: "Alice"})
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 'A', OptionGetOrElse(F.Constant(rune(0)))(got))
}

// TestOptionalCompose_Method_GetOption_OuterNone verifies None is returned when
// the outer optional produces None.
func TestOptionalCompose_Method_GetOption_OuterNone(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	assert.True(t, OptionIsNone(composed.GetOption(Person{Name: ""})))
}

// TestOptionalCompose_Method_GetOption_InnerNone verifies None is returned when
// the outer optional matches but the inner optional produces None.
func TestOptionalCompose_Method_GetOption_InnerNone(t *testing.T) {
	// inner optional that always returns None (simulates empty focus)
	alwaysNone := MakeOptional(
		func(_ string) Option[rune] { return OptionNone[rune]() },
		func(s string, _ rune) string { return s },
	)
	composed := makeNameOptional().Compose(alwaysNone)

	assert.True(t, OptionIsNone(composed.GetOption(Person{Name: "Alice"})))
}

// TestOptionalCompose_Method_Set_Match verifies that Set updates the focus when
// both optionals match.
func TestOptionalCompose_Method_Set_Match(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	updated := composed.Set('B')(Person{Name: "Alice", Age: 30})
	assert.Equal(t, "Blice", updated.Name)
	assert.Equal(t, 30, updated.Age)
}

// TestOptionalCompose_Method_Set_NoOpOnOuterNone verifies that Set is a no-op
// when the outer optional does not match (optional no-op law).
func TestOptionalCompose_Method_Set_NoOpOnOuterNone(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	original := Person{Name: "", Age: 42}
	result := composed.Set('Z')(original)
	assert.Equal(t, original, result)
}

// TestOptionalCompose_Method_OptionalLaws verifies the three optional laws on
// the method-composed optional.
func TestOptionalCompose_Method_OptionalLaws(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	matching := Person{Name: "Alice", Age: 30}
	nonMatching := Person{Name: "", Age: 30}

	t.Run("law GetSet: GetOption==None => Set is no-op", func(t *testing.T) {
		result := composed.Set('Z')(nonMatching)
		assert.Equal(t, nonMatching, result)
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := composed.Set('B')(matching)
		got := composed.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 'B', OptionGetOrElse(F.Constant(rune(0)))(got))
	})

	t.Run("law SetSet: Set(b)(Set(a)(s)) == Set(b)(s)", func(t *testing.T) {
		result1 := composed.Set('X')(composed.Set('Y')(matching))
		result2 := composed.Set('X')(matching)
		assert.Equal(t, result2, result1)
	})
}

// TestOptionalCompose_Method_Chained verifies that three .Compose calls chain
// correctly through three levels of optional structure.
func TestOptionalCompose_Method_Chained(t *testing.T) {
	// level 1: Response → *Info  (None when info is nil)
	// level 2: *Info → *Employment  (None when employment is nil)
	// level 3: *Employment → *Phone  (None when phone is nil)
	l1 := OptionalFromPredicateRef[Response](F.IsNonNil[Info])((*Response).GetInfo, (*Response).SetInfo)
	l2 := MakeOptionalRef(
		func(i *Info) Option[*Employment] {
			if i.employment != nil {
				return OptionSome(i.employment)
			}
			return OptionNone[*Employment]()
		},
		func(i *Info, e *Employment) *Info { i.employment = e; return i },
	)
	l3 := MakeOptionalRef(
		func(e *Employment) Option[*Phone] {
			if e.phone != nil {
				return OptionSome(e.phone)
			}
			return OptionNone[*Phone]()
		},
		func(e *Employment, p *Phone) *Employment { e.phone = p; return e },
	)

	composed := l1.Compose(l2).Compose(l3)

	phone := &Phone{number: "555-1234"}
	emp := &Employment{phone: phone}
	info := &Info{employment: emp}
	resp := &Response{info: info}

	t.Run("all three levels present", func(t *testing.T) {
		got := composed.GetOption(resp)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, phone, OptionGetOrElse(F.Constant[*Phone](nil))(got))
	})

	t.Run("middle level nil short-circuits", func(t *testing.T) {
		resp2 := &Response{info: &Info{employment: nil}}
		assert.True(t, OptionIsNone(composed.GetOption(resp2)))
	})

	t.Run("outermost level nil short-circuits", func(t *testing.T) {
		resp3 := &Response{info: nil}
		assert.True(t, OptionIsNone(composed.GetOption(resp3)))
	})

	t.Run("Set is no-op when chain breaks", func(t *testing.T) {
		resp4 := &Response{info: nil}
		result := composed.Set(&Phone{number: "999"})(resp4)
		assert.Nil(t, result.info)
	})
}
