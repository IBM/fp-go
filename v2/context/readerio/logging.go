package readerio

import (
	"context"
	"log/slog"

	"github.com/IBM/fp-go/v2/logging"
)

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

//go:inline
func SLog[A any](message string) Kleisli[A, A] {
	return SLogWithCallback[A](slog.LevelInfo, logging.GetLoggerFromContext, message)
}
