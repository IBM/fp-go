---
sidebar_position: 9
---

# Function Utilities (v1)

Core function composition and manipulation utilities.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [Function v2](../v2/function).

**Key differences in v2:**
- Improved type inference
- More pipe/flow variants
- Better performance
- Additional utilities
:::

## Overview

The `function` package provides utilities for:
- Function composition
- Piping data through transformations
- Currying and partial application
- Identity and constant functions

```go
import F "github.com/IBM/fp-go/function"
```

## Basic Functions

### Identity

Returns its input unchanged:

```go
package main

import (
    "fmt"
    F "github.com/IBM/fp-go/function"
)

func main() {
    // Identity returns the value as-is
    result := F.Identity(42)
    fmt.Println(result) // 42
    
    // Useful in higher-order functions
    numbers := []int{1, 2, 3}
    // Map with identity does nothing
    same := Map(F.Identity[int])(numbers)
    fmt.Println(same) // [1 2 3]
}
```

### Constant

Returns a function that always returns the same value:

```go
package main

import (
    "fmt"
    F "github.com/IBM/fp-go/function"
)

func main() {
    // Create function that always returns 42
    alwaysFortyTwo := F.Constant(42)
    
    fmt.Println(alwaysFortyTwo()) // 42
    fmt.Println(alwaysFortyTwo()) // 42
    
    // Useful for default values
    getDefault := F.Constant("default")
    fmt.Println(getDefault()) // default
}
```

## Composition

### Pipe

Chain functions left-to-right:

```go
package main

import (
    "fmt"
    "strings"
    F "github.com/IBM/fp-go/function"
)

func double(n int) int {
    return n * 2
}

func addTen(n int) int {
    return n + 10
}

func toString(n int) string {
    return fmt.Sprintf("%d", n)
}

func main() {
    // Pipe2: value -> double -> addTen
    result1 := F.Pipe2(5, double, addTen)
    fmt.Println(result1) // 20
    
    // Pipe3: value -> double -> addTen -> toString
    result2 := F.Pipe3(5, double, addTen, toString)
    fmt.Println(result2) // "20"
}
```

### Flow

Compose functions into a single function:

```go
package main

import (
    "fmt"
    "strings"
    F "github.com/IBM/fp-go/function"
)

func trim(s string) string {
    return strings.TrimSpace(s)
}

func upper(s string) string {
    return strings.ToUpper(s)
}

func addPrefix(s string) string {
    return "Hello, " + s
}

func main() {
    // Flow2: compose two functions
    process := F.Flow2(trim, upper)
    result1 := process("  alice  ")
    fmt.Println(result1) // "ALICE"
    
    // Flow3: compose three functions
    greet := F.Flow3(trim, upper, addPrefix)
    result2 := greet("  bob  ")
    fmt.Println(result2) // "Hello, BOB"
}
```

### Compose

Compose functions right-to-left (mathematical composition):

```go
package main

import (
    "fmt"
    F "github.com/IBM/fp-go/function"
)

func double(n int) int {
    return n * 2
}

func addTen(n int) int {
    return n + 10
}

func main() {
    // Compose: (addTen ∘ double)(x) = addTen(double(x))
    composed := F.Compose2(addTen, double)
    
    result := composed(5)
    fmt.Println(result) // 20 (double first: 10, then add: 20)
}
```

## Practical Examples

### Data Transformation Pipeline

```go
package main

import (
    "fmt"
    "strings"
    F "github.com/IBM/fp-go/function"
)

type User struct {
    Name  string
    Email string
}

func extractName(u User) string {
    return u.Name
}

func normalize(s string) string {
    return strings.TrimSpace(strings.ToLower(s))
}

func validate(s string) string {
    if len(s) < 2 {
        return "invalid"
    }
    return s
}

func main() {
    user := User{Name: "  ALICE  ", Email: "alice@example.com"}
    
    // Pipeline: extract -> normalize -> validate
    processName := F.Flow3(extractName, normalize, validate)
    
    result := processName(user)
    fmt.Println(result) // "alice"
}
```

### Error Handling Chain

```go
package main

import (
    "fmt"
    "strconv"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
)

func parseNumber(s string) E.Either[error, int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return E.Left[int](err)
    }
    return E.Right[error](n)
}

func validatePositive(n int) E.Either[error, int] {
    if n <= 0 {
        return E.Left[int](fmt.Errorf("must be positive"))
    }
    return E.Right[error](n)
}

func double(n int) int {
    return n * 2
}

func main() {
    // Chain with Either
    process := func(s string) E.Either[error, int] {
        return F.Pipe2(
            parseNumber(s),
            E.Chain(validatePositive),
            E.Map(double),
        )
    }
    
    result1 := process("5")
    result2 := process("-3")
    
    fmt.Println(E.IsRight(result1)) // true
    fmt.Println(E.IsLeft(result2))  // true
}
```

### Configuration Builder

```go
package main

import (
    "fmt"
    F "github.com/IBM/fp-go/function"
)

type Config struct {
    Host    string
    Port    int
    Timeout int
}

func setHost(host string) func(Config) Config {
    return func(cfg Config) Config {
        cfg.Host = host
        return cfg
    }
}

func setPort(port int) func(Config) Config {
    return func(cfg Config) Config {
        cfg.Port = port
        return cfg
    }
}

func setTimeout(timeout int) func(Config) Config {
    return func(cfg Config) Config {
        cfg.Timeout = timeout
        return cfg
    }
}

func main() {
    // Build config with pipeline
    buildConfig := F.Flow3(
        setHost("localhost"),
        setPort(8080),
        setTimeout(30),
    )
    
    config := buildConfig(Config{})
    fmt.Printf("Config: %+v\n", config)
    // Output: Config: {Host:localhost Port:8080 Timeout:30}
}
```

### Validation Pipeline

```go
package main

import (
    "fmt"
    "strings"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
)

func notEmpty(s string) E.Either[string, string] {
    if len(strings.TrimSpace(s)) == 0 {
        return E.Left[string]("cannot be empty")
    }
    return E.Right[string](s)
}

func minLength(min int) func(string) E.Either[string, string] {
    return func(s string) E.Either[string, string] {
        if len(s) < min {
            return E.Left[string](fmt.Sprintf("must be at least %d chars", min))
        }
        return E.Right[string](s)
    }
}

func maxLength(max int) func(string) E.Either[string, string] {
    return func(s string) E.Either[string, string] {
        if len(s) > max {
            return E.Left[string](fmt.Sprintf("must be at most %d chars", max))
        }
        return E.Right[string](s)
    }
}

func main() {
    // Validation pipeline
    validateUsername := func(s string) E.Either[string, string] {
        return F.Pipe3(
            notEmpty(s),
            E.Chain(minLength(3)),
            E.Chain(maxLength(20)),
        )
    }
    
    result1 := validateUsername("alice")
    result2 := validateUsername("ab")
    result3 := validateUsername("")
    
    fmt.Println(E.IsRight(result1)) // true
    fmt.Println(E.IsLeft(result2))  // true (too short)
    fmt.Println(E.IsLeft(result3))  // true (empty)
}
```

### HTTP Middleware Chain

```go
package main

import (
    "fmt"
    "net/http"
    F "github.com/IBM/fp-go/function"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func logging() Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
            next(w, r)
        }
    }
}

func auth() Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization")
            if token == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            next(w, r)
        }
    }
}

func cors() Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            next(w, r)
        }
    }
}

func main() {
    handler := func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    }
    
    // Chain middleware
    withMiddleware := F.Pipe3(
        handler,
        logging(),
        auth(),
        cors(),
    )
    
    fmt.Println("Handler with middleware created")
    _ = withMiddleware
}
```

### Data Processing Pipeline

```go
package main

import (
    "fmt"
    "strings"
    F "github.com/IBM/fp-go/function"
)

type Record struct {
    ID    int
    Name  string
    Email string
}

func filterActive(records []Record) []Record {
    // Filter logic
    return records
}

func sortByName(records []Record) []Record {
    // Sort logic
    return records
}

func extractEmails(records []Record) []string {
    emails := make([]string, len(records))
    for i, r := range records {
        emails[i] = r.Email
    }
    return emails
}

func normalize(emails []string) []string {
    normalized := make([]string, len(emails))
    for i, e := range emails {
        normalized[i] = strings.ToLower(e)
    }
    return normalized
}

func main() {
    records := []Record{
        {ID: 1, Name: "Alice", Email: "ALICE@EXAMPLE.COM"},
        {ID: 2, Name: "Bob", Email: "BOB@EXAMPLE.COM"},
    }
    
    // Process pipeline
    processRecords := F.Flow4(
        filterActive,
        sortByName,
        extractEmails,
        normalize,
    )
    
    emails := processRecords(records)
    fmt.Println(emails)
    // Output: [alice@example.com bob@example.com]
}
```

## Advanced Patterns

### Lazy Evaluation

```go
package main

import (
    "fmt"
    F "github.com/IBM/fp-go/function"
)

type Lazy[A any] func() A

func delay[A any](f func() A) Lazy[A] {
    return f
}

func force[A any](lazy Lazy[A]) A {
    return lazy()
}

func main() {
    // Expensive computation
    expensive := delay(func() int {
        fmt.Println("Computing...")
        return 42
    })
    
    fmt.Println("Lazy value created")
    
    // Force evaluation
    result := force(expensive)
    fmt.Println("Result:", result)
    
    // Output:
    // Lazy value created
    // Computing...
    // Result: 42
}
```

### Memoization

```go
package main

import (
    "fmt"
    "sync"
)

func memoize[A comparable, B any](f func(A) B) func(A) B {
    cache := make(map[A]B)
    var mu sync.RWMutex
    
    return func(a A) B {
        mu.RLock()
        if val, ok := cache[a]; ok {
            mu.RUnlock()
            return val
        }
        mu.RUnlock()
        
        mu.Lock()
        defer mu.Unlock()
        
        // Double-check after acquiring write lock
        if val, ok := cache[a]; ok {
            return val
        }
        
        result := f(a)
        cache[a] = result
        return result
    }
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
    // Memoized fibonacci
    memoFib := memoize(fibonacci)
    
    fmt.Println(memoFib(10)) // Computed
    fmt.Println(memoFib(10)) // Cached
}
```

## Migration to v2

### Key Changes

1. **More pipe variants**:
```go
// v2 has Pipe1 through Pipe9
result := F.Pipe5(
    value,
    step1,
    step2,
    step3,
    step4,
    step5,
)
```

2. **Better type inference**:
```go
// v2 has improved type inference for composition
composed := F.Flow3(f1, f2, f3)
```

### Migration Example

```go
// v1 code
func processV1(s string) string {
    return F.Pipe3(
        s,
        strings.TrimSpace,
        strings.ToUpper,
        func(s string) string { return "Hello, " + s },
    )
}

// v2 equivalent (same)
func processV2(s string) string {
    return F.Pipe3(
        s,
        strings.TrimSpace,
        strings.ToUpper,
        func(s string) string { return "Hello, " + s },
    )
}
```

## See Also

- [Pipe v2](../v2/pipe) - Latest pipe utilities
- [Flow v2](../v2/flow) - Latest flow utilities
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2