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
	"context"
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/stretchr/testify/assert"
)

type OuterContext struct {
	Value  string
	Number int
}

type InnerContext struct {
	Value string
}

func TestLocal(t *testing.T) {
	t.Run("transforms context for inner effect", func(t *testing.T) {
		// Create an effect that uses InnerContext
		innerEffect := Of[InnerContext]("result")

		// Transform OuterContext to InnerContext
		accessor := func(outer OuterContext) InnerContext {
			return InnerContext{Value: outer.Value}
		}

		// Apply Local to transform the context
		kleisli := Local[OuterContext, InnerContext, string](accessor)
		outerEffect := kleisli(innerEffect)

		// Run with OuterContext
		ioResult := Provide[OuterContext, string](OuterContext{
			Value:  "test",
			Number: 42,
		})(outerEffect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "result", result)
	})

	t.Run("allows accessing outer context fields", func(t *testing.T) {
		// Create an effect that reads from InnerContext
		innerEffect := Chain(func(_ string) Effect[InnerContext, string] {
			return Of[InnerContext]("inner value")
		})(Of[InnerContext]("start"))

		// Transform context
		accessor := func(outer OuterContext) InnerContext {
			return InnerContext{Value: outer.Value + " transformed"}
		}

		kleisli := Local[OuterContext, InnerContext, string](accessor)
		outerEffect := kleisli(innerEffect)

		// Run with OuterContext
		ioResult := Provide[OuterContext, string](OuterContext{
			Value:  "original",
			Number: 100,
		})(outerEffect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "inner value", result)
	})

	t.Run("propagates errors from inner effect", func(t *testing.T) {
		expectedErr := assert.AnError
		innerEffect := Fail[InnerContext, string](expectedErr)

		accessor := func(outer OuterContext) InnerContext {
			return InnerContext{Value: outer.Value}
		}

		kleisli := Local[OuterContext, InnerContext, string](accessor)
		outerEffect := kleisli(innerEffect)

		ioResult := Provide[OuterContext, string](OuterContext{
			Value:  "test",
			Number: 42,
		})(outerEffect)
		readerResult := RunSync(ioResult)
		_, err := readerResult(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("chains multiple Local transformations", func(t *testing.T) {
		type Level1 struct {
			A string
		}
		type Level2 struct {
			B string
		}
		type Level3 struct {
			C string
		}

		// Effect at deepest level
		level3Effect := Of[Level3]("deep result")

		// Transform Level2 -> Level3
		local23 := Local[Level2, Level3, string](func(l2 Level2) Level3 {
			return Level3{C: l2.B + "-c"}
		})

		// Transform Level1 -> Level2
		local12 := Local[Level1, Level2, string](func(l1 Level1) Level2 {
			return Level2{B: l1.A + "-b"}
		})

		// Compose transformations
		level2Effect := local23(level3Effect)
		level1Effect := local12(level2Effect)

		// Run with Level1 context
		ioResult := Provide[Level1, string](Level1{A: "a"})(level1Effect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "deep result", result)
	})

	t.Run("works with complex context transformations", func(t *testing.T) {
		type DatabaseConfig struct {
			Host     string
			Port     int
			Database string
		}

		type AppConfig struct {
			DB      DatabaseConfig
			APIKey  string
			Timeout int
		}

		// Effect that needs only DatabaseConfig
		dbEffect := Of[DatabaseConfig]("connected")

		// Extract DB config from AppConfig
		accessor := func(app AppConfig) DatabaseConfig {
			return app.DB
		}

		kleisli := Local[AppConfig, DatabaseConfig, string](accessor)
		appEffect := kleisli(dbEffect)

		// Run with full AppConfig
		ioResult := Provide[AppConfig, string](AppConfig{
			DB: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				Database: "mydb",
			},
			APIKey:  "secret",
			Timeout: 30,
		})(appEffect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "connected", result)
	})
}

func TestContramap(t *testing.T) {
	t.Run("is equivalent to Local", func(t *testing.T) {
		innerEffect := Of[InnerContext](42)

		accessor := func(outer OuterContext) InnerContext {
			return InnerContext{Value: outer.Value}
		}

		// Test Local
		localKleisli := Local[OuterContext, InnerContext, int](accessor)
		localEffect := localKleisli(innerEffect)

		// Test Contramap
		contramapKleisli := Contramap[OuterContext, InnerContext, int](accessor)
		contramapEffect := contramapKleisli(innerEffect)

		outerCtx := OuterContext{Value: "test", Number: 100}

		// Run both
		localIO := Provide[OuterContext, int](outerCtx)(localEffect)
		localReader := RunSync(localIO)
		localResult, localErr := localReader(context.Background())

		contramapIO := Provide[OuterContext, int](outerCtx)(contramapEffect)
		contramapReader := RunSync(contramapIO)
		contramapResult, contramapErr := contramapReader(context.Background())

		assert.NoError(t, localErr)
		assert.NoError(t, contramapErr)
		assert.Equal(t, localResult, contramapResult)
	})

	t.Run("transforms context correctly", func(t *testing.T) {
		innerEffect := Of[InnerContext]("success")

		accessor := func(outer OuterContext) InnerContext {
			return InnerContext{Value: outer.Value + " modified"}
		}

		kleisli := Contramap[OuterContext, InnerContext, string](accessor)
		outerEffect := kleisli(innerEffect)

		ioResult := Provide[OuterContext, string](OuterContext{
			Value:  "original",
			Number: 50,
		})(outerEffect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("handles errors from inner effect", func(t *testing.T) {
		expectedErr := assert.AnError
		innerEffect := Fail[InnerContext, int](expectedErr)

		accessor := func(outer OuterContext) InnerContext {
			return InnerContext{Value: outer.Value}
		}

		kleisli := Contramap[OuterContext, InnerContext, int](accessor)
		outerEffect := kleisli(innerEffect)

		ioResult := Provide[OuterContext, int](OuterContext{
			Value:  "test",
			Number: 42,
		})(outerEffect)
		readerResult := RunSync(ioResult)
		_, err := readerResult(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestLocalAndContramapInteroperability(t *testing.T) {
	t.Run("can be used interchangeably", func(t *testing.T) {
		type Config1 struct {
			Value string
		}
		type Config2 struct {
			Data string
		}
		type Config3 struct {
			Info string
		}

		// Effect at deepest level
		effect3 := Of[Config3]("result")

		// Use Local for first transformation
		local23 := Local[Config2, Config3, string](func(c2 Config2) Config3 {
			return Config3{Info: c2.Data}
		})

		// Use Contramap for second transformation
		contramap12 := Contramap[Config1, Config2, string](func(c1 Config1) Config2 {
			return Config2{Data: c1.Value}
		})

		// Compose them
		effect2 := local23(effect3)
		effect1 := contramap12(effect2)

		// Run
		ioResult := Provide[Config1, string](Config1{Value: "test"})(effect1)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "result", result)
	})
}

func TestLocalEffectK(t *testing.T) {
	t.Run("transforms context using effectful function", func(t *testing.T) {
		type DatabaseConfig struct {
			ConnectionString string
		}

		type AppConfig struct {
			ConfigPath string
		}

		// Effect that needs DatabaseConfig
		dbEffect := Of[DatabaseConfig]("query result")

		// Transform AppConfig to DatabaseConfig effectfully
		loadConfig := func(app AppConfig) Effect[AppConfig, DatabaseConfig] {
			return Of[AppConfig](DatabaseConfig{
				ConnectionString: "loaded from " + app.ConfigPath,
			})
		}

		// Apply the transformation
		transform := LocalEffectK[string](loadConfig)
		appEffect := transform(dbEffect)

		// Run with AppConfig
		ioResult := Provide[AppConfig, string](AppConfig{
			ConfigPath: "/etc/app.conf",
		})(appEffect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "query result", result)
	})

	t.Run("propagates errors from context transformation", func(t *testing.T) {
		type InnerCtx struct {
			Value string
		}

		type OuterCtx struct {
			Path string
		}

		innerEffect := Of[InnerCtx]("success")

		expectedErr := assert.AnError
		// Context transformation that fails
		failingTransform := func(outer OuterCtx) Effect[OuterCtx, InnerCtx] {
			return Fail[OuterCtx, InnerCtx](expectedErr)
		}

		transform := LocalEffectK[string](failingTransform)
		outerEffect := transform(innerEffect)

		ioResult := Provide[OuterCtx, string](OuterCtx{Path: "test"})(outerEffect)
		readerResult := RunSync(ioResult)
		_, err := readerResult(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("propagates errors from inner effect", func(t *testing.T) {
		type InnerCtx struct {
			Value string
		}

		type OuterCtx struct {
			Path string
		}

		expectedErr := assert.AnError
		innerEffect := Fail[InnerCtx, string](expectedErr)

		// Successful context transformation
		transform := func(outer OuterCtx) Effect[OuterCtx, InnerCtx] {
			return Of[OuterCtx](InnerCtx{Value: outer.Path})
		}

		transformK := LocalEffectK[string](transform)
		outerEffect := transformK(innerEffect)

		ioResult := Provide[OuterCtx, string](OuterCtx{Path: "test"})(outerEffect)
		readerResult := RunSync(ioResult)
		_, err := readerResult(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("allows effectful context transformation with IO operations", func(t *testing.T) {
		type Config struct {
			Data string
		}

		type AppContext struct {
			ConfigFile string
		}

		// Effect that uses Config
		configEffect := Chain(func(cfg Config) Effect[Config, string] {
			return Of[Config]("processed: " + cfg.Data)
		})(readerreaderioresult.Ask[Config]())

		// Effectful transformation that simulates loading config
		loadConfigEffect := func(app AppContext) Effect[AppContext, Config] {
			// Simulate IO operation (e.g., reading file)
			return Of[AppContext](Config{
				Data: "loaded from " + app.ConfigFile,
			})
		}

		transform := LocalEffectK[string](loadConfigEffect)
		appEffect := transform(configEffect)

		ioResult := Provide[AppContext, string](AppContext{
			ConfigFile: "config.json",
		})(appEffect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "processed: loaded from config.json", result)
	})

	t.Run("chains multiple LocalEffectK transformations", func(t *testing.T) {
		type Level1 struct {
			A string
		}
		type Level2 struct {
			B string
		}
		type Level3 struct {
			C string
		}

		// Effect at deepest level
		level3Effect := Of[Level3]("deep result")

		// Transform Level2 -> Level3 effectfully
		transform23 := LocalEffectK[string](func(l2 Level2) Effect[Level2, Level3] {
			return Of[Level2](Level3{C: l2.B + "-c"})
		})

		// Transform Level1 -> Level2 effectfully
		transform12 := LocalEffectK[string](func(l1 Level1) Effect[Level1, Level2] {
			return Of[Level1](Level2{B: l1.A + "-b"})
		})

		// Compose transformations
		level2Effect := transform23(level3Effect)
		level1Effect := transform12(level2Effect)

		// Run with Level1 context
		ioResult := Provide[Level1, string](Level1{A: "a"})(level1Effect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "deep result", result)
	})

	t.Run("accesses outer context during transformation", func(t *testing.T) {
		type DatabaseConfig struct {
			Host string
			Port int
		}

		type AppConfig struct {
			Environment string
			DBHost      string
			DBPort      int
		}

		// Effect that needs DatabaseConfig
		dbEffect := Chain(func(cfg DatabaseConfig) Effect[DatabaseConfig, string] {
			return Of[DatabaseConfig](fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
		})(readerreaderioresult.Ask[DatabaseConfig]())

		// Transform using outer context
		transformWithContext := func(app AppConfig) Effect[AppConfig, DatabaseConfig] {
			// Access outer context to build inner context
			prefix := ""
			if app.Environment == "prod" {
				prefix = "prod-"
			}
			return Of[AppConfig](DatabaseConfig{
				Host: prefix + app.DBHost,
				Port: app.DBPort,
			})
		}

		transform := LocalEffectK[string](transformWithContext)
		appEffect := transform(dbEffect)

		ioResult := Provide[AppConfig, string](AppConfig{
			Environment: "prod",
			DBHost:      "localhost",
			DBPort:      5432,
		})(appEffect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Contains(t, result, "prod-localhost")
	})

	t.Run("validates context during transformation", func(t *testing.T) {
		type ValidatedConfig struct {
			APIKey string
		}

		type RawConfig struct {
			APIKey string
		}

		innerEffect := Of[ValidatedConfig]("success")

		// Validation that can fail
		validateConfig := func(raw RawConfig) Effect[RawConfig, ValidatedConfig] {
			if raw.APIKey == "" {
				return Fail[RawConfig, ValidatedConfig](assert.AnError)
			}
			return Of[RawConfig](ValidatedConfig{
				APIKey: raw.APIKey,
			})
		}

		transform := LocalEffectK[string](validateConfig)
		outerEffect := transform(innerEffect)

		// Test with invalid config
		ioResult := Provide[RawConfig, string](RawConfig{APIKey: ""})(outerEffect)
		readerResult := RunSync(ioResult)
		_, err := readerResult(context.Background())

		assert.Error(t, err)

		// Test with valid config
		ioResult2 := Provide[RawConfig, string](RawConfig{APIKey: "valid-key"})(outerEffect)
		readerResult2 := RunSync(ioResult2)
		result, err2 := readerResult2(context.Background())

		assert.NoError(t, err2)
		assert.Equal(t, "success", result)
	})

	t.Run("composes with other Local functions", func(t *testing.T) {
		type Level1 struct {
			Value string
		}
		type Level2 struct {
			Data string
		}
		type Level3 struct {
			Info string
		}

		// Effect at deepest level
		effect3 := Of[Level3]("result")

		// Use LocalEffectK for first transformation (effectful)
		localEffectK23 := LocalEffectK[string](func(l2 Level2) Effect[Level2, Level3] {
			return Of[Level2](Level3{Info: l2.Data})
		})

		// Use Local for second transformation (pure)
		local12 := Local[Level1, Level2, string](func(l1 Level1) Level2 {
			return Level2{Data: l1.Value}
		})

		// Compose them
		effect2 := localEffectK23(effect3)
		effect1 := local12(effect2)

		// Run
		ioResult := Provide[Level1, string](Level1{Value: "test"})(effect1)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "result", result)
	})

	t.Run("handles complex nested effects in transformation", func(t *testing.T) {
		type InnerCtx struct {
			Value int
		}

		type OuterCtx struct {
			Multiplier int
		}

		// Effect that uses InnerCtx
		innerEffect := Chain(func(ctx InnerCtx) Effect[InnerCtx, int] {
			return Of[InnerCtx](ctx.Value * 2)
		})(readerreaderioresult.Ask[InnerCtx]())

		// Complex transformation with nested effects
		complexTransform := func(outer OuterCtx) Effect[OuterCtx, InnerCtx] {
			return Of[OuterCtx](InnerCtx{
				Value: outer.Multiplier * 10,
			})
		}

		transform := LocalEffectK[int](complexTransform)
		outerEffect := transform(innerEffect)

		ioResult := Provide[OuterCtx, int](OuterCtx{Multiplier: 3})(outerEffect)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 60, result) // 3 * 10 * 2
	})
}
