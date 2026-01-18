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

package readerresult

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/IBM/fp-go/v2/logging"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestSLogLogsSuccessValue tests that SLog logs successful Result values
func TestSLogLogsSuccessValue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()

	// Create a Result and log it
	res1 := result.Of(42)
	logged := SLog[int]("Result value")(res1)(ctx)

	assert.Equal(t, result.Of(42), logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Result value")
	assert.Contains(t, logOutput, "value=42")
}

// TestSLogLogsErrorValue tests that SLog logs error Result values
func TestSLogLogsErrorValue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()
	testErr := errors.New("test error")

	// Create an error Result and log it
	res1 := result.Left[int](testErr)
	logged := SLog[int]("Result value")(res1)(ctx)

	assert.Equal(t, res1, logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Result value")
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "test error")
}

// TestSLogInPipeline tests SLog in a functional pipeline
func TestSLogInPipeline(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()

	// SLog takes a Result[A] and returns ReaderResult[A]
	// So we need to start with a Result, apply SLog, then execute with context
	res1 := result.Of(10)
	logged := SLog[int]("Initial value")(res1)(ctx)

	assert.Equal(t, result.Of(10), logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Initial value")
	assert.Contains(t, logOutput, "value=10")
}

// TestSLogWithContextLogger tests SLog using logger from context
func TestSLogWithContextLogger(t *testing.T) {
	var buf bytes.Buffer
	contextLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx := logging.WithLogger(contextLogger)(t.Context())

	res1 := result.Of("test value")
	logged := SLog[string]("Context logger test")(res1)(ctx)

	assert.Equal(t, result.Of("test value"), logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Context logger test")
	assert.Contains(t, logOutput, `value="test value"`)
}

// TestSLogDisabled tests that SLog respects logger level
func TestSLogDisabled(t *testing.T) {
	var buf bytes.Buffer
	// Create logger with level that disables info logs
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError, // Only log errors
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()

	res1 := result.Of(42)
	logged := SLog[int]("This should not be logged")(res1)(ctx)

	assert.Equal(t, result.Of(42), logged)

	// Should have no logs since level is ERROR
	logOutput := buf.String()
	assert.Empty(t, logOutput, "Should have no logs when logging is disabled")
}

// TestSLogWithStruct tests SLog with structured data
func TestSLogWithStruct(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	type User struct {
		ID   int
		Name string
	}

	ctx := t.Context()
	user := User{ID: 123, Name: "Alice"}

	res1 := result.Of(user)
	logged := SLog[User]("User data")(res1)(ctx)

	assert.Equal(t, result.Of(user), logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "User data")
	assert.Contains(t, logOutput, "ID:123")
	assert.Contains(t, logOutput, "Name:Alice")
}

// TestSLogWithCallbackCustomLevel tests SLogWithCallback with custom log level
func TestSLogWithCallbackCustomLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	customCallback := func(ctx context.Context) *slog.Logger {
		return logger
	}

	ctx := t.Context()

	// Create a Result and log it with custom callback
	res1 := result.Of(42)
	logged := SLogWithCallback[int](slog.LevelDebug, customCallback, "Debug result")(res1)(ctx)

	assert.Equal(t, result.Of(42), logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Debug result")
	assert.Contains(t, logOutput, "value=42")
	assert.Contains(t, logOutput, "level=DEBUG")
}

// TestSLogWithCallbackLogsError tests SLogWithCallback logs errors
func TestSLogWithCallbackLogsError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	customCallback := func(ctx context.Context) *slog.Logger {
		return logger
	}

	ctx := t.Context()
	testErr := errors.New("warning error")

	// Create an error Result and log it with custom callback
	res1 := result.Left[int](testErr)
	logged := SLogWithCallback[int](slog.LevelWarn, customCallback, "Warning result")(res1)(ctx)

	assert.Equal(t, res1, logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Warning result")
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "warning error")
	assert.Contains(t, logOutput, "level=WARN")
}

// TestSLogChainedOperations tests SLog in chained operations
func TestSLogChainedOperations(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()

	// First log step 1
	res1 := result.Of(5)
	logged1 := SLog[int]("Step 1")(res1)(ctx)

	// Then log step 2 with doubled value
	res2 := result.Map(N.Mul(2))(logged1)
	logged2 := SLog[int]("Step 2")(res2)(ctx)

	assert.Equal(t, result.Of(10), logged2)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Step 1")
	assert.Contains(t, logOutput, "value=5")
	assert.Contains(t, logOutput, "Step 2")
	assert.Contains(t, logOutput, "value=10")
}

// TestSLogPreservesError tests that SLog preserves error through the pipeline
func TestSLogPreservesError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()
	testErr := errors.New("original error")

	res1 := result.Left[int](testErr)
	logged := SLog[int]("Logging error")(res1)(ctx)

	// Apply map to verify error is preserved
	res2 := result.Map(N.Mul(2))(logged)

	assert.Equal(t, res1, res2)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Logging error")
	assert.Contains(t, logOutput, "original error")
}

// TestSLogMultipleValues tests logging multiple different values
func TestSLogMultipleValues(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()

	// Test with different types
	intRes := SLog[int]("Integer")(result.Of(42))(ctx)
	assert.Equal(t, result.Of(42), intRes)

	strRes := SLog[string]("String")(result.Of("hello"))(ctx)
	assert.Equal(t, result.Of("hello"), strRes)

	boolRes := SLog[bool]("Boolean")(result.Of(true))(ctx)
	assert.Equal(t, result.Of(true), boolRes)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Integer")
	assert.Contains(t, logOutput, "value=42")
	assert.Contains(t, logOutput, "String")
	assert.Contains(t, logOutput, "value=hello")
	assert.Contains(t, logOutput, "Boolean")
	assert.Contains(t, logOutput, "value=true")
}
