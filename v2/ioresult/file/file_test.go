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
	"path/filepath"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	t.Run("successful open", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-open-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpPath)

		err = os.WriteFile(tmpPath, []byte("test content"), 0o644)
		require.NoError(t, err)

		result := Open(tmpPath)()
		file, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.NotNil(t, file)
		file.Close()
	})

	t.Run("open non-existent file", func(t *testing.T) {
		result := Open("/path/that/does/not/exist.txt")()
		_, err := E.UnwrapError(result)

		assert.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "new-file.txt")

		result := Create(testPath)()
		file, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.NotNil(t, file)

		_, statErr := os.Stat(testPath)
		assert.NoError(t, statErr)

		file.Close()
	})

	t.Run("create in non-existent directory", func(t *testing.T) {
		result := Create("/non/existent/directory/file.txt")()
		_, err := E.UnwrapError(result)

		assert.Error(t, err)
	})
}

func TestReadFile(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-read-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		expectedContent := []byte("Hello, World!")
		_, err = tmpFile.Write(expectedContent)
		require.NoError(t, err)
		tmpFile.Close()

		result := ReadFile(tmpPath)()
		content, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, content)
	})

	t.Run("read non-existent file", func(t *testing.T) {
		result := ReadFile("/non/existent/file.txt")()
		_, err := E.UnwrapError(result)

		assert.Error(t, err)
	})

	t.Run("read empty file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-empty-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpPath)

		result := ReadFile(tmpPath)()
		content, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.Empty(t, content)
	})
}

func TestWriteFile(t *testing.T) {
	t.Run("successful write", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "write-test.txt")
		testData := []byte("test data")

		result := WriteFile(testPath, 0o644)(testData)()
		returnedData, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.Equal(t, testData, returnedData)

		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("write to invalid path", func(t *testing.T) {
		testData := []byte("test data")
		result := WriteFile("/non/existent/dir/file.txt", 0o644)(testData)()
		_, err := E.UnwrapError(result)

		assert.Error(t, err)
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-overwrite-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpPath)

		err = os.WriteFile(tmpPath, []byte("initial"), 0o644)
		require.NoError(t, err)

		newData := []byte("overwritten")
		result := WriteFile(tmpPath, 0o644)(newData)()
		returnedData, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.Equal(t, newData, returnedData)

		content, err := os.ReadFile(tmpPath)
		require.NoError(t, err)
		assert.Equal(t, newData, content)
	})
}

func TestRemove(t *testing.T) {
	t.Run("successful remove", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-remove-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		tmpFile.Close()

		result := Remove(tmpPath)()
		name, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.Equal(t, tmpPath, name)

		_, statErr := os.Stat(tmpPath)
		assert.True(t, os.IsNotExist(statErr))
	})

	t.Run("remove non-existent file", func(t *testing.T) {
		result := Remove("/non/existent/file.txt")()
		_, err := E.UnwrapError(result)

		assert.Error(t, err)
	})
}

func TestClose(t *testing.T) {
	t.Run("successful close", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-close-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		result := Close(tmpFile)()
		_, err = E.UnwrapError(result)

		assert.NoError(t, err)

		_, writeErr := tmpFile.WriteString("test")
		assert.Error(t, writeErr)
	})

	t.Run("close already closed file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-close-twice-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		tmpFile.Close()

		result := Close(tmpFile)()
		_, err = E.UnwrapError(result)

		assert.Error(t, err)
	})
}
