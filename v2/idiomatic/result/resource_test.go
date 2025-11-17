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

	"github.com/stretchr/testify/assert"
)

func TestWithResource(t *testing.T) {

	onCreate := func() (*os.File, error) {
		return os.CreateTemp("", "*")
	}
	onDelete := func(f *os.File) (any, error) {
		return Chain(func(name string) (any, error) {
			return any(name), os.Remove(name)
		})(f.Name(), f.Close())
	}

	onHandler := func(f *os.File) (string, error) {
		return Of(f.Name())
	}

	tempFile := WithResource[*os.File, string](onCreate, onDelete)

	res, err := tempFile(onHandler)

	assert.True(t, IsRight(res, err))
}
