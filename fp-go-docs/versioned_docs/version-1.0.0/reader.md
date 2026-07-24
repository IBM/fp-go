---
sidebar_position: 8
---

# Reader (v1)

The `Reader` type represents a computation that depends on a shared environment.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [Reader v2](../v2/reader).

**Key differences in v2:**
- Simplified API
- Better type inference
- Improved composition
- More utility functions
:::

## Overview

`Reader` is a function that:
- Takes an environment/context as input
- Returns a computed value
- Enables dependency injection
- Supports composition

```go
type Reader[R, A any] func(R) A
```

## Creating Reader Values

### Of (Pure Value)

Wrap a value that ignores the environment:

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

func main() {
    // Create Reader that returns constant
    reader := R.Of[string, int](42)
    
    // Run with any environment
    result := reader("any context")
    fmt.Println(result) // 42
}
```

### Ask

Access the environment:

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

func main() {
    // Create Reader that returns the environment
    reader := R.Ask[string]()
    
    // Run with environment
    result := reader("Hello, World!")
    fmt.Println(result) // Hello, World!
}
```

### From Function

Create from a function:

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

type Config struct {
    Host string
    Port int
}

func main() {
    // Create Reader that extracts from config
    getHost := R.Reader[Config, string](func(cfg Config) string {
        return cfg.Host
    })
    
    config := Config{Host: "localhost", Port: 8080}
    host := getHost(config)
    fmt.Println(host) // localhost
}
```

## Basic Operations

### Map

Transform the result:

```go
package main

import (
    "fmt"
    "strings"
    R "github.com/IBM/fp-go/reader"
)

type Config struct {
    Name string
}

func main() {
    // Reader that gets name
    getName := R.Reader[Config, string](func(cfg Config) string {
        return cfg.Name
    })
    
    // Map to uppercase
    getUpperName := R.Map(strings.ToUpper)(getName)
    
    config := Config{Name: "alice"}
    result := getUpperName(config)
    fmt.Println(result) // ALICE
}
```

### Chain (FlatMap)

Chain Reader operations:

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

type Config struct {
    Host string
    Port int
}

func getHost() R.Reader[Config, string] {
    return func(cfg Config) string {
        return cfg.Host
    }
}

func getURL(host string) R.Reader[Config, string] {
    return func(cfg Config) string {
        return fmt.Sprintf("http://%s:%d", host, cfg.Port)
    }
}

func main() {
    // Chain readers
    getFullURL := R.Chain(getURL)(getHost())
    
    config := Config{Host: "localhost", Port: 8080}
    url := getFullURL(config)
    fmt.Println(url) // http://localhost:8080
}
```

### Local

Transform the environment:

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

type AppConfig struct {
    Database string
    Cache    string
}

type DBConfig struct {
    ConnectionString string
}

func main() {
    // Reader that needs DBConfig
    getConnection := R.Reader[DBConfig, string](func(cfg DBConfig) string {
        return cfg.ConnectionString
    })
    
    // Transform AppConfig to DBConfig
    getDBConnection := R.Local(func(app AppConfig) DBConfig {
        return DBConfig{ConnectionString: app.Database}
    })(getConnection)
    
    appConfig := AppConfig{
        Database: "postgres://localhost",
        Cache:    "redis://localhost",
    }
    
    conn := getDBConnection(appConfig)
    fmt.Println(conn) // postgres://localhost
}
```

## Dependency Injection

### Service Pattern

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

type Logger interface {
    Log(msg string)
}

type ConsoleLogger struct{}

func (l ConsoleLogger) Log(msg string) {
    fmt.Println("[LOG]", msg)
}

type Env struct {
    Logger Logger
}

func logMessage(msg string) R.Reader[Env, struct{}] {
    return func(env Env) struct{} {
        env.Logger.Log(msg)
        return struct{}{}
    }
}

func main() {
    env := Env{Logger: ConsoleLogger{}}
    
    // Run with environment
    logMessage("Application started")(env)
    // Output: [LOG] Application started
}
```

### Database Access

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

type User struct {
    ID   int
    Name string
}

type DB interface {
    FindUser(id int) (*User, error)
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

type Env struct {
    DB DB
}

func getUser(id int) R.Reader[Env, *User] {
    return func(env Env) *User {
        user, _ := env.DB.FindUser(id)
        return user
    }
}

func main() {
    env := Env{
        DB: MockDB{
            users: map[int]User{
                1: {ID: 1, Name: "Alice"},
            },
        },
    }
    
    user := getUser(1)(env)
    if user != nil {
        fmt.Printf("Found: %s\n", user.Name)
    }
}
```

## Composition

### Sequential Operations

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
    F "github.com/IBM/fp-go/function"
)

type Config struct {
    Prefix string
    Suffix string
}

func addPrefix(s string) R.Reader[Config, string] {
    return func(cfg Config) string {
        return cfg.Prefix + s
    }
}

func addSuffix(s string) R.Reader[Config, string] {
    return func(cfg Config) string {
        return s + cfg.Suffix
    }
}

func main() {
    // Compose operations
    process := F.Pipe2(
        R.Of[Config, string]("Hello"),
        R.Chain(addPrefix),
        R.Chain(addSuffix),
    )
    
    config := Config{Prefix: "[", Suffix: "]"}
    result := process(config)
    fmt.Println(result) // [Hello]
}
```

### Combining Readers

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

type Config struct {
    FirstName string
    LastName  string
}

func getFirstName() R.Reader[Config, string] {
    return func(cfg Config) string {
        return cfg.FirstName
    }
}

func getLastName() R.Reader[Config, string] {
    return func(cfg Config) string {
        return cfg.LastName
    }
}

func getFullName() R.Reader[Config, string] {
    return R.Chain(func(first string) R.Reader[Config, string] {
        return R.Map(func(last string) string {
            return first + " " + last
        })(getLastName())
    })(getFirstName())
}

func main() {
    config := Config{
        FirstName: "John",
        LastName:  "Doe",
    }
    
    fullName := getFullName()(config)
    fmt.Println(fullName) // John Doe
}
```

## Practical Examples

### Configuration Management

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

type AppConfig struct {
    DatabaseURL string
    APIKey      string
    Debug       bool
}

func getDatabaseURL() R.Reader[AppConfig, string] {
    return func(cfg AppConfig) string {
        return cfg.DatabaseURL
    }
}

func getAPIKey() R.Reader[AppConfig, string] {
    return func(cfg AppConfig) string {
        return cfg.APIKey
    }
}

func isDebugMode() R.Reader[AppConfig, bool] {
    return func(cfg AppConfig) bool {
        return cfg.Debug
    }
}

func main() {
    config := AppConfig{
        DatabaseURL: "postgres://localhost/mydb",
        APIKey:      "secret-key",
        Debug:       true,
    }
    
    dbURL := getDatabaseURL()(config)
    apiKey := getAPIKey()(config)
    debug := isDebugMode()(config)
    
    fmt.Println("Database:", dbURL)
    fmt.Println("API Key:", apiKey)
    fmt.Println("Debug:", debug)
}
```

### HTTP Handler with Dependencies

```go
package main

import (
    "fmt"
    "net/http"
    R "github.com/IBM/fp-go/reader"
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

func handleGetUser(userID string) R.Reader[Deps, http.HandlerFunc] {
    return func(deps Deps) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            deps.Logger.Info(fmt.Sprintf("Getting user %s", userID))
            
            user, err := deps.UserService.GetUser(userID)
            if err != nil {
                http.Error(w, err.Error(), http.StatusNotFound)
                return
            }
            
            fmt.Fprintf(w, "User: %s", user)
        }
    }
}

func main() {
    deps := Deps{
        Logger:      ConsoleLogger{},
        UserService: MockUserService{},
    }
    
    handler := handleGetUser("123")(deps)
    
    // Use handler with http.Server
    fmt.Println("Handler created with dependencies")
}
```

### Multi-Layer Application

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/reader"
)

type Config struct {
    Environment string
}

type Logger interface {
    Log(msg string)
}

type Repository interface {
    Save(data string) error
}

type AppDeps struct {
    Config     Config
    Logger     Logger
    Repository Repository
}

type ConsoleLogger struct{}

func (l ConsoleLogger) Log(msg string) {
    fmt.Println("[LOG]", msg)
}

type MemoryRepo struct {
    data []string
}

func (r *MemoryRepo) Save(data string) error {
    r.data = append(r.data, data)
    return nil
}

func saveData(data string) R.Reader[AppDeps, error] {
    return func(deps AppDeps) error {
        deps.Logger.Log(fmt.Sprintf("Saving data in %s", deps.Config.Environment))
        return deps.Repository.Save(data)
    }
}

func main() {
    deps := AppDeps{
        Config:     Config{Environment: "production"},
        Logger:     ConsoleLogger{},
        Repository: &MemoryRepo{},
    }
    
    err := saveData("important data")(deps)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

## Migration to v2

### Key Changes

1. **Simplified API**:
```go
// v1 and v2 are very similar for Reader
// Main improvements are in type inference

// v1
reader := R.Reader[Config, string](func(cfg Config) string {
    return cfg.Name
})

// v2 (same)
reader := R.Reader[Config, string](func(cfg Config) string {
    return cfg.Name
})
```

2. **Better composition**:
```go
// v2 has improved pipe utilities
result := F.Pipe2(
    getConfig(),
    R.Chain(processConfig),
    R.Map(formatOutput),
)
```

### Migration Example

```go
// v1 code
func getUserV1(id int) R.Reader[Env, *User] {
    return func(env Env) *User {
        user, _ := env.DB.FindUser(id)
        return user
    }
}

// v2 equivalent (mostly the same)
func getUserV2(id int) R.Reader[Env, *User] {
    return func(env Env) *User {
        user, _ := env.DB.FindUser(id)
        return user
    }
}
```

## See Also

- [ReaderEither v1](./readereither) - Reader with error handling
- [ReaderIO v1](./readerio) - Reader with side effects
- [Reader v2](../v2/reader) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2