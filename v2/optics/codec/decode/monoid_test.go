package decode

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	MO "github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestApplicativeMonoid(t *testing.T) {
	t.Run("with string monoid", func(t *testing.T) {
		m := ApplicativeMonoid[string](S.Monoid)

		t.Run("empty returns decoder that succeeds with empty string", func(t *testing.T) {
			empty := m.Empty()
			result := empty("any input")

			assert.Equal(t, validation.Of(""), result)
		})

		t.Run("concat combines successful decoders", func(t *testing.T) {
			decoder1 := Of[string]("Hello")
			decoder2 := Of[string](" World")

			combined := m.Concat(decoder1, decoder2)
			result := combined("input")

			assert.Equal(t, validation.Of("Hello World"), result)
		})

		t.Run("concat with failure returns failure", func(t *testing.T) {
			decoder1 := Of[string]("Hello")
			decoder2 := func(input string) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: input, Messsage: "decode failed"},
				})
			}

			combined := m.Concat(decoder1, decoder2)
			result := combined("input")

			assert.True(t, either.IsLeft(result))
			errors := either.MonadFold(result,
				F.Identity[Errors],
				func(string) Errors { return nil },
			)
			assert.Len(t, errors, 1)
			assert.Equal(t, "decode failed", errors[0].Messsage)
		})

		t.Run("concat accumulates all errors from both failures", func(t *testing.T) {
			decoder1 := func(input string) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: input, Messsage: "error 1"},
				})
			}
			decoder2 := func(input string) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: input, Messsage: "error 2"},
				})
			}

			combined := m.Concat(decoder1, decoder2)
			result := combined("input")

			assert.True(t, either.IsLeft(result))
			errors := either.MonadFold(result,
				F.Identity[Errors],
				func(string) Errors { return nil },
			)
			assert.Len(t, errors, 2)
			messages := []string{errors[0].Messsage, errors[1].Messsage}
			assert.Contains(t, messages, "error 1")
			assert.Contains(t, messages, "error 2")
		})

		t.Run("concat with empty preserves decoder", func(t *testing.T) {
			decoder := Of[string]("test")
			empty := m.Empty()

			result1 := m.Concat(decoder, empty)("input")
			result2 := m.Concat(empty, decoder)("input")

			val1 := either.MonadFold(result1,
				func(Errors) string { return "" },
				F.Identity[string],
			)
			val2 := either.MonadFold(result2,
				func(Errors) string { return "" },
				F.Identity[string],
			)

			assert.Equal(t, "test", val1)
			assert.Equal(t, "test", val2)
		})
	})

	t.Run("with int addition monoid", func(t *testing.T) {
		intMonoid := MO.MakeMonoid(
			func(a, b int) int { return a + b },
			0,
		)
		m := ApplicativeMonoid[string](intMonoid)

		t.Run("empty returns decoder with zero", func(t *testing.T) {
			empty := m.Empty()
			result := empty("input")

			value := either.MonadFold(result,
				func(Errors) int { return -1 },
				F.Identity[int],
			)
			assert.Equal(t, 0, value)
		})

		t.Run("concat adds decoded values", func(t *testing.T) {
			decoder1 := Of[string](10)
			decoder2 := Of[string](32)

			combined := m.Concat(decoder1, decoder2)
			result := combined("input")

			value := either.MonadFold(result,
				func(Errors) int { return 0 },
				F.Identity[int],
			)
			assert.Equal(t, 42, value)
		})

		t.Run("multiple concat operations", func(t *testing.T) {
			decoder1 := Of[string](1)
			decoder2 := Of[string](2)
			decoder3 := Of[string](3)
			decoder4 := Of[string](4)

			combined := m.Concat(m.Concat(m.Concat(decoder1, decoder2), decoder3), decoder4)
			result := combined("input")

			value := either.MonadFold(result,
				func(Errors) int { return 0 },
				F.Identity[int],
			)
			assert.Equal(t, 10, value)
		})
	})

	t.Run("with map input type", func(t *testing.T) {
		m := ApplicativeMonoid[map[string]any](S.Monoid)

		t.Run("combines decoders with different inputs", func(t *testing.T) {
			decoder1 := func(data map[string]any) Validation[string] {
				if name, ok := data["firstName"].(string); ok {
					return validation.Of(name)
				}
				return either.Left[string](validation.Errors{
					{Messsage: "missing firstName"},
				})
			}

			decoder2 := func(data map[string]any) Validation[string] {
				if name, ok := data["lastName"].(string); ok {
					return validation.Of(" " + name)
				}
				return either.Left[string](validation.Errors{
					{Messsage: "missing lastName"},
				})
			}

			combined := m.Concat(decoder1, decoder2)

			// Test success case
			result1 := combined(map[string]any{
				"firstName": "John",
				"lastName":  "Doe",
			})
			value1 := either.MonadFold(result1,
				func(Errors) string { return "" },
				F.Identity[string],
			)
			assert.Equal(t, "John Doe", value1)

			// Test failure case - both fields missing
			result2 := combined(map[string]any{})
			assert.True(t, either.IsLeft(result2))
			errors := either.MonadFold(result2,
				F.Identity[Errors],
				func(string) Errors { return nil },
			)
			assert.Len(t, errors, 2)
		})
	})
}

func TestMonoidLaws(t *testing.T) {
	t.Run("ApplicativeMonoid satisfies monoid laws", func(t *testing.T) {
		m := ApplicativeMonoid[string](S.Monoid)

		decoder1 := Of[string]("a")
		decoder2 := Of[string]("b")

		t.Run("left identity", func(t *testing.T) {
			// empty + a = a
			result := m.Concat(m.Empty(), decoder1)("input")
			value := either.MonadFold(result,
				func(Errors) string { return "" },
				F.Identity[string],
			)
			assert.Equal(t, "a", value)
		})

		t.Run("right identity", func(t *testing.T) {
			// a + empty = a
			result := m.Concat(decoder1, m.Empty())("input")
			value := either.MonadFold(result,
				func(Errors) string { return "" },
				F.Identity[string],
			)
			assert.Equal(t, "a", value)
		})

		t.Run("associativity", func(t *testing.T) {
			decoder3 := Of[string]("c")
			// (a + b) + c = a + (b + c)
			left := m.Concat(m.Concat(decoder1, decoder2), decoder3)("input")
			right := m.Concat(decoder1, m.Concat(decoder2, decoder3))("input")

			leftVal := either.MonadFold(left,
				func(Errors) string { return "" },
				F.Identity[string],
			)
			rightVal := either.MonadFold(right,
				func(Errors) string { return "" },
				F.Identity[string],
			)

			assert.Equal(t, "abc", leftVal)
			assert.Equal(t, "abc", rightVal)
		})
	})
}

func TestApplicativeMonoidWithFailures(t *testing.T) {
	m := ApplicativeMonoid[string](S.Monoid)

	t.Run("failure propagates through concat", func(t *testing.T) {
		decoder1 := Of[string]("a")
		decoder2 := func(input string) Validation[string] {
			return either.Left[string](validation.Errors{
				{Value: input, Messsage: "error"},
			})
		}
		decoder3 := Of[string]("c")

		combined := m.Concat(m.Concat(decoder1, decoder2), decoder3)
		result := combined("input")

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 1)
	})

	t.Run("multiple failures accumulate", func(t *testing.T) {
		decoder1 := func(input string) Validation[string] {
			return either.Left[string](validation.Errors{
				{Value: input, Messsage: "error 1"},
			})
		}
		decoder2 := func(input string) Validation[string] {
			return either.Left[string](validation.Errors{
				{Value: input, Messsage: "error 2"},
			})
		}
		decoder3 := func(input string) Validation[string] {
			return either.Left[string](validation.Errors{
				{Value: input, Messsage: "error 3"},
			})
		}

		combined := m.Concat(m.Concat(decoder1, decoder2), decoder3)
		result := combined("input")

		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 3)
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "error 1")
		assert.Contains(t, messages, "error 2")
		assert.Contains(t, messages, "error 3")
	})
}

func TestApplicativeMonoidEdgeCases(t *testing.T) {
	t.Run("with custom struct monoid", func(t *testing.T) {
		type Counter struct{ Count int }

		counterMonoid := MO.MakeMonoid(
			func(a, b Counter) Counter { return Counter{Count: a.Count + b.Count} },
			Counter{Count: 0},
		)

		m := ApplicativeMonoid[string](counterMonoid)

		decoder1 := Of[string](Counter{Count: 5})
		decoder2 := Of[string](Counter{Count: 10})

		combined := m.Concat(decoder1, decoder2)
		result := combined("input")

		value := either.MonadFold(result,
			func(Errors) Counter { return Counter{} },
			F.Identity[Counter],
		)
		assert.Equal(t, 15, value.Count)
	})

	t.Run("empty concat empty", func(t *testing.T) {
		m := ApplicativeMonoid[string](S.Monoid)

		combined := m.Concat(m.Empty(), m.Empty())
		result := combined("input")

		value := either.MonadFold(result,
			func(Errors) string { return "ERROR" },
			F.Identity[string],
		)
		assert.Equal(t, "", value)
	})

	t.Run("with different input types", func(t *testing.T) {
		intMonoid := MO.MakeMonoid(
			func(a, b int) int { return a + b },
			0,
		)
		m := ApplicativeMonoid[int](intMonoid)

		decoder1 := func(input int) Validation[int] {
			return validation.Of(input * 2)
		}
		decoder2 := func(input int) Validation[int] {
			return validation.Of(input + 10)
		}

		combined := m.Concat(decoder1, decoder2)
		result := combined(5)

		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		// (5 * 2) + (5 + 10) = 10 + 15 = 25
		assert.Equal(t, 25, value)
	})
}

func TestApplicativeMonoidRealWorldScenarios(t *testing.T) {
	t.Run("combining configuration from multiple sources", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}

		// Monoid that combines configs (last non-empty wins for strings, sum for ints)
		configMonoid := MO.MakeMonoid(
			func(a, b Config) Config {
				host := a.Host
				if b.Host != "" {
					host = b.Host
				}
				return Config{
					Host: host,
					Port: a.Port + b.Port,
				}
			},
			Config{Host: "", Port: 0},
		)

		m := ApplicativeMonoid[map[string]any](configMonoid)

		decoder1 := func(data map[string]any) Validation[Config] {
			if host, ok := data["host"].(string); ok {
				return validation.Of(Config{Host: host, Port: 0})
			}
			return either.Left[Config](validation.Errors{
				{Messsage: "missing host"},
			})
		}

		decoder2 := func(data map[string]any) Validation[Config] {
			if port, ok := data["port"].(int); ok {
				return validation.Of(Config{Host: "", Port: port})
			}
			return either.Left[Config](validation.Errors{
				{Messsage: "missing port"},
			})
		}

		combined := m.Concat(decoder1, decoder2)

		// Success case
		result := combined(map[string]any{
			"host": "localhost",
			"port": 8080,
		})

		config := either.MonadFold(result,
			func(Errors) Config { return Config{} },
			F.Identity[Config],
		)
		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, 8080, config.Port)
	})

	t.Run("aggregating validation results", func(t *testing.T) {
		intMonoid := MO.MakeMonoid(
			func(a, b int) int { return a + b },
			0,
		)
		m := ApplicativeMonoid[string](intMonoid)

		// Decoder that extracts and validates a number
		makeDecoder := func(value int, shouldFail bool) Decode[string, int] {
			return func(input string) Validation[int] {
				if shouldFail {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "validation failed"},
					})
				}
				return validation.Of(value)
			}
		}

		// All succeed - values are summed
		decoder1 := makeDecoder(10, false)
		decoder2 := makeDecoder(20, false)
		decoder3 := makeDecoder(12, false)

		combined := m.Concat(m.Concat(decoder1, decoder2), decoder3)
		result := combined("input")

		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)

		// Some fail - errors are accumulated
		decoder4 := makeDecoder(10, true)
		decoder5 := makeDecoder(20, true)

		combinedFail := m.Concat(decoder4, decoder5)
		resultFail := combinedFail("input")

		assert.True(t, either.IsLeft(resultFail))
		errors := either.MonadFold(resultFail,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2)
	})
}

// Made with Bob
