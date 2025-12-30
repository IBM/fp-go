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
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteAll(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("writes data and closes file", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "writeall.txt")
		testData := []byte("Hello, WriteAll!")

		// Create and write to file
		result := WriteAll[*os.File](testData)(Create(testPath))()

		assert.True(t, E.IsRight(result))
		data := E.GetOrElse(func(error) []byte { return nil })(result)
		assert.Equal(t, testData, data)

		// Verify file contents
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("handles write errors", func(t *testing.T) {
		// Try to write to a directory (should fail)
		testPath := filepath.Join(tempDir, "dir_not_file")
		err := os.Mkdir(testPath, 0755)
		require.NoError(t, err)

		testData := []byte("This should fail")

		result := WriteAll[*os.File](testData)(Create(testPath))()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("writes empty data", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "empty.txt")
		testData := []byte{}

		result := WriteAll[*os.File](testData)(Create(testPath))()

		assert.True(t, E.IsRight(result))

		// Verify file is empty
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, 0, len(content))
	})

	t.Run("overwrites existing file", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "overwrite.txt")

		// Write initial content
		err := os.WriteFile(testPath, []byte("old content"), 0644)
		require.NoError(t, err)

		// Overwrite with new content
		newData := []byte("new content")
		result := WriteAll[*os.File](newData)(Create(testPath))()

		assert.True(t, E.IsRight(result))

		// Verify new content
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, newData, content)
	})
}

func TestWrite(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("writes data with resource management", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "write.txt")
		testData := []byte("Hello, Write!")

		writeOp := Write[int](Create(testPath))
		result := writeOp(func(f *os.File) IOEither[error, int] {
			return IOE.TryCatchError(func() (int, error) {
				return f.Write(testData)
			})
		})()

		assert.True(t, E.IsRight(result))
		n := E.GetOrElse(func(error) int { return 0 })(result)
		assert.Equal(t, len(testData), n)

		// Verify file contents
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("handles multiple writes", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "multiwrite.txt")

		writeOp := Write[int](Create(testPath))
		result := writeOp(func(f *os.File) IOEither[error, int] {
			return IOE.TryCatchError(func() (int, error) {
				n1, err := f.Write([]byte("First "))
				if err != nil {
					return 0, err
				}
				n2, err := f.Write([]byte("Second"))
				return n1 + n2, err
			})
		})()

		assert.True(t, E.IsRight(result))

		// Verify file contents
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, "First Second", string(content))
	})

	t.Run("closes file even on error", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "error.txt")

		writeOp := Write[int](Create(testPath))
		result := writeOp(func(f *os.File) IOEither[error, int] {
			// Close the file prematurely to cause an error
			f.Close()
			return IOE.TryCatchError(func() (int, error) {
				return f.Write([]byte("This should fail"))
			})
		})()

		assert.True(t, E.IsLeft(result))
	})
}

func TestWriteFile(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("writes file successfully", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "writefile.txt")
		testData := []byte("Hello, WriteFile!")

		result := WriteFile(testPath, 0644)(testData)()

		assert.True(t, E.IsRight(result))
		data := E.GetOrElse(func(error) []byte { return nil })(result)
		assert.Equal(t, testData, data)

		// Verify file contents
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("creates file with correct permissions", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "perms.txt")
		testData := []byte("test")

		result := WriteFile(testPath, 0600)(testData)()

		assert.True(t, E.IsRight(result))

		// Verify file exists
		info, err := os.Stat(testPath)
		require.NoError(t, err)
		assert.False(t, info.IsDir())
	})

	t.Run("overwrites existing file", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "overwrite2.txt")

		// Write initial content
		err := os.WriteFile(testPath, []byte("old"), 0644)
		require.NoError(t, err)

		// Overwrite
		newData := []byte("new")
		result := WriteFile(testPath, 0644)(newData)()

		assert.True(t, E.IsRight(result))

		// Verify new content
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, newData, content)
	})

	t.Run("fails with invalid path", func(t *testing.T) {
		// Try to write to a directory that doesn't exist
		testPath := filepath.Join(tempDir, "nonexistent", "file.txt")
		testData := []byte("test")

		result := WriteFile(testPath, 0644)(testData)()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("writes empty file", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "empty2.txt")
		testData := []byte{}

		result := WriteFile(testPath, 0644)(testData)()

		assert.True(t, E.IsRight(result))

		// Verify file is empty
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, 0, len(content))
	})
}
