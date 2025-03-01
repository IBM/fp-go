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

package json

import (
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

type TestType struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func TestToType(t *testing.T) {

	generic := map[string]any{"a": "value", "b": 1}

	assert.True(t, E.IsRight(ToTypeE[TestType](generic)))
	assert.True(t, E.IsRight(ToTypeE[TestType](&generic)))

	assert.Equal(t, E.Right[error](&TestType{A: "value", B: 1}), ToTypeE[*TestType](&generic))
	assert.Equal(t, E.Right[error](TestType{A: "value", B: 1}), F.Pipe1(ToTypeE[*TestType](&generic), E.Map[error](F.Deref[TestType])))
}
