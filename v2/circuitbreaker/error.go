// Package circuitbreaker provides error types and utilities for circuit breaker implementations.
package circuitbreaker

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"time"

	E "github.com/IBM/fp-go/v2/errors"
	FH "github.com/IBM/fp-go/v2/http"
	"github.com/IBM/fp-go/v2/option"
)

// CircuitBreakerError represents an error that occurs when a circuit breaker is in the open state.
//
// When a circuit breaker opens due to too many failures, it prevents further operations
// from executing until a reset time is reached. This error type communicates that state
// and provides information about when the circuit breaker will attempt to close again.
//
// Fields:
//   - Name: The name identifying this circuit breaker instance
//   - ResetAt: The time at which the circuit breaker will transition from open to half-open state
//
// Thread Safety: This type is immutable and safe for concurrent use.
type CircuitBreakerError struct {
	// Name: The name identifying this circuit breaker instance
	Name string

	// ResetAt: The time at which the circuit breaker will transition from open to half-open state
	ResetAt time.Time
}

// Error implements the error interface for CircuitBreakerError.
//
// Returns a formatted error message indicating that the circuit breaker is open
// and when it will attempt to close.
//
// Returns:
//   - A string describing the circuit breaker state and reset time
//
// Thread Safety: This method is safe for concurrent use as it only reads immutable fields.
//
// Example:
//
//	err := &CircuitBreakerError{Name: "API", ResetAt: time.Now().Add(30 * time.Second)}
//	fmt.Println(err.Error())
//	// Output: circuit breaker is open [API], will close at 2026-01-09 12:20:47.123 +0100 CET
func (e *CircuitBreakerError) Error() string {
	return fmt.Sprintf("circuit breaker is open [%s], will close at %s", e.Name, e.ResetAt)
}

// MakeCircuitBreakerErrorWithName creates a circuit breaker error constructor with a custom name.
//
// This function returns a constructor that creates CircuitBreakerError instances with a specific
// circuit breaker name. This is useful when you have multiple circuit breakers in your system
// and want to identify which one is open in error messages.
//
// Parameters:
//   - name: The name to identify this circuit breaker in error messages
//
// Returns:
//   - A function that takes a reset time and returns a CircuitBreakerError with the specified name
//
// Thread Safety: The returned function is safe for concurrent use as it creates new error
// instances on each call.
//
// Example:
//
//	makeDBError := MakeCircuitBreakerErrorWithName("Database Circuit Breaker")
//	err := makeDBError(time.Now().Add(30 * time.Second))
//	fmt.Println(err.Error())
//	// Output: circuit breaker is open [Database Circuit Breaker], will close at 2026-01-09 12:20:47.123 +0100 CET
func MakeCircuitBreakerErrorWithName(name string) func(time.Time) error {
	return func(resetTime time.Time) error {
		return &CircuitBreakerError{Name: name, ResetAt: resetTime}
	}
}

// MakeCircuitBreakerError creates a new CircuitBreakerError with the specified reset time.
//
// This constructor function creates a circuit breaker error that indicates when the
// circuit breaker will transition from the open state to the half-open state, allowing
// test requests to determine if the underlying service has recovered.
//
// Parameters:
//   - resetTime: The time at which the circuit breaker will attempt to close
//
// Returns:
//   - An error representing the circuit breaker open state
//
// Thread Safety: This function is safe for concurrent use as it creates new error
// instances on each call.
//
// Example:
//
//	resetTime := time.Now().Add(30 * time.Second)
//	err := MakeCircuitBreakerError(resetTime)
//	if cbErr, ok := err.(*CircuitBreakerError); ok {
//	    fmt.Printf("Circuit breaker will reset at: %s\n", cbErr.ResetAt)
//	}
var MakeCircuitBreakerError = MakeCircuitBreakerErrorWithName("Generic Circuit Breaker")

// AnyError converts an error to an Option, wrapping non-nil errors in Some and nil errors in None.
//
// This variable provides a functional way to handle errors by converting them to Option types.
// It's particularly useful in functional programming contexts where you want to treat errors
// as optional values rather than using traditional error handling patterns.
//
// Behavior:
//   - If the error is non-nil, returns Some(error)
//   - If the error is nil, returns None
//
// Thread Safety: This function is pure and safe for concurrent use.
//
// Example:
//
//	err := errors.New("something went wrong")
//	optErr := AnyError(err)  // Some(error)
//
//	var noErr error = nil
//	optNoErr := AnyError(noErr)  // None
//
//	// Using in functional pipelines
//	result := F.Pipe2(
//	    someOperation(),
//	    AnyError,
//	    O.Map(func(e error) string { return e.Error() }),
//	)
var AnyError = option.FromPredicate(E.IsNonNil)

// shouldOpenCircuit determines if an error should cause a circuit breaker to open.
//
// This function checks if an error represents an infrastructure or server problem
// that indicates the service is unhealthy and should trigger circuit breaker protection.
// It examines both the error type and, for HTTP errors, the status code.
//
// Errors that should open the circuit include:
//   - HTTP 5xx server errors (500-599) indicating server-side problems
//   - Network errors (connection refused, connection reset, timeouts)
//   - DNS resolution errors
//   - TLS/certificate errors
//   - Other infrastructure-related errors
//
// Errors that should NOT open the circuit include:
//   - HTTP 4xx client errors (bad request, unauthorized, not found, etc.)
//   - Application-level validation errors
//   - Business logic errors
//
// The function unwraps error chains to find the root cause, making it compatible
// with wrapped errors created by fmt.Errorf with %w or errors.Join.
//
// Parameters:
//   - err: The error to evaluate (may be nil)
//
// Returns:
//   - true if the error should cause the circuit to open, false otherwise
//
// Thread Safety: This function is pure and safe for concurrent use. It does not
// modify any state.
//
// Example:
//
//	// HTTP 500 error - should open circuit
//	httpErr := &FH.HttpError{...} // status 500
//	if shouldOpenCircuit(httpErr) {
//	    // Open circuit breaker
//	}
//
//	// HTTP 404 error - should NOT open circuit (client error)
//	notFoundErr := &FH.HttpError{...} // status 404
//	if !shouldOpenCircuit(notFoundErr) {
//	    // Don't open circuit, this is a client error
//	}
//
//	// Network timeout - should open circuit
//	timeoutErr := &net.OpError{Op: "dial", Err: syscall.ETIMEDOUT}
//	if shouldOpenCircuit(timeoutErr) {
//	    // Open circuit breaker
//	}
func shouldOpenCircuit(err error) bool {
	if err == nil {
		return false
	}

	// Check for HTTP errors with server status codes (5xx)
	var httpErr *FH.HttpError
	if errors.As(err, &httpErr) {
		statusCode := httpErr.StatusCode()
		// Only 5xx errors should open the circuit
		// 4xx errors are client errors and shouldn't affect circuit state
		return statusCode >= http.StatusInternalServerError && statusCode < 600
	}

	// Check for network operation errors
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		// Network timeouts should open the circuit
		if opErr.Timeout() {
			return true
		}
		// Check the underlying error
		if opErr.Err != nil {
			return isInfrastructureError(opErr.Err)
		}
		return true
	}

	// Check for DNS errors
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}

	// Check for URL errors (often wrap network errors)
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		if urlErr.Timeout() {
			return true
		}
		// Recursively check the wrapped error
		return shouldOpenCircuit(urlErr.Err)
	}

	// Check for specific syscall errors that indicate infrastructure problems
	return isInfrastructureError(err) || isTLSError(err)
}

// isInfrastructureError checks if an error is a low-level infrastructure error
// that should cause the circuit to open.
//
// This function examines syscall errors to identify network and system-level failures
// that indicate the service is unavailable or unreachable.
//
// Infrastructure errors include:
//   - ECONNREFUSED: Connection refused (service not listening)
//   - ECONNRESET: Connection reset by peer (service crashed or network issue)
//   - ECONNABORTED: Connection aborted (network issue)
//   - ENETUNREACH: Network unreachable (routing problem)
//   - EHOSTUNREACH: Host unreachable (host down or network issue)
//   - EPIPE: Broken pipe (connection closed unexpectedly)
//   - ETIMEDOUT: Operation timed out (service not responding)
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - true if the error is an infrastructure error, false otherwise
//
// Thread Safety: This function is pure and safe for concurrent use.
func isInfrastructureError(err error) bool {

	var syscallErr *syscall.Errno

	if errors.As(err, &syscallErr) {
		switch *syscallErr {
		case syscall.ECONNREFUSED,
			syscall.ECONNRESET,
			syscall.ECONNABORTED,
			syscall.ENETUNREACH,
			syscall.EHOSTUNREACH,
			syscall.EPIPE,
			syscall.ETIMEDOUT:
			return true
		}

	}
	return false
}

// isTLSError checks if an error is a TLS/certificate error that should cause the circuit to open.
//
// TLS errors typically indicate infrastructure or configuration problems that prevent
// secure communication with the service. These errors suggest the service is not properly
// configured or accessible.
//
// TLS errors include:
//   - Certificate verification failures (invalid, expired, or malformed certificates)
//   - Unknown certificate authority errors (untrusted CA)
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - true if the error is a TLS/certificate error, false otherwise
//
// Thread Safety: This function is pure and safe for concurrent use.
func isTLSError(err error) bool {
	// Certificate verification failed
	var certErr *x509.CertificateInvalidError
	if errors.As(err, &certErr) {
		return true
	}

	// Unknown authority
	var unknownAuthErr *x509.UnknownAuthorityError
	if errors.As(err, &unknownAuthErr) {
		return true
	}

	return false
}

// InfrastructureError is a predicate that converts errors to Options based on whether
// they should trigger circuit breaker opening.
//
// This variable provides a functional way to filter errors that represent infrastructure
// failures (network issues, server errors, timeouts, etc.) from application-level errors
// (validation errors, business logic errors, client errors).
//
// Behavior:
//   - Returns Some(error) if the error should open the circuit (infrastructure failure)
//   - Returns None if the error should not open the circuit (application error)
//
// Thread Safety: This function is pure and safe for concurrent use.
//
// Use this in circuit breaker configurations to determine which errors should count
// toward the failure threshold.
//
// Example:
//
//	// In a circuit breaker configuration
//	breaker := MakeCircuitBreaker(
//	    ...,
//	    checkError: InfrastructureError,  // Only infrastructure errors open the circuit
//	    ...,
//	)
//
//	// HTTP 500 error - returns Some(error)
//	result := InfrastructureError(&FH.HttpError{...}) // Some(error)
//
//	// HTTP 404 error - returns None
//	result := InfrastructureError(&FH.HttpError{...}) // None
var InfrastructureError = option.FromPredicate(shouldOpenCircuit)
