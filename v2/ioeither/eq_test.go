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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {

	r1 := Of[string](1)
	r2 := Of[string](1)
	r3 := Of[string](2)

	e1 := Left[int]("a")
	e2 := Left[int]("a")
	e3 := Left[int]("b")

	eq := FromStrictEquals[string, int]()

	assert.True(t, eq.Equals(r1, r1))
	assert.True(t, eq.Equals(r1, r2))
	assert.False(t, eq.Equals(r1, r3))
	assert.False(t, eq.Equals(r1, e1))

	assert.True(t, eq.Equals(e1, e1))
	assert.True(t, eq.Equals(e1, e2))
	assert.False(t, eq.Equals(e1, e3))
	assert.False(t, eq.Equals(e2, r2))
}
