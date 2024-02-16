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

package ioeither

import (
	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestMapSeq(t *testing.T) {
	var results []string

	handler := func(value string) IOEither[error, string] {
		return func() E.Either[error, string] {
			results = append(results, value)
			return E.Of[error](value)
		}
	}

	src := A.From("a", "b", "c")

	res := F.Pipe2(
		src,
		TraverseArraySeq(handler),
		Map[error](func(data []string) bool {
			return assert.Equal(t, data, results)
		}),
	)

	assert.Equal(t, E.Of[error](true), res())
}
