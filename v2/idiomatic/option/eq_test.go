// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {

	r1, r1ok := Of(1)
	r2, r2ok := Of(1)
	r3, r3ok := Of(2)

	n1, n1ok := None[int]()

	eq := FromStrictEquals[int]()

	assert.True(t, eq(r1, r1ok)(r1, r1ok))
	assert.True(t, eq(r1, r1ok)(r2, r2ok))
	assert.False(t, eq(r1, r1ok)(r3, r3ok))
	assert.False(t, eq(r1, r1ok)(n1, n1ok))

	assert.True(t, eq(n1, n1ok)(n1, n1ok))
	assert.False(t, eq(n1, n1ok)(r2, r2ok))
}
