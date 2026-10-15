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

package ioresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// tapPositive fails for numbers that are not positive.
func tapPositive(n int) Result[int] {
	if n > 0 {
		return result.Of(n)
	}
	return result.Left[int](errors.New("not positive"))
}

func TestChainFirstResultK(t *testing.T) {
	all := func(start IOResult[int]) []IOResult[int] {
		return []IOResult[int]{
			ChainFirstResultK(tapPositive)(start),
			TapResultK(tapPositive)(start),
			MonadChainFirstResultK(start, tapPositive),
			MonadTapResultK(start, tapPositive),
		}
	}
	t.Run("keeps the value", func(t *testing.T) {
		for _, m := range all(Of(7)) {
			assert.Equal(t, result.Of(7), m())
		}
	})
	t.Run("failing check replaces the value", func(t *testing.T) {
		for _, m := range all(Of(-7)) {
			assert.Equal(t, result.Left[int](errors.New("not positive")), m())
		}
	})
	t.Run("upstream error is preserved", func(t *testing.T) {
		err := errors.New("upstream")
		for _, m := range all(Left[int](err)) {
			assert.Equal(t, result.Left[int](err), m())
		}
	})
}

// ExampleChainFirstResultK validates a value without changing it.
func ExampleChainFirstResultK() {
	fmt.Println(result.Unwrap(F.Pipe1(Of(3), ChainFirstResultK(tapPositive))()))
	fmt.Println(result.Unwrap(F.Pipe1(Of(-3), ChainFirstResultK(tapPositive))()))
	// Output:
	// 3 <nil>
	// 0 not positive
}

// ExampleTapResultK is an alias for ChainFirstResultK.
func ExampleTapResultK() {
	fmt.Println(result.Unwrap(F.Pipe1(Of(3), TapResultK(tapPositive))()))
	// Output: 3 <nil>
}
