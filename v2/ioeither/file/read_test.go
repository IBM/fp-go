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
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockReadCloser is a mock implementation of io.ReadCloser for testing
type mockReadCloser struct {
	data        []byte
	readErr     error
	closeErr    error
	readPos     int
	closeCalled bool
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	if m.readErr != nil {
		return 0, m.readErr
	}
	if m.readPos >= len(m.data) {
		return 0, io.EOF
	}
	n = copy(p, m.data[m.readPos:])
	m.readPos += n
	if m.readPos >= len(m.data) {
		return n, io.EOF
	}
	return n, nil
}

func (m *mockReadCloser) Close() error {
	m.closeCalled = true
	return m.closeErr
}

// TestReadSuccessfulRead tests reading data successfully from a ReadCloser
func TestReadSuccessfulRead(t *testing.T) {
	testData := []byte("Hello, World!")
	mock := &mockReadCloser{data: testData}

	// Create an acquire function that returns our mock
	acquire := ioeither.Of[error](mock)

	// Create a reader function that reads all data
	reader := func(rc *mockReadCloser) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return io.ReadAll(rc)
		})
	}

	// Execute the Read operation
	result := Read[[]byte](acquire)(reader)
	either := result()

	// Assertions
	assert.True(t, E.IsRight(either))
	data := E.GetOrElse(func(error) []byte { return nil })(either)
	assert.Equal(t, testData, data)
	assert.True(t, mock.closeCalled, "Close should have been called")
}

// TestReadWithRealFile tests reading from an actual file
func TestReadWithRealFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "read_test_*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	testContent := []byte("Test file content for Read function")
	_, err = tmpFile.Write(testContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Use Read to read the file
	reader := func(f *os.File) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return io.ReadAll(f)
		})
	}

	result := Read[[]byte](Open(tmpFile.Name()))(reader)
	either := result()

	assert.True(t, E.IsRight(either))
	data := E.GetOrElse(func(error) []byte { return nil })(either)
	assert.Equal(t, testContent, data)
}

// TestReadPartialRead tests reading only part of the data
func TestReadPartialRead(t *testing.T) {
	testData := []byte("Hello, World! This is a longer message.")
	mock := &mockReadCloser{data: testData}

	acquire := ioeither.Of[error](mock)

	// Reader that only reads first 13 bytes
	reader := func(rc *mockReadCloser) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			buf := make([]byte, 13)
			n, err := rc.Read(buf)
			if err != nil && err != io.EOF {
				return nil, err
			}
			return buf[:n], nil
		})
	}

	result := Read[[]byte](acquire)(reader)
	either := result()

	assert.True(t, E.IsRight(either))
	data := E.GetOrElse(func(error) []byte { return nil })(either)
	assert.Equal(t, []byte("Hello, World!"), data)
	assert.True(t, mock.closeCalled, "Close should have been called")
}

// TestReadErrorDuringRead tests that errors during reading are propagated
func TestReadErrorDuringRead(t *testing.T) {
	readError := errors.New("read error")
	mock := &mockReadCloser{
		data:    []byte("data"),
		readErr: readError,
	}

	acquire := ioeither.Of[error](mock)

	reader := func(rc *mockReadCloser) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return io.ReadAll(rc)
		})
	}

	result := Read[[]byte](acquire)(reader)
	either := result()

	assert.True(t, E.IsLeft(either))
	err := E.Fold(func(e error) error { return e }, func([]byte) error { return nil })(either)
	assert.Equal(t, readError, err)
	assert.True(t, mock.closeCalled, "Close should be called even on read error")
}

// TestReadErrorDuringClose tests that errors during close are handled
func TestReadErrorDuringClose(t *testing.T) {
	closeError := errors.New("close error")
	mock := &mockReadCloser{
		data:     []byte("Hello"),
		closeErr: closeError,
	}

	acquire := ioeither.Of[error](mock)

	reader := func(rc *mockReadCloser) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return io.ReadAll(rc)
		})
	}

	result := Read[[]byte](acquire)(reader)
	either := result()

	// The close error should be propagated
	assert.True(t, E.IsLeft(either))
	assert.True(t, mock.closeCalled, "Close should have been called")
}

// TestReadErrorDuringAcquire tests that errors during resource acquisition are propagated
func TestReadErrorDuringAcquire(t *testing.T) {
	acquireError := errors.New("acquire error")
	acquire := ioeither.Left[*mockReadCloser](acquireError)

	reader := func(rc *mockReadCloser) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return io.ReadAll(rc)
		})
	}

	result := Read[[]byte](acquire)(reader)
	either := result()

	assert.True(t, E.IsLeft(either))
	err := E.Fold(func(e error) error { return e }, func([]byte) error { return nil })(either)
	assert.Equal(t, acquireError, err)
}

// TestReadWithStringReader tests reading and transforming to a different type
func TestReadWithStringReader(t *testing.T) {
	testData := []byte("Hello, World!")
	mock := &mockReadCloser{data: testData}

	acquire := ioeither.Of[error](mock)

	// Reader that converts bytes to uppercase string
	reader := func(rc *mockReadCloser) IOEither[error, string] {
		return ioeither.TryCatchError(func() (string, error) {
			data, err := io.ReadAll(rc)
			if err != nil {
				return "", err
			}
			return strings.ToUpper(string(data)), nil
		})
	}

	result := Read[string](acquire)(reader)
	either := result()

	assert.True(t, E.IsRight(either))
	str := E.GetOrElse(func(error) string { return "" })(either)
	assert.Equal(t, "HELLO, WORLD!", str)
	assert.True(t, mock.closeCalled, "Close should have been called")
}

// TestReadComposition tests composing Read with other operations
func TestReadComposition(t *testing.T) {
	testData := []byte("42")
	mock := &mockReadCloser{data: testData}

	acquire := ioeither.Of[error](mock)

	// Reader that parses the content as an integer
	reader := func(rc *mockReadCloser) IOEither[error, int] {
		return ioeither.TryCatchError(func() (int, error) {
			data, err := io.ReadAll(rc)
			if err != nil {
				return 0, err
			}
			// Simple parsing
			num := int(data[0]-'0')*10 + int(data[1]-'0')
			return num, nil
		})
	}

	result := F.Pipe1(
		acquire,
		Read[int],
	)(reader)

	either := result()

	assert.True(t, E.IsRight(either))
	num := E.GetOrElse(func(error) int { return 0 })(either)
	assert.Equal(t, 42, num)
	assert.True(t, mock.closeCalled, "Close should have been called")
}

// TestReadMultipleOperations tests that Read can be used multiple times
func TestReadMultipleOperations(t *testing.T) {
	// Create a function that creates a new mock each time
	createMock := func() *mockReadCloser {
		return &mockReadCloser{data: []byte("test data")}
	}

	reader := func(rc *mockReadCloser) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return io.ReadAll(rc)
		})
	}

	// First read
	mock1 := createMock()
	result1 := Read[[]byte](ioeither.Of[error](mock1))(reader)
	either1 := result1()

	assert.True(t, E.IsRight(either1))
	data1 := E.GetOrElse(func(error) []byte { return nil })(either1)
	assert.Equal(t, []byte("test data"), data1)
	assert.True(t, mock1.closeCalled)

	// Second read with a new mock
	mock2 := createMock()
	result2 := Read[[]byte](ioeither.Of[error](mock2))(reader)
	either2 := result2()

	assert.True(t, E.IsRight(either2))
	data2 := E.GetOrElse(func(error) []byte { return nil })(either2)
	assert.Equal(t, []byte("test data"), data2)
	assert.True(t, mock2.closeCalled)
}

// TestReadEnsuresCloseOnPanic tests that Close is called even if reader panics
// Note: This is more of a conceptual test as the actual panic handling depends on
// the implementation of WithResource
func TestReadWithEmptyData(t *testing.T) {
	mock := &mockReadCloser{data: []byte{}}

	acquire := ioeither.Of[error](mock)

	reader := func(rc *mockReadCloser) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return io.ReadAll(rc)
		})
	}

	result := Read[[]byte](acquire)(reader)
	either := result()

	assert.True(t, E.IsRight(either))
	data := E.GetOrElse(func(error) []byte { return nil })(either)
	assert.Empty(t, data)
	assert.True(t, mock.closeCalled, "Close should be called even with empty data")
}

// TestReadIntegrationWithEither tests integration with Either operations
func TestReadIntegrationWithEither(t *testing.T) {
	testData := []byte("integration test")
	mock := &mockReadCloser{data: testData}

	acquire := ioeither.Of[error](mock)

	reader := func(rc *mockReadCloser) IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return io.ReadAll(rc)
		})
	}

	result := Read[[]byte](acquire)(reader)
	either := result()

	// Test with Either operations
	assert.True(t, E.IsRight(either))

	folded := E.Fold(
		func(err error) string { return "error: " + err.Error() },
		func(data []byte) string { return "success: " + string(data) },
	)(either)

	assert.Equal(t, "success: integration test", folded)
	assert.True(t, mock.closeCalled)
}
