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

	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

func TestSequenceT(t *testing.T) {
	// one argumemt
	s1 := SequenceT1[int]
	assert.Equal(t, Of(T.MakeTuple1(1)), s1(Of(1)))

	// two arguments
	s2 := SequenceT2[int, string]
	assert.Equal(t, Of(T.MakeTuple2(1, "a")), s2(Of(1), Of("a")))

	// three arguments
	s3 := SequenceT3[int, string, bool]
	assert.Equal(t, Of(T.MakeTuple3(1, "a", true)), s3(Of(1), Of("a"), Of(true)))

	// four arguments
	s4 := SequenceT4[int, string, bool, int]
	assert.Equal(t, Of(T.MakeTuple4(1, "a", true, 2)), s4(Of(1), Of("a"), Of(true), Of(2)))

	// three with one none
	assert.Equal(t, None[T.Tuple3[int, string, bool]](), s3(Of(1), Of("a"), None[bool]()))
}
