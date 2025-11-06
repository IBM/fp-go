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

package io

import (
	"math/rand"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, 2, F.Pipe1(Of(1), Map(utils.Double))())
}

func TestChain(t *testing.T) {
	f := func(n int) IO[int] {
		return Of(n * 2)
	}
	assert.Equal(t, 2, F.Pipe1(Of(1), Chain(f))())
}

func TestAp(t *testing.T) {
	assert.Equal(t, 2, F.Pipe1(Of(utils.Double), Ap[int](Of(1)))())
}

func TestFlatten(t *testing.T) {
	assert.Equal(t, 1, F.Pipe1(Of(Of(1)), Flatten[int])())
}

func TestMemoize(t *testing.T) {
	data := Memoize(rand.Int)

	value1 := data()
	value2 := data()

	assert.Equal(t, value1, value2)
}

func TestApFirst(t *testing.T) {

	x := F.Pipe1(
		Of("a"),
		ApFirst[string](Of("b")),
	)

	assert.Equal(t, "a", x())
}

func TestApSecond(t *testing.T) {

	x := F.Pipe1(
		Of("a"),
		ApSecond[string](Of("b")),
	)

	assert.Equal(t, "b", x())
}
