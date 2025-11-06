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

package function

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fromLibrary(data ...string) string {
	return strings.Join(data, "-")
}

func TestUnvariadic(t *testing.T) {

	res := Pipe1(
		[]string{"A", "B"},
		Unvariadic0(fromLibrary),
	)

	assert.Equal(t, "A-B", res)
}

func TestVariadicArity(t *testing.T) {

	f := Unsliced2(Unvariadic0(fromLibrary))

	res := f("A", "B")
	assert.Equal(t, "A-B", res)
}
