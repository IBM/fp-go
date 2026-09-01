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

package option

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

// testCounter is a simple struct used in lens-based operation tests.
type testCounter struct {
	Value int
	Name  string
}

// counterValueLens focuses on the Value field of testCounter.
var counterValueLens = L.MakeLens(
	func(c testCounter) int { return c.Value },
	func(c testCounter, v int) testCounter { c.Value = v; return c },
)

func getLastName(s utils.Initial) Option[string] {
	return Of("Doe")
}

func getGivenName(s utils.WithLastName) Option[string] {
	return Of("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res, Of("John Doe"))
}

func TestApSL(t *testing.T) {
	t.Run("Some struct, Some value - sets field", func(t *testing.T) {
		res := F.Pipe1(
			Do(testCounter{Name: "test"}),
			ApSL(counterValueLens, Some(42)),
		)
		assert.Equal(t, Some(testCounter{Value: 42, Name: "test"}), res)
	})

	t.Run("Some struct, None value - returns None", func(t *testing.T) {
		res := F.Pipe1(
			Do(testCounter{Value: 10, Name: "test"}),
			ApSL(counterValueLens, None[int]()),
		)
		assert.Equal(t, None[testCounter](), res)
	})

	t.Run("None struct, Some value - returns None", func(t *testing.T) {
		res := F.Pipe1(
			None[testCounter](),
			ApSL(counterValueLens, Some(42)),
		)
		assert.Equal(t, None[testCounter](), res)
	})
}

func TestBindL(t *testing.T) {
	increment := func(v int) Option[int] {
		if v < 100 {
			return Some(v + 1)
		}
		return None[int]()
	}

	t.Run("increments field value", func(t *testing.T) {
		res := F.Pipe1(
			Some(testCounter{Value: 42, Name: "test"}),
			BindL(counterValueLens, increment),
		)
		assert.Equal(t, Some(testCounter{Value: 43, Name: "test"}), res)
	})

	t.Run("returns None when computation fails", func(t *testing.T) {
		res := F.Pipe1(
			Some(testCounter{Value: 100, Name: "test"}),
			BindL(counterValueLens, increment),
		)
		assert.Equal(t, None[testCounter](), res)
	})

	t.Run("returns None for None input", func(t *testing.T) {
		res := F.Pipe1(
			None[testCounter](),
			BindL(counterValueLens, increment),
		)
		assert.Equal(t, None[testCounter](), res)
	})
}

func TestLetL(t *testing.T) {
	double := func(v int) int { return v * 2 }

	t.Run("doubles field value", func(t *testing.T) {
		res := F.Pipe1(
			Some(testCounter{Value: 21, Name: "test"}),
			LetL(counterValueLens, double),
		)
		assert.Equal(t, Some(testCounter{Value: 42, Name: "test"}), res)
	})

	t.Run("returns None for None input", func(t *testing.T) {
		res := F.Pipe1(
			None[testCounter](),
			LetL(counterValueLens, double),
		)
		assert.Equal(t, None[testCounter](), res)
	})
}

func TestLetToL(t *testing.T) {
	t.Run("sets field to constant value", func(t *testing.T) {
		res := F.Pipe1(
			Some(testCounter{Value: 42, Name: "test"}),
			LetToL(counterValueLens, 100),
		)
		assert.Equal(t, Some(testCounter{Value: 100, Name: "test"}), res)
	})

	t.Run("returns None for None input", func(t *testing.T) {
		res := F.Pipe1(
			None[testCounter](),
			LetToL(counterValueLens, 100),
		)
		assert.Equal(t, None[testCounter](), res)
	})
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		ApS(utils.SetLastName, Of("Doe")),
		ApS(utils.SetGivenName, Of("John")),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res, Of("John Doe"))
}
