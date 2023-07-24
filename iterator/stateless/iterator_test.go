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
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {

	result := F.Pipe2(
		A.From(1, 2, 3),
		FromArray[int],
		Reduce(utils.Sum, 0),
	)

	assert.Equal(t, 6, result)
}

func TestChain(t *testing.T) {

	outer := FromArray[int](A.From(1, 2, 3))

	inner := func(data int) Iterator[string] {
		return F.Pipe2(
			A.From(0, 1),
			FromArray[int],
			Map(func(idx int) string {
				return fmt.Sprintf("item[%d][%d]", data, idx)
			}),
		)
	}

	total := F.Pipe2(
		outer,
		Chain(inner),
		ToArray[string],
	)

	assert.Equal(t, A.From("item[1][0]", "item[1][1]", "item[2][0]", "item[2][1]", "item[3][0]", "item[3][1]"), total)
}
