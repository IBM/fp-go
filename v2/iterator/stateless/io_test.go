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

package stateless

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorFromLazy(t *testing.T) {
	num := rand.Int

	cit := FromLazy(num)

	// create arrays twice
	c1 := ToArray(cit)
	c2 := ToArray(cit)

	assert.Equal(t, 1, len(c1))
	assert.Equal(t, 1, len(c2))

	assert.NotEqual(t, c1, c2)
}
