//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache LicensVersion 2.0 (the "License");
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

package ioresult

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestApplicativeMonoid(t *testing.T) {
	m := ApplicativeMonoid(S.Monoid)

	// good cases
	assert.Equal(t, result.Of("ab"), m.Concat(Of("a"), Of("b"))())
	assert.Equal(t, result.Of("a"), m.Concat(Of("a"), m.Empty())())
	assert.Equal(t, result.Of("b"), m.Concat(m.Empty(), Of("b"))())

	// bad cases
	e1 := fmt.Errorf("e1")
	e2 := fmt.Errorf("e1")

	assert.Equal(t, result.Left[string](e1), m.Concat(Left[string](e1), Of("b"))())
	assert.Equal(t, result.Left[string](e1), m.Concat(Left[string](e1), Left[string](e2))())
	assert.Equal(t, result.Left[string](e2), m.Concat(Of("a"), Left[string](e2))())
}
