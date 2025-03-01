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

package either

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompactRecord(t *testing.T) {
	// make the map
	m := make(map[string]Either[string, int])
	m["foo"] = Left[int]("error")
	m["bar"] = Right[string](1)
	// compact it
	m1 := CompactRecord(m)
	// check expected
	exp := map[string]int{
		"bar": 1,
	}

	assert.Equal(t, exp, m1)
}
