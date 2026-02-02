// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validate

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	t.Run("creates successful validation with empty state", func(t *testing.T) {
		type State struct {
			x int
			y string
		}
		validator := Do[string](State{})
		result := validator("input")(nil)

		assert.Equal(t, either.Of[Errors](State{}), result)
	})

	t.Run("creates successful validation with initialized state", func(t *testing.T) {
		type State struct {
			x int
			y string
		}
		initial := State{x: 42, y: "hello"}
		validator := Do[string](initial)
		result := validator("input")(nil)

		assert.Equal(t, either.Of[Errors](initial), result)
	})

	t.Run("works with different input types", func(t *testing.T) {
		intValidator := Do[int](0)
		assert.Equal(t, either.Of[Errors](0), intValidator(42)(nil))

		strValidator := Do[string]("")
		assert.Equal(t, either.Of[Errors](""), strValidator("test")(nil))

		type Custom struct{ Value int }
		customValidator := Do[[]byte](Custom{Value: 100})
		assert.Equal(t, either.Of[Errors](Custom{Value: 100}), customValidator([]byte("data"))(nil))
	})
}

func TestBind(t *testing.T) {
	type State struct {
		x int
		y int
	}

	t.Run("binds successful validation to state", func(t *testing.T) {
		validator := F.Pipe2(
			Do[string](State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Validate[string, int] {
				return Of[string](42)
			}),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Validate[string, int] {
				return Of[string](10)
			}),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](State{x: 42, y: 10}), result)
	})

	t.Run("propagates failure", func(t *testing.T) {
		validator := F.Pipe2(
			Do[string](State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Validate[string, int] {
				return Of[string](42)
			}),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Validate[string, int] {
				return func(input string) Reader[Context, Validation[int]] {
					return func(ctx Context) Validation[int] {
						return validation.Failures[int](Errors{&validation.ValidationError{Messsage: "y failed"}})
					}
				}
			}),
		)

		result := validator("input")(nil)
		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(State) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "y failed", errors[0].Messsage)
	})

	t.Run("can access previous state values", func(t *testing.T) {
		validator := F.Pipe2(
			Do[string](State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Validate[string, int] {
				return Of[string](10)
			}),
			Bind(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, func(s State) Validate[string, int] {
				// y depends on x
				return Of[string](s.x * 2)
			}),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](State{x: 10, y: 20}), result)
	})

	t.Run("can access input value", func(t *testing.T) {
		validator := F.Pipe1(
			Do[int](State{}),
			Bind(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, func(s State) Validate[int, int] {
				return func(input int) Reader[Context, Validation[int]] {
					return func(ctx Context) Validation[int] {
						return validation.Success(input * 2)
					}
				}
			}),
		)

		result := validator(21)(nil)
		assert.Equal(t, either.Of[Errors](State{x: 42}), result)
	})
}

func TestLet(t *testing.T) {
	type State struct {
		x        int
		computed int
	}

	t.Run("attaches pure computation result to state", func(t *testing.T) {
		validator := F.Pipe1(
			Do[string](State{x: 5}),
			Let[string](func(c int) func(State) State {
				return func(s State) State { s.computed = c; return s }
			}, func(s State) int { return s.x * 2 }),
		)

		result := validator("input")(nil)
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 5, computed: 10}, value)
	})

	t.Run("preserves failure", func(t *testing.T) {
		failure := func(input string) Reader[Context, Validation[State]] {
			return func(ctx Context) Validation[State] {
				return validation.Failures[State](Errors{&validation.ValidationError{Messsage: "error"}})
			}
		}
		validator := Let[string](func(c int) func(State) State {
			return func(s State) State { s.computed = c; return s }
		}, func(s State) int { return s.x * 2 })

		result := validator(failure)("input")(nil)
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
		validator := F.Pipe3(
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

		result := validator("input")(nil)
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) State { return State{} },
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
		validator := F.Pipe1(
			Do[string](State{x: 5}),
			LetTo[string](func(n string) func(State) State {
				return func(s State) State { s.name = n; return s }
			}, "example"),
		)

		result := validator("input")(nil)
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{x: 5, name: "example"}, value)
	})

	t.Run("preserves failure", func(t *testing.T) {
		failure := func(input string) Reader[Context, Validation[State]] {
			return func(ctx Context) Validation[State] {
				return validation.Failures[State](Errors{&validation.ValidationError{Messsage: "error"}})
			}
		}
		validator := LetTo[string](func(n string) func(State) State {
			return func(s State) State { s.name = n; return s }
		}, "example")

		result := validator(failure)("input")(nil)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("sets multiple constant values", func(t *testing.T) {
		type State struct {
			name    string
			version int
			active  bool
		}
		validator := F.Pipe3(
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

		result := validator("input")(nil)
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) State { return State{} },
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
		validator := F.Pipe1(
			Of[string](42),
			BindTo[string](func(x int) State { return State{value: x} }),
		)

		result := validator("input")(nil)
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) State { return State{} },
			F.Identity[State],
		)
		assert.Equal(t, State{value: 42}, value)
	})

	t.Run("preserves failure", func(t *testing.T) {
		failure := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return validation.Failures[int](Errors{&validation.ValidationError{Messsage: "error"}})
			}
		}
		validator := BindTo[string](func(x int) State { return State{value: x} })

		result := validator(failure)("input")(nil)
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
		validator := F.Pipe1(
			Of[int]("hello"),
			BindTo[int](func(s string) StringState { return StringState{text: s} }),
		)

		result := validator(42)(nil)
		assert.Equal(t, either.Of[Errors](StringState{text: "hello"}), result)
	})
}

func TestApS(t *testing.T) {
	type State struct {
		x int
		y int
	}

	t.Run("attaches value using applicative pattern", func(t *testing.T) {
		validator := F.Pipe1(
			Do[string](State{}),
			ApS(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, Of[string](42)),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](State{x: 42}), result)
	})

	t.Run("accumulates errors from both validations", func(t *testing.T) {
		stateFailure := func(input string) Reader[Context, Validation[State]] {
			return func(ctx Context) Validation[State] {
				return validation.Failures[State](Errors{&validation.ValidationError{Messsage: "state error"}})
			}
		}
		valueFailure := func(input string) Reader[Context, Validation[int]] {
			return func(ctx Context) Validation[int] {
				return validation.Failures[int](Errors{&validation.ValidationError{Messsage: "value error"}})
			}
		}

		validator := ApS(func(x int) func(State) State {
			return func(s State) State { s.x = x; return s }
		}, valueFailure)

		result := validator(stateFailure)("input")(nil)
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
		validator := F.Pipe2(
			Do[string](State{}),
			ApS(func(x int) func(State) State {
				return func(s State) State { s.x = x; return s }
			}, Of[string](10)),
			ApS(func(y int) func(State) State {
				return func(s State) State { s.y = y; return s }
			}, Of[string](20)),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](State{x: 10, y: 20}), result)
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

		validator := F.Pipe1(
			Of[string](Person{Name: "Alice"}),
			ApSL(
				addressLens,
				Of[string](Address{Street: "Main St", City: "NYC"}),
			),
		)

		result := validator("input")(nil)
		expected := Person{
			Name:    "Alice",
			Address: Address{Street: "Main St", City: "NYC"},
		}
		assert.Equal(t, either.Of[Errors](expected), result)
	})

	t.Run("accumulates errors", func(t *testing.T) {
		addressLens := L.MakeLens(
			func(p Person) Address { return p.Address },
			func(p Person, a Address) Person { p.Address = a; return p },
		)

		personFailure := func(input string) Reader[Context, Validation[Person]] {
			return func(ctx Context) Validation[Person] {
				return validation.Failures[Person](Errors{&validation.ValidationError{Messsage: "person error"}})
			}
		}
		addressFailure := func(input string) Reader[Context, Validation[Address]] {
			return func(ctx Context) Validation[Address] {
				return validation.Failures[Address](Errors{&validation.ValidationError{Messsage: "address error"}})
			}
		}

		validator := ApSL(addressLens, addressFailure)
		result := validator(personFailure)("input")(nil)

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
		increment := func(v int) Validate[string, int] {
			return Of[string](v + 1)
		}

		validator := F.Pipe1(
			Of[string](Counter{Value: 42}),
			BindL(valueLens, increment),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](Counter{Value: 43}), result)
	})

	t.Run("fails validation based on current value", func(t *testing.T) {
		increment := func(v int) Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					if v >= 100 {
						return validation.Failures[int](Errors{&validation.ValidationError{Messsage: "exceeds limit"}})
					}
					return validation.Success(v + 1)
				}
			}
		}

		validator := F.Pipe1(
			Of[string](Counter{Value: 100}),
			BindL(valueLens, increment),
		)

		result := validator("input")(nil)
		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(Counter) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "exceeds limit", errors[0].Messsage)
	})

	t.Run("preserves failure", func(t *testing.T) {
		increment := func(v int) Validate[string, int] {
			return Of[string](v + 1)
		}

		failure := func(input string) Reader[Context, Validation[Counter]] {
			return func(ctx Context) Validation[Counter] {
				return validation.Failures[Counter](Errors{&validation.ValidationError{Messsage: "error"}})
			}
		}
		validator := BindL(valueLens, increment)
		result := validator(failure)("input")(nil)

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

		validator := F.Pipe1(
			Of[string](Counter{Value: 21}),
			LetL[string](valueLens, double),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](Counter{Value: 42}), result)
	})

	t.Run("preserves failure", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		failure := func(input string) Reader[Context, Validation[Counter]] {
			return func(ctx Context) Validation[Counter] {
				return validation.Failures[Counter](Errors{&validation.ValidationError{Messsage: "error"}})
			}
		}
		validator := LetL[string](valueLens, double)
		result := validator(failure)("input")(nil)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		add10 := func(v int) int { return v + 10 }
		double := func(v int) int { return v * 2 }

		validator := F.Pipe2(
			Of[string](Counter{Value: 5}),
			LetL[string](valueLens, add10),
			LetL[string](valueLens, double),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](Counter{Value: 30}), result)
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
		validator := F.Pipe1(
			Of[string](Config{Debug: true, Timeout: 30}),
			LetToL[string](debugLens, false),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](Config{Debug: false, Timeout: 30}), result)
	})

	t.Run("preserves failure", func(t *testing.T) {
		failure := func(input string) Reader[Context, Validation[Config]] {
			return func(ctx Context) Validation[Config] {
				return validation.Failures[Config](Errors{&validation.ValidationError{Messsage: "error"}})
			}
		}
		validator := LetToL[string](debugLens, false)
		result := validator(failure)("input")(nil)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("sets multiple fields", func(t *testing.T) {
		timeoutLens := L.MakeLens(
			func(c Config) int { return c.Timeout },
			func(c Config, t int) Config { c.Timeout = t; return c },
		)

		validator := F.Pipe2(
			Of[string](Config{Debug: true, Timeout: 30}),
			LetToL[string](debugLens, false),
			LetToL[string](timeoutLens, 60),
		)

		result := validator("input")(nil)
		assert.Equal(t, either.Of[Errors](Config{Debug: false, Timeout: 60}), result)
	})
}

func TestBindOperationsComposition(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	t.Run("combines Do, Bind, Let, and LetTo", func(t *testing.T) {
		validator := F.Pipe4(
			Do[string](User{}),
			LetTo[string](func(n string) func(User) User {
				return func(u User) User { u.Name = n; return u }
			}, "Alice"),
			Bind(func(a int) func(User) User {
				return func(u User) User { u.Age = a; return u }
			}, func(u User) Validate[string, int] {
				// Age validation
				if len(u.Name) > 0 {
					return Of[string](25)
				}
				return func(input string) Reader[Context, Validation[int]] {
					return func(ctx Context) Validation[int] {
						return validation.Failures[int](Errors{&validation.ValidationError{Messsage: "name required"}})
					}
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
			}, func(u User) Validate[string, int] {
				// Validate age is positive
				if u.Age > 0 {
					return Of[string](u.Age)
				}
				return func(input string) Reader[Context, Validation[int]] {
					return func(ctx Context) Validation[int] {
						return validation.Failures[int](Errors{&validation.ValidationError{Messsage: "age must be positive"}})
					}
				}
			}),
		)

		result := validator("input")(nil)
		expected := User{
			Name:  "Alice",
			Age:   25,
			Email: "Alice@example.com",
		}
		assert.Equal(t, either.Of[Errors](expected), result)
	})

	t.Run("validates with input-dependent logic", func(t *testing.T) {
		type Config struct {
			MaxValue int
			Value    int
		}

		validator := F.Pipe2(
			Do[int](Config{}),
			Bind(func(max int) func(Config) Config {
				return func(c Config) Config { c.MaxValue = max; return c }
			}, func(c Config) Validate[int, int] {
				// Extract max from input
				return func(input int) Reader[Context, Validation[int]] {
					return func(ctx Context) Validation[int] {
						return validation.Success(input)
					}
				}
			}),
			Bind(func(val int) func(Config) Config {
				return func(c Config) Config { c.Value = val; return c }
			}, func(c Config) Validate[int, int] {
				// Validate value against max
				return func(input int) Reader[Context, Validation[int]] {
					return func(ctx Context) Validation[int] {
						if input/2 <= c.MaxValue {
							return validation.Success(input / 2)
						}
						return validation.Failures[int](Errors{&validation.ValidationError{Messsage: "value exceeds max"}})
					}
				}
			}),
		)

		result := validator(100)(nil)
		assert.Equal(t, either.Of[Errors](Config{MaxValue: 100, Value: 50}), result)
	})
}
