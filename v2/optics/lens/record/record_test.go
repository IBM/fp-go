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

package record

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type (
	S = map[string]int
)

func TestAtKey(t *testing.T) {
	sa := F.Pipe1(
		L.Id[S](),
		AtKey[S, int]("a"),
	)

	assert.Equal(t, O.Some(1), sa.Get(S{"a": 1}))
	assert.Equal(t, S{"a": 2, "b": 2}, sa.Set(O.Some(2))(S{"a": 1, "b": 2}))
	assert.Equal(t, S{"a": 1, "b": 2}, sa.Set(O.Some(1))(S{"b": 2}))
	assert.Equal(t, S{"b": 2}, sa.Set(O.None[int]())(S{"a": 1, "b": 2}))
}
