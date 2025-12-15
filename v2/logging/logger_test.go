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
	"log"
	"strings"
	"testing"

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
