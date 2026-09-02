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

package either

import (
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) Either[error, string] {
	return Of[error]("Doe")
}

func getGivenName(s utils.WithLastName) Either[error, string] {
	return Of[error]("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do[error](utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[error](utils.GetFullName),
	)

	assert.Equal(t, res, Of[error]("John Doe"))
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do[error](utils.Empty),
		ApS(utils.SetLastName, Of[error]("Doe")),
		ApS(utils.SetGivenName, Of[error]("John")),
		Map[error](utils.GetFullName),
	)

	assert.Equal(t, res, Of[error]("John Doe"))
}

// bindTestState is a simple struct used by the lens-based bind tests.
type bindTestState struct {
	Name  string
	Value int
}

var (
	nameLens  = L.MakeLens(func(s bindTestState) string { return s.Name }, func(s bindTestState, v string) bindTestState { s.Name = v; return s })
	valueLens = L.MakeLens(func(s bindTestState) int { return s.Value }, func(s bindTestState, v int) bindTestState { s.Value = v; return s })
)

// TestApSL verifies that ApSL sets a field via a lens using applicative semantics.
func TestApSL(t *testing.T) {
	t.Run("ApSL with Right value", func(t *testing.T) {
		result := F.Pipe1(
			Of[error](bindTestState{Name: "Alice"}),
			ApSL[error](valueLens, Of[error](42)),
		)
		assert.Equal(t, Of[error](bindTestState{Name: "Alice", Value: 42}), result)
	})

	t.Run("ApSL with Left in context", func(t *testing.T) {
		err := errors.New("context error")
		result := F.Pipe1(
			Left[bindTestState](err),
			ApSL[error](valueLens, Of[error](42)),
		)
		assert.Equal(t, Left[bindTestState](err), result)
	})

	t.Run("ApSL with Left in value", func(t *testing.T) {
		err := errors.New("value error")
		result := F.Pipe1(
			Of[error](bindTestState{Name: "Alice"}),
			ApSL[error](valueLens, Left[int](err)),
		)
		assert.Equal(t, Left[bindTestState](err), result)
	})
}

// TestBindL verifies that BindL reads the focused field, applies f, and writes the result back.
func TestBindL(t *testing.T) {
	increment := func(v int) Either[error, int] {
		if v >= 100 {
			return Left[int](errors.New("overflow"))
		}
		return Of[error](v + 1)
	}

	t.Run("BindL with successful transformation", func(t *testing.T) {
		result := F.Pipe1(
			Of[error](bindTestState{Name: "x", Value: 41}),
			BindL[error](valueLens, increment),
		)
		assert.Equal(t, Of[error](bindTestState{Name: "x", Value: 42}), result)
	})

	t.Run("BindL with error from f", func(t *testing.T) {
		result := F.Pipe1(
			Of[error](bindTestState{Name: "x", Value: 100}),
			BindL[error](valueLens, increment),
		)
		assert.True(t, IsLeft(result))
	})

	t.Run("BindL with Left context", func(t *testing.T) {
		err := errors.New("upstream error")
		result := F.Pipe1(
			Left[bindTestState](err),
			BindL[error](valueLens, increment),
		)
		assert.Equal(t, Left[bindTestState](err), result)
	})
}

// TestLetL verifies that LetL applies a pure transformation to the focused field.
func TestLetL(t *testing.T) {
	double := func(v int) int { return v * 2 }

	t.Run("LetL transforms field", func(t *testing.T) {
		result := F.Pipe1(
			Of[error](bindTestState{Name: "x", Value: 21}),
			LetL[error](valueLens, double),
		)
		assert.Equal(t, Of[error](bindTestState{Name: "x", Value: 42}), result)
	})

	t.Run("LetL short-circuits on Left", func(t *testing.T) {
		err := errors.New("upstream error")
		result := F.Pipe1(
			Left[bindTestState](err),
			LetL[error](valueLens, double),
		)
		assert.Equal(t, Left[bindTestState](err), result)
	})
}

// TestLetToL verifies that LetToL replaces the focused field with a constant value.
func TestLetToL(t *testing.T) {
	t.Run("LetToL sets field to constant", func(t *testing.T) {
		result := F.Pipe1(
			Of[error](bindTestState{Name: "old", Value: 1}),
			LetToL[error](nameLens, "new"),
		)
		assert.Equal(t, Of[error](bindTestState{Name: "new", Value: 1}), result)
	})

	t.Run("LetToL short-circuits on Left", func(t *testing.T) {
		err := errors.New("upstream error")
		result := F.Pipe1(
			Left[bindTestState](err),
			LetToL[error](nameLens, "new"),
		)
		assert.Equal(t, Left[bindTestState](err), result)
	})
}
