---
sidebar_position: 12
---

# ReaderIOEither (v1)

The `ReaderIOEither` type combines `Reader`, `IO`, and `Either` for lazy computations that depend on an environment, perform side effects, and can fail.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [ReaderIOEither v2](../v2/readerioeither).

**Key differences in v2:**
- Simplified API
- Better type inference
- Improved composition
- More utility functions
:::

## Overview

`ReaderIOEither` represents a computation that:
- Depends on a shared environment (Reader)
- Performs side effects (IO)
- Can fail with an error (Either Left)
- Can succeed with a value (Either Right)
- Is lazy (deferred execution)

```go
type ReaderIOEither[R, E, A any] func(R) IOEither[E, A]
```

This is the most powerful combination in fp-go, suitable for real-world applications.

## Creating ReaderIOEither Values

### Of (Success)

```go
package main

import (
    "fmt"
    RIOE "github.com/IBM/fp-go/readerioeither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    Debug bool
}

func main() {
    // Create successful ReaderIOEither
    rioe := RIOE.Of[Config, error, int](42)
    
    config := Config{Debug: true}
    result := rioe(config)()
    
    fmt.Println(E.IsRight(result)) // true
}
```

### Left (Error)

```go
package main

import (
    "errors"
    "fmt"
    RIOE "github.com/IBM/fp-go/readerioeither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    Debug bool
}

func main() {
    // Create failed ReaderIOEither
    rioe := RIOE.Left[Config, int](errors.New("failed"))
    
    config := Config{Debug: true}
    result := rioe(config)()
    
    fmt.Println(E.IsLeft(result)) // true
}
```

### Ask

Access the environment:

```go
package main

import (
    "fmt"
    RIOE "github.com/IBM/fp-go/readerioeither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    APIKey string
}

func main() {
    // Get environment as Right value
    rioe := RIOE.Ask[Config, error]()
    
    config := Config{APIKey: "secret"}
    result := rioe(config)()
    
    if E.IsRight(result) {
        cfg := E.GetOrElse(func() Config { return Config{} })(result)
        fmt.Println("API Key:", cfg.APIKey)
    }
}
```

## Basic Operations

### Map

Transform the success value:

```go
package main

import (
    "fmt"
    RIOE "github.com/IBM/fp-go/readerioeither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    Multiplier int
}

func main() {
    rioe := RIOE.Of[Config, error, int](5)
    
    // Map transforms the Right value
    doubled := RIOE.Map(func(n int) int {
        return n * 2
    })(rioe)
    
    config := Config{Multiplier: 2}
    result := doubled(config)()
    
    value := E.GetOrElse(func() int { return 0 })(result)
    fmt.Println(value) // 10
}
```

### Chain

Chain operations that return ReaderIOEither:

```go
package main

import (
    "fmt"
    RIOE "github.com/IBM/fp-go/readerioeither"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    MaxValue int
}

func validate(n int) RIOE.ReaderIOEither[Config, error, int] {
    return func(cfg Config) IOE.IOEither[error, int] {
        return func() E.Either[error, int] {
            if n > cfg.MaxValue {
                return E.Left[int](fmt.Errorf("exceeds max: %d", cfg.MaxValue))
            }
            return E.Right[error](n)
        }
    }
}

func main() {
    rioe := RIOE.Of[Config, error, int](100)
    
    validated := RIOE.Chain(validate)(rioe)
    
    config := Config{MaxValue: 50}
    result := validated(config)()
    
    fmt.Println(E.IsLeft(result)) // true (exceeds max)
}
```

## Real-World Application

### Complete Web Service

```go
package main

import (
    "fmt"
    "net/http"
    RIOE "github.com/IBM/fp-go/readerioeither"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
)

type Logger interface {
    Info(msg string)
    Error(msg string)
}

type DB interface {
    FindUser(id int) (*User, error)
    SaveUser(user User) error
}

type Cache interface {
    Get(key string) (*User, bool)
    Set(key string, user User)
}

type Deps struct {
    Logger Logger
    DB     DB
    Cache  Cache
}

type User struct {
    ID   int
    Name string
}

type ConsoleLogger struct{}

func (l ConsoleLogger) Info(msg string) {
    fmt.Println("[INFO]", msg)
}

func (l ConsoleLogger) Error(msg string) {
    fmt.Println("[ERROR]", msg)
}

type MockDB struct {
    users map[int]User
}

func (db MockDB) FindUser(id int) (*User, error) {
    if user, ok := db.users[id]; ok {
        return &user, nil
    }
    return nil, fmt.Errorf("user not found")
}

func (db MockDB) SaveUser(user User) error {
    db.users[user.ID] = user
    return nil
}

type MemoryCache struct {
    data map[string]User
}

func (c *MemoryCache) Get(key string) (*User, bool) {
    if user, ok := c.data[key]; ok {
        return &user, true
    }
    return nil, false
}

func (c *MemoryCache) Set(key string, user User) {
    c.data[key] = user
}

// Log operation
func logInfo(msg string) RIOE.ReaderIOEither[Deps, error, struct{}] {
    return func(deps Deps) IOE.IOEither[error, struct{}] {
        return func() E.Either[error, struct{}] {
            deps.Logger.Info(msg)
            return E.Right[error](struct{}{})
        }
    }
}

// Get user from cache
func getCached(id int) RIOE.ReaderIOEither[Deps, error, *User] {
    return func(deps Deps) IOE.IOEither[error, *User] {
        return func() E.Either[error, *User] {
            key := fmt.Sprintf("user:%d", id)
            if user, ok := deps.Cache.Get(key); ok {
                deps.Logger.Info(fmt.Sprintf("Cache hit for user %d", id))
                return E.Right[error](user)
            }
            return E.Left[*User](fmt.Errorf("cache miss"))
        }
    }
}

// Get user from database
func getFromDB(id int) RIOE.ReaderIOEither[Deps, error, *User] {
    return func(deps Deps) IOE.IOEither[error, *User] {
        return func() E.Either[error, *User] {
            deps.Logger.Info(fmt.Sprintf("Fetching user %d from DB", id))
            user, err := deps.DB.FindUser(id)
            if err != nil {
                return E.Left[*User](err)
            }
            return E.Right[error](user)
        }
    }
}

// Cache user
func cacheUser(user *User) RIOE.ReaderIOEither[Deps, error, *User] {
    return func(deps Deps) IOE.IOEither[error, *User] {
        return func() E.Either[error, *User] {
            key := fmt.Sprintf("user:%d", user.ID)
            deps.Cache.Set(key, *user)
            deps.Logger.Info(fmt.Sprintf("Cached user %d", user.ID))
            return E.Right[error](user)
        }
    }
}

// Get user with caching
func getUser(id int) RIOE.ReaderIOEither[Deps, error, *User] {
    return RIOE.OrElse(func(err error) RIOE.ReaderIOEither[Deps, error, *User] {
        // Cache miss, get from DB and cache
        return F.Pipe2(
            getFromDB(id),
            RIOE.Chain(cacheUser),
        )
    })(getCached(id))
}

func main() {
    deps := Deps{
        Logger: ConsoleLogger{},
        DB: MockDB{
            users: map[int]User{
                1: {ID: 1, Name: "Alice"},
            },
        },
        Cache: &MemoryCache{data: make(map[string]User)},
    }
    
    // First call - cache miss, fetches from DB
    result1 := getUser(1)(deps)()
    
    // Second call - cache hit
    result2 := getUser(1)(deps)()
    
    if E.IsRight(result1) && E.IsRight(result2) {
        fmt.Println("Both calls succeeded")
    }
}
```

## HTTP API with Full Error Handling

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    RIOE "github.com/IBM/fp-go/readerioeither"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
)

type APIConfig struct {
    BaseURL string
    APIKey  string
    Timeout int
}

type APIResponse struct {
    Data  string `json:"data"`
    Error string `json:"error"`
}

func validateConfig() RIOE.ReaderIOEither[APIConfig, error, APIConfig] {
    return func(cfg APIConfig) IOE.IOEither[error, APIConfig] {
        return func() E.Either[error, APIConfig] {
            if cfg.BaseURL == "" {
                return E.Left[APIConfig](fmt.Errorf("base URL required"))
            }
            if cfg.APIKey == "" {
                return E.Left[APIConfig](fmt.Errorf("API key required"))
            }
            return E.Right[error](cfg)
        }
    }
}

func makeRequest(path string) RIOE.ReaderIOEither[APIConfig, error, string] {
    return func(cfg APIConfig) IOE.IOEither[error, string] {
        return func() E.Either[error, string] {
            url := cfg.BaseURL + path
            
            req, err := http.NewRequest("GET", url, nil)
            if err != nil {
                return E.Left[string](err)
            }
            
            req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
            
            client := &http.Client{}
            resp, err := client.Do(req)
            if err != nil {
                return E.Left[string](err)
            }
            defer resp.Body.Close()
            
            if resp.StatusCode != http.StatusOK {
                return E.Left[string](fmt.Errorf("status: %d", resp.StatusCode))
            }
            
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                return E.Left[string](err)
            }
            
            return E.Right[error](string(body))
        }
    }
}

func parseResponse(body string) RIOE.ReaderIOEither[APIConfig, error, APIResponse] {
    return func(cfg APIConfig) IOE.IOEither[error, APIResponse] {
        return func() E.Either[error, APIResponse] {
            var resp APIResponse
            if err := json.Unmarshal([]byte(body), &resp); err != nil {
                return E.Left[APIResponse](err)
            }
            return E.Right[error](resp)
        }
    }
}

func fetchData(path string) RIOE.ReaderIOEither[APIConfig, error, APIResponse] {
    return F.Pipe3(
        validateConfig(),
        RIOE.Chain(func(cfg APIConfig) RIOE.ReaderIOEither[APIConfig, error, string] {
            return makeRequest(path)
        }),
        RIOE.Chain(parseResponse),
    )
}

func main() {
    config := APIConfig{
        BaseURL: "https://api.example.com",
        APIKey:  "secret-key",
        Timeout: 30,
    }
    
    result := fetchData("/users")(config)()
    
    response := E.Fold(
        func(err error) string {
            return fmt.Sprintf("Error: %v", err)
        },
        func(resp APIResponse) string {
            return fmt.Sprintf("Success: %s", resp.Data)
        },
    )(result)
    
    fmt.Println(response)
}
```

## Database Transaction Pattern

```go
package main

import (
    "fmt"
    RIOE "github.com/IBM/fp-go/readerioeither"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
)

type Transaction interface {
    Commit() error
    Rollback() error
}

type DB interface {
    Begin() (Transaction, error)
    Execute(tx Transaction, sql string) error
}

type Deps struct {
    DB DB
}

type MockTx struct {
    committed bool
}

func (tx *MockTx) Commit() error {
    tx.committed = true
    return nil
}

func (tx *MockTx) Rollback() error {
    tx.committed = false
    return nil
}

type MockDB struct{}

func (db MockDB) Begin() (Transaction, error) {
    return &MockTx{}, nil
}

func (db MockDB) Execute(tx Transaction, sql string) error {
    fmt.Println("Executing:", sql)
    return nil
}

func beginTx() RIOE.ReaderIOEither[Deps, error, Transaction] {
    return func(deps Deps) IOE.IOEither[error, Transaction] {
        return func() E.Either[error, Transaction] {
            tx, err := deps.DB.Begin()
            if err != nil {
                return E.Left[Transaction](err)
            }
            return E.Right[error](tx)
        }
    }
}

func execute(tx Transaction, sql string) RIOE.ReaderIOEither[Deps, error, Transaction] {
    return func(deps Deps) IOE.IOEither[error, Transaction] {
        return func() E.Either[error, Transaction] {
            if err := deps.DB.Execute(tx, sql); err != nil {
                return E.Left[Transaction](err)
            }
            return E.Right[error](tx)
        }
    }
}

func commit(tx Transaction) RIOE.ReaderIOEither[Deps, error, struct{}] {
    return func(deps Deps) IOE.IOEither[error, struct{}] {
        return func() E.Either[error, struct{}] {
            if err := tx.Commit(); err != nil {
                return E.Left[struct{}](err)
            }
            return E.Right[error](struct{}{})
        }
    }
}

func runTransaction() RIOE.ReaderIOEither[Deps, error, struct{}] {
    return F.Pipe4(
        beginTx(),
        RIOE.Chain(func(tx Transaction) RIOE.ReaderIOEither[Deps, error, Transaction] {
            return execute(tx, "INSERT INTO users VALUES (1, 'Alice')")
        }),
        RIOE.Chain(func(tx Transaction) RIOE.ReaderIOEither[Deps, error, Transaction] {
            return execute(tx, "INSERT INTO posts VALUES (1, 'Hello')")
        }),
        RIOE.Chain(commit),
    )
}

func main() {
    deps := Deps{DB: MockDB{}}
    
    result := runTransaction()(deps)()
    
    if E.IsRight(result) {
        fmt.Println("Transaction committed successfully")
    }
}
```

## Migration to v2

### Key Changes

```go
// v1 and v2 are very similar for ReaderIOEither
// Main improvements are in type inference and utilities

// v1
func processV1(id int) RIOE.ReaderIOEither[Deps, error, User] {
    return func(deps Deps) IOE.IOEither[error, User] {
        return func() E.Either[error, User] {
            // implementation
            return E.Right[error](User{})
        }
    }
}

// v2 (same pattern, better inference)
func processV2(id int) RIOE.ReaderIOEither[Deps, error, User] {
    return func(deps Deps) IOE.IOEither[error, User] {
        return func() E.Either[error, User] {
            // implementation
            return E.Right[error](User{})
        }
    }
}
```

## See Also

- [Reader v1](./reader) - For environment only
- [IOEither v1](./ioeither) - For IO with errors
- [ReaderIOEither v2](../v2/readerioeither) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2