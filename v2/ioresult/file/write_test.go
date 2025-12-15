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
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
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
		)()

		returnedData, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.Equal(t, testData, returnedData)

		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("write all ensures file is closed", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "writeall-close-test.txt")
		testData := []byte("test data")

		var capturedFile *os.File
		acquire := func() E.Either[error, *os.File] {
			f, err := os.Create(testPath)
			capturedFile = f
			if err != nil {
				return E.Left[*os.File](err)
			}
			return E.Right[error](f)
		}

		result := F.Pipe1(
			acquire,
			WriteAll[*os.File](testData),
		)()

		_, err := E.UnwrapError(result)
		assert.NoError(t, err)

		_, writeErr := capturedFile.WriteString("more")
		assert.Error(t, writeErr)
	})

	t.Run("write all with acquire failure", func(t *testing.T) {
		testData := []byte("test data")

		acquire := Create("/non/existent/dir/file.txt")
		result := F.Pipe1(
			acquire,
			WriteAll[*os.File](testData),
		)()

		_, err := E.UnwrapError(result)
		assert.Error(t, err)
	})

	t.Run("write all with empty data", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "empty-writeall.txt")

		acquire := Create(testPath)
		result := F.Pipe1(
			acquire,
			WriteAll[*os.File]([]byte{}),
		)()

		returnedData, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.Empty(t, returnedData)

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
			return ioeither.TryCatchError(func() (int, error) {
				return f.Write(testData)
			})
		}

		result := Write[int](acquire)(useFile)()
		bytesWritten, err := E.UnwrapError(result)

		assert.NoError(t, err)
		assert.Equal(t, len(testData), bytesWritten)

		content, err := os.ReadFile(testPath)
		require.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("write ensures cleanup on success", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "write-cleanup-test.txt")

		var capturedFile *os.File
		acquire := func() E.Either[error, *os.File] {
			f, err := os.Create(testPath)
			capturedFile = f
			if err != nil {
				return E.Left[*os.File](err)
			}
			return E.Right[error](f)
		}

		useFile := func(f *os.File) IOResult[string] {
			return ioeither.TryCatchError(func() (string, error) {
				_, err := f.WriteString("data")
				return "success", err
			})
		}

		result := Write[string](acquire)(useFile)()
		_, err := E.UnwrapError(result)

		assert.NoError(t, err)

		_, writeErr := capturedFile.WriteString("more")
		assert.Error(t, writeErr)
	})

	t.Run("write ensures cleanup on failure", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "write-fail-test.txt")

		var capturedFile *os.File
		acquire := func() E.Either[error, *os.File] {
			f, err := os.Create(testPath)
			capturedFile = f
			if err != nil {
				return E.Left[*os.File](err)
			}
			return E.Right[error](f)
		}

		useFile := func(f *os.File) IOResult[string] {
			return ioeither.Left[string](assert.AnError)
		}

		result := Write[string](acquire)(useFile)()
		_, err := E.UnwrapError(result)

		assert.Error(t, err)

		_, writeErr := capturedFile.WriteString("more")
		assert.Error(t, writeErr)
	})

	t.Run("write with acquire failure", func(t *testing.T) {
		useFile := func(f *os.File) IOResult[string] {
			return ioeither.Of[error]("should not run")
		}

		result := Write[string](Create("/non/existent/dir/file.txt"))(useFile)()
		_, err := E.UnwrapError(result)

		assert.Error(t, err)
	})
}
