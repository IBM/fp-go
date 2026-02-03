package validation

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	t.Run("creates successful validation with empty state", func(t *testing.T) {
		type State struct {
			x int
			y string
		}
		result := Do(State{})

		assert.Equal(t, Of(State{}), result)
	})

	t.Run("creates successful validation with initialized state", func(t *testing.T) {
		type State struct {
			x int
			y string
		}
		initial := State{x: 42, y: "hello"}
		result := Do(initial)

		assert.Equal(t, Of(initial), result)
	})

	t.Run("works with different types", func(t *testing.T) {
		intResult := Do(0)
		assert.Equal(t, Of(0), intResult)

		strResult := Do("")
		assert.Equal(t, Of(""), strResult)

		type Custom struct{ Value int }
		customResult := Do(Custom{Value: 100})
		assert.Equal(t, Of(Custom{Value: 100}), customResult)
	})
}

func TestBind(t *testing.T) {
	type State struct {
		x int
		y int
	}

	t.Run("binds successful validation to state", func(t *testing.T) {
		result := F.Pipe2(
			Do(State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Validation[int] { return Success(42) }),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Validation[int] { return Success(10) }),
		)

		assert.Equal(t, Of(State{x: 42, y: 10}), result)
	})

	t.Run("propagates failure", func(t *testing.T) {
		result := F.Pipe2(
			Do(State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Validation[int] { return Success(42) }),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Validation[int] {
				return Failures[int](Errors{&ValidationError{Messsage: "y failed"}})
			}),
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(State) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "y failed", errors[0].Messsage)
	})

	t.Run("can access previous state values", func(t *testing.T) {
		result := F.Pipe2(
			Do(State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Validation[int] { return Success(10) }),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Validation[int] {
				// y depends on x
				return Success(s.x * 2)
			}),
		)

		assert.Equal(t, Success(State{x: 10, y: 20}), result)
	})
}

func TestLet(t *testing.T) {
	type State struct {
		x        int
		computed int
	}

	t.Run("attaches pure computation result to state", func(t *testing.T) {
		result := F.Pipe1(
			Do(State{x: 5}),
			Let(func(c int) func(State) State {
				return func(s State) State { s.computed = c; return s }
			}, func(s State) int { return s.x * 2 }),
		)

		assert.Equal(t, Of(State{x: 5, computed: 10}), result)
	})

	t.Run("preserves failure", func(t *testing.T) {
		failure := Failures[State](Errors{&ValidationError{Messsage: "error"}})
		result := Let(func(c int) func(State) State {
			return func(s State) State { s.computed = c; return s }
		}, func(s State) int { return s.x * 2 })(failure)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(State) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "error", errors[0].Messsage)
	})

	t.Run("chains multiple Let operations", func(t *testing.T) {
		type State struct {
			x int
			y int
			z int
		}
		result := F.Pipe3(
			Do(State{x: 5}),
			Let(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) int { return s.x * 2 }),
			Let(func(z int) func(State) State {
				return func(s State) State { s.z = z; return s }
			}, func(s State) int { return s.y + 10 }),
			Let(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) int { return s.z * 3 }),
		)

		assert.Equal(t, Of(State{x: 60, y: 10, z: 20}), result)
	})
}

func TestLetTo(t *testing.T) {
	type State struct {
		x    int
		name string
	}

	t.Run("attaches constant value to state", func(t *testing.T) {
		result := F.Pipe1(
			Do(State{x: 5}),
			LetTo(func(n string) func(State) State {
				return func(s State) State { s.name = n; return s }
			}, "example"),
		)

		assert.Equal(t, Of(State{x: 5, name: "example"}), result)
	})

	t.Run("preserves failure", func(t *testing.T) {
		failure := Failures[State](Errors{&ValidationError{Messsage: "error"}})
		result := LetTo(func(n string) func(State) State {
			return func(s State) State { s.name = n; return s }
		}, "example")(failure)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("sets multiple constant values", func(t *testing.T) {
		type State struct {
			name    string
			version int
			active  bool
		}
		result := F.Pipe3(
			Do(State{}),
			LetTo(func(n string) func(State) State {
				return func(s State) State { s.name = n; return s }
			}, "app"),
			LetTo(func(v int) func(State) State {
				return func(s State) State { s.version = v; return s }
			}, 2),
			LetTo(func(a bool) func(State) State {
				return func(s State) State { s.active = a; return s }
			}, true),
		)

		assert.Equal(t, Of(State{name: "app", version: 2, active: true}), result)
	})
}

func TestBindTo(t *testing.T) {
	type State struct {
		value int
	}

	t.Run("initializes state from value", func(t *testing.T) {
		result := F.Pipe1(
			Success(42),
			BindTo(func(x int) State { return State{value: x} }),
		)

		assert.Equal(t, Of(State{value: 42}), result)
	})

	t.Run("preserves failure", func(t *testing.T) {
		failure := Failures[int](Errors{&ValidationError{Messsage: "error"}})
		result := BindTo(func(x int) State { return State{value: x} })(failure)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(State) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "error", errors[0].Messsage)
	})

	t.Run("works with different types", func(t *testing.T) {
		type StringState struct {
			text string
		}
		result := F.Pipe1(
			Success("hello"),
			BindTo(func(s string) StringState { return StringState{text: s} }),
		)

		assert.Equal(t, Of(StringState{text: "hello"}), result)
	})
}

func TestApS(t *testing.T) {
	type State struct {
		x int
		y int
	}

	t.Run("attaches value using applicative pattern", func(t *testing.T) {
		result := F.Pipe1(
			Do(State{}),
			ApS(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, Success(42)),
		)

		assert.Equal(t, Of(State{x: 42}), result)
	})

	t.Run("accumulates errors from both validations", func(t *testing.T) {
		stateFailure := Failures[State](Errors{&ValidationError{Messsage: "state error"}})
		valueFailure := Failures[int](Errors{&ValidationError{Messsage: "value error"}})

		result := ApS(func(x int) func(State) State {
			return func(s State) State { s.x = x; return s }
		}, valueFailure)(stateFailure)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(State) Errors { return nil },
		)
		assert.Len(t, errors, 2)
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "state error")
		assert.Contains(t, messages, "value error")
	})

	t.Run("combines multiple ApS operations", func(t *testing.T) {
		result := F.Pipe2(
			Do(State{}),
			ApS(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, Success(10)),
			ApS(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, Success(20)),
		)

		assert.Equal(t, Of(State{x: 10, y: 20}), result)
	})
}

func TestApSL(t *testing.T) {
	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name    string
		Address Address
	}

	t.Run("updates nested structure using lens", func(t *testing.T) {
		addressLens := L.MakeLens(
			func(p Person) Address { return p.Address },
			func(p Person, a Address) Person { p.Address = a; return p },
		)

		result := F.Pipe1(
			Success(Person{Name: "Alice"}),
			ApSL(
				addressLens,
				Success(Address{Street: "Main St", City: "NYC"}),
			),
		)

		expected := Person{
			Name:    "Alice",
			Address: Address{Street: "Main St", City: "NYC"},
		}
		assert.Equal(t, Of(expected), result)
	})

	t.Run("accumulates errors", func(t *testing.T) {
		addressLens := L.MakeLens(
			func(p Person) Address { return p.Address },
			func(p Person, a Address) Person { p.Address = a; return p },
		)

		personFailure := Failures[Person](Errors{&ValidationError{Messsage: "person error"}})
		addressFailure := Failures[Address](Errors{&ValidationError{Messsage: "address error"}})

		result := ApSL(addressLens, addressFailure)(personFailure)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(Person) Errors { return nil },
		)
		assert.Len(t, errors, 2)
	})
}

func TestBindL(t *testing.T) {
	type Counter struct {
		Value int
	}

	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("updates field based on current value", func(t *testing.T) {
		increment := func(v int) Validation[int] {
			return Success(v + 1)
		}

		result := F.Pipe1(
			Success(Counter{Value: 42}),
			BindL(valueLens, increment),
		)

		assert.Equal(t, Of(Counter{Value: 43}), result)
	})

	t.Run("fails validation based on current value", func(t *testing.T) {
		increment := func(v int) Validation[int] {
			if v >= 100 {
				return Failures[int](Errors{&ValidationError{Messsage: "exceeds limit"}})
			}
			return Success(v + 1)
		}

		result := F.Pipe1(
			Success(Counter{Value: 100}),
			BindL(valueLens, increment),
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(Counter) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "exceeds limit", errors[0].Messsage)
	})

	t.Run("preserves failure", func(t *testing.T) {
		increment := func(v int) Validation[int] {
			return Success(v + 1)
		}

		failure := Failures[Counter](Errors{&ValidationError{Messsage: "error"}})
		result := BindL(valueLens, increment)(failure)

		assert.True(t, either.IsLeft(result))
	})
}

func TestLetL(t *testing.T) {
	type Counter struct {
		Value int
	}

	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("transforms field with pure function", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		result := F.Pipe1(
			Success(Counter{Value: 21}),
			LetL(valueLens, double),
		)

		assert.Equal(t, Of(Counter{Value: 42}), result)
	})

	t.Run("preserves failure", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		failure := Failures[Counter](Errors{&ValidationError{Messsage: "error"}})
		result := LetL(valueLens, double)(failure)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		add10 := func(v int) int { return v + 10 }
		double := func(v int) int { return v * 2 }

		result := F.Pipe2(
			Success(Counter{Value: 5}),
			LetL(valueLens, add10),
			LetL(valueLens, double),
		)

		assert.Equal(t, Of(Counter{Value: 30}), result)
	})
}

func TestLetToL(t *testing.T) {
	type Config struct {
		Debug   bool
		Timeout int
	}

	debugLens := L.MakeLens(
		func(c Config) bool { return c.Debug },
		func(c Config, d bool) Config { c.Debug = d; return c },
	)

	t.Run("sets field to constant value", func(t *testing.T) {
		result := F.Pipe1(
			Success(Config{Debug: true, Timeout: 30}),
			LetToL(debugLens, false),
		)

		assert.Equal(t, Of(Config{Debug: false, Timeout: 30}), result)
	})

	t.Run("preserves failure", func(t *testing.T) {
		failure := Failures[Config](Errors{&ValidationError{Messsage: "error"}})
		result := LetToL(debugLens, false)(failure)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("sets multiple fields", func(t *testing.T) {
		timeoutLens := L.MakeLens(
			func(c Config) int { return c.Timeout },
			func(c Config, t int) Config { c.Timeout = t; return c },
		)

		result := F.Pipe2(
			Success(Config{Debug: true, Timeout: 30}),
			LetToL(debugLens, false),
			LetToL(timeoutLens, 60),
		)

		assert.Equal(t, Of(Config{Debug: false, Timeout: 60}), result)
	})
}

func TestBindOperationsComposition(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	t.Run("combines Do, Bind, Let, and LetTo", func(t *testing.T) {
		result := F.Pipe4(
			Do(User{}),
			LetTo(func(n string) func(User) User {
				return func(u User) User { u.Name = n; return u }
			}, "Alice"),
			Bind(func(a int) func(User) User {
				return func(u User) User { u.Age = a; return u }
			}, func(u User) Validation[int] {
				// Age validation
				if len(u.Name) > 0 {
					return Success(25)
				}
				return Failures[int](Errors{&ValidationError{Messsage: "name required"}})
			}),
			Let(func(e string) func(User) User {
				return func(u User) User { u.Email = e; return u }
			}, func(u User) string {
				// Derive email from name
				return u.Name + "@example.com"
			}),
			Bind(func(a int) func(User) User {
				return func(u User) User { u.Age = a; return u }
			}, func(u User) Validation[int] {
				// Validate age is positive
				if u.Age > 0 {
					return Success(u.Age)
				}
				return Failures[int](Errors{&ValidationError{Messsage: "age must be positive"}})
			}),
		)

		expected := User{Name: "Alice", Age: 25, Email: "Alice@example.com"}
		assert.Equal(t, Of(expected), result)
	})
}
