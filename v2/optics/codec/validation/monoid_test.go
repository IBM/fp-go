package validation

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	MO "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestErrorsMonoid(t *testing.T) {
	m := ErrorsMonoid()

	t.Run("empty returns empty array", func(t *testing.T) {
		empty := m.Empty()
		assert.NotNil(t, empty)
		assert.Len(t, empty, 0)
	})

	t.Run("concat combines error arrays", func(t *testing.T) {
		errs1 := Errors{
			&ValidationError{Messsage: "error 1"},
			&ValidationError{Messsage: "error 2"},
		}
		errs2 := Errors{
			&ValidationError{Messsage: "error 3"},
		}

		result := m.Concat(errs1, errs2)

		assert.Len(t, result, 3)
		assert.Equal(t, "error 1", result[0].Messsage)
		assert.Equal(t, "error 2", result[1].Messsage)
		assert.Equal(t, "error 3", result[2].Messsage)
	})

	t.Run("concat with empty preserves errors", func(t *testing.T) {
		errs := Errors{
			&ValidationError{Messsage: "error"},
		}
		empty := m.Empty()

		result1 := m.Concat(errs, empty)
		result2 := m.Concat(empty, errs)

		assert.Equal(t, errs, result1)
		assert.Equal(t, errs, result2)
	})

	t.Run("concat is associative", func(t *testing.T) {
		errs1 := Errors{&ValidationError{Messsage: "1"}}
		errs2 := Errors{&ValidationError{Messsage: "2"}}
		errs3 := Errors{&ValidationError{Messsage: "3"}}

		// (a + b) + c
		left := m.Concat(m.Concat(errs1, errs2), errs3)
		// a + (b + c)
		right := m.Concat(errs1, m.Concat(errs2, errs3))

		assert.Len(t, left, 3)
		assert.Len(t, right, 3)
		for i := 0; i < 3; i++ {
			assert.Equal(t, left[i].Messsage, right[i].Messsage)
		}
	})
}

func TestApplicativeMonoid(t *testing.T) {
	t.Run("with string monoid", func(t *testing.T) {
		m := ApplicativeMonoid(S.Monoid)

		t.Run("empty returns successful validation with empty string", func(t *testing.T) {
			empty := m.Empty()

			assert.Equal(t, Success(""), empty)
		})

		t.Run("concat combines successful validations", func(t *testing.T) {
			v1 := Success("Hello")
			v2 := Success(" World")

			result := m.Concat(v1, v2)

			assert.Equal(t, Success("Hello World"), result)
		})

		t.Run("concat with failure returns failure", func(t *testing.T) {
			v1 := Success("Hello")
			v2 := Failures[string](Errors{
				&ValidationError{Messsage: "error"},
			})

			result := m.Concat(v1, v2)

			assert.True(t, either.IsLeft(result))
			errors := either.MonadFold(result,
				F.Identity[Errors],
				func(string) Errors { return nil },
			)
			assert.Len(t, errors, 1)
			assert.Equal(t, "error", errors[0].Messsage)
		})

		t.Run("concat accumulates all errors from both failures", func(t *testing.T) {
			v1 := Failures[string](Errors{
				&ValidationError{Messsage: "error 1"},
			})
			v2 := Failures[string](Errors{
				&ValidationError{Messsage: "error 2"},
			})

			result := m.Concat(v1, v2)

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

		t.Run("concat with empty preserves validation", func(t *testing.T) {
			v := Success("test")
			empty := m.Empty()

			result1 := m.Concat(v, empty)
			result2 := m.Concat(empty, v)

			assert.Equal(t, Of("test"), result1)
			assert.Equal(t, Of("test"), result2)
		})
	})

	t.Run("with int addition monoid", func(t *testing.T) {
		intMonoid := MO.MakeMonoid(
			func(a, b int) int { return a + b },
			0,
		)

		m := ApplicativeMonoid(intMonoid)

		t.Run("empty returns zero", func(t *testing.T) {
			empty := m.Empty()

			assert.Equal(t, Of(0), empty)
		})

		t.Run("concat adds values", func(t *testing.T) {
			v1 := Success(10)
			v2 := Success(32)

			result := m.Concat(v1, v2)

			assert.Equal(t, Of(42), result)
		})

		t.Run("multiple concat operations", func(t *testing.T) {
			v1 := Success(1)
			v2 := Success(2)
			v3 := Success(3)
			v4 := Success(4)

			result := m.Concat(m.Concat(m.Concat(v1, v2), v3), v4)

			assert.Equal(t, Of(10), result)
		})
	})
}

func TestMonoidLaws(t *testing.T) {
	t.Run("ErrorsMonoid satisfies monoid laws", func(t *testing.T) {
		m := ErrorsMonoid()

		errs1 := Errors{&ValidationError{Messsage: "1"}}
		errs2 := Errors{&ValidationError{Messsage: "2"}}

		t.Run("left identity", func(t *testing.T) {
			// empty + a = a
			result := m.Concat(m.Empty(), errs1)
			assert.Equal(t, errs1, result)
		})

		t.Run("right identity", func(t *testing.T) {
			// a + empty = a
			result := m.Concat(errs1, m.Empty())
			assert.Equal(t, errs1, result)
		})

		t.Run("associativity", func(t *testing.T) {
			errs3 := Errors{&ValidationError{Messsage: "3"}}
			// (a + b) + c = a + (b + c)
			left := m.Concat(m.Concat(errs1, errs2), errs3)
			right := m.Concat(errs1, m.Concat(errs2, errs3))

			assert.Len(t, left, 3)
			assert.Len(t, right, 3)
			for i := 0; i < 3; i++ {
				assert.Equal(t, left[i].Messsage, right[i].Messsage)
			}
		})
	})

	t.Run("ApplicativeMonoid satisfies monoid laws", func(t *testing.T) {
		m := ApplicativeMonoid(S.Monoid)

		v1 := Success("a")
		v2 := Success("b")

		t.Run("left identity", func(t *testing.T) {
			// empty + a = a
			result := m.Concat(m.Empty(), v1)
			assert.Equal(t, Of("a"), result)
		})

		t.Run("right identity", func(t *testing.T) {
			// a + empty = a
			result := m.Concat(v1, m.Empty())
			assert.Equal(t, Of("a"), result)
		})

		t.Run("associativity", func(t *testing.T) {
			v3 := Success("c")
			// (a + b) + c = a + (b + c)
			left := m.Concat(m.Concat(v1, v2), v3)
			right := m.Concat(v1, m.Concat(v2, v3))

			assert.Equal(t, Of("abc"), left)
			assert.Equal(t, Of("abc"), right)
		})
	})
}

func TestApplicativeMonoidWithFailures(t *testing.T) {
	m := ApplicativeMonoid(S.Monoid)

	t.Run("failure propagates through concat", func(t *testing.T) {
		v1 := Success("a")
		v2 := Failures[string](Errors{&ValidationError{Messsage: "error"}})
		v3 := Success("c")

		result := m.Concat(m.Concat(v1, v2), v3)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 1)
	})

	t.Run("multiple failures accumulate", func(t *testing.T) {
		v1 := Failures[string](Errors{&ValidationError{Messsage: "error 1"}})
		v2 := Failures[string](Errors{&ValidationError{Messsage: "error 2"}})
		v3 := Failures[string](Errors{&ValidationError{Messsage: "error 3"}})

		result := m.Concat(m.Concat(v1, v2), v3)

		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 3)
	})
}

func TestApplicativeMonoidEdgeCases(t *testing.T) {
	t.Run("with custom struct monoid", func(t *testing.T) {
		type Counter struct{ Count int }

		counterMonoid := MO.MakeMonoid(
			func(a, b Counter) Counter { return Counter{Count: a.Count + b.Count} },
			Counter{Count: 0},
		)

		m := ApplicativeMonoid(counterMonoid)

		v1 := Success(Counter{Count: 5})
		v2 := Success(Counter{Count: 10})

		result := m.Concat(v1, v2)

		assert.Equal(t, Of(Counter{Count: 15}), result)
	})

	t.Run("empty concat empty", func(t *testing.T) {
		m := ApplicativeMonoid(S.Monoid)

		result := m.Concat(m.Empty(), m.Empty())

		assert.Equal(t, Of(""), result)
	})
}
