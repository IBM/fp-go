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

package tuple

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {

	value := MakeTuple2("Carsten", 1)

	assert.Equal(t, "Tuple2[string, int](Carsten, 1)", value.String())

}

func TestMarshal(t *testing.T) {

	value := MakeTuple3("Carsten", 1, true)

	data, err := json.Marshal(value)
	require.NoError(t, err)

	var unmarshaled Tuple3[string, int, bool]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, value, unmarshaled)
}

func TestMarshalSmallArray(t *testing.T) {

	value := `["Carsten"]`

	var unmarshaled Tuple3[string, int, bool]
	err := json.Unmarshal([]byte(value), &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, MakeTuple3("Carsten", 0, false), unmarshaled)
}
