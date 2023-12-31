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

package number

import (
	"testing"

	"github.com/stretchr/testify/assert"

	M "github.com/IBM/fp-go/magma"
)

func TestSemigroupIsMagma(t *testing.T) {
	sum := SemigroupSum[int]()

	var magma M.Magma[int] = sum

	assert.Equal(t, sum.Concat(1, 2), magma.Concat(1, 2))
	assert.Equal(t, sum.Concat(1, 2), sum.Concat(2, 1))
}
