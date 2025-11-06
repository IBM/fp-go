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

package io

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	TST "github.com/IBM/fp-go/v2/internal/testing"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestMapSeq(t *testing.T) {
	var results []string

	handler := func(value string) IO[string] {
		return func() string {
			results = append(results, value)
			return value
		}
	}

	src := A.From("a", "b", "c")

	res := F.Pipe2(
		src,
		TraverseArraySeq(handler),
		Map(func(data []string) bool {
			return assert.Equal(t, data, results)
		}),
	)

	assert.True(t, res())
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
