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

	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// it('getOrd', () => {
//     const OS = _.getOrd(S.Ord)
//     U.deepStrictEqual(OS.compare(_.none, _.none), 0)
//     U.deepStrictEqual(OS.compare(_.some('a'), _.none), 1)
//     U.deepStrictEqual(OS.compare(_.none, _.some('a')), -1)
//     U.deepStrictEqual(OS.compare(_.some('a'), _.some('a')), 0)
//     U.deepStrictEqual(OS.compare(_.some('a'), _.some('b')), -1)
//     U.deepStrictEqual(OS.compare(_.some('b'), _.some('a')), 1)
//   })

func TestOrd(t *testing.T) {

	os := Ord(S.Ord)

	assert.Equal(t, 0, os((None[string]()))(None[string]()))
	assert.Equal(t, +1, os(Some("a"))(None[string]()))
	assert.Equal(t, -1, os(None[string]())(Some("a")))
	assert.Equal(t, 0, os(Some("a"))(Some("a")))
	assert.Equal(t, -1, os(Some("a"))(Some("b")))
	assert.Equal(t, +1, os(Some("b"))(Some("a")))

}
