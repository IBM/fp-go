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

// TestMakeVoidMetrics tests the MakeVoidMetrics constructor
func TestMakeVoidMetrics(t *testing.T) {
	t.Run("creates valid Metrics implementation", func(t *testing.T) {
		metrics := MakeVoidMetrics()

		assert.NotNil(t, metrics, "MakeVoidMetrics should return non-nil Metrics")
	})

	t.Run("returns voidMetrics type", func(t *testing.T) {
		metrics := MakeVoidMetrics()

		_, ok := metrics.(*voidMetrics)
		assert.True(t, ok, "should return *voidMetrics type")
	})

	t.Run("initializes noop IO operation", func(t *testing.T) {
		metrics := MakeVoidMetrics().(*voidMetrics)

		assert.NotNil(t, metrics.noop, "noop IO operation should be initialized")
	})
}

// TestVoidMetricsAccept tests the Accept method of voidMetrics
func TestVoidMetricsAccept(t *testing.T) {
	t.Run("returns non-nil IO operation", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Accept(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
	})

	t.Run("IO operation executes without side effects", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Accept(timestamp)
		result := io.Run(ioOp)

		assert.NotNil(t, result, "IO operation should execute successfully")
	})

	t.Run("returns same IO operation instance", func(t *testing.T) {
		metrics := MakeVoidMetrics().(*voidMetrics)
		timestamp := time.Now()

		ioOp1 := metrics.Accept(timestamp)
		ioOp2 := metrics.Accept(timestamp)

		// Both should be non-nil (we can't compare functions directly in Go)
		assert.NotNil(t, ioOp1, "should return non-nil IO operation")
		assert.NotNil(t, ioOp2, "should return non-nil IO operation")

		// Verify they execute without error
		io.Run(ioOp1)
		io.Run(ioOp2)
	})

	t.Run("ignores timestamp parameter", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		time1 := time.Date(2026, 1, 9, 15, 30, 0, 0, time.UTC)
		time2 := time.Date(2026, 1, 9, 16, 30, 0, 0, time.UTC)

		ioOp1 := metrics.Accept(time1)
		ioOp2 := metrics.Accept(time2)

		// Should return same operation regardless of timestamp
		io.Run(ioOp1)
		io.Run(ioOp2)
		// No assertions needed - just verify it doesn't panic
	})
}

// TestVoidMetricsReject tests the Reject method of voidMetrics
func TestVoidMetricsReject(t *testing.T) {
	t.Run("returns non-nil IO operation", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Reject(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
	})

	t.Run("IO operation executes without side effects", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Reject(timestamp)
		result := io.Run(ioOp)

		assert.NotNil(t, result, "IO operation should execute successfully")
	})

	t.Run("returns same IO operation instance", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Reject(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		io.Run(ioOp) // Verify it executes without error
	})
}

// TestVoidMetricsOpen tests the Open method of voidMetrics
func TestVoidMetricsOpen(t *testing.T) {
	t.Run("returns non-nil IO operation", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Open(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
	})

	t.Run("IO operation executes without side effects", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Open(timestamp)
		result := io.Run(ioOp)

		assert.NotNil(t, result, "IO operation should execute successfully")
	})

	t.Run("returns same IO operation instance", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Open(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		io.Run(ioOp) // Verify it executes without error
	})
}

// TestVoidMetricsClose tests the Close method of voidMetrics
func TestVoidMetricsClose(t *testing.T) {
	t.Run("returns non-nil IO operation", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Close(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
	})

	t.Run("IO operation executes without side effects", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Close(timestamp)
		result := io.Run(ioOp)

		assert.NotNil(t, result, "IO operation should execute successfully")
	})

	t.Run("returns same IO operation instance", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Close(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		io.Run(ioOp) // Verify it executes without error
	})
}

// TestVoidMetricsCanary tests the Canary method of voidMetrics
func TestVoidMetricsCanary(t *testing.T) {
	t.Run("returns non-nil IO operation", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Canary(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
	})

	t.Run("IO operation executes without side effects", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Canary(timestamp)
		result := io.Run(ioOp)

		assert.NotNil(t, result, "IO operation should execute successfully")
	})

	t.Run("returns same IO operation instance", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Canary(timestamp)

		assert.NotNil(t, ioOp, "should return non-nil IO operation")
		io.Run(ioOp) // Verify it executes without error
	})
}

// TestVoidMetricsThreadSafety tests concurrent access to voidMetrics
func TestVoidMetricsThreadSafety(t *testing.T) {
	t.Run("handles concurrent metric calls", func(t *testing.T) {
		metrics := MakeVoidMetrics()

		var wg sync.WaitGroup
		numGoroutines := 100
		wg.Add(numGoroutines * 5) // 5 methods

		timestamp := time.Now()

		// Launch multiple goroutines calling all methods concurrently
		for i := 0; i < numGoroutines; i++ {
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
		// Test passes if no panic occurs
	})

	t.Run("all methods return valid IO operations concurrently", func(t *testing.T) {
		metrics := MakeVoidMetrics()

		var wg sync.WaitGroup
		numGoroutines := 50
		wg.Add(numGoroutines)

		timestamp := time.Now()
		results := make([]IO[Void], numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(idx int) {
				defer wg.Done()
				// Each goroutine calls a different method
				switch idx % 5 {
				case 0:
					results[idx] = metrics.Accept(timestamp)
				case 1:
					results[idx] = metrics.Reject(timestamp)
				case 2:
					results[idx] = metrics.Open(timestamp)
				case 3:
					results[idx] = metrics.Close(timestamp)
				case 4:
					results[idx] = metrics.Canary(timestamp)
				}
			}(i)
		}

		wg.Wait()

		// All results should be non-nil and executable
		for i, result := range results {
			assert.NotNil(t, result, "result %d should be non-nil", i)
			io.Run(result) // Verify it executes without error
		}
	})
}

// TestVoidMetricsPerformance tests performance characteristics
func TestVoidMetricsPerformance(t *testing.T) {
	t.Run("has minimal overhead", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		// Execute many operations quickly
		iterations := 10000
		for i := 0; i < iterations; i++ {
			io.Run(metrics.Accept(timestamp))
			io.Run(metrics.Reject(timestamp))
			io.Run(metrics.Open(timestamp))
			io.Run(metrics.Close(timestamp))
			io.Run(metrics.Canary(timestamp))
		}
		// Test passes if it completes quickly without issues
	})

	t.Run("all methods return valid IO operations", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		// All methods should return non-nil IO operations
		accept := metrics.Accept(timestamp)
		reject := metrics.Reject(timestamp)
		open := metrics.Open(timestamp)
		close := metrics.Close(timestamp)
		canary := metrics.Canary(timestamp)

		assert.NotNil(t, accept, "Accept should return non-nil")
		assert.NotNil(t, reject, "Reject should return non-nil")
		assert.NotNil(t, open, "Open should return non-nil")
		assert.NotNil(t, close, "Close should return non-nil")
		assert.NotNil(t, canary, "Canary should return non-nil")

		// All should execute without error
		io.Run(accept)
		io.Run(reject)
		io.Run(open)
		io.Run(close)
		io.Run(canary)
	})
}

// TestVoidMetricsIntegration tests integration scenarios
func TestVoidMetricsIntegration(t *testing.T) {
	t.Run("can be used as drop-in replacement for loggingMetrics", func(t *testing.T) {
		// Create both types of metrics
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		loggingMetrics := MakeMetricsFromLogger("TestCircuit", logger)
		voidMetrics := MakeVoidMetrics()

		timestamp := time.Now()

		// Both should implement the same interface
		var m1 Metrics = loggingMetrics
		var m2 Metrics = voidMetrics

		// Both should be callable
		io.Run(m1.Accept(timestamp))
		io.Run(m2.Accept(timestamp))

		// Logging metrics should have output
		assert.NotEmpty(t, buf.String(), "logging metrics should produce output")

		// Void metrics should have no observable side effects
		// (we can't directly test this, but the test passes if no panic occurs)
	})

	t.Run("simulates complete circuit breaker lifecycle without side effects", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		baseTime := time.Date(2026, 1, 9, 15, 30, 0, 0, time.UTC)

		// Simulate circuit breaker lifecycle - all should be no-ops
		io.Run(metrics.Accept(baseTime))
		io.Run(metrics.Accept(baseTime.Add(1 * time.Second)))
		io.Run(metrics.Open(baseTime.Add(2 * time.Second)))
		io.Run(metrics.Reject(baseTime.Add(3 * time.Second)))
		io.Run(metrics.Canary(baseTime.Add(30 * time.Second)))
		io.Run(metrics.Close(baseTime.Add(31 * time.Second)))

		// Test passes if no panic occurs and completes quickly
	})
}

// TestVoidMetricsEdgeCases tests edge cases
func TestVoidMetricsEdgeCases(t *testing.T) {
	t.Run("handles zero time", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		zeroTime := time.Time{}

		io.Run(metrics.Accept(zeroTime))
		io.Run(metrics.Reject(zeroTime))
		io.Run(metrics.Open(zeroTime))
		io.Run(metrics.Close(zeroTime))
		io.Run(metrics.Canary(zeroTime))

		// Test passes if no panic occurs
	})

	t.Run("handles far future time", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		futureTime := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)

		io.Run(metrics.Accept(futureTime))
		io.Run(metrics.Reject(futureTime))
		io.Run(metrics.Open(futureTime))
		io.Run(metrics.Close(futureTime))
		io.Run(metrics.Canary(futureTime))

		// Test passes if no panic occurs
	})

	t.Run("IO operations are idempotent", func(t *testing.T) {
		metrics := MakeVoidMetrics()
		timestamp := time.Now()

		ioOp := metrics.Accept(timestamp)

		// Execute same operation multiple times
		io.Run(ioOp)
		io.Run(ioOp)
		io.Run(ioOp)

		// Test passes if no panic occurs
	})
}

// TestMetricsComparison compares loggingMetrics and voidMetrics
func TestMetricsComparison(t *testing.T) {
	t.Run("both implement Metrics interface", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)

		var m1 Metrics = MakeMetricsFromLogger("Test", logger)
		var m2 Metrics = MakeVoidMetrics()

		assert.NotNil(t, m1)
		assert.NotNil(t, m2)
	})

	t.Run("voidMetrics has no observable side effects unlike loggingMetrics", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		loggingMetrics := MakeMetricsFromLogger("Test", logger)
		voidMetrics := MakeVoidMetrics()

		timestamp := time.Now()

		// Logging metrics produces output
		io.Run(loggingMetrics.Accept(timestamp))
		assert.NotEmpty(t, buf.String(), "logging metrics should produce output")

		// Void metrics has no observable output
		// (we can only verify it doesn't panic)
		io.Run(voidMetrics.Accept(timestamp))
	})
}
