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

package readerioeither

import (
	"context"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

func TestSequence2(t *testing.T) {
	// two readers of heterogeneous types
	first := Of[context.Context, error]("a")
	second := Of[context.Context, error](1)

	// compose
	s2 := SequenceT2[context.Context, error, string, int]
	res := s2(first, second)

	ctx := context.Background()
	assert.Equal(t, either.Right[error](T.MakeTuple2("a", 1)), res(ctx)())
}
