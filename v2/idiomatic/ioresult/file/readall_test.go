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
	"github.com/stretchr/testify/require"
)

func TestReadAll(t *testing.T) {
	t.Run("successful read all", func(t *testing.T) {

		tmpFile, err := os.CreateTemp("", "test-readall-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		expectedContent := []byte("Hello, ReadAll!")
		_, err = tmpFile.Write(expectedContent)
		require.NoError(t, err)
		tmpFile.Close()

		result := ReadAll(Open(tmpPath))
		content, err := result()

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, content)
	})

	t.Run("read all ensures file is closed", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-readall-close-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		testContent := []byte("test data for close")
		_, err = tmpFile.Write(testContent)
		require.NoError(t, err)
		tmpFile.Close()

		var capturedFile *os.File
		acquire := func() (*os.File, error) {
			f, err := os.Open(tmpPath)
			capturedFile = f
			return f, err
		}

		result := ReadAll(acquire)
		content, err := result()

		assert.NoError(t, err)
		assert.Equal(t, testContent, content)

		// Verify file is closed by trying to read
		buf := make([]byte, 10)
		_, readErr := capturedFile.Read(buf)
		assert.Error(t, readErr)
	})

	t.Run("read all with open failure", func(t *testing.T) {
		result := ReadAll(Open("/non/existent/file.txt"))
		_, err := result()

		assert.Error(t, err)
	})

	t.Run("read all empty file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-readall-empty-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpPath)

		result := ReadAll(Open(tmpPath))
		content, err := result()

		assert.NoError(t, err)
		assert.Empty(t, content)
	})

	t.Run("read all large file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "large-file.txt")

		// Create a larger file
		largeContent := make([]byte, 10000)
		for i := range largeContent {
			largeContent[i] = byte('A' + (i % 26))
		}

		err := os.WriteFile(testPath, largeContent, 0o644)
		require.NoError(t, err)

		result := ReadAll(Open(testPath))
		content, err := result()

		assert.NoError(t, err)
		assert.Equal(t, largeContent, content)
		assert.Len(t, content, 10000)
	})

	t.Run("read all with binary data", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-readall-binary-*.bin")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		// Write binary data
		binaryContent := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
		_, err = tmpFile.Write(binaryContent)
		require.NoError(t, err)
		tmpFile.Close()

		result := ReadAll(Open(tmpPath))
		content, err := result()

		assert.NoError(t, err)
		assert.Equal(t, binaryContent, content)
	})
}
