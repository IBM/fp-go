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

package itereither

import (
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/stretchr/testify/assert"
)

func TestMonadAp(t *testing.T) {
	t.Run("applies function to value", func(t *testing.T) {
		fab := iter.From(E.Right[string](func(x int) int { return x * 2 }))
		fa := iter.From(E.Right[string](21))
		result := collectEithers(MonadAp(fab, fa))
		assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
	})

	t.Run("preserves Left in function", func(t *testing.T) {
		fab := iter.From(E.Left[func(int) int]("error"))
		fa := iter.From(E.Right[string](21))
		result := collectEithers(MonadAp(fab, fa))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error")}, result)
	})

	t.Run("preserves Left in value", func(t *testing.T) {
		fab := iter.From(E.Right[string](func(x int) int { return x * 2 }))
		fa := iter.From(E.Left[int]("error"))
		result := collectEithers(MonadAp(fab, fa))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error")}, result)
	})
}

func TestAp(t *testing.T) {
	fab := iter.From(E.Right[string](func(x int) int { return x * 2 }))
	fa := iter.From(E.Right[string](21))
	result := F.Pipe1(fab, Ap[int](fa))
	assert.Equal(t, []Either[string, int]{E.Right[string](42)}, collectEithers(result))
}

func TestMonadApFirst(t *testing.T) {
	t.Run("both Right values - returns first", func(t *testing.T) {
		first := iter.From(E.Right[string](10))
		second := iter.From(E.Right[string](20))
		result := collectEithers(MonadApFirst(first, second))
		assert.Equal(t, []Either[string, int]{E.Right[string](10)}, result)
	})

	t.Run("first Left - returns Left", func(t *testing.T) {
		first := iter.From(E.Left[int]("error1"))
		second := iter.From(E.Right[string](20))
		result := collectEithers(MonadApFirst(first, second))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error1")}, result)
	})

	t.Run("second Left - returns Left", func(t *testing.T) {
		first := iter.From(E.Right[string](10))
		second := iter.From(E.Left[int]("error2"))
		result := collectEithers(MonadApFirst(first, second))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error2")}, result)
	})

	t.Run("both Left - returns first Left", func(t *testing.T) {
		first := iter.From(E.Left[int]("error1"))
		second := iter.From(E.Left[int]("error2"))
		result := collectEithers(MonadApFirst(first, second))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error1")}, result)
	})
}

func TestApFirst(t *testing.T) {
	t.Run("composition with Map", func(t *testing.T) {
		result := F.Pipe2(
			iter.From(E.Right[string](10)),
			ApFirst[int](iter.From(E.Right[string](20))),
			Map[string](func(x int) int { return x * 2 }),
		)
		assert.Equal(t, []Either[string, int]{E.Right[string](20)}, collectEithers(result))
	})

	t.Run("with different types", func(t *testing.T) {
		result := F.Pipe1(
			iter.From(E.Right[string]("text")),
			ApFirst[string](iter.From(E.Right[string](42))),
		)
		assert.Equal(t, []Either[string, string]{E.Right[string]("text")}, collectEithers(result))
	})
}

func TestMonadApSecond(t *testing.T) {
	t.Run("both Right values - returns second", func(t *testing.T) {
		first := iter.From(E.Right[string](10))
		second := iter.From(E.Right[string](20))
		result := collectEithers(MonadApSecond(first, second))
		assert.Equal(t, []Either[string, int]{E.Right[string](20)}, result)
	})

	t.Run("first Left - returns Left", func(t *testing.T) {
		first := iter.From(E.Left[int]("error1"))
		second := iter.From(E.Right[string](20))
		result := collectEithers(MonadApSecond(first, second))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error1")}, result)
	})

	t.Run("second Left - returns Left", func(t *testing.T) {
		first := iter.From(E.Right[string](10))
		second := iter.From(E.Left[int]("error2"))
		result := collectEithers(MonadApSecond(first, second))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error2")}, result)
	})

	t.Run("both Left - returns first Left", func(t *testing.T) {
		first := iter.From(E.Left[int]("error1"))
		second := iter.From(E.Left[int]("error2"))
		result := collectEithers(MonadApSecond(first, second))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error1")}, result)
	})
}

func TestApSecond(t *testing.T) {
	t.Run("composition with Map", func(t *testing.T) {
		result := F.Pipe2(
			iter.From(E.Right[string](10)),
			ApSecond[int](iter.From(E.Right[string](20))),
			Map[string](func(x int) int { return x * 2 }),
		)
		assert.Equal(t, []Either[string, int]{E.Right[string](40)}, collectEithers(result))
	})

	t.Run("with different types", func(t *testing.T) {
		result := F.Pipe1(
			iter.From(E.Right[string]("text")),
			ApSecond[string](iter.From(E.Right[string](42))),
		)
		assert.Equal(t, []Either[string, int]{E.Right[string](42)}, collectEithers(result))
	})
}

func TestApSequence(t *testing.T) {
	t.Run("sequence of operations", func(t *testing.T) {
		result := F.Pipe3(
			iter.From(E.Right[string](5)),
			ApFirst[int](iter.From(E.Right[string](10))),
			ApSecond[int](iter.From(E.Right[string](15))),
			Map[string](func(x int) int { return x * 2 }),
		)
		assert.Equal(t, []Either[string, int]{E.Right[string](30)}, collectEithers(result))
	})
}

// Made with Bob
