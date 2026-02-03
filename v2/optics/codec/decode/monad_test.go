package decode

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestOf tests the Of function
func TestOf(t *testing.T) {
	t.Run("creates decoder that always succeeds", func(t *testing.T) {
		decoder := Of[string](42)
		res := decoder("any input")

		assert.Equal(t, validation.Of(42), res)
	})

	t.Run("works with different input types", func(t *testing.T) {
		decoder := Of[int]("hello")
		res := decoder(123)

		assert.Equal(t, validation.Of("hello"), res)
	})

	t.Run("works with complex types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		person := Person{Name: "Alice", Age: 30}
		decoder := Of[string](person)
		res := decoder("input")

		assert.Equal(t, validation.Of(person), res)
	})

	t.Run("ignores input value", func(t *testing.T) {
		decoder := Of[string](100)

		res1 := decoder("input1")
		res2 := decoder("input2")

		assert.Equal(t, res1, res2)
		assert.Equal(t, validation.Of(100), res1)
	})
}

// TestMonadChain tests the MonadChain function
func TestMonadChain(t *testing.T) {
	t.Run("chains successful decoders", func(t *testing.T) {
		decoder1 := Of[string](42)
		decoder2 := MonadChain(decoder1, func(n int) Decode[string, string] {
			return Of[string](fmt.Sprintf("Number: %d", n))
		})

		res := decoder2("input")
		assert.Equal(t, validation.Of("Number: 42"), res)
	})

	t.Run("chains multiple operations", func(t *testing.T) {
		decoder1 := Of[string](10)
		decoder2 := MonadChain(decoder1, func(n int) Decode[string, int] {
			return Of[string](n * 2)
		})
		decoder3 := MonadChain(decoder2, func(n int) Decode[string, string] {
			return Of[string](fmt.Sprintf("Result: %d", n))
		})

		res := decoder3("input")
		assert.Equal(t, validation.Of("Result: 20"), res)
	})

	t.Run("propagates validation errors", func(t *testing.T) {
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "decode failed"},
			})
		}

		decoder1 := failingDecoder
		decoder2 := MonadChain(decoder1, func(n int) Decode[string, string] {
			return Of[string](fmt.Sprintf("Number: %d", n))
		})

		res := decoder2("input")
		assert.True(t, either.IsLeft(res))
	})

	t.Run("short-circuits on first error", func(t *testing.T) {
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "first error"},
			})
		}

		chainCalled := false
		decoder := MonadChain(failingDecoder, func(n int) Decode[string, string] {
			chainCalled = true
			return Of[string]("should not be called")
		})

		res := decoder("input")
		assert.True(t, either.IsLeft(res))
		assert.False(t, chainCalled, "Chain function should not be called on error")
	})
}

// TestChain tests the Chain function
func TestChain(t *testing.T) {
	t.Run("creates chainable operator", func(t *testing.T) {
		chainOp := Chain(func(n int) Decode[string, string] {
			return Of[string](fmt.Sprintf("Number: %d", n))
		})

		decoder := chainOp(Of[string](42))
		res := decoder("input")

		assert.Equal(t, validation.Of("Number: 42"), res)
	})

	t.Run("can be composed", func(t *testing.T) {
		double := Chain(func(n int) Decode[string, int] {
			return Of[string](n * 2)
		})

		toString := Chain(func(n int) Decode[string, string] {
			return Of[string](fmt.Sprintf("Value: %d", n))
		})

		decoder := toString(double(Of[string](21)))
		res := decoder("input")

		assert.Equal(t, validation.Of("Value: 42"), res)
	})
}

// TestMonadMap tests the MonadMap function
func TestMonadMap(t *testing.T) {
	t.Run("maps successful decoder", func(t *testing.T) {
		decoder := Of[string](42)
		mapped := MonadMap(decoder, S.Format[int]("Number: %d"))

		res := mapped("input")
		assert.Equal(t, validation.Of("Number: 42"), res)
	})

	t.Run("transforms value type", func(t *testing.T) {
		decoder := Of[string]("hello")
		mapped := MonadMap(decoder, S.Size)

		res := mapped("input")
		assert.Equal(t, validation.Of(5), res)
	})

	t.Run("preserves validation errors", func(t *testing.T) {
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "decode failed"},
			})
		}

		mapped := MonadMap(failingDecoder, S.Format[int]("Number: %d"))

		res := mapped("input")
		assert.True(t, either.IsLeft(res))
	})

	t.Run("does not call function on error", func(t *testing.T) {
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "error"},
			})
		}

		mapCalled := false
		mapped := MonadMap(failingDecoder, func(n int) string {
			mapCalled = true
			return "should not be called"
		})

		res := mapped("input")
		assert.True(t, either.IsLeft(res))
		assert.False(t, mapCalled, "Map function should not be called on error")
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		decoder := Of[string](10)
		mapped1 := MonadMap(decoder, N.Mul(2))
		mapped2 := MonadMap(mapped1, N.Add(5))
		mapped3 := MonadMap(mapped2, S.Format[int]("Result: %d"))

		res := mapped3("input")
		assert.Equal(t, validation.Of("Result: 25"), res)
	})
}

// TestMap tests the Map function
func TestMap(t *testing.T) {
	t.Run("creates mappable operator", func(t *testing.T) {
		mapOp := Map[string](S.Format[int]("Number: %d"))

		decoder := mapOp(Of[string](42))
		res := decoder("input")

		assert.Equal(t, validation.Of("Number: 42"), res)
	})

	t.Run("can be composed", func(t *testing.T) {
		double := Map[string](N.Mul(2))
		toString := Map[string](S.Format[int]("Value: %d"))

		decoder := toString(double(Of[string](21)))
		res := decoder("input")

		assert.Equal(t, validation.Of("Value: 42"), res)
	})
}

// TestMonadAp tests the MonadAp function
func TestMonadAp(t *testing.T) {
	t.Run("applies function decoder to value decoder", func(t *testing.T) {
		decoderFn := Of[string](S.Format[int]("Number: %d"))
		decoderVal := Of[string](42)

		res := MonadAp(decoderFn, decoderVal)("input")
		assert.Equal(t, validation.Of("Number: 42"), res)
	})

	t.Run("works with different transformations", func(t *testing.T) {
		decoderFn := Of[string](N.Mul(2))
		decoderVal := Of[string](21)

		res := MonadAp(decoderFn, decoderVal)("input")
		assert.Equal(t, validation.Of(42), res)
	})

	t.Run("propagates function decoder error", func(t *testing.T) {
		failingFnDecoder := func(input string) Validation[func(int) string] {
			return either.Left[func(int) string](validation.Errors{
				{Value: input, Messsage: "function decode failed"},
			})
		}
		decoderVal := Of[string](42)

		res := MonadAp(failingFnDecoder, decoderVal)("input")
		assert.True(t, either.IsLeft(res))
	})

	t.Run("propagates value decoder error", func(t *testing.T) {
		decoderFn := Of[string](S.Format[int]("Number: %d"))
		failingValDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "value decode failed"},
			})
		}

		res := MonadAp(decoderFn, failingValDecoder)("input")
		assert.True(t, either.IsLeft(res))
	})

	t.Run("combines multiple values", func(t *testing.T) {
		// Create a function that takes two arguments
		decoderFn := Of[string](N.Add[int])
		decoderVal1 := Of[string](10)
		decoderVal2 := Of[string](32)

		// Apply first value
		partial := MonadAp(decoderFn, decoderVal1)
		// Apply second value
		result := MonadAp(partial, decoderVal2)

		res := result("input")
		assert.Equal(t, validation.Of(42), res)
	})
}

// TestAp tests the Ap function
func TestAp(t *testing.T) {
	t.Run("creates applicable operator", func(t *testing.T) {
		decoderVal := Of[string](42)
		apOp := Ap[string](decoderVal)

		decoderFn := Of[string](S.Format[int]("Number: %d"))

		res := apOp(decoderFn)("input")
		assert.Equal(t, validation.Of("Number: 42"), res)
	})

	t.Run("can be composed", func(t *testing.T) {
		val1 := Of[string](10)
		val2 := Of[string](32)

		apOp1 := Ap[func(int) int](val1)
		apOp2 := Ap[int](val2)

		fnDecoder := Of[string](N.Add[int])

		result := apOp2(apOp1(fnDecoder))
		res := result("input")

		assert.Equal(t, validation.Of(42), res)
	})
}

// TestMonadLaws tests that the monad operations satisfy monad laws
func TestMonadLaws(t *testing.T) {
	t.Run("left identity: Of(a) >>= f === f(a)", func(t *testing.T) {
		a := 42
		f := func(n int) Decode[string, string] {
			return Of[string](fmt.Sprintf("Number: %d", n))
		}

		left := MonadChain(Of[string](a), f)
		right := f(a)

		input := "test"
		assert.Equal(t, right(input), left(input))
	})

	t.Run("right identity: m >>= Of === m", func(t *testing.T) {
		m := Of[string](42)

		left := MonadChain(m, func(a int) Decode[string, int] {
			return Of[string](a)
		})

		input := "test"
		assert.Equal(t, m(input), left(input))
	})

	t.Run("associativity: (m >>= f) >>= g === m >>= (\\x -> f(x) >>= g)", func(t *testing.T) {
		m := Of[string](10)
		f := func(n int) Decode[string, int] {
			return Of[string](n * 2)
		}
		g := func(n int) Decode[string, string] {
			return Of[string](fmt.Sprintf("Result: %d", n))
		}

		// (m >>= f) >>= g
		left := MonadChain(MonadChain(m, f), g)

		// m >>= (\x -> f(x) >>= g)
		right := MonadChain(m, func(x int) Decode[string, string] {
			return MonadChain(f(x), g)
		})

		input := "test"
		assert.Equal(t, right(input), left(input))
	})
}

// TestFunctorLaws tests that the functor operations satisfy functor laws
func TestFunctorLaws(t *testing.T) {
	t.Run("identity: map(id) === id", func(t *testing.T) {
		decoder := Of[string](42)
		mapped := MonadMap(decoder, func(a int) int { return a })

		input := "test"
		assert.Equal(t, decoder(input), mapped(input))
	})

	t.Run("composition: map(f . g) === map(f) . map(g)", func(t *testing.T) {
		decoder := Of[string](10)
		f := N.Mul(2)
		g := N.Add(5)

		// map(f . g)
		left := MonadMap(decoder, func(n int) int {
			return f(g(n))
		})

		// map(f) . map(g)
		right := MonadMap(MonadMap(decoder, g), f)

		input := "test"
		assert.Equal(t, right(input), left(input))
	})
}

// TestChainLeft tests the ChainLeft function
func TestChainLeft(t *testing.T) {
	t.Run("transforms failures while preserving successes", func(t *testing.T) {
		// Create a failing decoder
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "decode failed"},
			})
		}

		// Handler that recovers from specific errors
		handler := ChainLeft(func(errs Errors) Decode[string, int] {
			for _, err := range errs {
				if err.Messsage == "decode failed" {
					return Of[string](0) // recover with default
				}
			}
			return func(input string) Validation[int] {
				return either.Left[int](errs)
			}
		})

		decoder := handler(failingDecoder)
		res := decoder("input")

		assert.Equal(t, validation.Of(0), res, "Should recover from failure")
	})

	t.Run("preserves success values unchanged", func(t *testing.T) {
		successDecoder := Of[string](42)

		handler := ChainLeft(func(errs Errors) Decode[string, int] {
			return func(input string) Validation[int] {
				return either.Left[int](validation.Errors{
					{Messsage: "should not be called"},
				})
			}
		})

		decoder := handler(successDecoder)
		res := decoder("input")

		assert.Equal(t, validation.Of(42), res, "Success should pass through unchanged")
	})

	t.Run("aggregates errors when transformation also fails", func(t *testing.T) {
		failingDecoder := func(input string) Validation[string] {
			return either.Left[string](validation.Errors{
				{Value: input, Messsage: "original error"},
			})
		}

		handler := ChainLeft(func(errs Errors) Decode[string, string] {
			return func(input string) Validation[string] {
				return either.Left[string](validation.Errors{
					{Messsage: "additional error"},
				})
			}
		})

		decoder := handler(failingDecoder)
		res := decoder("input")

		assert.True(t, either.IsLeft(res))
		errors := either.MonadFold(res,
			func(e Errors) Errors { return e },
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
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "invalid format"},
			})
		}

		addContext := ChainLeft(func(errs Errors) Decode[string, int] {
			return func(input string) Validation[int] {
				return either.Left[int](validation.Errors{
					{
						Context:  validation.Context{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
						Messsage: "failed to decode user age",
					},
				})
			}
		})

		decoder := addContext(failingDecoder)
		res := decoder("abc")

		assert.True(t, either.IsLeft(res))
		errors := either.MonadFold(res,
			func(e Errors) Errors { return e },
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should have both original and context errors")
	})

	t.Run("can be composed in pipeline", func(t *testing.T) {
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "error1"},
			})
		}

		handler1 := ChainLeft(func(errs Errors) Decode[string, int] {
			return func(input string) Validation[int] {
				return either.Left[int](validation.Errors{
					{Messsage: "error2"},
				})
			}
		})

		handler2 := ChainLeft(func(errs Errors) Decode[string, int] {
			// Check if we can recover
			for _, err := range errs {
				if err.Messsage == "error1" {
					return Of[string](100) // recover
				}
			}
			return func(input string) Validation[int] {
				return either.Left[int](errs)
			}
		})

		// Compose handlers
		decoder := handler2(handler1(failingDecoder))
		res := decoder("input")

		// Should recover because error1 is present
		assert.Equal(t, validation.Of(100), res)
	})

	t.Run("works with different input types", func(t *testing.T) {
		type Config struct {
			Port int
		}

		failingDecoder := func(cfg Config) Validation[string] {
			return either.Left[string](validation.Errors{
				{Value: cfg.Port, Messsage: "invalid port"},
			})
		}

		handler := ChainLeft(func(errs Errors) Decode[Config, string] {
			return Of[Config]("default-value")
		})

		decoder := handler(failingDecoder)
		res := decoder(Config{Port: 9999})

		assert.Equal(t, validation.Of("default-value"), res)
	})
}

// TestOrElse tests the OrElse function
func TestOrElse(t *testing.T) {
	t.Run("OrElse is equivalent to ChainLeft - Success case", func(t *testing.T) {
		successDecoder := Of[string](42)

		handler := func(errs Errors) Decode[string, int] {
			return func(input string) Validation[int] {
				return either.Left[int](validation.Errors{
					{Messsage: "should not be called"},
				})
			}
		}

		// Test with OrElse
		resultOrElse := OrElse(handler)(successDecoder)("input")
		// Test with ChainLeft
		resultChainLeft := ChainLeft(handler)(successDecoder)("input")

		assert.Equal(t, resultChainLeft, resultOrElse, "OrElse and ChainLeft should produce identical results for Success")
		assert.Equal(t, validation.Of(42), resultOrElse)
	})

	t.Run("OrElse is equivalent to ChainLeft - Failure recovery", func(t *testing.T) {
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "not found"},
			})
		}

		handler := func(errs Errors) Decode[string, int] {
			for _, err := range errs {
				if err.Messsage == "not found" {
					return Of[string](0) // recover with default
				}
			}
			return func(input string) Validation[int] {
				return either.Left[int](errs)
			}
		}

		// Test with OrElse
		resultOrElse := OrElse(handler)(failingDecoder)("input")
		// Test with ChainLeft
		resultChainLeft := ChainLeft(handler)(failingDecoder)("input")

		assert.Equal(t, resultChainLeft, resultOrElse, "OrElse and ChainLeft should produce identical results for recovery")
		assert.Equal(t, validation.Of(0), resultOrElse)
	})

	t.Run("OrElse is equivalent to ChainLeft - Error aggregation", func(t *testing.T) {
		failingDecoder := func(input string) Validation[string] {
			return either.Left[string](validation.Errors{
				{Value: input, Messsage: "original error"},
			})
		}

		handler := func(errs Errors) Decode[string, string] {
			return func(input string) Validation[string] {
				return either.Left[string](validation.Errors{
					{Messsage: "additional error"},
				})
			}
		}

		// Test with OrElse
		resultOrElse := OrElse(handler)(failingDecoder)("input")
		// Test with ChainLeft
		resultChainLeft := ChainLeft(handler)(failingDecoder)("input")

		assert.Equal(t, resultChainLeft, resultOrElse, "OrElse and ChainLeft should produce identical results for error aggregation")

		// Verify both aggregate errors
		assert.True(t, either.IsLeft(resultOrElse))
		errors := either.MonadFold(resultOrElse,
			func(e Errors) Errors { return e },
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

	t.Run("OrElse semantic meaning - fallback decoder", func(t *testing.T) {
		// OrElse provides a semantic name for fallback/alternative decoding
		// It reads naturally: "try this decoder, or else try this alternative"

		primaryDecoder := func(input string) Validation[int] {
			if input == "valid" {
				return validation.Of(42)
			}
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "primary decode failed"},
			})
		}

		// Use OrElse to provide a fallback: if decoding fails, use default value
		withDefault := OrElse(func(errs Errors) Decode[string, int] {
			return Of[string](0) // default to 0 if decoding fails
		})

		decoder := withDefault(primaryDecoder)

		// Test success case
		resSuccess := decoder("valid")
		assert.Equal(t, validation.Of(42), resSuccess, "Should use primary decoder on success")

		// Test fallback case
		resFallback := decoder("invalid")
		assert.Equal(t, validation.Of(0), resFallback, "OrElse provides fallback value")
	})

	t.Run("OrElse in pipeline composition", func(t *testing.T) {
		failingDecoder := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "database error"},
			})
		}

		addContext := OrElse(func(errs Errors) Decode[string, int] {
			return func(input string) Validation[int] {
				return either.Left[int](validation.Errors{
					{Messsage: "context added"},
				})
			}
		})

		recoverFromNotFound := OrElse(func(errs Errors) Decode[string, int] {
			for _, err := range errs {
				if err.Messsage == "not found" {
					return Of[string](0)
				}
			}
			return func(input string) Validation[int] {
				return either.Left[int](errs)
			}
		})

		// Test error aggregation in pipeline
		decoder1 := recoverFromNotFound(addContext(failingDecoder))
		res1 := decoder1("input")

		assert.True(t, either.IsLeft(res1))
		errors := either.MonadFold(res1,
			func(e Errors) Errors { return e },
			func(int) Errors { return nil },
		)
		// Errors accumulate through the pipeline
		assert.Greater(t, len(errors), 1, "Should aggregate errors from pipeline")

		// Test recovery in pipeline
		failingDecoder2 := func(input string) Validation[int] {
			return either.Left[int](validation.Errors{
				{Value: input, Messsage: "not found"},
			})
		}

		decoder2 := recoverFromNotFound(addContext(failingDecoder2))
		res2 := decoder2("input")

		assert.Equal(t, validation.Of(0), res2, "Should recover from 'not found' error")
	})

	t.Run("OrElse vs ChainLeft - identical behavior verification", func(t *testing.T) {
		// Create various test scenarios
		scenarios := []struct {
			name    string
			decoder Decode[string, int]
			handler func(Errors) Decode[string, int]
		}{
			{
				name:    "Success value",
				decoder: Of[string](42),
				handler: func(errs Errors) Decode[string, int] {
					return func(input string) Validation[int] {
						return either.Left[int](validation.Errors{{Messsage: "error"}})
					}
				},
			},
			{
				name: "Failure with recovery",
				decoder: func(input string) Validation[int] {
					return either.Left[int](validation.Errors{{Messsage: "error"}})
				},
				handler: func(errs Errors) Decode[string, int] {
					return Of[string](0)
				},
			},
			{
				name: "Failure with error transformation",
				decoder: func(input string) Validation[int] {
					return either.Left[int](validation.Errors{{Messsage: "error1"}})
				},
				handler: func(errs Errors) Decode[string, int] {
					return func(input string) Validation[int] {
						return either.Left[int](validation.Errors{{Messsage: "error2"}})
					}
				},
			},
			{
				name: "Multiple errors aggregation",
				decoder: func(input string) Validation[int] {
					return either.Left[int](validation.Errors{
						{Messsage: "error1"},
						{Messsage: "error2"},
					})
				},
				handler: func(errs Errors) Decode[string, int] {
					return func(input string) Validation[int] {
						return either.Left[int](validation.Errors{
							{Messsage: "error3"},
							{Messsage: "error4"},
						})
					}
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.name, func(t *testing.T) {
				resultOrElse := OrElse(scenario.handler)(scenario.decoder)("test-input")
				resultChainLeft := ChainLeft(scenario.handler)(scenario.decoder)("test-input")

				assert.Equal(t, resultChainLeft, resultOrElse,
					"OrElse and ChainLeft must produce identical results for: %s", scenario.name)
			})
		}
	})
}
