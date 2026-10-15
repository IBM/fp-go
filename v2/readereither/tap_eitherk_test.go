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

package readereither

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

type tapCfg struct{}

// tapPositive fails for numbers that are not positive.
func tapPositive(n int) Either[string, int] {
	if n > 0 {
		return either.Right[string](n)
	}
	return either.Left[int]("not positive")
}

func TestTapEitherK(t *testing.T) {
	all := func(start ReaderEither[tapCfg, string, int]) []ReaderEither[tapCfg, string, int] {
		return []ReaderEither[tapCfg, string, int]{
			TapEitherK[tapCfg](tapPositive)(start),
			MonadTapEitherK(start, tapPositive),
			ChainFirstEitherK[tapCfg](tapPositive)(start),
		}
	}
	t.Run("keeps the value", func(t *testing.T) {
		for _, m := range all(Of[tapCfg, string](7)) {
			assert.Equal(t, either.Right[string](7), m(tapCfg{}))
		}
	})
	t.Run("failing check replaces the value", func(t *testing.T) {
		for _, m := range all(Of[tapCfg, string](-7)) {
			assert.Equal(t, either.Left[int]("not positive"), m(tapCfg{}))
		}
	})
	t.Run("upstream error is preserved", func(t *testing.T) {
		for _, m := range all(Left[tapCfg, int]("upstream")) {
			assert.Equal(t, either.Left[int]("upstream"), m(tapCfg{}))
		}
	})
}

// ExampleTapEitherK validates a value without changing it.
func ExampleTapEitherK() {
	fmt.Println(F.Pipe1(Of[tapCfg, string](3), TapEitherK[tapCfg](tapPositive))(tapCfg{}))
	fmt.Println(F.Pipe1(Of[tapCfg, string](-3), TapEitherK[tapCfg](tapPositive))(tapCfg{}))
	// Output:
	// Right[int](3)
	// Left[string](not positive)
}

// ExampleMonadTapEitherK is the non curried version of TapEitherK.
func ExampleMonadTapEitherK() {
	fmt.Println(MonadTapEitherK(Of[tapCfg, string](3), tapPositive)(tapCfg{}))
	// Output: Right[int](3)
}
