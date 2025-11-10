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

package result

import (
	"os"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestWithResource(t *testing.T) {
	onCreate := func() Result[*os.File] {
		return TryCatchError(os.CreateTemp("", "*"))
	}
	onDelete := F.Flow2(
		func(f *os.File) Result[string] {
			return TryCatchError(f.Name(), f.Close())
		},
		Chain(func(name string) Result[any] {
			return TryCatchError(any(name), os.Remove(name))
		}),
	)

	onHandler := func(f *os.File) Result[string] {
		return Of(f.Name())
	}

	tempFile := WithResource[*os.File, string](onCreate, onDelete)

	resE := tempFile(onHandler)

	assert.True(t, IsRight(resE))
}
