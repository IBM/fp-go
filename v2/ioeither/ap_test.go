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
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestMonadApFirstExtended(t *testing.T) {
	t.Run("both Right values - returns first", func(t *testing.T) {
		first := Of[error]("first")
		second := Of[error]("second")
		result := MonadApFirst(first, second)
		assert.Equal(t, E.Of[error]("first"), result())
	})

	t.Run("first Left - returns Left", func(t *testing.T) {
		first := Left[string](errors.New("error1"))
		second := Of[error]("second")
		result := MonadApFirst(first, second)
		assert.True(t, E.IsLeft(result()))
	})

	t.Run("second Left - returns Left", func(t *testing.T) {
		first := Of[error]("first")
		second := Left[string](errors.New("error2"))
		result := MonadApFirst(first, second)
		assert.True(t, E.IsLeft(result()))
	})

	t.Run("both Left - returns first Left", func(t *testing.T) {
		first := Left[string](errors.New("error1"))
		second := Left[string](errors.New("error2"))
		result := MonadApFirst(first, second)
		assert.True(t, E.IsLeft(result()))
	})
}

func TestApFirstExtended(t *testing.T) {
	t.Run("composition with Map", func(t *testing.T) {
		result := F.Pipe2(
			Of[error](10),
			ApFirst[int](Of[error](20)),
			Map[error](func(x int) int { return x * 2 }),
		)
		assert.Equal(t, E.Of[error](20), result())
	})

	t.Run("with different types", func(t *testing.T) {
		result := F.Pipe1(
			Of[error]("text"),
			ApFirst[string](Of[error](42)),
		)
		assert.Equal(t, E.Of[error]("text"), result())
	})
}

func TestMonadApSecondExtended(t *testing.T) {
	t.Run("both Right values - returns second", func(t *testing.T) {
		first := Of[error]("first")
		second := Of[error]("second")
		result := MonadApSecond(first, second)
		assert.Equal(t, E.Of[error]("second"), result())
	})

	t.Run("first Left - returns Left", func(t *testing.T) {
		first := Left[string](errors.New("error1"))
		second := Of[error]("second")
		result := MonadApSecond(first, second)
		assert.True(t, E.IsLeft(result()))
	})

	t.Run("second Left - returns Left", func(t *testing.T) {
		first := Of[error]("first")
		second := Left[string](errors.New("error2"))
		result := MonadApSecond(first, second)
		assert.True(t, E.IsLeft(result()))
	})

	t.Run("both Left - returns first Left", func(t *testing.T) {
		first := Left[string](errors.New("error1"))
		second := Left[string](errors.New("error2"))
		result := MonadApSecond(first, second)
		assert.True(t, E.IsLeft(result()))
	})
}

func TestApSecondExtended(t *testing.T) {
	t.Run("composition with Map", func(t *testing.T) {
		result := F.Pipe2(
			Of[error](10),
			ApSecond[int](Of[error](20)),
			Map[error](func(x int) int { return x * 2 }),
		)
		assert.Equal(t, E.Of[error](40), result())
	})

	t.Run("sequence of operations", func(t *testing.T) {
		result := F.Pipe3(
			Of[error](1),
			ApSecond[int](Of[error](2)),
			ApSecond[int](Of[error](3)),
			ApSecond[int](Of[error](4)),
		)
		assert.Equal(t, E.Of[error](4), result())
	})

	t.Run("with different types", func(t *testing.T) {
		result := F.Pipe1(
			Of[error]("text"),
			ApSecond[string](Of[error](42)),
		)
		assert.Equal(t, E.Of[error](42), result())
	})
}
