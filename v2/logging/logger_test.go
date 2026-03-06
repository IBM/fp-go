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

package logging

import (
	"bytes"
	"context"
	"log"
	"log/slog"
	"strings"
	"testing"

	"github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
)

// TestLoggingCallbacks_NoLoggers tests the case when no loggers are provided.
// It should return two callbacks using the default logger.
func TestLoggingCallbacks_NoLoggers(t *testing.T) {
	infoLog, errLog := LoggingCallbacks()

	if infoLog == nil {
		t.Error("Expected infoLog to be non-nil")
	}
	if errLog == nil {
		t.Error("Expected errLog to be non-nil")
	}

	// Verify both callbacks work
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	infoLog("test info: %s", "message")
	if !strings.Contains(buf.String(), "test info: message") {
		t.Errorf("Expected log to contain 'test info: message', got: %s", buf.String())
	}

	buf.Reset()
	errLog("test error: %s", "message")
	if !strings.Contains(buf.String(), "test error: message") {
		t.Errorf("Expected log to contain 'test error: message', got: %s", buf.String())
	}
}

// TestLoggingCallbacks_OneLogger tests the case when one logger is provided.
// Both callbacks should use the same logger.
func TestLoggingCallbacks_OneLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "TEST: ", 0)

	infoLog, errLog := LoggingCallbacks(logger)

	if infoLog == nil {
		t.Error("Expected infoLog to be non-nil")
	}
	if errLog == nil {
		t.Error("Expected errLog to be non-nil")
	}

	// Test info callback
	infoLog("info message: %d", 42)
	output := buf.String()
	if !strings.Contains(output, "TEST: info message: 42") {
		t.Errorf("Expected log to contain 'TEST: info message: 42', got: %s", output)
	}

	// Test error callback uses same logger
	buf.Reset()
	errLog("error message: %s", "failed")
	output = buf.String()
	if !strings.Contains(output, "TEST: error message: failed") {
		t.Errorf("Expected log to contain 'TEST: error message: failed', got: %s", output)
	}
}

// TestLoggingCallbacks_TwoLoggers tests the case when two loggers are provided.
// First callback should use first logger, second callback should use second logger.
func TestLoggingCallbacks_TwoLoggers(t *testing.T) {
	var infoBuf, errBuf bytes.Buffer
	infoLogger := log.New(&infoBuf, "INFO: ", 0)
	errorLogger := log.New(&errBuf, "ERROR: ", 0)

	infoLog, errLog := LoggingCallbacks(infoLogger, errorLogger)

	if infoLog == nil {
		t.Error("Expected infoLog to be non-nil")
	}
	if errLog == nil {
		t.Error("Expected errLog to be non-nil")
	}

	// Test info callback uses first logger
	infoLog("success: %s", "operation completed")
	infoOutput := infoBuf.String()
	if !strings.Contains(infoOutput, "INFO: success: operation completed") {
		t.Errorf("Expected info log to contain 'INFO: success: operation completed', got: %s", infoOutput)
	}
	if errBuf.Len() != 0 {
		t.Errorf("Expected error buffer to be empty, got: %s", errBuf.String())
	}

	// Test error callback uses second logger
	errLog("failure: %s", "operation failed")
	errOutput := errBuf.String()
	if !strings.Contains(errOutput, "ERROR: failure: operation failed") {
		t.Errorf("Expected error log to contain 'ERROR: failure: operation failed', got: %s", errOutput)
	}
}

// TestLoggingCallbacks_MultipleLoggers tests the case when more than two loggers are provided.
// Should use first two loggers and ignore the rest.
func TestLoggingCallbacks_MultipleLoggers(t *testing.T) {
	var buf1, buf2, buf3 bytes.Buffer
	logger1 := log.New(&buf1, "LOG1: ", 0)
	logger2 := log.New(&buf2, "LOG2: ", 0)
	logger3 := log.New(&buf3, "LOG3: ", 0)

	infoLog, errLog := LoggingCallbacks(logger1, logger2, logger3)

	if infoLog == nil {
		t.Error("Expected infoLog to be non-nil")
	}
	if errLog == nil {
		t.Error("Expected errLog to be non-nil")
	}

	// Test that first logger is used for info
	infoLog("message 1")
	if !strings.Contains(buf1.String(), "LOG1: message 1") {
		t.Errorf("Expected first logger to be used, got: %s", buf1.String())
	}

	// Test that second logger is used for error
	errLog("message 2")
	if !strings.Contains(buf2.String(), "LOG2: message 2") {
		t.Errorf("Expected second logger to be used, got: %s", buf2.String())
	}

	// Test that third logger is not used
	if buf3.Len() != 0 {
		t.Errorf("Expected third logger to not be used, got: %s", buf3.String())
	}
}

// TestLoggingCallbacks_FormattingWithMultipleArgs tests that formatting works correctly
// with multiple arguments.
func TestLoggingCallbacks_FormattingWithMultipleArgs(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	infoLog, _ := LoggingCallbacks(logger)

	infoLog("test %s %d %v", "string", 123, true)
	output := buf.String()
	if !strings.Contains(output, "test string 123 true") {
		t.Errorf("Expected formatted output 'test string 123 true', got: %s", output)
	}
}

// TestLoggingCallbacks_NoFormatting tests logging without format specifiers.
func TestLoggingCallbacks_NoFormatting(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "PREFIX: ", 0)

	infoLog, _ := LoggingCallbacks(logger)

	infoLog("simple message")
	output := buf.String()
	if !strings.Contains(output, "PREFIX: simple message") {
		t.Errorf("Expected 'PREFIX: simple message', got: %s", output)
	}
}

// TestLoggingCallbacks_EmptyMessage tests logging with empty message.
func TestLoggingCallbacks_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	infoLog, _ := LoggingCallbacks(logger)

	infoLog("")
	output := buf.String()
	// Should still produce output (newline at minimum)
	if S.IsEmpty(output) {
		t.Error("Expected some output even with empty message")
	}
}

// TestLoggingCallbacks_NilLogger tests behavior when nil logger is passed.
// This tests edge case handling.
func TestLoggingCallbacks_NilLogger(t *testing.T) {
	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LoggingCallbacks panicked with nil logger: %v", r)
		}
	}()

	infoLog, errLog := LoggingCallbacks(nil)

	// The callbacks should still be created
	if infoLog == nil {
		t.Error("Expected infoLog to be non-nil even with nil logger")
	}
	if errLog == nil {
		t.Error("Expected errLog to be non-nil even with nil logger")
	}
}

// TestLoggingCallbacks_ConsecutiveCalls tests that callbacks can be called multiple times.
func TestLoggingCallbacks_ConsecutiveCalls(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	infoLog, errLog := LoggingCallbacks(logger)

	// Multiple calls to info
	infoLog("call 1")
	infoLog("call 2")
	infoLog("call 3")

	output := buf.String()
	if !strings.Contains(output, "call 1") || !strings.Contains(output, "call 2") || !strings.Contains(output, "call 3") {
		t.Errorf("Expected all three calls to be logged, got: %s", output)
	}

	buf.Reset()

	// Multiple calls to error
	errLog("error 1")
	errLog("error 2")

	output = buf.String()
	if !strings.Contains(output, "error 1") || !strings.Contains(output, "error 2") {
		t.Errorf("Expected both error calls to be logged, got: %s", output)
	}
}

// BenchmarkLoggingCallbacks_NoLoggers benchmarks the no-logger case.
func BenchmarkLoggingCallbacks_NoLoggers(b *testing.B) {
	for b.Loop() {
		LoggingCallbacks()
	}
}

// BenchmarkLoggingCallbacks_OneLogger benchmarks the single-logger case.
func BenchmarkLoggingCallbacks_OneLogger(b *testing.B) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	b.ResetTimer()
	for b.Loop() {
		LoggingCallbacks(logger)
	}
}

// BenchmarkLoggingCallbacks_TwoLoggers benchmarks the two-logger case.
func BenchmarkLoggingCallbacks_TwoLoggers(b *testing.B) {
	var buf1, buf2 bytes.Buffer
	logger1 := log.New(&buf1, "", 0)
	logger2 := log.New(&buf2, "", 0)

	b.ResetTimer()
	for b.Loop() {
		LoggingCallbacks(logger1, logger2)
	}
}

// BenchmarkLoggingCallbacks_Logging benchmarks actual logging operations.
func BenchmarkLoggingCallbacks_Logging(b *testing.B) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	infoLog, _ := LoggingCallbacks(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		infoLog("benchmark message %d", i)
	}
}

// TestSetLogger_Success tests setting a new global logger and verifying it returns the old one.
func TestSetLogger_Success(t *testing.T) {
	// Save original logger to restore later
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	// Create a new logger
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	newLogger := slog.New(handler)

	// Set the new logger
	oldLogger := SetLogger(newLogger)

	// Verify old logger was returned
	if oldLogger == nil {
		t.Error("Expected SetLogger to return the previous logger")
	}

	// Verify new logger is now active
	currentLogger := GetLogger()
	if currentLogger != newLogger {
		t.Error("Expected GetLogger to return the newly set logger")
	}
}

// TestSetLogger_Multiple tests setting logger multiple times.
func TestSetLogger_Multiple(t *testing.T) {
	// Save original logger to restore later
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	// Create three loggers
	logger1 := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	logger2 := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	logger3 := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

	// Set first logger
	old1 := SetLogger(logger1)
	if GetLogger() != logger1 {
		t.Error("Expected logger1 to be active")
	}

	// Set second logger
	old2 := SetLogger(logger2)
	if old2 != logger1 {
		t.Error("Expected SetLogger to return logger1")
	}
	if GetLogger() != logger2 {
		t.Error("Expected logger2 to be active")
	}

	// Set third logger
	old3 := SetLogger(logger3)
	if old3 != logger2 {
		t.Error("Expected SetLogger to return logger2")
	}
	if GetLogger() != logger3 {
		t.Error("Expected logger3 to be active")
	}

	// Restore to original
	restored := SetLogger(old1)
	if restored != logger3 {
		t.Error("Expected SetLogger to return logger3")
	}
}

// TestGetLogger_Default tests that GetLogger returns a valid logger by default.
func TestGetLogger_Default(t *testing.T) {
	logger := GetLogger()

	if logger == nil {
		t.Error("Expected GetLogger to return a non-nil logger")
	}

	// Verify it's usable
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	testLogger := slog.New(handler)

	oldLogger := SetLogger(testLogger)
	defer SetLogger(oldLogger)

	GetLogger().Info("test message")
	if !strings.Contains(buf.String(), "test message") {
		t.Errorf("Expected logger to log message, got: %s", buf.String())
	}
}

// TestGetLogger_AfterSet tests that GetLogger returns the logger set by SetLogger.
func TestGetLogger_AfterSet(t *testing.T) {
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	customLogger := slog.New(handler)

	SetLogger(customLogger)

	retrievedLogger := GetLogger()
	if retrievedLogger != customLogger {
		t.Error("Expected GetLogger to return the custom logger")
	}

	// Verify it's the same instance by logging
	retrievedLogger.Info("test")
	if !strings.Contains(buf.String(), "test") {
		t.Error("Expected retrieved logger to be the same instance")
	}
}

// TestGetLoggerFromContext_WithLogger tests retrieving a logger from context.
func TestGetLoggerFromContext_WithLogger(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	contextLogger := slog.New(handler)

	// Create context with logger using WithLogger
	ctx := context.Background()
	kleisli := WithLogger(contextLogger)
	result := kleisli(ctx)
	ctxWithLogger := pair.Second(result)

	// Retrieve logger from context
	retrievedLogger := GetLoggerFromContext(ctxWithLogger)

	if retrievedLogger != contextLogger {
		t.Error("Expected to retrieve the context logger")
	}

	// Verify it's the same instance by logging
	retrievedLogger.Info("context test")
	if !strings.Contains(buf.String(), "context test") {
		t.Error("Expected retrieved logger to be the same instance")
	}
}

// TestGetLoggerFromContext_WithoutLogger tests that it returns global logger when context has no logger.
func TestGetLoggerFromContext_WithoutLogger(t *testing.T) {
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	globalLogger := slog.New(handler)
	SetLogger(globalLogger)

	// Create context without logger
	ctx := context.Background()

	// Should return global logger
	retrievedLogger := GetLoggerFromContext(ctx)

	if retrievedLogger != globalLogger {
		t.Error("Expected to retrieve the global logger when context has no logger")
	}

	// Verify it's the same instance
	retrievedLogger.Info("global test")
	if !strings.Contains(buf.String(), "global test") {
		t.Error("Expected retrieved logger to be the global logger")
	}
}

// TestGetLoggerFromContext_NilContext tests behavior with nil context value.
func TestGetLoggerFromContext_NilContext(t *testing.T) {
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	globalLogger := slog.New(handler)
	SetLogger(globalLogger)

	// Create context with wrong type value
	ctx := context.WithValue(context.Background(), loggerInContextKey, "not a logger")

	// Should return global logger when type assertion fails
	retrievedLogger := GetLoggerFromContext(ctx)

	if retrievedLogger != globalLogger {
		t.Error("Expected to retrieve the global logger when context value is wrong type")
	}
}

// TestWithLogger_CreatesContextWithLogger tests that WithLogger adds logger to context.
func TestWithLogger_CreatesContextWithLogger(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	testLogger := slog.New(handler)

	// Create Kleisli arrow
	kleisli := WithLogger(testLogger)

	// Apply to context
	ctx := context.Background()
	result := kleisli(ctx)

	// Verify result is a ContextCancel pair
	cancelFunc := pair.First(result)
	newCtx := pair.Second(result)

	if cancelFunc == nil {
		t.Error("Expected cancel function to be non-nil")
	}

	if newCtx == nil {
		t.Error("Expected new context to be non-nil")
	}

	// Verify logger is in context
	retrievedLogger := GetLoggerFromContext(newCtx)
	if retrievedLogger != testLogger {
		t.Error("Expected logger to be in the new context")
	}
}

// TestWithLogger_CancelFuncIsNoop tests that the cancel function is a no-op.
func TestWithLogger_CancelFuncIsNoop(t *testing.T) {
	testLogger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	kleisli := WithLogger(testLogger)

	ctx := context.Background()
	result := kleisli(ctx)
	cancelFunc := pair.First(result)

	// Calling cancel should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Cancel function panicked: %v", r)
		}
	}()

	cancelFunc()
}

// TestWithLogger_PreservesOriginalContext tests that original context is not modified.
func TestWithLogger_PreservesOriginalContext(t *testing.T) {
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	globalLogger := slog.New(handler)
	SetLogger(globalLogger)

	testLogger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	kleisli := WithLogger(testLogger)

	// Original context without logger
	originalCtx := context.Background()

	// Apply transformation
	result := kleisli(originalCtx)
	newCtx := pair.Second(result)

	// Original context should still return global logger
	originalCtxLogger := GetLoggerFromContext(originalCtx)
	if originalCtxLogger != globalLogger {
		t.Error("Expected original context to still use global logger")
	}

	// New context should have the test logger
	newCtxLogger := GetLoggerFromContext(newCtx)
	if newCtxLogger != testLogger {
		t.Error("Expected new context to have the test logger")
	}
}

// TestWithLogger_Composition tests composing multiple WithLogger calls.
func TestWithLogger_Composition(t *testing.T) {
	logger1 := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	logger2 := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

	kleisli1 := WithLogger(logger1)
	kleisli2 := WithLogger(logger2)

	ctx := context.Background()

	// Apply first transformation
	result1 := kleisli1(ctx)
	ctx1 := pair.Second(result1)

	// Verify first logger
	if GetLoggerFromContext(ctx1) != logger1 {
		t.Error("Expected first logger in context after first transformation")
	}

	// Apply second transformation (should override)
	result2 := kleisli2(ctx1)
	ctx2 := pair.Second(result2)

	// Verify second logger (should override first)
	if GetLoggerFromContext(ctx2) != logger2 {
		t.Error("Expected second logger to override first logger")
	}
}

// BenchmarkSetLogger benchmarks setting the global logger.
func BenchmarkSetLogger(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetLogger(logger)
	}
}

// BenchmarkGetLogger benchmarks getting the global logger.
func BenchmarkGetLogger(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetLogger()
	}
}

// BenchmarkGetLoggerFromContext_WithLogger benchmarks retrieving logger from context.
func BenchmarkGetLoggerFromContext_WithLogger(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	kleisli := WithLogger(logger)
	ctx := pair.Second(kleisli(context.Background()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetLoggerFromContext(ctx)
	}
}

// BenchmarkGetLoggerFromContext_WithoutLogger benchmarks retrieving global logger from context.
func BenchmarkGetLoggerFromContext_WithoutLogger(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetLoggerFromContext(ctx)
	}
}

// BenchmarkWithLogger benchmarks creating context with logger.
func BenchmarkWithLogger(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	kleisli := WithLogger(logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kleisli(ctx)
	}
}
