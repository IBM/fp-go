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

package result

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {

	r1v, r1e := Of(1)
	r2v, r2e := Of(2)

	e1v, e1e := Left[int](errorString("a"))
	e2v, e2e := Left[int](errorString("a"))
	e3v, e3e := Left[int](errorString("b"))

	eq := FromStrictEquals[int]()

	// Right values
	assert.True(t, eq(r1v, r1e)(r1v, r1e))
	assert.False(t, eq(r1v, r1e)(r2v, r2e))

	// Left values (errors)
	assert.True(t, eq(e1v, e1e)(e1v, e1e))
	assert.True(t, eq(e1v, e1e)(e2v, e2e))
	assert.False(t, eq(e1v, e1e)(e3v, e3e))

	// Mixed Left and Right
	assert.False(t, eq(r1v, r1e)(e1v, e1e))
	assert.False(t, eq(e1v, e1e)(r1v, r1e))
}

// errorString is a simple error type for testing
type errorString string

func (e errorString) Error() string {
	return string(e)
}
