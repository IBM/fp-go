package validate

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	MO "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestAlternativeMonoid tests the AlternativeMonoid function
func TestAlternativeMonoid(t *testing.T) {
	t.Run("with string monoid", func(t *testing.T) {
		m := AlternativeMonoid[string, string](S.Monoid)

		t.Run("empty returns validator that succeeds with empty string", func(t *testing.T) {
			empty := m.Empty()
			result := empty("input")(nil)

			assert.Equal(t, validation.Of(""), result)
		})

		t.Run("concat combines successful validators using monoid", func(t *testing.T) {
			validator1 := Of[string, string]("Hello")
			validator2 := Of[string, string](" World")

			combined := m.Concat(validator1, validator2)
			result := combined("input")(nil)

			assert.Equal(t, validation.Of("Hello World"), result)
		})

		t.Run("concat uses second as fallback when first fails", func(t *testing.T) {
			failing := func(input string) Reader[Context, Validation[string]] {
				return func(ctx Context) Validation[string] {
					return either.Left[string](validation.Errors{
						{Value: input, Messsage: "first failed"},
					})
				}
			}
			succeeding := Of[string, string]("fallback")

			combined := m.Concat(failing, succeeding)
			result := combined("input")(nil)

			assert.Equal(t, validation.Of("fallback"), result)
		})

		t.Run("concat aggregates errors when both fail", func(t *testing.T) {
			failing1 := func(input string) Reader[Context, Validation[string]] {
				return func(ctx Context) Validation[string] {
					return either.Left[string](validation.Errors{
						{Value: input, Messsage: "error 1"},
					})
				}
			}
			failing2 := func(input string) Reader[Context, Validation[string]] {
				return func(ctx Context) Validation[string] {
					return either.Left[string](validation.Errors{
						{Value: input, Messsage: "error 2"},
					})
				}
			}

			combined := m.Concat(failing1, failing2)
			result := combined("input")(nil)

			assert.True(t, either.IsLeft(result))
			errors := either.MonadFold(result,
				F.Identity[Errors],
				func(string) Errors { return nil },
			)
			assert.GreaterOrEqual(t, len(errors), 2, "Should aggregate errors from both validators")

			messages := make([]string, len(errors))
			for i, err := range errors {
				messages[i] = err.Messsage
			}
			assert.Contains(t, messages, "error 1")
			assert.Contains(t, messages, "error 2")
		})

		t.Run("concat with empty preserves validator", func(t *testing.T) {
			validator := Of[string, string]("test")
			empty := m.Empty()

			result1 := m.Concat(validator, empty)("input")(nil)
			result2 := m.Concat(empty, validator)("input")(nil)

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
		m := AlternativeMonoid[string, int](intMonoid)

		t.Run("empty returns validator with zero", func(t *testing.T) {
			empty := m.Empty()
			result := empty("input")(nil)

			value := either.MonadFold(result,
				func(Errors) int { return -1 },
				F.Identity[int],
			)
			assert.Equal(t, 0, value)
		})

		t.Run("concat combines decoded values when both succeed", func(t *testing.T) {
			validator1 := Of[string, int](10)
			validator2 := Of[string, int](32)

			combined := m.Concat(validator1, validator2)
			result := combined("input")(nil)

			value := either.MonadFold(result,
				func(Errors) int { return 0 },
				F.Identity[int],
			)
			assert.Equal(t, 42, value)
		})

		t.Run("concat uses fallback when first fails", func(t *testing.T) {
			failing := func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "failed"},
					})
				}
			}
			succeeding := Of[string, int](42)

			combined := m.Concat(failing, succeeding)
			result := combined("input")(nil)

			value := either.MonadFold(result,
				func(Errors) int { return 0 },
				F.Identity[int],
			)
			assert.Equal(t, 42, value)
		})

		t.Run("multiple concat operations", func(t *testing.T) {
			validator1 := Of[string, int](1)
			validator2 := Of[string, int](2)
			validator3 := Of[string, int](3)
			validator4 := Of[string, int](4)

			combined := m.Concat(m.Concat(m.Concat(validator1, validator2), validator3), validator4)
			result := combined("input")(nil)

			value := either.MonadFold(result,
				func(Errors) int { return 0 },
				F.Identity[int],
			)
			assert.Equal(t, 10, value)
		})
	})

	t.Run("satisfies monoid laws", func(t *testing.T) {
		m := AlternativeMonoid[string, string](S.Monoid)

		validator1 := Of[string, string]("a")
		validator2 := Of[string, string]("b")
		validator3 := Of[string, string]("c")

		t.Run("left identity", func(t *testing.T) {
			result := m.Concat(m.Empty(), validator1)("input")(nil)
			value := either.MonadFold(result,
				func(Errors) string { return "" },
				F.Identity[string],
			)
			assert.Equal(t, "a", value)
		})

		t.Run("right identity", func(t *testing.T) {
			result := m.Concat(validator1, m.Empty())("input")(nil)
			value := either.MonadFold(result,
				func(Errors) string { return "" },
				F.Identity[string],
			)
			assert.Equal(t, "a", value)
		})

		t.Run("associativity", func(t *testing.T) {
			left := m.Concat(m.Concat(validator1, validator2), validator3)("input")(nil)
			right := m.Concat(validator1, m.Concat(validator2, validator3))("input")(nil)

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

// TestAltMonoid tests the AltMonoid function
func TestAltMonoid(t *testing.T) {
	t.Run("with default value as zero", func(t *testing.T) {
		m := AltMonoid(func() Validate[string, int] {
			return Of[string, int](0)
		})

		t.Run("empty returns the provided zero validator", func(t *testing.T) {
			empty := m.Empty()
			result := empty("input")(nil)

			assert.Equal(t, validation.Of(0), result)
		})

		t.Run("concat returns first validator when it succeeds", func(t *testing.T) {
			validator1 := Of[string, int](42)
			validator2 := Of[string, int](100)

			combined := m.Concat(validator1, validator2)
			result := combined("input")(nil)

			assert.Equal(t, validation.Of(42), result)
		})

		t.Run("concat uses second as fallback when first fails", func(t *testing.T) {
			failing := func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "failed"},
					})
				}
			}
			succeeding := Of[string, int](42)

			combined := m.Concat(failing, succeeding)
			result := combined("input")(nil)

			assert.Equal(t, validation.Of(42), result)
		})

		t.Run("concat aggregates errors when both fail", func(t *testing.T) {
			failing1 := func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 1"},
					})
				}
			}
			failing2 := func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 2"},
					})
				}
			}

			combined := m.Concat(failing1, failing2)
			result := combined("input")(nil)

			assert.True(t, either.IsLeft(result))
			errors := either.MonadFold(result,
				F.Identity[Errors],
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
	})

	t.Run("with failing zero", func(t *testing.T) {
		m := AltMonoid(func() Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Messsage: "no default available"},
					})
				}
			}
		})

		t.Run("empty returns the failing zero validator", func(t *testing.T) {
			empty := m.Empty()
			result := empty("input")(nil)

			assert.True(t, either.IsLeft(result))
		})

		t.Run("concat with all failures aggregates errors", func(t *testing.T) {
			failing1 := func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 1"},
					})
				}
			}
			failing2 := func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					return either.Left[int](validation.Errors{
						{Value: input, Messsage: "error 2"},
					})
				}
			}

			combined := m.Concat(failing1, failing2)
			result := combined("input")(nil)

			assert.True(t, either.IsLeft(result))
			errors := either.MonadFold(result,
				F.Identity[Errors],
				func(int) Errors { return nil },
			)
			assert.GreaterOrEqual(t, len(errors), 2, "Should aggregate errors")
		})
	})

	t.Run("chaining multiple fallbacks", func(t *testing.T) {
		m := AltMonoid(func() Validate[string, string] {
			return Of[string, string]("default")
		})

		primary := func(input string) Reader[Context, Validation[string]] {
			return func(ctx Context) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: input, Messsage: "primary failed"},
				})
			}
		}
		secondary := func(input string) Reader[Context, Validation[string]] {
			return func(ctx Context) Validation[string] {
				return either.Left[string](validation.Errors{
					{Value: input, Messsage: "secondary failed"},
				})
			}
		}
		tertiary := Of[string, string]("tertiary value")

		combined := m.Concat(m.Concat(primary, secondary), tertiary)
		result := combined("input")(nil)

		assert.Equal(t, validation.Of("tertiary value"), result)
	})

	t.Run("difference from AlternativeMonoid", func(t *testing.T) {
		// AltMonoid - first success wins
		altM := AltMonoid(func() Validate[string, int] {
			return Of[string, int](0)
		})

		// AlternativeMonoid - combines successes
		altMonoid := AlternativeMonoid[string, int](N.MonoidSum[int]())

		validator1 := Of[string, int](10)
		validator2 := Of[string, int](32)

		// AltMonoid: returns first success (10)
		result1 := altM.Concat(validator1, validator2)("input")(nil)
		value1 := either.MonadFold(result1,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 10, value1, "AltMonoid returns first success")

		// AlternativeMonoid: combines both successes (10 + 32 = 42)
		result2 := altMonoid.Concat(validator1, validator2)("input")(nil)
		value2 := either.MonadFold(result2,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value2, "AlternativeMonoid combines successes")
	})
}
