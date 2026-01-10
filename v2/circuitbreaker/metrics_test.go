package circuitbreaker

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/io"
	"github.com/stretchr/testify/assert"
)

// TestMakeMetricsFromLogger tests the MakeMetricsFromLogger constructor
func TestMakeMetricsFromLogger(t *testing.T) {
	t.Run("creates valid Metrics implementation", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)

		metrics := MakeMetricsFromLogger("TestCircuit", logger)

		assert.NotNil(t, metrics, "MakeMetricsFromLogger should return non-nil Metrics")
	})

	t.Run("returns loggingMetrics type", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)

		metrics := MakeMetricsFromLogger("TestCircuit", logger)

		_, ok := metrics.(*loggingMetrics)
		assert.True(t, ok, "should return *loggingMetrics type")
	})

	t.Run("stores name correctly", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		name := "MyCircuitBreaker"

		metrics := MakeMetricsFromLogger(name, logger).(*loggingMetrics)

		assert.Equal(t, name, metrics.name, "name should be stored correctly")
	})

	t.Run("stores logger correctly", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)

		metrics := MakeMetricsFromLogger("TestCircuit", logger).(*loggingMetrics)

		assert.Equal(t, logger, metrics.logger, "logger should be stored correctly")
	})
}

// TestLoggingMetricsAccept tests the Accept method
func TestLoggingMetricsAccept(t *testing.T) {
	t.Run("logs accept event with correct format", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Date(2026, 1, 9, 15, 30, 45, 0, time.UTC)

		io.Run(metrics.Accept(timestamp))

		output := buf.String()
		assert.Contains(t, output, "Accept:", "should contain Accept prefix")
		assert.Contains(t, output, "TestCircuit", "should contain circuit name")
		assert.Contains(t, output, timestamp.String(), "should contain timestamp")
	})

	t.Run("returns IO[Void] that can be executed", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		ioOp := metrics.Accept(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		result := io.Run(ioOp)
		assert.NotNil(t, result, "IO operation should execute successfully")
	})

	t.Run("logs multiple accept events", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		time1 := time.Date(2026, 1, 9, 15, 30, 0, 0, time.UTC)
		time2 := time.Date(2026, 1, 9, 15, 31, 0, 0, time.UTC)

		io.Run(metrics.Accept(time1))
		io.Run(metrics.Accept(time2))

		output := buf.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Len(t, lines, 2, "should have 2 log lines")
		assert.Contains(t, lines[0], time1.String())
		assert.Contains(t, lines[1], time2.String())
	})
}

// TestLoggingMetricsReject tests the Reject method
func TestLoggingMetricsReject(t *testing.T) {
	t.Run("logs reject event with correct format", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Date(2026, 1, 9, 15, 30, 45, 0, time.UTC)

		io.Run(metrics.Reject(timestamp))

		output := buf.String()
		assert.Contains(t, output, "Reject:", "should contain Reject prefix")
		assert.Contains(t, output, "TestCircuit", "should contain circuit name")
		assert.Contains(t, output, timestamp.String(), "should contain timestamp")
	})

	t.Run("returns IO[Void] that can be executed", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		ioOp := metrics.Reject(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		result := io.Run(ioOp)
		assert.NotNil(t, result, "IO operation should execute successfully")
	})
}

// TestLoggingMetricsOpen tests the Open method
func TestLoggingMetricsOpen(t *testing.T) {
	t.Run("logs open event with correct format", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Date(2026, 1, 9, 15, 30, 45, 0, time.UTC)

		io.Run(metrics.Open(timestamp))

		output := buf.String()
		assert.Contains(t, output, "Open:", "should contain Open prefix")
		assert.Contains(t, output, "TestCircuit", "should contain circuit name")
		assert.Contains(t, output, timestamp.String(), "should contain timestamp")
	})

	t.Run("returns IO[Void] that can be executed", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		ioOp := metrics.Open(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		result := io.Run(ioOp)
		assert.NotNil(t, result, "IO operation should execute successfully")
	})
}

// TestLoggingMetricsClose tests the Close method
func TestLoggingMetricsClose(t *testing.T) {
	t.Run("logs close event with correct format", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Date(2026, 1, 9, 15, 30, 45, 0, time.UTC)

		io.Run(metrics.Close(timestamp))

		output := buf.String()
		assert.Contains(t, output, "Close:", "should contain Close prefix")
		assert.Contains(t, output, "TestCircuit", "should contain circuit name")
		assert.Contains(t, output, timestamp.String(), "should contain timestamp")
	})

	t.Run("returns IO[Void] that can be executed", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		ioOp := metrics.Close(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		result := io.Run(ioOp)
		assert.NotNil(t, result, "IO operation should execute successfully")
	})
}

// TestLoggingMetricsCanary tests the Canary method
func TestLoggingMetricsCanary(t *testing.T) {
	t.Run("logs canary event with correct format", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Date(2026, 1, 9, 15, 30, 45, 0, time.UTC)

		io.Run(metrics.Canary(timestamp))

		output := buf.String()
		assert.Contains(t, output, "Canary:", "should contain Canary prefix")
		assert.Contains(t, output, "TestCircuit", "should contain circuit name")
		assert.Contains(t, output, timestamp.String(), "should contain timestamp")
	})

	t.Run("returns IO[Void] that can be executed", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		ioOp := metrics.Canary(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		result := io.Run(ioOp)
		assert.NotNil(t, result, "IO operation should execute successfully")
	})
}

// TestLoggingMetricsDoLog tests the doLog helper method
func TestLoggingMetricsDoLog(t *testing.T) {
	t.Run("formats log message correctly", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := &loggingMetrics{name: "TestCircuit", logger: logger}
		timestamp := time.Date(2026, 1, 9, 15, 30, 45, 0, time.UTC)

		io.Run(metrics.doLog("CustomEvent", timestamp))

		output := buf.String()
		assert.Contains(t, output, "CustomEvent:", "should contain custom prefix")
		assert.Contains(t, output, "TestCircuit", "should contain circuit name")
		assert.Contains(t, output, timestamp.String(), "should contain timestamp")
	})

	t.Run("handles different prefixes", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := &loggingMetrics{name: "TestCircuit", logger: logger}
		timestamp := time.Now()

		prefixes := []string{"Accept", "Reject", "Open", "Close", "Canary", "Custom"}
		for _, prefix := range prefixes {
			buf.Reset()
			io.Run(metrics.doLog(prefix, timestamp))
			output := buf.String()
			assert.Contains(t, output, prefix+":", "should contain prefix: "+prefix)
		}
	})
}

// TestMetricsIntegration tests integration scenarios
func TestMetricsIntegration(t *testing.T) {
	t.Run("logs complete circuit breaker lifecycle", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("APICircuit", logger)
		baseTime := time.Date(2026, 1, 9, 15, 30, 0, 0, time.UTC)

		// Simulate circuit breaker lifecycle
		io.Run(metrics.Accept(baseTime))                       // Request accepted
		io.Run(metrics.Accept(baseTime.Add(1 * time.Second)))  // Another request
		io.Run(metrics.Open(baseTime.Add(2 * time.Second)))    // Circuit opens
		io.Run(metrics.Reject(baseTime.Add(3 * time.Second)))  // Request rejected
		io.Run(metrics.Canary(baseTime.Add(30 * time.Second))) // Canary attempt
		io.Run(metrics.Close(baseTime.Add(31 * time.Second)))  // Circuit closes

		output := buf.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Len(t, lines, 6, "should have 6 log lines")

		assert.Contains(t, lines[0], "Accept:")
		assert.Contains(t, lines[1], "Accept:")
		assert.Contains(t, lines[2], "Open:")
		assert.Contains(t, lines[3], "Reject:")
		assert.Contains(t, lines[4], "Canary:")
		assert.Contains(t, lines[5], "Close:")
	})

	t.Run("distinguishes between multiple circuit breakers", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics1 := MakeMetricsFromLogger("Circuit1", logger)
		metrics2 := MakeMetricsFromLogger("Circuit2", logger)
		timestamp := time.Now()

		io.Run(metrics1.Accept(timestamp))
		io.Run(metrics2.Accept(timestamp))

		output := buf.String()
		assert.Contains(t, output, "Circuit1", "should contain first circuit name")
		assert.Contains(t, output, "Circuit2", "should contain second circuit name")
	})
}

// TestMetricsThreadSafety tests concurrent access to metrics
func TestMetricsThreadSafety(t *testing.T) {
	t.Run("handles concurrent metric recording", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("ConcurrentCircuit", logger)

		var wg sync.WaitGroup
		numGoroutines := 100
		wg.Add(numGoroutines)

		// Launch multiple goroutines recording metrics concurrently
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				timestamp := time.Now()
				io.Run(metrics.Accept(timestamp))
			}(i)
		}

		wg.Wait()

		output := buf.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Len(t, lines, numGoroutines, "should have logged all events")
	})

	t.Run("handles concurrent different event types", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("ConcurrentCircuit", logger)

		var wg sync.WaitGroup
		numIterations := 20
		wg.Add(numIterations * 5) // 5 event types

		timestamp := time.Now()

		for i := 0; i < numIterations; i++ {
			go func() {
				defer wg.Done()
				io.Run(metrics.Accept(timestamp))
			}()
			go func() {
				defer wg.Done()
				io.Run(metrics.Reject(timestamp))
			}()
			go func() {
				defer wg.Done()
				io.Run(metrics.Open(timestamp))
			}()
			go func() {
				defer wg.Done()
				io.Run(metrics.Close(timestamp))
			}()
			go func() {
				defer wg.Done()
				io.Run(metrics.Canary(timestamp))
			}()
		}

		wg.Wait()

		output := buf.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Len(t, lines, numIterations*5, "should have logged all events")
	})
}

// TestMetricsEdgeCases tests edge cases and special scenarios
func TestMetricsEdgeCases(t *testing.T) {
	t.Run("handles empty circuit breaker name", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("", logger)
		timestamp := time.Now()

		io.Run(metrics.Accept(timestamp))

		output := buf.String()
		assert.NotEmpty(t, output, "should still log even with empty name")
	})

	t.Run("handles very long circuit breaker name", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		longName := strings.Repeat("VeryLongCircuitBreakerName", 100)
		metrics := MakeMetricsFromLogger(longName, logger)
		timestamp := time.Now()

		io.Run(metrics.Accept(timestamp))

		output := buf.String()
		assert.Contains(t, output, longName, "should handle long names")
	})

	t.Run("handles special characters in name", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		specialName := "Circuit-Breaker_123!@#$%^&*()"
		metrics := MakeMetricsFromLogger(specialName, logger)
		timestamp := time.Now()

		io.Run(metrics.Accept(timestamp))

		output := buf.String()
		assert.Contains(t, output, specialName, "should handle special characters")
	})

	t.Run("handles zero time", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		zeroTime := time.Time{}

		io.Run(metrics.Accept(zeroTime))

		output := buf.String()
		assert.NotEmpty(t, output, "should handle zero time")
		assert.Contains(t, output, "Accept:")
	})

	t.Run("handles far future time", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		futureTime := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)

		io.Run(metrics.Accept(futureTime))

		output := buf.String()
		assert.NotEmpty(t, output, "should handle far future time")
		assert.Contains(t, output, "9999")
	})
}

// TestMetricsWithCustomLogger tests metrics with different logger configurations
func TestMetricsWithCustomLogger(t *testing.T) {
	t.Run("works with logger with custom prefix", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "[CB] ", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		io.Run(metrics.Accept(timestamp))

		output := buf.String()
		assert.Contains(t, output, "[CB]", "should include custom prefix")
		assert.Contains(t, output, "Accept:")
	})

	t.Run("works with logger with flags", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", log.Ldate|log.Ltime)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		io.Run(metrics.Accept(timestamp))

		output := buf.String()
		assert.NotEmpty(t, output, "should log with flags")
		assert.Contains(t, output, "Accept:")
	})
}

// TestMetricsIOOperations tests IO operation behavior
func TestMetricsIOOperations(t *testing.T) {
	t.Run("IO operations are lazy", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		// Create IO operation but don't execute it
		_ = metrics.Accept(timestamp)

		// Buffer should be empty because IO wasn't executed
		assert.Empty(t, buf.String(), "IO operation should be lazy")
	})

	t.Run("IO operations execute when run", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		ioOp := metrics.Accept(timestamp)
		io.Run(ioOp)

		assert.NotEmpty(t, buf.String(), "IO operation should execute when run")
	})

	t.Run("same IO operation can be executed multiple times", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		metrics := MakeMetricsFromLogger("TestCircuit", logger)
		timestamp := time.Now()

		ioOp := metrics.Accept(timestamp)
		io.Run(ioOp)
		io.Run(ioOp)
		io.Run(ioOp)

		output := buf.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Len(t, lines, 3, "should execute multiple times")
	})
}

// Made with Bob
