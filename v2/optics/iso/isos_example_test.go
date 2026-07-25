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
	"strings"

	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
)

// ExampleFromStrictEquals demonstrates basic usage: mapping "no"/"yes" to false/true.
func ExampleFromStrictEquals() {
	iso := FromStrictEquals(lazy.Of("no"), lazy.Of("yes"))

	fmt.Println(iso.Get("yes"))
	fmt.Println(iso.Get("no"))
	fmt.Println(iso.ReverseGet(true))
	fmt.Println(iso.ReverseGet(false))
	// Output:
	// true
	// false
	// yes
	// no
}

// ExampleFromStrictEquals_integer demonstrates FromStrictEquals with integer sentinels.
func ExampleFromStrictEquals_integer() {
	iso := FromStrictEquals(lazy.Of(0), lazy.Of(1))

	fmt.Println(iso.Get(1))
	fmt.Println(iso.Get(0))
	fmt.Println(iso.ReverseGet(true))
	fmt.Println(iso.ReverseGet(false))
	// Output:
	// true
	// false
	// 1
	// 0
}

// ExampleFromStrictEquals_roundTrip demonstrates the round-trip isomorphism laws.
func ExampleFromStrictEquals_roundTrip() {
	iso := FromStrictEquals(lazy.Of(0), lazy.Of(1))

	// Law 1: ReverseGet(Get(sentinel)) == sentinel
	fmt.Println(iso.ReverseGet(iso.Get(1)))
	fmt.Println(iso.ReverseGet(iso.Get(0)))

	// Law 2: Get(ReverseGet(bool)) == bool
	fmt.Println(iso.Get(iso.ReverseGet(true)))
	fmt.Println(iso.Get(iso.ReverseGet(false)))
	// Output:
	// 1
	// 0
	// true
	// false
}

// ExampleFromStrictEquals_modify demonstrates using FromStrictEquals with Modify
// to toggle a sentinel value through bool space.
func ExampleFromStrictEquals_modify() {
	iso := FromStrictEquals(lazy.Of("inactive"), lazy.Of("active"))
	toggle := Modify[string](func(b bool) bool { return !b })(iso)

	fmt.Println(toggle("active"))
	fmt.Println(toggle("inactive"))
	// Output:
	// inactive
	// active
}

// ExampleFromEquals demonstrates FromEquals with a custom case-insensitive Eq.
func ExampleFromEquals() {
	caseInsensitiveEq := eq.FromEquals(func(a, b string) bool {
		return strings.EqualFold(a, b)
	})

	iso := FromEquals(lazy.Of("NO"), lazy.Of("YES"))(caseInsensitiveEq)

	// Any casing of "YES" maps to true
	fmt.Println(iso.Get("yes"))
	fmt.Println(iso.Get("YES"))
	fmt.Println(iso.Get("Yes"))

	// Any casing of "NO" maps to false
	fmt.Println(iso.Get("no"))

	// ReverseGet always returns the canonical sentinel
	fmt.Println(iso.ReverseGet(true))
	fmt.Println(iso.ReverseGet(false))
	// Output:
	// true
	// true
	// true
	// false
	// YES
	// NO
}

// ExampleFromEquals_roundTrip demonstrates round-trip laws for FromEquals.
func ExampleFromEquals_roundTrip() {
	iso := FromEquals(lazy.Of("off"), lazy.Of("on"))(eq.FromStrictEquals[string]())

	// Get(ReverseGet(b)) == b
	fmt.Println(iso.Get(iso.ReverseGet(true)))
	fmt.Println(iso.Get(iso.ReverseGet(false)))
	// Output:
	// true
	// false
}

// ExampleNegate demonstrates basic Negate usage: NOT in both directions.
func ExampleNegate() {
	iso := Negate()

	fmt.Println(iso.Get(true))
	fmt.Println(iso.Get(false))
	fmt.Println(iso.ReverseGet(true))
	fmt.Println(iso.ReverseGet(false))
	// Output:
	// false
	// true
	// false
	// true
}

// ExampleNegate_roundTrip demonstrates the involution law: double negation is identity.
func ExampleNegate_roundTrip() {
	iso := Negate()

	fmt.Println(iso.Get(iso.Get(true)))
	fmt.Println(iso.Get(iso.Get(false)))
	// Output:
	// true
	// false
}

// ExampleNegate_modify demonstrates using Negate with Modify to double-negate a bool
// that was lifted from a string sentinel.
func ExampleNegate_modify() {
	// Map "disabled"/"enabled" → bool, then negate to flip the meaning
	sentinelIso := FromStrictEquals(lazy.Of("disabled"), lazy.Of("enabled"))
	neg := Negate()
	toggle := Modify[string](neg.Get)(sentinelIso)

	fmt.Println(toggle("enabled"))
	fmt.Println(toggle("disabled"))
	// Output:
	// disabled
	// enabled
}

// ExampleRef demonstrates basic Ref usage: wrapping a value into a pointer,
// and returning the lazy fallback when ReverseGet receives nil.
func ExampleRef() {
	iso := Ref(lazy.Of(0))

	p := iso.Get(42)
	fmt.Println(*p)
	fmt.Println(iso.ReverseGet(p))
	fmt.Println(iso.ReverseGet(nil))
	// Output:
	// 42
	// 42
	// 0
}

// ExampleRef_roundTrip demonstrates the round-trip isomorphism laws for Ref.
func ExampleRef_roundTrip() {
	iso := Ref(lazy.Of(""))

	// ReverseGet(Get(v)) == v
	fmt.Println(iso.ReverseGet(iso.Get("hello")))

	// *Get(ReverseGet(p)) == *p
	original := "world"
	fmt.Println(*iso.Get(iso.ReverseGet(&original)))
	// Output:
	// hello
	// world
}

// ExampleDeref demonstrates basic Deref usage: dereferencing a pointer to its value,
// and returning the lazy fallback when Get receives nil.
func ExampleDeref() {
	iso := Deref(lazy.Of(0))

	v := 42
	fmt.Println(iso.Get(&v))
	fmt.Println(iso.Get(nil))

	p := iso.ReverseGet(99)
	fmt.Println(*p)
	// Output:
	// 42
	// 0
	// 99
}

// ExampleDeref_roundTrip demonstrates the round-trip isomorphism laws for Deref.
func ExampleDeref_roundTrip() {
	iso := Deref(lazy.Of(""))

	// Get(ReverseGet(v)) == v
	fmt.Println(iso.Get(iso.ReverseGet("hello")))

	// *ReverseGet(Get(p)) == *p
	original := "world"
	fmt.Println(*iso.ReverseGet(iso.Get(&original)))
	// Output:
	// hello
	// world
}

// ExampleToNillable demonstrates basic ToNillable usage: converting between
// Option[A] and nullable *A in both directions.
func ExampleToNillable() {
	iso := ToNillable[int]()

	// Some(v) → non-nil pointer
	p := iso.Get(O.Some(42))
	fmt.Println(*p)

	// None → nil
	fmt.Println(iso.Get(O.None[int]()) == nil)

	// non-nil pointer → Some
	v := 99
	fmt.Println(iso.ReverseGet(&v))

	// nil → None
	fmt.Println(iso.ReverseGet(nil))
	// Output:
	// 42
	// true
	// Some[int](99)
	// None[int]
}

// ExampleToNillable_roundTrip demonstrates the round-trip isomorphism laws for ToNillable.
func ExampleToNillable_roundTrip() {
	iso := ToNillable[int]()

	// ReverseGet(Get(Some(v))) == Some(v)
	fmt.Println(iso.ReverseGet(iso.Get(O.Some(7))))

	// ReverseGet(Get(None)) == None
	fmt.Println(iso.ReverseGet(iso.Get(O.None[int]())))
	// Output:
	// Some[int](7)
	// None[int]
}

// ExampleFromNillable demonstrates basic FromNillable usage: converting between
// nullable *A and Option[A] in both directions.
func ExampleFromNillable() {
	iso := FromNillable[int]()

	// non-nil pointer → Some
	v := 42
	fmt.Println(iso.Get(&v))

	// nil → None
	fmt.Println(iso.Get(nil))

	// Some(v) → non-nil pointer
	p := iso.ReverseGet(O.Some(99))
	fmt.Println(*p)

	// None → nil
	fmt.Println(iso.ReverseGet(O.None[int]()) == nil)
	// Output:
	// Some[int](42)
	// None[int]
	// 99
	// true
}

// ExampleFromNillable_roundTrip demonstrates the round-trip isomorphism laws for FromNillable.
func ExampleFromNillable_roundTrip() {
	iso := FromNillable[int]()

	// Get(ReverseGet(Some(v))) == Some(v)
	fmt.Println(iso.Get(iso.ReverseGet(O.Some(7))))

	// Get(ReverseGet(None)) == None
	fmt.Println(iso.Get(iso.ReverseGet(O.None[int]())))
	// Output:
	// Some[int](7)
	// None[int]
}
