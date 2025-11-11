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

package constant

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	S "github.com/IBM/fp-go/v2/string"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	fa := Make[string, int]("foo")
	assert.Equal(t, fa, F.Pipe1(fa, Map[string](utils.Double)))
}

func TestOf(t *testing.T) {
	assert.Equal(t, Make[string, int](""), Of[string, int](S.Monoid)(1))
}

func TestAp(t *testing.T) {
	fab := Make[string, int]("bar")
	assert.Equal(t, Make[string, int]("foobar"), Ap[string, int, int](S.Monoid)(fab)(Make[string, func(int) int]("foo")))
}
