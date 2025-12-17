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
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/reader"
	RES "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestDoInit(t *testing.T) {
	initial := SimpleState{Value: 42}
	result := Do(initial)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, initial, state)
}

func TestBind(t *testing.T) {
	t.Run("successful bind", func(t *testing.T) {
		// Effectful function that depends on context
		fetchValue := func(s SimpleState) ReaderResult[int] {
			return func(ctx context.Context) (int, error) {
				return s.Value * 2, nil
			}
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 21}),
			Bind(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				fetchValue,
			),
		)

		state, err := result(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 42, state.Value)
	})

	t.Run("bind with error", func(t *testing.T) {
		fetchValue := func(s SimpleState) ReaderResult[int] {
			return func(ctx context.Context) (int, error) {
				return 0, errors.New("fetch failed")
			}
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 21}),
			Bind(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				fetchValue,
			),
		)

		_, err := result(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "fetch failed", err.Error())
	})
}

func TestLet(t *testing.T) {
	// Pure function that doesn't depend on context
	double := func(s SimpleState) int {
		return s.Value * 2
	}

	result := F.Pipe1(
		Do(SimpleState{Value: 21}),
		Let(
			func(v int) func(SimpleState) SimpleState {
				return func(s SimpleState) SimpleState {
					s.Value = v
					return s
				}
			},
			double,
		),
	)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 42, state.Value)
}

func TestLetTo(t *testing.T) {
	result := F.Pipe1(
		Do(SimpleState{}),
		LetTo(
			func(v int) func(SimpleState) SimpleState {
				return func(s SimpleState) SimpleState {
					s.Value = v
					return s
				}
			},
			100,
		),
	)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 100, state.Value)
}

func TestBindToInit(t *testing.T) {
	getValue := func(ctx context.Context) (int, error) {
		return 42, nil
	}

	result := F.Pipe1(
		getValue,
		BindTo(func(v int) SimpleState {
			return SimpleState{Value: v}
		}),
	)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 42, state.Value)
}

func TestApS(t *testing.T) {
	t.Run("successful ApS", func(t *testing.T) {
		getValue := func(ctx context.Context) (int, error) {
			return 100, nil
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 42}),
			ApS(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				getValue,
			),
		)

		state, err := result(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 100, state.Value)
	})

	t.Run("ApS with error", func(t *testing.T) {
		getValue := func(ctx context.Context) (int, error) {
			return 0, errors.New("failed")
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 42}),
			ApS(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				getValue,
			),
		)

		_, err := result(context.Background())
		assert.Error(t, err)
	})
}

func TestApSL(t *testing.T) {
	lenses := MakeSimpleStateLenses()

	getValue := func(ctx context.Context) (int, error) {
		return 100, nil
	}

	result := F.Pipe1(
		Do(SimpleState{Value: 42}),
		ApSL(lenses.Value, getValue),
	)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 100, state.Value)
}

func TestBindL(t *testing.T) {
	lenses := MakeSimpleStateLenses()

	// Effectful function
	increment := func(v int) ReaderResult[int] {
		return func(ctx context.Context) (int, error) {
			return v + 1, nil
		}
	}

	result := F.Pipe1(
		Do(SimpleState{Value: 41}),
		BindL(lenses.Value, increment),
	)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 42, state.Value)
}

func TestLetL(t *testing.T) {
	lenses := MakeSimpleStateLenses()

	result := F.Pipe1(
		Do(SimpleState{Value: 21}),
		LetL(lenses.Value, N.Mul(2)),
	)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 42, state.Value)
}

func TestLetToL(t *testing.T) {
	lenses := MakeSimpleStateLenses()

	result := F.Pipe1(
		Do(SimpleState{}),
		LetToL(lenses.Value, 42),
	)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 42, state.Value)
}

func TestBindReaderK(t *testing.T) {
	t.Run("successful BindReaderK", func(t *testing.T) {
		// Context-dependent function that doesn't return error
		getFromContext := func(s SimpleState) reader.Reader[context.Context, int] {
			return func(ctx context.Context) int {
				if val := ctx.Value("multiplier"); val != nil {
					return s.Value * val.(int)
				}
				return s.Value
			}
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 21}),
			BindReaderK(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				getFromContext,
			),
		)

		ctx := context.WithValue(context.Background(), "multiplier", 2)
		state, err := result(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 42, state.Value)
	})
}

func TestBindEitherK(t *testing.T) {
	t.Run("successful BindEitherK", func(t *testing.T) {
		// Pure function returning Result
		validate := func(s SimpleState) RES.Result[int] {
			if s.Value > 0 {
				return RES.Of(s.Value * 2)
			}
			return RES.Left[int](errors.New("value must be positive"))
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 21}),
			BindEitherK(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				validate,
			),
		)

		state, err := result(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 42, state.Value)
	})

	t.Run("BindEitherK with error", func(t *testing.T) {
		validate := func(s SimpleState) RES.Result[int] {
			return RES.Left[int](errors.New("validation failed"))
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 21}),
			BindEitherK(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				validate,
			),
		)

		_, err := result(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "validation failed", err.Error())
	})
}

func TestBindResultK(t *testing.T) {
	t.Run("successful BindResultK", func(t *testing.T) {
		// Pure function returning (value, error)
		parse := func(s SimpleState) (int, error) {
			return s.Value * 2, nil
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 21}),
			BindResultK(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				parse,
			),
		)

		state, err := result(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 42, state.Value)
	})

	t.Run("BindResultK with error", func(t *testing.T) {
		parse := func(s SimpleState) (int, error) {
			return 0, errors.New("parse failed")
		}

		result := F.Pipe1(
			Do(SimpleState{Value: 21}),
			BindResultK(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				parse,
			),
		)

		_, err := result(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "parse failed", err.Error())
	})
}

func TestBindToReader(t *testing.T) {
	getFromContext := func(ctx context.Context) int {
		if val := ctx.Value("value"); val != nil {
			return val.(int)
		}
		return 0
	}

	result := F.Pipe1(
		getFromContext,
		BindToReader(func(v int) SimpleState {
			return SimpleState{Value: v}
		}),
	)

	ctx := context.WithValue(context.Background(), "value", 42)
	state, err := result(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 42, state.Value)
}

func TestBindToEither(t *testing.T) {
	t.Run("successful BindToEither", func(t *testing.T) {
		resultValue := RES.Of(42)

		result := F.Pipe1(
			resultValue,
			BindToEither(func(v int) SimpleState {
				return SimpleState{Value: v}
			}),
		)

		state, err := result(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 42, state.Value)
	})

	t.Run("BindToEither with error", func(t *testing.T) {
		resultValue := RES.Left[int](errors.New("failed"))

		result := F.Pipe1(
			resultValue,
			BindToEither(func(v int) SimpleState {
				return SimpleState{Value: v}
			}),
		)

		_, err := result(context.Background())
		assert.Error(t, err)
	})
}

func TestBindToResult(t *testing.T) {
	t.Run("successful BindToResult", func(t *testing.T) {
		value, err := 42, error(nil)

		result := F.Pipe1(
			BindToResult(func(v int) SimpleState {
				return SimpleState{Value: v}
			}),
			func(f func(int, error) ReaderResult[SimpleState]) ReaderResult[SimpleState] {
				return f(value, err)
			},
		)

		state, resultErr := result(context.Background())
		assert.NoError(t, resultErr)
		assert.Equal(t, 42, state.Value)
	})

	t.Run("BindToResult with error", func(t *testing.T) {
		value, err := 0, errors.New("failed")

		result := F.Pipe1(
			BindToResult(func(v int) SimpleState {
				return SimpleState{Value: v}
			}),
			func(f func(int, error) ReaderResult[SimpleState]) ReaderResult[SimpleState] {
				return f(value, err)
			},
		)

		_, resultErr := result(context.Background())
		assert.Error(t, resultErr)
	})
}

func TestApReaderS(t *testing.T) {
	getFromContext := func(ctx context.Context) int {
		if val := ctx.Value("value"); val != nil {
			return val.(int)
		}
		return 0
	}

	result := F.Pipe1(
		Do(SimpleState{}),
		ApReaderS(
			func(v int) func(SimpleState) SimpleState {
				return func(s SimpleState) SimpleState {
					s.Value = v
					return s
				}
			},
			getFromContext,
		),
	)

	ctx := context.WithValue(context.Background(), "value", 42)
	state, err := result(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 42, state.Value)
}

func TestApResultS(t *testing.T) {
	t.Run("successful ApResultS", func(t *testing.T) {
		value, err := 42, error(nil)

		result := F.Pipe1(
			Do(SimpleState{}),
			func(rr ReaderResult[SimpleState]) ReaderResult[SimpleState] {
				return F.Pipe1(
					rr,
					ApResultS(
						func(v int) func(SimpleState) SimpleState {
							return func(s SimpleState) SimpleState {
								s.Value = v
								return s
							}
						},
					)(value, err),
				)
			},
		)

		state, resultErr := result(context.Background())
		assert.NoError(t, resultErr)
		assert.Equal(t, 42, state.Value)
	})

	t.Run("ApResultS with error", func(t *testing.T) {
		value, err := 0, errors.New("failed")

		result := F.Pipe1(
			Do(SimpleState{}),
			func(rr ReaderResult[SimpleState]) ReaderResult[SimpleState] {
				return F.Pipe1(
					rr,
					ApResultS(
						func(v int) func(SimpleState) SimpleState {
							return func(s SimpleState) SimpleState {
								s.Value = v
								return s
							}
						},
					)(value, err),
				)
			},
		)

		_, resultErr := result(context.Background())
		assert.Error(t, resultErr)
	})
}

func TestApEitherS(t *testing.T) {
	t.Run("successful ApEitherS", func(t *testing.T) {
		resultValue := RES.Of(42)

		result := F.Pipe1(
			Do(SimpleState{}),
			ApEitherS(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				resultValue,
			),
		)

		state, err := result(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 42, state.Value)
	})

	t.Run("ApEitherS with error", func(t *testing.T) {
		resultValue := RES.Left[int](errors.New("failed"))

		result := F.Pipe1(
			Do(SimpleState{}),
			ApEitherS(
				func(v int) func(SimpleState) SimpleState {
					return func(s SimpleState) SimpleState {
						s.Value = v
						return s
					}
				},
				resultValue,
			),
		)

		_, err := result(context.Background())
		assert.Error(t, err)
	})
}

func TestComplexPipeline(t *testing.T) {
	lenses := MakeSimpleStateLenses()

	// Complex pipeline combining multiple operations
	result := F.Pipe3(
		Do(SimpleState{}),
		LetToL(lenses.Value, 10),
		LetL(lenses.Value, N.Mul(2)),
		BindResultK(
			func(v int) func(SimpleState) SimpleState {
				return func(s SimpleState) SimpleState {
					s.Value = v
					return s
				}
			},
			func(s SimpleState) (int, error) {
				return s.Value + 22, nil
			},
		),
	)

	state, err := result(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 42, state.Value)
}
