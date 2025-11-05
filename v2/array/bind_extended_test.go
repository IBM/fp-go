// Copyright (c) 2023 IBM Corp.
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

package array

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

type TestState1 struct {
	X int
}

type TestState2 struct {
	X int
	Y int
}

func TestLet(t *testing.T) {
	result := F.Pipe2(
		Do(TestState1{}),
		Let(
			func(y int) func(s TestState1) TestState2 {
				return func(s TestState1) TestState2 {
					return TestState2{X: s.X, Y: y}
				}
			},
			func(s TestState1) int { return s.X * 2 },
		),
		Map(func(s TestState2) int { return s.X + s.Y }),
	)

	assert.Equal(t, []int{0}, result)
}

func TestLetTo(t *testing.T) {
	result := F.Pipe2(
		Do(TestState1{X: 5}),
		LetTo(
			func(y int) func(s TestState1) TestState2 {
				return func(s TestState1) TestState2 {
					return TestState2{X: s.X, Y: y}
				}
			},
			42,
		),
		Map(func(s TestState2) int { return s.X + s.Y }),
	)

	assert.Equal(t, []int{47}, result)
}

func TestBindTo(t *testing.T) {
	result := F.Pipe1(
		[]int{1, 2, 3},
		BindTo(func(x int) TestState1 {
			return TestState1{X: x}
		}),
	)

	expected := []TestState1{{X: 1}, {X: 2}, {X: 3}}
	assert.Equal(t, expected, result)
}
