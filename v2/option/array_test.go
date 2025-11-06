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
	"fmt"
	"testing"

	TST "github.com/IBM/fp-go/v2/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestCompactArray(t *testing.T) {
	ar := []Option[string]{
		Of("ok"),
		None[string](),
		Of("ok"),
	}

	res := CompactArray(ar)
	assert.Equal(t, 2, len(res))
}

func TestSequenceArray(t *testing.T) {

	s := TST.SequenceArrayTest(
		FromStrictEquals[bool](),
		Pointed[string](),
		Pointed[bool](),
		Functor[[]string, bool](),
		SequenceArray[string],
	)

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("TestSequenceArray %d", i), s(i))
	}
}
