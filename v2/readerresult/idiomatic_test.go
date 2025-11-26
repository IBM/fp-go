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
	RRI "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	"github.com/IBM/fp-go/v2/internal/utils"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

var (
	idiomaticTestError = errors.New("idiomatic test error")
)

// TestFromResultI tests the FromResultI function
func TestFromResultI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		rr := FromResultI[MyContext](42, nil)
		assert.Equal(t, result.Of(42), rr(defaultContext))
	})

	t.Run("error case", func(t *testing.T) {
		rr := FromResultI[MyContext](0, idiomaticTestError)
		assert.Equal(t, result.Left[int](idiomaticTestError), rr(defaultContext))
	})
}

// TestFromReaderResultI tests the FromReaderResultI function
func TestFromReaderResultI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		idiomaticRR := func(ctx MyContext) (int, error) {
			return 42, nil
		}
		rr := FromReaderResultI(idiomaticRR)
		assert.Equal(t, result.Of(42), rr(defaultContext))
	})

	t.Run("error case", func(t *testing.T) {
		idiomaticRR := func(ctx MyContext) (int, error) {
			return 0, idiomaticTestError
		}
		rr := FromReaderResultI(idiomaticRR)
		assert.Equal(t, result.Left[int](idiomaticTestError), rr(defaultContext))
	})
}

// TestMonadChainI tests the MonadChainI function
func TestMonadChainI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		idiomaticKleisli := func(x int) RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return x + 1, nil
			}
		}

		rr := Of[MyContext](5)
		res := MonadChainI(rr, idiomaticKleisli)
		assert.Equal(t, result.Of(6), res(defaultContext))
	})

	t.Run("error in first computation", func(t *testing.T) {
		idiomaticKleisli := func(x int) RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return x + 1, nil
			}
		}

		rr := Left[MyContext, int](idiomaticTestError)
		res := MonadChainI(rr, idiomaticKleisli)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})

	t.Run("error in second computation", func(t *testing.T) {
		idiomaticKleisli := func(x int) RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 0, idiomaticTestError
			}
		}

		rr := Of[MyContext](5)
		res := MonadChainI(rr, idiomaticKleisli)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})
}

// TestChainI tests the ChainI function
func TestChainI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		idiomaticKleisli := func(x int) RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return x * 2, nil
			}
		}

		res := F.Pipe1(
			Of[MyContext](5),
			ChainI(idiomaticKleisli),
		)
		assert.Equal(t, result.Of(10), res(defaultContext))
	})

	t.Run("error case", func(t *testing.T) {
		idiomaticKleisli := func(x int) RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 0, idiomaticTestError
			}
		}

		res := F.Pipe1(
			Of[MyContext](5),
			ChainI(idiomaticKleisli),
		)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})
}

// TestMonadApI tests the MonadApI function
func TestMonadApI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		add := func(x int) func(int) int {
			return func(y int) int { return x + y }
		}
		fabr := Of[MyContext](add(5))
		fa := func(ctx MyContext) (int, error) {
			return 3, nil
		}
		res := MonadApI(fabr, fa)
		assert.Equal(t, result.Of(8), res(defaultContext))
	})

	t.Run("error in idiomatic computation", func(t *testing.T) {
		add := func(x int) func(int) int {
			return func(y int) int { return x + y }
		}
		fabr := Of[MyContext](add(5))
		fa := func(ctx MyContext) (int, error) {
			return 0, idiomaticTestError
		}
		res := MonadApI(fabr, fa)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})
}

// TestApI tests the ApI function
func TestApI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		fa := func(ctx MyContext) (int, error) {
			return 10, nil
		}

		res := F.Pipe1(
			Of[MyContext](utils.Double),
			ApI[int](fa),
		)
		assert.Equal(t, result.Of(20), res(defaultContext))
	})
}

// TestOrElseI tests the OrElseI function
func TestOrElseI(t *testing.T) {
	t.Run("success case - doesn't use fallback", func(t *testing.T) {
		fallback := func(err error) RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 99, nil
			}
		}

		res := F.Pipe1(
			Of[MyContext](42),
			OrElseI(fallback),
		)
		assert.Equal(t, result.Of(42), res(defaultContext))
	})

	t.Run("error case - uses fallback", func(t *testing.T) {
		fallback := func(err error) RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 99, nil
			}
		}

		res := F.Pipe1(
			Left[MyContext, int](idiomaticTestError),
			OrElseI(fallback),
		)
		assert.Equal(t, result.Of(99), res(defaultContext))
	})

	t.Run("error case - fallback also fails", func(t *testing.T) {
		fallbackError := errors.New("fallback error")
		fallback := func(err error) RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 0, fallbackError
			}
		}

		res := F.Pipe1(
			Left[MyContext, int](idiomaticTestError),
			OrElseI(fallback),
		)
		assert.Equal(t, result.Left[int](fallbackError), res(defaultContext))
	})
}

// TestMonadChainEitherIK tests the MonadChainEitherIK function
func TestMonadChainEitherIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		idiomaticKleisli := func(x int) (int, error) {
			return x * 2, nil
		}

		rr := Of[MyContext](5)
		res := MonadChainEitherIK(rr, idiomaticKleisli)
		assert.Equal(t, result.Of(10), res(defaultContext))
	})

	t.Run("error case", func(t *testing.T) {
		idiomaticKleisli := func(x int) (int, error) {
			return 0, idiomaticTestError
		}

		rr := Of[MyContext](5)
		res := MonadChainEitherIK(rr, idiomaticKleisli)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})
}

// TestChainEitherIK tests the ChainEitherIK function
func TestChainEitherIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		idiomaticKleisli := func(x int) (int, error) {
			if x < 0 {
				return 0, errors.New("negative value")
			}
			return x * 2, nil
		}

		res := F.Pipe1(
			Of[MyContext](5),
			ChainEitherIK[MyContext](idiomaticKleisli),
		)
		assert.Equal(t, result.Of(10), res(defaultContext))
	})

	t.Run("error case", func(t *testing.T) {
		idiomaticKleisli := func(x int) (int, error) {
			if x < 0 {
				return 0, errors.New("negative value")
			}
			return x * 2, nil
		}

		res := F.Pipe1(
			Of[MyContext](-5),
			ChainEitherIK[MyContext](idiomaticKleisli),
		)
		assert.True(t, result.IsLeft(res(defaultContext)))
	})
}

// TestChainOptionIK tests the ChainOptionIK function
func TestChainOptionIK(t *testing.T) {
	t.Run("success case - Some", func(t *testing.T) {
		idiomaticKleisli := func(x int) (int, bool) {
			if x%2 == 0 {
				return x, true
			}
			return 0, false
		}

		notFound := func() error { return errors.New("not even") }

		res := F.Pipe1(
			Of[MyContext](4),
			ChainOptionIK[MyContext, int, int](notFound)(idiomaticKleisli),
		)
		assert.Equal(t, result.Of(4), res(defaultContext))
	})

	t.Run("None case", func(t *testing.T) {
		idiomaticKleisli := func(x int) (int, bool) {
			if x%2 == 0 {
				return x, true
			}
			return 0, false
		}

		notFound := func() error { return errors.New("not even") }

		res := F.Pipe1(
			Of[MyContext](3),
			ChainOptionIK[MyContext, int, int](notFound)(idiomaticKleisli),
		)
		assert.True(t, result.IsLeft(res(defaultContext)))
	})
}

// TestFlattenI tests the FlattenI function
func TestFlattenI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		nested := Of[MyContext](func(ctx MyContext) (int, error) {
			return 42, nil
		})
		flat := FlattenI(nested)
		assert.Equal(t, result.Of(42), flat(defaultContext))
	})

	t.Run("error case", func(t *testing.T) {
		nested := Of[MyContext](func(ctx MyContext) (int, error) {
			return 0, idiomaticTestError
		})
		flat := FlattenI(nested)
		assert.Equal(t, result.Left[int](idiomaticTestError), flat(defaultContext))
	})
}

// TestMonadAltI tests the MonadAltI function
func TestMonadAltI(t *testing.T) {
	t.Run("first succeeds", func(t *testing.T) {
		first := Of[MyContext](42)
		alternative := func() RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 99, nil
			}
		}

		res := MonadAltI(first, alternative)
		assert.Equal(t, result.Of(42), res(defaultContext))
	})

	t.Run("first fails, second succeeds", func(t *testing.T) {
		first := Left[MyContext, int](idiomaticTestError)
		alternative := func() RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 99, nil
			}
		}

		res := MonadAltI(first, alternative)
		assert.Equal(t, result.Of(99), res(defaultContext))
	})

	t.Run("both fail", func(t *testing.T) {
		first := Left[MyContext, int](idiomaticTestError)
		altError := errors.New("alternative error")
		alternative := func() RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 0, altError
			}
		}

		res := MonadAltI(first, alternative)
		assert.Equal(t, result.Left[int](altError), res(defaultContext))
	})
}

// TestAltI tests the AltI function
func TestAltI(t *testing.T) {
	t.Run("first succeeds", func(t *testing.T) {
		alternative := func() RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 99, nil
			}
		}

		res := F.Pipe1(
			Of[MyContext](42),
			AltI(alternative),
		)
		assert.Equal(t, result.Of(42), res(defaultContext))
	})

	t.Run("first fails, second succeeds", func(t *testing.T) {
		alternative := func() RRI.ReaderResult[MyContext, int] {
			return func(ctx MyContext) (int, error) {
				return 99, nil
			}
		}

		res := F.Pipe1(
			Left[MyContext, int](idiomaticTestError),
			AltI(alternative),
		)
		assert.Equal(t, result.Of(99), res(defaultContext))
	})
}

// TestBindI tests the BindI function
func TestBindI(t *testing.T) {
	type State struct {
		LastName  string
		GivenName string
	}

	t.Run("success case", func(t *testing.T) {
		getLastName := func(s State) RRI.ReaderResult[context.Context, string] {
			return func(ctx context.Context) (string, error) {
				return "Doe", nil
			}
		}

		res := F.Pipe2(
			Do[context.Context](State{}),
			BindI(
				func(name string) func(State) State {
					return func(s State) State {
						s.LastName = name
						return s
					}
				},
				getLastName,
			),
			Map[context.Context](func(s State) string { return s.LastName }),
		)

		assert.Equal(t, result.Of("Doe"), res(context.Background()))
	})

	t.Run("error case", func(t *testing.T) {
		getLastName := func(s State) RRI.ReaderResult[context.Context, string] {
			return func(ctx context.Context) (string, error) {
				return "", idiomaticTestError
			}
		}

		res := F.Pipe2(
			Do[context.Context](State{}),
			BindI(
				func(name string) func(State) State {
					return func(s State) State {
						s.LastName = name
						return s
					}
				},
				getLastName,
			),
			Map[context.Context](func(s State) string { return s.LastName }),
		)

		assert.True(t, result.IsLeft(res(context.Background())))
	})
}

// TestApIS tests the ApIS function
func TestApIS(t *testing.T) {
	type State struct {
		LastName  string
		GivenName string
	}

	t.Run("success case", func(t *testing.T) {
		getLastName := func(ctx context.Context) (string, error) {
			return "Doe", nil
		}

		res := F.Pipe2(
			Do[context.Context](State{}),
			ApIS(
				func(name string) func(State) State {
					return func(s State) State {
						s.LastName = name
						return s
					}
				},
				getLastName,
			),
			Map[context.Context](func(s State) string { return s.LastName }),
		)

		assert.Equal(t, result.Of("Doe"), res(context.Background()))
	})

	t.Run("error case", func(t *testing.T) {
		getLastName := func(ctx context.Context) (string, error) {
			return "", idiomaticTestError
		}

		res := F.Pipe2(
			Do[context.Context](State{}),
			ApIS(
				func(name string) func(State) State {
					return func(s State) State {
						s.LastName = name
						return s
					}
				},
				getLastName,
			),
			Map[context.Context](func(s State) string { return s.LastName }),
		)

		assert.True(t, result.IsLeft(res(context.Background())))
	})
}

// TestApISL tests the ApISL function
func TestApISL(t *testing.T) {
	type State struct {
		LastName  string
		GivenName string
	}

	lastNameLens := L.MakeLens(
		func(s State) string { return s.LastName },
		func(s State, name string) State {
			s.LastName = name
			return s
		},
	)

	t.Run("success case", func(t *testing.T) {
		getLastName := func(ctx context.Context) (string, error) {
			return "Doe", nil
		}

		res := F.Pipe2(
			Do[context.Context](State{}),
			ApISL(lastNameLens, getLastName),
			Map[context.Context](func(s State) string { return s.LastName }),
		)

		assert.Equal(t, result.Of("Doe"), res(context.Background()))
	})

	t.Run("error case", func(t *testing.T) {
		getLastName := func(ctx context.Context) (string, error) {
			return "", idiomaticTestError
		}

		res := F.Pipe2(
			Do[context.Context](State{}),
			ApISL(lastNameLens, getLastName),
			Map[context.Context](func(s State) string { return s.LastName }),
		)

		assert.True(t, result.IsLeft(res(context.Background())))
	})
}

// TestBindIL tests the BindIL function
func TestBindIL(t *testing.T) {
	type State struct {
		Value int
	}

	valueLens := L.MakeLens(
		func(s State) int { return s.Value },
		func(s State, v int) State {
			s.Value = v
			return s
		},
	)

	t.Run("success case", func(t *testing.T) {
		doubleValue := func(v int) RRI.ReaderResult[context.Context, int] {
			return func(ctx context.Context) (int, error) {
				return v * 2, nil
			}
		}

		res := F.Pipe2(
			Do[context.Context](State{Value: 21}),
			BindIL(valueLens, doubleValue),
			Map[context.Context](func(s State) int { return s.Value }),
		)

		assert.Equal(t, result.Of(42), res(context.Background()))
	})

	t.Run("error case", func(t *testing.T) {
		failValue := func(v int) RRI.ReaderResult[context.Context, int] {
			return func(ctx context.Context) (int, error) {
				return 0, idiomaticTestError
			}
		}

		res := F.Pipe2(
			Do[context.Context](State{Value: 21}),
			BindIL(valueLens, failValue),
			Map[context.Context](func(s State) int { return s.Value }),
		)

		assert.True(t, result.IsLeft(res(context.Background())))
	})
}

// TestBindEitherIK tests the BindEitherIK function
func TestBindEitherIK(t *testing.T) {
	type State struct {
		Value       int
		ParsedValue int
	}

	t.Run("success case", func(t *testing.T) {
		parseValue := func(s State) (int, error) {
			if s.Value < 0 {
				return 0, errors.New("negative value")
			}
			return s.Value * 2, nil
		}

		res := F.Pipe2(
			Do[context.Context](State{Value: 5}),
			BindEitherIK[context.Context](
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
		parseValue := func(s State) (int, error) {
			if s.Value < 0 {
				return 0, errors.New("negative value")
			}
			return s.Value * 2, nil
		}

		res := F.Pipe2(
			Do[context.Context](State{Value: -5}),
			BindEitherIK[context.Context](
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

// TestBindResultIK tests the BindResultIK function (alias of BindEitherIK)
func TestBindResultIK(t *testing.T) {
	type State struct {
		Value       int
		ParsedValue int
	}

	t.Run("success case", func(t *testing.T) {
		parseValue := func(s State) (int, error) {
			return s.Value * 2, nil
		}

		res := F.Pipe2(
			Do[context.Context](State{Value: 21}),
			BindResultIK[context.Context](
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

		assert.Equal(t, result.Of(42), res(context.Background()))
	})
}

// TestBindToEitherI tests the BindToEitherI function
func TestBindToEitherI(t *testing.T) {
	type State struct {
		Value int
	}

	t.Run("success case", func(t *testing.T) {
		bindTo := BindToEitherI[context.Context](func(value int) State {
			return State{Value: value}
		})
		computation := F.Pipe1(
			bindTo(42, nil),
			Map[context.Context](func(s State) int { return s.Value }),
		)

		assert.Equal(t, result.Of(42), computation(context.Background()))
	})

	t.Run("error case", func(t *testing.T) {
		bindTo := BindToEitherI[context.Context](func(value int) State {
			return State{Value: value}
		})
		computation := F.Pipe1(
			bindTo(0, idiomaticTestError),
			Map[context.Context](func(s State) int { return s.Value }),
		)

		assert.True(t, result.IsLeft(computation(context.Background())))
	})
}

// TestBindToResultI tests the BindToResultI function (alias of BindToEitherI)
func TestBindToResultI(t *testing.T) {
	type State struct {
		Value int
	}

	t.Run("success case", func(t *testing.T) {
		bindTo := BindToResultI[context.Context](func(value int) State {
			return State{Value: value}
		})
		computation := F.Pipe1(
			bindTo(42, nil),
			Map[context.Context](func(s State) int { return s.Value }),
		)

		assert.Equal(t, result.Of(42), computation(context.Background()))
	})
}

// TestApEitherIS tests the ApEitherIS function
func TestApEitherIS(t *testing.T) {
	type State struct {
		Value1 int
		Value2 int
	}

	t.Run("success case", func(t *testing.T) {
		computation := F.Pipe2(
			Do[context.Context](State{}),
			ApEitherIS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
			)(42, nil),
			Map[context.Context](func(s State) int { return s.Value1 }),
		)

		assert.Equal(t, result.Of(42), computation(context.Background()))
	})

	t.Run("error case", func(t *testing.T) {
		computation := F.Pipe2(
			Do[context.Context](State{}),
			ApEitherIS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value1 = v
						return s
					}
				},
			)(0, idiomaticTestError),
			Map[context.Context](func(s State) int { return s.Value1 }),
		)

		assert.True(t, result.IsLeft(computation(context.Background())))
	})
}

// TestApResultIS tests the ApResultIS function (alias of ApEitherIS)
func TestApResultIS(t *testing.T) {
	type State struct {
		Value int
	}

	t.Run("success case", func(t *testing.T) {
		computation := F.Pipe2(
			Do[context.Context](State{}),
			ApResultIS[context.Context](
				func(v int) func(State) State {
					return func(s State) State {
						s.Value = v
						return s
					}
				},
			)(42, nil),
			Map[context.Context](func(s State) int { return s.Value }),
		)

		assert.Equal(t, result.Of(42), computation(context.Background()))
	})
}

// TestMonadApResult tests the MonadApResult function
func TestMonadApResult(t *testing.T) {
	t.Run("success case - both succeed", func(t *testing.T) {
		add := func(x int) func(int) int {
			return func(y int) int { return x + y }
		}
		fabr := Of[MyContext](add(5))
		fa := result.Of(3)
		res := MonadApResult(fabr, fa)
		assert.Equal(t, result.Of(8), res(defaultContext))
	})

	t.Run("function is error", func(t *testing.T) {
		fabr := Left[MyContext, func(int) int](idiomaticTestError)
		fa := result.Of(3)
		res := MonadApResult(fabr, fa)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})

	t.Run("value is error", func(t *testing.T) {
		double := N.Mul(2)
		fabr := Of[MyContext](double)
		fa := result.Left[int](idiomaticTestError)
		res := MonadApResult(fabr, fa)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})

	t.Run("both are errors", func(t *testing.T) {
		funcError := errors.New("function error")
		valueError := errors.New("value error")
		fabr := Left[MyContext, func(int) int](funcError)
		fa := result.Left[int](valueError)
		res := MonadApResult(fabr, fa)
		// When both fail, the function error takes precedence in Applicative semantics
		assert.True(t, result.IsLeft(res(defaultContext)))
	})
}

// TestApResult tests the ApResult function
func TestApResult(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		fa := result.Of(10)
		res := F.Pipe1(
			Of[MyContext](utils.Double),
			ApResult[int, MyContext](fa),
		)
		assert.Equal(t, result.Of(20), res(defaultContext))
	})

	t.Run("function error", func(t *testing.T) {
		fa := result.Of(10)
		res := F.Pipe1(
			Left[MyContext, func(int) int](idiomaticTestError),
			ApResult[int, MyContext](fa),
		)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})

	t.Run("value error", func(t *testing.T) {
		fa := result.Left[int](idiomaticTestError)
		res := F.Pipe1(
			Of[MyContext](utils.Double),
			ApResult[int, MyContext](fa),
		)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})

	t.Run("with triple composition", func(t *testing.T) {
		triple := N.Mul(3)
		fa := result.Of(7)
		res := F.Pipe1(
			Of[MyContext](triple),
			ApResult[int, MyContext](fa),
		)
		assert.Equal(t, result.Of(21), res(defaultContext))
	})
}

// TestApResultI tests the ApResultI function
func TestApResultI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		value := 10
		var err error = nil
		res := F.Pipe1(
			Of[MyContext](utils.Double),
			ApResultI[int, MyContext](value, err),
		)
		assert.Equal(t, result.Of(20), res(defaultContext))
	})

	t.Run("function error", func(t *testing.T) {
		value := 10
		var err error = nil
		res := F.Pipe1(
			Left[MyContext, func(int) int](idiomaticTestError),
			ApResultI[int, MyContext](value, err),
		)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})

	t.Run("value error", func(t *testing.T) {
		value := 0
		err := idiomaticTestError
		res := F.Pipe1(
			Of[MyContext](utils.Double),
			ApResultI[int, MyContext](value, err),
		)
		assert.Equal(t, result.Left[int](idiomaticTestError), res(defaultContext))
	})

	t.Run("realistic example with strconv", func(t *testing.T) {
		// Simulate parsing a string to int
		parseValue := func(s string) (int, error) {
			if s == "42" {
				return 42, nil
			}
			return 0, errors.New("parse error")
		}

		addTen := N.Add(10)

		t.Run("parse success", func(t *testing.T) {
			value, err := parseValue("42")
			res := F.Pipe1(
				Of[MyContext](addTen),
				ApResultI[int, MyContext](value, err),
			)
			assert.Equal(t, result.Of(52), res(defaultContext))
		})

		t.Run("parse error", func(t *testing.T) {
			value, err := parseValue("invalid")
			res := F.Pipe1(
				Of[MyContext](addTen),
				ApResultI[int, MyContext](value, err),
			)
			assert.True(t, result.IsLeft(res(defaultContext)))
		})
	})

	t.Run("with curried function", func(t *testing.T) {
		// Test with a curried addition function
		add := func(x int) func(int) int {
			return func(y int) int { return x + y }
		}

		// First apply 5, get a function (int -> int)
		partialAdd := F.Pipe1(
			Of[MyContext](add),
			ApResultI[func(int) int, MyContext](5, nil),
		)

		// Then apply 3 to the result
		finalResult := F.Pipe1(
			partialAdd,
			ApResultI[int, MyContext](3, nil),
		)

		assert.Equal(t, result.Of(8), finalResult(defaultContext))
	})
}

// Test a complex scenario combining multiple idiomatic functions
func TestComplexIdiomaticScenario(t *testing.T) {
	type Env struct {
		Multiplier int
	}

	type State struct {
		Input  int
		Result int
	}

	// Idiomatic function that reads from environment
	multiply := func(s State) RRI.ReaderResult[Env, int] {
		return func(env Env) (int, error) {
			if env.Multiplier == 0 {
				return 0, errors.New("multiplier cannot be zero")
			}
			return s.Input * env.Multiplier, nil
		}
	}

	// Idiomatic function that validates
	validate := func(x int) (int, error) {
		if x < 0 {
			return 0, errors.New("result cannot be negative")
		}
		return x, nil
	}

	env := Env{Multiplier: 3}

	t.Run("success case", func(t *testing.T) {
		computation := F.Pipe3(
			Do[Env](State{Input: 10}),
			BindI(
				func(res int) func(State) State {
					return func(s State) State {
						s.Result = res
						return s
					}
				},
				multiply,
			),
			BindEitherIK[Env](
				func(validated int) func(State) State {
					return func(s State) State {
						s.Result = validated
						return s
					}
				},
				func(s State) (int, error) {
					return validate(s.Result)
				},
			),
			Map[Env](func(s State) int { return s.Result }),
		)

		assert.Equal(t, result.Of(30), computation(env))
	})

	t.Run("validation error", func(t *testing.T) {
		computation := F.Pipe3(
			Do[Env](State{Input: -10}),
			BindI(
				func(res int) func(State) State {
					return func(s State) State {
						s.Result = res
						return s
					}
				},
				multiply,
			),
			BindEitherIK[Env](
				func(validated int) func(State) State {
					return func(s State) State {
						s.Result = validated
						return s
					}
				},
				func(s State) (int, error) {
					return validate(s.Result)
				},
			),
			Map[Env](func(s State) int { return s.Result }),
		)

		assert.True(t, result.IsLeft(computation(env)))
	})

	t.Run("multiplier error", func(t *testing.T) {
		badEnv := Env{Multiplier: 0}
		computation := F.Pipe3(
			Do[Env](State{Input: 10}),
			BindI(
				func(res int) func(State) State {
					return func(s State) State {
						s.Result = res
						return s
					}
				},
				multiply,
			),
			BindEitherIK[Env](
				func(validated int) func(State) State {
					return func(s State) State {
						s.Result = validated
						return s
					}
				},
				func(s State) (int, error) {
					return validate(s.Result)
				},
			),
			Map[Env](func(s State) int { return s.Result }),
		)

		assert.True(t, result.IsLeft(computation(badEnv)))
	})
}
