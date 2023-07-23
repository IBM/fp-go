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

package reader

import (
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"

	"github.com/IBM/fp-go/internal/utils"
)

func TestMap(t *testing.T) {

	assert.Equal(t, 2, F.Pipe1(Of[string](1), Map[string](utils.Double))(""))
}

func TestAp(t *testing.T) {
	assert.Equal(t, 2, F.Pipe1(Of[int](utils.Double), Ap[int, int, int](Of[int](1)))(0))
}
