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

package readerresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type tapCfg struct{}

// tapPositive fails for numbers that are not positive.
func tapPositive(n int) Result[int] {
	if n > 0 {
		return result.Of(n)
	}
	return result.Left[int](errors.New("not positive"))
}

func TestTapResultK(t *testing.T) {
	all := func(start ReaderResult[tapCfg, int]) []ReaderResult[tapCfg, int] {
		return []ReaderResult[tapCfg, int]{
			TapResultK[tapCfg](tapPositive)(start),
			MonadTapResultK(start, tapPositive),
			ChainFirstEitherK[tapCfg](tapPositive)(start),
			MonadChainFirstEitherK(start, tapPositive),
			TapEitherK[tapCfg](tapPositive)(start),
			MonadTapEitherK(start, tapPositive),
		}
	}
	t.Run("keeps the value", func(t *testing.T) {
		for _, m := range all(Of[tapCfg](7)) {
			assert.Equal(t, result.Of(7), m(tapCfg{}))
		}
	})
	t.Run("failing check replaces the value", func(t *testing.T) {
		for _, m := range all(Of[tapCfg](-7)) {
			assert.Equal(t, result.Left[int](errors.New("not positive")), m(tapCfg{}))
		}
	})
	t.Run("upstream error is preserved", func(t *testing.T) {
		err := errors.New("upstream")
		for _, m := range all(Left[tapCfg, int](err)) {
			assert.Equal(t, result.Left[int](err), m(tapCfg{}))
		}
	})
}

// ExampleTapResultK validates a value without changing it.
func ExampleTapResultK() {
	fmt.Println(result.Unwrap(F.Pipe1(Of[tapCfg](3), TapResultK[tapCfg](tapPositive))(tapCfg{})))
	fmt.Println(result.Unwrap(F.Pipe1(Of[tapCfg](-3), TapResultK[tapCfg](tapPositive))(tapCfg{})))
	// Output:
	// 3 <nil>
	// 0 not positive
}

// ExampleMonadTapResultK is the non curried version of TapResultK.
func ExampleMonadTapResultK() {
	fmt.Println(result.Unwrap(MonadTapResultK(Of[tapCfg](3), tapPositive)(tapCfg{})))
	// Output: 3 <nil>
}

// ExampleChainFirstEitherK is an alias for ChainFirstResultK.
func ExampleChainFirstEitherK() {
	fmt.Println(result.Unwrap(F.Pipe1(Of[tapCfg](3), ChainFirstEitherK[tapCfg](tapPositive))(tapCfg{})))
	// Output: 3 <nil>
}

// ExampleMonadChainFirstEitherK is an alias for MonadChainFirstResultK.
func ExampleMonadChainFirstEitherK() {
	fmt.Println(result.Unwrap(MonadChainFirstEitherK(Of[tapCfg](3), tapPositive)(tapCfg{})))
	// Output: 3 <nil>
}

// ExampleMonadTapEitherK is an alias for MonadTapResultK.
func ExampleMonadTapEitherK() {
	fmt.Println(result.Unwrap(MonadTapEitherK(Of[tapCfg](-3), tapPositive)(tapCfg{})))
	// Output: 0 not positive
}

// ExampleTapEitherK is an alias for TapResultK.
func ExampleTapEitherK() {
	fmt.Println(result.Unwrap(F.Pipe1(Of[tapCfg](3), TapEitherK[tapCfg](tapPositive))(tapCfg{})))
	// Output: 3 <nil>
}
