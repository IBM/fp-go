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

package file

import (
	"os"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	"github.com/stretchr/testify/assert"
)

func TestWithTempFile(t *testing.T) {

	res := F.Pipe2(
		[]byte("Carsten"),
		onWriteAll[*os.File],
		WithTempFile,
	)

	result, err := res()
	assert.NoError(t, err)
	assert.Equal(t, []byte("Carsten"), result)
}

func TestWithTempFileOnClosedFile(t *testing.T) {

	res := WithTempFile(func(f *os.File) IOResult[[]byte] {
		return F.Pipe2(
			f,
			onWriteAll[*os.File]([]byte("Carsten")),
			ioresult.ChainFirst(F.Constant1[[]byte](Close(f))),
		)
	})

	result, err := res()
	assert.NoError(t, err)
	assert.Equal(t, []byte("Carsten"), result)
}
