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

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteAll(t *testing.T) {
	t.Run("successful write all", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "writeall-test.txt")
		testData := []byte("Hello, WriteAll!")

		acquire := Create(testPath)
		result := F.Pipe1(
			acquire,
			WriteAll[*os.File](testData),
		)

		returnedData, err := result()

		assert.NoError(t, err)
		assert.Equal(t, testData, returnedData)

		// Verify file content
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("write all ensures file is closed", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "writeall-close-test.txt")
		testData := []byte("test data")

		var capturedFile *os.File
		acquire := func() (*os.File, error) {
			f, err := os.Create(testPath)
			capturedFile = f
			return f, err
		}

		result := F.Pipe1(
			ioresult.FromResult(acquire()),
			WriteAll[*os.File](testData),
		)

		_, err := result()
		assert.NoError(t, err)

		// Verify file is closed by trying to write to it
		_, writeErr := capturedFile.WriteString("more")
		assert.Error(t, writeErr)
	})

	t.Run("write all with acquire failure", func(t *testing.T) {
		testData := []byte("test data")

		acquire := Create("/non/existent/dir/file.txt")
		result := F.Pipe1(
			acquire,
			WriteAll[*os.File](testData),
		)

		_, err := result()
		assert.Error(t, err)
	})

	t.Run("write all with empty data", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "empty-writeall.txt")

		acquire := Create(testPath)
		result := F.Pipe1(
			acquire,
			WriteAll[*os.File]([]byte{}),
		)

		returnedData, err := result()

		assert.NoError(t, err)
		assert.Empty(t, returnedData)

		// Verify file exists and is empty
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Empty(t, content)
	})
}

func TestWrite(t *testing.T) {
	t.Run("successful write with resource management", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "write-test.txt")
		testData := []byte("test content")

		acquire := Create(testPath)
		useFile := func(f *os.File) IOResult[int] {
			return func() (int, error) {
				return f.Write(testData)
			}
		}

		result := Write[int](acquire)(useFile)
		bytesWritten, err := result()

		assert.NoError(t, err)
		assert.Equal(t, len(testData), bytesWritten)

		// Verify file content
		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("write ensures cleanup on success", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "write-cleanup-test.txt")

		var capturedFile *os.File
		acquire := func() (*os.File, error) {
			f, err := os.Create(testPath)
			capturedFile = f
			return f, err
		}

		useFile := func(f *os.File) IOResult[string] {
			return func() (string, error) {
				_, err := f.WriteString("data")
				return "success", err
			}
		}

		result := Write[string](ioresult.FromResult(acquire()))(useFile)
		_, err := result()

		assert.NoError(t, err)

		// Verify file is closed
		_, writeErr := capturedFile.WriteString("more")
		assert.Error(t, writeErr)
	})

	t.Run("write ensures cleanup on failure", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "write-fail-test.txt")

		var capturedFile *os.File
		acquire := func() (*os.File, error) {
			f, err := os.Create(testPath)
			capturedFile = f
			return f, err
		}

		useFile := func(f *os.File) IOResult[string] {
			return ioresult.Left[string](assert.AnError)
		}

		result := Write[string](ioresult.FromResult(acquire()))(useFile)
		_, err := result()

		assert.Error(t, err)

		// Verify file is still closed even on error
		_, writeErr := capturedFile.WriteString("more")
		assert.Error(t, writeErr)
	})

	t.Run("write with acquire failure", func(t *testing.T) {
		useFile := func(f *os.File) IOResult[string] {
			return ioresult.Of("should not run")
		}

		result := Write[string](Create("/non/existent/dir/file.txt"))(useFile)
		_, err := result()

		assert.Error(t, err)
	})
}
