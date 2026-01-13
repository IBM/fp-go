// Copyright (c) 2025 IBM Corp.
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

package state

import (
	"strconv"
	"testing"

	"github.com/IBM/fp-go/v2/optics/iso"
	P "github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestIMapBasic tests basic IMap functionality
func TestIMapBasic(t *testing.T) {
	t.Run("transform state and value using isomorphism", func(t *testing.T) {
		// State that increments an int state and returns the old value
		increment := func(s int) P.Pair[int, int] {
			return P.MakePair(s+1, s)
		}

		// Isomorphism between string and int (string length <-> int)
		stringIntIso := iso.MakeIso(
			S.Size,
			strconv.Itoa,
		)

		// Transform int to string
		toString := strconv.Itoa

		adapted := IMap(stringIntIso, toString)(increment)
		result := adapted("hello") // length is 5

		// State should be "6" (5+1), value should be "5"
		assert.Equal(t, P.MakePair("6", "5"), result)
	})
}

// TestMapStateBasic tests basic MapState functionality
func TestMapStateBasic(t *testing.T) {
	t.Run("transform only state using isomorphism", func(t *testing.T) {
		// State that doubles the state and returns it
		double := func(s int) P.Pair[int, int] {
			doubled := s * 2
			return P.MakePair(doubled, doubled)
		}

		// Isomorphism between string and int
		stringIntIso := iso.MakeIso(
			func(s string) int { n, _ := strconv.Atoi(s); return n },
			strconv.Itoa,
		)

		adapted := MapState[int](stringIntIso)(double)
		result := adapted("5")

		// State should be "10" (5*2), value should be 10
		assert.Equal(t, P.MakePair("10", 10), result)
	})
}

// TestIMapComposition tests composing IMap transformations
func TestIMapComposition(t *testing.T) {
	t.Run("compose two IMap transformations", func(t *testing.T) {
		// Simple state that returns the state unchanged
		identity := func(s int) P.Pair[int, int] {
			return P.MakePair(s, s)
		}

		// First isomorphism: bool <-> int (false=0, true=1)
		boolIntIso := iso.MakeIso(
			func(b bool) int {
				if b {
					return 1
				}
				return 0
			},
			func(n int) bool { return n != 0 },
		)

		// Transform value
		addOne := func(n int) int { return n + 1 }

		adapted := IMap(boolIntIso, addOne)(identity)
		result := adapted(true) // true -> 1

		// State should be true (1 -> true), value should be 2 (1+1)
		assert.Equal(t, P.MakePair(true, 2), result)
	})
}
