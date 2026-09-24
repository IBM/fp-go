---
sidebar_position: 7
---

# IOEither (v1)

The `IOEither` type combines `IO` and `Either` for lazy computations that can fail.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [IOEither v2](../v2/ioeither).

**Key differences in v2:**
- Simplified API
- Better type inference
- Improved error handling
- More utility functions
:::

## Overview

`IOEither` represents a lazy computation that:
- Defers execution until explicitly run
- Can fail with a Left value (error)
- Can succeed with a Right value (result)

```go
type IOEither[E, A any] func() Either[E, A]
```

## Creating IOEither Values

### Of (Success)

Wrap a successful value:

```go
package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func main() {
    // Create successful IOEither
    io := IOE.Of[error, int](42)
    
    // Execute to get Either
    result := io()
    
    if E.IsRight(result) {
        fmt.Println("Success:", E.GetOrElse(func() int { return 0 })(result))
    }
}
```

### Left (Error)

Create a failed computation:

```go
package main

import (
    "errors"
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func main() {
    // Create failed IOEither
    io := IOE.Left[int](errors.New("something went wrong"))
    
    // Execute to get Either
    result := io()
    
    if E.IsLeft(result) {
        fmt.Println("Error occurred")
    }
}
```

### TryCatch

Wrap a function that might panic:

```go
package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func riskyOperation() int {
    // Might panic
    return 42 / 0
}

func main() {
    // Catch panics and convert to Left
    io := IOE.TryCatch(
        riskyOperation,
        func(err any) error {
            return fmt.Errorf("panic: %v", err)
        },
    )
    
    result := io()
    fmt.Println("Is error:", E.IsLeft(result))
}
```

## Basic Operations

### Map

Transform the success value:

```go
package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func main() {
    io := IOE.Of[error, int](5)
    
    // Map transforms the Right value
    doubled := IOE.Map(func(n int) int {
        return n * 2
    })(io)
    
    result := doubled()
    value := E.GetOrElse(func() int { return 0 })(result)
    fmt.Println(value) // 10
}
```

### Chain (FlatMap)

Chain operations that return IOEither:

```go
package main

import (
    "errors"
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func divide(a, b int) IOE.IOEither[error, int] {
    return func() E.Either[error, int] {
        if b == 0 {
            return E.Left[int](errors.New("division by zero"))
        }
        return E.Right[error](a / b)
    }
}

func main() {
    io := IOE.Of[error, int](10)
    
    result := IOE.Chain(func(n int) IOE.IOEither[error, int] {
        return divide(n, 2)
    })(io)
    
    either := result()
    value := E.GetOrElse(func() int { return 0 })(either)
    fmt.Println(value) // 5
}
```

### MapLeft

Transform the error value:

```go
package main

import (
    "errors"
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func main() {
    io := IOE.Left[int](errors.New("original error"))
    
    // Transform the error
    mapped := IOE.MapLeft(func(err error) error {
        return fmt.Errorf("wrapped: %w", err)
    })(io)
    
    result := mapped()
    fmt.Println("Is error:", E.IsLeft(result))
}
```

## Error Handling

### Fold

Handle both success and error:

```go
package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func processResult(io IOE.IOEither[error, int]) string {
    either := io()
    
    return E.Fold(
        func(err error) string {
            return fmt.Sprintf("Error: %v", err)
        },
        func(value int) string {
            return fmt.Sprintf("Success: %d", value)
        },
    )(either)
}

func main() {
    success := IOE.Of[error, int](42)
    failure := IOE.Left[int](fmt.Errorf("failed"))
    
    fmt.Println(processResult(success)) // Success: 42
    fmt.Println(processResult(failure)) // Error: failed
}
```

### OrElse

Provide fallback on error:

```go
package main

import (
    "errors"
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func main() {
    failing := IOE.Left[int](errors.New("error"))
    fallback := IOE.Of[error, int](42)
    
    // Use fallback if first fails
    result := IOE.OrElse(func(err error) IOE.IOEither[error, int] {
        return fallback
    })(failing)
    
    either := result()
    value := E.GetOrElse(func() int { return 0 })(either)
    fmt.Println(value) // 42
}
```

## File Operations

### Reading Files

```go
package main

import (
    "fmt"
    "os"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func readFile(path string) IOE.IOEither[error, string] {
    return func() E.Either[error, string] {
        data, err := os.ReadFile(path)
        if err != nil {
            return E.Left[string](err)
        }
        return E.Right[error](string(data))
    }
}

func main() {
    io := readFile("config.json")
    result := io()
    
    content := E.Fold(
        func(err error) string {
            return fmt.Sprintf("Error: %v", err)
        },
        func(data string) string {
            return data
        },
    )(result)
    
    fmt.Println(content)
}
```

### Writing Files

```go
package main

import (
    "fmt"
    "os"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func writeFile(path, content string) IOE.IOEither[error, struct{}] {
    return func() E.Either[error, struct{}] {
        err := os.WriteFile(path, []byte(content), 0644)
        if err != nil {
            return E.Left[struct{}](err)
        }
        return E.Right[error](struct{}{})
    }
}

func main() {
    io := writeFile("output.txt", "Hello, World!")
    result := io()
    
    if E.IsRight(result) {
        fmt.Println("File written successfully")
    } else {
        fmt.Println("Failed to write file")
    }
}
```

## HTTP Requests

### GET Request

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func httpGet(url string) IOE.IOEither[error, string] {
    return func() E.Either[error, string] {
        resp, err := http.Get(url)
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
    io := httpGet("https://api.example.com/data")
    result := io()
    
    response := E.GetOrElse(func() string {
        return "Failed to fetch"
    })(result)
    
    fmt.Println(response)
}
```

### POST Request

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func httpPost(url, body string) IOE.IOEither[error, string] {
    return func() E.Either[error, string] {
        resp, err := http.Post(
            url,
            "application/json",
            bytes.NewBufferString(body),
        )
        if err != nil {
            return E.Left[string](err)
        }
        defer resp.Body.Close()
        
        respBody, err := io.ReadAll(resp.Body)
        if err != nil {
            return E.Left[string](err)
        }
        
        return E.Right[error](string(respBody))
    }
}

func main() {
    io := httpPost(
        "https://api.example.com/users",
        `{"name": "Alice"}`,
    )
    
    result := io()
    
    if E.IsRight(result) {
        fmt.Println("User created")
    }
}
```

## Composition

### Sequential Operations

```go
package main

import (
    "fmt"
    "strconv"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
)

func parseNumber(s string) IOE.IOEither[error, int] {
    return func() E.Either[error, int] {
        n, err := strconv.Atoi(s)
        if err != nil {
            return E.Left[int](err)
        }
        return E.Right[error](n)
    }
}

func validatePositive(n int) IOE.IOEither[error, int] {
    return func() E.Either[error, int] {
        if n <= 0 {
            return E.Left[int](fmt.Errorf("must be positive"))
        }
        return E.Right[error](n)
    }
}

func main() {
    // Compose operations
    result := F.Pipe2(
        parseNumber("42"),
        IOE.Chain(validatePositive),
    )
    
    either := result()
    value := E.GetOrElse(func() int { return 0 })(either)
    fmt.Println(value) // 42
}
```

### Parallel Operations

```go
package main

import (
    "fmt"
    "sync"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func fetchUser(id int) IOE.IOEither[error, string] {
    return func() E.Either[error, string] {
        return E.Right[error](fmt.Sprintf("User-%d", id))
    }
}

func fetchPosts(userId string) IOE.IOEither[error, []string] {
    return func() E.Either[error, []string] {
        return E.Right[error]([]string{"Post1", "Post2"})
    }
}

func main() {
    var wg sync.WaitGroup
    var userResult E.Either[error, string]
    var postsResult E.Either[error, []string]
    
    wg.Add(2)
    
    go func() {
        defer wg.Done()
        userResult = fetchUser(1)()
    }()
    
    go func() {
        defer wg.Done()
        postsResult = fetchPosts("User-1")()
    }()
    
    wg.Wait()
    
    if E.IsRight(userResult) && E.IsRight(postsResult) {
        fmt.Println("Both succeeded")
    }
}
```

## Practical Examples

### Database Operations

```go
package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

type User struct {
    ID   int
    Name string
}

type DB struct {
    users map[int]User
}

func (db *DB) findUser(id int) IOE.IOEither[error, User] {
    return func() E.Either[error, User] {
        if user, ok := db.users[id]; ok {
            return E.Right[error](user)
        }
        return E.Left[User](fmt.Errorf("user not found"))
    }
}

func (db *DB) saveUser(user User) IOE.IOEither[error, User] {
    return func() E.Either[error, User] {
        db.users[user.ID] = user
        return E.Right[error](user)
    }
}

func main() {
    db := &DB{users: make(map[int]User)}
    
    // Save and retrieve user
    result := IOE.Chain(func(user User) IOE.IOEither[error, User] {
        return db.findUser(user.ID)
    })(db.saveUser(User{ID: 1, Name: "Alice"}))
    
    either := result()
    
    if E.IsRight(either) {
        fmt.Println("User saved and retrieved")
    }
}
```

### Retry Logic

```go
package main

import (
    "fmt"
    "time"
    IOE "github.com/IBM/fp-go/ioeither"
    E "github.com/IBM/fp-go/either"
)

func retryable(attempt int) IOE.IOEither[error, string] {
    return func() E.Either[error, string] {
        if attempt < 3 {
            return E.Left[string](fmt.Errorf("attempt %d failed", attempt))
        }
        return E.Right[error]("success")
    }
}

func retry(io IOE.IOEither[error, string], maxAttempts int) IOE.IOEither[error, string] {
    return func() E.Either[error, string] {
        for i := 0; i < maxAttempts; i++ {
            result := io()
            if E.IsRight(result) {
                return result
            }
            time.Sleep(100 * time.Millisecond)
        }
        return E.Left[string](fmt.Errorf("max retries exceeded"))
    }
}

func main() {
    io := retryable(3)
    result := retry(io, 5)()
    
    value := E.GetOrElse(func() string {
        return "failed"
    })(result)
    
    fmt.Println(value)
}
```

## Migration to v2

### Key Changes

1. **Simplified constructors**:
```go
// v1
IOE.Of[error, int](42)
IOE.Left[int](errors.New("error"))

// v2 (same pattern, better inference)
IOE.Of[error, int](42)
IOE.Left[int](errors.New("error"))
```

2. **Better composition**:
```go
// v2 has improved pipe and flow utilities
result := F.Pipe2(
    parseInput(input),
    IOE.Chain(validate),
    IOE.Chain(process),
)
```

### Migration Example

```go
// v1 code
func processV1(path string) IOE.IOEither[error, string] {
    return IOE.Chain(func(content string) IOE.IOEither[error, string] {
        return IOE.Of[error, string](content)
    })(readFile(path))
}

// v2 equivalent (mostly the same)
func processV2(path string) IOE.IOEither[error, string] {
    return F.Pipe1(
        readFile(path),
        IOE.Chain(func(content string) IOE.IOEither[error, string] {
            return IOE.Of[error, string](content)
        }),
    )
}
```

## See Also

- [IO v1](./io) - For side effects without errors
- [Either v1](./either) - For error handling without IO
- [IOEither v2](../v2/ioeither) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2