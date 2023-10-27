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

package either

import (
	"os"
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestWithResource(t *testing.T) {
	onCreate := func() Either[error, *os.File] {
		return TryCatchError(os.CreateTemp("", "*"))
	}
	onDelete := F.Flow2(
		func(f *os.File) Either[error, string] {
			return TryCatchError(f.Name(), f.Close())
		},
		Chain(func(name string) Either[error, any] {
			return TryCatchError(any(name), os.Remove(name))
		}),
	)

	onHandler := func(f *os.File) Either[error, string] {
		return Of[error](f.Name())
	}

	tempFile := WithResource[error, *os.File, string](onCreate, onDelete)

	resE := tempFile(onHandler)

	assert.True(t, IsRight(resE))
}
