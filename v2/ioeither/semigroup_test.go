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

package ioeither

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/assert"
)

func TestAltSemigroup(t *testing.T) {
	s := AltSemigroup[error, int]()

	t.Run("first Right returns first", func(t *testing.T) {
		first := Of[error](1)
		second := Of[error](2)
		result := s.Concat(first, second)()
		assert.Equal(t, E.Of[error](1), result)
	})

	t.Run("first Left tries second Right", func(t *testing.T) {
		first := Left[int](errors.New("error1"))
		second := Of[error](2)
		result := s.Concat(first, second)()
		assert.Equal(t, E.Of[error](2), result)
	})

	t.Run("both Left returns second Left", func(t *testing.T) {
		err1 := errors.New("error1")
		err2 := errors.New("error2")
		first := Left[int](err1)
		second := Left[int](err2)
		result := s.Concat(first, second)()
		assert.True(t, E.IsLeft(result))
		_, leftVal := E.Unwrap(result)
		assert.Equal(t, err2, leftVal)
	})

	t.Run("chaining multiple alternatives", func(t *testing.T) {
		first := Left[int](errors.New("error1"))
		second := Left[int](errors.New("error2"))
		third := Of[error](3)
		result := s.Concat(s.Concat(first, second), third)()
		assert.Equal(t, E.Of[error](3), result)
	})

	t.Run("first Right short-circuits", func(t *testing.T) {
		first := Of[error](1)
		second := Of[error](2)

		// When first succeeds, it returns immediately
		result := s.Concat(first, second)()
		assert.Equal(t, E.Of[error](1), result)
	})
}
