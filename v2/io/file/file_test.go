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
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	IO "github.com/IBM/fp-go/v2/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCloser is a mock implementation of io.Closer for testing
type mockCloser struct {
	closed    bool
	closeErr  error
	closeFunc func() error
}

func (m *mockCloser) Close() error {
	m.closed = true
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return m.closeErr
}

// TestClose_WithMockCloser tests the Close function with a mock closer
func TestClose_WithMockCloser(t *testing.T) {
	t.Run("closes resource successfully", func(t *testing.T) {
		mock := &mockCloser{}
		closeIO := Close(mock)

		result := closeIO()

		assert.True(t, mock.closed, "resource should be closed")
		assert.Equal(t, mock, result, "should return the same resource")
	})

	t.Run("ignores close error", func(t *testing.T) {
		mock := &mockCloser{
			closeErr: fmt.Errorf("close error"),
		}
		closeIO := Close(mock)

		// Should not panic even with error
		result := closeIO()

		assert.True(t, mock.closed, "resource should be closed despite error")
		assert.Equal(t, mock, result, "should return the same resource")
	})

	t.Run("can be called multiple times", func(t *testing.T) {
		mock := &mockCloser{}
		closeIO := Close(mock)

		result1 := closeIO()
		result2 := closeIO()

		assert.True(t, mock.closed, "resource should be closed")
		assert.Equal(t, result1, result2, "should return same resource each time")
	})
}

// TestClose_WithBytesBuffer tests Close with bytes.Buffer (implements io.Closer)
func TestClose_WithBytesBuffer(t *testing.T) {
	t.Run("closes bytes.Buffer", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte("test data"))
		closeIO := Close(io.NopCloser(buf))

		result := closeIO()

		assert.NotNil(t, result, "should return the closer")
	})
}

// TestClose_WithFile tests Close with actual file
func TestClose_WithFile(t *testing.T) {
	t.Run("closes real file", func(t *testing.T) {
		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test-close-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		// Write some data
		_, err = tmpFile.WriteString("test data")
		require.NoError(t, err)

		// Close using our function
		closeIO := Close(tmpFile)
		result := closeIO()

		assert.Equal(t, tmpFile, result, "should return the same file")

		// Verify file is closed by trying to write (should fail)
		_, err = tmpFile.WriteString("more data")
		assert.Error(t, err, "writing to closed file should fail")
	})
}

// TestClose_Composition tests Close in IO composition
func TestClose_Composition(t *testing.T) {
	t.Run("composes with other IO operations", func(t *testing.T) {
		mock := &mockCloser{}

		// Create a pipeline that uses the resource and then closes it
		step1 := IO.Of(mock)
		step2 := IO.Map(func(m *mockCloser) *mockCloser {
			// Simulate using the resource
			return m
		})(step1)
		pipeline := IO.Chain(Close[*mockCloser])(step2)

		result := pipeline()

		assert.True(t, mock.closed, "resource should be closed in pipeline")
		assert.Equal(t, mock, result, "should return the resource")
	})

	t.Run("works with ChainFirst", func(t *testing.T) {
		mock := &mockCloser{}
		data := "test data"

		// Process data and close resource as side effect
		pipeline := IO.ChainFirst(func(string) IO.IO[*mockCloser] {
			return Close(mock)
		})(IO.Of(data))

		result := pipeline()

		assert.True(t, mock.closed, "resource should be closed")
		assert.Equal(t, data, result, "should return original data")
	})
}

// TestRemove_BasicOperation tests basic Remove functionality
func TestRemove_BasicOperation(t *testing.T) {
	t.Run("removes existing file", func(t *testing.T) {
		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test-remove-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		tmpFile.Close()

		// Verify file exists
		_, err = os.Stat(tmpPath)
		require.NoError(t, err, "file should exist before removal")

		// Remove using our function
		removeIO := Remove(tmpPath)
		result := removeIO()

		assert.Equal(t, tmpPath, result, "should return the file path")

		// Verify file is removed
		_, err = os.Stat(tmpPath)
		assert.True(t, os.IsNotExist(err), "file should not exist after removal")
	})

	t.Run("ignores error for non-existent file", func(t *testing.T) {
		nonExistentPath := filepath.Join(os.TempDir(), "non-existent-file-12345.txt")

		// Should not panic even if file doesn't exist
		removeIO := Remove(nonExistentPath)
		result := removeIO()

		assert.Equal(t, nonExistentPath, result, "should return the path")
	})

	t.Run("removes empty directory", func(t *testing.T) {
		// Create a temporary directory
		tmpDir, err := os.MkdirTemp("", "test-remove-dir-*")
		require.NoError(t, err)

		// Verify directory exists
		_, err = os.Stat(tmpDir)
		require.NoError(t, err, "directory should exist before removal")

		// Remove using our function
		removeIO := Remove(tmpDir)
		result := removeIO()

		assert.Equal(t, tmpDir, result, "should return the directory path")

		// Verify directory is removed
		_, err = os.Stat(tmpDir)
		assert.True(t, os.IsNotExist(err), "directory should not exist after removal")
	})

	t.Run("ignores error for non-empty directory", func(t *testing.T) {
		// Create a temporary directory with a file
		tmpDir, err := os.MkdirTemp("", "test-remove-nonempty-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir) // Cleanup

		tmpFile := filepath.Join(tmpDir, "file.txt")
		err = os.WriteFile(tmpFile, []byte("data"), 0644)
		require.NoError(t, err)

		// Should not panic even if directory is not empty
		removeIO := Remove(tmpDir)
		result := removeIO()

		assert.Equal(t, tmpDir, result, "should return the path")

		// Directory should still exist (os.Remove doesn't remove non-empty dirs)
		_, err = os.Stat(tmpDir)
		assert.NoError(t, err, "non-empty directory should still exist")
	})
}

// TestRemove_Composition tests Remove in IO composition
func TestRemove_Composition(t *testing.T) {
	t.Run("composes with other IO operations", func(t *testing.T) {
		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test-compose-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		tmpFile.Close()

		// Create a pipeline that processes and removes the file
		step1 := IO.Of(tmpPath)
		step2 := IO.Map(func(path string) string {
			// Simulate processing
			return path
		})(step1)
		pipeline := IO.Chain(Remove)(step2)

		result := pipeline()

		assert.Equal(t, tmpPath, result, "should return the path")

		// Verify file is removed
		_, err = os.Stat(tmpPath)
		assert.True(t, os.IsNotExist(err), "file should be removed")
	})

	t.Run("removes multiple files in sequence", func(t *testing.T) {
		// Create temporary files
		tmpFile1, err := os.CreateTemp("", "test-multi-1-*.txt")
		require.NoError(t, err)
		tmpPath1 := tmpFile1.Name()
		tmpFile1.Close()

		tmpFile2, err := os.CreateTemp("", "test-multi-2-*.txt")
		require.NoError(t, err)
		tmpPath2 := tmpFile2.Name()
		tmpFile2.Close()

		// Remove both files in sequence
		pipeline := IO.ChainTo[string](Remove(tmpPath2))(Remove(tmpPath1))

		result := pipeline()

		assert.Equal(t, tmpPath2, result, "should return last path")

		// Verify both files are removed
		_, err = os.Stat(tmpPath1)
		assert.True(t, os.IsNotExist(err), "first file should be removed")

		_, err = os.Stat(tmpPath2)
		assert.True(t, os.IsNotExist(err), "second file should be removed")
	})
}

// TestRemove_CanBeCalledMultipleTimes tests idempotency
func TestRemove_CanBeCalledMultipleTimes(t *testing.T) {
	t.Run("calling remove multiple times is safe", func(t *testing.T) {
		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test-idempotent-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()
		tmpFile.Close()

		removeIO := Remove(tmpPath)

		// First call removes the file
		result1 := removeIO()
		assert.Equal(t, tmpPath, result1)

		// Second call should not panic (file already removed)
		result2 := removeIO()
		assert.Equal(t, tmpPath, result2)

		// Verify file is removed
		_, err = os.Stat(tmpPath)
		assert.True(t, os.IsNotExist(err), "file should be removed")
	})
}

// TestCloseAndRemove_Together tests using both functions together
func TestCloseAndRemove_Together(t *testing.T) {
	t.Run("close and remove file in sequence", func(t *testing.T) {
		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test-close-remove-*.txt")
		require.NoError(t, err)
		tmpPath := tmpFile.Name()

		// Write some data
		_, err = tmpFile.WriteString("test data")
		require.NoError(t, err)

		// Close and remove in sequence
		pipeline := IO.Chain(func(f *os.File) IO.IO[string] {
			return Remove(f.Name())
		})(Close(tmpFile))

		result := pipeline()

		assert.Equal(t, tmpPath, result, "should return the path")

		// Verify file is removed
		_, err = os.Stat(tmpPath)
		assert.True(t, os.IsNotExist(err), "file should be removed")
	})
}

// TestClose_TypeSafety tests that Close works with different io.Closer types
func TestClose_TypeSafety(t *testing.T) {
	t.Run("works with different closer types", func(t *testing.T) {
		// Test with different types that implement io.Closer
		types := []io.Closer{
			&mockCloser{},
			io.NopCloser(bytes.NewBuffer(nil)),
		}

		for _, closer := range types {
			closeIO := Close(closer)
			result := closeIO()
			assert.Equal(t, closer, result, "should return the same closer")
		}
	})
}

// Example_close demonstrates basic usage of Close
func Example_close() {
	// Create a mock closer
	mock := &mockCloser{}

	// Create an IO that closes the resource
	closeIO := Close(mock)

	// Execute the IO
	result := closeIO()

	fmt.Printf("Closed: %v\n", result.closed)
	// Output: Closed: true
}

// Example_remove demonstrates basic usage of Remove
func Example_remove() {
	// Create a temporary file
	tmpFile, _ := os.CreateTemp("", "example-*.txt")
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// Create an IO that removes the file
	removeIO := Remove(tmpPath)

	// Execute the IO
	path := removeIO()

	// Check if file exists
	_, err := os.Stat(path)
	fmt.Printf("File removed: %v\n", os.IsNotExist(err))
	// Output: File removed: true
}

// Example_closeAndRemove demonstrates using Close and Remove together
func Example_closeAndRemove() {
	// Create a temporary file
	tmpFile, _ := os.CreateTemp("", "example-*.txt")

	// Create a pipeline that closes and removes the file
	pipeline := IO.Chain(func(f *os.File) IO.IO[string] {
		return Remove(f.Name())
	})(Close(tmpFile))

	// Execute the pipeline
	path := pipeline()

	// Check if file exists
	_, err := os.Stat(path)
	fmt.Printf("File removed: %v\n", os.IsNotExist(err))
	// Output: File removed: true
}
