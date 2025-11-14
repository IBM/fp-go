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

package stateless

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {

	src := From("a", "b", "c")

	dst := F.Pipe1(
		src,
		Scan(func(cur Pair[int, string], val string) Pair[int, string] {
			return P.MakePair(P.Head(cur)+1, val)
		}, P.MakePair(0, "")),
	)

	assert.Equal(t, ToArray(From(
		P.MakePair(1, "a"),
		P.MakePair(2, "b"),
		P.MakePair(3, "c"),
	)), ToArray(dst))
}
