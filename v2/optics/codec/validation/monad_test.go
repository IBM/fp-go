package validation

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	t.Run("creates successful validation", func(t *testing.T) {
		result := Of(42)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("works with different types", func(t *testing.T) {
		strResult := Of("hello")
		assert.True(t, either.IsRight(strResult))

		boolResult := Of(true)
		assert.True(t, either.IsRight(boolResult))

		type Custom struct{ Value int }
		customResult := Of(Custom{Value: 100})
		assert.True(t, either.IsRight(customResult))
	})

	t.Run("is equivalent to Success", func(t *testing.T) {
		value := 42
		ofResult := Of(value)
		successResult := Success(value)

		assert.Equal(t, ofResult, successResult)
	})
}

func TestMap(t *testing.T) {
	t.Run("transforms successful validation", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		result := Map(double)(Of(21))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("preserves failure", func(t *testing.T) {
		errs := Errors{&ValidationError{Messsage: "error"}}
		failure := Failures[int](errs)

		double := func(x int) int { return x * 2 }
		result := Map(double)(failure)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "error", errors[0].Messsage)
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		add10 := func(x int) int { return x + 10 }
		double := func(x int) int { return x * 2 }
		toString := func(x int) string { return fmt.Sprintf("%d", x) }

		result := F.Pipe3(
			Of(5),
			Map(add10),
			Map(double),
			Map(toString),
		)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "30", value)
	})

	t.Run("type transformation", func(t *testing.T) {
		length := func(s string) int { return len(s) }
		result := Map(length)(Of("hello"))

		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 5, value)
	})
}

func TestAp(t *testing.T) {
	t.Run("applies function to value when both succeed", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		funcValidation := Of(double)
		valueValidation := Of(21)

		result := Ap[int, int](valueValidation)(funcValidation)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("accumulates errors when value fails", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		funcValidation := Of(double)
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := Ap[int, int](valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "value error", errors[0].Messsage)
	})

	t.Run("accumulates errors when function fails", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := Of(21)

		result := Ap[int, int](valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "function error", errors[0].Messsage)
	})

	t.Run("accumulates all errors when both fail", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := Ap[int, int](valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2)
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "function error")
		assert.Contains(t, messages, "value error")
	})

	t.Run("applies with string transformation", func(t *testing.T) {
		toUpper := func(s string) string { return fmt.Sprintf("UPPER:%s", s) }
		funcValidation := Of(toUpper)
		valueValidation := Of("hello")

		result := Ap[string, string](valueValidation)(funcValidation)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "UPPER:hello", value)
	})

	t.Run("accumulates multiple validation errors from different sources", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error 1"},
			&ValidationError{Messsage: "function error 2"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error 1"},
		})

		result := Ap[int, int](valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 3)
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "function error 1")
		assert.Contains(t, messages, "function error 2")
		assert.Contains(t, messages, "value error 1")
	})
}

func TestMonadLaws(t *testing.T) {
	t.Run("functor identity law", func(t *testing.T) {
		// Map(id) == id
		value := Of(42)
		mapped := Map(F.Identity[int])(value)

		assert.Equal(t, value, mapped)
	})

	t.Run("functor composition law", func(t *testing.T) {
		// Map(f . g) == Map(f) . Map(g)
		f := func(x int) int { return x * 2 }
		g := func(x int) int { return x + 10 }
		composed := func(x int) int { return f(g(x)) }

		value := Of(5)
		left := Map(composed)(value)
		right := F.Pipe2(value, Map(g), Map(f))

		assert.Equal(t, left, right)
	})

	t.Run("applicative identity law", func(t *testing.T) {
		// Ap(v)(Of(id)) == v
		v := Of(42)
		result := Ap[int, int](v)(Of(F.Identity[int]))

		assert.Equal(t, v, result)
	})

	t.Run("applicative homomorphism law", func(t *testing.T) {
		// Ap(Of(x))(Of(f)) == Of(f(x))
		f := func(x int) int { return x * 2 }
		x := 21

		left := Ap[int, int](Of(x))(Of(f))
		right := Of(f(x))

		assert.Equal(t, left, right)
	})
}

func TestMapWithOperator(t *testing.T) {
	t.Run("Map returns an Operator", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		operator := Map(double)

		// Operator can be applied to different validations
		result1 := operator(Of(10))
		result2 := operator(Of(20))

		val1 := either.MonadFold(result1,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		val2 := either.MonadFold(result2,
			func(Errors) int { return 0 },
			F.Identity[int],
		)

		assert.Equal(t, 20, val1)
		assert.Equal(t, 40, val2)
	})
}

func TestApWithOperator(t *testing.T) {
	t.Run("Ap returns an Operator", func(t *testing.T) {
		valueValidation := Of(21)
		operator := Ap[int, int](valueValidation)

		// Operator can be applied to different function validations
		double := func(x int) int { return x * 2 }
		triple := func(x int) int { return x * 3 }

		result1 := operator(Of(double))
		result2 := operator(Of(triple))

		val1 := either.MonadFold(result1,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		val2 := either.MonadFold(result2,
			func(Errors) int { return 0 },
			F.Identity[int],
		)

		assert.Equal(t, 42, val1)
		assert.Equal(t, 63, val2)
	})
}

func TestApplicative(t *testing.T) {
	t.Run("returns non-nil instance", func(t *testing.T) {
		app := Applicative[int, string]()
		assert.NotNil(t, app)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		app1 := Applicative[int, string]()
		app2 := Applicative[int, string]()

		// Both should work independently
		result1 := app1.Of(42)
		result2 := app2.Of(43)

		assert.True(t, either.IsRight(result1))
		assert.True(t, either.IsRight(result2))

		val1 := either.MonadFold(result1, func(Errors) int { return 0 }, F.Identity[int])
		val2 := either.MonadFold(result2, func(Errors) int { return 0 }, F.Identity[int])

		assert.Equal(t, 42, val1)
		assert.Equal(t, 43, val2)
	})
}

func TestApplicativeOf(t *testing.T) {
	app := Applicative[int, string]()

	t.Run("wraps a value in Validation context", func(t *testing.T) {
		result := app.Of(42)
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("wraps string value", func(t *testing.T) {
		app := Applicative[string, int]()
		result := app.Of("hello")
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "hello", value)
	})

	t.Run("wraps zero value", func(t *testing.T) {
		result := app.Of(0)
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(Errors) int { return -1 },
			F.Identity[int],
		)
		assert.Equal(t, 0, value)
	})

	t.Run("wraps complex types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		app := Applicative[User, string]()
		user := User{Name: "Alice", Age: 30}
		result := app.Of(user)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) User { return User{} },
			F.Identity[User],
		)
		assert.Equal(t, user, value)
	})
}

func TestApplicativeMap(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("maps a function over successful validation", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		result := app.Map(double)(app.Of(21))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("maps type conversion", func(t *testing.T) {
		app := Applicative[int, string]()
		toString := func(x int) string { return fmt.Sprintf("%d", x) }
		result := app.Map(toString)(app.Of(42))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "42", value)
	})

	t.Run("maps identity function", func(t *testing.T) {
		result := app.Map(F.Identity[int])(app.Of(42))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("preserves failure", func(t *testing.T) {
		errs := Errors{&ValidationError{Messsage: "error"}}
		failure := Failures[int](errs)

		double := func(x int) int { return x * 2 }
		result := app.Map(double)(failure)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "error", errors[0].Messsage)
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		app := Applicative[int, string]()
		add10 := func(x int) int { return x + 10 }
		double := func(x int) int { return x * 2 }
		toString := func(x int) string { return fmt.Sprintf("%d", x) }

		result := F.Pipe3(
			app.Of(5),
			Map(add10),
			Map(double),
			app.Map(toString),
		)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "30", value)
	})
}

func TestApplicativeAp(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("applies wrapped function to wrapped value when both succeed", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		funcValidation := Of(double)
		valueValidation := app.Of(21)

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("accumulates errors when value fails", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		funcValidation := Of(double)
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "value error", errors[0].Messsage)
	})

	t.Run("accumulates errors when function fails", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := app.Of(21)

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "function error", errors[0].Messsage)
	})

	t.Run("accumulates all errors when both fail", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2)
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "function error")
		assert.Contains(t, messages, "value error")
	})

	t.Run("applies with type conversion", func(t *testing.T) {
		app := Applicative[int, string]()
		toString := func(x int) string { return fmt.Sprintf("value:%d", x) }
		funcValidation := Of(toString)
		valueValidation := app.Of(42)

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "value:42", value)
	})

	t.Run("accumulates multiple errors from different sources", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error 1"},
			&ValidationError{Messsage: "function error 2"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error 1"},
			&ValidationError{Messsage: "value error 2"},
		})

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 4)
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "function error 1")
		assert.Contains(t, messages, "function error 2")
		assert.Contains(t, messages, "value error 1")
		assert.Contains(t, messages, "value error 2")
	})
}

func TestApplicativeComposition(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("composes Map and Of", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		result := F.Pipe1(
			app.Of(21),
			app.Map(double),
		)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("composes multiple Map operations", func(t *testing.T) {
		app := Applicative[int, string]()
		double := func(x int) int { return x * 2 }
		toString := func(x int) string { return fmt.Sprintf("%d", x) }

		result := F.Pipe2(
			app.Of(21),
			Map(double),
			app.Map(toString),
		)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "42", value)
	})

	t.Run("composes Map and Ap", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		ioFunc := F.Pipe1(
			app.Of(5),
			Map(add),
		)
		valueValidation := app.Of(16)

		result := app.Ap(valueValidation)(ioFunc)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 21, value)
	})
}

func TestApplicativeLawsWithInstance(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("identity law: Ap(Of(id))(v) == v", func(t *testing.T) {
		identity := func(x int) int { return x }
		v := app.Of(42)

		left := app.Ap(v)(Of(identity))
		right := v

		assert.Equal(t, right, left)
	})

	t.Run("homomorphism law: Ap(Of(x))(Of(f)) == Of(f(x))", func(t *testing.T) {
		f := func(x int) int { return x * 2 }
		x := 21

		left := app.Ap(app.Of(x))(Of(f))
		right := app.Of(f(x))

		assert.Equal(t, right, left)
	})

	t.Run("interchange law: Ap(Of(y))(u) == Ap(u)(Of($ y))", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		u := Of(double)
		y := 21

		left := app.Ap(app.Of(y))(u)

		applyY := func(f func(int) int) int { return f(y) }
		right := Ap[int](u)(Of(applyY))

		assert.Equal(t, right, left)
	})

	t.Run("identity law with failure", func(t *testing.T) {
		identity := func(x int) int { return x }
		v := Failures[int](Errors{&ValidationError{Messsage: "error"}})

		left := app.Ap(v)(Of(identity))
		right := v

		assert.Equal(t, right, left)
	})
}

func TestApplicativeMultipleArguments(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("applies curried two-argument function", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		// Create validation with curried function
		funcValidation := F.Pipe1(
			app.Of(10),
			Map(add),
		)

		// Apply to second argument
		result := app.Ap(app.Of(32))(funcValidation)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("applies curried three-argument function", func(t *testing.T) {
		add3 := func(a int) func(int) func(int) int {
			return func(b int) func(int) int {
				return func(c int) int {
					return a + b + c
				}
			}
		}

		// Build up the computation step by step
		funcValidation1 := F.Pipe1(
			app.Of(10),
			Map(add3),
		)

		funcValidation2 := Ap[func(int) int](app.Of(20))(funcValidation1)
		result := Ap[int](app.Of(12))(funcValidation2)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("accumulates errors from multiple arguments", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		// First argument fails
		arg1 := Failures[int](Errors{&ValidationError{Messsage: "arg1 error"}})
		// Second argument fails
		arg2 := Failures[int](Errors{&ValidationError{Messsage: "arg2 error"}})

		funcValidation := F.Pipe1(arg1, Map(add))
		result := app.Ap(arg2)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2)
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "arg1 error")
		assert.Contains(t, messages, "arg2 error")
	})
}

func TestApplicativeWithDifferentTypes(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		app := Applicative[int, string]()
		toString := func(x int) string { return fmt.Sprintf("%d", x) }
		result := app.Map(toString)(app.Of(42))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "42", value)
	})

	t.Run("string to int", func(t *testing.T) {
		app := Applicative[string, int]()
		toLength := func(s string) int { return len(s) }
		result := app.Map(toLength)(app.Of("hello"))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 5, value)
	})

	t.Run("bool to string", func(t *testing.T) {
		app := Applicative[bool, string]()
		toString := func(b bool) string {
			if b {
				return "true"
			}
			return "false"
		}
		result := app.Map(toString)(app.Of(true))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "true", value)
	})
}

func TestApplicativeRealWorldScenario(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	t.Run("validates user with all valid fields", func(t *testing.T) {
		validateName := func(name string) Validation[string] {
			if len(name) < 3 {
				return Failures[string](Errors{&ValidationError{Messsage: "Name must be at least 3 characters"}})
			}
			return Success(name)
		}

		validateAge := func(age int) Validation[int] {
			if age < 18 {
				return Failures[int](Errors{&ValidationError{Messsage: "Must be 18 or older"}})
			}
			return Success(age)
		}

		validateEmail := func(email string) Validation[string] {
			if len(email) == 0 {
				return Failures[string](Errors{&ValidationError{Messsage: "Email is required"}})
			}
			return Success(email)
		}

		makeUser := func(name string) func(int) func(string) User {
			return func(age int) func(string) User {
				return func(email string) User {
					return User{Name: name, Age: age, Email: email}
				}
			}
		}

		name := validateName("Alice")
		age := validateAge(25)
		email := validateEmail("alice@example.com")

		// Use the standalone Ap function with proper type parameters
		result := Ap[User](email)(Ap[func(string) User](age)(Ap[func(int) func(string) User](name)(Of(makeUser))))

		assert.True(t, either.IsRight(result))
		user := either.MonadFold(result,
			func(Errors) User { return User{} },
			F.Identity[User],
		)
		assert.Equal(t, "Alice", user.Name)
		assert.Equal(t, 25, user.Age)
		assert.Equal(t, "alice@example.com", user.Email)
	})

	t.Run("accumulates all validation errors", func(t *testing.T) {
		validateName := func(name string) Validation[string] {
			if len(name) < 3 {
				return Failures[string](Errors{&ValidationError{Messsage: "Name must be at least 3 characters"}})
			}
			return Success(name)
		}

		validateAge := func(age int) Validation[int] {
			if age < 18 {
				return Failures[int](Errors{&ValidationError{Messsage: "Must be 18 or older"}})
			}
			return Success(age)
		}

		validateEmail := func(email string) Validation[string] {
			if len(email) == 0 {
				return Failures[string](Errors{&ValidationError{Messsage: "Email is required"}})
			}
			return Success(email)
		}

		makeUser := func(name string) func(int) func(string) User {
			return func(age int) func(string) User {
				return func(email string) User {
					return User{Name: name, Age: age, Email: email}
				}
			}
		}

		// All validations fail
		name := validateName("ab")
		age := validateAge(16)
		email := validateEmail("")

		// Use the standalone Ap function with proper type parameters
		result := Ap[User](email)(Ap[func(string) User](age)(Ap[func(int) func(string) User](name)(Of(makeUser))))

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(User) Errors { return nil },
		)
		assert.Len(t, errors, 3)
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "Name must be at least 3 characters")
		assert.Contains(t, messages, "Must be 18 or older")
		assert.Contains(t, messages, "Email is required")
	})
}
