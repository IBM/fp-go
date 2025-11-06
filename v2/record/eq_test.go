// Copyright (c) 2024 - 2025 IBM Corp.
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

package record

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromStrictEquals(t *testing.T) {
	m1 := map[string]string{
		"a": "A",
		"b": "B",
	}
	m2 := map[string]string{
		"a": "A",
		"b": "C",
	}
	m3 := map[string]string{
		"a": "A",
		"b": "B",
	}
	m4 := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}

	e := FromStrictEquals[string, string]()
	assert.True(t, e.Equals(m1, m1))
	assert.True(t, e.Equals(m1, m3))
	assert.False(t, e.Equals(m1, m2))
	assert.False(t, e.Equals(m1, m4))
}
