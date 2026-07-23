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
