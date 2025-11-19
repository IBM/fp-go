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

package ioresult

import (
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray(t *testing.T) {

	src := A.From("A", "B")

	trfrm := TraverseArrayWithIndex(func(idx int, data string) IOResult[string] {
		return Of(fmt.Sprintf("idx: %d, data: %s", idx, data))
	})

	result, err := trfrm(src)()
	assert.NoError(t, err)
	assert.Equal(t, A.From("idx: 0, data: A", "idx: 1, data: B"), result)

}
