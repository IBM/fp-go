// Package circuitbreaker provides metrics collection for circuit breaker state transitions and events.
package circuitbreaker

import (
	"log"
	"time"

	"github.com/IBM/fp-go/v2/function"
)

type (
	// Metrics defines the interface for collecting circuit breaker metrics and events.
	// Implementations can use this interface to track circuit breaker behavior for
	// monitoring, alerting, and debugging purposes.
	//
	// All methods accept a time.Time parameter representing when the event occurred,
	// and return an IO[Void] operation that performs the metric recording when executed.
	//
	// Thread Safety: Implementations must be thread-safe as circuit breakers may be
	// accessed concurrently from multiple goroutines.
	//
	// Example Usage:
	//
	//	logger := log.New(os.Stdout, "[CircuitBreaker] ", log.LstdFlags)
	//	metrics := MakeMetricsFromLogger("API-Service", logger)
	//
	//	// In circuit breaker implementation
	//	io.Run(metrics.Accept(time.Now()))  // Record accepted request
	//	io.Run(metrics.Reject(time.Now()))  // Record rejected request
	//	io.Run(metrics.Open(time.Now()))    // Record circuit opening
	//	io.Run(metrics.Close(time.Now()))   // Record circuit closing
	//	io.Run(metrics.Canary(time.Now()))  // Record canary request
	Metrics interface {
		// Accept records that a request was accepted and allowed through the circuit breaker.
		// This is called when the circuit is closed or in half-open state (canary request).
		//
		// Parameters:
		//   - time.Time: The timestamp when the request was accepted
		//
		// Returns:
		//   - IO[Void]: An IO operation that records the acceptance when executed
		//
		// Thread Safety: Must be safe to call concurrently.
		Accept(time.Time) IO[Void]

		// Reject records that a request was rejected because the circuit breaker is open.
		// This is called when a request is blocked due to the circuit being in open state
		// and the reset time has not been reached.
		//
		// Parameters:
		//   - time.Time: The timestamp when the request was rejected
		//
		// Returns:
		//   - IO[Void]: An IO operation that records the rejection when executed
		//
		// Thread Safety: Must be safe to call concurrently.
		Reject(time.Time) IO[Void]

		// Open records that the circuit breaker transitioned to the open state.
		// This is called when the failure threshold is exceeded and the circuit opens
		// to prevent further requests from reaching the failing service.
		//
		// Parameters:
		//   - time.Time: The timestamp when the circuit opened
		//
		// Returns:
		//   - IO[Void]: An IO operation that records the state transition when executed
		//
		// Thread Safety: Must be safe to call concurrently.
		Open(time.Time) IO[Void]

		// Close records that the circuit breaker transitioned to the closed state.
		// This is called when:
		//   - A canary request succeeds in half-open state
		//   - The circuit is manually reset
		//   - The circuit breaker is initialized
		//
		// Parameters:
		//   - time.Time: The timestamp when the circuit closed
		//
		// Returns:
		//   - IO[Void]: An IO operation that records the state transition when executed
		//
		// Thread Safety: Must be safe to call concurrently.
		Close(time.Time) IO[Void]

		// Canary records that a canary (test) request is being attempted.
		// This is called when the circuit is in half-open state and a single test request
		// is allowed through to check if the service has recovered.
		//
		// Parameters:
		//   - time.Time: The timestamp when the canary request was initiated
		//
		// Returns:
		//   - IO[Void]: An IO operation that records the canary attempt when executed
		//
		// Thread Safety: Must be safe to call concurrently.
		Canary(time.Time) IO[Void]
	}

	// loggingMetrics is a simple implementation of the Metrics interface that logs
	// circuit breaker events using Go's standard log.Logger.
	//
	// This implementation is thread-safe as log.Logger is safe for concurrent use.
	//
	// Fields:
	//   - name: A human-readable name identifying the circuit breaker instance
	//   - logger: The log.Logger instance used for writing log messages
	loggingMetrics struct {
		name   string
		logger *log.Logger
	}
)

// doLog is a helper method that creates an IO operation for logging a circuit breaker event.
// It formats the log message with the event prefix, circuit breaker name, and timestamp.
//
// Parameters:
//   - prefix: The event type (e.g., "Accept", "Reject", "Open", "Close", "Canary")
//   - ct: The timestamp when the event occurred
//
// Returns:
//   - IO[Void]: An IO operation that logs the event when executed
//
// Thread Safety: Safe for concurrent use as log.Logger is thread-safe.
//
// Log Format: "<prefix>: <name>, <timestamp>"
// Example: "Open: API-Service, 2026-01-09 15:30:45.123 +0100 CET"
func (m *loggingMetrics) doLog(prefix string, ct time.Time) IO[Void] {
	return func() Void {
		m.logger.Printf("%s: %s, %s\n", prefix, m.name, ct)
		return function.VOID
	}
}

// Accept implements the Metrics interface for loggingMetrics.
// Logs when a request is accepted through the circuit breaker.
//
// Thread Safety: Safe for concurrent use.
func (m *loggingMetrics) Accept(ct time.Time) IO[Void] {
	return m.doLog("Accept", ct)
}

// Open implements the Metrics interface for loggingMetrics.
// Logs when the circuit breaker transitions to open state.
//
// Thread Safety: Safe for concurrent use.
func (m *loggingMetrics) Open(ct time.Time) IO[Void] {
	return m.doLog("Open", ct)
}

// Close implements the Metrics interface for loggingMetrics.
// Logs when the circuit breaker transitions to closed state.
//
// Thread Safety: Safe for concurrent use.
func (m *loggingMetrics) Close(ct time.Time) IO[Void] {
	return m.doLog("Close", ct)
}

// Reject implements the Metrics interface for loggingMetrics.
// Logs when a request is rejected because the circuit breaker is open.
//
// Thread Safety: Safe for concurrent use.
func (m *loggingMetrics) Reject(ct time.Time) IO[Void] {
	return m.doLog("Reject", ct)
}

// Canary implements the Metrics interface for loggingMetrics.
// Logs when a canary (test) request is attempted in half-open state.
//
// Thread Safety: Safe for concurrent use.
func (m *loggingMetrics) Canary(ct time.Time) IO[Void] {
	return m.doLog("Canary", ct)
}

// MakeMetricsFromLogger creates a Metrics implementation that logs circuit breaker events
// using the provided log.Logger.
//
// This is a simple metrics implementation suitable for development, debugging, and
// basic production monitoring. For more sophisticated metrics collection (e.g., Prometheus,
// StatsD), implement the Metrics interface with a custom type.
//
// Parameters:
//   - name: A human-readable name identifying the circuit breaker instance.
//     This name appears in all log messages to distinguish between multiple circuit breakers.
//   - logger: The log.Logger instance to use for writing log messages.
//     If nil, this will panic when metrics are recorded.
//
// Returns:
//   - Metrics: A thread-safe Metrics implementation that logs events
//
// Thread Safety: The returned Metrics implementation is safe for concurrent use
// as log.Logger is thread-safe.
//
// Example:
//
//	logger := log.New(os.Stdout, "[CB] ", log.LstdFlags)
//	metrics := MakeMetricsFromLogger("UserService", logger)
//
//	// Use with circuit breaker
//	io.Run(metrics.Open(time.Now()))
//	// Output: [CB] 2026/01/09 15:30:45 Open: UserService, 2026-01-09 15:30:45.123 +0100 CET
//
//	io.Run(metrics.Reject(time.Now()))
//	// Output: [CB] 2026/01/09 15:30:46 Reject: UserService, 2026-01-09 15:30:46.456 +0100 CET
func MakeMetricsFromLogger(name string, logger *log.Logger) Metrics {
	return &loggingMetrics{name: name, logger: logger}
}
