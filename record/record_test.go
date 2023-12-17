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

package record

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/internal/utils"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	data := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	keys := Keys(data)
	sort.Strings(keys)

	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestValues(t *testing.T) {
	data := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	keys := Values(data)
	sort.Strings(keys)

	assert.Equal(t, []string{"A", "B", "C"}, keys)
}

func TestMap(t *testing.T) {
	data := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	expected := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	assert.Equal(t, expected, Map[string](utils.Upper)(data))
}

func TestLookup(t *testing.T) {
	data := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	assert.Equal(t, O.Some("a"), Lookup[string]("a")(data))
	assert.Equal(t, O.None[string](), Lookup[string]("a1")(data))
}

func TestFilterChain(t *testing.T) {
	src := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	f := func(k string, value int) O.Option[map[string]string] {
		if value%2 != 0 {
			return O.Of(map[string]string{
				k: fmt.Sprintf("%s%d", k, value),
			})
		}
		return O.None[map[string]string]()
	}

	// monoid
	monoid := MergeMonoid[string, string]()

	res := FilterChainWithIndex[int](monoid)(f)(src)

	assert.Equal(t, map[string]string{
		"a": "a1",
		"c": "c3",
	}, res)
}

func ExampleFoldMap() {
	src := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}

	fold := FoldMapOrd[string, string](S.Ord)(S.Monoid)(strings.ToUpper)

	fmt.Println(fold(src))

	// Output: ABC

}

func ExampleValuesOrd() {
	src := map[string]string{
		"c": "a",
		"b": "b",
		"a": "c",
	}

	getValues := ValuesOrd[string](S.Ord)

	fmt.Println(getValues(src))

	// Output: [c b a]

}

func TestCopyVsClone(t *testing.T) {
	slc := []string{"b", "c"}
	src := map[string][]string{
		"a": slc,
	}
	// make a shallow copy
	cpy := Copy(src)
	// make a deep copy
	cln := Clone[string](A.Copy[string])(src)

	assert.Equal(t, cpy, cln)
	// make a modification to the original slice
	slc[0] = "d"
	assert.NotEqual(t, cpy, cln)
	assert.Equal(t, src, cpy)
}
