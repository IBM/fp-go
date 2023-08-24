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

package types

import (
	"reflect"
	"testing"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/assert"
)

func TestArray(t *testing.T) {
	stringArray := Array(String.Validate)

	validData := A.From(
		reflect.ValueOf(A.From("a", "b", "c")),
		reflect.ValueOf(A.Empty[string]()),
		reflect.ValueOf([]string{"a", "b"}),
		reflect.ValueOf(A.From(1, 2, 3)),
	)

	for i := 0; i < len(validData); i++ {
		assert.True(t, E.IsRight(stringArray.Decode(validData[i])))
	}

}
