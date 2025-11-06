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

package lens

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	OPT "github.com/IBM/fp-go/v2/optics/optional"
	OPTP "github.com/IBM/fp-go/v2/optics/optional/prism"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type Inner struct {
	A int
}

type State = O.Option[*Inner]

func (inner *Inner) getA() int {
	return inner.A
}

func (inner *Inner) setA(a int) *Inner {
	inner.A = a
	return inner
}

func TestCompose(t *testing.T) {

	inner1 := Inner{1}

	lensa := L.MakeLensRef((*Inner).getA, (*Inner).setA)

	sa := F.Pipe1(
		OPT.Id[State](),
		OPTP.Some[State, *Inner],
	)
	ab := F.Pipe1(
		L.IdRef[Inner](),
		L.ComposeRef[Inner](lensa),
	)
	sb := F.Pipe1(
		sa,
		Compose[State](ab),
	)
	// check get access
	assert.Equal(t, O.None[int](), sb.GetOption(O.None[*Inner]()))
	assert.Equal(t, O.Of(1), sb.GetOption(O.Of(&inner1)))

	// check set access
	res := F.Pipe1(
		sb.Set(2)(O.Of(&inner1)),
		O.Map(func(i *Inner) int {
			return i.A
		}),
	)
	assert.Equal(t, O.Of(2), res)
	assert.Equal(t, 1, inner1.A)

	assert.Equal(t, O.None[*Inner](), sb.Set(2)(O.None[*Inner]()))

}
