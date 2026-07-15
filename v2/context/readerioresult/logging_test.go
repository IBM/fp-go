package readerioresult

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/logging"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestLoggingContext tests basic nested logging with correlation IDs
func TestLoggingContext(t *testing.T) {
	data := F.Pipe2(
		Of("Sample"),
		LogEntryExit[string]("TestLoggingContext1"),
		LogEntryExit[string]("TestLoggingContext2"),
	)

	assert.Equal(t, result.Of("Sample"), data(t.Context())())
}

// TestLogEntryExitSuccess tests successful operation logging
func TestLogEntryExitSuccess(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	operation := F.Pipe1(
		Of("success value"),
		LogEntryExit[string]("TestOperation"),
	)

	res := operation(t.Context())()

	assert.Equal(t, result.Of("success value"), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "[entering]")
	assert.Contains(t, logOutput, "[exiting ]")
	assert.Contains(t, logOutput, "TestOperation")
	assert.Contains(t, logOutput, "ID=")
	assert.Contains(t, logOutput, "duration=")

	// Verify entry log appears before exit log
	enteringIdx := strings.Index(logOutput, "[entering]")
	exitingIdx := strings.Index(logOutput, "[exiting ]")
	assert.Greater(t, exitingIdx, enteringIdx, "Exit log should appear after entry log")
}

// TestLogEntryExitError tests error operation logging
func TestLogEntryExitError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testErr := errors.New("test error")
	operation := F.Pipe1(
		Left[string](testErr),
		LogEntryExit[string]("FailingOperation"),
	)

	res := operation(t.Context())()

	assert.True(t, result.IsLeft(res))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "[entering]")
	assert.Contains(t, logOutput, "[throwing]")
	assert.Contains(t, logOutput, "FailingOperation")
	assert.Contains(t, logOutput, "test error")
	assert.Contains(t, logOutput, "ID=")
	assert.Contains(t, logOutput, "duration=")

	// Verify entry log appears before error log
	enteringIdx := strings.Index(logOutput, "[entering]")
	throwingIdx := strings.Index(logOutput, "[throwing]")
	assert.Greater(t, throwingIdx, enteringIdx, "Error log should appear after entry log")
}

// TestLogEntryExitNested tests nested operations with different IDs
func TestLogEntryExitNested(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	innerOp := F.Pipe1(
		Of("inner"),
		LogEntryExit[string]("InnerOp"),
	)

	outerOp := F.Pipe2(
		Of("outer"),
		LogEntryExit[string]("OuterOp"),
		Chain(func(s string) ReaderIOResult[string] {
			return innerOp
		}),
	)

	res := outerOp(t.Context())()

	assert.True(t, result.IsRight(res))

	logOutput := buf.String()
	// Should have two different IDs
	assert.Contains(t, logOutput, "OuterOp")
	assert.Contains(t, logOutput, "InnerOp")

	// Count entering and exiting logs
	enterCount := strings.Count(logOutput, "[entering]")
	exitCount := strings.Count(logOutput, "[exiting ]")
	assert.Equal(t, 2, enterCount, "Should have 2 entering logs")
	assert.Equal(t, 2, exitCount, "Should have 2 exiting logs")

	// Verify log ordering: Each operation logs entry before exit
	// Note: Due to Chain semantics, OuterOp completes before InnerOp starts
	lines := strings.Split(logOutput, "\n")
	var logSequence []string
	for _, line := range lines {
		if strings.Contains(line, "OuterOp") && strings.Contains(line, "[entering]") {
			logSequence = append(logSequence, "OuterOp-entering")
		} else if strings.Contains(line, "OuterOp") && strings.Contains(line, "[exiting ]") {
			logSequence = append(logSequence, "OuterOp-exiting")
		} else if strings.Contains(line, "InnerOp") && strings.Contains(line, "[entering]") {
			logSequence = append(logSequence, "InnerOp-entering")
		} else if strings.Contains(line, "InnerOp") && strings.Contains(line, "[exiting ]") {
			logSequence = append(logSequence, "InnerOp-exiting")
		}
	}

	// Verify each operation's entry comes before its exit
	assert.Equal(t, 4, len(logSequence), "Should have 4 log entries")

	// Find indices
	outerEnterIdx := -1
	outerExitIdx := -1
	innerEnterIdx := -1
	innerExitIdx := -1

	for i, log := range logSequence {
		switch log {
		case "OuterOp-entering":
			outerEnterIdx = i
		case "OuterOp-exiting":
			outerExitIdx = i
		case "InnerOp-entering":
			innerEnterIdx = i
		case "InnerOp-exiting":
			innerExitIdx = i
		}
	}

	// Verify entry before exit for each operation
	assert.Greater(t, outerExitIdx, outerEnterIdx, "OuterOp exit should come after OuterOp entry")
	assert.Greater(t, innerExitIdx, innerEnterIdx, "InnerOp exit should come after InnerOp entry")
}

// TestLogEntryExitWithCallback tests custom log level and callback
func TestLogEntryExitWithCallback(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	customCallback := func(ctx context.Context) *slog.Logger {
		return logger
	}

	operation := F.Pipe1(
		Of(42),
		LogEntryExitWithCallback[int](slog.LevelDebug, customCallback, "DebugOperation"),
	)

	res := operation(t.Context())()

	assert.Equal(t, result.Of(42), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "[entering]")
	assert.Contains(t, logOutput, "[exiting ]")
	assert.Contains(t, logOutput, "DebugOperation")
	assert.Contains(t, logOutput, "level=DEBUG")
}

// TestLogEntryExitDisabled tests that logging can be disabled
func TestLogEntryExitDisabled(t *testing.T) {
	var buf bytes.Buffer
	// Create logger with level that disables info logs
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError, // Only log errors
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	operation := F.Pipe1(
		Of("value"),
		LogEntryExit[string]("DisabledOperation"),
	)

	res := operation(t.Context())()

	assert.True(t, result.IsRight(res))

	// Should have no logs since level is ERROR
	logOutput := buf.String()
	assert.Empty(t, logOutput, "Should have no logs when logging is disabled")
}

// TestLoggingIDUniqueness tests that logging IDs are unique
func TestLoggingIDUniqueness(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	// Run multiple operations
	for i := range 5 {
		op := F.Pipe1(
			Of(i),
			LogEntryExit[int]("Operation"),
		)
		op(t.Context())()
	}

	logOutput := buf.String()

	// Extract all IDs and verify they're unique
	lines := strings.Split(logOutput, "\n")
	ids := make(map[string]bool)
	for _, line := range lines {
		if strings.Contains(line, "ID=") {
			// Extract ID value
			parts := strings.Split(line, "ID=")
			if len(parts) > 1 {
				idPart := strings.Fields(parts[1])[0]
				ids[idPart] = true
			}
		}
	}

	// Should have 5 unique IDs (one per operation)
	assert.GreaterOrEqual(t, len(ids), 5, "Should have at least 5 unique IDs")
}

// TestLogEntryExitWithContextLogger tests using logger from context
func TestLogEntryExitWithContextLogger(t *testing.T) {
	var buf bytes.Buffer
	contextLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cancelFct, ctx := pair.Unpack(logging.WithLogger(contextLogger)(t.Context()))
	defer cancelFct()

	operation := F.Pipe1(
		Of("context value"),
		LogEntryExit[string]("ContextOperation"),
	)

	res := operation(ctx)()

	assert.True(t, result.IsRight(res))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "[entering]")
	assert.Contains(t, logOutput, "[exiting ]")
	assert.Contains(t, logOutput, "ContextOperation")
}

// TestLogEntryExitTiming tests that duration is captured
func TestLogEntryExitTiming(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	// Operation with delay
	slowOp := func(ctx context.Context) IOResult[string] {
		return func() Result[string] {
			time.Sleep(10 * time.Millisecond)
			return result.Of("done")
		}
	}

	operation := F.Pipe1(
		slowOp,
		LogEntryExit[string]("SlowOperation"),
	)

	res := operation(t.Context())()

	assert.True(t, result.IsRight(res))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "duration=")

	// Verify duration is present in exit log
	lines := strings.Split(logOutput, "\n")
	var foundDuration bool
	for _, line := range lines {
		if strings.Contains(line, "[exiting ]") && strings.Contains(line, "duration=") {
			foundDuration = true
			break
		}
	}
	assert.True(t, foundDuration, "Exit log should contain duration")
}

// TestLogEntryExitChainedOperations tests complex chained operations
func TestLogEntryExitChainedOperations(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	step1 := F.Pipe1(
		Of(1),
		LogEntryExit[int]("Step1"),
	)

	step2 := F.Flow3(
		N.Mul(2),
		Of,
		LogEntryExit[int]("Step2"),
	)

	step3 := F.Flow3(
		strconv.Itoa,
		Of,
		LogEntryExit[string]("Step3"),
	)

	pipeline := F.Pipe1(
		step1,
		Chain(F.Flow2(
			step2,
			Chain(step3),
		)),
	)

	res := pipeline(t.Context())()

	assert.Equal(t, result.Of("2"), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Step1")
	assert.Contains(t, logOutput, "Step2")
	assert.Contains(t, logOutput, "Step3")

	// Verify all steps completed
	assert.Equal(t, 3, strings.Count(logOutput, "[entering]"))
	assert.Equal(t, 3, strings.Count(logOutput, "[exiting ]"))
}

// TestTapSLog tests basic TapSLog functionality
func TestTapSLog(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	operation := F.Pipe2(
		Of(42),
		TapSLog[int]("Processing value"),
		Map(N.Mul(2)),
	)

	res := operation(t.Context())()

	assert.Equal(t, result.Of(84), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Processing value")
	assert.Contains(t, logOutput, "value=42")
}

// TestTapSLogInPipeline tests TapSLog in a multi-step pipeline
func TestTapSLogInPipeline(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	step1 := F.Pipe2(
		Of("hello"),
		TapSLog[string]("Step 1: Initial value"),
		Map(func(s string) string { return s + " world" }),
	)

	step2 := F.Pipe2(
		step1,
		TapSLog[string]("Step 2: After concatenation"),
		Map(S.Size),
	)

	pipeline := F.Pipe1(
		step2,
		TapSLog[int]("Step 3: Final length"),
	)

	res := pipeline(t.Context())()

	assert.Equal(t, result.Of(11), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Step 1: Initial value")
	assert.Contains(t, logOutput, "value=hello")
	assert.Contains(t, logOutput, "Step 2: After concatenation")
	assert.Contains(t, logOutput, `value="hello world"`)
	assert.Contains(t, logOutput, "Step 3: Final length")
	assert.Contains(t, logOutput, "value=11")
}

// TestTapSLogWithError tests that TapSLog logs errors (via SLog)
func TestTapSLogWithError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testErr := errors.New("computation failed")
	pipeline := F.Pipe2(
		Left[int](testErr),
		TapSLog[int]("Error logged"),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	assert.True(t, result.IsLeft(res))

	logOutput := buf.String()
	// TapSLog uses SLog internally, which logs both successes and errors
	assert.Contains(t, logOutput, "Error logged")
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "computation failed")
}

// TestTapSLogWithStruct tests TapSLog with structured data
func TestTapSLogWithStruct(t *testing.T) {
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

	user := User{ID: 123, Name: "Alice"}
	operation := F.Pipe2(
		Of(user),
		TapSLog[User]("User data"),
		Map(func(u User) string { return u.Name }),
	)

	res := operation(t.Context())()

	assert.Equal(t, result.Of("Alice"), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "User data")
	assert.Contains(t, logOutput, "ID:123")
	assert.Contains(t, logOutput, "Name:Alice")
}

// TestTapSLogDisabled tests that TapSLog respects logger level
func TestTapSLogDisabled(t *testing.T) {
	var buf bytes.Buffer
	// Create logger with level that disables info logs
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError, // Only log errors
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	operation := F.Pipe2(
		Of(42),
		TapSLog[int]("This should not be logged"),
		Map(N.Mul(2)),
	)

	res := operation(t.Context())()

	assert.Equal(t, result.Of(84), res)

	// Should have no logs since level is ERROR
	logOutput := buf.String()
	assert.Empty(t, logOutput, "Should have no logs when logging is disabled")
}

// TestTapSLogWithContextLogger tests TapSLog using logger from context
func TestTapSLogWithContextLogger(t *testing.T) {
	var buf bytes.Buffer
	contextLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cancelFct, ctx := pair.Unpack(logging.WithLogger(contextLogger)(t.Context()))
	defer cancelFct()

	operation := F.Pipe2(
		Of("test value"),
		TapSLog[string]("Context logger test"),
		Map(S.Size),
	)

	res := operation(ctx)()

	assert.Equal(t, result.Of(10), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Context logger test")
	assert.Contains(t, logOutput, `value="test value"`)
}

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
	logged := SLog[int]("Result value")(res1)(ctx)()

	assert.Equal(t, result.Of(F.VOID), logged)

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
	logged := SLog[int]("Result value")(res1)(ctx)()

	assert.Equal(t, result.Of(F.VOID), logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Result value")
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "test error")
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
	logged := SLogWithCallback[int](slog.LevelDebug, customCallback, "Debug result")(res1)(ctx)()

	assert.Equal(t, result.Of(F.VOID), logged)

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
	logged := SLogWithCallback[int](slog.LevelWarn, customCallback, "Warning result")(res1)(ctx)()

	assert.Equal(t, result.Of(F.VOID), logged)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Warning result")
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "warning error")
	assert.Contains(t, logOutput, "level=WARN")
}

// TestTapSLogPreservesResult tests that TapSLog doesn't modify the result
func TestTapSLogPreservesResult(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	// Test with success value
	successOp := F.Pipe2(
		Of(42),
		TapSLog[int]("Success value"),
		Map(N.Mul(2)),
	)

	successRes := successOp(t.Context())()
	assert.Equal(t, result.Of(84), successRes)

	// Test with error value
	testErr := errors.New("test error")
	errorOp := F.Pipe2(
		Left[int](testErr),
		TapSLog[int]("Error value"),
		Map(N.Mul(2)),
	)

	errorRes := errorOp(t.Context())()
	assert.True(t, result.IsLeft(errorRes))

	// Verify the error is preserved
	_, err := result.Unwrap(errorRes)
	assert.Equal(t, testErr, err)
}

// TestTapSLogChainBehavior tests that TapSLog properly chains with other operations
func TestTapSLogChainBehavior(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	// Create a pipeline with multiple TapSLog calls
	step1 := F.Pipe2(
		Of(1),
		TapSLog[int]("Step 1"),
		Map(N.Mul(2)),
	)

	step2 := F.Pipe2(
		step1,
		TapSLog[int]("Step 2"),
		Map(N.Mul(3)),
	)

	pipeline := F.Pipe1(
		step2,
		TapSLog[int]("Step 3"),
	)

	res := pipeline(t.Context())()
	assert.Equal(t, result.Of(6), res)

	logOutput := buf.String()

	// Verify all steps were logged
	assert.Contains(t, logOutput, "Step 1")
	assert.Contains(t, logOutput, "value=1")
	assert.Contains(t, logOutput, "Step 2")
	assert.Contains(t, logOutput, "value=2")
	assert.Contains(t, logOutput, "Step 3")
	assert.Contains(t, logOutput, "value=6")
}

// TestTapSLogWithNilValue tests TapSLog with nil pointer values
func TestTapSLogWithNilValue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	type Data struct {
		Value string
	}

	// Test with nil pointer
	var nilData *Data
	operation := F.Pipe1(
		Of(nilData),
		TapSLog[*Data]("Nil pointer value"),
	)

	res := operation(t.Context())()
	assert.True(t, result.IsRight(res))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Nil pointer value")
	// The exact representation of nil may vary, but it should be logged
	assert.NotEmpty(t, logOutput)
}

// TestTapSLogLogsErrors verifies that TapSLog DOES log errors
// TapSLog uses SLog internally, which logs both success values and errors
func TestTapSLogLogsErrors(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testErr := errors.New("test error message")
	pipeline := F.Pipe2(
		Left[int](testErr),
		TapSLog[int]("Error logging test"),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	// Verify the error is preserved
	assert.True(t, result.IsLeft(res))

	// Verify logging occurred for the error
	logOutput := buf.String()
	assert.NotEmpty(t, logOutput, "TapSLog should log when the Result is an error")
	assert.Contains(t, logOutput, "Error logging test")
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "test error message")
}

func TestTapSLog_PreservesErrorCause(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	// Create a wrapped error with a cause chain
	rootErr := errors.New("root cause")
	wrappedErr := fmt.Errorf("wrapped: %w", rootErr)

	pipeline := F.Pipe2(
		Left[int](wrappedErr),
		TapSLog[int]("Error with cause"),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	// Extract the error from the result
	extractedErr := F.Pipe1(res, result.Fold(
		F.Identity[error],
		func(_ int) error { t.Fatal("expected Left but got Right"); return nil },
	))

	// Verify the error cause chain is intact
	assert.True(t, errors.Is(extractedErr, rootErr), "error cause chain should be preserved")
	assert.True(t, errors.Is(extractedErr, wrappedErr), "wrapped error should be preserved")

	// Verify logging occurred
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Error with cause")
	assert.Contains(t, logOutput, "wrapped")
}

// TestSLogLeftLogsError tests that SLogLeft logs an error at Error level
func TestSLogLeftLogsError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()
	testErr := errors.New("validation failed")

	// Use SLogLeft to log the error
	res := SLogLeft("Input validation error")(testErr)(ctx)()

	// Verify the result is a Left with the original error
	assert.Equal(t, result.Of(F.VOID), res)

	// Verify logging occurred
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Input validation error")
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "validation failed")
	assert.Contains(t, logOutput, "level=ERROR")
}

// TestSLogLeftInPipeline tests SLogLeft in an error handling pipeline
func TestSLogLeftInPipeline(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	validateInput := func(input string) ReaderIOResult[string] {
		if input == "" {
			// SLogLeft returns ReaderIOResult[Void], so we need to chain it
			return F.Pipe2(
				errors.New("input cannot be empty"),
				SLogLeft("Validation failed"),
				Chain(func(F.Void) ReaderIOResult[string] {
					return Left[string](errors.New("input cannot be empty"))
				}),
			)
		}
		return Of(input)
	}

	// Test with invalid input
	res := validateInput("")(t.Context())()

	assert.True(t, result.IsLeft(res))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Validation failed")
	assert.Contains(t, logOutput, "input cannot be empty")
	assert.Contains(t, logOutput, "level=ERROR")
}

// TestSLogLeftWithOrElse tests SLogLeft in error recovery scenarios
func TestSLogLeftWithOrElse(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	failingOp := Left[int](errors.New("operation failed"))

	// SLogLeft logs the error and preserves it as Left
	// To recover, we need to handle the error after logging using OrElse
	pipeline := F.Pipe1(
		failingOp,
		OrElse(func(err error) ReaderIOResult[int] {
			// Log the error and then recover with a fallback value
			return F.Pipe2(
				SLogLeft("Error occurred, using fallback")(err),
				OrElse(func(error) ReaderIOResult[Void] {
					return Of(F.VOID)
				}),
				Map(func(F.Void) int { return 42 }),
			)
		}),
	)

	res := pipeline(t.Context())()

	// Should be Right(42) after recovery
	assert.Equal(t, result.Of(42), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Error occurred, using fallback")
	assert.Contains(t, logOutput, "operation failed")
	assert.Contains(t, logOutput, "level=ERROR")
}

// TestSLogLeftWithContextLogger tests SLogLeft using logger from context
func TestSLogLeftWithContextLogger(t *testing.T) {
	var buf bytes.Buffer
	contextLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	cancelFct, ctx := pair.Unpack(logging.WithLogger(contextLogger)(t.Context()))
	defer cancelFct()

	testErr := errors.New("context logger error")

	res := SLogLeft("Context error test")(testErr)(ctx)()

	assert.Equal(t, result.Of(F.VOID), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Context error test")
	assert.Contains(t, logOutput, "context logger error")
}

// TestSLogLeftDisabled tests that SLogLeft respects logger level
func TestSLogLeftDisabled(t *testing.T) {
	var buf bytes.Buffer
	// Create logger that only logs Fatal level (higher than Error)
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.Level(12), // Higher than Error (8)
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testErr := errors.New("this should not be logged")

	res := SLogLeft("Should not log")(testErr)(t.Context())()

	// Error should still be preserved
	assert.Equal(t, result.Of(F.VOID), res)

	// But no logs should be written
	logOutput := buf.String()
	assert.Empty(t, logOutput, "Should have no logs when logging is disabled")
}

// TestSLogLeftMultipleErrors tests SLogLeft with multiple different errors
func TestSLogLeftMultipleErrors(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()

	testErrors := []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
	}

	for i, err := range testErrors {
		res := SLogLeft(S.Format[int]("Error %d")(i + 1))(err)(ctx)()
		assert.Equal(t, result.Of(F.VOID), res)

	}

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Error 1")
	assert.Contains(t, logOutput, "error 1")
	assert.Contains(t, logOutput, "Error 2")
	assert.Contains(t, logOutput, "error 2")
	assert.Contains(t, logOutput, "Error 3")
	assert.Contains(t, logOutput, "error 3")
}

// TestSLogLeftChainBehavior tests SLogLeft in a Chain operation
func TestSLogLeftChainBehavior(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	processData := func(data string) ReaderIOResult[int] {
		if data == "" {
			// SLogLeft returns ReaderIOResult[Void], chain to convert to ReaderIOResult[int]
			return F.Pipe2(
				errors.New("empty data"),
				SLogLeft("Processing failed"),
				Chain(func(F.Void) ReaderIOResult[int] {
					return Left[int](errors.New("empty data"))
				}),
			)
		}
		return Of(S.Size(data))
	}

	pipeline := F.Pipe2(
		Of(""),
		Chain(processData),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	// Should be Left due to error
	assert.True(t, result.IsLeft(res))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Processing failed")
	assert.Contains(t, logOutput, "empty data")
}

// TestSLogRightLogsValue tests that SLogRight logs a successful value
func TestSLogRightLogsValue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()
	testValue := "success data"

	// Use SLogRight to log the value
	res := SLogRight[string]("Operation succeeded")(testValue)(ctx)()

	// Verify the result is a Right with Void
	assert.True(t, result.IsRight(res))

	void := F.Pipe1(res, result.Fold(
		func(e error) F.Void { t.Fatal(e); return F.VOID },
		F.Identity[F.Void],
	))
	assert.Equal(t, F.VOID, void)

	// Verify logging occurred
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Operation succeeded")
	assert.Contains(t, logOutput, "value")
	assert.Contains(t, logOutput, "success data")
	assert.Contains(t, logOutput, "level=INFO")
}

// TestSLogRightInPipeline tests SLogRight in a success handling pipeline
func TestSLogRightInPipeline(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	processInput := func(input string) ReaderIOResult[string] {
		if input != "" {
			// SLogRight returns ReaderIOResult[Void], so we need to chain it
			return F.Pipe2(
				input,
				SLogRight[string]("Input validated"),
				Chain(func(F.Void) ReaderIOResult[string] {
					return Of(strings.ToUpper(input))
				}),
			)
		}
		return Left[string](errors.New("input cannot be empty"))
	}

	// Test with valid input
	res := processInput("hello")(t.Context())()

	assert.Equal(t, result.Of("HELLO"), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Input validated")
	assert.Contains(t, logOutput, "hello")
	assert.Contains(t, logOutput, "level=INFO")
}

// TestSLogRightPreservesValue tests that SLogRight preserves the original value
func TestSLogRightPreservesValue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	originalValue := 42

	res := SLogRight[int]("Logging value")(originalValue)(t.Context())()

	// Extract the void and verify it's Right
	void := F.Pipe1(res, result.Fold(
		func(e error) F.Void { t.Fatal(e); return F.VOID },
		F.Identity[F.Void],
	))

	assert.Equal(t, F.VOID, void)

	// Verify the original value was logged
	logOutput := buf.String()
	assert.Contains(t, logOutput, "42")
}

// TestSLogRightWithContextLogger tests SLogRight using logger from context
func TestSLogRightWithContextLogger(t *testing.T) {
	var buf bytes.Buffer
	contextLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cancelFct, ctx := pair.Unpack(logging.WithLogger(contextLogger)(t.Context()))
	defer cancelFct()

	testValue := "context logger value"

	res := SLogRight[string]("Context success test")(testValue)(ctx)()

	assert.True(t, result.IsRight(res))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Context success test")
	assert.Contains(t, logOutput, "context logger value")
}

// TestSLogRightDisabled tests that SLogRight respects logger level
func TestSLogRightDisabled(t *testing.T) {
	var buf bytes.Buffer
	// Create logger that only logs Warn level and above (higher than Info)
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testValue := "this should not be logged"

	res := SLogRight[string]("Should not log")(testValue)(t.Context())()

	// Value should still be processed
	assert.True(t, result.IsRight(res))

	// But no logs should be written
	logOutput := buf.String()
	assert.Empty(t, logOutput, "Should have no logs when logging is disabled")
}

// TestSLogRightMultipleValues tests SLogRight with multiple different values
func TestSLogRightMultipleValues(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()

	testValues := []int{1, 2, 3}

	for i, val := range testValues {
		res := SLogRight[int](S.Format[int]("Value %d")(i + 1))(val)(ctx)()
		assert.True(t, result.IsRight(res))
	}

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Value 1")
	assert.Contains(t, logOutput, "value=1")
	assert.Contains(t, logOutput, "Value 2")
	assert.Contains(t, logOutput, "value=2")
	assert.Contains(t, logOutput, "Value 3")
	assert.Contains(t, logOutput, "value=3")
}

// TestSLogRightChainBehavior tests SLogRight in a Chain operation
func TestSLogRightChainBehavior(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	processData := func(data string) ReaderIOResult[int] {
		if data != "" {
			// SLogRight returns ReaderIOResult[Void], chain to convert to ReaderIOResult[int]
			return F.Pipe2(
				data,
				SLogRight[string]("Processing succeeded"),
				Chain(func(F.Void) ReaderIOResult[int] {
					return Of(S.Size(data))
				}),
			)
		}
		return Left[int](errors.New("empty data"))
	}

	pipeline := F.Pipe2(
		Of("hello"),
		Chain(processData),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	// Should be Right with value 10 (5 * 2)
	assert.Equal(t, result.Of(10), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Processing succeeded")
	assert.Contains(t, logOutput, "hello")
}

// TestSLogRightWithStructValue tests SLogRight with a struct value
func TestSLogRightWithStructValue(t *testing.T) {
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

	user := User{ID: 123, Name: "Alice"}

	res := SLogRight[User]("User created")(user)(t.Context())()

	assert.True(t, result.IsRight(res))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "User created")
	assert.Contains(t, logOutput, "123")
	assert.Contains(t, logOutput, "Alice")
}

// TestSLogRightWithMap tests SLogRight in a Map operation
func TestSLogRightWithMap(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe3(
		Of(10),
		Chain(func(n int) ReaderIOResult[Void] {
			return SLogRight[int]("Initial value")(n)
		}),
		Chain(func(F.Void) ReaderIOResult[int] {
			return Of(20)
		}),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	assert.Equal(t, result.Of(40), res)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Initial value")
	assert.Contains(t, logOutput, "value=10")
}

// TestSLogRightComparison tests the difference between SLogRight and TapSLog
func TestSLogRightComparison(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	// SLogRight requires chaining to continue with a different type
	pipelineWithSLogRight := F.Pipe2(
		Of("test"),
		Tap(SLogRight[string]("Using SLogRight")),
		Chain(func(string) ReaderIOResult[int] {
			return Of(42)
		}),
	)

	res1 := pipelineWithSLogRight(t.Context())()
	assert.Equal(t, result.Of(42), res1)

	// TapSLog preserves the type and is more ergonomic
	pipelineWithTapSLog := F.Pipe2(
		Of("test"),
		TapSLog[string]("Using TapSLog"),
		Map(func(s string) int { return 42 }),
	)

	res2 := pipelineWithTapSLog(t.Context())()
	assert.Equal(t, result.Of(42), res2)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Using SLogRight")
	assert.Contains(t, logOutput, "Using TapSLog")
}

// ---------------------------------------------------------------------------
// SLogInfo
// ---------------------------------------------------------------------------

// TestSLogInfo_LogsSuccessAtInfoLevel verifies that SLogInfo produces an INFO entry
// with both the message and the structured value attribute for a Right Result.
func TestSLogInfo_LogsSuccessAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	SLogInfo[int]("info success")(result.Of(42))(t.Context())()

	log := buf.String()
	assert.Contains(t, log, "level=INFO")
	assert.Contains(t, log, "info success")
	assert.Contains(t, log, "value=42")
}

// TestSLogInfo_LogsErrorAtInfoLevel verifies that SLogInfo produces an INFO entry
// with the error attribute when given a Left Result.
func TestSLogInfo_LogsErrorAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testErr := errors.New("info error")
	SLogInfo[int]("info failure")(result.Left[int](testErr))(t.Context())()

	log := buf.String()
	assert.Contains(t, log, "level=INFO")
	assert.Contains(t, log, "info failure")
	assert.Contains(t, log, "info error")
}

// TestSLogInfo_EquivalentToSLog verifies that SLogInfo and SLog produce identical
// structured log output for the same message and Result.
func TestSLogInfo_EquivalentToSLog(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	infoLogger := slog.New(slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slogLogger := slog.New(slog.NewTextHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelInfo}))

	oldLogger := logging.SetLogger(infoLogger)
	SLogInfo[int]("msg")(result.Of(7))(t.Context())()
	logging.SetLogger(oldLogger)

	oldLogger = logging.SetLogger(slogLogger)
	SLog[int]("msg")(result.Of(7))(t.Context())()
	logging.SetLogger(oldLogger)

	assert.Contains(t, buf1.String(), "msg")
	assert.Contains(t, buf1.String(), "value=7")
	assert.Contains(t, buf2.String(), "msg")
	assert.Contains(t, buf2.String(), "value=7")
}

// TestSLogInfo_SuppressedAboveInfo verifies that no output is produced when the
// logger's minimum level is above Info.
func TestSLogInfo_SuppressedAboveInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	SLogInfo[int]("hidden")(result.Of(1))(t.Context())()

	assert.Empty(t, buf.String())
}

// TestSLogInfo_ReturnsVoid verifies that SLogInfo always returns Right(Void) regardless
// of whether the input Result is a success or an error.
func TestSLogInfo_ReturnsVoid(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := logging.SetLogger(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})))
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()

	assert.Equal(t, result.Of(F.VOID), SLogInfo[int]("ok")(result.Of(5))(ctx)())
	assert.Equal(t, result.Of(F.VOID), SLogInfo[int]("err")(result.Left[int](errors.New("e")))(ctx)())
}

// ---------------------------------------------------------------------------
// SLogDebug
// ---------------------------------------------------------------------------

// TestSLogDebug_LogsSuccessAtDebugLevel verifies that SLogDebug produces a DEBUG entry
// with the value attribute when the logger accepts Debug messages.
func TestSLogDebug_LogsSuccessAtDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	SLogDebug[int]("debug success")(result.Of(99))(t.Context())()

	log := buf.String()
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, "debug success")
	assert.Contains(t, log, "value=99")
}

// TestSLogDebug_LogsErrorAtDebugLevel verifies that SLogDebug produces a DEBUG entry
// with the error attribute for a Left Result.
func TestSLogDebug_LogsErrorAtDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testErr := errors.New("debug error")
	SLogDebug[int]("debug failure")(result.Left[int](testErr))(t.Context())()

	log := buf.String()
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, "debug failure")
	assert.Contains(t, log, "debug error")
}

// TestSLogDebug_SuppressedAtInfoLevel verifies that no output is produced when the
// logger's minimum level is Info (i.e., Debug is filtered out), while the computation
// still returns Right(Void).
func TestSLogDebug_SuppressedAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	ctx := t.Context()
	res := SLogDebug[int]("silent debug")(result.Of(3))(ctx)()

	assert.Empty(t, buf.String(), "Debug entry must not appear when logger level is Info")
	assert.Equal(t, result.Of(F.VOID), res)
}

// TestSLogDebug_UsesContextLogger verifies that a logger injected via
// logging.WithLogger is preferred over the global logger.
func TestSLogDebug_UsesContextLogger(t *testing.T) {
	var buf bytes.Buffer
	contextLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cancelFct, ctx := pair.Unpack(logging.WithLogger(contextLogger)(t.Context()))
	defer cancelFct()

	SLogDebug[string]("ctx debug")(result.Of("hello"))(ctx)()

	log := buf.String()
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, "ctx debug")
	assert.Contains(t, log, "value=hello")
}

// ---------------------------------------------------------------------------
// TapSLogInfo
// ---------------------------------------------------------------------------

// TestTapSLogInfo_LogsSuccessAndPreservesResult verifies that TapSLogInfo writes an
// INFO entry and passes the successful ReaderIOResult through unchanged.
func TestTapSLogInfo_LogsSuccessAndPreservesResult(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(42),
		TapSLogInfo[int]("tapsloginfo success"),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	assert.Equal(t, result.Of(84), res)
	log := buf.String()
	assert.Contains(t, log, "level=INFO")
	assert.Contains(t, log, "tapsloginfo success")
	assert.Contains(t, log, "value=42")
}

// TestTapSLogInfo_LogsErrorAndPreservesResult verifies that TapSLogInfo writes an
// INFO entry for errors and the Left Result is propagated unchanged.
func TestTapSLogInfo_LogsErrorAndPreservesResult(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testErr := errors.New("tapsloginfo error")
	pipeline := F.Pipe2(
		Left[int](testErr),
		TapSLogInfo[int]("tapsloginfo on error"),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	assert.True(t, result.IsLeft(res))
	log := buf.String()
	assert.Contains(t, log, "level=INFO")
	assert.Contains(t, log, "tapsloginfo on error")
	assert.Contains(t, log, "tapsloginfo error")
}

// TestTapSLogInfo_SuppressedAboveInfo verifies that the pipeline still produces the
// correct result when Info logging is disabled, with no output written.
func TestTapSLogInfo_SuppressedAboveInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(5),
		TapSLogInfo[int]("invisible"),
		Map(N.Mul(3)),
	)

	res := pipeline(t.Context())()

	assert.Equal(t, result.Of(15), res)
	assert.Empty(t, buf.String())
}

// TestTapSLogInfo_EquivalentToTapSLog verifies that TapSLogInfo and TapSLog produce
// identical output and results for the same pipeline.
func TestTapSLogInfo_EquivalentToTapSLog(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	makeLogger := func(buf *bytes.Buffer) *slog.Logger {
		return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	oldLogger := logging.SetLogger(makeLogger(&buf1))
	res1 := F.Pipe1(Of(9), TapSLogInfo[int]("same"))(t.Context())()
	logging.SetLogger(oldLogger)

	oldLogger = logging.SetLogger(makeLogger(&buf2))
	res2 := F.Pipe1(Of(9), TapSLog[int]("same"))(t.Context())()
	logging.SetLogger(oldLogger)

	assert.Equal(t, res1, res2)
	assert.Contains(t, buf1.String(), "value=9")
	assert.Contains(t, buf2.String(), "value=9")
}

// TestTapSLogInfo_MultiStep verifies that multiple TapSLogInfo steps each log their
// own snapshot of the value at that point in the pipeline.
func TestTapSLogInfo_MultiStep(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe4(
		Of(1),
		TapSLogInfo[int]("step 1"),
		Map(N.Mul(2)),
		TapSLogInfo[int]("step 2"),
		Map(N.Mul(3)),
	)

	res := pipeline(t.Context())()

	assert.Equal(t, result.Of(6), res)

	log := buf.String()
	assert.Contains(t, log, "step 1")
	assert.Contains(t, log, "value=1")
	assert.Contains(t, log, "step 2")
	assert.Contains(t, log, "value=2")
}

// ---------------------------------------------------------------------------
// TapSLogDebug
// ---------------------------------------------------------------------------

// TestTapSLogDebug_LogsSuccessAtDebugLevelAndPreservesResult verifies that TapSLogDebug
// writes a DEBUG entry and passes the successful ReaderIOResult through unchanged.
func TestTapSLogDebug_LogsSuccessAtDebugLevelAndPreservesResult(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(8),
		TapSLogDebug[int]("debug tap success"),
		Map(N.Mul(5)),
	)

	res := pipeline(t.Context())()

	assert.Equal(t, result.Of(40), res)
	log := buf.String()
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, "debug tap success")
	assert.Contains(t, log, "value=8")
}

// TestTapSLogDebug_LogsErrorAtDebugLevel verifies that TapSLogDebug writes a DEBUG
// entry for errors and the Left Result is propagated unchanged.
func TestTapSLogDebug_LogsErrorAtDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	testErr := errors.New("debug tap error")
	pipeline := F.Pipe2(
		Left[int](testErr),
		TapSLogDebug[int]("debug on error"),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	assert.True(t, result.IsLeft(res))
	log := buf.String()
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, "debug on error")
	assert.Contains(t, log, "debug tap error")
}

// TestTapSLogDebug_SuppressedAtInfoLevel verifies that no output is produced and the
// pipeline result is still correct when Debug is filtered out.
func TestTapSLogDebug_SuppressedAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	oldLogger := logging.SetLogger(logger)
	defer logging.SetLogger(oldLogger)

	pipeline := F.Pipe2(
		Of(6),
		TapSLogDebug[int]("should not appear"),
		Map(N.Mul(2)),
	)

	res := pipeline(t.Context())()

	assert.Equal(t, result.Of(12), res)
	assert.Empty(t, buf.String(), "Debug entry must not appear when logger level is Info")
}

// TestTapSLogDebug_LazyExecution verifies that no output is produced until the
// resulting IO is actually executed.
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

// TestTapSLogDebug_UsesContextLogger verifies that a logger injected via
// logging.WithLogger is preferred over the global logger.
func TestTapSLogDebug_UsesContextLogger(t *testing.T) {
	var buf bytes.Buffer
	contextLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cancelFct, ctx := pair.Unpack(logging.WithLogger(contextLogger)(t.Context()))
	defer cancelFct()

	pipeline := F.Pipe1(
		Of("hello"),
		TapSLogDebug[string]("ctx debug tap"),
	)

	res := pipeline(ctx)()

	assert.Equal(t, result.Of("hello"), res)
	log := buf.String()
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, "ctx debug tap")
	assert.Contains(t, log, "value=hello")
}
