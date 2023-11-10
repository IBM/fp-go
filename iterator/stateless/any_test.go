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

package stateless

import (
	"testing"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {

	anyBool := Any(F.Identity[bool])

	i1 := FromArray(A.From(false, true, false))
	assert.True(t, anyBool(i1))

	i2 := FromArray(A.From(false, false, false))
	assert.False(t, anyBool(i2))

	i3 := Empty[bool]()
	assert.False(t, anyBool(i3))
}
