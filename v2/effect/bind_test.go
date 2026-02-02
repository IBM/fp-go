// Copyright (c) 2023 - 2025 IBM Corp.
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

package effect

import (
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/stretchr/testify/assert"
)

type BindState struct {
	Name  string
	Age   int
	Email string
}

func TestDo(t *testing.T) {
	t.Run("creates effect with initial state", func(t *testing.T) {
		initial := BindState{Name: "Alice", Age: 30}
		eff := Do[TestContext](initial)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, initial, result)
	})

	t.Run("creates effect with empty struct", func(t *testing.T) {
		type Empty struct{}
		eff := Do[TestContext](Empty{})

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, Empty{}, result)
	})
}

func TestBind(t *testing.T) {
	t.Run("binds effect result to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := Bind(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) Effect[TestContext, int] {
				return Of[TestContext](30)
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("chains multiple binds", func(t *testing.T) {
		initial := BindState{}

		eff := Bind(
			func(email string) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Email = email
					return s
				}
			},
			func(s BindState) Effect[TestContext, string] {
				return Of[TestContext]("alice@example.com")
			},
		)(Bind(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) Effect[TestContext, int] {
				return Of[TestContext](30)
			},
		)(Bind(
			func(name string) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Name = name
					return s
				}
			},
			func(s BindState) Effect[TestContext, string] {
				return Of[TestContext]("Alice")
			},
		)(Do[TestContext](initial))))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
		assert.Equal(t, "alice@example.com", result.Email)
	})

	t.Run("propagates errors", func(t *testing.T) {
		expectedErr := errors.New("bind error")
		initial := BindState{Name: "Alice"}

		eff := Bind(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) Effect[TestContext, int] {
				return Fail[TestContext, int](expectedErr)
			},
		)(Do[TestContext](initial))

		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestLet(t *testing.T) {
	t.Run("computes value and binds to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := Let[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) int {
				return len(s.Name) * 10
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 50, result.Age) // len("Alice") * 10
	})

	t.Run("chains with Bind", func(t *testing.T) {
		initial := BindState{Name: "Bob"}

		eff := Let[TestContext](
			func(email string) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Email = email
					return s
				}
			},
			func(s BindState) string {
				return s.Name + "@example.com"
			},
		)(Bind(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) Effect[TestContext, int] {
				return Of[TestContext](25)
			},
		)(Do[TestContext](initial)))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Bob", result.Name)
		assert.Equal(t, 25, result.Age)
		assert.Equal(t, "Bob@example.com", result.Email)
	})
}

func TestLetTo(t *testing.T) {
	t.Run("binds constant value to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := LetTo[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			42,
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 42, result.Age)
	})

	t.Run("chains multiple LetTo", func(t *testing.T) {
		initial := BindState{}

		eff := LetTo[TestContext](
			func(email string) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Email = email
					return s
				}
			},
			"test@example.com",
		)(LetTo[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			30,
		)(LetTo[TestContext](
			func(name string) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Name = name
					return s
				}
			},
			"Alice",
		)(Do[TestContext](initial))))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
		assert.Equal(t, "test@example.com", result.Email)
	})
}

func TestBindTo(t *testing.T) {
	t.Run("wraps value in state", func(t *testing.T) {
		type SimpleState struct {
			Value int
		}

		eff := BindTo[TestContext](func(v int) SimpleState {
			return SimpleState{Value: v}
		})(Of[TestContext](42))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 42, result.Value)
	})

	t.Run("starts a bind chain", func(t *testing.T) {
		type State struct {
			X int
			Y string
		}

		eff := Let[TestContext](
			func(y string) func(State) State {
				return func(s State) State {
					s.Y = y
					return s
				}
			},
			func(s State) string {
				return "computed"
			},
		)(BindTo[TestContext](func(x int) State {
			return State{X: x}
		})(Of[TestContext](10)))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 10, result.X)
		assert.Equal(t, "computed", result.Y)
	})
}

func TestApS(t *testing.T) {
	t.Run("applies effect and binds result to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}
		ageEffect := Of[TestContext](30)

		eff := ApS(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			ageEffect,
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("propagates errors from applied effect", func(t *testing.T) {
		expectedErr := errors.New("aps error")
		initial := BindState{Name: "Alice"}
		ageEffect := Fail[TestContext, int](expectedErr)

		eff := ApS(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			ageEffect,
		)(Do[TestContext](initial))

		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestBindIOK(t *testing.T) {
	t.Run("binds IO operation to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := BindIOK[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) io.IO[int] {
				return func() int {
					return 30
				}
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})
}

func TestBindIOEitherK(t *testing.T) {
	t.Run("binds successful IOEither to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := BindIOEitherK[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) ioeither.IOEither[error, int] {
				return ioeither.Of[error](30)
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("propagates IOEither error", func(t *testing.T) {
		expectedErr := errors.New("ioeither error")
		initial := BindState{Name: "Alice"}

		eff := BindIOEitherK[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) ioeither.IOEither[error, int] {
				return ioeither.Left[int](expectedErr)
			},
		)(Do[TestContext](initial))

		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestBindIOResultK(t *testing.T) {
	t.Run("binds successful IOResult to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := BindIOResultK[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) ioresult.IOResult[int] {
				return ioresult.Of(30)
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})
}

func TestBindReaderK(t *testing.T) {
	t.Run("binds Reader operation to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := BindReaderK(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) reader.Reader[TestContext, int] {
				return func(ctx TestContext) int {
					return 30
				}
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})
}

func TestBindReaderIOK(t *testing.T) {
	t.Run("binds ReaderIO operation to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := BindReaderIOK(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) readerio.ReaderIO[TestContext, int] {
				return func(ctx TestContext) io.IO[int] {
					return func() int {
						return 30
					}
				}
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})
}

func TestBindEitherK(t *testing.T) {
	t.Run("binds successful Either to state", func(t *testing.T) {
		initial := BindState{Name: "Alice"}

		eff := BindEitherK[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) either.Either[error, int] {
				return either.Of[error](30)
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("propagates Either error", func(t *testing.T) {
		expectedErr := errors.New("either error")
		initial := BindState{Name: "Alice"}

		eff := BindEitherK[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			func(s BindState) either.Either[error, int] {
				return either.Left[int](expectedErr)
			},
		)(Do[TestContext](initial))

		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestLensOperations(t *testing.T) {
	// Create lenses for BindState
	nameLens := lens.MakeLens(
		func(s BindState) string { return s.Name },
		func(s BindState, name string) BindState {
			s.Name = name
			return s
		},
	)

	ageLens := lens.MakeLens(
		func(s BindState) int { return s.Age },
		func(s BindState, age int) BindState {
			s.Age = age
			return s
		},
	)

	t.Run("ApSL applies effect using lens", func(t *testing.T) {
		initial := BindState{Name: "Alice", Age: 25}
		ageEffect := Of[TestContext](30)

		eff := ApSL(ageLens, ageEffect)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("BindL binds effect using lens", func(t *testing.T) {
		initial := BindState{Name: "Alice", Age: 25}

		eff := BindL(
			ageLens,
			func(age int) Effect[TestContext, int] {
				return Of[TestContext](age + 5)
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("LetL computes value using lens", func(t *testing.T) {
		initial := BindState{Name: "Alice", Age: 25}

		eff := LetL[TestContext](
			ageLens,
			func(age int) int {
				return age * 2
			},
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 50, result.Age)
	})

	t.Run("LetToL sets constant using lens", func(t *testing.T) {
		initial := BindState{Name: "Alice", Age: 25}

		eff := LetToL[TestContext](ageLens, 100)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 100, result.Age)
	})

	t.Run("chains lens operations", func(t *testing.T) {
		initial := BindState{}

		eff := LetToL[TestContext](
			ageLens,
			30,
		)(LetToL[TestContext](
			nameLens,
			"Bob",
		)(Do[TestContext](initial)))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Bob", result.Name)
		assert.Equal(t, 30, result.Age)
	})
}

func TestApOperations(t *testing.T) {
	t.Run("ApIOS applies IO effect", func(t *testing.T) {
		initial := BindState{Name: "Alice"}
		ioEffect := func() int { return 30 }

		eff := ApIOS[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			ioEffect,
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("ApReaderS applies Reader effect", func(t *testing.T) {
		initial := BindState{Name: "Alice"}
		readerEffect := func(ctx TestContext) int { return 30 }

		eff := ApReaderS(
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			readerEffect,
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("ApEitherS applies Either effect", func(t *testing.T) {
		initial := BindState{Name: "Alice"}
		eitherEffect := either.Of[error](30)

		eff := ApEitherS[TestContext](
			func(age int) func(BindState) BindState {
				return func(s BindState) BindState {
					s.Age = age
					return s
				}
			},
			eitherEffect,
		)(Do[TestContext](initial))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 30, result.Age)
	})
}

func TestComplexBindChain(t *testing.T) {
	t.Run("builds complex state with multiple operations", func(t *testing.T) {
		type ComplexState struct {
			Name    string
			Age     int
			Email   string
			IsAdmin bool
			Score   int
		}

		eff := LetTo[TestContext](
			func(score int) func(ComplexState) ComplexState {
				return func(s ComplexState) ComplexState {
					s.Score = score
					return s
				}
			},
			100,
		)(Let[TestContext](
			func(isAdmin bool) func(ComplexState) ComplexState {
				return func(s ComplexState) ComplexState {
					s.IsAdmin = isAdmin
					return s
				}
			},
			func(s ComplexState) bool {
				return s.Age >= 18
			},
		)(Let[TestContext](
			func(email string) func(ComplexState) ComplexState {
				return func(s ComplexState) ComplexState {
					s.Email = email
					return s
				}
			},
			func(s ComplexState) string {
				return s.Name + "@example.com"
			},
		)(Bind(
			func(age int) func(ComplexState) ComplexState {
				return func(s ComplexState) ComplexState {
					s.Age = age
					return s
				}
			},
			func(s ComplexState) Effect[TestContext, int] {
				return Of[TestContext](25)
			},
		)(BindTo[TestContext](func(name string) ComplexState {
			return ComplexState{Name: name}
		})(Of[TestContext]("Alice"))))))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 25, result.Age)
		assert.Equal(t, "Alice@example.com", result.Email)
		assert.True(t, result.IsAdmin)
		assert.Equal(t, 100, result.Score)
	})
}
