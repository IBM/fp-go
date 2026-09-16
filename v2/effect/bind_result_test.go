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

package effect

import (
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/ioresult"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// bindResultAgeLens focuses on BindState.Age.
var bindResultAgeLens = common.MakeLens(
	func(s BindState) int { return s.Age },
	func(s BindState, age int) BindState {
		s.Age = age
		return s
	},
)

// bindResultSetAge is the setter form of bindResultAgeLens.
func bindResultSetAge(age int) func(BindState) BindState {
	return func(s BindState) BindState {
		s.Age = age
		return s
	}
}

var bindResultErr = errors.New("bind result error")

func TestBindIOResultKL(t *testing.T) {
	t.Run("refines the focused field in place", func(t *testing.T) {
		double := func(age int) ioresult.IOResult[int] { return ioresult.Of(age * 2) }

		state, err := runEffect(BindIOResultKL[TestContext](bindResultAgeLens, double)(Do[TestContext](BindState{Name: "Alice", Age: 15})), testCtx)

		assert.NoError(t, err)
		assert.Equal(t, BindState{Name: "Alice", Age: 30}, state)
	})

	t.Run("propagates the IOResult error", func(t *testing.T) {
		fail := func(int) ioresult.IOResult[int] { return ioresult.Left[int](bindResultErr) }

		_, err := runEffect(BindIOResultKL[TestContext](bindResultAgeLens, fail)(Do[TestContext](BindState{Age: 15})), testCtx)

		assert.Equal(t, bindResultErr, err)
	})
}

func TestApIOResultS(t *testing.T) {
	t.Run("binds a successful IOResult", func(t *testing.T) {
		state, err := runEffect(ApIOResultS[TestContext](bindResultSetAge, ioresult.Of(30))(Do[TestContext](BindState{Name: "Alice"})), testCtx)

		assert.NoError(t, err)
		assert.Equal(t, BindState{Name: "Alice", Age: 30}, state)
	})

	t.Run("propagates the IOResult error", func(t *testing.T) {
		_, err := runEffect(ApIOResultS[TestContext](bindResultSetAge, ioresult.Left[int](bindResultErr))(Do[TestContext](BindState{})), testCtx)

		assert.Equal(t, bindResultErr, err)
	})
}

func TestApResultS(t *testing.T) {
	t.Run("propagates the Result error", func(t *testing.T) {
		_, err := runEffect(ApResultS[TestContext](bindResultSetAge, R.Left[int](bindResultErr))(Do[TestContext](BindState{})), testCtx)

		assert.Equal(t, bindResultErr, err)
	})
}

func TestApResultSL(t *testing.T) {
	t.Run("sets the focused field from a successful Result", func(t *testing.T) {
		state, err := runEffect(ApResultSL[TestContext](bindResultAgeLens, R.Of(30))(Do[TestContext](BindState{Name: "Alice"})), testCtx)

		assert.NoError(t, err)
		assert.Equal(t, BindState{Name: "Alice", Age: 30}, state)
	})

	t.Run("propagates the Result error", func(t *testing.T) {
		_, err := runEffect(ApResultSL[TestContext](bindResultAgeLens, R.Left[int](bindResultErr))(Do[TestContext](BindState{})), testCtx)

		assert.Equal(t, bindResultErr, err)
	})
}

func TestApIOResultSL(t *testing.T) {
	t.Run("sets the focused field from a successful IOResult", func(t *testing.T) {
		state, err := runEffect(ApIOResultSL[TestContext](bindResultAgeLens, ioresult.Of(30))(Do[TestContext](BindState{Name: "Alice"})), testCtx)

		assert.NoError(t, err)
		assert.Equal(t, BindState{Name: "Alice", Age: 30}, state)
	})

	t.Run("propagates the IOResult error", func(t *testing.T) {
		_, err := runEffect(ApIOResultSL[TestContext](bindResultAgeLens, ioresult.Left[int](bindResultErr))(Do[TestContext](BindState{})), testCtx)

		assert.Equal(t, bindResultErr, err)
	})
}
