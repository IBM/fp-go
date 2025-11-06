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

	"github.com/stretchr/testify/assert"
)

func TestOnError(t *testing.T) {
	onError := OnError("Failed to process [%s]", "filename")

	err := fmt.Errorf("Cause")

	err1 := onError(err)

	assert.Equal(t, "Failed to process [filename], Caused By: Cause", err1.Error())
}

func TestOnNone(t *testing.T) {
	t.Run("creates error without args", func(t *testing.T) {
		getError := OnNone("value not found")
		err := getError()

		assert.NotNil(t, err)
		assert.Equal(t, "value not found", err.Error())
	})

	t.Run("creates error with format args", func(t *testing.T) {
		getError := OnNone("failed to load %s from %s", "config", "file.json")
		err := getError()

		assert.NotNil(t, err)
		assert.Equal(t, "failed to load config from file.json", err.Error())
	})

	t.Run("can be called multiple times", func(t *testing.T) {
		getError := OnNone("repeated error")
		err1 := getError()
		err2 := getError()

		assert.Equal(t, err1.Error(), err2.Error())
	})
}

func TestOnSome(t *testing.T) {
	t.Run("creates error with value only", func(t *testing.T) {
		makeError := OnSome[int]("invalid value: %d")
		err := makeError(42)

		assert.NotNil(t, err)
		assert.Equal(t, "invalid value: 42", err.Error())
	})

	t.Run("creates error with value and additional args", func(t *testing.T) {
		makeError := OnSome[string]("failed to process %s in file %s", "data.txt")
		err := makeError("record123")

		assert.NotNil(t, err)
		assert.Equal(t, "failed to process record123 in file data.txt", err.Error())
	})

	t.Run("works with different types", func(t *testing.T) {
		makeIntError := OnSome[int]("number: %d")
		makeStringError := OnSome[string]("text: %s")
		makeBoolError := OnSome[bool]("flag: %t")

		assert.Equal(t, "number: 100", makeIntError(100).Error())
		assert.Equal(t, "text: hello", makeStringError("hello").Error())
		assert.Equal(t, "flag: true", makeBoolError(true).Error())
	})

	t.Run("works with struct types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		makeError := OnSome[User]("invalid user: %+v")
		user := User{Name: "Alice", Age: 30}
		err := makeError(user)

		assert.Contains(t, err.Error(), "invalid user:")
		assert.Contains(t, err.Error(), "Alice")
	})
}

func TestToString(t *testing.T) {
	t.Run("converts simple error to string", func(t *testing.T) {
		err := fmt.Errorf("simple error")
		str := ToString(err)

		assert.Equal(t, "simple error", str)
	})

	t.Run("converts wrapped error to string", func(t *testing.T) {
		rootErr := fmt.Errorf("root cause")
		wrappedErr := fmt.Errorf("wrapped: %w", rootErr)
		str := ToString(wrappedErr)

		assert.Equal(t, "wrapped: root cause", str)
	})

	t.Run("converts custom error to string", func(t *testing.T) {
		type CustomError struct {
			Code int
			Msg  string
		}
		customErr := &CustomError{Code: 404, Msg: "not found"}
		// Need to implement Error() method
		err := fmt.Errorf("custom: code=%d, msg=%s", customErr.Code, customErr.Msg)
		str := ToString(err)

		assert.Contains(t, str, "404")
		assert.Contains(t, str, "not found")
	})
}
