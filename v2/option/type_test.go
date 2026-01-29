// Copyright (c) 2025 IBM Corp.
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

	"github.com/stretchr/testify/assert"
)

func TestTypeConversion(t *testing.T) {

	var src any = "Carsten"

	dst := InstanceOf[string](src)
	assert.Equal(t, Some("Carsten"), dst)
}

func TestInvalidConversion(t *testing.T) {
	var src any = make(map[string]string)

	dst := InstanceOf[int](src)
	assert.Equal(t, None[int](), dst)
}
