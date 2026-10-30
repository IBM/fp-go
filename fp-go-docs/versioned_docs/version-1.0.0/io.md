---
sidebar_position: 6
---

# IO (v1)

The `IO` type represents a lazy computation that performs side effects.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [IO v2](../v2/io).

**Key differences in v2:**
- Simplified API
- Better error handling integration
- Improved performance
- More consistent naming
:::

## Overview

`IO` is a lazy computation that:
- Defers execution until explicitly run
- Can perform side effects (I/O, mutations, etc.)
- Is composable and testable

```go
type IO[A any] func() A
```

## Creating IO Values

### Of (Pure Value)

Wrap a pure value:

```go
package main

import (
    "fmt"
    IO "github.com/IBM/fp-go/io"
)

func main() {
    // Create IO that returns a value
    io := IO.Of(42)
    
    // Execute to get the value
    result := io()
    fmt.Println(result) // 42
}
```

### From Function

Create from a function:

```go
package main

import (
    "fmt"
    "time"
    IO "github.com/IBM/fp-go/io"
)

func main() {
    // Create IO that gets current time
    getTime := func() time.Time {
        return time.Now()
    }
    
    io := IO.IO[time.Time](getTime)
    
    // Execute multiple times
    time1 := io()
    time.Sleep(100 * time.Millisecond)
    time2 := io()
    
    fmt.Println("Different times:", !time1.Equal(time2))
}
```

## Basic Operations

### Map

Transform the result:

```go
package main

import (
    "fmt"
    IO "github.com/IBM/fp-go/io"
)

func main() {
    io := IO.Of(5)
    
    // Map transforms the result
    doubled := IO.Map(func(n int) int {
        return n * 2
    })(io)
    
    result := doubled()
    fmt.Println(result) // 10
}
```

### Chain (FlatMap)

Chain IO operations:

```go
package main

import (
    "fmt"
    IO "github.com/IBM/fp-go/io"
)

func getUser(id int) IO.IO[string] {
    return func() string {
        return fmt.Sprintf("User-%d", id)
    }
}

func main() {
    io := IO.Of(42)
    
    // Chain to another IO operation
    result := IO.Chain(func(id int) IO.IO[string] {
        return getUser(id)
    })(io)
    
    user := result()
    fmt.Println(user) // User-42
}
```

### Ap (Apply)

Apply a function wrapped in IO:

```go
package main

import (
    "fmt"
    IO "github.com/IBM/fp-go/io"
)

func main() {
    // IO containing a function
    ioFunc := IO.Of(func(n int) int {
        return n * 2
    })
    
    // IO containing a value
    ioValue := IO.Of(21)
    
    // Apply the function to the value
    result := IO.Ap[int, int](ioValue)(ioFunc)
    
    fmt.Println(result()) // 42
}
```

## Side Effects

### Console I/O

```go
package main

import (
    "fmt"
    IO "github.com/IBM/fp-go/io"
)

func readLine() IO.IO[string] {
    return func() string {
        var input string
        fmt.Scanln(&input)
        return input
    }
}

func writeLine(s string) IO.IO[struct{}] {
    return func() struct{} {
        fmt.Println(s)
        return struct{}{}
    }
}

func main() {
    // Compose I/O operations
    program := IO.Chain(func(_ struct{}) IO.IO[string] {
        return readLine()
    })(writeLine("Enter your name:"))
    
    name := program()
    fmt.Printf("Hello, %s!\n", name)
}
```

### File Operations

```go
package main

import (
    "fmt"
    "os"
    IO "github.com/IBM/fp-go/io"
)

func readFile(path string) IO.IO[string] {
    return func() string {
        data, err := os.ReadFile(path)
        if err != nil {
            return ""
        }
        return string(data)
    }
}

func writeFile(path, content string) IO.IO[error] {
    return func() error {
        return os.WriteFile(path, []byte(content), 0644)
    }
}

func main() {
    // Read file
    content := readFile("input.txt")()
    fmt.Println("Content:", content)
    
    // Write file
    err := writeFile("output.txt", "Hello, World!")()
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

## Composition

### Sequential Execution

```go
package main

import (
    "fmt"
    IO "github.com/IBM/fp-go/io"
    F "github.com/IBM/fp-go/function"
)

func step1() IO.IO[int] {
    return func() int {
        fmt.Println("Step 1")
        return 1
    }
}

func step2(n int) IO.IO[int] {
    return func() int {
        fmt.Println("Step 2")
        return n + 1
    }
}

func step3(n int) IO.IO[int] {
    return func() int {
        fmt.Println("Step 3")
        return n * 2
    }
}

func main() {
    // Compose steps sequentially
    program := F.Pipe2(
        step1(),
        IO.Chain(step2),
        IO.Chain(step3),
    )
    
    result := program()
    fmt.Println("Result:", result) // 4
}
```

### Parallel Execution

```go
package main

import (
    "fmt"
    "sync"
    IO "github.com/IBM/fp-go/io"
)

func fetchUser(id int) IO.IO[string] {
    return func() string {
        return fmt.Sprintf("User-%d", id)
    }
}

func fetchPosts(userId string) IO.IO[[]string] {
    return func() []string {
        return []string{"Post1", "Post2"}
    }
}

func main() {
    var wg sync.WaitGroup
    var user string
    var posts []string
    
    wg.Add(2)
    
    // Execute in parallel
    go func() {
        defer wg.Done()
        user = fetchUser(1)()
    }()
    
    go func() {
        defer wg.Done()
        posts = fetchPosts("User-1")()
    }()
    
    wg.Wait()
    
    fmt.Println("User:", user)
    fmt.Println("Posts:", posts)
}
```

## Practical Examples

### Database Operations

```go
package main

import (
    "fmt"
    IO "github.com/IBM/fp-go/io"
)

type User struct {
    ID   int
    Name string
}

type DB struct {
    users map[int]User
}

func (db *DB) findUser(id int) IO.IO[*User] {
    return func() *User {
        if user, ok := db.users[id]; ok {
            return &user
        }
        return nil
    }
}

func (db *DB) saveUser(user User) IO.IO[error] {
    return func() error {
        db.users[user.ID] = user
        return nil
    }
}

func main() {
    db := &DB{users: make(map[int]User)}
    
    // Save user
    saveIO := db.saveUser(User{ID: 1, Name: "Alice"})
    saveIO()
    
    // Find user
    findIO := db.findUser(1)
    user := findIO()
    
    if user != nil {
        fmt.Printf("Found: %s\n", user.Name)
    }
}
```

### HTTP Requests

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    IO "github.com/IBM/fp-go/io"
)

func httpGet(url string) IO.IO[string] {
    return func() string {
        resp, err := http.Get(url)
        if err != nil {
            return ""
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return ""
        }
        
        return string(body)
    }
}

func main() {
    // Create HTTP request IO
    getIO := httpGet("https://api.example.com/data")
    
    // Execute when needed
    response := getIO()
    fmt.Println("Response:", response)
}
```

### Caching

```go
package main

import (
    "fmt"
    "sync"
    IO "github.com/IBM/fp-go/io"
)

type Cache struct {
    mu    sync.RWMutex
    data  map[string]string
}

func (c *Cache) get(key string) IO.IO[*string] {
    return func() *string {
        c.mu.RLock()
        defer c.mu.RUnlock()
        
        if val, ok := c.data[key]; ok {
            return &val
        }
        return nil
    }
}

func (c *Cache) set(key, value string) IO.IO[struct{}] {
    return func() struct{} {
        c.mu.Lock()
        defer c.mu.Unlock()
        
        c.data[key] = value
        return struct{}{}
    }
}

func main() {
    cache := &Cache{data: make(map[string]string)}
    
    // Set value
    cache.set("key1", "value1")()
    
    // Get value
    result := cache.get("key1")()
    if result != nil {
        fmt.Println("Cached:", *result)
    }
}
```

## Migration to v2

### Key Changes

1. **Error handling**:
```go
// v1: No built-in error handling
func readFileV1(path string) IO.IO[string] {
    return func() string {
        data, _ := os.ReadFile(path)
        return string(data)
    }
}

// v2: Use IOEither for errors
func readFileV2(path string) IOE.IOEither[error, string] {
    return func() E.Either[error, string] {
        data, err := os.ReadFile(path)
        if err != nil {
            return E.Left[string](err)
        }
        return E.Right[error](string(data))
    }
}
```

2. **Simplified API**:
```go
// v2 has more consistent naming and better type inference
result := IO.Map(transform)(io)
```

### Migration Example

```go
// v1 code
func processV1() IO.IO[int] {
    return F.Pipe2(
        IO.Of(5),
        IO.Map(func(n int) int { return n * 2 }),
    )
}

// v2 equivalent
func processV2() IO.IO[int] {
    return F.Pipe1(
        IO.Of(5),
        IO.Map(func(n int) int { return n * 2 }),
    )
}
```

## See Also

- [IOEither v1](./ioeither) - For IO with error handling
- [IO v2](../v2/io) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2