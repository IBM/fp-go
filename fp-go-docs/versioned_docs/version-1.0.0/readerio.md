---
sidebar_position: 11
---

# ReaderIO (v1)

The `ReaderIO` type combines `Reader` and `IO` for lazy computations that depend on an environment and perform side effects.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [ReaderIO v2](../v2/readerio).

**Key differences in v2:**
- Simplified API
- Better type inference
- Improved composition
:::

## Overview

`ReaderIO` represents a computation that:
- Depends on a shared environment (Reader)
- Performs side effects (IO)
- Is lazy (deferred execution)

```go
type ReaderIO[R, A any] func(R) IO[A]
```

## Creating ReaderIO Values

### Of (Pure Value)

```go
package main

import (
    "fmt"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
)

type Config struct {
    Debug bool
}

func main() {
    // Create ReaderIO that returns constant
    rio := RIO.Of[Config, int](42)
    
    config := Config{Debug: true}
    io := rio(config)
    result := io()
    
    fmt.Println(result) // 42
}
```

### Ask

Access the environment:

```go
package main

import (
    "fmt"
    RIO "github.com/IBM/fp-go/readerio"
)

type Config struct {
    AppName string
}

func main() {
    // Get environment
    rio := RIO.Ask[Config]()
    
    config := Config{AppName: "MyApp"}
    io := rio(config)
    result := io()
    
    fmt.Println("App:", result.AppName)
}
```

### From IO

Lift an IO into ReaderIO:

```go
package main

import (
    "fmt"
    "time"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
)

type Config struct {
    Debug bool
}

func main() {
    // Lift IO to ReaderIO
    getTime := IO.IO[time.Time](func() time.Time {
        return time.Now()
    })
    
    rio := RIO.FromIO[Config](getTime)
    
    config := Config{Debug: true}
    io := rio(config)
    result := io()
    
    fmt.Println("Time:", result)
}
```

## Basic Operations

### Map

Transform the result:

```go
package main

import (
    "fmt"
    RIO "github.com/IBM/fp-go/readerio"
)

type Config struct {
    Multiplier int
}

func main() {
    rio := RIO.Of[Config, int](5)
    
    // Map transforms the result
    doubled := RIO.Map(func(n int) int {
        return n * 2
    })(rio)
    
    config := Config{Multiplier: 2}
    result := doubled(config)()
    
    fmt.Println(result) // 10
}
```

### Chain

Chain operations that return ReaderIO:

```go
package main

import (
    "fmt"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
)

type Config struct {
    Prefix string
}

func addPrefix(s string) RIO.ReaderIO[Config, string] {
    return func(cfg Config) IO.IO[string] {
        return func() string {
            return cfg.Prefix + s
        }
    }
}

func main() {
    rio := RIO.Of[Config, string]("Hello")
    
    withPrefix := RIO.Chain(addPrefix)(rio)
    
    config := Config{Prefix: "[LOG] "}
    result := withPrefix(config)()
    
    fmt.Println(result) // [LOG] Hello
}
```

## Side Effects with Environment

### Logging with Config

```go
package main

import (
    "fmt"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
)

type Config struct {
    Debug   bool
    LogFile string
}

func log(msg string) RIO.ReaderIO[Config, struct{}] {
    return func(cfg Config) IO.IO[struct{}] {
        return func() struct{} {
            if cfg.Debug {
                fmt.Printf("[DEBUG] %s\n", msg)
            }
            return struct{}{}
        }
    }
}

func main() {
    config := Config{Debug: true, LogFile: "app.log"}
    
    // Execute logging
    log("Application started")(config)()
    log("Processing request")(config)()
    
    // Output:
    // [DEBUG] Application started
    // [DEBUG] Processing request
}
```

### Database Operations

```go
package main

import (
    "fmt"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
)

type User struct {
    ID   int
    Name string
}

type DB interface {
    Query(sql string) ([]User, error)
}

type MockDB struct{}

func (db MockDB) Query(sql string) ([]User, error) {
    return []User{{ID: 1, Name: "Alice"}}, nil
}

type Env struct {
    DB DB
}

func getUsers() RIO.ReaderIO[Env, []User] {
    return func(env Env) IO.IO[[]User] {
        return func() []User {
            users, _ := env.DB.Query("SELECT * FROM users")
            return users
        }
    }
}

func main() {
    env := Env{DB: MockDB{}}
    
    users := getUsers()(env)()
    
    for _, user := range users {
        fmt.Printf("User: %s\n", user.Name)
    }
}
```

## Practical Examples

### HTTP Server with Dependencies

```go
package main

import (
    "fmt"
    "net/http"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
)

type Logger interface {
    Info(msg string)
}

type UserService interface {
    GetUser(id string) (string, error)
}

type Deps struct {
    Logger      Logger
    UserService UserService
}

type ConsoleLogger struct{}

func (l ConsoleLogger) Info(msg string) {
    fmt.Println("[INFO]", msg)
}

type MockUserService struct{}

func (s MockUserService) GetUser(id string) (string, error) {
    return fmt.Sprintf("User-%s", id), nil
}

func handleRequest(userID string) RIO.ReaderIO[Deps, http.HandlerFunc] {
    return func(deps Deps) IO.IO[http.HandlerFunc] {
        return func() http.HandlerFunc {
            return func(w http.ResponseWriter, r *http.Request) {
                deps.Logger.Info(fmt.Sprintf("Handling request for user %s", userID))
                
                user, err := deps.UserService.GetUser(userID)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusNotFound)
                    return
                }
                
                fmt.Fprintf(w, "User: %s", user)
            }
        }
    }
}

func main() {
    deps := Deps{
        Logger:      ConsoleLogger{},
        UserService: MockUserService{},
    }
    
    handler := handleRequest("123")(deps)()
    
    fmt.Println("Handler created with dependencies")
    _ = handler
}
```

### File Operations with Config

```go
package main

import (
    "fmt"
    "os"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
)

type Config struct {
    DataDir string
}

func readFile(filename string) RIO.ReaderIO[Config, string] {
    return func(cfg Config) IO.IO[string] {
        return func() string {
            path := cfg.DataDir + "/" + filename
            data, err := os.ReadFile(path)
            if err != nil {
                return ""
            }
            return string(data)
        }
    }
}

func writeFile(filename, content string) RIO.ReaderIO[Config, error] {
    return func(cfg Config) IO.IO[error] {
        return func() error {
            path := cfg.DataDir + "/" + filename
            return os.WriteFile(path, []byte(content), 0644)
        }
    }
}

func main() {
    config := Config{DataDir: "/tmp"}
    
    // Write file
    writeFile("test.txt", "Hello, World!")(config)()
    
    // Read file
    content := readFile("test.txt")(config)()
    fmt.Println("Content:", content)
}
```

### Caching with Environment

```go
package main

import (
    "fmt"
    "sync"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
)

type Cache interface {
    Get(key string) (string, bool)
    Set(key, value string)
}

type MemoryCache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *MemoryCache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}

func (c *MemoryCache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

type Env struct {
    Cache Cache
}

func getCached(key string) RIO.ReaderIO[Env, *string] {
    return func(env Env) IO.IO[*string] {
        return func() *string {
            if val, ok := env.Cache.Get(key); ok {
                return &val
            }
            return nil
        }
    }
}

func setCached(key, value string) RIO.ReaderIO[Env, struct{}] {
    return func(env Env) IO.IO[struct{}] {
        return func() struct{} {
            env.Cache.Set(key, value)
            return struct{}{}
        }
    }
}

func main() {
    env := Env{
        Cache: &MemoryCache{data: make(map[string]string)},
    }
    
    // Set value
    setCached("key1", "value1")(env)()
    
    // Get value
    result := getCached("key1")(env)()
    if result != nil {
        fmt.Println("Cached:", *result)
    }
}
```

## Composition

### Sequential Operations

```go
package main

import (
    "fmt"
    RIO "github.com/IBM/fp-go/readerio"
    IO "github.com/IBM/fp-go/io"
    F "github.com/IBM/fp-go/function"
)

type Config struct {
    Step1 string
    Step2 string
}

func step1() RIO.ReaderIO[Config, string] {
    return func(cfg Config) IO.IO[string] {
        return func() string {
            fmt.Println("Step 1:", cfg.Step1)
            return cfg.Step1
        }
    }
}

func step2(input string) RIO.ReaderIO[Config, string] {
    return func(cfg Config) IO.IO[string] {
        return func() string {
            fmt.Println("Step 2:", cfg.Step2)
            return input + " -> " + cfg.Step2
        }
    }
}

func main() {
    config := Config{Step1: "Init", Step2: "Process"}
    
    // Chain steps
    pipeline := F.Pipe2(
        step1(),
        RIO.Chain(step2),
    )
    
    result := pipeline(config)()
    fmt.Println("Result:", result)
}
```

## Migration to v2

### Key Changes

```go
// v1 and v2 are very similar for ReaderIO
// Main improvements are in type inference

// v1
func processV1() RIO.ReaderIO[Config, int] {
    return func(cfg Config) IO.IO[int] {
        return func() int {
            return 42
        }
    }
}

// v2 (same pattern)
func processV2() RIO.ReaderIO[Config, int] {
    return func(cfg Config) IO.IO[int] {
        return func() int {
            return 42
        }
    }
}
```

## See Also

- [Reader v1](./reader) - For environment without IO
- [IO v1](./io) - For side effects without environment
- [ReaderIOEither v1](./readerioeither) - For environment, IO, and errors
- [ReaderIO v2](../v2/readerio) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2