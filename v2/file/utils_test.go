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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopReadCloser_Success(t *testing.T) {
	t.Run("wraps strings.Reader successfully", func(t *testing.T) {
		// Arrange
		content := "hello world"
		reader := strings.NewReader(content)

		// Act
		readCloser := NopReadCloser(reader)

		// Assert
		assert.NotNil(t, readCloser)
		data, err := io.ReadAll(readCloser)
		assert.NoError(t, err)
		assert.Equal(t, content, string(data))
	})

	t.Run("wraps bytes.Buffer successfully", func(t *testing.T) {
		// Arrange
		content := []byte("test data")
		buf := bytes.NewBuffer(content)

		// Act
		readCloser := NopReadCloser(buf)

		// Assert
		assert.NotNil(t, readCloser)
		data, err := io.ReadAll(readCloser)
		assert.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("reads empty content", func(t *testing.T) {
		// Arrange
		reader := strings.NewReader("")

		// Act
		readCloser := NopReadCloser(reader)

		// Assert
		data, err := io.ReadAll(readCloser)
		assert.NoError(t, err)
		assert.Empty(t, data)
	})

	t.Run("reads large content", func(t *testing.T) {
		// Arrange
		content := strings.Repeat("x", 10000)
		reader := strings.NewReader(content)

		// Act
		readCloser := NopReadCloser(reader)

		// Assert
		data, err := io.ReadAll(readCloser)
		assert.NoError(t, err)
		assert.Equal(t, 10000, len(data))
		assert.Equal(t, content, string(data))
	})
}

func TestNopReadCloser_Close(t *testing.T) {
	t.Run("Close returns nil", func(t *testing.T) {
		// Arrange
		reader := strings.NewReader("test")
		readCloser := NopReadCloser(reader)

		// Act
		err := readCloser.Close()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Close can be called multiple times", func(t *testing.T) {
		// Arrange
		reader := strings.NewReader("test")
		readCloser := NopReadCloser(reader)

		// Act & Assert
		assert.NoError(t, readCloser.Close())
		assert.NoError(t, readCloser.Close())
		assert.NoError(t, readCloser.Close())
	})

	t.Run("Close does not affect reading", func(t *testing.T) {
		// Arrange
		content := "test data"
		reader := strings.NewReader(content)
		readCloser := NopReadCloser(reader)

		// Act
		err := readCloser.Close()
		assert.NoError(t, err)

		// Assert - can still read after close
		data, err := io.ReadAll(readCloser)
		assert.NoError(t, err)
		assert.Equal(t, content, string(data))
	})
}

func TestNopReadCloser_Integration(t *testing.T) {
	t.Run("works with defer pattern", func(t *testing.T) {
		// Arrange
		content := "deferred close"
		reader := strings.NewReader(content)

		// Act
		func() {
			readCloser := NopReadCloser(reader)
			defer readCloser.Close() // Should not panic

			data, err := io.ReadAll(readCloser)
			assert.NoError(t, err)
			assert.Equal(t, content, string(data))
		}()
	})

	t.Run("satisfies io.ReadCloser interface", func(t *testing.T) {
		// Arrange
		reader := strings.NewReader("interface test")

		// Act
		var readCloser io.ReadCloser = NopReadCloser(reader)

		// Assert
		assert.NotNil(t, readCloser)
		data, err := io.ReadAll(readCloser)
		assert.NoError(t, err)
		assert.Equal(t, "interface test", string(data))
	})
}

func TestNopWriteCloser_Success(t *testing.T) {
	t.Run("wraps bytes.Buffer successfully", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		content := []byte("hello world")

		// Act
		writeCloser := NopWriteCloser(&buf)

		// Assert
		assert.NotNil(t, writeCloser)
		n, err := writeCloser.Write(content)
		assert.NoError(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, buf.Bytes())
	})

	t.Run("writes empty content", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer

		// Act
		writeCloser := NopWriteCloser(&buf)
		n, err := writeCloser.Write([]byte{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, n)
		assert.Empty(t, buf.Bytes())
	})

	t.Run("writes large content", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		content := bytes.Repeat([]byte("x"), 10000)

		// Act
		writeCloser := NopWriteCloser(&buf)
		n, err := writeCloser.Write(content)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 10000, n)
		assert.Equal(t, content, buf.Bytes())
	})

	t.Run("writes multiple times", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		writeCloser := NopWriteCloser(&buf)

		// Act
		n1, err1 := writeCloser.Write([]byte("hello "))
		n2, err2 := writeCloser.Write([]byte("world"))

		// Assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 6, n1)
		assert.Equal(t, 5, n2)
		assert.Equal(t, "hello world", buf.String())
	})
}

func TestNopWriteCloser_Close(t *testing.T) {
	t.Run("Close returns nil", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		writeCloser := NopWriteCloser(&buf)

		// Act
		err := writeCloser.Close()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Close can be called multiple times", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		writeCloser := NopWriteCloser(&buf)

		// Act & Assert
		assert.NoError(t, writeCloser.Close())
		assert.NoError(t, writeCloser.Close())
		assert.NoError(t, writeCloser.Close())
	})

	t.Run("Close does not affect writing", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		writeCloser := NopWriteCloser(&buf)
		content := []byte("test data")

		// Act
		err := writeCloser.Close()
		assert.NoError(t, err)

		// Assert - can still write after close
		n, err := writeCloser.Write(content)
		assert.NoError(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, buf.Bytes())
	})
}

func TestNopWriteCloser_Integration(t *testing.T) {
	t.Run("works with defer pattern", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		content := []byte("deferred close")

		// Act
		func() {
			writeCloser := NopWriteCloser(&buf)
			defer writeCloser.Close() // Should not panic

			n, err := writeCloser.Write(content)
			assert.NoError(t, err)
			assert.Equal(t, len(content), n)
		}()

		// Assert
		assert.Equal(t, content, buf.Bytes())
	})

	t.Run("satisfies io.WriteCloser interface", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		content := []byte("interface test")

		// Act
		var writeCloser io.WriteCloser = NopWriteCloser(&buf)

		// Assert
		assert.NotNil(t, writeCloser)
		n, err := writeCloser.Write(content)
		assert.NoError(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, buf.Bytes())
	})

	t.Run("works with io.Copy", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		content := "copy test"
		reader := strings.NewReader(content)
		writeCloser := NopWriteCloser(&buf)

		// Act
		n, err := io.Copy(writeCloser, reader)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(len(content)), n)
		assert.Equal(t, content, buf.String())
	})
}

func TestNopWriteCloser_EdgeCases(t *testing.T) {
	t.Run("handles nil write gracefully", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		writeCloser := NopWriteCloser(&buf)

		// Act
		n, err := writeCloser.Write(nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, n)
		assert.Empty(t, buf.Bytes())
	})

	t.Run("preserves underlying writer behavior", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		buf.Grow(100) // Pre-allocate
		writeCloser := NopWriteCloser(&buf)
		content := []byte("test")

		// Act
		n, err := writeCloser.Write(content)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, buf.Bytes())
		assert.GreaterOrEqual(t, buf.Cap(), 100) // Capacity preserved
	})
}

func TestNopReadCloser_NopWriteCloser_Together(t *testing.T) {
	t.Run("can pipe from NopReadCloser to NopWriteCloser", func(t *testing.T) {
		// Arrange
		content := "pipe test data"
		reader := strings.NewReader(content)
		var buf bytes.Buffer

		readCloser := NopReadCloser(reader)
		writeCloser := NopWriteCloser(&buf)

		// Act
		n, err := io.Copy(writeCloser, readCloser)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(len(content)), n)
		assert.Equal(t, content, buf.String())

		// Cleanup
		assert.NoError(t, readCloser.Close())
		assert.NoError(t, writeCloser.Close())
	})
}
