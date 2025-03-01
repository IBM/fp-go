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

package string

import (
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	assert.True(t, IsEmpty(""))
	assert.False(t, IsEmpty("Carsten"))
}

func TestJoin(t *testing.T) {

	x := Join(",")(A.From("a", "b", "c"))
	assert.Equal(t, x, x)

	assert.Equal(t, "a,b,c", Join(",")(A.From("a", "b", "c")))
	assert.Equal(t, "a", Join(",")(A.From("a")))
	assert.Equal(t, "", Join(",")(A.Empty[string]()))
}

func TestEquals(t *testing.T) {
	assert.True(t, Equals("a")("a"))
	assert.False(t, Equals("a")("b"))
	assert.False(t, Equals("b")("a"))
}

func TestIncludes(t *testing.T) {
	assert.True(t, Includes("a")("bab"))
	assert.False(t, Includes("bab")("a"))
	assert.False(t, Includes("b")("a"))
}
