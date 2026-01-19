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

package readerreaderioresult

import (
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type AppConfig struct {
	DatabaseURL string
	LogLevel    string
}

var defaultConfig = AppConfig{
	DatabaseURL: "postgres://localhost",
	LogLevel:    "info",
}

func getLastName(s utils.Initial) ReaderReaderIOResult[AppConfig, string] {
	return Of[AppConfig]("Doe")
}

func getGivenName(s utils.WithLastName) ReaderReaderIOResult[AppConfig, string] {
	return Of[AppConfig]("John")
}

func TestDo(t *testing.T) {
	res := Do[AppConfig](utils.Empty)
	outcome := res(defaultConfig)(t.Context())()

	assert.True(t, result.IsRight(outcome))
	assert.Equal(t, result.Of(utils.Empty), outcome)
}

func TestBind(t *testing.T) {
	res := F.Pipe3(
		Do[AppConfig](utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[AppConfig](utils.GetFullName),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of("John Doe"), outcome)
}

func TestBindWithError(t *testing.T) {
	testErr := errors.New("bind error")

	getLastNameErr := func(s utils.Initial) ReaderReaderIOResult[AppConfig, string] {
		return Left[AppConfig, string](testErr)
	}

	res := F.Pipe3(
		Do[AppConfig](utils.Empty),
		Bind(utils.SetLastName, getLastNameErr),
		Bind(utils.SetGivenName, getGivenName),
		Map[AppConfig](utils.GetFullName),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.True(t, result.IsLeft(outcome))
}

func TestLet(t *testing.T) {
	type State struct {
		FirstName string
		LastName  string
		FullName  string
	}

	res := F.Pipe2(
		Do[AppConfig](State{FirstName: "John", LastName: "Doe"}),
		Let[AppConfig](
			func(fullName string) func(State) State {
				return func(s State) State {
					s.FullName = fullName
					return s
				}
			},
			func(s State) string {
				return s.FirstName + " " + s.LastName
			},
		),
		Map[AppConfig](func(s State) string { return s.FullName }),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of("John Doe"), outcome)
}

func TestLetTo(t *testing.T) {
	type State struct {
		Status string
	}

	res := F.Pipe2(
		Do[AppConfig](State{}),
		LetTo[AppConfig](
			func(status string) func(State) State {
				return func(s State) State {
					s.Status = status
					return s
				}
			},
			"active",
		),
		Map[AppConfig](func(s State) string { return s.Status }),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of("active"), outcome)
}

func TestBindTo(t *testing.T) {
	type State struct {
		Count int
	}

	res := F.Pipe2(
		Of[AppConfig](42),
		BindTo[AppConfig](func(n int) State { return State{Count: n} }),
		Map[AppConfig](func(s State) int { return s.Count }),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestApS(t *testing.T) {
	res := F.Pipe3(
		Do[AppConfig](utils.Empty),
		ApS(utils.SetLastName, Of[AppConfig]("Doe")),
		ApS(utils.SetGivenName, Of[AppConfig]("John")),
		Map[AppConfig](utils.GetFullName),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of("John Doe"), outcome)
}

func TestApSWithError(t *testing.T) {
	testErr := errors.New("aps error")

	res := F.Pipe3(
		Do[AppConfig](utils.Empty),
		ApS(utils.SetLastName, Left[AppConfig, string](testErr)),
		ApS(utils.SetGivenName, Of[AppConfig]("John")),
		Map[AppConfig](utils.GetFullName),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.True(t, result.IsLeft(outcome))
}

func TestBindReaderK(t *testing.T) {
	type State struct {
		Config string
	}

	getConfigPath := func(s State) func(AppConfig) string {
		return func(cfg AppConfig) string {
			return cfg.DatabaseURL
		}
	}

	res := F.Pipe2(
		Do[AppConfig](State{}),
		BindReaderK(
			func(path string) func(State) State {
				return func(s State) State {
					s.Config = path
					return s
				}
			},
			getConfigPath,
		),
		Map[AppConfig](func(s State) string { return s.Config }),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of("postgres://localhost"), outcome)
}

func TestBindIOResultK(t *testing.T) {
	type State struct {
		Value       int
		ParsedValue int
	}

	parseValue := func(s State) ioresult.IOResult[int] {
		return func() result.Result[int] {
			if s.Value < 0 {
				return result.Left[int](errors.New("negative value"))
			}
			return result.Of(s.Value * 2)
		}
	}

	t.Run("success case", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 5}),
			BindIOResultK[AppConfig](
				func(parsed int) func(State) State {
					return func(s State) State {
						s.ParsedValue = parsed
						return s
					}
				},
				parseValue,
			),
			Map[AppConfig](func(s State) int { return s.ParsedValue }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(10), outcome)
	})

	t.Run("error case", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: -5}),
			BindIOResultK[AppConfig](
				func(parsed int) func(State) State {
					return func(s State) State {
						s.ParsedValue = parsed
						return s
					}
				},
				parseValue,
			),
			Map[AppConfig](func(s State) int { return s.ParsedValue }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestBindIOK(t *testing.T) {
	type State struct {
		Value int
	}

	getValue := func(s State) io.IO[int] {
		return func() int {
			return s.Value * 2
		}
	}

	res := F.Pipe2(
		Do[AppConfig](State{Value: 21}),
		BindIOK[AppConfig](
			func(v int) func(State) State {
				return func(s State) State {
					s.Value = v
					return s
				}
			},
			getValue,
		),
		Map[AppConfig](func(s State) int { return s.Value }),
	)

	outcome := res(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestBindReaderIOK(t *testing.T) {
	type State struct {
		Value int
	}

	getValue := func(s State) readerio.ReaderIO[AppConfig, int] {
		return func(cfg AppConfig) io.IO[int] {
			return func() int {
				return s.Value + len(cfg.DatabaseURL)
			}
		}
	}

	res := F.Pipe2(
		Do[AppConfig](State{Value: 10}),
		BindReaderIOK(
			func(v int) func(State) State {
				return func(s State) State {
					s.Value = v
					return s
				}
			},
			getValue,
		),
		Map[AppConfig](func(s State) int { return s.Value }),
	)

	outcome := res(defaultConfig)(t.Context())()
	// 10 + len("postgres://localhost") = 10 + 20 = 30
	assert.Equal(t, result.Of(30), outcome)
}

func TestBindEitherK(t *testing.T) {
	type State struct {
		Value int
	}

	parseValue := func(s State) either.Either[error, int] {
		if s.Value < 0 {
			return either.Left[int](errors.New("negative"))
		}
		return either.Right[error](s.Value * 2)
	}

	t.Run("success case", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 5}),
			BindEitherK[AppConfig](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value = v
						return s
					}
				},
				parseValue,
			),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(10), outcome)
	})

	t.Run("error case", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: -5}),
			BindEitherK[AppConfig](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value = v
						return s
					}
				},
				parseValue,
			),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestBindIOEitherK(t *testing.T) {
	type State struct {
		Value int
	}

	parseValue := func(s State) ioeither.IOEither[error, int] {
		return func() either.Either[error, int] {
			if s.Value < 0 {
				return either.Left[int](errors.New("negative"))
			}
			return either.Right[error](s.Value * 2)
		}
	}

	t.Run("success case", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 5}),
			BindIOEitherK[AppConfig](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value = v
						return s
					}
				},
				parseValue,
			),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(10), outcome)
	})

	t.Run("error case", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: -5}),
			BindIOEitherK[AppConfig](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value = v
						return s
					}
				},
				parseValue,
			),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestLensOperations(t *testing.T) {
	type State struct {
		Value int
	}

	valueLens := lens.MakeLens(
		func(s State) int { return s.Value },
		func(s State, v int) State {
			s.Value = v
			return s
		},
	)

	t.Run("ApSL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApSL(valueLens, Of[AppConfig](42)),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("BindL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 10}),
			BindL(valueLens, func(v int) ReaderReaderIOResult[AppConfig, int] {
				return Of[AppConfig](v * 2)
			}),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(20), outcome)
	})

	t.Run("LetL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 10}),
			LetL[AppConfig](valueLens, func(v int) int { return v * 3 }),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(30), outcome)
	})

	t.Run("LetToL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			LetToL[AppConfig](valueLens, 99),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(99), outcome)
	})

	t.Run("BindIOEitherKL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 5}),
			BindIOEitherKL[AppConfig](valueLens, func(v int) ioeither.IOEither[error, int] {
				return ioeither.Of[error](v * 2)
			}),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(10), outcome)
	})

	t.Run("BindIOKL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 7}),
			BindIOKL[AppConfig](valueLens, func(v int) io.IO[int] {
				return func() int { return v * 3 }
			}),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(21), outcome)
	})

	t.Run("BindReaderKL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 5}),
			BindReaderKL(valueLens, func(v int) reader.Reader[AppConfig, int] {
				return func(cfg AppConfig) int {
					return v + len(cfg.LogLevel)
				}
			}),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		// 5 + len("info") = 5 + 4 = 9
		assert.Equal(t, result.Of(9), outcome)
	})

	t.Run("BindReaderIOKL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{Value: 10}),
			BindReaderIOKL(valueLens, func(v int) readerio.ReaderIO[AppConfig, int] {
				return func(cfg AppConfig) io.IO[int] {
					return func() int {
						return v + len(cfg.DatabaseURL)
					}
				}
			}),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		// 10 + len("postgres://localhost") = 10 + 20 = 30
		assert.Equal(t, result.Of(30), outcome)
	})

	t.Run("ApIOEitherSL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApIOEitherSL[AppConfig](valueLens, ioeither.Of[error](42)),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("ApIOSL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApIOSL[AppConfig](valueLens, func() int { return 99 }),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(99), outcome)
	})

	t.Run("ApReaderSL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApReaderSL(valueLens, func(cfg AppConfig) int {
				return len(cfg.LogLevel)
			}),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(4), outcome)
	})

	t.Run("ApReaderIOSL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApReaderIOSL(valueLens, func(cfg AppConfig) io.IO[int] {
				return func() int { return len(cfg.DatabaseURL) }
			}),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(20), outcome)
	})

	t.Run("ApEitherSL", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApEitherSL[AppConfig](valueLens, either.Right[error](77)),
			Map[AppConfig](func(s State) int { return s.Value }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(77), outcome)
	})
}

func TestApOperations(t *testing.T) {
	type State struct {
		Value1 int
		Value2 int
	}

	t.Run("ApIOEitherS", func(t *testing.T) {
		res := F.Pipe3(
			Do[AppConfig](State{}),
			ApIOEitherS[AppConfig](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
				ioeither.Of[error](10),
			),
			ApIOEitherS[AppConfig](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value2 = v
						return s
					}
				},
				ioeither.Of[error](20),
			),
			Map[AppConfig](func(s State) int { return s.Value1 + s.Value2 }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(30), outcome)
	})

	t.Run("ApIOS", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApIOS[AppConfig](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
				func() int { return 42 },
			),
			Map[AppConfig](func(s State) int { return s.Value1 }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("ApReaderS", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApReaderS(
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
				func(cfg AppConfig) int { return len(cfg.LogLevel) },
			),
			Map[AppConfig](func(s State) int { return s.Value1 }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(4), outcome)
	})

	t.Run("ApReaderIOS", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApReaderIOS(
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
				func(cfg AppConfig) io.IO[int] {
					return func() int { return len(cfg.DatabaseURL) }
				},
			),
			Map[AppConfig](func(s State) int { return s.Value1 }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(20), outcome)
	})

	t.Run("ApEitherS", func(t *testing.T) {
		res := F.Pipe2(
			Do[AppConfig](State{}),
			ApEitherS[AppConfig](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
				either.Right[error](99),
			),
			Map[AppConfig](func(s State) int { return s.Value1 }),
		)

		outcome := res(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(99), outcome)
	})
}
