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

package readerio

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/IBM/fp-go/v2/io"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/reader"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Test environment types
type DetailedConfig struct {
	Host  string
	Port  int
	Debug bool
}

type SimpleConfig struct {
	Host string
	Port int
}

type AppConfig struct {
	Database DatabaseConfig
	Server   ServerConfig
	LogLevel string
}

type DatabaseConfig struct {
	Host string
	Port int
}

type ServerConfig struct {
	Port    int
	Timeout int
}

type UserEnv struct {
	UserID int
}

type FullEnv struct {
	UserID int
	Role   string
}

// TestPromapBasic tests basic Promap functionality
func TestPromapBasic(t *testing.T) {
	t.Run("transform both input and output", func(t *testing.T) {
		// ReaderIO that reads port from SimpleConfig
		getPort := func(c SimpleConfig) IO[int] {
			return io.Of(c.Port)
		}

		// Adapt DetailedConfig to SimpleConfig and convert int to string
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Host: d.Host, Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getPort)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080, Debug: true})()

		assert.Equal(t, "8080", result)
	})

	t.Run("identity transformations", func(t *testing.T) {
		getValue := func(n int) IO[int] {
			return io.Of(n * 2)
		}

		// Identity functions should not change behavior
		identity := reader.Ask[int]()
		adapted := Promap(identity, identity)(getValue)
		result := adapted(5)()

		assert.Equal(t, 10, result)
	})

	t.Run("compose multiple transformations", func(t *testing.T) {
		getPort := func(c SimpleConfig) IO[int] {
			return io.Of(c.Port)
		}

		// First transformation
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Host: d.Host, Port: d.Port}
		}
		double := N.Mul(2)

		step1 := Promap(simplify, double)(getPort)

		// Second transformation
		addDebug := func(d DetailedConfig) DetailedConfig {
			d.Debug = true
			return d
		}
		toString := S.Format[int]("Port: %d")

		step2 := Promap(addDebug, toString)(step1)
		result := step2(DetailedConfig{Host: "localhost", Port: 8080})()

		assert.Equal(t, "Port: 16160", result)
	})
}

// TestPromapWithIO tests Promap with actual IO effects
func TestPromapWithIO(t *testing.T) {
	t.Run("transform IO result", func(t *testing.T) {
		counter := 0
		getAndIncrement := func(n int) IO[int] {
			return func() int {
				counter++
				return n + counter
			}
		}

		double := reader.Ask[int]()
		toString := S.Format[int]("Result: %d")

		adapted := Promap(double, toString)(getAndIncrement)
		result := adapted(10)()

		assert.Equal(t, "Result: 11", result)
		assert.Equal(t, 1, counter)
	})

	t.Run("environment transformation with side effects", func(t *testing.T) {
		var log []string

		logAndReturn := func(msg string) IO[string] {
			return func() string {
				log = append(log, msg)
				return msg
			}
		}

		addPrefix := S.Prepend("Input: ")
		addSuffix := S.Append(" [processed]")

		adapted := Promap(addPrefix, addSuffix)(logAndReturn)
		result := adapted("test")()

		assert.Equal(t, "Input: test [processed]", result)
		assert.Equal(t, []string{"Input: test"}, log)
	})
}

// TestPromapEnvironmentExtraction tests extracting subsets of environments
func TestPromapEnvironmentExtraction(t *testing.T) {
	t.Run("extract database config", func(t *testing.T) {
		connectDB := func(cfg DatabaseConfig) IO[string] {
			return io.Of(fmt.Sprintf("Connected to %s:%d", cfg.Host, cfg.Port))
		}

		extractDB := func(app AppConfig) DatabaseConfig {
			return app.Database
		}
		identity := reader.Ask[string]()

		adapted := Promap(extractDB, identity)(connectDB)
		result := adapted(AppConfig{
			Database: DatabaseConfig{Host: "localhost", Port: 5432},
			Server:   ServerConfig{Port: 8080, Timeout: 30},
		})()

		assert.Equal(t, "Connected to localhost:5432", result)
	})

	t.Run("extract and transform", func(t *testing.T) {
		getServerPort := func(cfg ServerConfig) IO[int] {
			return io.Of(cfg.Port)
		}

		extractServer := func(app AppConfig) ServerConfig {
			return app.Server
		}
		formatPort := func(port int) string {
			return fmt.Sprintf("Server listening on port %d", port)
		}

		adapted := Promap(extractServer, formatPort)(getServerPort)
		result := adapted(AppConfig{
			Server: ServerConfig{Port: 8080, Timeout: 30},
		})()

		assert.Equal(t, "Server listening on port 8080", result)
	})
}

// TestLocalBasic tests basic Local functionality
func TestLocalBasic(t *testing.T) {
	t.Run("extract subset of environment", func(t *testing.T) {
		connectDB := func(cfg DatabaseConfig) IO[string] {
			return io.Of(fmt.Sprintf("Connected to %s:%d", cfg.Host, cfg.Port))
		}

		extractDB := func(app AppConfig) DatabaseConfig {
			return app.Database
		}

		adapted := Local[string](extractDB)(connectDB)
		result := adapted(AppConfig{
			Database: DatabaseConfig{Host: "localhost", Port: 5432},
		})()

		assert.Equal(t, "Connected to localhost:5432", result)
	})

	t.Run("transform environment type", func(t *testing.T) {
		getUserData := func(env UserEnv) IO[string] {
			return io.Of(fmt.Sprintf("User: %d", env.UserID))
		}

		toUserEnv := func(full FullEnv) UserEnv {
			return UserEnv{UserID: full.UserID}
		}

		adapted := Local[string](toUserEnv)(getUserData)
		result := adapted(FullEnv{UserID: 42, Role: "admin"})()

		assert.Equal(t, "User: 42", result)
	})

	t.Run("identity transformation", func(t *testing.T) {
		getValue := func(n int) IO[int] {
			return io.Of(n * 2)
		}

		identity := reader.Ask[int]()
		adapted := Local[int](identity)(getValue)
		result := adapted(5)()

		assert.Equal(t, 10, result)
	})
}

// TestLocalComposition tests composing Local transformations
func TestLocalComposition(t *testing.T) {
	t.Run("compose two Local transformations", func(t *testing.T) {
		getPort := func(cfg DatabaseConfig) IO[int] {
			return io.Of(cfg.Port)
		}

		extractDB := func(app AppConfig) DatabaseConfig {
			return app.Database
		}

		// First Local
		step1 := Local[int](extractDB)(getPort)

		// Second Local - add default values
		addDefaults := func(app AppConfig) AppConfig {
			if app.Database.Host == "" {
				app.Database.Host = "localhost"
			}
			return app
		}

		step2 := Local[int](addDefaults)(step1)
		result := step2(AppConfig{
			Database: DatabaseConfig{Host: "", Port: 5432},
		})()

		assert.Equal(t, 5432, result)
	})

	t.Run("chain multiple environment transformations", func(t *testing.T) {
		getHost := func(cfg SimpleConfig) IO[string] {
			return io.Of(cfg.Host)
		}

		// Transform DetailedConfig -> SimpleConfig
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Host: d.Host, Port: d.Port}
		}

		adapted := Local[string](simplify)(getHost)
		result := adapted(DetailedConfig{Host: "example.com", Port: 8080, Debug: true})()

		assert.Equal(t, "example.com", result)
	})
}

// TestLocalWithIO tests Local with IO effects
func TestLocalWithIO(t *testing.T) {
	t.Run("environment transformation with side effects", func(t *testing.T) {
		var accessLog []int

		logAccess := func(id int) IO[string] {
			return func() string {
				accessLog = append(accessLog, id)
				return fmt.Sprintf("Accessed: %d", id)
			}
		}

		extractUserID := func(env FullEnv) int {
			return env.UserID
		}

		adapted := Local[string](extractUserID)(logAccess)
		result := adapted(FullEnv{UserID: 123, Role: "user"})()

		assert.Equal(t, "Accessed: 123", result)
		assert.Equal(t, []int{123}, accessLog)
	})

	t.Run("multiple executions with different environments", func(t *testing.T) {
		counter := 0
		increment := func(n int) IO[int] {
			return func() int {
				counter++
				return n + counter
			}
		}

		double := N.Mul(2)
		adapted := Local[int](double)(increment)

		result1 := adapted(5)()  // 10 + 1 = 11
		result2 := adapted(10)() // 20 + 2 = 22

		assert.Equal(t, 11, result1)
		assert.Equal(t, 22, result2)
		assert.Equal(t, 2, counter)
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("environment adaptation", func(t *testing.T) {
		readConfig := func(env SimpleConfig) IO[string] {
			return io.Of(fmt.Sprintf("%s:%d", env.Host, env.Port))
		}

		simplify := func(detailed DetailedConfig) SimpleConfig {
			return SimpleConfig{Host: detailed.Host, Port: detailed.Port}
		}

		adapted := Contramap[string](simplify)(readConfig)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080, Debug: true})()

		assert.Equal(t, "localhost:8080", result)
	})

	t.Run("extract field from larger structure", func(t *testing.T) {
		getPort := func(port int) IO[string] {
			return io.Of(fmt.Sprintf("Port: %d", port))
		}

		extractPort := func(cfg SimpleConfig) int {
			return cfg.Port
		}

		adapted := Contramap[string](extractPort)(getPort)
		result := adapted(SimpleConfig{Host: "localhost", Port: 9000})()

		assert.Equal(t, "Port: 9000", result)
	})
}

// TestContramapVsLocal verifies Contramap and Local are equivalent
func TestContramapVsLocal(t *testing.T) {
	t.Run("same behavior as Local", func(t *testing.T) {
		getValue := func(n int) IO[int] {
			return io.Of(n * 3)
		}

		double := N.Mul(2)

		localResult := Local[int](double)(getValue)(5)()
		contramapResult := Contramap[int](double)(getValue)(5)()

		assert.Equal(t, localResult, contramapResult)
		assert.Equal(t, 30, localResult) // (5 * 2) * 3 = 30
	})

	t.Run("environment extraction equivalence", func(t *testing.T) {
		getHost := func(cfg SimpleConfig) IO[string] {
			return io.Of(cfg.Host)
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Host: d.Host, Port: d.Port}
		}

		env := DetailedConfig{Host: "example.com", Port: 8080, Debug: false}

		localResult := Local[string](simplify)(getHost)(env)()
		contramapResult := Contramap[string](simplify)(getHost)(env)()

		assert.Equal(t, localResult, contramapResult)
		assert.Equal(t, "example.com", localResult)
	})
}

// TestProfunctorLaws tests profunctor laws
func TestProfunctorLaws(t *testing.T) {
	t.Run("identity law", func(t *testing.T) {
		getValue := func(n int) IO[int] {
			return io.Of(n + 10)
		}

		identity := reader.Ask[int]()

		// Promap(id, id) should be equivalent to id
		adapted := Promap(identity, identity)(getValue)
		original := getValue(5)()
		transformed := adapted(5)()

		assert.Equal(t, original, transformed)
		assert.Equal(t, 15, transformed)
	})

	t.Run("composition law", func(t *testing.T) {
		getValue := func(n int) IO[int] {
			return io.Of(n * 2)
		}

		f1 := N.Add(1)
		f2 := N.Mul(3)
		g1 := N.Sub(5)
		g2 := N.Mul(2)

		// Promap(f1, g2) . Promap(f2, g1) should equal Promap(f2 . f1, g2 . g1)
		// Note: composition order is reversed for contravariant part
		step1 := Promap(f2, g1)(getValue)
		composed1 := Promap(f1, g2)(step1)

		composed2 := Promap(
			func(x int) int { return f2(f1(x)) },
			func(x int) int { return g2(g1(x)) },
		)(getValue)

		result1 := composed1(10)()
		result2 := composed2(10)()

		assert.Equal(t, result1, result2)
	})
}

// TestEdgeCases tests edge cases and special scenarios
func TestEdgeCases(t *testing.T) {
	t.Run("empty struct environment", func(t *testing.T) {
		type Empty struct{}

		getValue := func(e Empty) IO[int] {
			return io.Of(42)
		}

		identity := reader.Ask[Empty]()
		adapted := Local[int](identity)(getValue)
		result := adapted(Empty{})()

		assert.Equal(t, 42, result)
	})

	t.Run("function type handling", func(t *testing.T) {
		getFunc := func(n int) IO[func(int) int] {
			return io.Of(N.Mul(2))
		}

		double := N.Mul(2)
		applyFunc := reader.Read[int](5)

		adapted := Promap(double, applyFunc)(getFunc)
		result := adapted(3)() // (3 * 2) = 6, then func(5) = 10

		assert.Equal(t, 10, result)
	})

	t.Run("complex nested transformations", func(t *testing.T) {
		type Level3 struct{ Value int }
		type Level2 struct{ L3 Level3 }
		type Level1 struct{ L2 Level2 }

		getValue := func(l3 Level3) IO[int] {
			return io.Of(l3.Value)
		}

		extract := func(l1 Level1) Level3 {
			return l1.L2.L3
		}
		multiply := N.Mul(10)

		adapted := Promap(extract, multiply)(getValue)
		result := adapted(Level1{L2: Level2{L3: Level3{Value: 7}}})()

		assert.Equal(t, 70, result)
	})
}

// TestRealWorldScenarios tests practical use cases
func TestRealWorldScenarios(t *testing.T) {
	t.Run("database connection with config extraction", func(t *testing.T) {
		type DBConfig struct {
			ConnectionString string
		}

		type AppSettings struct {
			DB      DBConfig
			APIKey  string
			Timeout int
		}

		connect := func(cfg DBConfig) IO[string] {
			return io.Of("Connected: " + cfg.ConnectionString)
		}

		extractDB := func(settings AppSettings) DBConfig {
			return settings.DB
		}

		adapted := Local[string](extractDB)(connect)
		result := adapted(AppSettings{
			DB:      DBConfig{ConnectionString: "postgres://localhost"},
			APIKey:  "secret",
			Timeout: 30,
		})()

		assert.Equal(t, "Connected: postgres://localhost", result)
	})

	t.Run("logging with environment transformation", func(t *testing.T) {
		type LogContext struct {
			RequestID string
			UserID    int
		}

		type RequestContext struct {
			RequestID string
			UserID    int
			Path      string
			Method    string
		}

		var logs []string
		logMessage := func(ctx LogContext) IO[func()] {
			return func() func() {
				return func() {
					logs = append(logs, fmt.Sprintf("[%s] User %d", ctx.RequestID, ctx.UserID))
				}
			}
		}

		extractLogContext := func(req RequestContext) LogContext {
			return LogContext{RequestID: req.RequestID, UserID: req.UserID}
		}

		adapted := Local[func()](extractLogContext)(logMessage)
		result := adapted(RequestContext{
			RequestID: "req-123",
			UserID:    42,
			Path:      "/api/users",
			Method:    "GET",
		})()

		result()
		assert.Equal(t, []string{"[req-123] User 42"}, logs)
	})

	t.Run("API response transformation", func(t *testing.T) {
		type APIResponse struct {
			StatusCode int
			Body       string
		}

		type EnrichedResponse struct {
			Response  APIResponse
			Timestamp int64
			RequestID string
		}

		formatResponse := func(resp APIResponse) IO[string] {
			return io.Of(fmt.Sprintf("Status: %d, Body: %s", resp.StatusCode, resp.Body))
		}

		extractResponse := func(enriched EnrichedResponse) APIResponse {
			return enriched.Response
		}
		addMetadata := func(s string) string {
			return "[API] " + s
		}

		adapted := Promap(extractResponse, addMetadata)(formatResponse)
		result := adapted(EnrichedResponse{
			Response:  APIResponse{StatusCode: 200, Body: "OK"},
			Timestamp: 1234567890,
			RequestID: "req-456",
		})()

		assert.Equal(t, "[API] Status: 200, Body: OK", result)
	})
}

// TestLocalIOK tests LocalIOK functionality
func TestLocalIOK(t *testing.T) {
	t.Run("basic IO transformation", func(t *testing.T) {
		// IO effect that loads config from a path
		loadConfig := func(path string) io.IO[SimpleConfig] {
			return func() SimpleConfig {
				// Simulate loading config
				return SimpleConfig{Host: "localhost", Port: 8080}
			}
		}

		// ReaderIO that uses the config
		useConfig := func(cfg SimpleConfig) io.IO[string] {
			return io.Of(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
		}

		// Compose using LocalIOK
		adapted := LocalIOK[string, SimpleConfig, string](loadConfig)(useConfig)
		result := adapted("config.json")()

		assert.Equal(t, "localhost:8080", result)
	})

	t.Run("IO transformation with side effects", func(t *testing.T) {
		var loadLog []string

		loadData := func(key string) io.IO[int] {
			return func() int {
				loadLog = append(loadLog, "Loading: "+key)
				return len(key) * 10
			}
		}

		processData := func(n int) io.IO[string] {
			return io.Of(fmt.Sprintf("Processed: %d", n))
		}

		adapted := LocalIOK[string, int, string](loadData)(processData)
		result := adapted("test")()

		assert.Equal(t, "Processed: 40", result)
		assert.Equal(t, []string{"Loading: test"}, loadLog)
	})

	t.Run("compose multiple LocalIOK", func(t *testing.T) {
		// First transformation: string -> int
		parseID := func(s string) io.IO[int] {
			return func() int {
				id, _ := strconv.Atoi(s)
				return id
			}
		}

		// Second transformation: int -> UserEnv
		loadUser := func(id int) io.IO[UserEnv] {
			return func() UserEnv {
				return UserEnv{UserID: id}
			}
		}

		// Use the UserEnv
		formatUser := func(env UserEnv) io.IO[string] {
			return io.Of(fmt.Sprintf("User ID: %d", env.UserID))
		}

		// Compose transformations
		step1 := LocalIOK[string, UserEnv, int](loadUser)(formatUser)
		step2 := LocalIOK[string, int, string](parseID)(step1)

		result := step2("42")()
		assert.Equal(t, "User ID: 42", result)
	})

	t.Run("environment extraction with IO", func(t *testing.T) {
		// Extract database config from app config
		extractDB := func(app AppConfig) io.IO[DatabaseConfig] {
			return func() DatabaseConfig {
				// Could perform validation or default setting here
				cfg := app.Database
				if cfg.Host == "" {
					cfg.Host = "localhost"
				}
				return cfg
			}
		}

		// Use the database config
		connectDB := func(cfg DatabaseConfig) io.IO[string] {
			return io.Of(fmt.Sprintf("Connected to %s:%d", cfg.Host, cfg.Port))
		}

		adapted := LocalIOK[string, DatabaseConfig, AppConfig](extractDB)(connectDB)
		result := adapted(AppConfig{
			Database: DatabaseConfig{Host: "", Port: 5432},
		})()

		assert.Equal(t, "Connected to localhost:5432", result)
	})

	t.Run("real-world: load and parse config file", func(t *testing.T) {
		type ConfigFile struct {
			Path string
		}

		// Simulate reading file content
		readFile := func(cf ConfigFile) io.IO[string] {
			return func() string {
				return `{"host":"example.com","port":9000}`
			}
		}

		// Parse the content
		parseConfig := func(content string) io.IO[SimpleConfig] {
			return io.Of(SimpleConfig{Host: "example.com", Port: 9000})
		}

		// Use the parsed config
		useConfig := func(cfg SimpleConfig) io.IO[string] {
			return io.Of(fmt.Sprintf("Using %s:%d", cfg.Host, cfg.Port))
		}

		// Compose the pipeline
		step1 := LocalIOK[string, SimpleConfig, string](parseConfig)(useConfig)
		step2 := LocalIOK[string, string, ConfigFile](readFile)(step1)

		result := step2(ConfigFile{Path: "app.json"})()
		assert.Equal(t, "Using example.com:9000", result)
	})
}
