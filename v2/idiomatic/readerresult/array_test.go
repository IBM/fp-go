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

package readerresult

import (
	"context"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestSequenceArray(t *testing.T) {

	n := 10

	readers := A.MakeBy(n, Of[context.Context, int])
	exp := A.MakeBy(n, F.Identity[int])

	g := F.Pipe1(
		readers,
		SequenceArray[context.Context, int],
	)

	v, err := g(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, exp, v)
}
