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
	"testing"

	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsNone(t *testing.T) {
	assert.True(t, IsNone(None[int]()))
	assert.False(t, IsNone(Of(1)))
}

func TestIsSome(t *testing.T) {
	assert.True(t, IsSome(Of(1)))
	assert.False(t, IsSome(None[int]()))
}

func TestMapOption(t *testing.T) {

	AssertEq(Map(utils.Double)(Some(2)))(Some(4))(t)

	AssertEq(Map(utils.Double)(None[int]()))(None[int]())(t)
}

func TestAp(t *testing.T) {
	AssertEq(Some(4))(Ap[int](Some(2))(Some(utils.Double)))(t)
	AssertEq(None[int]())(Ap[int](None[int]())(Some(utils.Double)))(t)
	AssertEq(None[int]())(Ap[int](Some(2))(None[func(int) int]()))(t)
	AssertEq(None[int]())(Ap[int](None[int]())(None[func(int) int]()))(t)
}

// func TestChain(t *testing.T) {
// 	f := func(n int) Option[int] { return Some(n * 2) }
// 	g := func(_ int) Option[int] { return None[int]() }

// 	assert.Equal(t, Some(2), F.Pipe1(
// 		Some(1),
// 		Chain(f),
// 	))

// 	assert.Equal(t, None[int](), F.Pipe1(
// 		None[int](),
// 		Chain(f),
// 	))

// 	assert.Equal(t, None[int](), F.Pipe1(
// 		Some(1),
// 		Chain(g),
// 	))

// 	assert.Equal(t, None[int](), F.Pipe1(
// 		None[int](),
// 		Chain(g),
// 	))
// }

// func TestFlatten(t *testing.T) {
// 	assert.Equal(t, Of(1), F.Pipe1(Of(Of(1)), Flatten[int]))
// }

// func TestFold(t *testing.T) {
// 	f := F.Constant("none")
// 	g := func(s string) string { return fmt.Sprintf("some%d", len(s)) }

// 	fold := Fold(f, g)

// 	assert.Equal(t, "none", fold(None[string]()))
// 	assert.Equal(t, "some3", fold(Some("abc")))
// }

// func TestFromPredicate(t *testing.T) {
// 	p := func(n int) bool { return n > 2 }
// 	f := FromPredicate(p)

// 	assert.Equal(t, None[int](), f(1))
// 	assert.Equal(t, Some(3), f(3))
// }

// func TestAlt(t *testing.T) {
// 	assert.Equal(t, Some(1), F.Pipe1(Some(1), Alt(F.Constant(Some(2)))))
// 	assert.Equal(t, Some(2), F.Pipe1(Some(2), Alt(F.Constant(None[int]()))))
// 	assert.Equal(t, Some(1), F.Pipe1(None[int](), Alt(F.Constant(Some(1)))))
// 	assert.Equal(t, None[int](), F.Pipe1(None[int](), Alt(F.Constant(None[int]()))))
// }
