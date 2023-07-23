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

package option

import (
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestSequenceArray(t *testing.T) {

	one := Of(1)
	two := Of(2)

	res := F.Pipe1(
		[]Option[int]{one, two},
		SequenceArray[int],
	)

	assert.Equal(t, res, Of([]int{1, 2}))
}
