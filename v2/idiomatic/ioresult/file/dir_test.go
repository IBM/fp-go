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

	"github.com/stretchr/testify/assert"
)

func TestMkdir(t *testing.T) {
	t.Run("successful mkdir", func(t *testing.T) {
		tmpDir := t.TempDir()
		newDir := filepath.Join(tmpDir, "testdir")

		result := Mkdir(newDir, 0o755)
		path, err := result()

		assert.NoError(t, err)
		assert.Equal(t, newDir, path)

		// Verify directory was created
		info, err := os.Stat(newDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("mkdir with existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		result := Mkdir(tmpDir, 0o755)
		_, err := result()

		assert.Error(t, err)
	})

	t.Run("mkdir with parent directory not existing", func(t *testing.T) {
		result := Mkdir("/non/existent/parent/child", 0o755)
		_, err := result()

		assert.Error(t, err)
	})
}

func TestMkdirAll(t *testing.T) {
	t.Run("successful mkdir all", func(t *testing.T) {
		tmpDir := t.TempDir()
		nestedDir := filepath.Join(tmpDir, "level1", "level2", "level3")

		result := MkdirAll(nestedDir, 0o755)
		path, err := result()

		assert.NoError(t, err)
		assert.Equal(t, nestedDir, path)

		// Verify all directories were created
		info, err := os.Stat(nestedDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())

		// Verify intermediate directories
		level1 := filepath.Join(tmpDir, "level1")
		info1, err := os.Stat(level1)
		assert.NoError(t, err)
		assert.True(t, info1.IsDir())

		level2 := filepath.Join(tmpDir, "level1", "level2")
		info2, err := os.Stat(level2)
		assert.NoError(t, err)
		assert.True(t, info2.IsDir())
	})

	t.Run("mkdirall with existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		result := MkdirAll(tmpDir, 0o755)
		path, err := result()

		// MkdirAll should succeed even if directory exists
		assert.NoError(t, err)
		assert.Equal(t, tmpDir, path)
	})

	t.Run("mkdirall single level", func(t *testing.T) {
		tmpDir := t.TempDir()
		newDir := filepath.Join(tmpDir, "single")

		result := MkdirAll(newDir, 0o755)
		path, err := result()

		assert.NoError(t, err)
		assert.Equal(t, newDir, path)

		info, err := os.Stat(newDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("mkdirall with file in path", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "file.txt")

		// Create a file
		err := os.WriteFile(filePath, []byte("content"), 0o644)
		assert.NoError(t, err)

		// Try to create a directory where file exists
		result := MkdirAll(filepath.Join(filePath, "subdir"), 0o755)
		_, err = result()

		assert.Error(t, err)
	})
}
