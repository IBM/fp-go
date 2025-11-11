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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mToKm = MakeIso(
		func(m float32) float32 {
			return m / 1000
		},
		func(km float32) float32 {
			return km * 1000
		},
	)

	kmToMile = MakeIso(
		func(km float32) float32 {
			return km * 0.621371
		},
		func(mile float32) float32 {
			return mile / 0.621371
		},
	)
)

func TestGet(t *testing.T) {
	assert.Equal(t, mToKm.Get(100), float32(0.1))
	assert.Equal(t, Unwrap[float32](float32(100))(mToKm), float32(0.1))
	assert.Equal(t, To[float32](float32(100))(mToKm), float32(0.1))
}

func TestReverseGet(t *testing.T) {
	assert.Equal(t, mToKm.ReverseGet(1.2), float32(1200))
	assert.Equal(t, Wrap[float32](float32(1.2))(mToKm), float32(1200))
	assert.Equal(t, From[float32](float32(1.2))(mToKm), float32(1200))
}

func TestModify(t *testing.T) {

	double := func(x float32) float32 {
		return x * 2
	}

	assert.Equal(t, float32(2000), Modify[float32](double)(mToKm)(float32(1000)))
}

func TestReverse(t *testing.T) {

	double := func(x float32) float32 {
		return x * 2
	}

	assert.Equal(t, float32(4000), Modify[float32](double)(Reverse(mToKm))(float32(2000)))
}

func TestCompose(t *testing.T) {
	comp := Compose[float32](mToKm)(kmToMile)

	assert.InDelta(t, 0.93, comp.Get(1500), 0.01)
	assert.InDelta(t, 1609.34, comp.ReverseGet(1), 0.01)
}

func TestId(t *testing.T) {
	idIso := Id[int]()

	assert.Equal(t, 42, idIso.Get(42))
	assert.Equal(t, 42, idIso.ReverseGet(42))
}

func TestIMap(t *testing.T) {
	// Start with meters to kilometers
	localMToKm := MakeIso(
		func(m float32) float32 { return m / 1000 },
		func(km float32) float32 { return km * 1000 },
	)

	// Map to a different representation (string)
	kmToString := IMap[float32](
		func(km float32) string { return fmt.Sprintf("%.2f km", km) },
		func(s string) float32 {
			var km float32
			fmt.Sscanf(s, "%f km", &km)
			return km
		},
	)(localMToKm)

	assert.Equal(t, "1.50 km", kmToString.Get(1500))
	assert.InDelta(t, 2000, kmToString.ReverseGet("2.00 km"), 0.01)
}

func TestRoundTripLaws(t *testing.T) {
	// Test that isomorphisms satisfy round-trip laws

	// Law 1: ReverseGet(Get(s)) == s
	meters := float32(1500)
	assert.InDelta(t, meters, mToKm.ReverseGet(mToKm.Get(meters)), 0.001)

	// Law 2: Get(ReverseGet(a)) == a
	km := float32(1.5)
	assert.InDelta(t, km, mToKm.Get(mToKm.ReverseGet(km)), 0.001)
}

func TestComposeAssociativity(t *testing.T) {
	// Test that composition is associative
	// Compose left-to-right: (mToKm . kmToMile)
	leftCompose := Compose[float32](kmToMile)(mToKm)

	// Compose right-to-left should give same result
	rightCompose := Compose[float32](Compose[float32](kmToMile)(mToKm))(Id[float32]())

	meters := float32(1609.34)
	assert.InDelta(t, leftCompose.Get(meters), rightCompose.Get(meters), 0.01)
}
