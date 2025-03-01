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

package function

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	var count int

	withSideEffect := func(n int) int {
		count++
		return n
	}

	cached := Memoize(withSideEffect)

	assert.Equal(t, 0, count)

	assert.Equal(t, 10, cached(10))
	assert.Equal(t, 1, count)

	assert.Equal(t, 10, cached(10))
	assert.Equal(t, 1, count)

	assert.Equal(t, 20, cached(20))
	assert.Equal(t, 2, count)

	assert.Equal(t, 20, cached(20))
	assert.Equal(t, 2, count)

	assert.Equal(t, 10, cached(10))
	assert.Equal(t, 2, count)
}

func TestSingleElementCache(t *testing.T) {
	f := func(key string) string {
		return fmt.Sprintf("%s: %d", key, rand.Int())
	}
	cb := CacheCallback(func(s string) string { return s }, SingleElementCache[string, string]())
	cf := cb(f)

	v1 := cf("1")
	v2 := cf("1")
	v3 := cf("2")
	v4 := cf("1")

	assert.Equal(t, v1, v2)
	assert.NotEqual(t, v2, v3)
	assert.NotEqual(t, v3, v4)
	assert.NotEqual(t, v1, v4)
}
