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

package errors

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type MyError struct{}

func (m *MyError) Error() string {
	return "boom"
}

func TestAs(t *testing.T) {
	root := &MyError{}
	err := fmt.Errorf("This is my custom error, %w", root)

	errO := F.Pipe1(
		err,
		As[*MyError](),
	)

	assert.Equal(t, O.Of(root), errO)
}

func TestNotAs(t *testing.T) {
	err := fmt.Errorf("This is my custom error")

	errO := F.Pipe1(
		err,
		As[*MyError](),
	)

	assert.Equal(t, O.None[*MyError](), errO)
}
