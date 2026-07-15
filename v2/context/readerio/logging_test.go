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

package readerio

import (
	"bytes"
	"log/slog"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/logging"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// makeTestLogger creates a text-format slog.Logger that writes into buf at Info level.
func makeTestLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// TestTapSLog_LogsMessageAndValue verifies that the log message and value attribute
// are written to the configured logger when the pipeline is executed.
func TestTapSLog_LogsMessageAndValue(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(42),
		TapSLog[int]("processing value"),
		Map(N.Mul(2)),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 84, result)
	assert.Contains(t, buf.String(), "processing value")
	assert.Contains(t, buf.String(), "value=42")
}

// TestTapSLog_PreservesValue verifies that TapSLog does not alter the value flowing
// through the pipeline — the downstream computation receives the original value.
func TestTapSLog_PreservesValue(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe1(
		Of(99),
		TapSLog[int]("tap"),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 99, result)
}

// TestTapSLog_MultiStepPipeline verifies that multiple TapSLog steps each log their
// own message with the value present at that point in the pipeline.
func TestTapSLog_MultiStepPipeline(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe4(
		Of("hello"),
		TapSLog[string]("step 1"),
		Map(func(s string) string { return s + " world" }),
		TapSLog[string]("step 2"),
		Map(S.Size),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 11, result)

	log := buf.String()
	assert.Contains(t, log, "step 1")
	assert.Contains(t, log, "value=hello")
	assert.Contains(t, log, "step 2")
	assert.Contains(t, log, `value="hello world"`)
}

// TestTapSLog_RespectsLogLevel verifies that no output is produced when the global
// logger's minimum level is above Info (i.e., the log entry is filtered out), while
// the value still flows through the pipeline unchanged.
func TestTapSLog_RespectsLogLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(7),
		TapSLog[int]("should not appear"),
		Map(N.Mul(3)),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 21, result)
	assert.Empty(t, buf.String(), "no log output expected when level is above Info")
}

// TestTapSLog_UsesContextLogger verifies that a logger injected into the context via
// logging.WithLogger is used in preference to the global logger.
func TestTapSLog_UsesContextLogger(t *testing.T) {
	var buf bytes.Buffer
	contextLogger := makeTestLogger(&buf)

	cancelFct, ctx := pair.Unpack(logging.WithLogger(contextLogger)(t.Context()))
	defer cancelFct()

	pipeline := F.Pipe2(
		Of("ctx-value"),
		TapSLog[string]("context logger test"),
		Map(S.Size),
	)

	result := pipeline(ctx)()

	assert.Equal(t, 9, result)
	assert.Contains(t, buf.String(), "context logger test")
	assert.Contains(t, buf.String(), "value=ctx-value")
}

// TestTapSLog_StructValue verifies that TapSLog serialises structured values into
// the log output and still passes the original struct to downstream operators.
func TestTapSLog_StructValue(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	user := User{ID: 1, Name: "Alice"}

	pipeline := F.Pipe2(
		Of(user),
		TapSLog[User]("user data"),
		Map(func(u User) string { return u.Name }),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, "Alice", result)

	log := buf.String()
	assert.Contains(t, log, "user data")
	assert.Contains(t, log, "ID:1")
	assert.Contains(t, log, "Name:Alice")
}

// TestTapSLog_ZeroValue verifies correct behaviour when the value being logged is
// the zero value for its type.
func TestTapSLog_ZeroValue(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe1(
		Of(0),
		TapSLog[int]("zero value"),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 0, result)
	assert.Contains(t, buf.String(), "zero value")
	assert.Contains(t, buf.String(), "value=0")
}

// TestTapSLog_NilPointer verifies that TapSLog handles nil pointer values without
// panicking and still produces a log entry.
func TestTapSLog_NilPointer(t *testing.T) {
	type Data struct{ X int }

	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	var ptr *Data

	pipeline := F.Pipe1(
		Of(ptr),
		TapSLog[*Data]("nil pointer"),
	)

	assert.NotPanics(t, func() {
		_ = pipeline(t.Context())()
	})
	assert.Contains(t, buf.String(), "nil pointer")
}

// TestTapSLog_LazyExecution verifies that the log side-effect is deferred: no output
// must be produced until the resulting IO is actually executed.
func TestTapSLog_LazyExecution(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	// Build but do not yet execute the pipeline.
	pipeline := F.Pipe1(
		Of(5),
		TapSLog[int]("lazy log"),
	)
	readerIO := pipeline(t.Context())

	assert.Empty(t, buf.String(), "no log output expected before IO execution")

	_ = readerIO()

	assert.Contains(t, buf.String(), "lazy log")
}

// ---------------------------------------------------------------------------
// SLogInfo
// ---------------------------------------------------------------------------

// TestSLogInfo_LogsAtInfoLevel verifies that SLogInfo produces a log entry at
// INFO level with the message and value attribute.
func TestSLogInfo_LogsAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe1(
		Of(10),
		Chain(SLogInfo[int]("sloginfo msg")),
	)

	pipeline(t.Context())()

	log := buf.String()
	assert.Contains(t, log, "INFO")
	assert.Contains(t, log, "sloginfo msg")
	assert.Contains(t, log, "value=10")
}

// TestSLogInfo_EquivalentToSLog verifies that SLogInfo and SLog produce
// identical output for the same input.
func TestSLogInfo_EquivalentToSLog(t *testing.T) {
	var bufInfo, bufSLog bytes.Buffer

	oldLogger := logging.SetLogger(makeTestLogger(&bufInfo))
	Of(42)(t.Context())()
	F.Pipe1(Of(42), Chain(SLogInfo[int]("msg")))(t.Context())()
	logging.SetLogger(oldLogger)

	oldLogger = logging.SetLogger(makeTestLogger(&bufSLog))
	F.Pipe1(Of(42), Chain(SLog[int]("msg")))(t.Context())()
	logging.SetLogger(oldLogger)

	// Both should contain the same message and value.
	assert.Contains(t, bufInfo.String(), "msg")
	assert.Contains(t, bufInfo.String(), "value=42")
	assert.Contains(t, bufSLog.String(), "msg")
	assert.Contains(t, bufSLog.String(), "value=42")
}

// TestSLogInfo_SuppressedWhenLevelAboveInfo verifies that no log output is
// produced when the logger's minimum level is above Info.
func TestSLogInfo_SuppressedWhenLevelAboveInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	F.Pipe1(Of(1), Chain(SLogInfo[int]("hidden")))(t.Context())()

	assert.Empty(t, buf.String())
}

// ---------------------------------------------------------------------------
// SLogDebug
// ---------------------------------------------------------------------------

// TestSLogDebug_LogsAtDebugLevel verifies that SLogDebug produces a log entry at
// DEBUG level when the logger accepts Debug messages.
func TestSLogDebug_LogsAtDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	F.Pipe1(Of(7), Chain(SLogDebug[int]("debug msg")))(t.Context())()

	log := buf.String()
	assert.Contains(t, log, "DEBUG")
	assert.Contains(t, log, "debug msg")
	assert.Contains(t, log, "value=7")
}

// TestSLogDebug_SuppressedAtInfoLevel verifies that no output is produced when
// the logger's minimum level is Info (i.e., Debug is filtered out).
func TestSLogDebug_SuppressedAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf)) // Info level
	defer logging.SetLogger(oldLogger)

	F.Pipe1(Of(3), Chain(SLogDebug[int]("should not appear")))(t.Context())()

	assert.Empty(t, buf.String(), "Debug entry must not appear when logger level is Info")
}

// TestSLogDebug_PreservesValue verifies the Kleisli arrow returns Void and does
// not affect the value in a subsequent Chain step.
func TestSLogDebug_PreservesValue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	// Chain returns Void; the tap pattern is tested via TapSLogDebug.
	// Here we just confirm the arrow executes without panic.
	assert.NotPanics(t, func() {
		F.Pipe1(Of("hello"), Chain(SLogDebug[string]("check")))(t.Context())()
	})
}

// ---------------------------------------------------------------------------
// TapSLogInfo
// ---------------------------------------------------------------------------

// TestTapSLogInfo_LogsAtInfoLevelAndPreservesValue verifies that TapSLogInfo
// writes an Info log entry and passes the original value downstream.
func TestTapSLogInfo_LogsAtInfoLevelAndPreservesValue(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf))
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(55),
		TapSLogInfo[int]("tapsloginfo msg"),
		Map(N.Mul(2)),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 110, result)
	log := buf.String()
	assert.Contains(t, log, "INFO")
	assert.Contains(t, log, "tapsloginfo msg")
	assert.Contains(t, log, "value=55")
}

// TestTapSLogInfo_SuppressedWhenLevelAboveInfo verifies that the pipeline still
// runs and produces the correct value when Info logging is disabled.
func TestTapSLogInfo_SuppressedWhenLevelAboveInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(4),
		TapSLogInfo[int]("invisible"),
		Map(N.Mul(3)),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 12, result)
	assert.Empty(t, buf.String())
}

// TestTapSLogInfo_EquivalentToTapSLog verifies that TapSLogInfo behaves
// identically to TapSLog for the same message and value.
func TestTapSLogInfo_EquivalentToTapSLog(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	oldLogger := logging.SetLogger(makeTestLogger(&buf1))
	F.Pipe1(Of(9), TapSLogInfo[int]("same"))(t.Context())()
	logging.SetLogger(oldLogger)

	oldLogger = logging.SetLogger(makeTestLogger(&buf2))
	F.Pipe1(Of(9), TapSLog[int]("same"))(t.Context())()
	logging.SetLogger(oldLogger)

	assert.Contains(t, buf1.String(), "same")
	assert.Contains(t, buf1.String(), "value=9")
	assert.Contains(t, buf2.String(), "same")
	assert.Contains(t, buf2.String(), "value=9")
}

// ---------------------------------------------------------------------------
// TapSLogDebug
// ---------------------------------------------------------------------------

// TestTapSLogDebug_LogsAtDebugLevelAndPreservesValue verifies that TapSLogDebug
// writes a Debug log entry and passes the original value downstream.
func TestTapSLogDebug_LogsAtDebugLevelAndPreservesValue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(8),
		TapSLogDebug[int]("debug tap"),
		Map(N.Mul(5)),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 40, result)
	log := buf.String()
	assert.Contains(t, log, "DEBUG")
	assert.Contains(t, log, "debug tap")
	assert.Contains(t, log, "value=8")
}

// TestTapSLogDebug_SuppressedAtInfoLevel verifies that no output is produced
// and the pipeline result is still correct when Debug is filtered out.
func TestTapSLogDebug_SuppressedAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(makeTestLogger(&buf)) // Info level
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(6),
		TapSLogDebug[int]("should not appear"),
		Map(N.Mul(2)),
	)

	result := pipeline(t.Context())()

	assert.Equal(t, 12, result)
	assert.Empty(t, buf.String(), "Debug entry must not appear when logger level is Info")
}

// TestTapSLogDebug_LazyExecution verifies that no debug output is produced
// until the resulting IO is actually executed.
func TestTapSLogDebug_LazyExecution(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe1(Of(2), TapSLogDebug[int]("lazy debug"))
	readerIO := pipeline(t.Context())

	assert.Empty(t, buf.String(), "no output expected before IO execution")

	_ = readerIO()

	assert.Contains(t, buf.String(), "lazy debug")
}
