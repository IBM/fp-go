package decode

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	t.Run("creates decoder with empty state", func(t *testing.T) {
		type State struct {
			x int
			y string
		}
		decoder := Do[string](State{})
		result := decoder("input")

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{}, value)
	})

	t.Run("creates decoder with initialized state", func(t *testing.T) {
		type State struct {
			x int
			y string
		}
		initial := State{x: 42, y: "hello"}
		decoder := Do[string](initial)
		result := decoder("input")

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, initial, value)
	})

	t.Run("works with different input types", func(t *testing.T) {
		intDecoder := Do[int](0)
		assert.True(t, either.IsRight(intDecoder(42)))

		strDecoder := Do[string]("")
		assert.True(t, either.IsRight(strDecoder("test")))

		type Custom struct{ Value int }
		customDecoder := Do[[]byte](Custom{Value: 100})
		assert.True(t, either.IsRight(customDecoder([]byte("data"))))
	})
}

func TestBind(t *testing.T) {
	type State struct {
		x int
		y int
	}

	t.Run("binds successful decode to state", func(t *testing.T) {
		decoder := F.Pipe2(
			Do[string](State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Decode[string, int] {
				return Of[string](42)
			}),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Decode[string, int] {
				return Of[string](10)
			}),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 42, y: 10}, value)
	})

	t.Run("propagates failure", func(t *testing.T) {
		decoder := F.Pipe2(
			Do[string](State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Decode[string, int] {
				return Of[string](42)
			}),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Decode[string, int] {
				return func(input string) validation.Validation[int] {
					return validation.Failures[int](validation.Errors{
						&validation.ValidationError{Messsage: "y failed"},
					})
				}
			}),
		)

		result := decoder("input")
		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(State) validation.Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "y failed", errors[0].Messsage)
	})

	t.Run("can access previous state values", func(t *testing.T) {
		decoder := F.Pipe2(
			Do[string](State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Decode[string, int] {
				return Of[string](10)
			}),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Decode[string, int] {
				// y depends on x
				return Of[string](s.x * 2)
			}),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 10, y: 20}, value)
	})

	t.Run("can access input in decoder", func(t *testing.T) {
		decoder := F.Pipe1(
			Do[string](State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Decode[string, int] {
				return func(input string) validation.Validation[int] {
					// Use input to determine value
					if input == "large" {
						return validation.Success(100)
					}
					return validation.Success(10)
				}
			}),
		)

		result1 := decoder("large")
		value1 := either.MonadFold(result1,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, 100, value1.x)

		result2 := decoder("small")
		value2 := either.MonadFold(result2,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, 10, value2.x)
	})
}

func TestLet(t *testing.T) {
	type State struct {
		x        int
		computed int
	}

	t.Run("attaches pure computation result to state", func(t *testing.T) {
		decoder := F.Pipe1(
			Do[string](State{x: 5}),
			Let[string](func(c int) func(State) State {
				return func(s State) State { s.computed = c; return s }
			}, func(s State) int { return s.x * 2 }),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 5, computed: 10}, value)
	})

	t.Run("chains multiple Let operations", func(t *testing.T) {
		type State struct {
			x int
			y int
			z int
		}
		decoder := F.Pipe3(
			Do[string](State{x: 5}),
			Let[string](func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) int { return s.x * 2 }),
			Let[string](func(z int) func(State) State {
				return func(s State) State { s.z = z; return s }
			}, func(s State) int { return s.y + 10 }),
			Let[string](func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) int { return s.z * 3 }),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 60, y: 10, z: 20}, value)
	})
}

func TestLetTo(t *testing.T) {
	type State struct {
		x    int
		name string
	}

	t.Run("attaches constant value to state", func(t *testing.T) {
		decoder := F.Pipe1(
			Do[string](State{x: 5}),
			LetTo[string](func(n string) func(State) State {
				return func(s State) State { s.name = n; return s }
			}, "example"),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 5, name: "example"}, value)
	})

	t.Run("sets multiple constant values", func(t *testing.T) {
		type State struct {
			name    string
			version int
			active  bool
		}
		decoder := F.Pipe3(
			Do[string](State{}),
			LetTo[string](func(n string) func(State) State {
				return func(s State) State { s.name = n; return s }
			}, "app"),
			LetTo[string](func(v int) func(State) State {
				return func(s State) State { s.version = v; return s }
			}, 2),
			LetTo[string](func(a bool) func(State) State {
				return func(s State) State { s.active = a; return s }
			}, true),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{name: "app", version: 2, active: true}, value)
	})
}

func TestBindTo(t *testing.T) {
	type State struct {
		value int
	}

	t.Run("initializes state from value", func(t *testing.T) {
		decoder := F.Pipe1(
			Of[string](42),
			BindTo[string](func(x int) State { return State{value: x} }),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{value: 42}, value)
	})

	t.Run("works with different types", func(t *testing.T) {
		type StringState struct {
			text string
		}
		decoder := F.Pipe1(
			Of[string]("hello"),
			BindTo[string](func(s string) StringState { return StringState{text: s} }),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) StringState { return StringState{} },
			F.Identity[StringState],
		)
		assert.Equal(t, StringState{text: "hello"}, value)
	})
}

func TestApS(t *testing.T) {
	type State struct {
		x int
		y int
	}

	t.Run("attaches value using applicative pattern", func(t *testing.T) {
		decoder := F.Pipe1(
			Do[string](State{}),
			ApS(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, Of[string](42)),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 42}, value)
	})

	t.Run("accumulates errors from both decoders", func(t *testing.T) {
		stateDecoder := func(input string) validation.Validation[State] {
			return validation.Failures[State](validation.Errors{
				&validation.ValidationError{Messsage: "state error"},
			})
		}
		valueDecoder := func(input string) validation.Validation[int] {
			return validation.Failures[int](validation.Errors{
				&validation.ValidationError{Messsage: "value error"},
			})
		}

		decoder := ApS(func(x int) func(State) State {
			return func(s State) State { s.x = x; return s }
		}, valueDecoder)(stateDecoder)

		result := decoder("input")
		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(State) validation.Errors { return nil },
		)
		assert.Len(t, errors, 2)
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "state error")
		assert.Contains(t, messages, "value error")
	})

	t.Run("combines multiple ApS operations", func(t *testing.T) {
		decoder := F.Pipe2(
			Do[string](State{}),
			ApS(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, Of[string](10)),
			ApS(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, Of[string](20)),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 10, y: 20}, value)
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

		decoder := F.Pipe1(
			Of[string](Person{Name: "Alice"}),
			ApSL(
				addressLens,
				Of[string](Address{Street: "Main St", City: "NYC"}),
			),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) Person { return Person{} },
			F.Identity[Person],
		)
		assert.Equal(t, "Alice", value.Name)
		assert.Equal(t, "Main St", value.Address.Street)
		assert.Equal(t, "NYC", value.Address.City)
	})

	t.Run("accumulates errors", func(t *testing.T) {
		addressLens := L.MakeLens(
			func(p Person) Address { return p.Address },
			func(p Person, a Address) Person { p.Address = a; return p },
		)

		personDecoder := func(input string) validation.Validation[Person] {
			return validation.Failures[Person](validation.Errors{
				&validation.ValidationError{Messsage: "person error"},
			})
		}
		addressDecoder := func(input string) validation.Validation[Address] {
			return validation.Failures[Address](validation.Errors{
				&validation.ValidationError{Messsage: "address error"},
			})
		}

		decoder := ApSL(addressLens, addressDecoder)(personDecoder)
		result := decoder("input")

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(Person) validation.Errors { return nil },
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
		increment := func(v int) Decode[string, int] {
			return Of[string](v + 1)
		}

		decoder := F.Pipe1(
			Of[string](Counter{Value: 42}),
			BindL(valueLens, increment),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) Counter { return Counter{} },
			F.Identity[Counter],
		)
		assert.Equal(t, Counter{Value: 43}, value)
	})

	t.Run("fails validation based on current value", func(t *testing.T) {
		increment := func(v int) Decode[string, int] {
			return func(input string) validation.Validation[int] {
				if v >= 100 {
					return validation.Failures[int](validation.Errors{
						&validation.ValidationError{Messsage: "exceeds limit"},
					})
				}
				return validation.Success(v + 1)
			}
		}

		decoder := F.Pipe1(
			Of[string](Counter{Value: 100}),
			BindL(valueLens, increment),
		)

		result := decoder("input")
		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(Counter) validation.Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "exceeds limit", errors[0].Messsage)
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

		decoder := F.Pipe1(
			Of[string](Counter{Value: 21}),
			LetL[string](valueLens, double),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) Counter { return Counter{} },
			F.Identity[Counter],
		)
		assert.Equal(t, Counter{Value: 42}, value)
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		add10 := func(v int) int { return v + 10 }
		double := func(v int) int { return v * 2 }

		decoder := F.Pipe2(
			Of[string](Counter{Value: 5}),
			LetL[string](valueLens, add10),
			LetL[string](valueLens, double),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) Counter { return Counter{} },
			F.Identity[Counter],
		)
		assert.Equal(t, Counter{Value: 30}, value)
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
		decoder := F.Pipe1(
			Of[string](Config{Debug: true, Timeout: 30}),
			LetToL[string](debugLens, false),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) Config { return Config{} },
			F.Identity[Config],
		)
		assert.Equal(t, Config{Debug: false, Timeout: 30}, value)
	})

	t.Run("sets multiple fields", func(t *testing.T) {
		timeoutLens := L.MakeLens(
			func(c Config) int { return c.Timeout },
			func(c Config, t int) Config { c.Timeout = t; return c },
		)

		decoder := F.Pipe2(
			Of[string](Config{Debug: true, Timeout: 30}),
			LetToL[string](debugLens, false),
			LetToL[string](timeoutLens, 60),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) Config { return Config{} },
			F.Identity[Config],
		)
		assert.Equal(t, Config{Debug: false, Timeout: 60}, value)
	})
}

func TestBindOperationsComposition(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	t.Run("combines Do, Bind, Let, and LetTo", func(t *testing.T) {
		decoder := F.Pipe4(
			Do[string](User{}),
			LetTo[string](func(n string) func(User) User {
				return func(u User) User { u.Name = n; return u }
			}, "Alice"),
			Bind(func(a int) func(User) User {
				return func(u User) User { u.Age = a; return u }
			}, func(u User) Decode[string, int] {
				// Age validation
				if len(u.Name) > 0 {
					return Of[string](25)
				}
				return func(input string) validation.Validation[int] {
					return validation.Failures[int](validation.Errors{
						&validation.ValidationError{Messsage: "name required"},
					})
				}
			}),
			Let[string](func(e string) func(User) User {
				return func(u User) User { u.Email = e; return u }
			}, func(u User) string {
				// Derive email from name
				return u.Name + "@example.com"
			}),
			Bind(func(a int) func(User) User {
				return func(u User) User { u.Age = a; return u }
			}, func(u User) Decode[string, int] {
				// Validate age is positive
				if u.Age > 0 {
					return Of[string](u.Age)
				}
				return func(input string) validation.Validation[int] {
					return validation.Failures[int](validation.Errors{
						&validation.ValidationError{Messsage: "age must be positive"},
					})
				}
			}),
		)

		result := decoder("input")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) User { return User{} },
			F.Identity[User],
		)
		assert.Equal(t, "Alice", value.Name)
		assert.Equal(t, 25, value.Age)
		assert.Equal(t, "Alice@example.com", value.Email)
	})
}
