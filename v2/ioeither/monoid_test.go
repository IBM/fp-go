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

package ioeither

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestApplicativeMonoid(t *testing.T) {
	m := ApplicativeMonoid[error](S.Monoid)

	// good cases
	assert.Equal(t, E.Of[error]("ab"), m.Concat(Of[error]("a"), Of[error]("b"))())
	assert.Equal(t, E.Of[error]("a"), m.Concat(Of[error]("a"), m.Empty())())
	assert.Equal(t, E.Of[error]("b"), m.Concat(m.Empty(), Of[error]("b"))())

	// bad cases
	e1 := fmt.Errorf("e1")
	e2 := fmt.Errorf("e1")

	assert.Equal(t, E.Left[string](e1), m.Concat(Left[string](e1), Of[error]("b"))())
	assert.Equal(t, E.Left[string](e1), m.Concat(Left[string](e1), Left[string](e2))())
	assert.Equal(t, E.Left[string](e2), m.Concat(Of[error]("a"), Left[string](e2))())
}
