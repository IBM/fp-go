package circuitbreaker

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	FH "github.com/IBM/fp-go/v2/http"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestCircuitBreakerError tests the CircuitBreakerError type
func TestCircuitBreakerError(t *testing.T) {
	t.Run("Error returns formatted message with reset time", func(t *testing.T) {
		resetTime := time.Date(2026, 1, 9, 12, 30, 0, 0, time.UTC)
		err := &CircuitBreakerError{ResetAt: resetTime}

		result := err.Error()

		assert.Contains(t, result, "circuit breaker is open")
		assert.Contains(t, result, "will close at")
		assert.Contains(t, result, resetTime.String())
	})

	t.Run("Error message includes full timestamp", func(t *testing.T) {
		resetTime := time.Now().Add(30 * time.Second)
		err := &CircuitBreakerError{ResetAt: resetTime}

		result := err.Error()

		assert.NotEmpty(t, result)
		assert.Contains(t, result, "circuit breaker is open")
	})
}

// TestMakeCircuitBreakerError tests the constructor function
func TestMakeCircuitBreakerError(t *testing.T) {
	t.Run("creates CircuitBreakerError with correct reset time", func(t *testing.T) {
		resetTime := time.Date(2026, 1, 9, 13, 0, 0, 0, time.UTC)

		err := MakeCircuitBreakerError(resetTime)

		assert.NotNil(t, err)
		cbErr, ok := err.(*CircuitBreakerError)
		assert.True(t, ok, "should return *CircuitBreakerError type")
		assert.Equal(t, resetTime, cbErr.ResetAt)
	})

	t.Run("returns error interface", func(t *testing.T) {
		resetTime := time.Now().Add(1 * time.Minute)

		err := MakeCircuitBreakerError(resetTime)

		// Should be assignable to error interface
		var _ error = err
		assert.NotNil(t, err)
	})

	t.Run("created error can be type asserted", func(t *testing.T) {
		resetTime := time.Now().Add(45 * time.Second)

		err := MakeCircuitBreakerError(resetTime)

		cbErr, ok := err.(*CircuitBreakerError)
		assert.True(t, ok)
		assert.Equal(t, resetTime, cbErr.ResetAt)
	})
}

// TestAnyError tests the AnyError function
func TestAnyError(t *testing.T) {
	t.Run("returns Some for non-nil error", func(t *testing.T) {
		err := errors.New("test error")

		result := AnyError(err)

		assert.True(t, option.IsSome(result), "should return Some for non-nil error")
		value := option.GetOrElse(func() error { return nil })(result)
		assert.Equal(t, err, value)
	})

	t.Run("returns None for nil error", func(t *testing.T) {
		var err error = nil

		result := AnyError(err)

		assert.True(t, option.IsNone(result), "should return None for nil error")
	})

	t.Run("works with different error types", func(t *testing.T) {
		err1 := fmt.Errorf("wrapped: %w", errors.New("inner"))
		err2 := &CircuitBreakerError{ResetAt: time.Now()}

		result1 := AnyError(err1)
		result2 := AnyError(err2)

		assert.True(t, option.IsSome(result1))
		assert.True(t, option.IsSome(result2))
	})
}

// TestShouldOpenCircuit tests the shouldOpenCircuit function
func TestShouldOpenCircuit(t *testing.T) {
	t.Run("returns false for nil error", func(t *testing.T) {
		result := shouldOpenCircuit(nil)
		assert.False(t, result)
	})

	t.Run("HTTP 5xx errors should open circuit", func(t *testing.T) {
		testCases := []struct {
			name       string
			statusCode int
			expected   bool
		}{
			{"500 Internal Server Error", 500, true},
			{"501 Not Implemented", 501, true},
			{"502 Bad Gateway", 502, true},
			{"503 Service Unavailable", 503, true},
			{"504 Gateway Timeout", 504, true},
			{"599 Custom Server Error", 599, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testURL, _ := url.Parse("http://example.com")
				resp := &http.Response{
					StatusCode: tc.statusCode,
					Request:    &http.Request{URL: testURL},
					Body:       http.NoBody,
				}
				httpErr := FH.StatusCodeError(resp)

				result := shouldOpenCircuit(httpErr)

				assert.Equal(t, tc.expected, result)
			})
		}
	})

	t.Run("HTTP 4xx errors should NOT open circuit", func(t *testing.T) {
		testCases := []struct {
			name       string
			statusCode int
			expected   bool
		}{
			{"400 Bad Request", 400, false},
			{"401 Unauthorized", 401, false},
			{"403 Forbidden", 403, false},
			{"404 Not Found", 404, false},
			{"422 Unprocessable Entity", 422, false},
			{"429 Too Many Requests", 429, false},
			{"499 Custom Client Error", 499, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testURL, _ := url.Parse("http://example.com")
				resp := &http.Response{
					StatusCode: tc.statusCode,
					Request:    &http.Request{URL: testURL},
					Body:       http.NoBody,
				}
				httpErr := FH.StatusCodeError(resp)

				result := shouldOpenCircuit(httpErr)

				assert.Equal(t, tc.expected, result)
			})
		}
	})

	t.Run("HTTP 2xx and 3xx should NOT open circuit", func(t *testing.T) {
		testCases := []int{200, 201, 204, 301, 302, 304}

		for _, statusCode := range testCases {
			t.Run(fmt.Sprintf("Status %d", statusCode), func(t *testing.T) {
				testURL, _ := url.Parse("http://example.com")
				resp := &http.Response{
					StatusCode: statusCode,
					Request:    &http.Request{URL: testURL},
					Body:       http.NoBody,
				}
				httpErr := FH.StatusCodeError(resp)

				result := shouldOpenCircuit(httpErr)

				assert.False(t, result)
			})
		}
	})

	t.Run("network timeout errors should open circuit", func(t *testing.T) {
		opErr := &net.OpError{
			Op:  "dial",
			Err: &timeoutError{},
		}

		result := shouldOpenCircuit(opErr)

		assert.True(t, result)
	})

	t.Run("DNS errors should open circuit", func(t *testing.T) {
		dnsErr := &net.DNSError{
			Err:  "no such host",
			Name: "example.com",
		}

		result := shouldOpenCircuit(dnsErr)

		assert.True(t, result)
	})

	t.Run("URL timeout errors should open circuit", func(t *testing.T) {
		urlErr := &url.Error{
			Op:  "Get",
			URL: "http://example.com",
			Err: &timeoutError{},
		}

		result := shouldOpenCircuit(urlErr)

		assert.True(t, result)
	})

	t.Run("URL errors with nested network timeout should open circuit", func(t *testing.T) {
		urlErr := &url.Error{
			Op:  "Get",
			URL: "http://example.com",
			Err: &net.OpError{
				Op:  "dial",
				Err: &timeoutError{},
			},
		}

		result := shouldOpenCircuit(urlErr)

		assert.True(t, result)
	})

	t.Run("OpError with nil Err should open circuit", func(t *testing.T) {
		opErr := &net.OpError{
			Op:  "dial",
			Err: nil,
		}

		result := shouldOpenCircuit(opErr)

		assert.True(t, result)
	})

	t.Run("wrapped HTTP 5xx error should open circuit", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com")
		resp := &http.Response{
			StatusCode: 503,
			Request:    &http.Request{URL: testURL},
			Body:       http.NoBody,
		}
		httpErr := FH.StatusCodeError(resp)
		wrappedErr := fmt.Errorf("service error: %w", httpErr)

		result := shouldOpenCircuit(wrappedErr)

		assert.True(t, result)
	})

	t.Run("wrapped HTTP 4xx error should NOT open circuit", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com")
		resp := &http.Response{
			StatusCode: 404,
			Request:    &http.Request{URL: testURL},
			Body:       http.NoBody,
		}
		httpErr := FH.StatusCodeError(resp)
		wrappedErr := fmt.Errorf("not found: %w", httpErr)

		result := shouldOpenCircuit(wrappedErr)

		assert.False(t, result)
	})

	t.Run("generic application error should NOT open circuit", func(t *testing.T) {
		err := errors.New("validation failed")

		result := shouldOpenCircuit(err)

		assert.False(t, result)
	})
}

// TestIsInfrastructureError tests infrastructure error detection through shouldOpenCircuit
func TestIsInfrastructureError(t *testing.T) {
	t.Run("network timeout is infrastructure error", func(t *testing.T) {
		opErr := &net.OpError{Op: "dial", Err: &timeoutError{}}
		result := shouldOpenCircuit(opErr)
		assert.True(t, result)
	})

	t.Run("OpError with nil Err is infrastructure error", func(t *testing.T) {
		opErr := &net.OpError{Op: "dial", Err: nil}
		result := shouldOpenCircuit(opErr)
		assert.True(t, result)
	})

	t.Run("generic error returns false", func(t *testing.T) {
		err := errors.New("generic error")
		result := shouldOpenCircuit(err)
		assert.False(t, result)
	})

	t.Run("wrapped network timeout is detected", func(t *testing.T) {
		opErr := &net.OpError{Op: "dial", Err: &timeoutError{}}
		wrappedErr := fmt.Errorf("connection failed: %w", opErr)
		result := shouldOpenCircuit(wrappedErr)
		assert.True(t, result)
	})
}

// TestIsTLSError tests the isTLSError function
func TestIsTLSError(t *testing.T) {
	t.Run("certificate invalid error is TLS error", func(t *testing.T) {
		certErr := &x509.CertificateInvalidError{
			Reason: x509.Expired,
		}

		result := isTLSError(certErr)

		assert.True(t, result)
	})

	t.Run("unknown authority error is TLS error", func(t *testing.T) {
		authErr := &x509.UnknownAuthorityError{}

		result := isTLSError(authErr)

		assert.True(t, result)
	})

	t.Run("generic error is not TLS error", func(t *testing.T) {
		err := errors.New("generic error")

		result := isTLSError(err)

		assert.False(t, result)
	})

	t.Run("wrapped certificate error is detected", func(t *testing.T) {
		certErr := &x509.CertificateInvalidError{
			Reason: x509.Expired,
		}
		wrappedErr := fmt.Errorf("TLS handshake failed: %w", certErr)

		result := isTLSError(wrappedErr)

		assert.True(t, result)
	})

	t.Run("wrapped unknown authority error is detected", func(t *testing.T) {
		authErr := &x509.UnknownAuthorityError{}
		wrappedErr := fmt.Errorf("certificate verification failed: %w", authErr)

		result := isTLSError(wrappedErr)

		assert.True(t, result)
	})
}

// TestInfrastructureError tests the InfrastructureError variable
func TestInfrastructureError(t *testing.T) {
	t.Run("returns Some for infrastructure errors", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com")
		resp := &http.Response{
			StatusCode: 503,
			Request:    &http.Request{URL: testURL},
			Body:       http.NoBody,
		}
		httpErr := FH.StatusCodeError(resp)

		result := InfrastructureError(httpErr)

		assert.True(t, option.IsSome(result))
	})

	t.Run("returns None for non-infrastructure errors", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com")
		resp := &http.Response{
			StatusCode: 404,
			Request:    &http.Request{URL: testURL},
			Body:       http.NoBody,
		}
		httpErr := FH.StatusCodeError(resp)

		result := InfrastructureError(httpErr)

		assert.True(t, option.IsNone(result))
	})

	t.Run("returns None for nil error", func(t *testing.T) {
		result := InfrastructureError(nil)

		assert.True(t, option.IsNone(result))
	})

	t.Run("returns Some for network timeout", func(t *testing.T) {
		opErr := &net.OpError{
			Op:  "dial",
			Err: &timeoutError{},
		}

		result := InfrastructureError(opErr)

		assert.True(t, option.IsSome(result))
	})
}

// TestComplexErrorScenarios tests complex real-world error scenarios
func TestComplexErrorScenarios(t *testing.T) {
	t.Run("deeply nested URL error with HTTP 5xx", func(t *testing.T) {
		testURL, _ := url.Parse("http://api.example.com")
		resp := &http.Response{
			StatusCode: 502,
			Request:    &http.Request{URL: testURL},
			Body:       http.NoBody,
		}
		httpErr := FH.StatusCodeError(resp)
		urlErr := &url.Error{
			Op:  "Get",
			URL: "http://api.example.com",
			Err: httpErr,
		}
		wrappedErr := fmt.Errorf("API call failed: %w", urlErr)

		result := shouldOpenCircuit(wrappedErr)

		assert.True(t, result, "should detect HTTP 5xx through multiple layers")
	})

	t.Run("URL error with timeout nested in OpError", func(t *testing.T) {
		opErr := &net.OpError{
			Op:  "dial",
			Err: &timeoutError{},
		}
		urlErr := &url.Error{
			Op:  "Post",
			URL: "http://api.example.com",
			Err: opErr,
		}

		result := shouldOpenCircuit(urlErr)

		assert.True(t, result, "should detect timeout through URL error")
	})

	t.Run("multiple wrapped errors with infrastructure error at core", func(t *testing.T) {
		coreErr := &net.OpError{Op: "dial", Err: &timeoutError{}}
		layer1 := fmt.Errorf("connection attempt failed: %w", coreErr)
		layer2 := fmt.Errorf("retry exhausted: %w", layer1)
		layer3 := fmt.Errorf("service unavailable: %w", layer2)

		result := shouldOpenCircuit(layer3)

		assert.True(t, result, "should unwrap to find infrastructure error")
	})

	t.Run("OpError with nil Err should open circuit", func(t *testing.T) {
		opErr := &net.OpError{
			Op:  "dial",
			Err: nil,
		}

		result := shouldOpenCircuit(opErr)

		assert.True(t, result, "OpError with nil Err should be treated as infrastructure error")
	})

	t.Run("mixed error types - HTTP 4xx with network error", func(t *testing.T) {
		// This tests that we correctly identify the error type
		testURL, _ := url.Parse("http://example.com")
		resp := &http.Response{
			StatusCode: 400,
			Request:    &http.Request{URL: testURL},
			Body:       http.NoBody,
		}
		httpErr := FH.StatusCodeError(resp)

		result := shouldOpenCircuit(httpErr)

		assert.False(t, result, "HTTP 4xx should not open circuit even if wrapped")
	})
}

// Helper type for testing timeout errors
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }
