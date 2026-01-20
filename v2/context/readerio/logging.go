package readerio

import (
	"context"
	"log/slog"

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
	message string) Kleisli[A, A] {
	return func(a A) ReaderIO[A] {
		return func(ctx context.Context) IO[A] {
			// logger
			logger := cb(ctx)
			return func() A {
				logger.LogAttrs(ctx, logLevel, message, slog.Any("value", a))
				return a
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
func SLog[A any](message string) Kleisli[A, A] {
	return SLogWithCallback[A](slog.LevelInfo, logging.GetLoggerFromContext, message)
}
