package readerio

import (
	"context"
	"log/slog"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/logging"
)

// SLogWithCallback creates a Kleisli arrow that logs a value with a custom logger and log level.
// The value is logged and then passed through unchanged, making this useful for debugging
// and monitoring values as they flow through a ReaderIO computation.
//
// Type Parameters:
//   - A: The type of value to log and pass through
//
// Parameters:
//   - logLevel: The slog.Level to use for logging (e.g., slog.LevelInfo, slog.LevelDebug)
//   - cb: Callback function to retrieve the *slog.Logger from the context
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - A Kleisli arrow that logs the value and returns it unchanged
//
// Example:
//
//	getMyLogger := func(ctx context.Context) *slog.Logger {
//	    if logger := ctx.Value("logger"); logger != nil {
//	        return logger.(*slog.Logger)
//	    }
//	    return slog.Default()
//	}
//
//	debugLog := SLogWithCallback[User](
//	    slog.LevelDebug,
//	    getMyLogger,
//	    "Processing user",
//	)
//
//	pipeline := F.Pipe2(
//	    fetchUser(123),
//	    Chain(debugLog),
//	)
func SLogWithCallback[A any](
	logLevel slog.Level,
	cb func(context.Context) *slog.Logger,
	message string) Kleisli[A, Void] {
	return func(a A) ReaderIO[Void] {
		return func(ctx context.Context) IO[Void] {
			// logger
			logger := cb(ctx)
			return func() Void {
				logger.LogAttrs(ctx, logLevel, message, slog.Any("value", a))
				return F.VOID
			}
		}
	}
}

// SLog creates a Kleisli arrow that logs a value at Info level and passes it through unchanged.
// This is a convenience wrapper around SLogWithCallback with standard settings.
//
// The value is logged with the provided message and then returned unchanged, making this
// useful for debugging and monitoring values in a ReaderIO computation pipeline.
//
// Type Parameters:
//   - A: The type of value to log and pass through
//
// Parameters:
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - A Kleisli arrow that logs the value at Info level and returns it unchanged
//
// Example:
//
//	pipeline := F.Pipe3(
//	    fetchUser(123),
//	    Chain(SLog[User]("Fetched user")),
//	    Map(func(u User) string { return u.Name }),
//	    Chain(SLog[string]("Extracted name")),
//	)
//
//	result := pipeline(t.Context())()
//	// Logs: "Fetched user" value={ID:123 Name:"Alice"}
//	// Logs: "Extracted name" value="Alice"
//
//go:inline
func SLog[A any](message string) Kleisli[A, Void] {
	return SLogWithCallback[A](slog.LevelInfo, logging.GetLoggerFromContext, message)
}

// SLogInfo creates a Kleisli arrow that logs a value at Info level and passes it through unchanged.
// It is an explicit alias for SLog, provided for symmetry with SLogDebug when the desired
// log level needs to be clear at the call site.
//
// Type Parameters:
//   - A: The type of value to log and pass through
//
// Parameters:
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - A Kleisli arrow that logs the value at Info level and returns it unchanged
//
// See Also:
//   - SLog: The function this delegates to
//   - SLogDebug: The Debug-level counterpart
//   - SLogWithCallback: For logging with a custom logger callback or log level
//
//go:inline
func SLogInfo[A any](message string) Kleisli[A, Void] {
	return SLog[A](message)
}

// SLogDebug creates a Kleisli arrow that logs a value at Debug level and passes it through unchanged.
// This is useful for high-frequency trace points that should only appear in debug builds or
// when the logger's minimum level is set to slog.LevelDebug.
//
// The logger is retrieved from the context via GetLoggerFromContext. If the logger's
// configured level is above Debug, no output is produced and the value flows through unchanged.
//
// Type Parameters:
//   - A: The type of value to log and pass through
//
// Parameters:
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - A Kleisli arrow that logs the value at Debug level and returns it unchanged
//
// See Also:
//   - SLogInfo: The Info-level counterpart
//   - SLogWithCallback: For logging with a custom logger callback or log level
//
//go:inline
func SLogDebug[A any](message string) Kleisli[A, Void] {
	return SLogWithCallback[A](slog.LevelDebug, logging.GetLoggerFromContext, message)
}

// TapSLog creates an Operator that logs the current value at Info level and passes it through unchanged.
// This is a convenience wrapper that combines SLog with Tap, making it suitable for
// inserting non-intrusive log points into a ReaderIO computation pipeline.
//
// The value is logged using the logger retrieved from the context via GetLoggerFromContext.
// After logging, the original ReaderIO[A] is returned unchanged, so the type and value
// are fully preserved across the log step.
//
// Type Parameters:
//   - A: The type of value to log and pass through
//
// Parameters:
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - An Operator that logs the value at Info level and returns the input ReaderIO unchanged
//
// See Also:
//   - SLog: The underlying Kleisli arrow used to perform the log
//   - Tap: The operator used to sequence the log side-effect while preserving the value
//   - SLogWithCallback: For logging with a custom logger callback or log level
func TapSLog[A any](message string) Operator[A, A] {
	return F.Pipe2(
		message,
		SLog[A],
		Tap,
	)
}

// TapSLogInfo creates an Operator that logs the current value at Info level and passes it through unchanged.
// It is an explicit alias for TapSLog, provided for symmetry with TapSLogDebug when the desired
// log level needs to be clear at the call site.
//
// Type Parameters:
//   - A: The type of value to log and pass through
//
// Parameters:
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - An Operator that logs the value at Info level and returns the input ReaderIO unchanged
//
// See Also:
//   - TapSLog: The function this delegates to
//   - TapSLogDebug: The Debug-level counterpart
//   - SLogInfo: The underlying Kleisli arrow used to perform the log
func TapSLogInfo[A any](message string) Operator[A, A] {
	return F.Pipe2(
		message,
		SLogInfo[A],
		Tap,
	)
}

// TapSLogDebug creates an Operator that logs the current value at Debug level and passes it through unchanged.
// This is useful for inserting high-frequency trace points into a ReaderIO pipeline without
// producing output unless the logger's minimum level is set to slog.LevelDebug.
//
// The logger is retrieved from the context via GetLoggerFromContext on each execution.
// If the logger's configured level is above Debug, no output is produced and the value
// flows through unchanged.
//
// Type Parameters:
//   - A: The type of value to log and pass through
//
// Parameters:
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - An Operator that logs the value at Debug level and returns the input ReaderIO unchanged
//
// See Also:
//   - TapSLogInfo: The Info-level counterpart
//   - TapSLog: The general Info-level tap operator
//   - SLogDebug: The underlying Kleisli arrow used to perform the log
func TapSLogDebug[A any](message string) Operator[A, A] {
	return F.Pipe2(
		message,
		SLogDebug[A],
		Tap,
	)
}
