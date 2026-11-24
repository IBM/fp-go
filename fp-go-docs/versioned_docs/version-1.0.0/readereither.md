---
sidebar_position: 10
---

# ReaderEither (v1)

The `ReaderEither` type combines `Reader` and `Either` for computations that depend on an environment and can fail.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [ReaderEither v2](../v2/readereither).

**Key differences in v2:**
- Simplified API
- Better type inference
- Improved composition
:::

## Overview

`ReaderEither` represents a computation that:
- Depends on a shared environment (Reader)
- Can fail with an error (Either Left)
- Can succeed with a value (Either Right)

```go
type ReaderEither[R, E, A any] func(R) Either[E, A]
```

## Creating ReaderEither Values

### Of (Success)

```go
package main

import (
    "fmt"
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    Debug bool
}

func main() {
    // Create successful ReaderEither
    re := RE.Of[Config, error, int](42)
    
    config := Config{Debug: true}
    result := re(config)
    
    fmt.Println(E.IsRight(result)) // true
}
```

### Left (Error)

```go
package main

import (
    "errors"
    "fmt"
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    Debug bool
}

func main() {
    // Create failed ReaderEither
    re := RE.Left[Config, int](errors.New("failed"))
    
    config := Config{Debug: true}
    result := re(config)
    
    fmt.Println(E.IsLeft(result)) // true
}
```

### Ask

Access the environment:

```go
package main

import (
    "fmt"
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    APIKey string
}

func main() {
    // Get environment as Right value
    re := RE.Ask[Config, error]()
    
    config := Config{APIKey: "secret"}
    result := re(config)
    
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
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    Multiplier int
}

func main() {
    re := RE.Of[Config, error, int](5)
    
    // Map transforms the Right value
    doubled := RE.Map(func(n int) int {
        return n * 2
    })(re)
    
    config := Config{Multiplier: 2}
    result := doubled(config)
    
    value := E.GetOrElse(func() int { return 0 })(result)
    fmt.Println(value) // 10
}
```

### Chain

Chain operations that return ReaderEither:

```go
package main

import (
    "fmt"
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    MaxValue int
}

func validate(n int) RE.ReaderEither[Config, error, int] {
    return func(cfg Config) E.Either[error, int] {
        if n > cfg.MaxValue {
            return E.Left[int](fmt.Errorf("exceeds max: %d", cfg.MaxValue))
        }
        return E.Right[error](n)
    }
}

func main() {
    re := RE.Of[Config, error, int](100)
    
    validated := RE.Chain(validate)(re)
    
    config := Config{MaxValue: 50}
    result := validated(config)
    
    fmt.Println(E.IsLeft(result)) // true (exceeds max)
}
```

## Dependency Injection

### Database Service

```go
package main

import (
    "fmt"
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
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

func getUser(id int) RE.ReaderEither[Env, error, *User] {
    return func(env Env) E.Either[error, *User] {
        user, err := env.DB.FindUser(id)
        if err != nil {
            return E.Left[*User](err)
        }
        return E.Right[error](user)
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
    
    result := getUser(1)(env)
    
    if E.IsRight(result) {
        user := E.GetOrElse(func() *User { return nil })(result)
        fmt.Printf("Found: %s\n", user.Name)
    }
}
```

### Configuration Validation

```go
package main

import (
    "fmt"
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
)

type Config struct {
    Host string
    Port int
}

func validateHost() RE.ReaderEither[Config, error, string] {
    return func(cfg Config) E.Either[error, string] {
        if cfg.Host == "" {
            return E.Left[string](fmt.Errorf("host is required"))
        }
        return E.Right[error](cfg.Host)
    }
}

func validatePort() RE.ReaderEither[Config, error, int] {
    return func(cfg Config) E.Either[error, int] {
        if cfg.Port <= 0 || cfg.Port > 65535 {
            return E.Left[int](fmt.Errorf("invalid port"))
        }
        return E.Right[error](cfg.Port)
    }
}

func main() {
    config := Config{Host: "localhost", Port: 8080}
    
    hostResult := validateHost()(config)
    portResult := validatePort()(config)
    
    fmt.Println("Host valid:", E.IsRight(hostResult))
    fmt.Println("Port valid:", E.IsRight(portResult))
}
```

## Practical Examples

### API Client with Config

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
)

type APIConfig struct {
    BaseURL string
    APIKey  string
}

func makeRequest(path string) RE.ReaderEither[APIConfig, error, string] {
    return func(cfg APIConfig) E.Either[error, string] {
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
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return E.Left[string](err)
        }
        
        return E.Right[error](string(body))
    }
}

func main() {
    config := APIConfig{
        BaseURL: "https://api.example.com",
        APIKey:  "secret-key",
    }
    
    result := makeRequest("/users")(config)
    
    response := E.Fold(
        func(err error) string {
            return fmt.Sprintf("Error: %v", err)
        },
        func(body string) string {
            return body
        },
    )(result)
    
    fmt.Println(response)
}
```

### Multi-Service Application

```go
package main

import (
    "fmt"
    RE "github.com/IBM/fp-go/readereither"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
)

type Logger interface {
    Info(msg string)
}

type UserRepo interface {
    FindByID(id int) (string, error)
}

type EmailService interface {
    Send(to, subject, body string) error
}

type Deps struct {
    Logger       Logger
    UserRepo     UserRepo
    EmailService EmailService
}

type ConsoleLogger struct{}

func (l ConsoleLogger) Info(msg string) {
    fmt.Println("[INFO]", msg)
}

type MockUserRepo struct{}

func (r MockUserRepo) FindByID(id int) (string, error) {
    return fmt.Sprintf("user%d@example.com", id), nil
}

type MockEmailService struct{}

func (s MockEmailService) Send(to, subject, body string) error {
    fmt.Printf("Email to %s: %s\n", to, subject)
    return nil
}

func getUserEmail(userID int) RE.ReaderEither[Deps, error, string] {
    return func(deps Deps) E.Either[error, string] {
        deps.Logger.Info(fmt.Sprintf("Getting email for user %d", userID))
        
        email, err := deps.UserRepo.FindByID(userID)
        if err != nil {
            return E.Left[string](err)
        }
        
        return E.Right[error](email)
    }
}

func sendWelcomeEmail(email string) RE.ReaderEither[Deps, error, struct{}] {
    return func(deps Deps) E.Either[error, struct{}] {
        deps.Logger.Info(fmt.Sprintf("Sending welcome email to %s", email))
        
        err := deps.EmailService.Send(email, "Welcome!", "Welcome to our service")
        if err != nil {
            return E.Left[struct{}](err)
        }
        
        return E.Right[error](struct{}{})
    }
}

func main() {
    deps := Deps{
        Logger:       ConsoleLogger{},
        UserRepo:     MockUserRepo{},
        EmailService: MockEmailService{},
    }
    
    // Chain operations
    sendWelcome := F.Pipe2(
        getUserEmail(1),
        RE.Chain(sendWelcomeEmail),
    )
    
    result := sendWelcome(deps)
    
    if E.IsRight(result) {
        fmt.Println("Welcome email sent successfully")
    }
}
```

## Migration to v2

### Key Changes

```go
// v1 and v2 are very similar for ReaderEither
// Main improvements are in type inference and composition

// v1
func processV1(id int) RE.ReaderEither[Env, error, User] {
    return func(env Env) E.Either[error, User] {
        // implementation
        return E.Right[error](User{})
    }
}

// v2 (same pattern)
func processV2(id int) RE.ReaderEither[Env, error, User] {
    return func(env Env) E.Either[error, User] {
        // implementation
        return E.Right[error](User{})
    }
}
```

## See Also

- [Reader v1](./reader) - For environment without errors
- [Either v1](./either) - For error handling
- [ReaderEither v2](../v2/readereither) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2