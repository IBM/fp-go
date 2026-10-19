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
	"errors"
	"fmt"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// codeError implements error with a value receiver.
type codeError struct{ code int }

func (e codeError) Error() string { return fmt.Sprintf("code %d", e.code) }

func TestErrorAs_Success(t *testing.T) {
	t.Run("selects an error of the type", func(t *testing.T) {
		pathErr := &fs.PathError{Op: "open", Path: "a", Err: fs.ErrNotExist}
		assert.Equal(t, Of(pathErr), ErrorAs[*fs.PathError](pathErr))
	})

	t.Run("selects a wrapped error", func(t *testing.T) {
		pathErr := &fs.PathError{Op: "open", Path: "a", Err: fs.ErrNotExist}
		assert.Equal(t, Of(pathErr), ErrorAs[*fs.PathError](fmt.Errorf("config: %w", pathErr)))
	})

	t.Run("selects an error from a join", func(t *testing.T) {
		assert.Equal(t, Of(codeError{code: 1}), ErrorAs[codeError](errors.Join(errSentinel, codeError{code: 1})))
	})
}

func TestErrorAs_Failure(t *testing.T) {
	t.Run("keeps an error of another type", func(t *testing.T) {
		assert.Equal(t, Left[codeError](errSentinel), ErrorAs[codeError](errSentinel))
	})

	t.Run("pointer type does not match a value receiver error", func(t *testing.T) {
		err := codeError{code: 1}
		assert.Equal(t, Left[*codeError](error(err)), ErrorAs[*codeError](err))
	})

	t.Run("value type does not match a pointer to a value receiver error", func(t *testing.T) {
		err := &codeError{code: 1}
		assert.Equal(t, Left[codeError](error(err)), ErrorAs[codeError](err))
	})
}

func TestErrorIs_Success(t *testing.T) {
	err := fmt.Errorf("read: %w", errSentinel)
	assert.Equal(t, Of(err), ErrorIs(errSentinel)(err))
	assert.Equal(t, Of(errSentinel), ErrorIs(errSentinel)(errSentinel))
}

func TestErrorIs_Failure(t *testing.T) {
	other := errors.New("other")
	assert.Equal(t, Left[error](other), ErrorIs(errSentinel)(other))
}
