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

package iterresult

import (
	"errors"
	"slices"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	L "github.com/IBM/fp-go/v2/optics/lens"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) SeqResult[string] {
	return Of("Doe")
}

func getGivenName(s utils.WithLastName) SeqResult[string] {
	return Of("John")
}

func TestBind(t *testing.T) {
	t.Run("successful bind chain", func(t *testing.T) {
		res := F.Pipe3(
			Do(utils.Empty),
			Bind(utils.SetLastName, getLastName),
			Bind(utils.SetGivenName, getGivenName),
			Map(utils.GetFullName),
		)

		collected := slices.Collect(F.Pipe1(res, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []string{"John Doe"}, collected)
	})

	t.Run("bind with error", func(t *testing.T) {
		getError := func(s utils.Initial) SeqResult[string] {
			return Left[string](errors.New("test error"))
		}

		res := F.Pipe2(
			Do(utils.Empty),
			Bind(utils.SetLastName, getError),
			Bind(utils.SetGivenName, getGivenName),
		)

		var err error
		for e := range res {
			err = R.MonadFold(e,
				F.Identity[error],
				func(s utils.WithGivenName) error { t.Fatal("expected Right"); return nil },
			)
			break
		}

		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
	})
}

func TestApS(t *testing.T) {
	t.Run("successful ApS chain", func(t *testing.T) {
		res := F.Pipe3(
			Do(utils.Empty),
			ApS(utils.SetLastName, Of("Doe")),
			ApS(utils.SetGivenName, Of("John")),
			Map(utils.GetFullName),
		)

		collected := slices.Collect(F.Pipe1(res, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []string{"John Doe"}, collected)
	})

	t.Run("ApS with error in first value", func(t *testing.T) {
		res := F.Pipe2(
			Do(utils.Empty),
			ApS(utils.SetLastName, Left[string](errors.New("error1"))),
			ApS(utils.SetGivenName, Of("John")),
		)

		var err error
		for e := range res {
			err = R.MonadFold(e,
				F.Identity[error],
				func(s utils.WithGivenName) error { t.Fatal("expected Left"); return nil },
			)
			break
		}

		assert.Error(t, err)
		assert.Equal(t, "error1", err.Error())
	})
}

type Counter struct {
	Value int
}

type Person struct {
	Name string
	Age  int
}

type Config struct {
	Debug   bool
	Timeout int
}

func TestApSL(t *testing.T) {
	ageLens := L.MakeLens(
		func(p Person) int { return p.Age },
		func(p Person, a int) Person { p.Age = a; return p },
	)

	t.Run("ApSL with Right value", func(t *testing.T) {
		result := F.Pipe1(
			Right(Person{Name: "Alice", Age: 25}),
			ApSL(ageLens, Right(30)),
		)

		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[Person] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []Person{{Name: "Alice", Age: 30}}, collected)
	})

	t.Run("ApSL with Left in context", func(t *testing.T) {
		result := F.Pipe1(
			Left[Person](assert.AnError),
			ApSL(ageLens, Right(30)),
		)

		var err error
		for e := range result {
			err = R.MonadFold(e,
				F.Identity[error],
				func(p Person) error { t.Fatal("expected Left"); return nil },
			)
			break
		}

		assert.Equal(t, assert.AnError, err)
	})

	t.Run("ApSL with Left in value", func(t *testing.T) {
		result := F.Pipe1(
			Right(Person{Name: "Alice", Age: 25}),
			ApSL(ageLens, Left[int](assert.AnError)),
		)

		var err error
		for e := range result {
			err = R.MonadFold(e,
				F.Identity[error],
				func(p Person) error { t.Fatal("expected Left"); return nil },
			)
			break
		}

		assert.Equal(t, assert.AnError, err)
	})
}

func TestBindL(t *testing.T) {
	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("BindL with successful transformation", func(t *testing.T) {
		increment := func(v int) SeqResult[int] {
			if v >= 100 {
				return Left[int](errors.New("overflow"))
			}
			return Right(v + 1)
		}

		result := F.Pipe1(
			Of(Counter{Value: 42}),
			BindL(valueLens, increment),
		)

		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[Counter] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []Counter{{Value: 43}}, collected)
	})

	t.Run("BindL with error", func(t *testing.T) {
		increment := func(v int) SeqResult[int] {
			if v >= 100 {
				return Left[int](errors.New("overflow"))
			}
			return Right(v + 1)
		}

		result := F.Pipe1(
			Of(Counter{Value: 100}),
			BindL(valueLens, increment),
		)

		var err error
		for e := range result {
			err = R.MonadFold(e,
				F.Identity[error],
				func(c Counter) error { t.Fatal("expected Left"); return nil },
			)
			break
		}

		assert.Error(t, err)
		assert.Equal(t, "overflow", err.Error())
	})
}

func TestLetL(t *testing.T) {
	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("LetL with pure transformation", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		result := F.Pipe1(
			Of(Counter{Value: 21}),
			LetL(valueLens, double),
		)

		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[Counter] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []Counter{{Value: 42}}, collected)
	})

	t.Run("LetL with Left value", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		result := F.Pipe1(
			Left[Counter](errors.New("error")),
			LetL(valueLens, double),
		)

		var err error
		for e := range result {
			err = R.MonadFold(e,
				F.Identity[error],
				func(c Counter) error { t.Fatal("expected Left"); return nil },
			)
			break
		}

		assert.Error(t, err)
	})
}

func TestLetToL(t *testing.T) {
	debugLens := L.MakeLens(
		func(c Config) bool { return c.Debug },
		func(c Config, d bool) Config { c.Debug = d; return c },
	)

	t.Run("LetToL sets constant value", func(t *testing.T) {
		result := F.Pipe1(
			Of(Config{Debug: true, Timeout: 30}),
			LetToL(debugLens, false),
		)

		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[Config] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []Config{{Debug: false, Timeout: 30}}, collected)
	})
}

func TestLet(t *testing.T) {
	t.Run("Let with pure function", func(t *testing.T) {
		type State struct {
			Value  int
			Double int
		}

		result := F.Pipe2(
			Do(State{Value: 21}),
			Let(
				func(d int) func(State) State {
					return func(s State) State { s.Double = d; return s }
				},
				func(s State) int { return s.Value * 2 },
			),
			Map(func(s State) int { return s.Double }),
		)

		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []int{42}, collected)
	})
}

func TestLetTo(t *testing.T) {
	t.Run("LetTo with constant value", func(t *testing.T) {
		type State struct {
			Value    int
			Constant string
		}

		result := F.Pipe2(
			Do(State{Value: 42}),
			LetTo(
				func(c string) func(State) State {
					return func(s State) State { s.Constant = c; return s }
				},
				"hello",
			),
			Map(func(s State) string { return s.Constant }),
		)

		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []string{"hello"}, collected)
	})
}

func TestBindTo(t *testing.T) {
	t.Run("BindTo initializes state", func(t *testing.T) {
		type State struct {
			Value int
		}

		result := F.Pipe2(
			Of(42),
			BindTo(func(v int) State { return State{Value: v} }),
			Map(func(s State) int { return s.Value }),
		)

		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []int{42}, collected)
	})
}
