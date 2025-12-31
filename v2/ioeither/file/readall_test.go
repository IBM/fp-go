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

func TestReadAll(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("reads entire file", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "readall.txt")
		testData := []byte("Hello, ReadAll!")

		// Create test file
		err := os.WriteFile(testPath, testData, 0644)
		require.NoError(t, err)

		// Read file
		result := ReadAll(Open(testPath))()

		assert.True(t, E.IsRight(result))
		data := E.GetOrElse(func(error) []byte { return nil })(result)
		assert.Equal(t, testData, data)
	})

	t.Run("reads empty file", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "empty.txt")

		// Create empty file
		err := os.WriteFile(testPath, []byte{}, 0644)
		require.NoError(t, err)

		// Read file
		result := ReadAll(Open(testPath))()

		assert.True(t, E.IsRight(result))
		data := E.GetOrElse(func(error) []byte { return nil })(result)
		assert.Equal(t, 0, len(data))
	})

	t.Run("reads large file", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "large.txt")

		// Create large file (1MB)
		largeData := make([]byte, 1024*1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}
		err := os.WriteFile(testPath, largeData, 0644)
		require.NoError(t, err)

		// Read file
		result := ReadAll(Open(testPath))()

		assert.True(t, E.IsRight(result))
		data := E.GetOrElse(func(error) []byte { return nil })(result)
		assert.Equal(t, len(largeData), len(data))
		assert.Equal(t, largeData, data)
	})

	t.Run("fails when file does not exist", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "nonexistent.txt")

		// Try to read non-existent file
		result := ReadAll(Open(testPath))()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("fails when trying to read directory", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "dir")
		err := os.Mkdir(testPath, 0755)
		require.NoError(t, err)

		// Try to read directory
		result := ReadAll(Open(testPath))()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("reads file with special characters", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "special.txt")
		testData := []byte("Hello\nWorld\t!\r\n")

		// Create test file
		err := os.WriteFile(testPath, testData, 0644)
		require.NoError(t, err)

		// Read file
		result := ReadAll(Open(testPath))()

		assert.True(t, E.IsRight(result))
		data := E.GetOrElse(func(error) []byte { return nil })(result)
		assert.Equal(t, testData, data)
	})

	t.Run("reads binary file", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "binary.bin")
		testData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}

		// Create binary file
		err := os.WriteFile(testPath, testData, 0644)
		require.NoError(t, err)

		// Read file
		result := ReadAll(Open(testPath))()

		assert.True(t, E.IsRight(result))
		data := E.GetOrElse(func(error) []byte { return nil })(result)
		assert.Equal(t, testData, data)
	})

	t.Run("closes file after reading", func(t *testing.T) {
		testPath := filepath.Join(tempDir, "close_test.txt")
		testData := []byte("test")

		// Create test file
		err := os.WriteFile(testPath, testData, 0644)
		require.NoError(t, err)

		// Read file
		result := ReadAll(Open(testPath))()
		assert.True(t, E.IsRight(result))

		// Verify we can delete the file (it's closed)
		err = os.Remove(testPath)
		assert.NoError(t, err)
	})
}
