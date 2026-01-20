package readerioresult

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/logging"
	N "github.com/IBM/fp-go/v2/number"
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

// TestLogEntryExitF tests custom entry/exit callbacks
func TestLogEntryExitF(t *testing.T) {
	var entryCount, exitCount int

	onEntry := func(ctx context.Context) IO[context.Context] {
		return func() context.Context {
			entryCount++
			return ctx
		}
	}

	onExit := func(res Result[string]) ReaderIO[any] {
		return func(ctx context.Context) IO[any] {
			return func() any {
				exitCount++
				return nil
			}
		}
	}

	operation := F.Pipe1(
		Of("test"),
		LogEntryExitF(onEntry, onExit),
	)

	res := operation(t.Context())()

	assert.True(t, result.IsRight(res))
	assert.Equal(t, 1, entryCount, "Entry callback should be called once")
	assert.Equal(t, 1, exitCount, "Exit callback should be called once")
}

// TestLogEntryExitFWithError tests custom callbacks with error
func TestLogEntryExitFWithError(t *testing.T) {
	var entryCount, exitCount int
	var capturedError error

	onEntry := func(ctx context.Context) IO[context.Context] {
		return func() context.Context {
			entryCount++
			return ctx
		}
	}

	onExit := func(res Result[string]) ReaderIO[any] {
		return func(ctx context.Context) IO[any] {
			return func() any {
				exitCount++
				if result.IsLeft(res) {
					_, capturedError = result.Unwrap(res)
				}
				return nil
			}
		}
	}

	testErr := errors.New("custom error")
	operation := F.Pipe1(
		Left[string](testErr),
		LogEntryExitF(onEntry, onExit),
	)

	res := operation(t.Context())()

	assert.True(t, result.IsLeft(res))
	assert.Equal(t, 1, entryCount, "Entry callback should be called once")
	assert.Equal(t, 1, exitCount, "Exit callback should be called once")
	assert.Equal(t, testErr, capturedError, "Should capture the error")
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

	ctx := logging.WithLogger(contextLogger)(t.Context())

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

	ctx := logging.WithLogger(contextLogger)(t.Context())

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
	logged := SLog[int]("Result value")(res1)(ctx)()

	assert.True(t, result.IsLeft(logged))

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
	logged := SLogWithCallback[int](slog.LevelWarn, customCallback, "Warning result")(res1)(ctx)()

	assert.True(t, result.IsLeft(logged))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Warning result")
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "warning error")
	assert.Contains(t, logOutput, "level=WARN")
}
