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

package either

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestTraverse(t *testing.T) {
	f := func(n int) Option[int] {
		if n >= 2 {
			return O.Of(n)
		}
		return O.None[int]()
	}
	trav := Traverse[int](
		O.Of[Either[string, int]],
		O.Map[int, Either[string, int]],
	)(f)

	assert.Equal(t, O.Of(Left[int]("a")), F.Pipe1(Left[int]("a"), trav))
	assert.Equal(t, O.None[Either[string, int]](), F.Pipe1(Right[string](1), trav))
	assert.Equal(t, O.Of(Right[string](3)), F.Pipe1(Right[string](3), trav))
}

func TestSequence(t *testing.T) {

	seq := Sequence(
		O.Of[Either[string, int]],
		O.Map[int, Either[string, int]],
	)

	assert.Equal(t, O.Of(Right[string](1)), seq(Right[string](O.Of(1))))
	assert.Equal(t, O.Of(Left[int]("a")), seq(Left[Option[int]]("a")))
	assert.Equal(t, O.None[Either[string, int]](), seq(Right[string](O.None[int]())))
}
