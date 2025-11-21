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

package readerresult

import (
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) ReaderResult[context.Context, string] {
	return Of[context.Context]("Doe")
}

func getGivenName(s utils.WithLastName) ReaderResult[context.Context, string] {
	return Of[context.Context]("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do[context.Context](utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[context.Context](utils.GetFullName),
	)

	assert.Equal(t, res(context.Background()), result.Of("John Doe"))
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do[context.Context](utils.Empty),
		ApS(utils.SetLastName, Of[context.Context]("Doe")),
		ApS(utils.SetGivenName, Of[context.Context]("John")),
		Map[context.Context](utils.GetFullName),
	)

	assert.Equal(t, res(context.Background()), result.Of("John Doe"))
}

func TestBindReaderK(t *testing.T) {
	type Env struct {
		ConfigPath string
	}
	type State struct {
		Config string
	}

	// A pure Reader computation
	getConfigPath := func(s State) func(Env) string {
		return func(env Env) string {
			return env.ConfigPath
		}
	}

	res := F.Pipe2(
		Do[Env](State{}),
		BindReaderK(
			func(path string) func(State) State {
				return func(s State) State {
					s.Config = path
					return s
				}
			},
			getConfigPath,
		),
		Map[Env](func(s State) string { return s.Config }),
	)

	env := Env{ConfigPath: "/etc/config"}
	assert.Equal(t, result.Of("/etc/config"), res(env))
}

func TestBindResultK(t *testing.T) {
	type State struct {
		Value       int
		ParsedValue int
	}

	// A Result computation that may fail
	parseValue := func(s State) result.Result[int] {
		if s.Value < 0 {
			return result.Left[int](assert.AnError)
		}
		return result.Of(s.Value * 2)
	}

	t.Run("success case", func(t *testing.T) {
		res := F.Pipe2(
			Do[context.Context](State{Value: 5}),
			BindResultK[context.Context](
				func(parsed int) func(State) State {
					return func(s State) State {
						s.ParsedValue = parsed
						return s
					}
				},
				parseValue,
			),
			Map[context.Context](func(s State) int { return s.ParsedValue }),
		)

		assert.Equal(t, result.Of(10), res(context.Background()))
	})

	t.Run("error case", func(t *testing.T) {
		res := F.Pipe2(
			Do[context.Context](State{Value: -5}),
			BindResultK[context.Context](
				func(parsed int) func(State) State {
					return func(s State) State {
						s.ParsedValue = parsed
						return s
					}
				},
				parseValue,
			),
			Map[context.Context](func(s State) int { return s.ParsedValue }),
		)

		assert.True(t, result.IsLeft(res(context.Background())))
	})
}

func TestBindToReader(t *testing.T) {
	type Env struct {
		ConfigPath string
	}
	type State struct {
		Config string
	}

	// A Reader that just reads from the environment
	getConfigPath := func(env Env) string {
		return env.ConfigPath
	}

	res := F.Pipe2(
		getConfigPath,
		BindToReader[Env](func(path string) State {
			return State{Config: path}
		}),
		Map[Env](func(s State) string { return s.Config }),
	)

	env := Env{ConfigPath: "/etc/config"}
	assert.Equal(t, result.Of("/etc/config"), res(env))
}

func TestBindToResult(t *testing.T) {
	type State struct {
		Value int
	}

	t.Run("success case", func(t *testing.T) {
		parseResult := result.Of(42)

		computation := F.Pipe2(
			parseResult,
			BindToResult[context.Context](func(value int) State {
				return State{Value: value}
			}),
			Map[context.Context](func(s State) int { return s.Value }),
		)

		assert.Equal(t, result.Of(42), computation(context.Background()))
	})

	t.Run("error case", func(t *testing.T) {
		parseResult := result.Left[int](assert.AnError)

		computation := F.Pipe2(
			parseResult,
			BindToResult[context.Context](func(value int) State {
				return State{Value: value}
			}),
			Map[context.Context](func(s State) int { return s.Value }),
		)

		assert.True(t, result.IsLeft(computation(context.Background())))
	})
}

func TestApReaderS(t *testing.T) {
	type Env struct {
		DefaultPort int
		DefaultHost string
	}
	type State struct {
		Port int
		Host string
	}

	getDefaultPort := func(env Env) int { return env.DefaultPort }
	getDefaultHost := func(env Env) string { return env.DefaultHost }

	res := F.Pipe3(
		Do[Env](State{}),
		ApReaderS(
			func(port int) func(State) State {
				return func(s State) State {
					s.Port = port
					return s
				}
			},
			getDefaultPort,
		),
		ApReaderS(
			func(host string) func(State) State {
				return func(s State) State {
					s.Host = host
					return s
				}
			},
			getDefaultHost,
		),
		Map[Env](func(s State) State { return s }),
	)

	env := Env{DefaultPort: 8080, DefaultHost: "localhost"}
	r := res(env)
	assert.True(t, result.IsRight(r))
	state := result.GetOrElse(func(error) State { return State{} })(r)
	assert.Equal(t, 8080, state.Port)
	assert.Equal(t, "localhost", state.Host)
}

func TestApResultS(t *testing.T) {
	type State struct {
		Value1 int
		Value2 int
	}

	t.Run("success case", func(t *testing.T) {
		parseValue1 := result.Of(42)
		parseValue2 := result.Of(100)

		computation := F.Pipe3(
			Do[context.Context](State{}),
			ApResultS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
				parseValue1,
			),
			ApResultS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value2 = v
						return s
					}
				},
				parseValue2,
			),
			Map[context.Context](func(s State) State { return s }),
		)

		r := computation(context.Background())
		assert.True(t, result.IsRight(r))
		state := result.GetOrElse(func(error) State { return State{} })(r)
		assert.Equal(t, 42, state.Value1)
		assert.Equal(t, 100, state.Value2)
	})

	t.Run("error in first value", func(t *testing.T) {
		parseValue1 := result.Left[int](assert.AnError)
		parseValue2 := result.Of(100)

		computation := F.Pipe3(
			Do[context.Context](State{}),
			ApResultS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
				parseValue1,
			),
			ApResultS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value2 = v
						return s
					}
				},
				parseValue2,
			),
			Map[context.Context](func(s State) State { return s }),
		)

		assert.True(t, result.IsLeft(computation(context.Background())))
	})

	t.Run("error in second value", func(t *testing.T) {
		parseValue1 := result.Of(42)
		parseValue2 := result.Left[int](assert.AnError)

		computation := F.Pipe3(
			Do[context.Context](State{}),
			ApResultS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
				parseValue1,
			),
			ApResultS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value2 = v
						return s
					}
				},
				parseValue2,
			),
			Map[context.Context](func(s State) State { return s }),
		)

		assert.True(t, result.IsLeft(computation(context.Background())))
	})
}
