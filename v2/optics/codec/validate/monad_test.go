package validate

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
)

// TestMonadChainLeft tests the MonadChainLeft function
func TestMonadChainLeft(t *testing.T) {
	t.Run("transforms failures while preserving successes", func(t *testing.T) {
		// Create a failing validator
		failingValidator := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "validation failed"},
				})
			}
		}

		// Handler that recovers from specific errors
		handler := func(errs Errors) Validate[string, int] {
			for _, err := range errs {
				if err.Messsage == "validation failed" {
					return Of[string, int](0) // recover with default
				}
			}
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](errs)
				}
			}
		}

		validator := MonadChainLeft(failingValidator, handler)
		res := validator("input")(nil)

		assert.Equal(t, validation.Of(0), res, "Should recover from failure")
	})

	t.Run("preserves success values unchanged", func(t *testing.T) {
		successValidator := Of[string, int](42)

		handler := func(errs Errors) Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Messsage: "should not be called"},
					})
				}
			}
		}

		validator := MonadChainLeft(successValidator, handler)
		res := validator("input")(nil)

		assert.Equal(t, validation.Of(42), res, "Success should pass through unchanged")
	})

	t.Run("aggregates errors when transformation also fails", func(t *testing.T) {
		failingValidator := func(input string) Reader[Context, Validation[string]] {
			return func(ctx Context) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: input, Messsage: "original error"},
				})
			}
		}

		handler := func(errs Errors) Validate[string, string] {
			return func(input string) Reader[Context, Validation[string]] {
				return func(ctx Context) Validation[string] {
					return either.Left[string](validation.Errors{
						{Messsage: "additional error"},
					})
				}
			}
		}

		validator := MonadChainLeft(failingValidator, handler)
		res := validator("input")(nil)

		assert.True(t, either.IsLeft(res))
		errors := either.MonadFold(res,
			reader.Ask[Errors](),
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should aggregate both errors")

		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "original error")
		assert.Contains(t, messages, "additional error")
	})

	t.Run("adds context to errors", func(t *testing.T) {
		failingValidator := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "invalid format"},
				})
			}
		}

		addContext := func(errs Errors) Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{
							Context:  validation.Context{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
							Messsage: "failed to validate user age",
						},
					})
				}
			}
		}

		validator := MonadChainLeft(failingValidator, addContext)
		res := validator("abc")(nil)

		assert.True(t, either.IsLeft(res))
		errors := either.MonadFold(res,
			reader.Ask[Errors](),
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should have both original and context errors")
	})

	t.Run("works with different input types", func(t *testing.T) {
		type Config struct {
			Port int
		}

		failingValidator := func(cfg Config) Reader[Context, Validation[string]] {
			return func(ctx Context) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: cfg.Port, Messsage: "invalid port"},
				})
			}
		}

		handler := func(errs Errors) Validate[Config, string] {
			return Of[Config, string]("default-value")
		}

		validator := MonadChainLeft(failingValidator, handler)
		res := validator(Config{Port: 9999})(nil)

		assert.Equal(t, validation.Of("default-value"), res)
	})

	t.Run("handler can access original input", func(t *testing.T) {
		failingValidator := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "parse failed"},
				})
			}
		}

		handler := func(errs Errors) Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					// Handler can use the original input to make decisions
					if input == "special" {
						return validation.Of(999)
					}
					return validation.Of(0)
				}
			}
		}

		validator := MonadChainLeft(failingValidator, handler)

		res1 := validator("special")(nil)
		assert.Equal(t, validation.Of(999), res1)

		res2 := validator("other")(nil)
		assert.Equal(t, validation.Of(0), res2)
	})

	t.Run("is equivalent to ChainLeft", func(t *testing.T) {
		failingValidator := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "error"},
				})
			}
		}

		handler := func(errs Errors) Validate[string, int] {
			return Of[string, int](42)
		}

		// MonadChainLeft - direct application
		result1 := MonadChainLeft(failingValidator, handler)("input")(nil)

		// ChainLeft - curried for pipelines
		result2 := ChainLeft(handler)(failingValidator)("input")(nil)

		assert.Equal(t, result1, result2, "MonadChainLeft and ChainLeft should produce identical results")
	})

	t.Run("chains multiple error transformations", func(t *testing.T) {
		failingValidator := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "error1"},
				})
			}
		}

		handler1 := func(errs Errors) Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Messsage: "error2"},
					})
				}
			}
		}

		handler2 := func(errs Errors) Validate[string, int] {
			// Check if we can recover
			for _, err := range errs {
				if err.Messsage == "error1" {
					return Of[string, int](100) // recover
				}
			}
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](errs)
				}
			}
		}

		// Chain handlers
		validator := MonadChainLeft(MonadChainLeft(failingValidator, handler1), handler2)
		res := validator("input")(nil)

		// Should recover because error1 is present
		assert.Equal(t, validation.Of(100), res)
	})

	t.Run("does not call handler on success", func(t *testing.T) {
		successValidator := Of[string, int](42)
		handlerCalled := false

		handler := func(errs Errors) Validate[string, int] {
			handlerCalled = true
			return Of[string, int](0)
		}

		validator := MonadChainLeft(successValidator, handler)
		res := validator("input")(nil)

		assert.Equal(t, validation.Of(42), res)
		assert.False(t, handlerCalled, "Handler should not be called on success")
	})
}

// TestMonadAlt tests the MonadAlt function
func TestMonadAlt(t *testing.T) {
	t.Run("returns first validator when it succeeds", func(t *testing.T) {
		validator1 := Of[string, int](42)
		validator2 := func() Validate[string, int] {
			return Of[string, int](100)
		}

		result := MonadAlt(validator1, validator2)("input")(nil)
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("returns second validator when first fails", func(t *testing.T) {
		failing := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "first failed"},
				})
			}
		}
		fallback := func() Validate[string, int] {
			return Of[string, int](42)
		}

		result := MonadAlt(failing, fallback)("input")(nil)
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("aggregates errors when both fail", func(t *testing.T) {
		failing1 := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "error 1"},
				})
			}
		}
		failing2 := func() Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 2"},
					})
				}
			}
		}

		result := MonadAlt(failing1, failing2)("input")(nil)
		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			reader.Ask[Errors](),
			func(int) Errors { return nil },
		)
		assert.GreaterOrEqual(t, len(errors), 2, "Should aggregate errors from both validators")

		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "error 1", "Should contain error from first validator")
		assert.Contains(t, messages, "error 2", "Should contain error from second validator")
	})

	t.Run("does not evaluate second validator when first succeeds", func(t *testing.T) {
		validator1 := Of[string, int](42)
		evaluated := false
		validator2 := func() Validate[string, int] {
			evaluated = true
			return Of[string, int](100)
		}

		result := MonadAlt(validator1, validator2)("input")(nil)
		assert.Equal(t, validation.Of(42), result)
		assert.False(t, evaluated, "Second validator should not be evaluated")
	})

	t.Run("works with different types", func(t *testing.T) {
		failing := func(input string) Reader[Context, Validation[string]] {
			return func(ctx Context) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: input, Messsage: "failed"},
				})
			}
		}
		fallback := func() Validate[string, string] {
			return Of[string, string]("fallback")
		}

		result := MonadAlt(failing, fallback)("input")(nil)
		assert.Equal(t, validation.Of("fallback"), result)
	})

	t.Run("chains multiple alternatives", func(t *testing.T) {
		failing1 := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "error 1"},
				})
			}
		}
		failing2 := func() Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 2"},
					})
				}
			}
		}
		succeeding := func() Validate[string, int] {
			return Of[string, int](42)
		}

		// Chain: try failing1, then failing2, then succeeding
		result := MonadAlt(MonadAlt(failing1, failing2), succeeding)("input")(nil)
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("works with complex input types", func(t *testing.T) {
		type Config struct {
			Port int
		}

		failing := func(cfg Config) Reader[Context, Validation[string]] {
			return func(ctx Context) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: cfg.Port, Messsage: "invalid port"},
				})
			}
		}
		fallback := func() Validate[Config, string] {
			return Of[Config, string]("default")
		}

		result := MonadAlt(failing, fallback)(Config{Port: 9999})(nil)
		assert.Equal(t, validation.Of("default"), result)
	})

	t.Run("preserves error context", func(t *testing.T) {
		failing1 := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{
						Value:    input,
						Messsage: "parse error",
						Context:  validation.Context{{Key: "field", Type: "int"}},
					},
				})
			}
		}
		failing2 := func() Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{
							Value:    input,
							Messsage: "validation error",
							Context:  validation.Context{{Key: "value", Type: "int"}},
						},
					})
				}
			}
		}

		result := MonadAlt(failing1, failing2)("abc")(nil)
		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			reader.Ask[Errors](),
			func(int) Errors { return nil },
		)
		assert.GreaterOrEqual(t, len(errors), 2, "Should have errors from both validators")

		// Verify that errors with context are present
		hasParseError := false
		hasValidationError := false
		for _, err := range errors {
			if err.Messsage == "parse error" {
				hasParseError = true
				assert.NotNil(t, err.Context)
			}
			if err.Messsage == "validation error" {
				hasValidationError = true
				assert.NotNil(t, err.Context)
			}
		}
		assert.True(t, hasParseError, "Should have parse error")
		assert.True(t, hasValidationError, "Should have validation error")
	})
}

// TestAlt tests the Alt function
func TestAlt(t *testing.T) {
	t.Run("returns first validator when it succeeds", func(t *testing.T) {
		validator1 := Of[string, int](42)
		validator2 := func() Validate[string, int] {
			return Of[string, int](100)
		}

		withAlt := Alt(validator2)
		result := withAlt(validator1)("input")(nil)
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("returns second validator when first fails", func(t *testing.T) {
		failing := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "first failed"},
				})
			}
		}
		fallback := func() Validate[string, int] {
			return Of[string, int](42)
		}

		withAlt := Alt(fallback)
		result := withAlt(failing)("input")(nil)
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("aggregates errors when both fail", func(t *testing.T) {
		failing1 := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "error 1"},
				})
			}
		}
		failing2 := func() Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 2"},
					})
				}
			}
		}

		withAlt := Alt(failing2)
		result := withAlt(failing1)("input")(nil)
		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			reader.Ask[Errors](),
			func(int) Errors { return nil },
		)
		assert.GreaterOrEqual(t, len(errors), 2, "Should aggregate errors from both validators")

		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "error 1")
		assert.Contains(t, messages, "error 2")
	})

	t.Run("does not evaluate second validator when first succeeds", func(t *testing.T) {
		validator1 := Of[string, int](42)
		evaluated := false
		validator2 := func() Validate[string, int] {
			evaluated = true
			return Of[string, int](100)
		}

		withAlt := Alt(validator2)
		result := withAlt(validator1)("input")(nil)
		assert.Equal(t, validation.Of(42), result)
		assert.False(t, evaluated, "Second validator should not be evaluated")
	})

	t.Run("can be used in pipelines", func(t *testing.T) {
		failing1 := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "error 1"},
				})
			}
		}
		failing2 := func() Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 2"},
					})
				}
			}
		}
		succeeding := func() Validate[string, int] {
			return Of[string, int](42)
		}

		// Use F.Pipe to chain alternatives
		validator := F.Pipe2(
			failing1,
			Alt(failing2),
			Alt(succeeding),
		)

		result := validator("input")(nil)
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("is equivalent to MonadAlt", func(t *testing.T) {
		failing := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "error"},
				})
			}
		}
		fallback := func() Validate[string, int] {
			return Of[string, int](42)
		}

		// Alt - curried for pipelines
		result1 := Alt(fallback)(failing)("input")(nil)

		// MonadAlt - direct application
		result2 := MonadAlt(failing, fallback)("input")(nil)

		assert.Equal(t, result1, result2, "Alt and MonadAlt should produce identical results")
	})
}

// TestMonadAltAndAltEquivalence tests that MonadAlt and Alt are equivalent
func TestMonadAltAndAltEquivalence(t *testing.T) {
	t.Run("both produce same results for success", func(t *testing.T) {
		validator1 := Of[string, int](42)
		validator2 := func() Validate[string, int] {
			return Of[string, int](100)
		}

		resultMonadAlt := MonadAlt(validator1, validator2)("input")(nil)
		resultAlt := Alt(validator2)(validator1)("input")(nil)

		assert.Equal(t, resultMonadAlt, resultAlt)
	})

	t.Run("both produce same results for fallback", func(t *testing.T) {
		failing := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "failed"},
				})
			}
		}
		fallback := func() Validate[string, int] {
			return Of[string, int](42)
		}

		resultMonadAlt := MonadAlt(failing, fallback)("input")(nil)
		resultAlt := Alt(fallback)(failing)("input")(nil)

		assert.Equal(t, resultMonadAlt, resultAlt)
	})

	t.Run("both produce same results for error aggregation", func(t *testing.T) {
		failing1 := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return either.Left[int](validation.Errors{
					{Value: input, Messsage: "error 1"},
				})
			}
		}
		failing2 := func() Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 2"},
					})
				}
			}
		}

		resultMonadAlt := MonadAlt(failing1, failing2)("input")(nil)
		resultAlt := Alt(failing2)(failing1)("input")(nil)

		// Both should fail
		assert.True(t, either.IsLeft(resultMonadAlt))
		assert.True(t, either.IsLeft(resultAlt))

		// Both should have same errors
		errorsMonadAlt := either.MonadFold(resultMonadAlt,
			reader.Ask[Errors](),
			func(int) Errors { return nil },
		)
		errorsAlt := either.MonadFold(resultAlt,
			reader.Ask[Errors](),
			func(int) Errors { return nil },
		)

		assert.Equal(t, len(errorsMonadAlt), len(errorsAlt))
	})
}
