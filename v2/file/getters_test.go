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
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	t.Run("joins simple paths", func(t *testing.T) {
		result := Join("config.json")("/etc/myapp")
		expected := filepath.Join("/etc/myapp", "config.json")
		assert.Equal(t, expected, result)
	})

	t.Run("joins with subdirectories", func(t *testing.T) {
		result := Join("logs/app.log")("/var")
		expected := filepath.Join("/var", "logs/app.log")
		assert.Equal(t, expected, result)
	})

	t.Run("handles empty root", func(t *testing.T) {
		result := Join("file.txt")("")
		assert.Equal(t, "file.txt", result)
	})

	t.Run("handles empty name", func(t *testing.T) {
		result := Join("")("/root")
		expected := filepath.Join("/root", "")
		assert.Equal(t, expected, result)
	})

	t.Run("handles relative paths", func(t *testing.T) {
		result := Join("config.json")("./app")
		expected := filepath.Join("./app", "config.json")
		assert.Equal(t, expected, result)
	})

	t.Run("normalizes path separators", func(t *testing.T) {
		result := Join("file.txt")("/root/path")
		// Should use OS-specific separator
		assert.Contains(t, result, "file.txt")
		assert.Contains(t, result, "root")
		assert.Contains(t, result, "path")
	})

	t.Run("works with Pipe", func(t *testing.T) {
		result := F.Pipe1("/var/log", Join("app.log"))
		expected := filepath.Join("/var/log", "app.log")
		assert.Equal(t, expected, result)
	})

	t.Run("chains multiple joins", func(t *testing.T) {
		result := F.Pipe2(
			"/root",
			Join("subdir"),
			Join("file.txt"),
		)
		expected := filepath.Join("/root", "subdir", "file.txt")
		assert.Equal(t, expected, result)
	})

	t.Run("handles special characters", func(t *testing.T) {
		result := Join("my file.txt")("/path with spaces")
		expected := filepath.Join("/path with spaces", "my file.txt")
		assert.Equal(t, expected, result)
	})

	t.Run("handles dots in path", func(t *testing.T) {
		result := Join("../config.json")("/app/current")
		expected := filepath.Join("/app/current", "../config.json")
		assert.Equal(t, expected, result)
	})
}

func TestToReader(t *testing.T) {
	t.Run("converts bytes.Buffer to io.Reader", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte("hello world"))
		reader := ToReader(buf)

		// Verify it's an io.Reader
		var _ io.Reader = reader

		// Verify it works
		data, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", string(data))
	})

	t.Run("converts bytes.Reader to io.Reader", func(t *testing.T) {
		bytesReader := bytes.NewReader([]byte("test data"))
		reader := ToReader(bytesReader)

		var _ io.Reader = reader

		data, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "test data", string(data))
	})

	t.Run("converts strings.Reader to io.Reader", func(t *testing.T) {
		strReader := strings.NewReader("string content")
		reader := ToReader(strReader)

		var _ io.Reader = reader

		data, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "string content", string(data))
	})

	t.Run("preserves reader functionality", func(t *testing.T) {
		original := bytes.NewBuffer([]byte("test"))
		reader := ToReader(original)

		// Read once
		buf1 := make([]byte, 2)
		n, err := reader.Read(buf1)
		assert.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.Equal(t, "te", string(buf1))

		// Read again
		buf2 := make([]byte, 2)
		n, err = reader.Read(buf2)
		assert.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.Equal(t, "st", string(buf2))
	})

	t.Run("handles empty reader", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		reader := ToReader(buf)

		data, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "", string(data))
	})
}

func TestToWriter(t *testing.T) {
	t.Run("converts bytes.Buffer to io.Writer", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writer := ToWriter(buf)

		// Verify it's an io.Writer
		var _ io.Writer = writer

		// Verify it works
		n, err := writer.Write([]byte("hello"))
		assert.NoError(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, "hello", buf.String())
	})

	t.Run("preserves writer functionality", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writer := ToWriter(buf)

		// Write multiple times
		writer.Write([]byte("hello "))
		writer.Write([]byte("world"))

		assert.Equal(t, "hello world", buf.String())
	})

	t.Run("handles empty writes", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writer := ToWriter(buf)

		n, err := writer.Write([]byte{})
		assert.NoError(t, err)
		assert.Equal(t, 0, n)
		assert.Equal(t, "", buf.String())
	})

	t.Run("handles large writes", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writer := ToWriter(buf)

		data := make([]byte, 10000)
		for i := range data {
			data[i] = byte('A' + (i % 26))
		}

		n, err := writer.Write(data)
		assert.NoError(t, err)
		assert.Equal(t, 10000, n)
		assert.Equal(t, 10000, buf.Len())
	})
}

func TestToCloser(t *testing.T) {
	t.Run("converts file to io.Closer", func(t *testing.T) {
		// Create a temporary file
		tmpfile, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		closer := ToCloser(tmpfile)

		// Verify it's an io.Closer
		var _ io.Closer = closer

		// Verify it works
		err = closer.Close()
		assert.NoError(t, err)
	})

	t.Run("converts nopCloser to io.Closer", func(t *testing.T) {
		// Use io.NopCloser which is a standard implementation
		reader := strings.NewReader("test")
		nopCloser := io.NopCloser(reader)

		closer := ToCloser(nopCloser)
		var _ io.Closer = closer

		err := closer.Close()
		assert.NoError(t, err)
	})

	t.Run("preserves close functionality", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		closer := ToCloser(tmpfile)

		// Close should work
		err = closer.Close()
		assert.NoError(t, err)

		// Subsequent operations should fail
		_, err = tmpfile.Write([]byte("test"))
		assert.Error(t, err)
	})
}

// Test type conversions work together
func TestIntegration(t *testing.T) {
	t.Run("reader and closer together", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		// Write some data
		tmpfile.Write([]byte("test content"))
		tmpfile.Seek(0, 0)

		// Convert to interfaces
		reader := ToReader(tmpfile)
		closer := ToCloser(tmpfile)

		// Use as reader
		data, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "test content", string(data))

		// Close
		err = closer.Close()
		assert.NoError(t, err)
	})

	t.Run("writer and closer together", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		// Convert to interfaces
		writer := ToWriter(tmpfile)
		closer := ToCloser(tmpfile)

		// Use as writer
		n, err := writer.Write([]byte("test data"))
		assert.NoError(t, err)
		assert.Equal(t, 9, n)

		// Close
		err = closer.Close()
		assert.NoError(t, err)

		// Verify data was written
		data, err := os.ReadFile(tmpfile.Name())
		assert.NoError(t, err)
		assert.Equal(t, "test data", string(data))
	})

	t.Run("all conversions with file", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		// File implements Reader, Writer, and Closer
		var reader io.Reader = ToReader(tmpfile)
		var writer io.Writer = ToWriter(tmpfile)
		var closer io.Closer = ToCloser(tmpfile)

		// All should be non-nil
		assert.NotNil(t, reader)
		assert.NotNil(t, writer)
		assert.NotNil(t, closer)

		// Write, read, close
		writer.Write([]byte("hello"))
		tmpfile.Seek(0, 0)
		data, _ := io.ReadAll(reader)
		assert.Equal(t, "hello", string(data))
		closer.Close()
	})
}

// Benchmark tests
func BenchmarkJoin(b *testing.B) {
	joiner := Join("config.json")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = joiner("/etc/myapp")
	}
}

func BenchmarkToReader(b *testing.B) {
	buf := bytes.NewBuffer([]byte("test data"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToReader(buf)
	}
}

func BenchmarkToWriter(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToWriter(buf)
	}
}

func BenchmarkToCloser(b *testing.B) {
	tmpfile, _ := os.CreateTemp("", "bench")
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToCloser(tmpfile)
	}
}
