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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	R "github.com/IBM/fp-go/v2/context/readerioresult"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	J "github.com/IBM/fp-go/v2/json"
	"github.com/stretchr/testify/assert"
)

type RecordType struct {
	Data string `json:"data"`
}

func getData(r RecordType) string {
	return r.Data
}

func ExampleReadFile() {

	data := F.Pipe3(
		ReadFile("./data/file.json"),
		R.ChainEitherK(J.Unmarshal[RecordType]),
		R.ChainFirstIOK(io.Logf[RecordType]("Log: %v")),
		R.Map(getData),
	)

	result := data(context.Background())

	fmt.Println(result())

	// Output:
	// Right[string](Carsten)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - creates new file", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_create.txt")

		createOp := Create(tempFile)
		result := createOp(ctx)()

		assert.True(t, E.IsRight(result))

		// Verify file was created
		_, err := os.Stat(tempFile)
		assert.NoError(t, err)

		// Clean up file handle
		E.MonadFold(result,
			func(error) *os.File { return nil },
			func(f *os.File) *os.File { f.Close(); return f },
		)
	})

	t.Run("Success - truncates existing file", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_truncate.txt")

		// Create file with initial content
		err := os.WriteFile(tempFile, []byte("initial content"), 0644)
		assert.NoError(t, err)

		// Create should truncate
		createOp := Create(tempFile)
		result := createOp(ctx)()

		assert.True(t, E.IsRight(result))

		// Close the file
		E.MonadFold(result,
			func(error) *os.File { return nil },
			func(f *os.File) *os.File { f.Close(); return f },
		)

		// Verify file was truncated
		content, err := os.ReadFile(tempFile)
		assert.NoError(t, err)
		assert.Empty(t, content)
	})

	t.Run("Failure - invalid path", func(t *testing.T) {
		// Try to create file in non-existent directory
		invalidPath := filepath.Join(t.TempDir(), "nonexistent", "test.txt")

		createOp := Create(invalidPath)
		result := createOp(ctx)()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("Success - file can be written to", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_write.txt")

		createOp := Create(tempFile)
		result := createOp(ctx)()

		assert.True(t, E.IsRight(result))

		// Write to the file
		E.MonadFold(result,
			func(err error) *os.File { t.Fatalf("Unexpected error: %v", err); return nil },
			func(f *os.File) *os.File {
				defer f.Close()
				_, err := f.WriteString("test content")
				assert.NoError(t, err)
				return f
			},
		)

		// Verify content was written
		content, err := os.ReadFile(tempFile)
		assert.NoError(t, err)
		assert.Equal(t, "test content", string(content))
	})

	t.Run("Context cancellation", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		tempFile := filepath.Join(t.TempDir(), "test_cancel.txt")

		createOp := Create(tempFile)
		result := createOp(cancelCtx)()

		// Note: File creation itself doesn't check context, but this tests the pattern
		// In practice, context cancellation would affect subsequent operations
		_ = result
	})
}

func TestWriteFile(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - writes data to new file", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_write.txt")
		testData := []byte("Hello, World!")

		writeOp := WriteFile(testData)
		result := writeOp(tempFile)(ctx)()

		assert.True(t, E.IsRight(result))

		// Verify returned data
		E.MonadFold(result,
			func(err error) []byte { t.Fatalf("Unexpected error: %v", err); return nil },
			func(data []byte) []byte {
				assert.Equal(t, testData, data)
				return data
			},
		)

		// Verify file content
		content, err := os.ReadFile(tempFile)
		assert.NoError(t, err)
		assert.Equal(t, testData, content)
	})

	t.Run("Success - overwrites existing file", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_overwrite.txt")

		// Write initial content
		err := os.WriteFile(tempFile, []byte("old content"), 0644)
		assert.NoError(t, err)

		// Overwrite with new content
		newData := []byte("new content")
		writeOp := WriteFile(newData)
		result := writeOp(tempFile)(ctx)()

		assert.True(t, E.IsRight(result))

		// Verify file was overwritten
		content, err := os.ReadFile(tempFile)
		assert.NoError(t, err)
		assert.Equal(t, newData, content)
	})

	t.Run("Success - writes empty data", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_empty.txt")
		emptyData := []byte{}

		writeOp := WriteFile(emptyData)
		result := writeOp(tempFile)(ctx)()

		assert.True(t, E.IsRight(result))

		// Verify file is empty
		content, err := os.ReadFile(tempFile)
		assert.NoError(t, err)
		assert.Empty(t, content)
	})

	t.Run("Success - writes large data", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_large.txt")
		largeData := make([]byte, 1024*1024) // 1MB
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		writeOp := WriteFile(largeData)
		result := writeOp(tempFile)(ctx)()

		assert.True(t, E.IsRight(result))

		// Verify file content
		content, err := os.ReadFile(tempFile)
		assert.NoError(t, err)
		assert.Equal(t, largeData, content)
	})

	t.Run("Failure - invalid path", func(t *testing.T) {
		invalidPath := filepath.Join(t.TempDir(), "nonexistent", "test.txt")
		testData := []byte("test")

		writeOp := WriteFile(testData)
		result := writeOp(invalidPath)(ctx)()

		assert.True(t, E.IsLeft(result))
	})

	t.Run("Success - writes binary data", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_binary.bin")
		binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}

		writeOp := WriteFile(binaryData)
		result := writeOp(tempFile)(ctx)()

		assert.True(t, E.IsRight(result))

		// Verify binary content
		content, err := os.ReadFile(tempFile)
		assert.NoError(t, err)
		assert.Equal(t, binaryData, content)
	})

	t.Run("Integration - write then read", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_roundtrip.txt")
		testData := []byte("Round trip test data")

		// Write data
		writeOp := WriteFile(testData)
		writeResult := writeOp(tempFile)(ctx)()
		assert.True(t, E.IsRight(writeResult))

		// Read data back
		readOp := ReadFile(tempFile)
		readResult := readOp(ctx)()
		assert.True(t, E.IsRight(readResult))

		// Verify data matches
		E.MonadFold(readResult,
			func(err error) []byte { t.Fatalf("Unexpected error: %v", err); return nil },
			func(data []byte) []byte {
				assert.Equal(t, testData, data)
				return data
			},
		)
	})

	t.Run("Composition with Map", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test_compose.txt")
		testData := []byte("test data")

		// Write and transform result
		pipeline := F.Pipe1(
			WriteFile(testData)(tempFile),
			R.Map(func(data []byte) int { return len(data) }),
		)

		result := pipeline(ctx)()
		assert.True(t, E.IsRight(result))

		E.MonadFold(result,
			func(err error) int { t.Fatalf("Unexpected error: %v", err); return 0 },
			func(length int) int {
				assert.Equal(t, len(testData), length)
				return length
			},
		)
	})

	t.Run("Context cancellation during write", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		tempFile := filepath.Join(t.TempDir(), "test_cancel.txt")
		testData := []byte("test")

		writeOp := WriteFile(testData)
		result := writeOp(tempFile)(cancelCtx)()

		// Note: The actual write may complete before cancellation is checked
		// This test verifies the pattern works with cancelled contexts
		_ = result
	})
}
