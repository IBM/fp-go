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

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
)

func TestWithTempFile(t *testing.T) {

	res := WithTempFile(onWriteAll[*os.File]([]byte("Carsten")))

	assert.Equal(t, E.Of[error]([]byte("Carsten")), res())
}

func TestWithTempFileOnClosedFile(t *testing.T) {

	res := WithTempFile(func(f *os.File) IOEither[error, []byte] {
		return F.Pipe2(
			f,
			onWriteAll[*os.File]([]byte("Carsten")),
			ioeither.ChainFirst(F.Constant1[[]byte](Close(f))),
		)
	})

	assert.Equal(t, E.Of[error]([]byte("Carsten")), res())
}
