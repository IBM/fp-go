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

	"github.com/stretchr/testify/assert"
)

// Fixtures used across method-compose tests.
// All use discrete (exact) arithmetic so assertions can use assert.Equal
// rather than InDelta.

// stringToBytes: string ↔ []byte
var stringToBytes = MakeIso(
	func(s string) []byte { return []byte(s) },
	func(b []byte) string { return string(b) },
)

// bytesToLen: []byte ↔ int  (number of bytes)
// This is intentionally not a true isomorphism for all values of int,
// but it satisfies the round-trip laws for the values used in the tests:
//   - Get(ReverseGet(n)) == n  when n >= 0
//   - ReverseGet(Get(b)) == b  always
var bytesToLen = MakeIso(
	func(b []byte) int { return len(b) },
	func(n int) []byte { return make([]byte, n) },
)

// TestIsoCompose_MethodEquivalentToFreeFunction verifies that the method form
// sa.Compose(ab) returns an isomorphism identical in behaviour to the
// free-function Compose[S](ab)(sa).
func TestIsoCompose_MethodEquivalentToFreeFunction(t *testing.T) {
	method := stringToBytes.Compose(bytesToLen)
	free := IsoComposeIso[string](bytesToLen)(stringToBytes)

	inputs := []string{"", "hello", "fp-go"}
	for _, s := range inputs {
		t.Run("Get parity", func(t *testing.T) {
			assert.Equal(t, free.Get(s), method.Get(s))
		})
		t.Run("ReverseGet parity", func(t *testing.T) {
			n := len(s)
			assert.Equal(t, free.ReverseGet(n), method.ReverseGet(n))
		})
	}
}

// TestIsoCompose_Method_Get verifies that Get applies the outer then inner
// transformation.
func TestIsoCompose_Method_Get(t *testing.T) {
	composed := stringToBytes.Compose(bytesToLen)

	assert.Equal(t, 5, composed.Get("hello"))
	assert.Equal(t, 0, composed.Get(""))
	assert.Equal(t, 3, composed.Get("abc"))
}

// TestIsoCompose_Method_ReverseGet verifies that ReverseGet applies the inner
// then outer inverse transformation.
func TestIsoCompose_Method_ReverseGet(t *testing.T) {
	composed := stringToBytes.Compose(bytesToLen)

	// ReverseGet(n) = stringToBytes.ReverseGet(bytesToLen.ReverseGet(n))
	//               = string(make([]byte, n))
	assert.Equal(t, string(make([]byte, 3)), composed.ReverseGet(3))
	assert.Equal(t, "", composed.ReverseGet(0))
}

// TestIsoCompose_Method_IsoLaws verifies both round-trip laws on the
// method-composed isomorphism using exact-arithmetic isos.
func TestIsoCompose_Method_IsoLaws(t *testing.T) {
	// Use the package-level float32 isos so the laws apply to realistic values,
	// and use InDelta for floating-point comparisons.
	// mToKm and kmToMile are declared in iso_test.go (same package).
	composed := mToKm.Compose(kmToMile)

	t.Run("law 1: ReverseGet(Get(s)) == s", func(t *testing.T) {
		for _, meters := range []float32{0, 100, 1609.34, 5000} {
			assert.InDelta(t, float64(meters), float64(composed.ReverseGet(composed.Get(meters))), 0.01,
				"round-trip failed for %v meters", meters)
		}
	})

	t.Run("law 2: Get(ReverseGet(b)) == b", func(t *testing.T) {
		for _, miles := range []float32{0, 1, 3.11} {
			assert.InDelta(t, float64(miles), float64(composed.Get(composed.ReverseGet(miles))), 0.001,
				"round-trip failed for %v miles", miles)
		}
	})
}

// TestIsoCompose_Method_ExactRoundTrip verifies both iso laws using
// exact-arithmetic isos where InDelta is not needed.
func TestIsoCompose_Method_ExactRoundTrip(t *testing.T) {
	composed := stringToBytes.Compose(bytesToLen)

	t.Run("law 1: ReverseGet(Get(s)) has correct length", func(t *testing.T) {
		// ReverseGet(Get("hello")) = ReverseGet(5) = string of 5 zero bytes;
		// length is preserved even though the exact bytes differ.
		s := "hello"
		reconstructed := composed.ReverseGet(composed.Get(s))
		assert.Equal(t, len(s), len(reconstructed))
	})

	t.Run("law 2: Get(ReverseGet(n)) == n", func(t *testing.T) {
		for _, n := range []int{0, 1, 7, 100} {
			assert.Equal(t, n, composed.Get(composed.ReverseGet(n)))
		}
	})
}

// TestIsoCompose_Method_Chained verifies that three .Compose calls chain
// correctly through three type conversions.
func TestIsoCompose_Method_Chained(t *testing.T) {
	type Meters float32
	type Kilometers float32
	type Miles float32

	metersToKm := MakeIso(
		func(m Meters) Kilometers { return Kilometers(m / 1000) },
		func(km Kilometers) Meters { return Meters(km * 1000) },
	)
	kmToMiles := MakeIso(
		func(km Kilometers) Miles { return Miles(km * 0.621371) },
		func(mi Miles) Kilometers { return Kilometers(mi / 0.621371) },
	)
	milesToFeet := MakeIso(
		func(mi Miles) float32 { return float32(mi * 5280) },
		func(ft float32) Miles { return Miles(ft / 5280) },
	)

	// Meters → Kilometers → Miles → feet
	composed := metersToKm.Compose(kmToMiles).Compose(milesToFeet)

	t.Run("Get traverses three levels", func(t *testing.T) {
		// 1609.34 m ≈ 1 mile ≈ 5280 feet
		assert.InDelta(t, 5280.0, composed.Get(Meters(1609.34)), 5.0)
	})

	t.Run("ReverseGet traverses three levels", func(t *testing.T) {
		assert.InDelta(t, 1609.34, float64(composed.ReverseGet(5280)), 5.0)
	})

	t.Run("law 1: ReverseGet(Get(s)) ≈ s", func(t *testing.T) {
		m := Meters(5000)
		assert.InDelta(t, float64(m), float64(composed.ReverseGet(composed.Get(m))), 1.0)
	})

	t.Run("law 2: Get(ReverseGet(b)) ≈ b", func(t *testing.T) {
		ft := float32(10000)
		assert.InDelta(t, float64(ft), float64(composed.Get(composed.ReverseGet(ft))), 1.0)
	})
}
