//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package record

import (
	"testing"

	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestUnionMonoid(t *testing.T) {
	m := UnionMonoid[string](S.Semigroup)

	e := Empty[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	res := map[string]string{
		"a": "a1",
		"b": "b1b2",
		"c": "c1c2",
		"d": "d2",
	}

	assert.Equal(t, x, m.Concat(x, m.Empty()))
	assert.Equal(t, x, m.Concat(m.Empty(), x))

	assert.Equal(t, x, m.Concat(x, e))
	assert.Equal(t, x, m.Concat(e, x))

	assert.Equal(t, res, m.Concat(x, y))
}

func TestUnionFirstMonoid(t *testing.T) {
	m := UnionFirstMonoid[string, string]()

	e := Empty[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	res := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
		"d": "d2",
	}

	assert.Equal(t, x, m.Concat(x, m.Empty()))
	assert.Equal(t, x, m.Concat(m.Empty(), x))

	assert.Equal(t, x, m.Concat(x, e))
	assert.Equal(t, x, m.Concat(e, x))

	assert.Equal(t, res, m.Concat(x, y))
}

func TestUnionLastMonoid(t *testing.T) {
	m := UnionLastMonoid[string, string]()

	e := Empty[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	res := map[string]string{
		"a": "a1",
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	assert.Equal(t, x, m.Concat(x, m.Empty()))
	assert.Equal(t, x, m.Concat(m.Empty(), x))

	assert.Equal(t, x, m.Concat(x, e))
	assert.Equal(t, x, m.Concat(e, x))

	assert.Equal(t, res, m.Concat(x, y))
}
