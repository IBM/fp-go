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

package stateless

import (
	"testing"

	F "github.com/IBM/fp-go/function"
	T "github.com/IBM/fp-go/tuple"
	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {

	src := From("a", "b", "c")

	dst := F.Pipe1(
		src,
		Scan(func(cur T.Tuple2[int, string], val string) T.Tuple2[int, string] {
			return T.MakeTuple2(cur.F1+1, val)
		}, T.MakeTuple2(0, "")),
	)

	assert.Equal(t, ToArray(From(
		T.MakeTuple2(1, "a"),
		T.MakeTuple2(2, "b"),
		T.MakeTuple2(3, "c"),
	)), ToArray(dst))
}
