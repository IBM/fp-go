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

	A "github.com/IBM/fp-go/v2/array"
	"github.com/stretchr/testify/assert"
)

func TestDropWhile(t *testing.T) {
	// sequence of 5 items
	data := Take[int](10)(Cycle(From(0, 1, 2, 3)))

	total := DropWhile(func(data int) bool {
		return data <= 2
	})(data)

	assert.Equal(t, A.From(3, 0, 1, 2, 3, 0, 1), ToArray(total))

}
