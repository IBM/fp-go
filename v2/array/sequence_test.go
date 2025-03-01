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

package array

import (
	"testing"

	"github.com/stretchr/testify/assert"

	O "github.com/IBM/fp-go/v2/option"
)

func TestSequenceOption(t *testing.T) {
	seq := ArrayOption[int]()

	assert.Equal(t, O.Of([]int{1, 3}), seq([]O.Option[int]{O.Of(1), O.Of(3)}))
	assert.Equal(t, O.None[[]int](), seq([]O.Option[int]{O.Of(1), O.None[int]()}))
}
