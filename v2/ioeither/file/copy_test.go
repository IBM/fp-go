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
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCopyFileSuccess tests successful file copying
func TestCopyFileSuccess(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file with test content
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := []byte("Hello, CopyFile! This is test content.")
	err := os.WriteFile(srcPath, testContent, 0644)
	require.NoError(t, err)

	// Copy to destination
	dstPath := filepath.Join(tempDir, "destination.txt")
	result := CopyFile(srcPath)(dstPath)()

	// Verify success
	assert.True(t, E.IsRight(result))
	returnedPath := E.GetOrElse(func(error) string { return "" })(result)
	assert.Equal(t, dstPath, returnedPath)

	// Verify destination file content matches source
	dstContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, testContent, dstContent)
}

// TestCopyFileEmptyFile tests copying an empty file
func TestCopyFileEmptyFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create empty source file
	srcPath := filepath.Join(tempDir, "empty_source.txt")
	err := os.WriteFile(srcPath, []byte{}, 0644)
	require.NoError(t, err)

	// Copy to destination
	dstPath := filepath.Join(tempDir, "empty_destination.txt")
	result := CopyFile(srcPath)(dstPath)()

	// Verify success
	assert.True(t, E.IsRight(result))

	// Verify destination file is also empty
	dstContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, 0, len(dstContent))
}

// TestCopyFileLargeFile tests copying a larger file
func TestCopyFileLargeFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file with larger content (1MB)
	srcPath := filepath.Join(tempDir, "large_source.txt")
	largeContent := make([]byte, 1024*1024) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	err := os.WriteFile(srcPath, largeContent, 0644)
	require.NoError(t, err)

	// Copy to destination
	dstPath := filepath.Join(tempDir, "large_destination.txt")
	result := CopyFile(srcPath)(dstPath)()

	// Verify success
	assert.True(t, E.IsRight(result))

	// Verify destination file content matches source
	dstContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, largeContent, dstContent)
}

// TestCopyFileSourceNotFound tests error when source file doesn't exist
func TestCopyFileSourceNotFound(t *testing.T) {
	tempDir := t.TempDir()

	srcPath := filepath.Join(tempDir, "nonexistent_source.txt")
	dstPath := filepath.Join(tempDir, "destination.txt")

	result := CopyFile(srcPath)(dstPath)()

	// Verify failure
	assert.True(t, E.IsLeft(result))
	err := E.Fold(func(e error) error { return e }, func(string) error { return nil })(result)
	assert.Error(t, err)
}

// TestCopyFileDestinationDirectoryNotFound tests error when destination directory doesn't exist
func TestCopyFileDestinationDirectoryNotFound(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	err := os.WriteFile(srcPath, []byte("test"), 0644)
	require.NoError(t, err)

	// Try to copy to non-existent directory
	dstPath := filepath.Join(tempDir, "nonexistent_dir", "destination.txt")
	result := CopyFile(srcPath)(dstPath)()

	// Verify failure
	assert.True(t, E.IsLeft(result))
}

// TestCopyFileOverwriteExisting tests overwriting an existing destination file
func TestCopyFileOverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	newContent := []byte("New content")
	err := os.WriteFile(srcPath, newContent, 0644)
	require.NoError(t, err)

	// Create existing destination file with different content
	dstPath := filepath.Join(tempDir, "destination.txt")
	oldContent := []byte("Old content that should be replaced")
	err = os.WriteFile(dstPath, oldContent, 0644)
	require.NoError(t, err)

	// Copy and overwrite
	result := CopyFile(srcPath)(dstPath)()

	// Verify success
	assert.True(t, E.IsRight(result))

	// Verify destination has new content
	dstContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, newContent, dstContent)
}

// TestCopyFileCurrying tests the curried nature of CopyFile (data-last pattern)
func TestCopyFileCurrying(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := []byte("Currying test content")
	err := os.WriteFile(srcPath, testContent, 0644)
	require.NoError(t, err)

	// Create a partially applied function
	copyFromSource := CopyFile(srcPath)

	// Use the partially applied function multiple times
	dst1 := filepath.Join(tempDir, "dest1.txt")
	dst2 := filepath.Join(tempDir, "dest2.txt")

	result1 := copyFromSource(dst1)()
	result2 := copyFromSource(dst2)()

	// Verify both copies succeeded
	assert.True(t, E.IsRight(result1))
	assert.True(t, E.IsRight(result2))

	// Verify both destinations have the same content
	content1, err := os.ReadFile(dst1)
	require.NoError(t, err)
	content2, err := os.ReadFile(dst2)
	require.NoError(t, err)
	assert.Equal(t, testContent, content1)
	assert.Equal(t, testContent, content2)
}

// TestCopyFileComposition tests composing CopyFile with other operations
func TestCopyFileComposition(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := []byte("Composition test")
	err := os.WriteFile(srcPath, testContent, 0644)
	require.NoError(t, err)

	dstPath := filepath.Join(tempDir, "destination.txt")

	// Compose CopyFile with Map to transform the result
	result := F.Pipe1(
		CopyFile(srcPath)(dstPath),
		IOE.Map[error](func(dst string) string {
			return "Successfully copied to: " + dst
		}),
	)()

	// Verify success and transformation
	assert.True(t, E.IsRight(result))
	message := E.GetOrElse(func(error) string { return "" })(result)
	assert.Equal(t, "Successfully copied to: "+dstPath, message)

	// Verify file was actually copied
	dstContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, testContent, dstContent)
}

// TestCopyFileChaining tests chaining multiple copy operations
func TestCopyFileChaining(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := []byte("Chaining test")
	err := os.WriteFile(srcPath, testContent, 0644)
	require.NoError(t, err)

	dst1Path := filepath.Join(tempDir, "dest1.txt")
	dst2Path := filepath.Join(tempDir, "dest2.txt")

	// Chain two copy operations
	result := F.Pipe1(
		CopyFile(srcPath)(dst1Path),
		IOE.Chain(func(string) IOEither[error, string] {
			return CopyFile(dst1Path)(dst2Path)
		}),
	)()

	// Verify success
	assert.True(t, E.IsRight(result))

	// Verify both files exist with correct content
	content1, err := os.ReadFile(dst1Path)
	require.NoError(t, err)
	assert.Equal(t, testContent, content1)

	content2, err := os.ReadFile(dst2Path)
	require.NoError(t, err)
	assert.Equal(t, testContent, content2)
}

// TestCopyFileWithBinaryContent tests copying binary files
func TestCopyFileWithBinaryContent(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file with binary content
	srcPath := filepath.Join(tempDir, "binary_source.bin")
	binaryContent := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD, 0x7F, 0x80}
	err := os.WriteFile(srcPath, binaryContent, 0644)
	require.NoError(t, err)

	// Copy to destination
	dstPath := filepath.Join(tempDir, "binary_destination.bin")
	result := CopyFile(srcPath)(dstPath)()

	// Verify success
	assert.True(t, E.IsRight(result))

	// Verify binary content is preserved
	dstContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, binaryContent, dstContent)
}

// TestCopyFileErrorHandling tests error handling with Either operations
func TestCopyFileErrorHandling(t *testing.T) {
	tempDir := t.TempDir()

	srcPath := filepath.Join(tempDir, "nonexistent.txt")
	dstPath := filepath.Join(tempDir, "destination.txt")

	result := CopyFile(srcPath)(dstPath)()

	// Test error handling with Fold
	message := E.Fold(
		func(err error) string { return "Error: " + err.Error() },
		func(dst string) string { return "Success: " + dst },
	)(result)

	assert.Contains(t, message, "Error:")
}

// TestCopyFileResourceCleanup tests that resources are properly cleaned up
func TestCopyFileResourceCleanup(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := []byte("Resource cleanup test")
	err := os.WriteFile(srcPath, testContent, 0644)
	require.NoError(t, err)

	dstPath := filepath.Join(tempDir, "destination.txt")

	// Perform copy
	result := CopyFile(srcPath)(dstPath)()
	assert.True(t, E.IsRight(result))

	// Verify we can immediately delete both files (no file handles left open)
	err = os.Remove(srcPath)
	assert.NoError(t, err, "Source file should be closed and deletable")

	err = os.Remove(dstPath)
	assert.NoError(t, err, "Destination file should be closed and deletable")
}

// TestCopyFileMultipleOperations tests using CopyFile multiple times independently
func TestCopyFileMultipleOperations(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple source files
	src1 := filepath.Join(tempDir, "source1.txt")
	src2 := filepath.Join(tempDir, "source2.txt")
	content1 := []byte("Content 1")
	content2 := []byte("Content 2")

	err := os.WriteFile(src1, content1, 0644)
	require.NoError(t, err)
	err = os.WriteFile(src2, content2, 0644)
	require.NoError(t, err)

	// Perform multiple independent copies
	dst1 := filepath.Join(tempDir, "dest1.txt")
	dst2 := filepath.Join(tempDir, "dest2.txt")

	result1 := CopyFile(src1)(dst1)()
	result2 := CopyFile(src2)(dst2)()

	// Verify both succeeded
	assert.True(t, E.IsRight(result1))
	assert.True(t, E.IsRight(result2))

	// Verify correct content in each destination
	dstContent1, err := os.ReadFile(dst1)
	require.NoError(t, err)
	assert.Equal(t, content1, dstContent1)

	dstContent2, err := os.ReadFile(dst2)
	require.NoError(t, err)
	assert.Equal(t, content2, dstContent2)
}
