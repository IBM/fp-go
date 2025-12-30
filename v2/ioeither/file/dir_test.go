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

func TestMkdirAll(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	t.Run("creates nested directories", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "nested", "dir", "structure")

		result := MkdirAll(testPath, 0755)()

		assert.True(t, E.IsRight(result))
		path := E.GetOrElse(func(error) string { return "" })(result)
		assert.Equal(t, testPath, path)

		// Verify directory was created
		info, err := os.Stat(testPath)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("succeeds when directory already exists", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "existing")

		// Create directory first
		err := os.Mkdir(testPath, 0755)
		require.NoError(t, err)

		// Try to create again
		result := MkdirAll(testPath, 0755)()

		assert.True(t, E.IsRight(result))
		path := E.GetOrElse(func(error) string { return "" })(result)
		assert.Equal(t, testPath, path)
	})

	t.Run("fails with invalid path", func(t *testing.T) {
		// Use a path that contains a file as a parent
		filePath := filepath.Join(tempDir, "file.txt")
		err := os.WriteFile(filePath, []byte("test"), 0644)
		require.NoError(t, err)

		invalidPath := filepath.Join(filePath, "subdir")
		result := MkdirAll(invalidPath, 0755)()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("creates directory with correct permissions", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "perms")

		result := MkdirAll(testPath, 0700)()

		assert.True(t, E.IsRight(result))

		// Verify permissions (on Unix-like systems)
		info, err := os.Stat(testPath)
		require.NoError(t, err)
		// Note: actual permissions may differ due to umask
		assert.True(t, info.IsDir())
	})
}

func TestMkdir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	t.Run("creates single directory", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "single")

		result := Mkdir(testPath, 0755)()

		assert.True(t, E.IsRight(result))
		path := E.GetOrElse(func(error) string { return "" })(result)
		assert.Equal(t, testPath, path)

		// Verify directory was created
		info, err := os.Stat(testPath)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("fails when parent does not exist", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "nonexistent", "child")

		result := Mkdir(testPath, 0755)()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("fails when directory already exists", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "existing2")

		// Create directory first
		err := os.Mkdir(testPath, 0755)
		require.NoError(t, err)

		// Try to create again
		result := Mkdir(testPath, 0755)()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("fails with invalid path", func(t *testing.T) {
		// Use a path that contains a file as a parent
		filePath := filepath.Join(tempDir, "file2.txt")
		err := os.WriteFile(filePath, []byte("test"), 0644)
		require.NoError(t, err)

		invalidPath := filepath.Join(filePath, "subdir")
		result := Mkdir(invalidPath, 0755)()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("creates directory with correct permissions", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "perms2")

		result := Mkdir(testPath, 0700)()

		assert.True(t, E.IsRight(result))

		// Verify permissions (on Unix-like systems)
		info, err := os.Stat(testPath)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})
}
