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

package readerioresult

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/logging"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestTapSLogComprehensive_Success verifies TapSLog logs successful values
func TestTapSLogComprehensive_Success(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	t.Run("logs integer success value", func(t *testing.T) {
		buf.Reset()

		pipeline := F.Pipe2(
			Of(42),
			TapSLog[int]("Integer value"),
			Map(N.Mul(2)),
		)

		res := pipeline(t.Context())()

		// Verify result is correct
		assert.Equal(t, 84, F.Pipe1(res, getOrZero))

		// Verify logging occurred
		logOutput := buf.String()
		assert.Contains(t, logOutput, "Integer value", "Should log the message")
		assert.Contains(t, logOutput, "value=42", "Should log the success value")
		assert.NotContains(t, logOutput, "error", "Should not contain error keyword for success")
	})

	t.Run("logs string success value", func(t *testing.T) {
		buf.Reset()

		pipeline := F.Pipe1(
			Of("hello world"),
			TapSLog[string]("String value"),
		)

		res := pipeline(t.Context())()

		// Verify result is correct
		assert.True(t, F.Pipe1(res, isRight[string]))

		// Verify logging occurred
		logOutput := buf.String()
		assert.Contains(t, logOutput, "String value")
		assert.Contains(t, logOutput, `value="hello world"`)
	})

	t.Run("logs struct success value", func(t *testing.T) {
		buf.Reset()

		type User struct {
			ID   int
			Name string
		}

		user := User{ID: 123, Name: "Alice"}
		pipeline := F.Pipe1(
			Of(user),
			TapSLog[User]("User struct"),
		)

		res := pipeline(t.Context())()

		// Verify result is correct
		assert.True(t, F.Pipe1(res, isRight[User]))

		// Verify logging occurred with struct fields
		logOutput := buf.String()
		assert.Contains(t, logOutput, "User struct")
		assert.Contains(t, logOutput, "ID:123")
		assert.Contains(t, logOutput, "Name:Alice")
	})

	t.Run("logs multiple success values in pipeline", func(t *testing.T) {
		buf.Reset()

		step1 := F.Pipe2(
			Of(10),
			TapSLog[int]("Initial value"),
			Map(N.Mul(2)),
		)

		pipeline := F.Pipe2(
			step1,
			TapSLog[int]("After doubling"),
			Map(N.Add(5)),
		)

		res := pipeline(t.Context())()

		// Verify result is correct
		assert.Equal(t, 25, getOrZero(res))

		// Verify both log entries
		logOutput := buf.String()
		assert.Contains(t, logOutput, "Initial value")
		assert.Contains(t, logOutput, "value=10")
		assert.Contains(t, logOutput, "After doubling")
		assert.Contains(t, logOutput, "value=20")
	})
}

// TestTapSLogComprehensive_Error verifies TapSLog behavior with errors
func TestTapSLogComprehensive_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	t.Run("logs error values", func(t *testing.T) {
		buf.Reset()

		testErr := errors.New("test error")
		pipeline := F.Pipe2(
			Left[int](testErr),
			TapSLog[int]("Error case"),
			Map(N.Mul(2)),
		)

		res := pipeline(t.Context())()

		// Verify error is preserved
		assert.True(t, F.Pipe1(res, isLeft[int]))

		// Verify logging occurred for error
		logOutput := buf.String()
		assert.Contains(t, logOutput, "Error case", "Should log the message")
		assert.Contains(t, logOutput, "error", "Should contain error keyword")
		assert.Contains(t, logOutput, "test error", "Should log the error message")
		assert.NotContains(t, logOutput, "value=", "Should not log 'value=' for errors")
	})

	t.Run("preserves error through pipeline", func(t *testing.T) {
		buf.Reset()

		originalErr := errors.New("original error")
		step1 := F.Pipe2(
			Left[int](originalErr),
			TapSLog[int]("First tap"),
			Map(N.Mul(2)),
		)

		pipeline := F.Pipe2(
			step1,
			TapSLog[int]("Second tap"),
			Map(N.Add(5)),
		)

		res := pipeline(t.Context())()

		// Verify error is preserved
		assert.True(t, isLeft(res))

		// Verify both taps logged the error
		logOutput := buf.String()
		errorCount := strings.Count(logOutput, "original error")
		assert.Equal(t, 2, errorCount, "Both TapSLog calls should log the error")
		assert.Contains(t, logOutput, "First tap")
		assert.Contains(t, logOutput, "Second tap")
	})

	t.Run("logs error after successful operation", func(t *testing.T) {
		buf.Reset()

		pipeline := F.Pipe3(
			Of(10),
			TapSLog[int]("Before error"),
			Chain(func(n int) ReaderIOResult[int] {
				return Left[int](errors.New("chain error"))
			}),
			TapSLog[int]("After error"),
		)

		res := pipeline(t.Context())()

		// Verify error is present
		assert.True(t, F.Pipe1(res, isLeft[int]))

		// Verify both logs
		logOutput := buf.String()
		assert.Contains(t, logOutput, "Before error")
		assert.Contains(t, logOutput, "value=10")
		assert.Contains(t, logOutput, "After error")
		assert.Contains(t, logOutput, "chain error")
	})
}

// TestTapSLogComprehensive_EdgeCases verifies TapSLog with edge cases
func TestTapSLogComprehensive_EdgeCases(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	t.Run("logs zero value", func(t *testing.T) {
		buf.Reset()

		pipeline := F.Pipe1(
			Of(0),
			TapSLog[int]("Zero value"),
		)

		res := pipeline(t.Context())()

		assert.Equal(t, 0, F.Pipe1(res, getOrZero))

		logOutput := buf.String()
		assert.Contains(t, logOutput, "Zero value")
		assert.Contains(t, logOutput, "value=0")
	})

	t.Run("logs empty string", func(t *testing.T) {
		buf.Reset()

		pipeline := F.Pipe1(
			Of(""),
			TapSLog[string]("Empty string"),
		)

		res := pipeline(t.Context())()

		assert.True(t, F.Pipe1(res, isRight[string]))

		logOutput := buf.String()
		assert.Contains(t, logOutput, "Empty string")
		assert.Contains(t, logOutput, `value=""`)
	})

	t.Run("logs nil pointer", func(t *testing.T) {
		buf.Reset()

		type Data struct {
			Value string
		}

		var nilData *Data
		pipeline := F.Pipe1(
			Of(nilData),
			TapSLog[*Data]("Nil pointer"),
		)

		res := pipeline(t.Context())()

		assert.True(t, F.Pipe1(res, isRight[*Data]))

		logOutput := buf.String()
		assert.Contains(t, logOutput, "Nil pointer")
		// Nil representation may vary, but should be logged
		assert.NotEmpty(t, logOutput)
	})

	t.Run("respects logger level - disabled", func(t *testing.T) {
		buf.Reset()

		// Create logger that only logs errors
		errorLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
		oldLogger := logging.SetLogger(errorLogger)
		defer logging.SetLogger(oldLogger)

		pipeline := F.Pipe1(
			Of(42),
			TapSLog[int]("Should not log"),
		)

		res := pipeline(t.Context())()

		assert.Equal(t, 42, F.Pipe1(res, getOrZero))

		// Should have no logs since level is ERROR
		logOutput := buf.String()
		assert.Empty(t, logOutput, "Should not log when level is disabled")
	})
}

// TestTapSLogComprehensive_Integration verifies TapSLog in realistic scenarios
func TestTapSLogComprehensive_Integration(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	t.Run("complex pipeline with mixed success and error", func(t *testing.T) {
		buf.Reset()

		// Simulate a data processing pipeline
		validatePositive := func(n int) ReaderIOResult[int] {
			if n > 0 {
				return Of(n)
			}
			return Left[int](errors.New("number must be positive"))
		}

		step1 := F.Pipe3(
			Of(5),
			TapSLog[int]("Input received"),
			Map(N.Mul(2)),
			TapSLog[int]("After multiplication"),
		)

		pipeline := F.Pipe2(
			step1,
			Chain(validatePositive),
			TapSLog[int]("After validation"),
		)

		res := pipeline(t.Context())()

		assert.Equal(t, 10, getOrZero(res))

		logOutput := buf.String()
		assert.Contains(t, logOutput, "Input received")
		assert.Contains(t, logOutput, "value=5")
		assert.Contains(t, logOutput, "After multiplication")
		assert.Contains(t, logOutput, "value=10")
		assert.Contains(t, logOutput, "After validation")
		assert.Contains(t, logOutput, "value=10")
	})

	t.Run("error propagation with logging", func(t *testing.T) {
		buf.Reset()

		validatePositive := func(n int) ReaderIOResult[int] {
			if n > 0 {
				return Of(n)
			}
			return Left[int](errors.New("number must be positive"))
		}

		step1 := F.Pipe3(
			Of(-5),
			TapSLog[int]("Input received"),
			Map(N.Mul(2)),
			TapSLog[int]("After multiplication"),
		)

		pipeline := F.Pipe2(
			step1,
			Chain(validatePositive),
			TapSLog[int]("After validation"),
		)

		res := pipeline(t.Context())()

		assert.True(t, isLeft(res))

		logOutput := buf.String()
		// First two taps should log success
		assert.Contains(t, logOutput, "Input received")
		assert.Contains(t, logOutput, "value=-5")
		assert.Contains(t, logOutput, "After multiplication")
		assert.Contains(t, logOutput, "value=-10")
		// Last tap should log error
		assert.Contains(t, logOutput, "After validation")
		assert.Contains(t, logOutput, "number must be positive")
	})
}

// Helper functions for tests

func getOrZero(res Result[int]) int {
	val, err := result.Unwrap(res)
	if err == nil {
		return val
	}
	return 0
}

func isRight[A any](res Result[A]) bool {
	return result.IsRight(res)
}

func isLeft[A any](res Result[A]) bool {
	return result.IsLeft(res)
}
