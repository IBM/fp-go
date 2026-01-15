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

package readerioresult

import (
	"context"
	"fmt"
	"log"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context](1),
		Map[context.Context](utils.Double),
	)

	assert.Equal(t, result.Of(2), g(context.Background())())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Right[context.Context](utils.Double),
		Ap[int](Right[context.Context](1)),
	)

	assert.Equal(t, result.Of(2), g(context.Background())())
}

func TestChainReaderK(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context](1),
		ChainReaderK(func(v int) R.Reader[context.Context, string] {
			return R.Of[context.Context](fmt.Sprintf("%d", v))
		}),
	)

	assert.Equal(t, result.Of("1"), g(context.Background())())
}

func TestTapReaderIOK(t *testing.T) {

	rdr := Of[int]("TestTapReaderIOK")

	x := F.Pipe1(
		rdr,
		TapReaderIOK(func(a string) ReaderIO[int, any] {
			return func(ctx int) IO[any] {
				return func() any {
					log.Printf("Context: %d, Value: %s", ctx, a)
					return nil
				}
			}
		}),
	)

	x(10)()
}

func TestReadIOEither(t *testing.T) {
	type Config struct {
		BaseURL string
	}

	t.Run("success case - environment and computation both succeed", func(t *testing.T) {
		// Create an IOResult that successfully produces a config
		getConfig := func() IOResult[Config] {
			return func() Result[Config] {
				return result.Of(Config{BaseURL: "https://api.example.com"})
			}
		}

		// Create a ReaderIOResult that uses the config
		computation := func(cfg Config) IOResult[string] {
			return func() Result[string] {
				return result.Of(cfg.BaseURL + "/users")
			}
		}

		// Execute using ReadIOEither
		ioResult := ReadIOEither[string](getConfig())(computation)
		res := ioResult()

		assert.True(t, result.IsRight(res))
		assert.Equal(t, "https://api.example.com/users", result.GetOrElse(func(error) string { return "" })(res))
	})

	t.Run("failure case - environment acquisition fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("config load failed")

		// Create an IOResult that fails to produce a config
		getConfig := func() IOResult[Config] {
			return func() Result[Config] {
				return result.Left[Config](expectedErr)
			}
		}

		// Create a ReaderIOResult (won't be executed)
		computation := func(cfg Config) IOResult[string] {
			return func() Result[string] {
				return result.Of("should not be called")
			}
		}

		// Execute using ReadIOEither
		ioResult := ReadIOEither[string](getConfig())(computation)
		res := ioResult()

		assert.True(t, result.IsLeft(res))
		leftVal := result.Fold(F.Identity[error], func(string) error { return nil })(res)
		assert.Equal(t, expectedErr, leftVal)
	})

	t.Run("failure case - computation fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("computation failed")

		// Create an IOResult that successfully produces a config
		getConfig := func() IOResult[Config] {
			return func() Result[Config] {
				return result.Of(Config{BaseURL: "https://api.example.com"})
			}
		}

		// Create a ReaderIOResult that fails
		computation := func(cfg Config) IOResult[string] {
			return func() Result[string] {
				return result.Left[string](expectedErr)
			}
		}

		// Execute using ReadIOEither
		ioResult := ReadIOEither[string](getConfig())(computation)
		res := ioResult()

		assert.True(t, result.IsLeft(res))
		leftVal := result.Fold(F.Identity[error], func(string) error { return nil })(res)
		assert.Equal(t, expectedErr, leftVal)
	})
}

func TestReadIOResult(t *testing.T) {
	type Database struct {
		ConnectionString string
	}

	t.Run("success case - database and query both succeed", func(t *testing.T) {
		// Create an IOResult that successfully produces a database
		getDB := func() IOResult[Database] {
			return func() Result[Database] {
				return result.Of(Database{ConnectionString: "localhost:5432"})
			}
		}

		// Create a ReaderIOResult that uses the database
		queryUsers := func(db Database) IOResult[int] {
			return func() Result[int] {
				// Simulate query returning user count
				return result.Of(42)
			}
		}

		// Execute using ReadIOResult
		ioResult := ReadIOResult[int](getDB())(queryUsers)
		res := ioResult()

		assert.True(t, result.IsRight(res))
		assert.Equal(t, 42, result.GetOrElse(func(error) int { return 0 })(res))
	})

	t.Run("failure case - database connection fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("connection failed")

		// Create an IOResult that fails to produce a database
		getDB := func() IOResult[Database] {
			return func() Result[Database] {
				return result.Left[Database](expectedErr)
			}
		}

		// Create a ReaderIOResult (won't be executed)
		queryUsers := func(db Database) IOResult[int] {
			return func() Result[int] {
				return result.Of(0)
			}
		}

		// Execute using ReadIOResult
		ioResult := ReadIOResult[int](getDB())(queryUsers)
		res := ioResult()

		assert.True(t, result.IsLeft(res))
		leftVal := result.Fold(F.Identity[error], func(int) error { return nil })(res)
		assert.Equal(t, expectedErr, leftVal)
	})

	t.Run("failure case - query fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("query failed")

		// Create an IOResult that successfully produces a database
		getDB := func() IOResult[Database] {
			return func() Result[Database] {
				return result.Of(Database{ConnectionString: "localhost:5432"})
			}
		}

		// Create a ReaderIOResult that fails
		queryUsers := func(db Database) IOResult[int] {
			return func() Result[int] {
				return result.Left[int](expectedErr)
			}
		}

		// Execute using ReadIOResult
		ioResult := ReadIOResult[int](getDB())(queryUsers)
		res := ioResult()

		assert.True(t, result.IsLeft(res))
		leftVal := result.Fold(F.Identity[error], func(int) error { return nil })(res)
		assert.Equal(t, expectedErr, leftVal)
	})
}

func TestReadIO(t *testing.T) {
	type Logger struct {
		Level string
	}

	t.Run("success case - logger and operation both succeed", func(t *testing.T) {
		// Create an IO that produces a logger (always succeeds)
		getLogger := func() IO[Logger] {
			return func() Logger {
				return Logger{Level: "INFO"}
			}
		}

		// Create a ReaderIOResult that uses the logger
		logMessage := func(logger Logger) IOResult[string] {
			return func() Result[string] {
				return result.Of(fmt.Sprintf("[%s] Message logged", logger.Level))
			}
		}

		// Execute using ReadIO
		ioResult := ReadIO[string](getLogger())(logMessage)
		res := ioResult()

		assert.True(t, result.IsRight(res))
		assert.Equal(t, "[INFO] Message logged", result.GetOrElse(func(error) string { return "" })(res))
	})

	t.Run("failure case - operation fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("logging failed")

		// Create an IO that produces a logger (always succeeds)
		getLogger := func() IO[Logger] {
			return func() Logger {
				return Logger{Level: "ERROR"}
			}
		}

		// Create a ReaderIOResult that fails
		logMessage := func(logger Logger) IOResult[string] {
			return func() Result[string] {
				return result.Left[string](expectedErr)
			}
		}

		// Execute using ReadIO
		ioResult := ReadIO[string](getLogger())(logMessage)
		res := ioResult()

		assert.True(t, result.IsLeft(res))
		leftVal := result.Fold(F.Identity[error], func(string) error { return nil })(res)
		assert.Equal(t, expectedErr, leftVal)
	})

	t.Run("success case - complex computation with context", func(t *testing.T) {
		type AppContext struct {
			UserID   int
			Username string
		}

		// Create an IO that produces an app context
		getContext := func() IO[AppContext] {
			return func() AppContext {
				return AppContext{UserID: 123, Username: "alice"}
			}
		}

		// Create a ReaderIOResult that uses the context
		processUser := func(ctx AppContext) IOResult[string] {
			return func() Result[string] {
				return result.Of(fmt.Sprintf("Processing user %s (ID: %d)", ctx.Username, ctx.UserID))
			}
		}

		// Execute using ReadIO
		ioResult := ReadIO[string](getContext())(processUser)
		res := ioResult()

		assert.True(t, result.IsRight(res))
		assert.Equal(t, "Processing user alice (ID: 123)", result.GetOrElse(func(error) string { return "" })(res))
	})
}
