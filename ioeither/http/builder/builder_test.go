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

package builder

import (
	"testing"

	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"github.com/stretchr/testify/assert"
)

func TestBuiler(t *testing.T) {

	name := "Content-type"
	withContentType := WithHeader(name)
	withoutContentType := WithoutHeader(name)

	b1 := F.Pipe1(
		Default,
		withContentType("application/json"),
	)

	b2 := F.Pipe1(
		b1,
		withContentType("text/plain"),
	)

	b3 := F.Pipe1(
		b2,
		withoutContentType,
	)

	assert.Equal(t, O.None[string](), Default.GetHeader(name))
	assert.Equal(t, O.Of("application/json"), b1.GetHeader(name))
	assert.Equal(t, O.Of("text/plain"), b2.GetHeader(name))
	assert.Equal(t, O.None[string](), b3.GetHeader(name))
}
