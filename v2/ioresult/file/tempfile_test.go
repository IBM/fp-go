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
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
)

func TestCreateTemp(t *testing.T) {
	t.Run("successful create temp", func(t *testing.T) {
		result := CreateTemp("", "test-*.txt")()
		file, err := E.UnwrapError(result)

		assert.NoError(t, err)

		assert.NotNil(t, file)

		tmpPath := file.Name()
		file.Close()
		defer os.Remove(tmpPath)

		_, statErr := os.Stat(tmpPath)
		assert.NoError(t, statErr)
	})

	t.Run("create temp with pattern", func(t *testing.T) {
		result := CreateTemp("", "prefix-*-suffix.txt")()
		file, err := E.UnwrapError(result)

		assert.NoError(t, err)

		assert.NotNil(t, file)

		tmpPath := file.Name()
		assert.Contains(t, tmpPath, "prefix-")
		assert.Contains(t, tmpPath, "-suffix.txt")

		file.Close()
		os.Remove(tmpPath)
	})

	t.Run("create temp in specific directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		result := CreateTemp(tmpDir, "test-*.txt")()
		file, err := E.UnwrapError(result)

		assert.NoError(t, err)

		assert.NotNil(t, file)

		tmpPath := file.Name()
		assert.Contains(t, tmpPath, tmpDir)

		file.Close()
		os.Remove(tmpPath)
	})
}

func TestWithTempFile(t *testing.T) {
	t.Run("successful temp file usage", func(t *testing.T) {
		testData := []byte("temp file content")

		useFile := func(f *os.File) IOResult[[]byte] {
			return ioeither.TryCatchError(func() ([]byte, error) {
				_, err := f.Write(testData)
				return testData, err
			})
		}

		result := WithTempFile(useFile)()
		returnedData, err := E.UnwrapError(result)

		assert.NoError(t, err)

		assert.Equal(t, testData, returnedData)
	})

	t.Run("temp file is removed after use", func(t *testing.T) {
		var tmpPath string

		useFile := func(f *os.File) IOResult[string] {
			return ioeither.TryCatchError(func() (string, error) {
				tmpPath = f.Name()
				_, err := f.WriteString("test")
				return tmpPath, err
			})
		}

		result := WithTempFile(useFile)()
		path, err := E.UnwrapError(result)

		assert.NoError(t, err)

		assert.Equal(t, tmpPath, path)

		_, statErr := os.Stat(tmpPath)
		assert.True(t, os.IsNotExist(statErr))
	})

	t.Run("temp file removed even on failure", func(t *testing.T) {
		var tmpPath string

		useFile := func(f *os.File) IOResult[string] {
			tmpPath = f.Name()
			return ioeither.Left[string](assert.AnError)
		}

		result := WithTempFile(useFile)()
		_, err := E.UnwrapError(result)

		assert.Error(t, err)

		_, statErr := os.Stat(tmpPath)
		assert.True(t, os.IsNotExist(statErr))
	})

	t.Run("temp file can be written and read", func(t *testing.T) {
		testContent := []byte("Hello, Temp File!")

		useFile := func(f *os.File) IOResult[[]byte] {
			return ioeither.TryCatchError(func() ([]byte, error) {
				_, err := f.Write(testContent)
				if err != nil {
					return nil, err
				}

				_, err = f.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				buf := make([]byte, len(testContent))
				_, err = f.Read(buf)
				return buf, err
			})
		}

		result := WithTempFile(useFile)()
		content, err := E.UnwrapError(result)

		assert.NoError(t, err)

		assert.Equal(t, testContent, content)
	})

	t.Run("multiple temp files", func(t *testing.T) {
		var paths []string

		for i := 0; i < 3; i++ {
			useFile := func(f *os.File) IOResult[string] {
				return ioeither.TryCatchError(func() (string, error) {
					return f.Name(), nil
				})
			}

			result := WithTempFile(useFile)()
			path, err := E.UnwrapError(result)

			assert.NoError(t, err)

			paths = append(paths, path)
		}

		// All paths should be different
		assert.NotEqual(t, paths[0], paths[1])
		assert.NotEqual(t, paths[1], paths[2])
		assert.NotEqual(t, paths[0], paths[2])

		// All files should be removed
		for _, path := range paths {
			_, statErr := os.Stat(path)
			assert.True(t, os.IsNotExist(statErr))
		}
	})

	t.Run("file is closed before removal", func(t *testing.T) {
		var capturedFile *os.File

		useFile := func(f *os.File) IOResult[string] {
			capturedFile = f
			return ioeither.TryCatchError(func() (string, error) {
				_, err := f.WriteString("test")
				return f.Name(), err
			})
		}

		result := WithTempFile(useFile)()
		_, err := E.UnwrapError(result)

		assert.NoError(t, err)

		// File should be closed
		_, writeErr := capturedFile.WriteString("more")
		assert.Error(t, writeErr)
	})

	t.Run("integration with write operations", func(t *testing.T) {
		testData := []byte("Integration test data")

		// Helper to write to file
		onWriteAll := func(data []byte) func(*os.File) IOResult[[]byte] {
			return func(w *os.File) IOResult[[]byte] {
				return ioeither.TryCatchError(func() ([]byte, error) {
					_, err := w.Write(data)
					return data, err
				})
			}
		}

		result := WithTempFile(onWriteAll(testData))()
		returnedData, err := E.UnwrapError(result)

		assert.NoError(t, err)

		assert.Equal(t, testData, returnedData)
	})

	t.Run("file closed even if use operation closes it first", func(t *testing.T) {
		useFile := func(f *os.File) IOResult[[]byte] {
			return F.Pipe2(
				f,
				func(file *os.File) IOResult[[]byte] {
					return ioeither.TryCatchError(func() ([]byte, error) {
						data := []byte("test")
						_, err := file.Write(data)
						return data, err
					})
				},
				ioeither.ChainFirst(F.Constant1[[]byte](Close(f))),
			)
		}

		result := WithTempFile(useFile)()
		_, err := E.UnwrapError(result)

		assert.NoError(t, err)
	})
}
